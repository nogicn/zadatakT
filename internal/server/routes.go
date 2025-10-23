package server

import (
	"bytes"
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

	"backendT/internal/server/handlers"

	"github.com/Treblle/treblle-go/v2"

	echoSwagger "github.com/swaggo/echo-swagger"

	_ "backendT/docs"
)

// @title Your API Name
// @version 1.0
// @description Your API Description
// @host localhost:8080
// @BasePath /
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
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

	// Wrap treblle's net/http middleware so it can be used as an Echo middleware.
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

	e.GET("/failure", s.simulateHorribleFailure)

	usersHandler := handlers.NewUsersHandler(s.db.GetRepositoryRW())
	e.GET("/users", usersHandler.GetAllUsers)
	e.POST("/users", usersHandler.CreateUser)
	// curl example command: curl -X POST http://localhost:8080/users -H "Content-Type: application/json" -d '{"username":"testuser","email":"test@aaaa.bbbb"}'

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
func (s *Server) simulateHorribleFailure(c echo.Context) error {
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
