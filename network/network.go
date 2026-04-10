package network

import (
	"net/http"
	"time"
)

type StatusUpdate struct {
	URL  string
	IsUp bool
}

// IsUp checks if a given URL is accessible and returns true if the status code is 2xx
func IsUp(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// MonitorStatus continuously monitors the status of a URL and sends updates through a channel
func MonitorStatus(url string, stop <-chan struct{}, status chan<- StatusUpdate) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	checkAndSend(url, status)

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			checkAndSend(url, status)
		}
	}
}

// checkAndSend checks the status of a URL and sends an update through the status channel
func checkAndSend(url string, status chan<- StatusUpdate) {
	status <- StatusUpdate{URL: url, IsUp: IsUp(url)}
}
