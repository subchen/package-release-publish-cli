package exec

import (
	"errors"

	"github.com/subchen/go-cli"
	"github.com/subchen/storm/pkg/sh"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:            "exec",
		Usage:           "execute command with GOPATH",
		SkipFlagParsing: true,
		Action:          runCommand,
	}
}

func runCommand(c *cli.Context) {
	args := c.Args()
	if len(args) < 1 {
		c.ShowError(errors.New("no command to run"))
	}

	sh.RunCmd(args[0], args[1:]...)
}
