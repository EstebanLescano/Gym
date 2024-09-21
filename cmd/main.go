package main

import (
	"Gym/internal/user"
	"Gym/pkg/bootstrap"

	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	server := http.NewServeMux()

	db := bootstrap.NewDB()

	logger := bootstrap.NewLogger()
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	server.HandleFunc("/users", user.MakeEndpoint(ctx, service))

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
