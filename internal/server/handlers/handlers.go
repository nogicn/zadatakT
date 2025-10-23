package handlers

import (
	"backendT/internal/database/repository"
	posts "backendT/internal/server/handlers/posts"
	users "backendT/internal/server/handlers/users"
	// add other handler packages here, e.g.
)

type Handlers struct {
	Users *users.UsersHandler
	Posts *posts.PostsHandler
}

func New(repo *repository.Queries) *Handlers {
	return &Handlers{
		Users: users.NewUsersHandler(repo),
		Posts: posts.NewPostsHandler(repo),
	}
}
