package source

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/attento/balancer/app"
	"github.com/attento/balancer/app/core"
)

type configuration []core.Server

type FileExtractor struct {
	Daemon app.DaemonInterface
	Path   string
}

func (e *FileExtractor) Extract() error {
	var c configuration
	content, err := ioutil.ReadFile(e.Path)
	if err != nil {
		return errors.New("unable to open file")
	}
	err = json.Unmarshal(content, &c)
	if err != nil {
		return err
	}
	for _, s := range c {
		var upstreams []*core.Upstream
		for _, u := range s.Upstreams {
			upstreams = append(upstreams, u)
		}
		e.Daemon.StartHttpServer(s.Address, s.Filter, upstreams)
	}
	return nil
}
