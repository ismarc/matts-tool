package main

import (
	"log"
	"os"

	"github.com/ismarc/policy-handler/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	var inputPolicyFile string
	var sourceConjurRC string
	var sourceVersion string
	var destinationConjurRC string
	var destinationVersion string
	var noAct bool

	policyFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "input-file",
			Aliases:     []string{"i"},
			Value:       "policy.yml",
			Usage:       "The base policy file to load",
			Destination: &inputPolicyFile,
		},
	}

	apiFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "source-conjurrc",
			Aliases:     []string{"s"},
			Value:       "",
			Usage:       "The conjurrc file to use as the source for syncing data",
			Destination: &sourceConjurRC,
		},
		&cli.StringFlag{
			Name:        "source-version",
			Value:       "4",
			Usage:       "The major API version of the source for syncing data",
			Destination: &sourceVersion,
		},
		&cli.StringFlag{
			Name:        "destination-conjurrc",
			Aliases:     []string{"d"},
			Value:       "",
			Usage:       "The conjurrc file to use as the destination for syncing data",
			Destination: &destinationConjurRC,
		},
		&cli.StringFlag{
			Name:        "destination-version",
			Value:       "5",
			Usage:       "The major API version of the destination for syncing data",
			Destination: &destinationVersion,
		},
		&cli.BoolFlag{
			Name:        "no-act",
			Aliases:     []string{"n"},
			Value:       false,
			Usage:       "Do not read or write variables of data, but fetch the resources that would be synced",
			Destination: &noAct,
		},
	}

	commands := []*cli.Command{
		{
			Name:  "policy",
			Usage: "Perform policy related operations",
			Flags: policyFlags,
			Action: func(c *cli.Context) error {
				app.RunPolicy(inputPolicyFile)
				return nil
			},
		},
		{
			Name:  "api",
			Usage: "Perform api related operations",
			Flags: apiFlags,
			Action: func(c *cli.Context) error {
				app.RunApi(sourceConjurRC, sourceVersion, destinationConjurRC, destinationVersion, noAct)
				return nil
			},
		},
	}

	cli := &cli.App{
		Commands: commands,
	}

	err := cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
