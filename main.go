package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

var (
	ready     bool       // Global variable to track readiness
	readyLock sync.Mutex // Mutex to synchronize access to readiness state
)

func main() {
	// Health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Respond with a simple health check message
		w.WriteHeader(http.StatusOK) // HTTP 200 OK
		fmt.Fprintf(w, "Server is healthy!")
	})

	// Readiness probe endpoint
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		readyLock.Lock()
		defer readyLock.Unlock()

		if ready {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Server is ready!")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Server is not ready yet.")
		}
	})

	// Endpoint to toggle readiness (for simulation/testing purposes)
	http.HandleFunc("/toggle-ready", func(w http.ResponseWriter, r *http.Request) {
		readyLock.Lock()
		defer readyLock.Unlock()

		ready = !ready // Toggle readiness state
		fmt.Fprintf(w, "Readiness state toggled. Now ready: %v", ready)
	})

	// Starting the server on port 8080
	fmt.Println("Health check server starting on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
