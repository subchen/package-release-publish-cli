package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/encoding/sha256"
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

func main() {
	app := cli.NewApp()
	app.Name = "sha256sum-files"
	app.Usage = "Generates SHA256 (256-bit) checksum files"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"
	app.UsageText = " <dir|file> ..."

	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			c.ShowHelp()
			os.Exit(0)
		}

		for _, f := range c.Args() {
			if fs.IsDir(f) {
				files, err := ioutil.ReadDir(f)
				runs.PanicIfErr(err)

				for _, file := range files {
					sha256sum(filepath.Join(f, file.Name()))
				}
			} else if fs.IsFile(f) {
				sha256sum(f)
			} else {
				panic("file not exists: " + f)
			}
		}

		fmt.Println("sha256sum: Completed!")
	}

	if buildVersion != "" {
		app.Version = buildVersion + "-" + buildGitRev
	}
	app.BuildGitCommit = buildGitCommit
	app.BuildDate = buildDate

	app.Run(os.Args)
}

func sha256sum(file string) {
	if strings.HasSuffix(file, ".sha256") {
		return
	}

	if !fs.IsFile(file) {
		return
	}

	fmt.Printf("sha256sum: %s ...\n", file)
	err := sha256.GenerateSumFile(file)
	runs.PanicIfErr(err)
}
