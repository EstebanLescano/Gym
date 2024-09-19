package main

import (
	"fmt"
	"log"
	"net/http"
)


func main() {

	http.HandleFunc("/users", UserServer)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func UserServer(w http.ResponseWriter, r *http.Request) {
	var status int
	switch r.Method {
	case http.MethodGet:
		status = 200
		w.WriteHeader(status)
		fmt.Fprintf(w,`{"status": %d, "message": %s}`, status,"success in get")

	case http.MethodPost:
		status = 200
		w.WriteHeader(status)
		fmt.Fprintf(w,`{"status": %d, "message": "%s"}`, status,"success in get")

	default:
		w.WriteHeader(http.StatusMethodNotAllowed):


	}
}