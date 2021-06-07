package main

import (
	"log"
	"os"

	"github.com/ismarc/matts-tool/internal/app"
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

	dbFlags := []cli.Flag{
		&cli.StringFlag{
			Name:    "source",
			Aliases: []string{"s"},
			Value:   "",
			Usage:   "Source filename containing appropriate pg_dump data to load",
		},
		&cli.StringFlag{
			Name:    "destination-dsn",
			Aliases: []string{"d"},
			Value:   "",
			Usage:   "DSN for connecting to the destination postgres database",
		},
		&cli.StringFlag{
			Name:    "account",
			Aliases: []string{"a"},
			Value:   "conjur",
			Usage:   "Conjur account to use for the destination values",
		},
		&cli.BoolFlag{
			Name:    "no-act",
			Aliases: []string{"n"},
			Value:   false,
			Usage:   "Only process the data and display what would be written to the database",
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
		&cli.BoolFlag{
			Name:  "skip-same-value",
			Value: true,
			Usage: "Skip writing variables that have the same value in both systems",
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
			Name: "db",
			Usage: `Perform db related operations.
			Decrypt values from pg_dump file and re-encrypt and add to running instance.
			IN_CONJUR_DATA_KEY -- key to use for decryption of values in file
			OUT_CONJUR_DATA_KEY -- key to use to encrypt values for insertion.

			Generate the input data file:
			pg_dump --data-only --schema="authn" --table="authn.users" > ~/data.sql
			`,
			Flags: dbFlags,
			Action: func(c *cli.Context) error {
				sourceDataKey := os.Getenv("IN_CONJUR_DATA_KEY")
				destinationDataKey := os.Getenv("OUT_CONJUR_DATA_KEY")
				config := app.DBConfig{
					SourceFilename:     c.String("source"),
					SourceDataKey:      sourceDataKey,
					DestinationDSN:     c.String("destination-dsn"),
					DestinationDataKey: destinationDataKey,
					DestinationAccount: c.String("account"),
					NoAct:              c.Bool("no-act"),
				}
				app.RunDB(config)
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
					SkipSameValue:       c.Bool("skip-same-value"),
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
