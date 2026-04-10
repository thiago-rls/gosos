package cmd

import (
	"sync"

	"git.thrls.net/thiagorls/gosos/network"
	"git.thrls.net/thiagorls/gosos/output"
)

type URLStatus struct {
	URL  string
	IsUp bool
}

func Run() {
	urlList, err := loadURLs()
	if err != nil {
		output.PrintError("Error loading URLs: " + err.Error())
		return
	}

	results := checkURLs(urlList.URLs)

	printResults(results)
}

func checkURLs(urls []string) <-chan URLStatus {
	results := make(chan URLStatus, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			isUp := network.IsUp(url)
			results <- URLStatus{URL: url, IsUp: isUp}
		}(url)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func printResults(results <-chan URLStatus) {
	output.PrintInfo("Checking URLs:")
	for result := range results {
		output.PrintURLStatus(result.URL, result.IsUp)
	}
}
