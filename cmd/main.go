package main

import (
	"context"
	"fmt"
	"github.com/EstebanLescano/Gym/internal/user"
	"github.com/EstebanLescano/Gym/pkg/bootstrap"
	"github.com/EstebanLescano/Gym/pkg/handler"
	"log"
	"net/http"
)

func main() {
	server := http.NewServeMux()

	db, err := bootstrap.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	logger := bootstrap.NewLogger()
	repo := user.NewRepo(db, logger)
	service := user.NewService(logger, repo)

	ctx := context.Background()

	handler.NewUserHttpServer(ctx, server, user.MakeEndpoint(ctx, service))

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080", server))
}
