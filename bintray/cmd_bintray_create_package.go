package bintray

import (
	"fmt"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/runs"
)

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

// POST /packages/:subject/:repo
func (c *bintrayClient) packageCreate(pkg *bintrayPackage, repo string, forceCreate bool) error {
	url := fmt.Sprintf("%s/packages/%s/%s", BINTRAY_API_PREFIX, c.subject, repo)
	req := c.newRequest()
	req.JSON = pkg
	resp, err := req.Post(url)
	return c.getRespErr(resp, err, forceCreate)
}

var pkgFlags = struct {
	repoName   string
	pkgName    string
	pkgLicense string
	pkgGithub  string
}{}

func bintrayCreatePackageCommand() *cli.Command {
	return &cli.Command{
		Name:  "create-package",
		Usage: "create package in repo",
		Flags: []*cli.Flag{
			{
				Name:   "repo",
				Usage:  "bintray repository name",
				Value:  &pkgFlags.repoName,
				EnvVar: "BINTRAY_REPO",
			},
			{
				Name:   "name",
				Usage:  "bintray package name",
				Value:  &pkgFlags.pkgName,
				EnvVar: "BINTRAY_PACKAGE",
			},
			{
				Name:     "license",
				Usage:    "bintray package license",
				Value:    &pkgFlags.pkgLicense,
				DefValue: "Apache-2.0",
			},
			{
				Name:  "github-repo",
				Usage: "github repo name: (:user/:repo)",
				Value: &pkgFlags.pkgGithub,
			},
		},
		Action: func(_ *cli.Context) {
			if bintrayFlags.subject == "" {
				panic("no --subject provided")
			}
			if bintrayFlags.apikey == "" {
				panic("no --apikey provided")
			}
			if pkgFlags.repoName == "" {
				panic("no --repo provided")
			}
			if pkgFlags.pkgName == "" {
				panic("no --name provided")
			}

			pkg := &bintrayPackage{
				Name:                  pkgFlags.pkgName,
				Licenses:              []string{pkgFlags.pkgLicense},
				PublicDownloadNumbers: true,
			}

			if pkgFlags.pkgGithub != "" {
				pkg.VcsURL = fmt.Sprintf("https://github.com/%s.git", pkgFlags.pkgGithub)
				pkg.IssueTrackerURL = fmt.Sprintf("https://github.com/%s/issues", pkgFlags.pkgGithub)
				//pkg.GithubRepo = pkgFlags.pkgGithub
			}

			c := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
			err := c.packageCreate(pkg, pkgFlags.repoName, bintrayFlags.force)
			runs.PanicIfErr(err)
		},
	}
}
