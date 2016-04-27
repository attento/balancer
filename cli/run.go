package cli

import "github.com/codegangsta/cli"
import "github.com/attento/balancer/api"

func run(c *cli.Context) {
	backend := ":9123"
	api.Run(backend, c.Bool("debug"))
}

