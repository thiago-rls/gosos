package main

import (
	"fmt"
	"git.thrls.net/thiagorls/gosos/cmd"
	"git.thrls.net/thiagorls/gosos/output"
	"git.thrls.net/thiagorls/gosos/utils"
	"os"
	"strconv"
)

func main() {
	// Check if any command-line arguments were provided
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	// Map commands to their respective handler functions
	// This makes it easy to add new commands
	commandFuncs := map[string]func([]string){
		"add":    cmd.Add,
		"remove": cmd.Remove,
		"list":   func([]string) { cmd.List() },
		"run":    func([]string) { cmd.Run() },
		"live":   handleLive,
		"help":   func([]string) { printHelp() },
	}

	// Check if the command exists in our map and execute it
	// If not, print an error and the help text
	if fn, ok := commandFuncs[command]; ok {
		fn(args)
	} else {
		output.PrintError("Unknown command: " + command)
		printHelp()
	}
}

// handleLive processes the "live" command
// It parses the interval argument if provided, otherwise uses a default value
func handleLive(args []string) {
	interval := 30 // default

	if len(args) > 0 {
		var err error
		interval, err = strconv.Atoi(args[0])

		if err != nil || interval <= 0 {
			fmt.Println("Error: interval must be a positive integer")
			return
		}
	}

	cmd.Live(interval)
}

func printHelp() {
	output.PrintInfo(utils.HelpText)
}
