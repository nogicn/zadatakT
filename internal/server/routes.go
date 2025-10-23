package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/coder/websocket"

	"github.com/Treblle/treblle-go/v2"

	echoSwagger "github.com/swaggo/echo-swagger"

	"backendT/internal/database/repository"
	"backendT/internal/server/handlers"

	_ "backendT/docs"
)

// @title Your API Name
// @version 1.0
// @description Your API Description
// @host localhost:8080
// @BasePath /
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Process the request
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// Create log entry after request is processed
			entry := repository.LogsCreateParams{
				RequestID:    sql.NullString{String: c.Response().Header().Get(echo.HeaderXRequestID), Valid: true},
				RemoteIp:     sql.NullString{String: c.RealIP(), Valid: true},
				Host:         sql.NullString{String: c.Request().Host, Valid: true},
				Method:       sql.NullString{String: c.Request().Method, Valid: true},
				Uri:          sql.NullString{String: c.Request().RequestURI, Valid: true},
				UserAgent:    sql.NullString{String: c.Request().UserAgent(), Valid: true},
				Status:       sql.NullInt64{Int64: int64(c.Response().Status), Valid: true},
				Error:        sql.NullString{String: fmt.Sprintf("%v", err), Valid: true},
				Latency:      sql.NullInt64{Int64: time.Since(start).Microseconds(), Valid: true},
				LatencyHuman: sql.NullString{String: time.Since(start).String(), Valid: true},
				BytesIn:      sql.NullInt64{Int64: c.Request().ContentLength, Valid: true},
				BytesOut:     sql.NullInt64{Int64: int64(c.Response().Size), Valid: true},
			}

			// Log to console
			logLine, _ := json.Marshal(entry)
			fmt.Fprintln(os.Stdout, string(logLine))

			// Save to database
			_, dbErr := s.db.GetRepositoryRW().LogsCreate(c.Request().Context(), entry)
			if dbErr != nil {
				log.Printf("Error saving log entry: %v", dbErr)
			}

			return err
		}

	})

	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	treblle.Configure(treblle.Configuration{
		SDK_TOKEN: os.Getenv("TREBLLE_SDK_TOKEN"),
		API_KEY:   os.Getenv("TREBLLE_API_KEY"),
		Debug:     false,
	})

	// Wrap treblle's net/http middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Read and buffer the request body so both Treblle and Echo handlers can consume it.
			var bodyBytes []byte
			if c.Request().Body != nil {
				var err error
				bodyBytes, err = io.ReadAll(c.Request().Body)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]string{
						"error": "failed to read request body",
					})
				}
			}
			// restore body for Echo handlers
			c.Request().Body = io.NopCloser(bytes.NewReader(bodyBytes))

			h := treblle.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// give Treblle a copy of the request body
				if len(bodyBytes) > 0 {
					r.Body = io.NopCloser(bytes.NewReader(bodyBytes))
				}

				// Replace Echo's underlying response writer with Treblle's wrapped writer
				originalWriter := c.Response().Writer
				c.Response().Writer = w

				if err := next(c); err != nil {
					// let Echo handle the error
					c.Error(err)
				}

				// restore original writer after the handler completes
				c.Response().Writer = originalWriter
			}))

			// Use Echo's underlying response writer and a request copy for Treblle
			reqForTreblle := c.Request().Clone(c.Request().Context())
			if len(bodyBytes) > 0 {
				reqForTreblle.Body = io.NopCloser(bytes.NewReader(bodyBytes))
			}
			h.ServeHTTP(c.Response().Writer, reqForTreblle)

			return nil
		}
	})

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	e.GET("/websocket", s.websocketHandler)

	e.GET("/failure", s.simulateHorribleFailureRandomly)

	handlersRW := handlers.New(s.db.GetRepositoryRW())
	//e.GET("/users", handlersRW.Users.GetAllUsers)
	e.POST("/users", handlersRW.Users.CreateUser)
	// curl example command: curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username":"testuser","email":"test@aaaa.bbbb"}'
	e.GET("/users/id/:id", handlersRW.Users.GetUserByID)
	e.GET("/users/username/:username", handlersRW.Users.GetUserByUsername)
	e.GET("/users/email/:email", handlersRW.Users.GetUserByEmail)

	//e.GET("/posts", handlersRW.Posts.GetAllPosts)
	// curl example command: curl http://localhost:8080/posts

	e.POST("/posts", handlersRW.Posts.CreatePost)
	// curl example command: curl -X POST http://localhost:8080/posts -H "Content-Type: application/json" -d '{"title":"Test Post","content":"This is a test post."}'

	e.GET("/posts/id/:id", handlersRW.Posts.GetPostByID)
	// curl example command: curl http://localhost:8080/posts/id/1

	e.GET("/posts/userid/:userid", handlersRW.Posts.GetPostByUserID)
	// curl example command: curl http://localhost:8080/posts/userid/1

	//e.GET("/logs", handlersRW.Logs.GetAllLogs)

	// Read-only handlers for greater speed where big data is read
	handlerRO := handlers.New(s.db.GetRepositoryRO())
	e.GET("/users", handlerRO.Users.GetAllUsers)
	e.GET("/posts", handlerRO.Posts.GetAllPosts)
	e.GET("/logs", handlerRO.Logs.GetAllLogs)

	e.GET("/logs/paginated", handlerRO.Logs.GetLogsWithPagination)
	e.GET("/logs/filtered", handlerRO.Logs.GetLogsAdvanced)
	// curl example command: curl "http://localhost:8080/logs/filtered?method=GET&response=200&timeRange=-1%20hour"

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) websocketHandler(c echo.Context) error {
	w := c.Response().Writer
	r := c.Request()
	socket, err := websocket.Accept(w, r, nil)

	if err != nil {
		log.Printf("could not open websocket: %v", err)
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}

	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")

	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)

	for {
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
	return nil
}

func (s *Server) simulateHorribleFailureRandomly(c echo.Context) error {
	if rand.Int()%9 == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "🚨🚨🚨🚨🚨🚨",
		})
	} else {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "All good! ",
		})
	}
}
