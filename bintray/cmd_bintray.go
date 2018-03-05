package bintray

import (
	"github.com/subchen/go-cli"
)

var bintrayFlags = struct {
	subject string
	apikey  string
	force   bool
}{}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "bintray",
		Usage: "bintray cli",
		Flags: []*cli.Flag{
			{
				Name:   "subject",
				Usage:  "bintray subject",
				EnvVar: "BINTRAY_SUBJECT",
				Value:  &bintrayFlags.subject,
			},
			{
				Name:   "apikey",
				Usage:  "bintray apikey",
				EnvVar: "BINTRAY_APIKEY",
				Value:  &bintrayFlags.apikey,
			},
			{
				Name:     "force",
				Usage:    "dont error if exists",
				Value:    &bintrayFlags.force,
				DefValue: "false",
			},
		},
		Commands: []*cli.Command{
			bintrayRepoCreateCommand(),
			bintrayPackageCreateCommand(),
			bintrayVersionCreateCommand(),
			bintrayFileUploadCommand(),
		},
	}
}
