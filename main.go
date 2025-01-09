package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type ReadinessManager struct {
	states map[string]bool // Tracks readiness states for components
	lock   sync.RWMutex    // Ensures thread-safe access to `ready`
}

func NewReadinessManager() *ReadinessManager {
	return &ReadinessManager{
		states: make(map[string]bool),
	}
}

// SetState safely updates the readiness state
// Pointer receiver is used because we modify the struct
func (rm *ReadinessManager) SetState(component string, value bool) {
	rm.lock.Lock()         // Lock to safely access 'ready'
	defer rm.lock.Unlock() // Unlock after checking 'ready'
	rm.states[component] = value
}

// IsReady safely retrieves the readiness state
// Pointer receiver is used because we modify the struct
func (rm *ReadinessManager) isReady(component string) (bool, bool) {
	rm.lock.RLock()
	defer rm.lock.RUnlock()
	state, exists := rm.states[component]
	return state, exists
}

// ListComponents lists all components and their readiness states
func (rm *ReadinessManager) ListComponents() map[string]bool {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	// Create a copy of the map to avoid concurrent modification issues
	copy := make(map[string]bool)
	for k, v := range rm.states {
		copy[k] = v
	}
	return copy
}

func main() {
	// Create a new ReadinessState instance and return a pointer to it
	manager := NewReadinessManager()

	// Health check handler
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Respond with a simple health check message
		w.WriteHeader(http.StatusOK) // HTTP 200 OK
		fmt.Fprintf(w, "Server is healthy!")
	})

	// Readiness probe endpoint
	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		component := r.URL.Query().Get("component")
		if component == "" {
			http.Error(w, "Missing 'component' query parameter", http.StatusBadRequest)
			return
		}

		state, exists := manager.isReady(component)
		if !exists {
			http.Error(w, fmt.Sprintf("Component '%s' not found", component), http.StatusNotFound)
			return
		}

		if state {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Component '%s' is ready!", component)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Component '%s' is not ready!", component)
		}
	})

	// Endpoint to toggle readiness (for simulation/testing purposes)
	http.HandleFunc("/set-ready", func(w http.ResponseWriter, r *http.Request) {
		component := r.URL.Query().Get("component")
		value := r.URL.Query().Get("value")

		if component == "" || value == "" {
			http.Error(w, "Missing 'component' or 'value' query parameter", http.StatusBadRequest)
			return
		}

		ready := value == "true"
		manager.SetState(component, ready)
		fmt.Fprintf(w, "Component '%s' readiness set to %v", component, ready)
	})

	// Endpoint to list all components and their readiness states
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		states := manager.ListComponents()
		for component, ready := range states {
			fmt.Fprintf(w, "Component: %s, Ready: %v\n", component, ready)
		}
	})

	// Starting the server on port 8080
	fmt.Println("Server starting on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
