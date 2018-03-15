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
	ID      int    `json:"id"`
	Name    string `json:"name"`
	TagName string `json:"tag_name"`
	Assets  []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Size int64  `json:"size"`
	} `json:"assets"`
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

		release := githubGetRelease()

		sourceFiles := c.Args()
		for _, f := range sourceFiles {
			if fs.IsDir(f) {
				files, err := ioutil.ReadDir(f)
				runs.PanicIfErr(err)

				for _, file := range files {
					githubUploadFile(release, filepath.Join(f, file.Name()))
				}
			} else if fs.IsFile(f) {
				githubUploadFile(release, f)
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

func githubGetRelease() *githubRelease {
	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)

	fmt.Printf("getting release from tag: %s ...\n", tag)
	releaseURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", user, repo, tag)
	resp, err := req.Get(releaseURL)
	runs.PanicIfErr(err)

	//fmt.Println(resp.Text())
	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}

	release := new(githubRelease)
	err = resp.JSONUnmarshal(release)
	runs.PanicIfErr(err)

	return release
}

func githubUploadFile(release *githubRelease, filename string) {
	name := filepath.Base(filename)
	if release.exists(filename) {
		if !override {
			panic(fmt.Sprintf("asset already exists: %s", name))
		}
		githubDeleteFile(release, name)
	}

	fmt.Printf("uploading asset: %s ...\n", name)
	uploadUrl := fmt.Sprintf(
		"https://uploads.github.com/repos/%s/%s/releases/%d/assets?name=%s",
		user,
		repo,
		release.ID,
		url.QueryEscape(name),
	)
	//fmt.Println(uploadUrl)

	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)
	req.WithHeader("Content-Type", "application/octet-stream")

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

func githubDeleteFile(release *githubRelease, name string) {
	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)

	fmt.Printf("deleting exists asset: %s ...\n", name)
	deleteUrl := fmt.Sprintf("/repos/%s/%s/releases/assets/%v", user, repo, release.ID)
	resp, err := req.Delete(deleteUrl)
	runs.PanicIfErr(err)

	//fmt.Println(resp.Text())
	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}
}

func (r *githubRelease) exists(filename string) bool {
	for _, asset := range r.Assets {
		if asset.Name == filename {
			return true
		}
	}
	return false
}
