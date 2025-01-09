package main

import (
	"testing"
)

func TestSetAndGetReadiness(t *testing.T) {
	manager := NewReadinessManager()

	// Set readiness for a component
	manager.SetState("db", true)

	// Check readiness for the component
	state, exists := manager.IsReady("db")
	if !exists {
		t.Errorf("Component 'db' should exist")
	}
	if state != true {
		t.Errorf("Expected 'db' readiness to be true, got %v", state)
	}

	// Set readiness to false
	manager.SetState("db", false)

	// Check updated readiness
	state, exists = manager.IsReady("db")
	if state != false {
		t.Errorf("Expected 'db' readiness to be false, got %v", state)
	}
}

func TestListComponents(t *testing.T) {
	manager := NewReadinessManager()

	// Add components
	manager.SetState("db", true)
	manager.SetState("cache", false)

	// List all components
	components := manager.ListComponents()

	// Validate the number of components
	if len(components) != 2 {
		t.Errorf("Expected 2 components, got %d", len(components))
	}

	// Validate individual components
	if components["db"] != true {
		t.Errorf("Expected 'db' readiness to be true, got %v", components["db"])
	}
	if components["cache"] != false {
		t.Errorf("Expected 'cache' readiness to be false, got %v", components["cache"])
	}
}
