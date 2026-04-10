package cmd

import (
	"flag"
	"fmt"
	"git.thrls.net/thiagorls/gosos/output"
	"git.thrls.net/thiagorls/gosos/storage"
	"net/url"

	"golang.org/x/exp/slices"
)

// Add function handles the 'add' command to add a new URL to the list
func Add(args []string) {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addCmd.Parse(args)

	if err := validateArgs(addCmd); err != nil {
		output.PrintError(err.Error())
		return
	}

	// Get the URL from the first argument
	urlStr := addCmd.Arg(0)
	if err := validateURL(urlStr); err != nil {
		output.PrintError(err.Error())
		return
	}

	urlList, err := loadURLs()
	if err != nil {
		output.PrintError("Error loading URLs: " + err.Error())
		return
	}

	if slices.Contains(urlList.URLs, urlStr) {
		output.PrintWarning("URL already exists in gosos.")
		return
	}

	if err := addURLToList(urlList, urlStr); err != nil {
		output.PrintError("Error saving URL: " + err.Error())
		return
	}

	output.PrintSuccess("URL added successfully")
}

// validateArgs checks if the correct number of arguments is provided
func validateArgs(cmd *flag.FlagSet) error {
	if cmd.NArg() < 1 {
		return fmt.Errorf("insufficient arguments\nUsage: gosos add <url>")
	}
	return nil
}

// validateURL checks if the provided URL is valid
func validateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil || !isValidURL(parsedURL) {
		return fmt.Errorf("invalid URL: %s", urlStr)
	}
	return nil
}

// addURLToList appends the new URL to the list and saves it to storage
func addURLToList(urlList *storage.URLList, urlStr string) error {
	urlList.URLs = append(urlList.URLs, urlStr)
	return storage.SaveURLs(urlList, storage.FileName)
}

// isValidURL checks if a parsed URL has both a scheme and a host
func isValidURL(u *url.URL) bool {
	return u.Scheme != "" && u.Host != ""
}
