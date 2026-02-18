package gotour

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFindDrones(t *testing.T) {
	ctx := context.Background()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/token":
			// Mock token response
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(TokenResponse{AccessToken: "mock-token"})	
		case "/drones/search":
		// Check Auth header
			if r.Header.Get("Authorization") != "Bearer mock-token" {
				t.Fatalf("Expected \"Bearer mock-token\", got %s", r.Header.Get("Authorization"))
			}
			// Check POST method
			if r.Method != "POST" {
				t.Fatalf("Expected POST method, got %s", r.Method)
			}
			// Mock the response from Picogrid
			serviceResponse := []Drone{
				{ID: "drone1", Lat: 37.7749, Lon: -122.4194},
				{ID: "drone2", Lat: 37.7750, Lon: -122.4195},
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(serviceResponse)
		default:
			http.NotFound(w, r)
			t.Fatalf("Unexpected endpoint: %s", r.URL.Path)
		}
	}))
	defer server.Close()	

	// Use the test server's URL for the DroneService
	service := NewDroneService("mock-apikey", 5*time.Second)
	service.BaseURL = server.URL // Point to the test server

	drones, err := service.FindDrones(ctx, 37.7749, -122.4194, 100)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(drones) != 2 {
		t.Fatalf("Expected 2 drones, got %d", len(drones))
	}
	if drones[0].ID != "drone1" || drones[1].ID != "drone2" {
		t.Errorf("Unexpected drone IDs: %s, %s", drones[0].ID, drones[1].ID)
	}
}

func TestFindDrones_Unauthorized(t *testing.T) {
	service := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "api error"}`))
	}))
	defer service.Close()

	ds := NewDroneService("invalid-key", 5*time.Second)
	ds.BaseURL = service.URL

	_, err := ds.FindDrones(context.Background(), 37.7749, -122.4194, 100)
	if err == nil {
		t.Fatal("Expected error for unauthorized access, got nil")
	}

	// Check if error message contains "api error"
	expected := "api error"
	if !contains(err.Error(), expected) {
		t.Fatalf("Expected error message to contain %s, got %s", expected, err.Error())
	}	
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestFindDrones_Caching(t *testing.T) {
    var authCount int // The counter
    
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/auth/token" {
            authCount++ // Increment every time auth is called
            w.Write([]byte(`{"access_token": "token123"}`))
            return
        }
        w.Write([]byte(`[]`))
    }))
    defer server.Close()

    service := NewDroneService("key", time.Second)
    service.BaseURL = server.URL

    // Call it twice!
    service.FindDrones(context.Background(), 0, 0, 0)
    service.FindDrones(context.Background(), 0, 0, 0)

    if authCount != 1 {
        t.Errorf("Expected 1 auth call due to caching, but got %d", authCount)
    }
}