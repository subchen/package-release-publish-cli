package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-curl"
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
	token string

	user     string
	repo     string
	tag      string
	override bool
)

// https://developer.github.com/v3/repos/releases/#get-a-release-by-tag-name
type githubRelease struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
}

func main() {
	app := cli.NewApp()
	app.Name = "github-release-upload"
	app.Usage = "Upload asset files into github release"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"
	app.UsageText = " [OPTIONS...] file..."

	app.Flags = []*cli.Flag{
		{
			Name:   "token",
			Usage:  "GitHub access token",
			EnvVar: "GITHUB_TOKEN",
			Value:  &token,
		},
		{
			Name:   "u, user",
			Usage:  "GitHub user",
			Value:  &user,
			EnvVar: "GITHUB_REPO",
		},
		{
			Name:   "r, repo",
			Usage:  "GitHub repo",
			Value:  &repo,
			EnvVar: "GITHUB_REPO",
		},
		{
			Name:  "t, tag",
			Usage: "GitHub release tag to upload",
			Value: &tag,
		},
		{
			Name:     "override",
			Usage:    "set to true to enable overriding existing asset files",
			Value:    &override,
			DefValue: "false",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			c.ShowHelp()
			os.Exit(0)
		}

		if token == "" {
			panic("no --token provided")
		}
		if user == "" {
			panic("no --user provided")
		}
		if repo == "" {
			panic("no --repo provided")
		}
		if tag == "" {
			panic("no --tag provided")
		}

		releaseId := githubGetReleaseId()

		sourceFiles := c.Args()
		for _, f := range sourceFiles {
			if fs.IsDir(f) {
				files, err := ioutil.ReadDir(f)
				runs.PanicIfErr(err)

				for _, file := range files {
					githubUploadFile(releaseId, filepath.Join(f, file.Name()))
				}
			} else if fs.IsFile(f) {
				githubUploadFile(releaseId, f)
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

func githubGetReleaseId() string {
	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)

	// get release id from tag
	releaseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", user, repo, tag)
	resp, err := req.Get(releaseURL)
	runs.PanicIfErr(err)

	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}

	release := new(githubRelease)
	err = resp.JSONUnmarshal(release)
	runs.PanicIfErr(err)

	return release.ID
}

func githubUploadFile(releaseId string, filename string) {
	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)

	// upload file into release
	uploadUrl := fmt.Sprintf(
		"https://uploads.github.com/repos/%s/%s/releases/%s/assets?name=%s",
		user,
		repo,
		releaseId,
		url.QueryEscape(filepath.Base(filename)),
	)
	body, err := curl.NewFilePayload(filename)
	runs.PanicIfErr(err)
	resp, err := req.Post(uploadUrl, body)
	runs.PanicIfErr(err)

	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}
}
