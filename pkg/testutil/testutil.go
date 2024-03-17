package testutil

import (
	"context"
	"echo-test/config"
	"echo-test/model"
	"echo-test/pkg/password"
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ResetDB() {
	testName := fmt.Sprintf("%s_test", config.DbName)
	if !strings.Contains(config.DatabaseURL, testName) {
		log.Fatalf("Database URL must contain %s", testName)
	}

	migrationsDir := fmt.Sprintf("file://%s/db/migrations", config.Root)
	mt, err := migrate.New(migrationsDir, config.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to create migration: %v\n", err)
	}

	// Migrate down the database
	if err := mt.Down(); err != nil {
		log.Printf("Unable to migrate down: %v\n", err)
	}

	// Migrate up the database
	if err := mt.Up(); err != nil {
		log.Fatalf("Unable to migrate up: %v\n", err)
	}
}

func CreateUser(email, pass string) *model.User {
	db, err := pgxpool.New(context.Background(), config.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	q := model.New(db)

	hash, err := password.Hash(pass)
	if err != nil {
		log.Fatalf("Unable to hash password: %v\n", err)
	}

	user, err := q.CreateUser(context.Background(), model.CreateUserParams{
		Email:     email,
		Password:  hash,
		FirstName: "Test",
		LastName:  "User",
		Verified:  true,
		Avatar:    "https://example.com/avatar.jpg",
	})

	if err != nil {
		log.Fatalf("Unable to create user: %v\n", err)
	}

	return user
}
