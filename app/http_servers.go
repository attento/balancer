package app

import (
	"net/http"
	"time"

	"gopkg.in/tylerb/graceful.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/liuggio/events"
	"github.com/attento/balancer/app/core"
)

type HttpServers struct {
	servers  map[core.Address]*graceful.Server
	events   events.Dispatcher
	proxies map[core.Address]Reverser
}

type HttpServersHandler interface {
	ListenAndServe(a core.Address, h http.Handler) error
	Stop(a core.Address, t time.Duration) error
}

func NewHttpServers(e events.Dispatcher) *HttpServers {
	return &HttpServers{
		events:  e,
		servers: make(map[core.Address]*graceful.Server),
		proxies: make(map[core.Address]Reverser),
	}
}

func (s *HttpServers) Stop(a core.Address, t time.Duration) error {
	if _, ok := s.servers[a]; ok {
		log.Info("Stopping server")
		s.servers[a].Stop(t)
	}
	return nil
}

func (s *HttpServers) createProxyIfNotExists(a core.Address, r core.ConfigRepository) (Reverser, error) {

	if _, ok := s.proxies[a]; !ok {
		s.proxies[a] = NewReverse(a, r)
	}

	return s.proxies[a], nil
}

func (s *HttpServers) ChangeRoutes(a core.Address, h http.Handler) error {

	if _, ok := s.servers[a]; !ok {
		// @todo check server status
		log.Warn("Trying to change Routes on Server but it doens't exists...",a)
		return nil
	}
	s.servers[a].Handler = h

	return nil
}

func (s *HttpServers) ListenAndServe(a core.Address, h http.Handler) error {

	if _, ok := s.servers[a]; ok {
		// @todo check server status
		log.Warn("Trying to run server but was already on...",a)
		return nil
	}
	s.servers[a] = &graceful.Server{
		Server: &http.Server{
			Addr: string(a),
			Handler: h,
		},
	}

	s.events.Raise(core.EventHttpServerStarted, a)
	err := s.servers[a].ListenAndServe()
	if err != nil {
		log.Error(err)
		s.events.Raise(core.EventHttpServerStoppedWithError, a)
		return err
	}

	s.events.Raise(core.EventHttpServerStopped, a)
	return nil
}

func addHosts(rr *mux.Router, s *core.Server, d Reverser) {

	if len(s.Filter().Hosts) <= 0 {
		rr.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return true
		}).HandlerFunc(d.ServeHTTP)
		return
	}

	for _, v := range s.Filter().Hosts {
		 rr.Host(v).HandlerFunc(d.ServeHTTP)
	}
}

