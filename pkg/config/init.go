package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/subchen/go-stack/findup"
	"github.com/subchen/go-stack/fs"
	"github.com/subchen/go-stack/lookup"
	"github.com/subchen/storm/pkg/yaml"
)

const (
	configFile = ".storm.yml"
)

var (
	ProjectRoot    string
	ProjectGopath  string
	ProjectWorkdir string

	config RootConfig

	Project = &config.Project
	Build   = &config.Build
)

func init() {
	initProjectRoot()
	initConfig()
	initGopathWorkdir()
}

func initProjectRoot() {
	file, err := findup.Find(configFile)
	if err != nil {
		log.Fatal("no .storm.yml found")
	}

	ProjectRoot = filepath.Dir(file)
}

func initConfig() {
	// project.name
	config.Project.Name = filepath.Base(ProjectRoot)

	// project.version
	versionFile, err := lookup.FindAt(ProjectRoot, "VERSION", "VERSION.txt")
	if err == nil {
		version, err := ioutil.ReadFile(versionFile)
		if err != nil {
			panic(err)
		}
		config.Project.Version = string(version)
	}

	// project.module
	if pos := strings.Index(ProjectRoot, "/src/"); pos > 0 {
		config.Project.Module = ProjectRoot[pos+5:]
	}

	// load yaml file
	yamlFile := filepath.Join(ProjectRoot, configFile)
	err = yaml.ReadFile(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	setupDefaultConfig()
}

func initGopathWorkdir() {
	if strings.HasSuffix(ProjectRoot, "/src/"+config.Project.Module) {
		ProjectGopath = strings.TrimSuffix(ProjectRoot, "/src/"+config.Project.Module)
		ProjectWorkdir = ProjectRoot
		return
	}

	// set GOPATH
	ProjectGopath = filepath.Join(ProjectRoot, ".build/go")
	ProjectWorkdir = filepath.Join(ProjectGopath, "src", config.Project.Module)

	// create workdir symlink
	if !fs.Exists(ProjectWorkdir) {
		p := filepath.Dir(ProjectWorkdir)
		if !fs.IsDir(p) {
			err := os.MkdirAll(p, 0755)
			if err != nil {
				panic(err)
			}
		}
		err := os.Symlink(ProjectRoot, ProjectWorkdir)
		if err != nil {
			panic(err)
		}
	}
}

func setupDefaultConfig() {
	if len(config.Build.Binaries) == 0 {
		config.Build.Binaries = []config.BuildBinaryConfig{
			{
				Name: config.Project.Name,
				Path: ".",
			},
		}
	}

	if len(config.Build.Platforms) == 0 {
		config.Build.Platforms = []string{
			fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH),
		}
	}
}
