package handlers

import (
	"encoding/json"
	"go-todo/internal/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo

	// 1. parse the body
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TODO: 2. insert into database
	log.Printf("Created Todo: %+v", todo)

	// 3. respond with the created todo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get the ID from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// TODO: 2. fetch the todo from the database
	todo := models.Todo{ID: id, UserID: 1, Title: "Sample Todo", Completed: false}

	// 3. respond with the fetched todo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get the ID from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 2. parse the request body
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	todo.ID = id    // Ensure ID is preserved
	todo.UserID = 1 // Ensure User ID is preserved

	// TODO 3. update the todo in the database
	log.Printf("Updated Todo: %+v", todo)

	// 4. respond with the updated todo
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get the ID from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// TODO 2. delete the todo from the database
	log.Printf("Deleted Todo ID: %d", id)

	// 3. respond with no content
	w.WriteHeader(http.StatusNoContent)
}
