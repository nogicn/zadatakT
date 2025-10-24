package database

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"backendT/internal/database/repository"

	_ "github.com/joho/godotenv/autoload"
	_ "modernc.org/sqlite"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// get repository of rw db
	GetRepositoryRW() *repository.Queries

	// get rw db
	GetReadWriteDB() *sql.DB

	// get repository of ro db
	GetRepositoryRO() *repository.Queries

	// get ro db
	GetReadOnlyDB() *sql.DB

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	dbro   *sql.DB
	dbrw   *sql.DB
	reporo *repository.Queries
	reporw *repository.Queries
}

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *service
)

func New(dburlOverride ...string) Service {
	var doesExist bool = false

	if dburl == "" && len(dburlOverride) == 0 {
		log.Fatal("BLUEPRINT_DB_URL is not set, check your .env file")
	}

	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Only used for testing purposes with sqlite
	if len(dburlOverride) > 0 && dburlOverride[0] != "" {
		dburl = dburlOverride[0]
	} else {
		dburl = "file:" + os.Getenv("BLUEPRINT_DB_URL")
	}

	// check if file exists
	if _, err := os.Stat(dburl); err != nil {
		doesExist = true
	}

	dbro, err := sql.Open("sqlite", dburl)
	if err != nil {
		log.Fatal(err)

		// Ensure the directory exists
		dir := filepath.Dir(dburl)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("failed to create database directory: %v", err)
		}
		dbro, err = sql.Open("sqlite", dburl)
		if err != nil {
			log.Fatalf("failed to open database even after trying to create it, check free disk space: %v", err)
			return nil
		}
	}

	dbro.Exec("PRAGMA journal_mode=WAL;")
	dbro.Exec("mode=ro;")

	dbrw, err := sql.Open("sqlite", dburl)
	if err != nil {
		log.Fatal(err)

		// Ensure the directory exists
		dir := filepath.Dir(dburl)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("failed to create database directory: %v", err)
		}
		dbrw, err = sql.Open("sqlite", dburl)
		if err != nil {
			log.Fatalf("failed to open database even after trying to create it, check free disk space: %v", err)
			return nil
		}
	}
	dbrw.SetMaxOpenConns(1)
	dbrw.Exec("PRAGMA journal_mode=WAL;")
	dbrw.Exec("mode=rw;")
	dbrw.Exec("_txlock=immediate;")

	if dburl == "" {
		dburl = "./data/sqlite.db"
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		log.Fatalf("goose set dialect failed: %v", err)
	}

	// Run migrations

	if err := goose.Up(dbrw, "migrations"); err != nil {
		log.Fatalf("goose up failed: %v", err)
	}

	queriesro := repository.New(dbro)
	queriesrw := repository.New(dbrw)
	dbInstance = &service{
		dbro:   dbro,
		dbrw:   dbrw,
		reporo: queriesro,
		reporw: queriesrw,
	}

	if !doesExist {
		FillWithData(dbInstance)
	}

	return dbInstance
}

func FillWithData(s Service) {
	repo := s.GetRepositoryRW()
	ctx := context.Background()

	_, err := repo.UsersCreate(ctx, repository.UsersCreateParams{
		Username: "test",
		Email:    "test@test.com",
	})
	if err != nil {
		log.Printf("Error creating test user: %v", err)
	}

	//fmt.Println("Created test user")
	//fmt.Println(repo.UsersGetAll(ctx))

	userID, _ := repo.UsersGetByUsername(ctx, "test")
	//fmt.Printf("Fetched test user: %+v\n", userID)

	_, err = repo.PostsCreate(ctx, repository.PostsCreateParams{
		Title:   "Hello World",
		Content: "This is the first post.",
		UserID:  userID.ID,
	})
	if err != nil {
		log.Printf("Error creating test post: %v", err)
	}

}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.dbro.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.dbro.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// GetRepository returns the database repository instance
func (s *service) GetRepositoryRW() *repository.Queries {
	return s.reporw
}

func (s *service) GetRepositoryRO() *repository.Queries {
	return s.reporo
}

func (s *service) GetReadOnlyDB() *sql.DB {
	return s.dbro
}

func (s *service) GetReadWriteDB() *sql.DB {
	return s.dbrw
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	errRO := s.dbro.Close()
	errRw := s.dbrw.Close()
	if errRO != nil || errRw != nil {
		return fmt.Errorf("failed to close database connection: ro=%v rw=%v", errRO, errRw)
	}
	return nil
}
