package main

import (
	"Gym/internal/domain"
	"Gym/internal/user"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	server := http.NewServeMux()

	db := user.DB{
		Users: []domain.User{{
			ID:       1,
			Name:     "Nahuel",
			LastName: "Costamagna",
			Email:    "nahuel@test.com",
		}, {
			ID:       1,
			Name:     "Esteban",
			LastName: "Costama",
			Email:    "Esteban@test.com",
		}, {
			ID:       1,
			Name:     "Maribel",
			LastName: "Costa",
			Email:    "maribel@test.com",
		}},
		MaxUserID: 3,
	}
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	server.HandleFunc("/users", user.MakeEndpoint(ctx, service))

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
