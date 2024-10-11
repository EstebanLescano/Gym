package main

import (
	"context"
	"fmt"
	"github.com/EstebanLescano/Gym/internal/user"
	"github.com/EstebanLescano/Gym/pkg/bootstrap"
	"github.com/EstebanLescano/Gym/pkg/handler"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	_ = godotenv.Load() //esta es una libreria que se usa y asi como esta ya sabe que va y busca el .env donde estan las variables y las ejecuta

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

	port := os.Getenv("PORT")
	fmt.Println("Server started at port ", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), server))
}
