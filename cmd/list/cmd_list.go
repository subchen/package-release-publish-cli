package list

import (
	"fmt"
	"os"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/filewalk"
	"github.com/subchen/storm/pkg/config"
	"github.com/subchen/storm/pkg/sh"
)

var (
	oGopkgs      bool
	oGofiles     bool
	oGotestfiles bool
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list files",
		Flags: []*cli.Flag{
			{
				Name:  "gopkgs",
				Usage: "list go packages",
				Value: &oGopkgs,
			},
			{
				Name:  "gofiles",
				Usage: "list go files without test files",
				Value: &oGofiles,
			},
			{
				Name:  "gotestfiles",
				Usage: "list go test files",
				Value: &oGotestfiles,
			},
		},
		Action: runCommand,
	}
}

func runCommand(c *cli.Context) {
	if oGopkgs {
		listGopkgs()
	}

	if oGofiles || oGotestfiles {
		listGofiles()
	}
}

func listGopkgs() {
	sh.RunCmd("go", "list", "./...")
}

func listGofiles() {
	acceptFn := func(path string, info os.FileInfo) bool {
		name := info.Name()
		if strings.HasSuffix(name, ".go") {
			if strings.HasSuffix(name, "_test.go") {
				return oGotestfiles
			}
			return oGofiles
		}
		return false
	}

	skipDirFn := func(path string, info os.FileInfo) bool {
		name := info.Name()
		return strings.HasPrefix(name, ".") || name == "vendor"
	}

	matches, err := filewalk.FindFiles(config.ProjectRoot, acceptFn, skipDirFn, true)
	if err != nil {
		panic(err)
	}

	for _, match := range matches {
		fmt.Println(match)
	}
}
