package core

import (
	"fmt"
	"net/url"
	"strconv"
)

type Upstream struct {
	Target   string
	Port     uint16
	Priority uint16
	Weight   uint16
}

const (
	EventConfigServerCreated string = "cnf-srv-created"

	EventHttpServerStarted          string = "http-srv-started"
	EventHttpServerStopped          string = "http-srv-stopped"
	EventHttpServerStoppedWithError string = "http-srv-stopped-err"

	EventConfigUpstreamsUpdated string = "ups-updated"
	EventConfigFilterUpdated    string = "flt-updated"
)

func (u *Upstream) ToUrl(scheme string) (url *url.URL, err error) {

	if scheme == "" {
		scheme = "http"
	}

	return url.Parse(scheme + "://" + u.Target + ":" + strconv.Itoa(int(u.Port)))
}

type Filter struct {
	Hosts      []string
	Schemes    [2]string
	PathPrefix string
}

type Server struct {
	Address   Address              `json:"address"`
	Filter    Filter               `json:"filter"`
	Upstreams map[string]*Upstream `json:"upstreams"`
}

func newServer(a Address) *Server {
	return &Server{Address: a}
}

func (s *Server) Upstream(target string, port uint16) (k *Upstream, ok bool) {
	k, ok = s.Upstreams[CreateUpstreamKey(target, port)]
	return
}

func (s *Server) addUpstreamProperty(target string, port uint16, priority uint16, weight uint16) {

	if s.Upstreams == nil {
		s.Upstreams = make(map[string]*Upstream)
	}

	s.Upstreams[CreateUpstreamKey(target, port)] = &Upstream{target, port, priority, weight}
}

func (s *Server) addUpstream(u *Upstream) {

	if s.Upstreams == nil {
		s.Upstreams = make(map[string]*Upstream)
	}

	s.Upstreams[CreateUpstreamKey(u.Target, u.Port)] = u
}

func (s *Server) setUpstreams(us map[string]*Upstream) {

	if s.Upstreams == nil {
		s.Upstreams = make(map[string]*Upstream)
	}

	s.Upstreams = us
}

func (s *Server) removeUpstream(target string, port uint16) {
	if _, ok := s.Upstream(target, port); ok {
		delete(s.Upstreams, CreateUpstreamKey(target, port))
	}
}

type Address string // ":80" golang address spec.

type Config struct {
	servers map[Address]*Server
}

func Create() Config {
	return Config{servers: make(map[Address]*Server)}
}

func (c Config) NewServer(a Address) bool {

	if _, ok := c.servers[a]; ok {
		return false
	}

	c.servers[a] = newServer(a)
	return true
}

// Put with empty values if you don't need filter
// eg. []string{}, [2]string{}, ""
func (c Config) PutFilterProperties(address Address, hosts []string, schemes [2]string, prefix string) {

	c.PutFilter(address, Filter{hosts, schemes, prefix})
}

func (c Config) PutFilter(address Address, f Filter) {

	if _, ok := c.servers[address]; !ok {
		c.NewServer(address)
	}
	c.servers[address].Filter = f
}

func (c Config) RemoveServer(address Address) {
	if _, ok := c.servers[address]; ok {
		delete(c.servers, address)
	}
}

func (c Config) AddUpstreamProperty(address Address, target string, port uint16, priority uint16, weight uint16) {

	if _, ok := c.servers[address]; !ok {
		c.NewServer(address)
	}

	c.servers[address].addUpstreamProperty(target, port, priority, weight)
}

func (c Config) AddUpstream(address Address, u *Upstream) {

	if _, ok := c.servers[address]; !ok {
		c.NewServer(address)
	}

	c.servers[address].addUpstream(u)
}

func (c Config) SetUpstreams(address Address, us map[string]*Upstream) {

	if _, ok := c.servers[address]; !ok {
		c.NewServer(address)
	}

	c.servers[address].setUpstreams(us)
}

func (c Config) Servers() map[Address]*Server {
	return c.servers
}

func (c Config) Server(address Address) (s *Server, ok bool) {
	s, ok = c.servers[address]
	return
}

func (c Config) RemoveUpstream(address Address, target string, port uint16) {

	if _, ok := c.servers[address]; !ok {
		return
	}

	c.servers[address].removeUpstream(target, port)
}

func CreateUpstreamKey(target string, port uint16) string {
	return fmt.Sprintf("%s:%d", target, port)
}
