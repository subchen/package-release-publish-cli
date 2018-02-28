package bintray

import (
	"fmt"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/runs"
)

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

// POST /repos/:subject/:repo
func (c *bintrayClient) repoCreate(repo *bintrayRepo, forceCreate bool) error {
	url := fmt.Sprintf("%s/repos/%s/%s", BINTRAY_API_PREFIX, c.subject, repo.Name)
	req := c.newRequest()
	req.JSON = repo
	resp, err := req.Post(url)
	return c.getRespErr(resp, err, forceCreate)
}

var repoFlags = struct {
	repoName string
	repoType string
}{}

func bintrayCreateRepoCommand() *cli.Command {
	return &cli.Command{
		Name:  "create-repo",
		Usage: "create repository",
		Flags: []*cli.Flag{
			{
				Name:   "name",
				Usage:  "bintray repository name",
				Value:  &repoFlags.repoName,
				EnvVar: "BINTRAY_REPO",
			},
			{
				Name:     "type",
				Usage:    "bintray repository type (maven, debian, rpm, docker, npm, generic, ...)",
				Value:    &repoFlags.repoType,
				DefValue: "generic",
			},
		},
		Action: func(_ *cli.Context) {
			if bintrayFlags.subject == "" {
				panic("no --subject provided")
			}
			if bintrayFlags.apikey == "" {
				panic("no --apikey provided")
			}
			if repoFlags.repoName == "" {
				panic("no --name provided")
			}

			repo := &bintrayRepo{
				Name: repoFlags.repoName,
				Type: repoFlags.repoType,
			}

			c := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
			err := c.repoCreate(repo, bintrayFlags.force)
			runs.PanicIfErr(err)
		},
	}
}
