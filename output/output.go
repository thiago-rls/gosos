package output

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pterm/pterm"
)

// PrintError prints an error message in a box
func PrintError(message string) {
	box := pterm.DefaultBox.WithTitle("ERROR").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Red(box))
}

// PrintWarning prints a warning message in a box
func PrintWarning(message string) {
	box := pterm.DefaultBox.WithTitle("WARNING").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Yellow(box))
}

// PrintSuccess prints a success message in a box
func PrintSuccess(message string) {
	box := pterm.DefaultBox.WithTitle("SUCCESS").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Green(box))
}

// PrintInfo prints an info message in a box
func PrintInfo(message string) {
	box := pterm.DefaultBox.WithTitle("INFO").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Cyan(box))
}

// PrintURLStatus prints the status of a URL in a box
func PrintURLStatus(url string, isUp bool) {
	status := pterm.Green("UP")
	if !isUp {
		status = pterm.Red("DOWN")
	}
	message := fmt.Sprintf("%s - %s", url, status)
	box := pterm.DefaultBox.Sprint(message)
	pterm.Println(box)
}

// PrintURLList prints a table of URLs in a box
func PrintURLList(urls []string) {
	table := pterm.TableData{
		{"Index", "URL"},
	}
	for i, url := range urls {
		table = append(table, []string{pterm.Sprint(i), url})
	}
	tableStr, err := pterm.DefaultTable.WithHasHeader().WithData(table).Srender()
	if err != nil {
		PrintError(fmt.Sprintf("Failed to render URL list: %v", err))
		return
	}
	box := pterm.DefaultBox.WithTitle("URL List").Sprint(tableStr)
	pterm.Println(box)
}

// LiveList is a thread-safe, self-contained live display of URL statuses.
// It replaces the previous package-level globals and the hand-rolled ANSI
// clear-screen sequences with a pterm.AreaPrinter, which handles portable
// in-place updates across terminals.
type LiveList struct {
	mu       sync.Mutex
	statuses []string
	area     *pterm.AreaPrinter
}

// NewLiveList initializes a live display for the provided URLs and starts
// rendering. Callers must call Stop when monitoring is done.
func NewLiveList(urls []string) (*LiveList, error) {
	area, err := pterm.DefaultArea.Start()
	if err != nil {
		return nil, err
	}

	l := &LiveList{
		statuses: make([]string, len(urls)),
		area:     area,
	}
	for i, url := range urls {
		l.statuses[i] = fmt.Sprintf("%s - Checking...", url)
	}
	l.render()
	return l, nil
}

// Update sets the status of the URL at the given index and re-renders.
func (l *LiveList) Update(index int, url string, isUp bool) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if index < 0 || index >= len(l.statuses) {
		return
	}

	status := pterm.Green("UP")
	if !isUp {
		status = pterm.Red("DOWN")
	}
	l.statuses[index] = fmt.Sprintf("%s - %s", url, status)
	l.render()
}

// Stop finalizes the live display.
func (l *LiveList) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.area != nil {
		_ = l.area.Stop()
		l.area = nil
	}
}

// render rebuilds the box contents and pushes them to the area. Callers
// must hold l.mu.
func (l *LiveList) render() {
	if l.area == nil {
		return
	}

	var content strings.Builder
	content.WriteString(pterm.Blue("Live Status:") + "\n")
	content.WriteString(strings.Repeat("-", 40) + "\n")
	for _, status := range l.statuses {
		content.WriteString(status + "\n")
	}
	content.WriteString(strings.Repeat("-", 40) + "\n")
	content.WriteString(pterm.Gray("Press Enter or Ctrl-C to stop monitoring"))

	box := pterm.DefaultBox.WithTitle("URL Monitoring").WithTitleBottomCenter().Sprint(content.String())
	l.area.Update(box)
}
