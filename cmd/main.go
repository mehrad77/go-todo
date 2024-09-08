package main

import (
	"go-todo/internal/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/hello", HelloHandler).Methods("GET")

	// User registration route
	router.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")

	log.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is running!"))
}
