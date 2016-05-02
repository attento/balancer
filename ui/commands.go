package ui

import "github.com/codegangsta/cli"

var (
	commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the balancer daemon",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "exporter, e",
					Value: "none",
					Usage: "Exporter",
				},
			},
			Action: run,
		},
	}
)
