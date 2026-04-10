package cmd

import (
	"bufio"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"git.thrls.net/thiagorls/gosos/network"
	"git.thrls.net/thiagorls/gosos/output"
)

// Live function manages the real-time monitoring of URLs
func Live(interval int) {
	urlList, err := loadURLs()
	if err != nil {
		output.PrintError("Error loading URLs: " + err.Error())
		return
	}

	liveList, err := output.NewLiveList(urlList.URLs)
	if err != nil {
		output.PrintError("Error initializing live display: " + err.Error())
		return
	}
	defer liveList.Stop()

	stopChan := make(chan struct{})
	statusChan := make(chan network.StatusUpdate, len(urlList.URLs))

	wg := launchMonitors(urlList.URLs, time.Duration(interval)*time.Second, stopChan, statusChan)

	inputChan := listenForUserInput()
	sigChan := listenForInterrupt()
	defer signal.Stop(sigChan)

	urlIndexMap := createURLIndexMap(urlList.URLs)

	exitMessage := monitorLoop(urlIndexMap, liveList, statusChan, inputChan, sigChan, stopChan)

	// Wait for all monitor goroutines to observe the stop signal and exit
	// before returning. We intentionally do not close statusChan — there is
	// no remaining reader, and closing it while goroutines might still send
	// would race.
	wg.Wait()

	// Tear down the live display before printing the exit message so the
	// warning box isn't overwritten by a final area update.
	liveList.Stop()
	output.PrintWarning(exitMessage)
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

// listenForInterrupt returns a channel that receives on SIGINT or SIGTERM so
// users can stop live monitoring with Ctrl-C in addition to pressing Enter.
func listenForInterrupt() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	return sigChan
}

// createURLIndexMap builds a map of URLs to their indices for quick lookups
func createURLIndexMap(urls []string) map[string]int {
	urlIndexMap := make(map[string]int, len(urls))
	for i, url := range urls {
		urlIndexMap[url] = i
	}
	return urlIndexMap
}

// monitorLoop handles incoming status updates and watches for stop signals
// (user input or interrupt). Returns the message to show after shutdown.
func monitorLoop(
	urlIndexMap map[string]int,
	liveList *output.LiveList,
	statusChan <-chan network.StatusUpdate,
	inputChan <-chan struct{},
	sigChan <-chan os.Signal,
	stopChan chan<- struct{},
) string {
	for {
		select {
		case status := <-statusChan:
			if index, exists := urlIndexMap[status.URL]; exists {
				liveList.Update(index, status.URL, status.IsUp)
			}
		case <-inputChan:
			close(stopChan)
			return "Monitoring stopped. Closing all connections."
		case <-sigChan:
			close(stopChan)
			return "Interrupted. Closing all connections."
		}
	}
}
