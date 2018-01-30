package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack"
	"github.com/subchen/go-stack/checksum"
	"github.com/subchen/go-stack/fs"
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
				if fs.DirExists(f) {
					files, err := ioutil.ReadDir(f)
					gstack.PanicIfErr(err)

					for _, file := range files {
						sha256sum(filepath.Join(f, file.Name()))
					}
				} else if fs.FileExists(f) {
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

	if !fs.FileExists(file) {
		return
	}

	fmt.Printf("sha256sum %s ...\n", file)
	err := checksum.Sha256sumAsFile(file)
	gstack.PanicIfErr(err)
}
