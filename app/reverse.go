package app

import (
	"net/http"
	"net/http/httputil"

	log "github.com/Sirupsen/logrus"

	"github.com/attento/balancer/app/core"
)

type Reverser interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

type reverse struct {
	a    core.Address
	repo core.ConfigRepository
}

func NewReverse(a core.Address, r core.ConfigRepository) *reverse {
	return &reverse{a, r}
}

// @todo Raise Event NewRequestProxied(r, upstream)
// @todo Raise Event NewProxiedResponse(w)
// because we could Elect another Upstream if the response contain is on error
// ... with timed blacklist
func (r *reverse) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	s, ok := r.getFreshConfigService()

	if !ok {
		log.WithFields(log.Fields{
			"where": "reverse",
			"who":   "server",
		}).Info("error server doesn't exist")
		return
	}
	// IoC get ElectionFunc
	upstream, err := s.Elect(core.RoundRobin)
	if nil != err {
		log.WithFields(log.Fields{
			"where": "reverse",
			"who":   "backend",
		}).Error("error electing upstream", err)

		return
	}
	doReverseProxy(upstream, w, req)
}

func doReverseProxy(u *core.Upstream, w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"req":      req,
		"upstream": u,
	}).Info("request:", req.RequestURI, req.Host, req.Proto)

	url, err := u.ToUrl(req.URL.Scheme)
	if nil != err {
		log.WithFields(log.Fields{
			"where": "reverse",
			"who":   "parsing",
		}).Error("Impossible to parse", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ServeHTTP(w, req)
}

func (r *reverse) getFreshConfigService() (s *core.Server, ok bool) {
	return r.repo.Server(r.a)
}
