package middleware

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
)

// AuthMiddleware extracts and validates JWT token to get the user ID
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}
		// Remove "Bearer " prefix if present
		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatalf("JWT_SECRET not set in environment")
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Invalid Token Claims", http.StatusUnauthorized)
			return
		}
		log.Println(claims)

		userId, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid Token ID", http.StatusUnauthorized)
			return
		}

		// **We have the user ID in JWT claims now, but we could also fetch the user from the database**
		// email, ok := claims["email"].(string)
		// if !ok {
		//     http.Error(w, "Invalid Token Email", http.StatusUnauthorized)
		//     return
		// }

		// // Get the user ID from the email
		// var userID int
		// row := database.DB.QueryRow("SELECT id FROM users WHERE email = ?", email)
		// err = row.Scan(&userID)
		// if err != nil {
		//     http.Error(w, "Unauthorized", http.StatusUnauthorized)
		//     return
		// }

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), "userID", userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
