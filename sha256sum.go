package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack"
	"github.com/ungerik/go-dry"
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
				if !dry.FileExists(f) {
					panic("file not exists: " + f)
				}
				if dry.FileIsDir(f) {
					files, err := ioutil.ReadDir(f)
					gstack.PanicIfErr(err)

					for _, file := range files {
						sha256sum(filepath.Join(f, file.Name()))
					}
				} else {
					sha256sum(f)
				}
			}
		},
	}
}

func sha256sum(file string) {
	if strings.HasSuffix(file, ".sha256") {
		return
	}

	if dry.FileIsDir(file) {
		return
	}

	fmt.Printf("sha256sum %s ...\n", file)
	err := gstack.Sha256sumFile(file)
	gstack.PanicIfErr(err)
}
