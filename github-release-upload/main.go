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
type RepositoryRelease struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	TagName   string         `json:"tag_name"`
	Assets    []ReleaseAsset `json:"assets"`
	URL       string         `json:"url"`
	AssetURL  string         `json:"assets_url"`
	UploadURL string         `json:"upload_url"`
}

type ReleaseAsset struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	URL  string `json:"url"`
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

		release := getRepositoryReleaseByTag(user, repo, tag)

		sourceFiles := c.Args()
		for _, f := range sourceFiles {
			if fs.IsDir(f) {
				files, err := ioutil.ReadDir(f)
				runs.PanicIfErr(err)

				for _, file := range files {
					release.uploadAsset(filepath.Join(f, file.Name()))
				}
			} else if fs.IsFile(f) {
				release.uploadAsset(f)
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

func getRepositoryReleaseByTag(user, repo, tag string) *RepositoryRelease {
	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)

	fmt.Printf("getting repository release from tag: %s ...\n", tag)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", user, repo, tag)
	resp, err := req.Get(url)
	runs.PanicIfErr(err)

	//fmt.Println(resp.Text())
	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}

	release := new(RepositoryRelease)
	err = resp.JSONUnmarshal(release)
	runs.PanicIfErr(err)

	return release
}

func (r *RepositoryRelease) getAsset(name string) *ReleaseAsset {
	for _, asset := range r.Assets {
		if asset.Name == name {
			return &asset
		}
	}
	return nil
}

func (r *RepositoryRelease) uploadAsset(filename string) {
	name := filepath.Base(filename)
	if asset := r.getAsset(name); asset != nil {
		if !override {
			panic(fmt.Sprintf("asset already exists: %s", name))
		}
		asset.deleteAsset()
	}

	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)

	url := strings.TrimSuffix(r.UploadURL, "{?name,label}")
	url = curl.NewURL(url, []string{"name", name})

	body, err := curl.NewFilePayload(filename)
	runs.PanicIfErr(err)

	fmt.Printf("uploading asset: %s ...\n", name)
	resp, err := req.Post(url, body)
	runs.PanicIfErr(err)

	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}
}

func (a *ReleaseAsset) deleteAsset() {
	req := curl.NewRequest(nil)
	req.WithTokenAuth(token)

	fmt.Printf("deleting exists asset: %s ...\n", a.Name)
	resp, err := req.Delete(a.URL)
	runs.PanicIfErr(err)

	//fmt.Println(resp.Text())
	if !resp.OK() {
		text, err := resp.Text()
		runs.PanicIfErr(err)
		panic(text)
	}
}
