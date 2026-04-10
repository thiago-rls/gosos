package output

import (
	"fmt"
	"github.com/pterm/pterm"
	"strings"
	"sync"
)

var (
	liveStatusList []string
	statusMutex    sync.Mutex
	spinner        *pterm.SpinnerPrinter
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

func InitLiveList(urls []string) error {
	liveStatusList = make([]string, len(urls))
	for i, url := range urls {
		liveStatusList[i] = fmt.Sprintf("%s - Checking...", url)
	}

	var err error
	spinner, err = pterm.DefaultSpinner.
		WithRemoveWhenDone(false).
		WithText("Monitoring...").
		Start()

	if err != nil {
		return err
	}

	renderLiveList()
	return nil
}

// UpdateURLStatus updates the status of a URL in the live list
func UpdateURLStatus(index int, url string, isUp bool) {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	status := pterm.Green("UP")
	if !isUp {
		status = pterm.Red("DOWN")
	}
	liveStatusList[index] = fmt.Sprintf("%s - %s", url, status)
	renderLiveList()
}

// renderLiveList clears the console and re-renders the entire list in a box
func renderLiveList() {
	// Clear the console
	pterm.Print("\033[2J")
	pterm.Print("\033[H")

	// Create content for the box
	content := pterm.Blue("Live Status:") + "\n"
	content += strings.Repeat("-", 40) + "\n"

	// Add each status to the content
	for _, status := range liveStatusList {
		content += status + "\n"
	}

	content += strings.Repeat("-", 40) + "\n"
	content += pterm.Gray("Press Enter to stop monitoring")

	// Create the box with the content
	box := pterm.DefaultBox.WithTitle("URL Monitoring").WithTitleBottomCenter().Sprint(content)

	// Print the box
	fmt.Println(box)

	// Update spinner text (this keeps the spinner running)
	if spinner != nil {
		spinner.UpdateText("Monitoring...")
	}
}

// StopLiveList stops the spinner
func StopLiveList() {
	if spinner != nil {
		_ = spinner.Stop()
	}
}
