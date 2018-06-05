package main

import (
	"os"

	"github.com/subchen/go-cli"
	"github.com/subchen/storm/cmd/build"
	"github.com/subchen/storm/cmd/exec"
	"github.com/subchen/storm/cmd/list"
	"github.com/subchen/storm/cmd/tool"
)

var buildVersion string
var buildInfo string

func main() {
	app := cli.NewApp()
	app.Name = "storm"
	app.Usage = "Fast and easily to build and release project"
	app.Version = buildVersion
	app.BuildInfo = cli.ParseBuildInfo(buildInfo)

	app.Commands = []*cli.Command{
		exec.NewCommand(),
		tool.NewCommand(),
		list.NewCommand(),
		build.NewCommand(),
	}

	app.Run(os.Args)
}
