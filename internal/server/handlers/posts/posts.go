package posts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"backendT/internal/database/repository"
)

type Repo interface {
	PostsCreate(ctx context.Context, params repository.PostsCreateParams) (repository.Post, error)
	PostsGetAll(ctx context.Context) ([]repository.Post, error)
	PostsGetByID(ctx context.Context, id interface{}) (repository.Post, error)
	PostsGetByUserID(ctx context.Context, userID int64) ([]repository.Post, error)
}

type PostsHandler struct {
	repo Repo
}

func NewPostsHandler(r *repository.Queries) *PostsHandler {
	return &PostsHandler{
		repo: r,
	}
}

// GetAllPosts handles HTTP GET requests to retrieve all users.
// @Summary Get all users
// @Description Returns a list of all users from the database.
// @Tags users
// @Produce json
// @Success 200 {array} repository.User "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
func (h *PostsHandler) GetAllPosts(c echo.Context) error {
	users, err := h.repo.PostsGetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch users",
		})
	}
	return c.JSON(http.StatusOK, users)
}

// CreateUser handles HTTP POST requests to create a new user.
// @Summary Create a new user
// @Description Creates a new user in the database. Expects a JSON body with the required fields.
// @Tags users
// @Accept json
// @Produce json
// @Param user body repository.UsersCreateParams true "New user payload"
// @Success 201 {object} repository.User "Created user"
// @Failure 400 {object} map[string]string "Bad request - invalid payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [post]

func (h *PostsHandler) CreatePost(c echo.Context) error {
	var newPost repository.Post
	if err := c.Bind(&newPost); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload" + err.Error(),
			"value": fmt.Sprintf("%+v", newPost),
		})
	}

	createdUser, err := h.repo.PostsCreate(c.Request().Context(), repository.PostsCreateParams{
		Title:   newPost.Title,
		Content: newPost.Content,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, createdUser)
}

// GetPostByID handles HTTP GET requests to retrieve a user by their ID.
// @Summary Get user by ID
// @Description Fetches a single user by numeric ID.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} repository.User "Found user"
// @Failure 400 {object} map[string]string "Bad request - invalid ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/{id} [get]
func (h *PostsHandler) GetPostByID(c echo.Context) error {
	var postID int64
	if err := c.Bind(&postID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload" + err.Error(),
		})
	}

	user, err := h.repo.PostsGetByID(c.Request().Context(), postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, user)

}

// GetPostByUserID handles HTTP GET requests to retrieve a user by username.
// @Summary Get user by username
// @Description Fetches a single user by their username.
// @Tags users
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} repository.User "Found user"
// @Failure 400 {object} map[string]string "Bad request - invalid username"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/username/{username} [get]
func (h *PostsHandler) GetPostByUserID(c echo.Context) error {
	var postID int64
	if err := c.Bind(&postID); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload" + err.Error(),
		})
	}

	user, err := h.repo.PostsGetByUserID(c.Request().Context(), postID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, user)

}
