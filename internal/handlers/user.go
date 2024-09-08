package handlers

import (
	"encoding/json"
	"go-todo/internal/models"
	"io"
	"log"
	"net/http"
)

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	// init a new variable with the User struct
	var user models.User

	// parse the request body using io.ReadAll
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request", http.StatusBadRequest)
		return
	}

	// unmarshal (TIL!!) the JSON request into the user struct
	err = json.Unmarshal(body, &user)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// validation
	if user.Email == "" || user.Password == "" || user.Name == "" {
		log.Printf("%+v", user)
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// TODO: save to database
	log.Printf("Registered User: %+v", user)

	// Send success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}
