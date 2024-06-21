package main

import (
	"errors"
	"fmt"
	"os"
)

type Command struct {
	Command     string
	Description string
	Handler     func() error
}

var commands = []Command{
	{
		Command:     "start",
		Description: "Start a timer for a project",
		Handler: func() error {
			if len(os.Args) < 3 {
				return errors.New("project name is requried")
			}
			projectName := os.Args[2]
			fmt.Println("Hello, start!", projectName)

			return nil
		},
	},
}

func HelpCmd() {
	fmt.Println("go-time-logger help command.\nThe following list shows the available commands")
	fmt.Printf("\nCommand\t\tDescription\n")
	for _, cmd := range commands {
		fmt.Printf("%s\t\t%s\n", cmd.Command, cmd.Description)
	}
}

func ProcessCommand() {
	if len(os.Args) < 2 {
		HelpCmd()
		return
	}
	cmdName := os.Args[1]

	var cmd *Command
	for _, c := range commands {
		if c.Command == cmdName {
			cmd = &c
			break
		}
	}
	if cmd == nil {
		HelpCmd()
		return
	}
	err := cmd.Handler()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running command: %s\n", err)
	}
}
