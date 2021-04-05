package main

import (
	"log"
	"os"

	"github.com/ismarc/policy-handler/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	var inputPolicyFile string
	var sourceDBConnUrl string
	var destinationDBConnUrl string

	policyFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "input-file",
			Aliases:     []string{"i"},
			Value:       "policy.yml",
			Usage:       "The base policy file to load",
			Destination: &inputPolicyFile,
		},
	}

	dbFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "source",
			Aliases:     []string{"s"},
			Usage:       "The source database connection url",
			Destination: &sourceDBConnUrl,
		},
		&cli.StringFlag{
			Name:        "destination",
			Aliases:     []string{"d"},
			Usage:       "The destination database connection url",
			Destination: &destinationDBConnUrl,
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
			Name:  "db",
			Usage: "Perform db related operations",
			Flags: dbFlags,
			Action: func(c *cli.Context) error {
				app.RunDB(sourceDBConnUrl, destinationDBConnUrl)
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
