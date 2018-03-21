package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-curl"
	"github.com/subchen/go-stack/fs"
	"github.com/subchen/go-stack/iif"
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
	subject string
	apikey  string

	targetLocation string
	publish        bool
	override       bool
	explode        bool
)

const BINTRAY_API_PREFIX = "https://api.bintray.com/"

//const BINTRAY_API_PREFIX = "http://127.0.0.1/"

func main() {
	app := cli.NewApp()
	app.Name = "bintray-upload"
	app.Usage = "Upload files into bintray repo"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"
	app.UsageText = " [OPTIONS...] target-location source-file..."

	app.Flags = []*cli.Flag{
		{
			Name:   "subject",
			Usage:  "bintray user/orig",
			EnvVar: "BINTRAY_SUBJECT",
			Value:  &subject,
		},
		{
			Name:   "apikey",
			Usage:  "bintray API key",
			EnvVar: "BINTRAY_APIKEY",
			Value:  &apikey,
		},
		{
			Name:     "publish",
			Usage:    "set to true to publish the uploaded files",
			Value:    &publish,
			DefValue: "true",
		},
		{
			Name:     "override",
			Usage:    "set to true to enable overriding existing published files",
			Value:    &override,
			DefValue: "false",
		},
		{
			Name:     "explode",
			Usage:    "set to true to explode archived files after upload",
			Value:    &explode,
			DefValue: "false",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() < 2 {
			c.ShowHelpAndExit(0)
		}

		if apikey == "" {
			panic("no --apikey provided")
		}

		targetLocation = c.Args()[0]
		if strings.Count(targetLocation, "/") < 4 {
			panic("target-location is in the form of subject/repository/package/version[/path/...]")
		}
		if subject == "" {
			subject = strings.Split(targetLocation, "/")[0]
		}

		sourceFiles := c.Args()[1:]
		for _, f := range sourceFiles {
			if fs.IsDir(f) {
				files, err := ioutil.ReadDir(f)
				runs.PanicIfErr(err)

				for _, file := range files {
					bintrayUploadFile(filepath.Join(f, file.Name()))
				}
			} else if fs.IsFile(f) {
				bintrayUploadFile(f)
			} else {
				panic("file not exists: " + f)
			}
		}

		fmt.Println("Completed!")
	}

	if buildVersion != "" {
		app.Version = buildVersion + "-" + buildGitRev
	}
	app.BuildGitCommit = buildGitCommit
	app.BuildDate = buildDate

	app.Run(os.Args)
}

func bintrayUploadFile(filename string) {
	req := curl.NewRequest(nil)
	req.WithBasicAuth(subject, apikey)
	req.Headers = map[string]string{
		"X-Bintray-Publish":  iif.String(publish, "1", "0"),
		"X-Bintray-Override": iif.String(override, "1", "0"),
		"X-Bintray-Explode":  iif.String(explode, "1", "0"),
	}

	url := BINTRAY_API_PREFIX + filepath.Join("content", targetLocation, filepath.Base(filename))
	fmt.Printf("uploading: %s ...\n", url)
	resp, err := req.Put(url, filename)
	runs.PanicIfErr(err)

	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}
}
