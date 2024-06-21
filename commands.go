package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/table"
)

type Command struct {
	Command     string
	Description string
	Handler     func(args []string, flags Flags) error
}

type Flags struct {
	ProjectName string
}

var commands = []Command{
	{
		Command:     "start",
		Description: "Start a timer for a project",
		Handler: func(args []string, flags Flags) error {
			if len(args) < 2 {
				return errors.New("project name is requried")
			}
			projectName := args[1]
			startTime := time.Now()
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Println("Enter q to quit")
			for scanner.Scan() {
				if strings.Compare(scanner.Text(), "q") == 0 {
					break
				}
			}
			endTime := time.Now()

			database.Add(projectName, startTime, endTime)

			return nil
		},
	},
	{
		Command:     "ls",
		Description: "list all of the time entries",
		Handler: func(args []string, flags Flags) error {

			filter := EntryFilter{}
			if flags.ProjectName != "" {
				filter.ProjectName = flags.ProjectName
			}
			entries := database.GetAll(filter)

			t := table.NewWriter()
			header := table.Row{"id", "project name", "date", "start", "end", "duration"}
			t.AppendHeader(header)

			var totalDuration time.Duration

			for _, entry := range entries {
				edate := entry.Start.Format("2006:01:02")
				estart := entry.Start.Format("15:04")
				eend := entry.End.Format("15:04")
				dur := entry.End.Sub(entry.Start).Round(time.Second)
				totalDuration += dur
				row := table.Row{entry.Id, entry.Name, edate, estart, eend, dur}
				t.AppendRow(row)
			}
			footer := table.Row{"total", "", "", "", "", totalDuration}
			t.AppendFooter(footer)
			fmt.Println(t.Render())
			return nil
		},
	},
	{
		Command:     "project",
		Description: "displays a table of all of the unique projects",
		Handler: func(args []string, flags Flags) error {
			projectCount := database.GetUniqueProjects()
			t := table.NewWriter()
			header := table.Row{"Name", "Count"}
			t.AppendHeader(header)
			for _, project := range projectCount {
				row := table.Row{project.Name, project.Count}
				t.AppendRow(row)
			}
			fmt.Println(t.Render())
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

func ProcessCommand(flags Flags) {
	args := flag.Args()
	if len(args) < 1 {
		HelpCmd()
		return
	}
	cmdName := args[0]
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
	err := cmd.Handler(args, flags)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running command: %s\n", err)
	}
}
