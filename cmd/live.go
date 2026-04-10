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

	// Wait for monitors to drain before returning. statusChan is intentionally
	// never closed: no one reads it after monitorLoop returns, and closing it
	// while goroutines might still send would race.
	wg.Wait()

	// Tear down the live area before printing so the warning isn't overwritten
	// by a final redraw.
	liveList.Stop()
	output.PrintWarning(exitMessage)
}

// launchMonitors spawns one monitor goroutine per URL. The returned WaitGroup
// completes when every monitor has exited.
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

func listenForUserInput() <-chan struct{} {
	inputChan := make(chan struct{})
	go func() {
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		close(inputChan)
	}()
	return inputChan
}

func listenForInterrupt() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	return sigChan
}

func createURLIndexMap(urls []string) map[string]int {
	urlIndexMap := make(map[string]int, len(urls))
	for i, url := range urls {
		urlIndexMap[url] = i
	}
	return urlIndexMap
}

// monitorLoop forwards status updates to the live display until inputChan
// or sigChan fires, then closes stopChan and returns the exit message.
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
