package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/go-stack"
	"github.com/subchen/go-stack/fs"
)

var gobuildFlags = struct {
	name      string
	version   string
	goos      string
	goarch    string
	archived  bool
	sourceDir string
	outputDir string
}{}

func gobuildCommand() *cli.Command {
	return &cli.Command{
		Name:  "gobuild",
		Usage: "build go sources",
		Flags: []*cli.Flags{
			{
				Name:  "name",
				Desc:  "application name",
				Value: &gobuildFlags.name,
			},
			{
				Name:  "version",
				Desc:  "application version",
				Value: &gobuildFlags.version,
			},
			{
				Name:     "goos",
				Desc:     "go build target os: GOOS",
				Value:    &gobuildFlags.goos,
				DefValue: "linux,darwin,windows",
			},
			{
				Name:     "goarch",
				Desc:     "go build target arch: GOARCH",
				Value:    &gobuildFlags.goarch,
				DefValue: "amd64",
			},
			{
				Name:     "z,zip",
				Desc:     "output archive file",
				Value:    &gobuildFlags.archived,
				DefValue: "false",
			},
			{
				Name:     "s,source-dir",
				Desc:     "go sources dir",
				Value:    &gobuildFlags.sourceDir,
				DefValue: ".",
			},
			{
				Name:     "o,output-dir",
				Desc:     "build target dir",
				Value:    &gobuildFlags.outputDir,
				DefValue: "./_releases",
			},
		},
		Action: func(c *cli.Context) {
			if gobuildFlags.name == "" {
				panic("no --name provided")
			}
			if gobuildFlags.version == "" {
				panic("no --version provided")
			}
			if !fs.DirExists(gobuildFlags.sourceDir) {
				panic("source-dir does not exists")
			}

			if !fs.DirExists(gobuildFlags.outputDir) {
				os.MkdirAll(gobuildFlags.outputDir)
			}

			gobuild()
		},
	}
}

func gobuild() {

}
