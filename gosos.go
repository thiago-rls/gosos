package main

import (
	"fmt"
	"os"
	"strconv"

	"git.thrls.net/thiagorls/gosos/cmd"
	"git.thrls.net/thiagorls/gosos/output"
	"git.thrls.net/thiagorls/gosos/utils"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	commandFuncs := map[string]func([]string){
		"add":    cmd.Add,
		"remove": cmd.Remove,
		"list":   func([]string) { cmd.List() },
		"run":    func([]string) { cmd.Run() },
		"live":   handleLive,
		"help":   func([]string) { printHelp() },
	}

	if fn, ok := commandFuncs[command]; ok {
		fn(args)
	} else {
		output.PrintError("Unknown command: " + command)
		printHelp()
	}
}

func handleLive(args []string) {
	interval := 30

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
