package middleware

import (
	"net/http"
)

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Atal-Matal-Tututle", "true")
		// Allow all origins
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Specify allowed methods
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		// Specify allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "*")
		// Allow credentials
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
