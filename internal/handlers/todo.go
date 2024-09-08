package handlers

import (
	"database/sql"
	"encoding/json"
	"go-todo/internal/database"
	"go-todo/internal/models"
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

	// 2. get the user ID from the context
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	todo.UserID = userID

	// 3. insert the todo into database
	result, err := database.DB.Exec("INSERT INTO todos (user_id, title, completed) VALUES (?, ?, ?)", todo.UserID, todo.Title, todo.Completed)
	if err != nil {
		http.Error(w, "Error inserting todo into database", http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error retrieving inserted todo ID", http.StatusInternalServerError)
		return
	}
	todo.ID = int(id)

	// 4. respond with the created todo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func GetAllTodoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get the user ID from the context
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 2. fetch all the todo from the database
	rows, err := database.DB.Query("SELECT id, title, completed FROM todos WHERE user_id = ?", userID)
	if err != nil {
		http.Error(w, "Error fetching todos from database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed); err != nil {
			http.Error(w, "Error reading todo data", http.StatusInternalServerError)
			return
		}
		todo.UserID = userID // ensure UserID is set
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error iterating over todos", http.StatusInternalServerError)
		return
	}

	// 4. respond with the fetched todo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. get the ID from the URL
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// 2. get the user ID from the context
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. fetch the todo from the database
	var todo models.Todo
	err = database.DB.QueryRow("SELECT id, user_id, title, completed FROM todos WHERE id = ? AND user_id = ?", id, userID).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Completed)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching todo from database", http.StatusInternalServerError)
		}
		return
	}

	// 4. respond with the fetched todo
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
	todo.ID = id // ensure ID is preserved

	// 3. get the user ID from the context
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	todo.UserID = userID

	// 3. update the todo in the database
	_, err = database.DB.Exec("UPDATE todos SET title = ?, completed = ? WHERE id = ? AND user_id = ?", todo.Title, todo.Completed, todo.ID, todo.UserID)
	if err != nil {
		http.Error(w, "Error updating todo in database", http.StatusInternalServerError)
		return
	}

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

	// 2. retrieve user ID from context
	userID, ok := r.Context().Value("userID").(float64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 3. delete the todo from the database
	_, err = database.DB.Exec("DELETE FROM todos WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		http.Error(w, "Error deleting todo from database", http.StatusInternalServerError)
		return
	}

	// 4. respond with no content
	w.WriteHeader(http.StatusNoContent)
}
