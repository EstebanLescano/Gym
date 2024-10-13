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

	h := handler.NewUserHttpServer(user.MakeEndpoint(ctx, service))

	port := os.Getenv("PORT")
	fmt.Println("Server started at port ", port)

	address := fmt.Sprintf("127.0.0.1:%s", port)

	srv := &http.Server{
		Handler: accessControl(h),
		Addr:    address,
	}
	log.Fatal(srv.ListenAndServe())
}

// esta funcion se hace para no tener problemas de cors y porder pegarle a la api desde cualquier lado
func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", " GET, POST, PATCH, OPTIONS, HEAD, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type,"+
			"DNT, If-Modified-Since, Keep-Alive, Origin, User-Agent,X-Request-With, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}
