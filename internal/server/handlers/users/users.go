package users

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"backendT/internal/database/repository"
)

type Repo interface {
	UsersCreate(ctx context.Context, params repository.UsersCreateParams) (repository.User, error)
	UsersGetAll(ctx context.Context) ([]repository.User, error)
	UsersGetByID(ctx context.Context, userID int64) (repository.User, error)
	UsersGetByUsername(ctx context.Context, username string) (repository.User, error)
	UsersGetByEmail(ctx context.Context, email string) (repository.User, error)
}

type UsersHandler struct {
	repo Repo
}

func NewUsersHandler(r *repository.Queries) *UsersHandler {
	return &UsersHandler{
		repo: r,
	}
}

// GetAllUsers handles HTTP GET requests to retrieve all users.
// @Summary Get all users
// @Description Returns a list of all users from the database.
// @Tags users
// @Produce json
// @Success 200 {array} repository.User "List of users"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users [get]
func (h *UsersHandler) GetAllUsers(c echo.Context) error {
	users, err := h.repo.UsersGetAll(c.Request().Context())
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
func (h *UsersHandler) CreateUser(c echo.Context) error {
	var newUser repository.User
	if err := c.Bind(&newUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload" + err.Error(),
			"value": fmt.Sprintf("%+v", newUser),
		})
	}

	createdUser, err := h.repo.UsersCreate(c.Request().Context(), repository.UsersCreateParams{
		Username: newUser.Username,
		Email:    newUser.Email,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, createdUser)
}

// GetUserByID handles HTTP GET requests to retrieve a user by their ID.
// @Summary Get user by ID
// @Description Fetches a single user by numeric ID.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} repository.User "Found user"
// @Failure 400 {object} map[string]string "Bad request - invalid ID"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/id/{id} [get]
func (h *UsersHandler) GetUserByID(c echo.Context) error {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID format",
		})
	}

	user, err := h.repo.UsersGetByID(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, user)

}

// GetUserByUsername handles HTTP GET requests to retrieve a user by username.
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
func (h *UsersHandler) GetUserByUsername(c echo.Context) error {
	username := c.Param("username")
	if username == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Username is required",
		})
	}

	user, err := h.repo.UsersGetByUsername(c.Request().Context(), username)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// GetUserByEmail handles HTTP GET requests to retrieve a user by email.
// @Summary Get user by email
// @Description Fetches a single user by their email address.
// @Tags users
// @Produce json
// @Param email path string true "Email address"
// @Success 200 {object} repository.User "Found user"
// @Failure 400 {object} map[string]string "Bad request - invalid email"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /users/email/{email} [get]
func (h *UsersHandler) GetUserByEmail(c echo.Context) error {
	email := c.Param("email")
	if email == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Email is required",
		})
	}

	user, err := h.repo.UsersGetByEmail(c.Request().Context(), email)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "User not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, user)
}
