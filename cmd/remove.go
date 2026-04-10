package cmd

import (
	"flag"
	"fmt"
	"git.thrls.net/thiagorls/gosos/output"
	"git.thrls.net/thiagorls/gosos/storage"
	"git.thrls.net/thiagorls/gosos/utils"

	"golang.org/x/exp/slices"
)

// Remove function handles the removal of a URL from the list
func Remove(args []string) {
	url, err := parseRemoveArgs(args)
	if err != nil {
		output.PrintError(err.Error())
		return
	}

	urlList, err := loadURLs()
	if err != nil {
		output.PrintError("Error loading URLs: " + err.Error())
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

	output.PrintSuccess("URL removed from list successfully")
}

// parseRemoveArgs parses and validates the command-line arguments for the remove command
func parseRemoveArgs(args []string) (string, error) {
	rmCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	rmCmd.Parse(args)

	if rmCmd.NArg() < 1 {
		return "", fmt.Errorf("insufficient arguments\nUsage: gosos remove <url>")
	}

	return rmCmd.Arg(0), nil
}

// removeURLFromList removes the specified URL from the URLList
func removeURLFromList(urlList *storage.URLList, url string) error {
	if !slices.Contains(urlList.URLs, url) {
		return fmt.Errorf("URL does not exist in the list")
	}

	urlList.URLs = utils.RemoveElement(urlList.URLs, url)
	return nil
}
