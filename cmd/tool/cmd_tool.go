package tool

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-log"
	"github.com/subchen/go-tableify"
)

var (
	oInstall bool
	oStatus  bool
)

var toolset = []struct {
	name string
	url  string
	path string
}{
	{"dep", "github.com/golang/dep", "cmd/dep"},
	{"golangci-lint", "github.com/golangci/golangci-lint", "cmd/golangci-lint"},
	{"goveralls", "github.com/mattn/goveralls", ""},
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "tool",
		Usage: "manage dependent tools",
		Flags: []*cli.Flag{
			{
				Name:  "install",
				Usage: "install tools",
				Value: &oInstall,
			},
			{
				Name:  "status",
				Usage: "list status of tools",
				Value: &oStatus,
			},
		},
		Action: runCommand,
	}
}

func runCommand(c *cli.Context) {
	if oInstall {
		gopath := "/tmp/gopath"

		defer func() {
			_ = os.RemoveAll(gopath)
		}()

		for _, t := range toolset {
			pkg := path.Join(t.url, t.path)

			if _, err := exec.LookPath(t.name); err == nil {
				log.Println(t.name, "installed")
				continue
			}

			cmd := exec.Command("go", "get", "-u", pkg)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Env = append(os.Environ(), []string{
				"GOPATH=" + gopath,
			}...)
			log.Println(t.name, "installing:", strings.Join(cmd.Args, " "))
			err := cmd.Run()
			if err != nil {
				panic(err)
			}

			src := filepath.Join(gopath, "bin", t.name)
			dest := filepath.Join("/usr/local/bin", t.name)
			err = os.Rename(src, dest)
			if err != nil {
				panic(err)
			}
		}

		return
	}

	table := tableify.New()
	table.SetHeaders("NAME", "URL", "INSTALLED")
	for _, t := range toolset {
		installed := "No"
		if _, err := exec.LookPath(t.name); err == nil {
			installed = "Yes"
		}
		table.AddRow(t.name, "https://"+t.url, installed)
	}
	table.Print()
}
