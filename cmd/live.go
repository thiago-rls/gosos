package cmd

import (
	"bufio"
	"git.thrls.net/thiagorls/gosos/network"
	"git.thrls.net/thiagorls/gosos/output"
	"os"
	"sync"
	"time"
)

const (
	updateInterval = 100 * time.Millisecond
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

	wg := launchMonitors(urlList.URLs, time.Duration(interval)*time.Second, stopChan, statusChan)

	// Listen for user input to stop the monitoring
	inputChan := listenForUserInput()

	// Create a map for efficient lookup of URL indices
	urlIndexMap := createURLIndexMap(urlList.URLs)

	// Start the main monitoring loop
	monitorLoop(urlIndexMap, statusChan, inputChan, stopChan)

	// Wait for all monitor goroutines to observe the stop signal and exit
	// before returning. We intentionally do not close statusChan — there is
	// no remaining reader, and closing it while goroutines might still send
	// would race.
	wg.Wait()
}

// initializeLiveDisplay sets up the live display for URL statuses
func initializeLiveDisplay(urls []string) error {
	if err := output.InitLiveList(urls); err != nil {
		output.PrintError("Error initializing live display: " + err.Error())
		return err
	}
	return nil
}

// launchMonitors starts a goroutine for each URL to monitor its status.
// Returns a WaitGroup that completes when every monitor has exited.
func launchMonitors(urls []string, interval time.Duration, stopChan <-chan struct{}, statusChan chan<- network.StatusUpdate) *sync.WaitGroup {
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			network.MonitorStatus(u, interval, stopChan, statusChan)
		}(url)
	}
	return &wg
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

