package app

import (
	"errors"
	"time"

	"github.com/attento/balancer/app/core"
	"github.com/attento/balancer/app/repository"
	"github.com/liuggio/events"
)

type DaemonInterface interface {
	Config() *core.Config
	ConfigServer(core.Address) (*core.Server, bool, error)
	PutConfigFilter(core.Address, core.Filter) error
	AddConfigUpstream(core.Address, *core.Upstream) error

	StartHttpServer(core.Address, core.Filter, []*core.Upstream) error
	StopHttpServer(core.Address, time.Duration) error
}

type daemon struct {
	repo        core.ConfigRepository
	httpServers *HttpServers
	e           events.Dispatcher
}

var ServerNotFound = errors.New("Server not found.")

func New(repo core.ConfigRepository, e events.Dispatcher) *daemon {

	d := &daemon{
		httpServers: NewHttpServers(e),
		repo:        repo,
		e:           e,
	}

	d.e.On(core.EventConfigServerCreated, d.onNewConfigServer)
	d.e.On(core.EventHttpServerStopped, d.onHttpServerStopped)
	d.e.On(core.EventHttpServerStoppedWithError, d.onHttpServerStoppedWithError)
	d.e.On(core.EventConfigFilterUpdated, d.onConfigServerChangedFilter)

	return d
}

func NewStandard() *daemon {
	return New(repository.NewInMemoryConfigRepository(), events.New())
}

func NewStandardRepo(repo core.ConfigRepository) *daemon {
	return New(repo, events.New())
}

func (d *daemon) Config() *core.Config {
	cnf := d.repo.Get()
	return &cnf
}

func (d *daemon) ConfigServer(a core.Address) (s *core.Server, ok bool, err error) {
	s, ok = d.repo.Server(a)
	err = nil
	return
}

// @todo updated filters change servers rr
func (d *daemon) StartHttpServer(a core.Address, f core.Filter, us []*core.Upstream) error {
	created := d.repo.NewServer(a)
	d.repo.SetUpstreams(a, us)
	d.repo.PutFilter(a, f)
	ss, ok := d.repo.Server(a)
	if !ok {
		return ServerNotFound
	}
	if created {
		d.e.Raise(core.EventConfigServerCreated, ss)
	} else {
		d.e.Raise(core.EventConfigUpstreamsUpdated, us)
		d.e.Raise(core.EventConfigFilterUpdated, us)
	}

	return nil
}

func (d *daemon) StopHttpServer(a core.Address, t time.Duration) error {
	d.httpServers.Stop(a, t)
	return nil
}

func (d *daemon) PutConfigFilter(a core.Address, f core.Filter) error {
	d.repo.PutFilter(a, f)
	srv, ok := d.repo.Server(a)
	if !ok {
		return errors.New("problem updating filter")
	}

	d.e.Raise(core.EventConfigFilterUpdated, srv)
	return nil
}

func (d *daemon) AddConfigUpstream(a core.Address, u *core.Upstream) error {
	d.repo.AddUpstream(a, u)
	srv, ok := d.repo.Server(a)
	if !ok {
		return errors.New("problem updating Upstream")
	}

	d.e.Raise(core.EventConfigUpstreamsUpdated, srv)
	return nil
}
