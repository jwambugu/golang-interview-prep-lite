package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/config"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/db"
	"github.com/matthewjamesboyle/golang-interview-prep/internal/user"
	"log"
	"net/http"
)

func main() {
	if err := config.Load(".env"); err != nil {
		log.Fatalln(err)
	}

	conn, err := db.NewConnection(config.Config.DbDSN)
	if err != nil {
		log.Fatalln(err)
	}

	defer func(conn *sql.DB) {
		_ = conn.Close()
	}(conn)

	if err = runMigrations(conn); err != nil {
		log.Fatalln(err)
	}

	svc, err := user.NewService(conn)
	if err != nil {
		log.Fatal(err)
	}

	h := user.Handler{Svc: *svc}

	http.HandleFunc("/user", h.AddUser)

	log.Println("starting http server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create postgres: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://internal/migrations", "postgres", driver)
	if err != nil {
		return fmt.Errorf("create migrate: %v", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate up: %v", err)
	}

	log.Println("Database migration complete.")
	return nil
}
