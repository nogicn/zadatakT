package database

import (
	"context"
	"testing"

	"backendT/internal/database/repository"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseIntegration(t *testing.T) {
	// Create a new in-memory database for testing
	db := New("file:memory:?mode=memory&cache=shared")
	defer db.Close()

	repo := db.GetRepositoryRW()
	ctx := context.Background()

	t.Run("Create and Get User", func(t *testing.T) {
		// Create a test user
		testUser := repository.UsersCreateParams{
			Username: "integration_test",
			Email:    "integration@test.com",
		}

		user, err := repo.UsersCreate(ctx, testUser)
		assert.NoError(t, err)
		assert.NotZero(t, user.ID)
		assert.Equal(t, testUser.Username, user.Username)
		assert.Equal(t, testUser.Email, user.Email)

		// Retrieve the user by username
		fetchedUser, err := repo.UsersGetByUsername(ctx, testUser.Username)
		assert.NoError(t, err)
		assert.Equal(t, user.ID, fetchedUser.ID)
		assert.Equal(t, testUser.Email, fetchedUser.Email)
	})

	t.Run("Create and Get Post", func(t *testing.T) {
		// First get our test user
		user, err := repo.UsersGetByUsername(ctx, "integration_test")
		assert.NoError(t, err)

		// Create a test post
		testPost := repository.PostsCreateParams{
			Title:   "Integration Test Post",
			Content: "This is a test post for integration testing",
			UserID:  user.ID,
		}

		post, err := repo.PostsCreate(ctx, testPost)
		assert.NoError(t, err)
		assert.NotZero(t, post.ID)
		assert.Equal(t, testPost.Title, post.Title)
		assert.Equal(t, testPost.Content, post.Content)
		assert.Equal(t, testPost.UserID, post.UserID)

		// Get posts by user ID
		posts, err := repo.PostsGetByUserID(ctx, user.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(posts), 1)

		// Verify the post content
		found := false
		for _, p := range posts {
			if p.ID == post.ID {
				assert.Equal(t, testPost.Title, p.Title)
				assert.Equal(t, testPost.Content, p.Content)
				found = true
				break
			}
		}
		assert.True(t, found, "Created post should be found in user's posts")
	})

	t.Run("Get All Users", func(t *testing.T) {
		users, err := repo.UsersGetAll(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, users, "Users list should not be empty")

		// Should include our test user
		found := false
		for _, u := range users {
			if u.Email == "integration@test.com" {
				found = true
				break
			}
		}
		assert.True(t, found, "Integration test user should be in the users list")
	})
}
