package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/subchen/go-tableify"
)

type bintrayClient struct {
	client  *http.Client
	subject string
	apikey  string
}

type bintrayRepo struct {
	Name            string   `json:"name" tableify:"-"`
	Owner           string   `json:"owner"`
	Type            string   `json:"type" tableify:"-"`
	Private         bool     `json:"private"`
	BusinessUnit    string   `json:"business_unit"`
	Premium         bool     `json:"premium"`
	Desc            string   `json:"desc"`
	Labels          []string `json:"labels"`
	Created         string   `json:"created"`
	PackageCount    int      `json:"package_count" tableify:"-"`
	GpgSignMetadata bool     `json:"gpg_sign_metadata"`
	GpgSignFiles    bool     `json:"gpg_sign_files"`
	GpgUseOwnerKey  bool     `json:"gpg_use_owner_key"`

	//
	LastUpdated string `json:"lastUpdated" tableify:"-"`

	// Enterprise version
	VersionUpdateMaxDays int `json:"version_update_max_days"`

	// type=debian
	DefaultDebianArchitecture string `json:"default_debian_architecture"`
	DefaultDebianComponent    string `json:"default_debian_component"`
	DefaultDebianDistribution string `json:"default_debian_distribution"`

	// type=rpm
	YumGroupsFile    string `json:"yum_groups_file"`
	YumMetadataDepth int    `json:"yum_metadata_depth"`
}

type bintrayPackage struct {
	AttributeNames         []string      `json:"attribute_names"`
	Attributes             string        `json:"attributes"`
	Created                string        `json:"created"`
	CustomLicenses         []string      `json:"custom_licenses"`
	Desc                   string        `json:"desc"`
	FollowersCount         int           `json:"followers_count"`
	GithubReleaseNotesFile string        `json:"github_release_notes_file"`
	GithubRepo             string        `json:"github_repo"`
	IssueTrackerURL        string        `json:"issue_tracker_url"`
	Labels                 []string      `json:"labels"`
	LatestVersion          string        `json:"latest_version"`
	Licenses               []string      `json:"licenses"`
	LinkedToRepos          []interface{} `json:"linked_to_repos"`
	Name                   string        `json:"name"`
	Owner                  string        `json:"owner"`
	Permissions            []interface{} `json:"permissions"`
	PublicDownloadNumbers  bool          `json:"public_download_numbers"`
	PublicStats            bool          `json:"public_stats"`
	Rating                 int           `json:"rating"`
	RatingCount            int           `json:"rating_count"`
	Repo                   string        `json:"repo"`
	SystemIds              []interface{} `json:"system_ids"`
	Updated                string        `json:"updated"`
	VcsURL                 string        `json:"vcs_url"`
	Versions               []string      `json:"versions"`
	WebsiteURL             string        `json:"website_url"`
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

func testBintray() {
	c := newBintrayClient("subchen", "3bd86a64193bd6ecb27b3189aec576bd0eb03d63")
	repolist, err := c.RepoList()
	if err != nil {
		fmt.Println(err)
	}

	t := tableify.New()
	t.SetHeadersFromStruct(new(bintrayRepo))
	t.AddRowObjectList(repolist)
	t.Print()
}
