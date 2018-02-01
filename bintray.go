package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mozillazg/request"
	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack"
)

const BINTRAY_API_PREFIX = "https://api.bintray.com"

//const BINTRAY_API_PREFIX = "http://127.0.0.1"

type bintrayClient struct {
	subject string
	apikey  string
}

type bintrayRepo struct {
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Private         bool     `json:"private"`
	BusinessUnit    string   `json:"business_unit,omitempty"`
	Desc            string   `json:"desc"`
	Labels          []string `json:"labels,omitempty"`
	GpgSignMetadata bool     `json:"gpg_sign_metadata,omitempty"`
	GpgSignFiles    bool     `json:"gpg_sign_files,omitempty"`
	GpgUseOwnerKey  bool     `json:"gpg_use_owner_key,omitempty"`
}

type bintrayPackage struct {
	Name                   string   `json:"name"`
	Desc                   string   `json:"desc"`
	Labels                 []string `json:"labels,omitempty"`
	Licenses               []string `json:"licenses"`
	VcsURL                 string   `json:"vcs_url,omitempty"`
	WebsiteURL             string   `json:"website_url,omitempty"`
	IssueTrackerURL        string   `json:"issue_tracker_url,omitempty"`
	GithubRepo             string   `json:"github_repo,omitempty"`
	GithubReleaseNotesFile string   `json:"github_release_notes_file,omitempty"`
	PublicDownloadNumbers  bool     `json:"public_download_numbers"`
}

type bintrayVersion struct {
	Name                     string `json:"name"`
	Released                 string `json:"released,omitempty"`
	Desc                     string `json:"desc"`
	GithubReleaseNotesFile   string `json:"github_release_notes_file,omitempty"`
	GithubUseTagReleaseNotes bool   `json:"github_use_tag_release_notes,omitempty"`
	VcsTag                   string `json:"vcs_tag,omitempty"`
}

func newBintrayClient(subject, apikey string) *bintrayClient {
	return &bintrayClient{
		subject: subject,
		apikey:  apikey,
	}
}

func (c *bintrayClient) newRequest() *request.Request {
	req := request.NewRequest(new(http.Client))
	req.BasicAuth = request.BasicAuth{c.subject, c.apikey}
	return req
}

// POST /repos/:subject/:repo
func (c *bintrayClient) repoCreate(repo *bintrayRepo, force bool) error {
	url := fmt.Sprintf("%s/repos/%s/%s", BINTRAY_API_PREFIX, c.subject, repo.Name)
	req := c.newRequest()
	req.Json = repo
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		if resp.StatusCode == 409 && force {
			return nil
		}
		json, err := resp.Json()
		if err != nil {
			return errors.New(resp.Reason())
		}
		return errors.New(json.Get("message").MustString())
	}
	return nil
}

// POST /packages/:subject/:repo
func (c *bintrayClient) packageCreate(pkg *bintrayPackage, repo string, force bool) error {
	url := fmt.Sprintf("%s/packages/%s/%s", BINTRAY_API_PREFIX, c.subject, repo)
	req := c.newRequest()
	req.Json = pkg
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		if resp.StatusCode == 409 && force {
			return nil
		}
		json, err := resp.Json()
		if err != nil {
			return errors.New(resp.Reason())
		}
		return errors.New(json.Get("message").MustString())
	}
	return nil
}

// POST /packages/:subject/:repo/:package/versions
func (c *bintrayClient) versionCreate(version *bintrayVersion, repo string, pkg string, force bool) error {
	url := fmt.Sprintf("%s/packages/%s/%s/%s/versions", BINTRAY_API_PREFIX, c.subject, repo, pkg)
	req := c.newRequest()
	req.Json = version
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		if resp.StatusCode == 409 && force {
			return nil
		}
		json, err := resp.Json()
		if err != nil {
			return errors.New(resp.Reason())
		}
		return errors.New(json.Get("message").MustString())
	}
	return nil
}

// PUT /content/:subject/:repo/:package/:version/:file_path[?publish=0/1][?override=0/1][?explode=0/1]
func (c *bintrayClient) bintrayUpload(repo string, pkg string, version string, path string, fileContent io.Reader) error {
	url := fmt.Sprintf("%s/content/%s/%s/%s/%s", BINTRAY_API_PREFIX, c.subject, repo, pkg, version, path)
	req := c.newRequest()
	req.Headers = map[string]string{
		"X-Bintray-Publish":  "1",
		"X-Bintray-Override": "1",
		"X-Bintray-Explode":  "1",
	}
	req.Body = fileContent
	resp, err := req.Put(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		json, err := resp.Json()
		if err != nil {
			return errors.New(resp.Reason())
		}
		return errors.New(json.Get("message").MustString())
	}
	return nil
}

