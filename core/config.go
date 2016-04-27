package core

import (
	"fmt"
	"strconv"
	"net/url"
)

type Upstream struct {
	Target   string
	Port     uint16
	Priority uint16
	Weight   uint16
}

func (u *Upstream) toUrl(scheme string) (url *url.URL, err error) {
	return url.Parse(scheme+"://"+u.Target+":"+ strconv.Itoa(int(u.Port)))
}

type Filter struct {
	Hosts      []string
	Schemes    [2]string
	PathPrefix string
}

type Server struct {
	address   Address
	filter    Filter
	upstreams  map[string]*Upstream
}

func newServer(a Address) *Server {
	return &Server{address: a}
}

func (s *Server) Address() Address {
	return s.address
}

func (s *Server) Upstreams() map[string]*Upstream {
	return s.upstreams
}

func (s *Server) Filter() Filter {
	return s.filter
}

func (s *Server) Upstream(target string, port uint16) (k *Upstream, ok bool) {
	k, ok = s.upstreams[createUpstreamKey(target, port)]
	return
}

func (s *Server) addUpstreamProperty(target string, port uint16, priority uint16, weight uint16) {

	if s.upstreams == nil {
		s.upstreams = make(map[string]*Upstream)
	}

	s.upstreams[createUpstreamKey(target, port)] = &Upstream{target, port, priority, weight}
}

func (s *Server) addUpstream(u *Upstream) {

	if s.upstreams == nil {
		s.upstreams = make(map[string]*Upstream)
	}

	s.upstreams[createUpstreamKey(u.Target, u.Port)] = u
}

func (s *Server) removeUpstream(target string, port uint16) {
	if _, ok := s.Upstream(target, port); ok {
		delete(s.upstreams, createUpstreamKey(target, port))
	}
}

type Address string        // ":80" golang address spec.
type Config map[Address]*Server

func Create() Config {
	return Config{}
}

func (c Config) newServer(a Address) {
	create := false
	if _, ok := c[a]; !ok {
		create = true
	}
	s := newServer(a)
	c[a] = s
	if create {
		go listeners.onConfigCreatedNewAddress(a, s)
	}
}

// Put with empty values if you don't need filter
// eg. []string{}, [2]string{}, ""
func (c Config) PutFilterProperties(address Address, hosts []string, schemes [2]string, prefix string) {

	if _, ok := c[address]; !ok {
		c.newServer(address)
	}
	c[address].filter = Filter{hosts, schemes, prefix}
}


func (c Config) PutFilter(address Address, f Filter) {

	if _, ok := c[address]; !ok {
		c.newServer(address)
	}
	c[address].filter = f
}


func (c Config) RemoveServer(address Address) {
	if _, ok := c[address]; !ok {
		go listeners.onConfigUpdateRemovedAddress(address)
		delete(c, address)
	}
}

func (c Config) AddUpstreamProperty(address Address, target string, port uint16, priority uint16, weight uint16) {

	if _, ok := c[address]; !ok {
		c.newServer(address)
	}

	c[address].addUpstreamProperty(target, port, priority, weight)
}

func (c Config) AddUpstream(address Address, u *Upstream) {

	if _, ok := c[address]; !ok {
		c.newServer(address)
	}

	c[address].addUpstream(u)
}

func (c Config) Server(address Address) (s *Server, ok bool){
	s, ok = c[address];
	return
}

func (c Config) RemoveUpstream(address Address, target string, port uint16) {

	if _, ok := c[address]; !ok {
		return
	}

	c[address].removeUpstream(target, port)
}

func createUpstreamKey(target string, port uint16) string {
	return fmt.Sprintf("%s:%d", target, port)
}