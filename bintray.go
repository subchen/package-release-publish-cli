package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	
	"github.com/mozillazg/request"
)

type bintrayClient struct {
	subject string
	apikey  string
}

type bintrayRepo struct {
	Name            string   `json:"name"`
	Type            string   `json:"type"`
	Private         bool     `json:"private"`
	BusinessUnit    string   `json:"business_unit"`
	Desc            string   `json:"desc"`
	Labels          []string `json:"labels"`
	GpgSignMetadata bool     `json:"gpg_sign_metadata,omitempty"`
	GpgSignFiles    bool     `json:"gpg_sign_files,omitempty"`
	GpgUseOwnerKey  bool     `json:"gpg_use_owner_key,omitempty"`
}

type bintrayPackage struct {
	Name                   string   `json:"name"`
	Desc                   string   `json:"desc"`
	Labels                 []string `json:"labels"`
	Licenses               []string `json:"licenses"`
	VcsURL                 string   `json:"vcs_url"`
	WebsiteURL             string   `json:"website_url"`
	IssueTrackerURL        string   `json:"issue_tracker_url"`
	GithubRepo             string   `json:"github_repo"`
	GithubReleaseNotesFile string   `json:"github_release_notes_file"`
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
		apikey : apikey,
	}
}

func (c *bintrayClient) newRequest() *request.Request {
	req := request.NewRequest(http.DefaultClient)
	req.BasicAuth = request.BasicAuth{c.subject, c.apikey}
	return req
}

// POST /repos/:subject/:repo
func (c *bintrayClient) repoCreate(repo *bintrayRepo, force bool) error {
	url := fmt.Sprintf("https://api.bintray.com/repos/%s/%s", c.subject, repo.Name)
	req := c.newRequest()
	req.Json = repo
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		return errors.New(resp.Reason())
	}
	return nil
}

// POST /packages/:subject/:repo
func (c *bintrayClient) packageCreate(pkg *bintrayPackage, repo string) error {
	url := fmt.Sprintf("https://api.bintray.com/packages/%s/%s", c.subject, repo)
	req := c.newRequest()
	req.Json = pkg
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		return errors.New(resp.Reason())
	}
	return nil
}

// POST /packages/:subject/:repo/:package/versions
func (c *bintrayClient) packageCreate(version *bintrayVersion, repo string, pkg string) error {
	url := fmt.Sprintf("https://api.bintray.com/packages/%s/%s/%s/versions", c.subject, repo, pkg)
	req := c.newRequest()
	req.Json = version
	resp, err := req.Post(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		return errors.New(resp.Reason())
	}
	return nil
}

// PUT /content/:subject/:repo/:package/:version/:file_path[?publish=0/1][?override=0/1][?explode=0/1]
func (c *bintrayClient) bintrayUpload(repo string, pkg string, version string, path string, fileContent io.Reader) error {
	url := fmt.Sprintf("https://api.bintray.com/content/%s/%s/%s/%s", c.subject, repo, pkg, version, path)
	req := c.newRequest()
	req.Headers = map[string]string{
		"X-Bintray-Publish": "1",
		"X-Bintray-Override": "1",
		"X-Bintray-Explode": "1",
	}
	req.Body = fileContent
	resp, err := req.Put(url)
	if err != nil {
		return err
	}
	if !resp.OK() {
		return errors.New(resp.Reason())
	}
	return nil
}

var bintrayFlags = struct {
	subject string
	apikey  string
}{}

func bintrayUpload() {
	c := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
}
