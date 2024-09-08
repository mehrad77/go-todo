package handlers

import (
	"encoding/json"
	"go-todo/internal/models"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("your-secret-key") // Replace with a secure key

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
		log.Printf("========[USER]==>", "%+v", user)
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// TODO: save to database
	log.Printf("Registered User: %+v", user)

	// Send success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// parse the request body with NewDecoder instead of ReadAll and Unmarshal
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// validate the fields
	if loginRequest.Email == "" || loginRequest.Password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// TODO: Retrieve user from database
	user := models.User{
		Email:    loginRequest.Email,
		Password: "$2a$10$5puhOzuMoavMyZ43GpLvNe8ZPyNkwctMbRywE9vNVZ.lDOYGMru9y",
	}

	// Check if the password is correct
	if !user.CheckPassword(loginRequest.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Success response in JSON encoded map
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
