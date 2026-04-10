package cmd

import (
	"flag"
	"fmt"
	"net/url"

	"golang.org/x/exp/slices"

	"git.thrls.net/thiagorls/gosos/output"
	"git.thrls.net/thiagorls/gosos/storage"
)

func Add(args []string) {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	if err := addCmd.Parse(args); err != nil {
		output.PrintError(err.Error())
		return
	}

	if err := validateArgs(addCmd); err != nil {
		output.PrintError(err.Error())
		return
	}

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

func validateArgs(cmd *flag.FlagSet) error {
	if cmd.NArg() < 1 {
		return fmt.Errorf("insufficient arguments\nUsage: gosos add <url>")
	}
	return nil
}

func validateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil || !isValidURL(parsedURL) {
		return fmt.Errorf("invalid URL: %s", urlStr)
	}
	return nil
}

func addURLToList(urlList *storage.URLList, urlStr string) error {
	urlList.URLs = append(urlList.URLs, urlStr)
	return storage.SaveURLs(urlList, storage.FileName)
}

// isValidURL rejects schemes other than http/https so non-HTTP URLs don't
// pass validation only to fail mysteriously later in IsUp.
func isValidURL(u *url.URL) bool {
	if u.Host == "" {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}
