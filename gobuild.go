package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack"
	"github.com/subchen/go-stack/fs"
	"github.com/subchen/go-stack/archive"
	"github.com/subchen/go-stack/cmd"
)

var gobuildFlags = struct {
	name      string
	version   string
	goos      string
	goarch    string
	archive   string
	sourceDir string
	outputDir string
}{}

func gobuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "gobuild",
		Usage: "build go sources",
		Flags: []*cli.Flags{
			{
				Name:  "name",
				Desc:  "application name",
				Value: &gobuildFlags.name,
			},
			{
				Name:  "version",
				Desc:  "application version",
				Value: &gobuildFlags.version,
			},
			{
				Name:     "goos",
				Desc:     "go build target os: GOOS",
				Value:    &gobuildFlags.goos,
				DefValue: "linux,darwin,windows",
			},
			{
				Name:     "goarch",
				Desc:     "go build target arch: GOARCH",
				Value:    &gobuildFlags.goarch,
				DefValue: "amd64",
			},
			{
				Name:     "archive",
				Desc:     "archive format: zip or tar.gz, default is not archived",
				Value:    &gobuildFlags.archive,
			},
			{
				Name:     "s,source-dir",
				Desc:     "go sources dir",
				Value:    &gobuildFlags.sourceDir,
				DefValue: ".",
			},
			{
				Name:     "o,output-dir",
				Desc:     "build target dir",
				Value:    &gobuildFlags.outputDir,
				DefValue: "./_releases",
			},
		},
		Action: func(c *cli.Context) {
			if gobuildFlags.name == "" {
				panic("no --name provided")
			}
			if gobuildFlags.version == "" {
				panic("no --version provided")
			}
			if !fs.DirExists(gobuildFlags.sourceDir) {
				panic("source-dir does not exists")
			}

			if !fs.DirExists(gobuildFlags.outputDir) {
				os.MkdirAll(gobuildFlags.outputDir)
			}

			gobuild()
		},
	}
}

func gobuild() {
	buildDate := time.Now().Format(time.RFC1123Z)
	buildGitRev := cmd.ExecOutput("git", "rev-list", "HEAD", "--count")
	buildGitCommit := cmd.ExecOutput("git", "describe", "--abbrev=0", "--always")

	ldflags := []string {
		"-s",
		"-w",
		fmt.Sprintf("-X 'main.buildVersion=%s'", gobuildFlags.version),
		fmt.Sprintf("-X 'main.buildDate=%s'", buildDate),
		fmt.Sprintf("-X 'main.BuildGitRev=%s'", buildGitRev),
		fmt.Sprintf("-X 'main.BuildGitCommit=%s'", buildGitCommit),
	}

	for _, goos := range strings.Split(gobuild.goos, ",") {
		goos = strings.TrimSpace(goos)
		for _, goarch := range strings.Split(gobuild.goarch, ",") {
			goarch = strings.TrimSpace(goarch)

			filename := fmt.Sprintf("%s-%s-%s-%s", gobuildFlags.name, gobuildFlags.version, goos, goarch)
			if runtime.GOOS == "windows" {
				filename += ".exe"
			}
			outputFilename := filepath.Join(gobuildFlags.outputDir, filename)

			cmdline := fmt.Sprintf(
				`cd "%s" && GOOS=%s GOARCH=%s go build -ldflags "%s" -o "%s"`,
				gobuildFlags.sourceDir,
				goos,
				goarch,
				strings.Join(ldflags, " "),
				outputFilename,
			)
			err := cmd.Shell(cmdline)
			gstack.PanicIfErr(err)
			
			// archive
			if gobuildFlags.archive != "" {
				archiveFilename := fs.BasenameWithoutExt(filename) + "." + gobuildFlags.archive
				a := archive.New(archiveFilename)
				defer a.Close()
				err := a.Add(gobuildFlags.name, outputFilename)
				gstack.PanicIfErr(err)
				
				err = os.Remove(outputFilename)
				gstack.PanicIfErr(err)
			}
		}
	}
}
