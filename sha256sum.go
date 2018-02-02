package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack/encoding/sha256"
	"github.com/subchen/go-stack/fs"
	"github.com/subchen/go-stack/runs"
)

func sha256sumCommand() *cli.Command {
	return &cli.Command{
		Name:  "sha256sum",
		Usage: "add .sha256 checksum file",
		Action: func(c *cli.Context) {
			if c.NArg() == 0 {
				panic("no dir or file provided")
			}

			for _, f := range c.Args() {
				if fs.IsDir(f) {
					files, err := ioutil.ReadDir(f)
					runs.PanicIfErr(err)

					for _, file := range files {
						sha256sum(filepath.Join(f, file.Name()))
					}
				} else if fs.IsFile(f) {
					sha256sum(f)
				} else {
					panic("file not exists: " + f)
				}
			}
		},
	}
}

func sha256sum(file string) {
	if strings.HasSuffix(file, ".sha256") {
		return
	}

	if !fs.IsFile(file) {
		return
	}

	fmt.Printf("sha256sum %s ...\n", file)
	err := sha256.GenerateSumFile(file)
	runs.PanicIfErr(err)
}
