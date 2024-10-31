package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Log request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	// Handle GET requests
	if r.Method == http.MethodGet {
		fmt.Println("GET request received")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("GET request received"))
		if err != nil {
			http.Error(w, "could not write request body", http.StatusInternalServerError)
			return
		}
	}

	// Handle POST requests
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusInternalServerError)
			return
		}
		fmt.Printf("POST request body: %s\n", string(body))

		w.WriteHeader(http.StatusOK)
	}
}

func lockHandler(w http.ResponseWriter, r *http.Request) {
	// Log request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	// Handle POST requests
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Assets locked. E2E ID: %s\n", string(body))

		w.WriteHeader(http.StatusOK)
	}
}

func unlockHandler(w http.ResponseWriter, r *http.Request) {
	// Log request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	// Handle POST requests
	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Assets unlocked. Settlement E2E ID: %s\n", string(body))

		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	// Register the handler for all incoming requests
	http.HandleFunc("/", handler)
	http.HandleFunc("/lock", lockHandler)
	http.HandleFunc("/unlock", unlockHandler)

	// Listen on port 1000
	port := ":1000"
	fmt.Printf("Server listening on port %s\n", port)
	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
