package bintray

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/fs"
	"github.com/subchen/go-stack/iif"
	"github.com/subchen/go-stack/runs"
)

// PUT /content/:subject/:repo/:package/:version/:file_path[?publish=0/1][?override=0/1][?explode=0/1]
func (c *bintrayClient) bintrayFileUpload(repo string, pkg string, version string, filepath string, fileContent io.Reader, forceCreate bool) error {
	url := fmt.Sprintf("%s/content/%s/%s/%s/%s/%s", BINTRAY_API_PREFIX, c.subject, repo, pkg, version, filepath)
	fmt.Println(url)
	req := c.newRequest()
	req.Headers = map[string]string{
		"X-Bintray-Publish":  "1",
		"X-Bintray-Override": iif.String(forceCreate, "1", "0"),
		"X-Bintray-Explode":  "1",
	}
	req.Body = fileContent
	resp, err := req.Put(url)
	return c.getRespErr(resp, err, forceCreate)
}

var uploadFlags = struct {
	repoName   string
	pkgName    string
	pkgVersion string
	path       string
}{}

func bintrayFileUploadCommand() *cli.Command {
	return &cli.Command{
		Name:  "upload",
		Usage: "upload files",
		Flags: []*cli.Flag{
			{
				Name:   "repo",
				Usage:  "bintray repository name",
				Value:  &uploadFlags.repoName,
				EnvVar: "BINTRAY_REPO",
			},
			{
				Name:   "package",
				Usage:  "bintray package name",
				Value:  &uploadFlags.pkgName,
				EnvVar: "BINTRAY_PACKAGE",
			},
			{
				Name:   "version",
				Usage:  "bintray package version",
				Value:  &uploadFlags.pkgVersion,
				EnvVar: "BINTRAY_VERSION",
			},
			{
				Name:  "path",
				Usage: "file path in url",
				Value: &uploadFlags.path,
			},
		},
		Action: func(c *cli.Context) {
			if bintrayFlags.subject == "" {
				panic("no --subject provided")
			}
			if bintrayFlags.apikey == "" {
				panic("no --apikey provided")
			}
			if uploadFlags.repoName == "" {
				panic("no --repo provided")
			}
			if uploadFlags.pkgName == "" {
				panic("no --package provided")
			}
			if uploadFlags.pkgVersion == "" {
				panic("no --version provided")
			}

			if c.NArg() == 0 {
				panic("no dir or file provided")
			}

			for _, f := range c.Args() {
				if fs.IsDir(f) {
					files, err := ioutil.ReadDir(f)
					runs.PanicIfErr(err)

					for _, file := range files {
						uploadFile(filepath.Join(f, file.Name()))
					}
				} else if fs.IsFile(f) {
					uploadFile(f)
				} else {
					panic("file not exists: " + f)
				}
			}
		},
	}
}

func uploadFile(file string) {
	f, err := os.Open(file)
	runs.PanicIfErr(err)
	defer f.Close()

	c := newBintrayClient(bintrayFlags.subject, bintrayFlags.apikey)
	err = c.bintrayFileUpload(
		uploadFlags.repoName,
		uploadFlags.pkgName,
		uploadFlags.pkgVersion,
		filepath.Join(uploadFlags.path, filepath.Base(file)),
		f,
		bintrayFlags.force,
	)
	runs.PanicIfErr(err)
}
