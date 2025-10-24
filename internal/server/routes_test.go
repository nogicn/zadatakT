package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backendT/internal/database"
	"backendT/internal/database/repository"
	"backendT/internal/server/handlers"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupTestDb() database.Service {
	dbService := database.New("file:memory:?mode=memory&cache=shared")
	database.FillWithData(dbService)
	return dbService
}

func setupPostsTestServer() (*echo.Echo, *repository.Queries) {
	e := echo.New()
	dbService := setupTestDb()
	repo := dbService.GetRepositoryRW()

	postsHandler := handlers.New(repo).Posts

	e.GET("/posts", postsHandler.GetAllPosts)
	e.POST("/posts", postsHandler.CreatePost)
	e.GET("/posts/id/:id", postsHandler.GetPostByID)
	e.GET("/posts/userid/:userid", postsHandler.GetPostByUserID)

	usersHandler := handlers.New(repo).Users

	e.GET("/users", usersHandler.GetAllUsers)
	e.GET("/users/username/:username", usersHandler.GetUserByUsername)

	return e, repo
}
func setupUsersTestServer() (*echo.Echo, *repository.Queries) {
	e := echo.New()
	dbService := setupTestDb()
	repo := dbService.GetRepositoryRW()

	usersHandler := handlers.New(repo).Users

	e.GET("/users", usersHandler.GetAllUsers)
	e.POST("/users", usersHandler.CreateUser)
	e.GET("/users/id/:id", usersHandler.GetUserByID)
	e.GET("/users/username/:username", usersHandler.GetUserByUsername)
	e.GET("/users/email/:email", usersHandler.GetUserByEmail)

	return e, repo
}

func setupLogsTestServer() (*echo.Echo, *repository.Queries) {
	e := echo.New()

	dbService := setupTestDb()
	repo := dbService.GetRepositoryRW()

	s := &Server{

		db: dbService,
	}
	e.Use(s.LoggingMiddleware())

	logsHandler := handlers.New(repo).Logs
	e.GET("/logs", logsHandler.GetAllLogs)
	e.GET("/logs/paginated", logsHandler.GetLogsWithPagination)
	e.GET("/logs/filtered", logsHandler.GetLogsAdvanced)

	userHandler := handlers.New(repo).Users

	e.GET("/users", userHandler.GetAllUsers)

	return e, repo
}

func TestPostEndpoints(t *testing.T) {
	e, _ := setupPostsTestServer()

	// Test CreatePost
	t.Run("Create Post", func(t *testing.T) {
		// First get the test user ID that was created by the database service
		req := httptest.NewRequest(http.MethodGet, "/users/username/test", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		//t.Logf("Response body: %s", rec.Body.String())

		var userResponse map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&userResponse)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		//t.Logf("User response: %+v", userResponse)

		if userResponse["id"] == nil {
			t.Fatal("User ID is nil in response")
		}

		userID := int(userResponse["id"].(float64))
		//t.Logf("Using user ID: %d", userID)

		// Now create a post for this user
		postJSON := fmt.Sprintf(`{"title":"Test Post 2","content":"This is a test post","user_id":%d}`, userID)
		req = httptest.NewRequest(http.MethodPost, "/posts", strings.NewReader(postJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err = json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Test Post 2", response["title"])
	})

	// Test GetAllPosts
	t.Run("Get All Posts", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})

	// Test GetPostByUserID
	t.Run("Get Posts by UserID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/posts/userid/1", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})
}

func TestUserEndpoints(t *testing.T) {
	e, _ := setupUsersTestServer()

	// Test CreateUser with new user
	t.Run("Create User", func(t *testing.T) {
		userJSON := `{"username":"ayoo","email":"ayoo@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "ayoo", response["username"])
	})

	// Test GetUserByUsername for existing test user
	t.Run("Get User by Username", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/username/ayoo", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "ayoo", response["username"])
	})

	// Test GetUserByEmail for the test user
	t.Run("Get User by Email", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/email/ayoo@example.com", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "ayoo@example.com", response["email"])
	})

	// Test GetAllUsers
	t.Run("Get All Users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})
}

func TestLogsEndpoints(t *testing.T) {
	e, _ := setupLogsTestServer()

	// Test Check Log Creation Middleware
	t.Run("Check Log Creation Middleware", func(t *testing.T) {
		// Create a sample log by making a request to another endpoint
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Now test the logs endpoint

		req = httptest.NewRequest(http.MethodGet, "/logs", nil)
		rec = httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response []map[string]interface{}
		err := json.NewDecoder(rec.Body).Decode(&response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})

}
