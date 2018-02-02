package gobuild

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/cmd"
	"github.com/subchen/go-stack/encoding/archive"
	"github.com/subchen/go-stack/fs"
	"github.com/subchen/go-stack/runs"
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

func Command() *cli.Command {
	return &cli.Command{
		Name:  "gobuild",
		Usage: "build go sources",
		Flags: []*cli.Flag{
			{
				Name:  "name",
				Usage: "application name",
				Value: &gobuildFlags.name,
			},
			{
				Name:  "version",
				Usage: "application version",
				Value: &gobuildFlags.version,
			},
			{
				Name:     "goos",
				Usage:    "go build target os: GOOS",
				Value:    &gobuildFlags.goos,
				DefValue: "linux,darwin,windows",
			},
			{
				Name:     "goarch",
				Usage:    "go build target arch: GOARCH",
				Value:    &gobuildFlags.goarch,
				DefValue: "amd64",
			},
			{
				Name:  "archive",
				Usage: "archive format: zip or tar.gz, default is not archived",
				Value: &gobuildFlags.archive,
			},
			{
				Name:     "s,source-dir",
				Usage:    "go sources dir",
				Value:    &gobuildFlags.sourceDir,
				DefValue: ".",
			},
			{
				Name:     "o,output-dir",
				Usage:    "build target dir",
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
			if !fs.IsDir(gobuildFlags.sourceDir) {
				panic("source-dir does not exists")
			}

			if !fs.IsDir(gobuildFlags.outputDir) {
				os.MkdirAll(gobuildFlags.outputDir, 0755)
			}

			gobuild()
		},
	}
}

func gobuild() {
	buildDate := time.Now().Format(time.RFC1123Z)
	buildGitRev, err := cmd.ExecOutput("git", "rev-list", "HEAD", "--count")
	runs.PanicIfErr(err)
	buildGitCommit, err := cmd.ExecOutput("git", "describe", "--abbrev=0", "--always")
	runs.PanicIfErr(err)

	ldflags := []string{
		"-s",
		"-w",
		fmt.Sprintf("-X 'main.buildVersion=%s'", gobuildFlags.version),
		fmt.Sprintf("-X 'main.buildDate=%s'", buildDate),
		fmt.Sprintf("-X 'main.BuildGitRev=%s'", buildGitRev),
		fmt.Sprintf("-X 'main.BuildGitCommit=%s'", buildGitCommit),
	}

	for _, goos := range strings.Split(gobuildFlags.goos, ",") {
		goos = strings.TrimSpace(goos)
		for _, goarch := range strings.Split(gobuildFlags.goarch, ",") {
			goarch = strings.TrimSpace(goarch)

			filename := fmt.Sprintf("%s-%s-%s-%s", gobuildFlags.name, gobuildFlags.version, goos, goarch)
			if goos == "windows" {
				filename += ".exe"
			}
			outputFilename := filepath.Join(gobuildFlags.outputDir, filename)

			fmt.Printf("go build: %s ...\n", outputFilename)
			cmdline := fmt.Sprintf(
				`cd "%s" && GOOS=%s GOARCH=%s go build -ldflags "%s" -o "%s"`,
				gobuildFlags.sourceDir,
				goos,
				goarch,
				strings.Join(ldflags, " "),
				outputFilename,
			)
			err := cmd.Shell(cmdline)
			runs.PanicIfErr(err)

			// archive
			if gobuildFlags.archive != "" {
				archiveFilename := strings.TrimSuffix(filename, ".exe") + "." + gobuildFlags.archive
				archiveFilename = filepath.Join(gobuildFlags.outputDir, archiveFilename)
				fmt.Printf("archived: %s ...\n", archiveFilename)

				a := archive.New(archiveFilename)
				defer a.Close()

				name := gobuildFlags.name
				if goos == "windows" {
					name += ".exe"
				}
				err := a.Add(name, outputFilename)
				runs.PanicIfErr(err)

				// remove binary
				err = os.Remove(outputFilename)
				runs.PanicIfErr(err)
			}
		}
	}

	fmt.Println("Completed!")
}
