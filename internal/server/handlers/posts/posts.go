package posts

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"backendT/internal/database/repository"
)

type Repo interface {
	PostsCreate(ctx context.Context, params repository.PostsCreateParams) (repository.Post, error)
	PostsGetAll(ctx context.Context) ([]repository.Post, error)
	PostsGetByID(ctx context.Context, userID int64) (repository.Post, error)
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

// GetAllPosts handles HTTP GET requests to retrieve all posts.
// @Summary Get all posts
// @Description Returns a list of all posts from the database.
// @Tags posts
// @Produce json
// @Success 200 {array} repository.Post "List of posts"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [get]
func (h *PostsHandler) GetAllPosts(c echo.Context) error {
	posts, err := h.repo.PostsGetAll(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch posts",
		})
	}
	return c.JSON(http.StatusOK, posts)
}

// CreatePost handles HTTP POST requests to create a new post.
// @Summary Create a new post
// @Description Creates a new post in the database. Expects a JSON body with the required fields.
// @Tags posts
// @Accept json
// @Produce json
// @Param post body repository.PostsCreateParams true "New post payload"
// @Success 201 {object} repository.Post "Created post"
// @Failure 400 {object} map[string]string "Bad request - invalid payload"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts [post]
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

// GetPostByID handles HTTP GET requests to retrieve a post by their ID.
// @Summary Get post by ID
// @Description Fetches a single post by numeric ID.
// @Tags posts
// @Produce json
// @Param id path int true "Post ID"
// @Success 200 {object} repository.Post "Found post"
// @Failure 400 {object} map[string]string "Bad request - invalid ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/id/{id} [get]
func (h *PostsHandler) GetPostByID(c echo.Context) error {
	idstr := c.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID format",
		})
	}

	post, err := h.repo.PostsGetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, post)

}

// GetPostByUserID handles HTTP GET requests to retrieve a post by user ID.
// @Summary Get post by user ID
// @Description Fetches a single post by its user ID.
// @Tags posts
// @Produce json
// @Param userid path int true "User ID"
// @Success 200 {object} repository.Post "Found post"
// @Failure 400 {object} map[string]string "Bad request - invalid user ID"
// @Failure 404 {object} map[string]string "Post not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /posts/userid/{userid} [get]
func (h *PostsHandler) GetPostByUserID(c echo.Context) error {
	userID := c.Param("userid")

	userIDInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid user ID format",
		})
	}

	user, err := h.repo.PostsGetByUserID(c.Request().Context(), userIDInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch user",
		})
	}

	return c.JSON(http.StatusOK, user)

}
