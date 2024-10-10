package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	// Log request method and URL
	fmt.Printf("Received %s request for %s\n", r.Method, r.URL.Path)

	// Handle GET requests
	if r.Method == http.MethodGet {
		fmt.Println("GET request received")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GET request received"))
	}

	// Handle POST requests
	if r.Method == http.MethodPost {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Could not read request body", http.StatusInternalServerError)
			return
		}
		fmt.Printf("POST request body: %s\n", string(body))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("POST request received"))
	}
}

func main() {
	// Register the handler for all incoming requests
	http.HandleFunc("/", handler)

	// Listen on port 1000
	port := ":1000"
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

