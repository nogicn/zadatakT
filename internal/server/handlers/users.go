package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"baseTemplate/internal/database/repository"
)

type Repo interface {
	UsersCreate(ctx context.Context, params repository.UsersCreateParams) (repository.User, error)
	UsersGetAll(ctx context.Context) ([]repository.User, error)
}

type UsersHandler struct {
	repo Repo
}

func NewUsersHandler(r Repo) *UsersHandler { return &UsersHandler{repo: r} }

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

// @Summary Get all users
// @Description Get all users from the database
// @Tags users
// @Produce json
// @Success 200 {array} repository.User
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
