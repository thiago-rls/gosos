package cmd

import (
	"bufio"
	"git.thrls.net/thiagorls/gosos/network"
	"git.thrls.net/thiagorls/gosos/output"
	"os"
	"time"
)

const (
	updateInterval = 100 * time.Millisecond
	shutdownDelay  = time.Second
)

// Live function manages the real-time monitoring of URLs
func Live(interval int) {
	urlList, err := loadURLs()
	if err != nil {
		return
	}

	if err := initializeLiveDisplay(urlList.URLs); err != nil {
		return
	}
	defer output.StopLiveList()

	stopChan := make(chan struct{})
	statusChan := make(chan network.StatusUpdate, len(urlList.URLs))

	launchMonitors(urlList.URLs, stopChan, statusChan)

	// Listen for user input to stop the monitoring
	inputChan := listenForUserInput()

	// Create a map for efficient lookup of URL indices
	urlIndexMap := createURLIndexMap(urlList.URLs)

	// Start the main monitoring loop
	monitorLoop(urlIndexMap, statusChan, inputChan, stopChan)

	shutdown(statusChan)
}

// initializeLiveDisplay sets up the live display for URL statuses
func initializeLiveDisplay(urls []string) error {
	if err := output.InitLiveList(urls); err != nil {
		output.PrintError("Error initializing live display: " + err.Error())
		return err
	}
	return nil
}

// launchMonitors starts a goroutine for each URL to monitor its status
func launchMonitors(urls []string, stopChan <-chan struct{}, statusChan chan<- network.StatusUpdate) {
	for _, url := range urls {
		go network.MonitorStatus(url, stopChan, statusChan)
	}
}

// listenForUserInput creates a channel that closes when user input is detected
func listenForUserInput() <-chan struct{} {
	inputChan := make(chan struct{})
	go func() {
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		close(inputChan)
	}()
	return inputChan
}

// createURLIndexMap builds a map of URLs to their indices for quick lookups
func createURLIndexMap(urls []string) map[string]int {
	urlIndexMap := make(map[string]int, len(urls))
	for i, url := range urls {
		urlIndexMap[url] = i
	}
	return urlIndexMap
}

// monitorLoop handles incoming status updates and checks for user input to stop monitoring
func monitorLoop(urlIndexMap map[string]int, statusChan <-chan network.StatusUpdate, inputChan <-chan struct{}, stopChan chan<- struct{}) {
	for {
		select {
		case status := <-statusChan:
			// Update the status of a URL when a status update is received
			if index, exists := urlIndexMap[status.URL]; exists {
				output.UpdateURLStatus(index, status.URL, status.IsUp)
			}
		case <-inputChan:
			// Stop monitoring when user input is detected
			close(stopChan)
			output.PrintWarning("Monitoring stopped. Closing all connections.")
			return
		case <-time.After(updateInterval):
			// This case prevents the select from blocking indefinitely
			// It allows the loop to check for new status updates or user input regularly
		}
	}
}

// shutdown performs cleanup operations before exiting the program
func shutdown(statusChan chan network.StatusUpdate) {
	// Allow some time for final status updates to be processed
	time.Sleep(shutdownDelay)
	close(statusChan)
}
