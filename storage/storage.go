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

// LoadURLs reads a URLList from the user's home directory. A missing file
// yields an empty list with a nil error; other errors return a nil list.
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

func getFilePath(filename string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, filename), nil
}
