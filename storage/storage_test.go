package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

const testFileName = ".gosos-urls-test.json"

// setupTest redirects the user's home directory to a per-test temp dir so
// that storage tests cannot read from or write to the real $HOME.
func setupTest(t *testing.T) {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	t.Setenv("USERPROFILE", tmp) // Windows
}

func TestSaveURLs(t *testing.T) {
	setupTest(t)

	// Test saving URLs to a file
	urls := &URLList{URLs: []string{"http://example.com", "http://test.com"}}
	err := SaveURLs(urls, testFileName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the content of the file
	filePath, err := getFilePath(testFileName)
	if err != nil {
		t.Fatalf("Failed to get test file path: %v", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var loadedURLs URLList
	err = json.Unmarshal(data, &loadedURLs)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(loadedURLs.URLs) != 2 || loadedURLs.URLs[0] != "http://example.com" || loadedURLs.URLs[1] != "http://test.com" {
		t.Fatalf("Expected %v, got %v", urls.URLs, loadedURLs.URLs)
	}
}

func TestLoadURLs(t *testing.T) {
	setupTest(t)

	// Test loading from a non-existent file
	urls, err := LoadURLs(testFileName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(urls.URLs) != 0 {
		t.Fatalf("Expected empty URL list, got %v", urls.URLs)
	}

	// Test loading from an existing file
	expectedURLs := &URLList{URLs: []string{"http://example.com", "http://test.com"}}
	err = SaveURLs(expectedURLs, testFileName)
	if err != nil {
		t.Fatalf("Failed to save URLs: %v", err)
	}

	urls, err = LoadURLs(testFileName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(urls.URLs) != 2 || urls.URLs[0] != "http://example.com" || urls.URLs[1] != "http://test.com" {
		t.Fatalf("Expected %v, got %v", expectedURLs.URLs, urls.URLs)
	}
}

func TestFilePath(t *testing.T) {
	setupTest(t)

	homeDir, _ := os.UserHomeDir()
	expectedPath := filepath.Join(homeDir, testFileName)

	filePath, err := getFilePath(testFileName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if filePath != expectedPath {
		t.Fatalf("Expected %v, got %v", expectedPath, filePath)
	}
}
