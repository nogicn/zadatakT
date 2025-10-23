package handlers

import (
	"backendT/internal/database/repository"
	logs "backendT/internal/server/handlers/logs"
	posts "backendT/internal/server/handlers/posts"
	users "backendT/internal/server/handlers/users"
	// add other handler packages here, e.g.
)

type Handlers struct {
	Users *users.UsersHandler
	Posts *posts.PostsHandler
	Logs  *logs.LogsHandler
}

func New(repo *repository.Queries) *Handlers {
	return &Handlers{
		Users: users.NewUsersHandler(repo),
		Posts: posts.NewPostsHandler(repo),
		Logs:  logs.NewLogsHandler(repo),
	}
}
