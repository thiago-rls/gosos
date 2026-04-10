package cmd

import (
	"flag"
	"fmt"
	"strconv"

	"golang.org/x/exp/slices"

	"git.thrls.net/thiagorls/gosos/output"
	"git.thrls.net/thiagorls/gosos/storage"
	"git.thrls.net/thiagorls/gosos/utils"
)

// Remove function handles the removal of a URL from the list.
// The target may be specified either as the full URL or as its index in
// `gosos list` output.
func Remove(args []string) {
	target, err := parseRemoveArgs(args)
	if err != nil {
		output.PrintError(err.Error())
		return
	}

	urlList, err := loadURLs()
	if err != nil {
		output.PrintError("Error loading URLs: " + err.Error())
		return
	}

	url, err := resolveTarget(target, urlList.URLs)
	if err != nil {
		output.PrintError(err.Error())
		return
	}

	if err := removeURLFromList(urlList, url); err != nil {
		output.PrintError(err.Error())
		return
	}

	if err := storage.SaveURLs(urlList, storage.FileName); err != nil {
		output.PrintError("Error saving URL list: " + err.Error())
		return
	}

	output.PrintSuccess("URL removed from list successfully: " + url)
}

// parseRemoveArgs parses and validates the command-line arguments for the remove command
func parseRemoveArgs(args []string) (string, error) {
	rmCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	if err := rmCmd.Parse(args); err != nil {
		return "", err
	}

	if rmCmd.NArg() < 1 {
		return "", fmt.Errorf("insufficient arguments\nUsage: gosos remove <url|index>")
	}

	return rmCmd.Arg(0), nil
}

// resolveTarget turns a user-supplied remove target into a concrete URL.
// A target that parses cleanly as a non-negative integer is treated as an
// index into urls (matching the numbering shown by `gosos list`); anything
// else is used as a literal URL.
func resolveTarget(target string, urls []string) (string, error) {
	if idx, err := strconv.Atoi(target); err == nil {
		if idx < 0 || idx >= len(urls) {
			return "", fmt.Errorf("index %d out of range (list has %d entries)", idx, len(urls))
		}
		return urls[idx], nil
	}
	return target, nil
}

// removeURLFromList removes the specified URL from the URLList
func removeURLFromList(urlList *storage.URLList, url string) error {
	if !slices.Contains(urlList.URLs, url) {
		return fmt.Errorf("URL does not exist in the list")
	}

	urlList.URLs = utils.RemoveElement(urlList.URLs, url)
	return nil
}
