package network

import (
	"net/http"
	"time"
)

type StatusUpdate struct {
	URL  string
	IsUp bool
}

const requestTimeout = 10 * time.Second

var httpClient = &http.Client{Timeout: requestTimeout}

// IsUp checks if a given URL is accessible and returns true if the status code is 2xx.
// It issues a HEAD request first (cheap, no body transfer) and falls back to GET if
// the server doesn't support HEAD (405 Method Not Allowed or 501 Not Implemented).
func IsUp(url string) bool {
	resp, err := httpClient.Head(url)
	if err != nil {
		return false
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusNotImplemented {
		resp, err = httpClient.Get(url)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
	}

	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// MonitorStatus continuously monitors the status of a URL and sends updates through a channel
func MonitorStatus(url string, interval time.Duration, stop <-chan struct{}, status chan<- StatusUpdate) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	checkAndSend(url, status, stop)

	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			checkAndSend(url, status, stop)
		}
	}
}

// checkAndSend checks the status of a URL and sends an update through the status channel.
// If stop is closed while waiting to send, it returns without sending so callers can
// shut down cleanly even when the receiver has stopped reading.
func checkAndSend(url string, status chan<- StatusUpdate, stop <-chan struct{}) {
	update := StatusUpdate{URL: url, IsUp: IsUp(url)}
	select {
	case status <- update:
	case <-stop:
	}
}
