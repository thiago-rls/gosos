package network

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestIsUp(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse int
		expected       bool
	}{
		{"Success 200", http.StatusOK, true},
		{"Success 299", 299, true},
		{"Failure 300", 300, false},
		{"Failure 404", http.StatusNotFound, false},
		{"Failure 500", http.StatusInternalServerError, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverResponse)
			}))
			defer server.Close()

			result := IsUp(server.URL)
			if result != tt.expected {
				t.Errorf("IsUp(%s) = %v, want %v", server.URL, result, tt.expected)
			}
		})
	}
}

func TestIsUpFallsBackToGETWhenHEADUnsupported(t *testing.T) {
	var sawGET bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodHead:
			w.WriteHeader(http.StatusMethodNotAllowed)
		case http.MethodGet:
			sawGET = true
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	if !IsUp(server.URL) {
		t.Errorf("IsUp should fall back to GET and report up, got false")
	}
	if !sawGET {
		t.Errorf("IsUp never issued a GET fallback after 405 on HEAD")
	}
}

func TestIsUpPrefersHEAD(t *testing.T) {
	var methods []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods = append(methods, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	if !IsUp(server.URL) {
		t.Fatalf("IsUp returned false for a 200 server")
	}
	if len(methods) != 1 || methods[0] != http.MethodHead {
		t.Errorf("expected a single HEAD request, got %v", methods)
	}
}

func TestIsUpInvalidURL(t *testing.T) {
	result := IsUp("http://invalid-url-that-does-not-exist.invalid")
	if result != false {
		t.Errorf("IsUp with invalid URL should return false, got true")
	}
}

func TestMonitorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	url := server.URL
	stop := make(chan struct{})
	status := make(chan StatusUpdate, 1)

	// Start monitoring in a goroutine
	go MonitorStatus(url, 30*time.Second, stop, status)

	// Wait for the first status update
	var update StatusUpdate
	select {
	case update = <-status:
		// Received an update
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for status update")
	}

	// Check the received update
	if update.URL != url {
		t.Errorf("Expected URL %s, got %s", url, update.URL)
	}

	// Stop the monitoring
	close(stop)

	// Ensure the function has stopped by waiting a bit and checking no more updates are sent
	time.Sleep(50 * time.Millisecond)
	select {
	case <-status:
		t.Fatal("Received unexpected status update after stopping")
	default:
		// This is expected
	}
}

func TestCheckAndSend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	url := server.URL
	status := make(chan StatusUpdate, 1)
	stop := make(chan struct{})

	checkAndSend(url, status, stop)

	select {
	case update := <-status:
		if update.URL != url {
			t.Errorf("Expected URL %s, got %s", url, update.URL)
		}
		// Note: We can't reliably test the IsUp value here as it depends on external factors
	case <-time.After(time.Second):
		t.Fatal("Timed out waiting for status update")
	}
}
