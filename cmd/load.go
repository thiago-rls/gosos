package cmd

import (
	"git.thrls.net/thiagorls/gosos/storage"
)

// loadURLs wraps storage.LoadURLs. Callers must treat the returned list
// as nil when err != nil and are responsible for reporting the error.
func loadURLs() (*storage.URLList, error) {
	return storage.LoadURLs(storage.FileName)
}
