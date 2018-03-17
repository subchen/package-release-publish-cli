package main

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

// version
var (
	buildVersion   string
	buildGitRev    string
	buildGitCommit string
	buildDate      string
)

var (
	appName    string
	appVersion string
	goos       string
	goarch     string
	archiveFmt string
	sourceDir  string
	outputDir  string
)

func main() {
	app := cli.NewApp()
	app.Name = "go-build"
	app.Usage = "Go build sources for multiple platforms"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"
	app.UsageText = " [OPTIONS...] [source-dir]"

	app.Flags = []*cli.Flag{
		{
			Name:  "n, app-name",
			Usage: "application name",
			Value: &appName,
		},
		{
			Name:  "v, app-version",
			Usage: "application version",
			Value: &appVersion,
		},
		{
			Name:     "goos",
			Usage:    "go build target os: GOOS",
			Value:    &goos,
			DefValue: "linux,darwin,windows",
		},
		{
			Name:     "goarch",
			Usage:    "go build target arch: GOARCH",
			Value:    &goarch,
			DefValue: "amd64",
		},
		{
			Name:  "f, archive",
			Usage: "archive format: zip or tar.gz, default is not archived",
			Value: &archiveFmt,
		},
		{
			Name:     "o, output-dir",
			Usage:    "build target dir",
			Value:    &outputDir,
			DefValue: "./_releases",
		},
	}

	app.Action = func(c *cli.Context) {
		if len(os.Args) == 0 {
			c.ShowHelp()
			os.Exit(0)
		}

		if appName == "" {
			panic("no --app-name provided")
		}
		if appVersion == "" {
			panic("no --app-version provided")
		}

		sourceDir = "."
		if c.NArg() == 1 {
			sourceDir = c.Args()[0]
		}

		if !fs.IsDir(sourceDir) {
			panic("source-dir does not exists")
		}

		if !fs.IsDir(outputDir) {
			os.MkdirAll(outputDir, 0755)
		}

		gobuild()
	}

	if buildVersion != "" {
		app.Version = buildVersion + "-" + buildGitRev
	}
	app.BuildGitCommit = buildGitCommit
	app.BuildDate = buildDate

	app.Run(os.Args)
}

func gobuild() {
	buildDate := time.Now().Format(time.RFC1123Z)
	buildGitRev, err := cmd.ShellOutput(fmt.Sprintf("cd %s && git rev-list HEAD --count", sourceDir))
	runs.PanicIfErr(err)
	buildGitCommit, err := cmd.ShellOutput(fmt.Sprintf("cd %s && git describe --abbrev=0 --always", sourceDir))
	runs.PanicIfErr(err)

	ldflags := []string{
		"-s",
		"-w",
		fmt.Sprintf("-X 'main.buildVersion=%s'", appVersion),
		fmt.Sprintf("-X 'main.buildDate=%s'", buildDate),
		fmt.Sprintf("-X 'main.BuildGitRev=%s'", buildGitRev),
		fmt.Sprintf("-X 'main.BuildGitCommit=%s'", buildGitCommit),
	}

	for _, goos := range strings.Split(goos, ",") {
		goos = strings.TrimSpace(goos)
		for _, goarch := range strings.Split(goarch, ",") {
			goarch = strings.TrimSpace(goarch)

			filename := fmt.Sprintf("%s-%s-%s-%s", appName, appVersion, goos, goarch)
			if goos == "windows" {
				filename += ".exe"
			}
			outputFilename := filepath.Join(outputDir, filename)

			fmt.Printf("go build: %s ...\n", outputFilename)
			cmdline := fmt.Sprintf(
				`cd "%s" && GOOS=%s GOARCH=%s go build -ldflags "%s" -o "%s"`,
				sourceDir,
				goos,
				goarch,
				strings.Join(ldflags, " "),
				outputFilename,
			)
			err := cmd.Shell(cmdline)
			runs.PanicIfErr(err)

			// archive
			if archiveFmt != "" {
				archiveFilename := strings.TrimSuffix(filename, ".exe") + "." + archiveFmt
				archiveFilename = filepath.Join(outputDir, archiveFilename)
				fmt.Printf("archived: %s ...\n", archiveFilename)

				a := archive.New(archiveFilename)
				defer a.Close()

				entryName := appName
				if goos == "windows" {
					entryName += ".exe"
				}
				err := a.Add(entryName, outputFilename)
				runs.PanicIfErr(err)

				// remove binary
				err = os.Remove(outputFilename)
				runs.PanicIfErr(err)
			}
		}
	}

	fmt.Println("go build: Completed!")
}
