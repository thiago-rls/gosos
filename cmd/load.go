package cmd

import (
	"git.thrls.net/thiagorls/gosos/storage"
)

// loadURLs retrieves the list of URLs from storage.
//
// Callers are responsible for reporting any error to the user, and must
// treat the returned *URLList as nil when err != nil — storage.LoadURLs
// may return a nil list on read or path errors. On a missing config file
// it returns an empty list with a nil error.
func loadURLs() (*storage.URLList, error) {
	return storage.LoadURLs(storage.FileName)
}
