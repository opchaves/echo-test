package main

import (
	"context"
	"fmt"
	"log"

	"echo-test/config"
	"echo-test/model"
	"echo-test/pkg/password"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.DatabaseURL)

	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	defer db.Close()
	query := model.New(db)

	tablesToDelete := []string{
		"transactions",
		"categories",
		"accounts",
		"workspaces_users",
		"workspaces",
		"tokens",
		"users",
	}

	for _, t := range tablesToDelete {
		_, err = db.Exec(ctx, fmt.Sprintf("DELETE FROM %s", t))
		if err != nil {
			log.Fatalf("Unable to delete table %v. err: %v\n", t, err)
		}
	}

	password, _ := password.Hash("password12")

	userParams := model.CreateUserParams{
		FirstName: "Paulo",
		LastName:  "Chaves",
		Email:     "paulo@example.com",
		Password:  password,
		Verified:  true,
		Avatar:    "https://example.com/avatar.jpg",
	}

	user2Params := model.CreateUserParams{
		FirstName: "Jonh",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  password,
		Verified:  true,
		Avatar:    "https://example.com/avatar.jpg",
	}

	user, err := query.CreateUser(ctx, userParams)
	if err != nil {
		log.Fatalf("Unable to create user: %v\n", err)
	}

	user2, err := query.CreateUser(ctx, user2Params)
	if err != nil {
		log.Fatalf("Unable to create user 2: %v\n", err)
	}

	log.Printf("User 1 created: %v\n", user.ID.String())
	log.Printf("User 2 created: %v\n", user2.ID.String())

	log.Println("Seeding completed successfully!")
}
