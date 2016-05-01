package ui

import "github.com/codegangsta/cli"

var (
	commands = []cli.Command{
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "Run the balancer daemon",
			//Flags: []cli.Flag{flHeartBeat,flTTL,flDiscoveryOpt},
			Action: run,
		},
	}
)

