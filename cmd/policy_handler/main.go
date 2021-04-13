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
	var continueOnError bool

	policyFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "input-file",
			Aliases:     []string{"i"},
			Value:       "policy.yml",
			Usage:       "The base policy file to load",
			Destination: &inputPolicyFile,
		},
		&cli.BoolFlag{
			Name:  "strip-annotations",
			Value: false,
			Usage: "Whether to strip any annotations from the resulting policy",
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
			Name:  "source-netrc",
			Value: "~/.netrc",
			Usage: "The netrc file to use for credentials for the source instance",
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
			Name:  "destination-netrc",
			Value: "~/.netrc",
			Usage: "The netrc file to use for credentials for the destination instance",
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
		&cli.BoolFlag{
			Name:        "continue-on-error",
			Aliases:     []string{"c"},
			Value:       true,
			Usage:       "Continue processing when receiving an error reading or writing a variable",
			Destination: &continueOnError,
		},
		&cli.IntFlag{
			Name:  "resource-batch-size",
			Value: 25,
			Usage: "Number of resources to fetch in a single call to conjur",
		},
		&cli.IntFlag{
			Name:  "variable-batch-size",
			Value: 10,
			Usage: "Number of variable values to fetch in a single call to conjur",
		},
	}

	commands := []*cli.Command{
		{
			Name:  "policy",
			Usage: "Perform policy related operations",
			Flags: policyFlags,
			Action: func(c *cli.Context) error {
				app.RunPolicy(inputPolicyFile, c.Bool("strip-annotations"))
				return nil
			},
		},
		{
			Name:  "api",
			Usage: "Perform api related operations",
			Flags: apiFlags,
			Action: func(c *cli.Context) error {
				config := app.APIConfig{
					SourceConjurRC:      sourceConjurRC,
					SourceNetRC:         c.String("source-netrc"),
					SourceVersion:       sourceVersion,
					DestinationConjurRC: destinationConjurRC,
					DestinationNetRC:    c.String("destination-netrc"),
					DestinationVersion:  destinationVersion,
					NoAct:               noAct,
					ContinueOnError:     continueOnError,
					ResourceBatchSize:   c.Int("resource-batch-size"),
					VariableBatchSize:   c.Int("variable-batch-size"),
				}
				app.RunAPI(config)
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
