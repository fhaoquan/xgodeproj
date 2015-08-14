package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/mpon/xgodeproj/command"
)

// GlobalFlags that uses cli flags
var GlobalFlags = []cli.Flag{}

// Commands have command list
var Commands = []cli.Command{
	{
		Name:   "list",
		Usage:  "",
		Action: command.CmdList,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound error
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}