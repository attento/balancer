package ui

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
)

import (
	"github.com/attento/balancer/app"
	"github.com/attento/balancer/ui/api"
	"github.com/attento/balancer/ui/source"
)

func run(c *cli.Context) {
	backend := ":9123"
	d := app.NewStandard()
	exporter := c.String("exporter")
	if exporter == "file" {
		e := source.FileExtractor{d, "/etc/balancer.json"}
		err := e.Extract()
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}
	api.Run(d, backend, c.Bool("debug"))
}
