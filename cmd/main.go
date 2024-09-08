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

	// User routes
	router.HandleFunc("/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginUserHandler).Methods("POST")

	// To-do routes
	router.HandleFunc("/todos", handlers.CreateTodoHandler).Methods("POST")
	router.HandleFunc("/todos/{id:[0-9]+}", handlers.GetTodoHandler).Methods("GET")
	router.HandleFunc("/todos/{id:[0-9]+}", handlers.UpdateTodoHandler).Methods("PUT")
	router.HandleFunc("/todos/{id:[0-9]+}", handlers.DeleteTodoHandler).Methods("DELETE")

	log.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is running!"))
}
