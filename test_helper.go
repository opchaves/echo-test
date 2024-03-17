package main

import (
	"echo-test/config"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func TestMain(m *testing.M) {
	log.Println("TestMain url", config.DatabaseURL)
	migrationsDir := fmt.Sprintf("file://%s/db/migrations", config.Root)
	mt, err := migrate.New(migrationsDir, config.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to create migration: %v\n", err)
	}

	// Migrate up the database
	if err := mt.Up(); err != nil {
		log.Fatalf("Unable to migrate up: %v\n", err)
	}

	// Run tests
	code := m.Run()

	// Migrate down the database
	if err := mt.Down(); err != nil {
		log.Fatalf("Unable to migrate down: %v\n", err)
	}

	os.Exit(code)
}