var bintrayFlags = struct {
	subject string
	apikey  string
	force   bool

	repoName string
	repoType string

	pkgName    string
	pkgLicense string
	pkgGithub  string
}{}

func bintrayCommand() *cli.Command {
	return &cli.Command{
		Name:  "bintray",
		Usage: "bintray cli",
		Flags: []*cli.Flag{
			{
				Name:   "subject",
				Usage:  "bintray subject",
				EnvVar: "BINTRAY_SUBJECT",
				Value:  &bintrayFlags.subject,
			},
			{
				Name:   "apikey",
				Usage:  "bintray apikey",
				EnvVar: "BINTRAY_APIKEY",
				Value:  &bintrayFlags.apikey,
			},
			{
				Name:     "force",
				Usage:    "dont error if exists",
				Value:    &bintrayFlags.force,
				DefValue: "false",
			},
		},
		Commands: []*cli.Command{
			bintrayCreateRepoCommand(),
			bintrayCreatePackageCommand(),
		},
	}
}

func bintrayCreateRepoCommand() *cli.Command {
	return &cli.Command{
		Name:  "create-repo",
		Usage: "create repository",
		Flags: []*cli.Flag{
			{
				Name:  "name",
				Usage: "bintray repository name",
				Value: &bintrayFlags.repoName,
			},
			{
				Name:     "type",
				Usage:    "bintray repository type (maven, debian, rpm, docker, npm, generic, ...)",
				Value:    &bintrayFlags.repoType,
				DefValue: "generic",
			},
		},
		Action: func(c *cli.Context) {
			if bintrayFlags.subject == "" {
				panic("no --subject provided")
			}
			if bintrayFlags.apikey == "" {
				panic("no --apikey provided")
			}
			if bintrayFlags.repoName == "" {
				panic("no --name provided")
			}

			repo := &bintrayRepo{
				Name: bintrayFlags.repoName,
				Type: bintrayFlags.repoType,
			}

			bc := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
			err := bc.repoCreate(repo, bintrayFlags.force)
			gstack.PanicIfErr(err)
		},
	}
}

func bintrayCreatePackageCommand() *cli.Command {
	return &cli.Command{
		Name:  "create-package",
		Usage: "create package in repo",
		Flags: []*cli.Flag{
			{
				Name:  "repo",
				Usage: "bintray repository name",
				Value: &bintrayFlags.repoName,
			},
			{
				Name:  "name",
				Usage: "bintray package name",
				Value: &bintrayFlags.pkgName,
			},
			{
				Name:     "license",
				Usage:    "bintray package license",
				Value:    &bintrayFlags.pkgLicense,
				DefValue: "Apache-2.0",
			},
			{
				Name:  "github-repo",
				Usage: "github repo name: (:user/:repo)",
				Value: &bintrayFlags.pkgGithub,
			},
		},
		Action: func(c *cli.Context) {
			if bintrayFlags.subject == "" {
				panic("no --subject provided")
			}
			if bintrayFlags.apikey == "" {
				panic("no --apikey provided")
			}
			if bintrayFlags.repoName == "" {
				panic("no --repo provided")
			}
			if bintrayFlags.pkgName == "" {
				panic("no --name provided")
			}

			pkg := &bintrayPackage{
				Name:                  bintrayFlags.pkgName,
				Licenses:              []string{bintrayFlags.pkgLicense},
				PublicDownloadNumbers: true,
			}

			if bintrayFlags.pkgGithub != "" {
				pkg.VcsURL = fmt.Sprintf("https://github.com/%s.git", bintrayFlags.pkgGithub)
				pkg.IssueTrackerURL = fmt.Sprintf("https://github.com/%s/issues", bintrayFlags.pkgGithub)
				//pkg.GithubRepo = bintrayFlags.pkgGithub
			}

			bc := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
			err := bc.packageCreate(pkg, bintrayFlags.repoName, bintrayFlags.force)
			gstack.PanicIfErr(err)
		},
	}
}
