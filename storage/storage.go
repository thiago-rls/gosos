package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type URLList struct {
	URLs []string `json:"urls"`
}

const FileName = ".gosos-urls.json"

// LoadURLs reads the list of URLs from a file and returns it as a URLList struct
func LoadURLs(filename string) (*URLList, error) {
	filePath, err := getFilePath(filename)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filePath)
	if os.IsNotExist(err) {
		return &URLList{}, nil
	} else if err != nil {
		return nil, err
	}

	var urls URLList
	err = json.Unmarshal(data, &urls)
	return &urls, err
}

// SaveURLs writes the provided URLList struct to a file in JSON format
func SaveURLs(urls *URLList, filename string) error {
	data, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return err
	}

	filePath, err := getFilePath(filename)
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0600)
}

// getFilePath returns the full path to the file where URLs are stored
func getFilePath(filename string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, filename), nil
}
