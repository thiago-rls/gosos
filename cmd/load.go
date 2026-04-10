package cmd

import (
	"git.thrls.net/thiagorls/gosos/storage"
)

// loadURLs retrieves the list of URLs from storage. Callers are responsible
// for reporting any error to the user.
func loadURLs() (*storage.URLList, error) {
	return storage.LoadURLs(storage.FileName)
}
