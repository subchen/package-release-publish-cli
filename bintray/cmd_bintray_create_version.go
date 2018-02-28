package bintray

import (
	"fmt"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/runs"
)

type bintrayVersion struct {
	Name                     string `json:"name"`
	Released                 string `json:"released,omitempty"`
	Desc                     string `json:"desc"`
	GithubReleaseNotesFile   string `json:"github_release_notes_file,omitempty"`
	GithubUseTagReleaseNotes bool   `json:"github_use_tag_release_notes,omitempty"`
	VcsTag                   string `json:"vcs_tag,omitempty"`
}

// POST /packages/:subject/:repo/:package/versions
func (c *bintrayClient) versionCreate(version *bintrayVersion, repo string, pkg string, forceCreate bool) error {
	url := fmt.Sprintf("%s/packages/%s/%s/%s/versions", BINTRAY_API_PREFIX, c.subject, repo, pkg)
	req := c.newRequest()
	req.JSON = version
	resp, err := req.Post(url)
	return c.getRespErr(resp, err, forceCreate)
}

var versionFlags = struct {
	repoName   string
	pkgName    string
	pkgVersion string
}{}

func bintrayCreateVersionCommand() *cli.Command {
	return &cli.Command{
		Name:  "create-version",
		Usage: "create package version in repo",
		Flags: []*cli.Flag{
			{
				Name:   "repo",
				Usage:  "bintray repository name",
				Value:  &versionFlags.repoName,
				EnvVar: "BINTRAY_REPO",
			},
			{
				Name:   "package",
				Usage:  "bintray package name",
				Value:  &versionFlags.pkgName,
				EnvVar: "BINTRAY_PACKAGE",
			},
			{
				Name:   "version",
				Usage:  "bintray package version",
				Value:  &versionFlags.pkgVersion,
				EnvVar: "BINTRAY_VERSION",
			},
		},
		Action: func(_ *cli.Context) {
			if bintrayFlags.subject == "" {
				panic("no --subject provided")
			}
			if bintrayFlags.apikey == "" {
				panic("no --apikey provided")
			}
			if versionFlags.repoName == "" {
				panic("no --repo provided")
			}
			if versionFlags.pkgName == "" {
				panic("no --package provided")
			}
			if versionFlags.pkgVersion == "" {
				panic("no --version provided")
			}

			version := &bintrayVersion{
				Name: versionFlags.pkgVersion,
			}

			c := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
			err := c.versionCreate(version, versionFlags.repoName, versionFlags.pkgName, bintrayFlags.force)
			runs.PanicIfErr(err)
		},
	}
}
