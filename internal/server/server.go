package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"backendT/internal/database"
)

type Server struct {
	port int

	db database.Service
}

/*func (s *Server) GetServer() (*http.Server, database.Service) {
	return NewServer()
}*/

func NewServer(databaseNameOverride ...string) (*http.Server, database.Service) {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: database.New(databaseNameOverride...),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, NewServer.db
}
