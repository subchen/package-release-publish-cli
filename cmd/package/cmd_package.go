package packages

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/storm/pkg/config"
	"github.com/subchen/storm/pkg/sh"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "build go project",
		Action: runCommand,
	}
}

func runCommand(c *cli.Context) {
	binaries := config.Build.Binaries
	if len(binaries) == 0 {
		binaries = []config.BuildBinaryConfig{
			{
				Name: config.Project.Name,
				Path: ".",
			},
		}
	}

	platforms := config.Build.Platforms
	if len(platforms) == 0 {
		platforms = []string{
			fmt.Sprintf("%v/%v", runtime.GOOS, runtime.GOARCH),
		}
	}

	for _, bin := range binaries {
		for _, platform := range platforms {
			p := strings.Split(platform, "/")
			build(&bin, p[0], p[1])
		}
	}
}

func build(bin *config.BuildBinaryConfig, goos string, goarch string) {
	dest := filepath.Join(".build", fmt.Sprintf("%s-%s-%s", bin.Name, goos, goarch))

	flags := strings.Split(config.Build.Flags, " ")

	ldflags := strings.Replace(config.Build.Ldflags, "\n", " ", -1)
	// if goos != "darwin" && goos != "solaris" && !strings.Contains(ldflags, "-static") {
	// 	ldflags += ` -extldflags '-static'`
	// }

	args := []string{"build", "-ldflags", ldflags, "-o", dest}
	if len(flags) > 0 {
		args = append(args, flags...)
	}
	args = append(args, bin.Path)

	fmt.Printf("building: %s ...\n", dest)
	cmd := &sh.Command{
		Command: "go",
		Args:    args,
		PipeOut: true,
		Dir:     config.ProjectWorkdir,
		Env: append(config.Build.Env, []string{
			"GOPATH=" + config.ProjectGopath,
			"GOOS=" + goos,
			"GOARCH=" + goarch,
		}...),
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}
