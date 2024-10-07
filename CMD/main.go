package main

import (
	"go-todo/internal/database"
	"go-todo/internal/handlers"
	"go-todo/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Check if the environment is development
	env := os.Getenv("GO_ENV")
	log.Printf("Environment is %s", env)
	if env == "development" {
		// Load environment variables from .env file only in development
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}

	}
	// initialize database
	database.Initialize()

	// initialize router and apply middleware
	router := mux.NewRouter()
	router.Use(middleware.ErrorHandlerMiddleware)

	// Add CORS middleware to allow all origins (for development purposes)
	router.Use(middleware.CORSMiddleware)

	// test route
	router.HandleFunc("/hello", HelloHandler).Methods("GET")

	// User routes
	router.HandleFunc("/user/register", handlers.RegisterUserHandler).Methods("POST")
	router.HandleFunc("/user/login", handlers.LoginUserHandler).Methods("POST")

	// To-do routes with authentication middleware
	todoRouter := router.PathPrefix("/todos").Subrouter()
	todoRouter.Use(middleware.AuthMiddleware)
	todoRouter.HandleFunc("", handlers.CreateTodoHandler).Methods("POST")
	todoRouter.HandleFunc("", handlers.GetAllTodoHandler).Methods("GET")
	todoRouter.HandleFunc("/{id:[0-9]+}", handlers.GetTodoHandler).Methods("GET")
	todoRouter.HandleFunc("/{id:[0-9]+}", handlers.UpdateTodoHandler).Methods("PUT")
	todoRouter.HandleFunc("/{id:[0-9]+}", handlers.DeleteTodoHandler).Methods("DELETE")

	// start the server
	log.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("API is running!"))
}
