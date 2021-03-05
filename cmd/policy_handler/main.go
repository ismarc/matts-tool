package main

import (
	"log"
	"os"

	"github.com/ismarc/policy-handler/internal/app"
	"github.com/urfave/cli/v2"
)

func main() {
	var inputPolicyFile string

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "input-file",
			Aliases:     []string{"i"},
			Value:       "policy.yml",
			Usage:       "The base policy file to load",
			Destination: &inputPolicyFile,
		},
	}

	cli := &cli.App{
		Name:  "policy-handler",
		Usage: "Run policy handler",
		Flags: flags,
		Action: func(c *cli.Context) error {
			app.Run(inputPolicyFile)
			return nil
		},
	}

	err := cli.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
