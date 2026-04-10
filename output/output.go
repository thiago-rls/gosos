package output

import (
	"fmt"
	"strings"
	"sync"

	"github.com/pterm/pterm"
)

func PrintError(message string) {
	box := pterm.DefaultBox.WithTitle("ERROR").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Red(box))
}

func PrintWarning(message string) {
	box := pterm.DefaultBox.WithTitle("WARNING").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Yellow(box))
}

func PrintSuccess(message string) {
	box := pterm.DefaultBox.WithTitle("SUCCESS").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Green(box))
}

func PrintInfo(message string) {
	box := pterm.DefaultBox.WithTitle("INFO").WithTitleBottomRight().Sprint(message)
	pterm.Println(pterm.Cyan(box))
}

func PrintURLStatus(url string, isUp bool) {
	status := pterm.Green("UP")
	if !isUp {
		status = pterm.Red("DOWN")
	}
	message := fmt.Sprintf("%s - %s", url, status)
	box := pterm.DefaultBox.Sprint(message)
	pterm.Println(box)
}

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

// LiveList is a thread-safe live display of URL statuses backed by a
// pterm.AreaPrinter for portable in-place updates.
type LiveList struct {
	mu       sync.Mutex
	statuses []string
	area     *pterm.AreaPrinter
}

// NewLiveList starts a live display for urls. Callers must call Stop.
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

func (l *LiveList) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.area != nil {
		_ = l.area.Stop()
		l.area = nil
	}
}

// render rebuilds and pushes the box contents. Callers must hold l.mu.
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
