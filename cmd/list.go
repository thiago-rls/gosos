package cmd

import (
	"git.thrls.net/thiagorls/gosos/output"
)

func List() {
	urlList, err := loadURLs()
	if err != nil {
		output.PrintError("Error loading URLs: " + err.Error())
		return
	}

	if len(urlList.URLs) == 0 {
		output.PrintInfo("No URLs found. Use 'gosos add <url>' to add URLs.")
		return
	}

	output.PrintURLList(urlList.URLs)
}
