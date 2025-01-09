package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ReadinessState struct {
	ready     bool       // Tracks readiness state
	readyLock sync.Mutex // Ensures thread-safe access to `ready`
}

// SetState safely updates the readiness state
// Pointer receiver is used because we modify the struct
func (s *ReadinessState) SetState(value bool) {
	s.readyLock.Lock()
	defer s.readyLock.Unlock()
	s.ready = value
}

// IsReady safely retrieves the readiness state
// Pointer receiver is used because we modify the struct
func (s *ReadinessState) isReady() bool {
	s.readyLock.Lock()         // Lock to safely access 'ready'
	defer s.readyLock.Unlock() // Unlock after checking 'ready'
	return s.ready
}

func main() {
	// Create a new ReadinessState instance and return a pointer to it
	state := &ReadinessState{}

	// Health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Respond with a simple health check message
		w.WriteHeader(http.StatusOK) // HTTP 200 OK
		fmt.Fprintf(w, "Server is healthy!")
	})

	// Readiness probe endpoint
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		if state.isReady() {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Service is ready!")
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Service is not ready yet.")
		}
	})

	// Endpoint to toggle readiness (for simulation/testing purposes)
	http.HandleFunc("/toggle-ready", func(w http.ResponseWriter, r *http.Request) {
		currentState := state.isReady()
		state.SetState(!currentState)
		fmt.Fprintf(w, "Readiness state toggled. Now ready: %v", !currentState)
	})

	// Starting the server on port 8080
	fmt.Println("Health check server starting on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
