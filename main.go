package main

import (
	"os"

	"github.com/subchen/go-cli"
	"github.com/subchen/publish-toolset/bintray"
	"github.com/subchen/publish-toolset/gobuild"
	"github.com/subchen/publish-toolset/sha256sum"
)

// version
var (
	buildVersion   string
	buildGitRev    string
	buildGitCommit string
	buildDate      string
)

func main() {
	app := cli.NewApp()
	app.Name = "pts"
	app.Usage = "Publish toolset for application"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"

	app.Commands = []*cli.Command{
		sha256sum.Command(),
		gobuild.Command(),
		bintray.Command(),
	}

	if buildVersion != "" {
		app.Version = buildVersion + "-" + buildGitRev
	}
	app.BuildGitCommit = buildGitCommit
	app.BuildDate = buildDate

	app.Run(os.Args)
}
