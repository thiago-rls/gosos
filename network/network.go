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

// IsUp reports whether url responds with a 2xx status. It issues HEAD first
// and falls back to GET only when the server refuses HEAD (405 or 501).
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

// checkAndSend probes url and sends an update, or abandons the send if stop
// is closed first so callers can shut down even when no one is reading.
func checkAndSend(url string, status chan<- StatusUpdate, stop <-chan struct{}) {
	update := StatusUpdate{URL: url, IsUp: IsUp(url)}
	select {
	case status <- update:
	case <-stop:
	}
}
