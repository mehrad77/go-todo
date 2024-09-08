package handlers

import (
	"database/sql"
	"encoding/json"
	"go-todo/internal/database"
	"go-todo/internal/models"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
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

	// validations
	if user.Email == "" || user.Password == "" || user.Name == "" {
		log.Printf("========[USER]==>", "%+v", user)
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}
	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// validate if email already exists in DB
	var existingUserID int
	err = database.DB.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&existingUserID)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "Error checking email existence", http.StatusInternalServerError)
		return
	}
	if existingUserID != 0 {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	// hash the password
	err = user.HashUserPassword()
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// insert the user into the database
	_, err = database.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Error inserting user into database", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully, now login.",
		"user": map[string]string{
			"email": user.Email,
			"name":  user.Name,
		},
	})
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

	// read the user from database
	var user models.User
	err := database.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", loginRequest.Email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		} else {
			http.Error(w, "Error retrieving user from database", http.StatusInternalServerError)
		}
		return
	}

	// Check if the password is correct
	if !user.CheckPassword(loginRequest.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatalf("JWT_SECRET not set in environment")
	}

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	// Success response in JSON encoded map
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"email": user.Email, "token": tokenString})
}

// ⇩⇩⇩⇩⇩ Utility Functions ⇩⇩⇩⇩⇩

func isValidEmail(email string) bool {
	// Simple regex for email validation
	return regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email)
}
