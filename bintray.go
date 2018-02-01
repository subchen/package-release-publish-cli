package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type bintrayClient struct {
	client  *http.Client
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
	PublicStats            bool     `json:"public_stats"`
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
		client:  http.DefaultClient,
		subject: subject,
		apikey:  apikey,
	}
}

func (c *bintrayClient) newRequestWithReader(method, url string, requestReader io.Reader, requestLength int64) (*http.Request, error) {
	req, err := http.NewRequest(method, url, requestReader)
	if err != nil {
		return nil, err
	}
	if requestLength > 0 {
		req.ContentLength = int64(requestLength)
	}
	if c.subject != "" {
		req.SetBasicAuth(c.subject, c.apikey)
	}
	return req, nil
}

// GET /repos/:subject
func (c *bintrayClient) RepoList() ([]*bintrayRepo, error) {
	url := "https://api.bintray.com/repos/" + c.subject

	req, err := c.newRequestWithReader("GET", url, nil, 0)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		}

		v := make([]*bintrayRepo, 0)
		err = json.Unmarshal(body, &v)
		return v, err
	}

	return nil, errors.New("status_code is not 200")
}


var bintrayFlags = struct {
	subject string
	apikey  string
}{}

func bintrayUpload() {
	c := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
}
