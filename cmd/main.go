package main

import (
	"log"
	"net/http"
)

func main() {
	// Start the server
	log.Println("Starting the server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
