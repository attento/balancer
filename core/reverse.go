package core

import (
	"net/http"
	log "github.com/Sirupsen/logrus"
	"net/http/httputil"
)

type Reverser interface  {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var ok bool
	s, ok = InMemoryRepository.Refresh(s)
	if !ok {
		log.WithFields(log.Fields{
			"where":  "reverse",
			"who":  "server",
		}).Info("error server doesn't exist")
		return;
	}
	// IoC get ElectionFunc
	upstream, err := s.Elect(RoundRobin)
	if nil != err {
		log.WithFields(log.Fields{
			"where":  "reverse",
			"who":  "backend",
		}).Error("error electing upstream", err)

		return
	}

	// @todo Raise Event NewRequestProxied(r, upstream)
	doReverseProxy(upstream, w, req)
	// @todo Raise Event NewProxiedResponse(w)
	// because we could Elect another Upstream if the response contain is on error
	// ... with timed blacklist
}

func doReverseProxy(u *Upstream, w http.ResponseWriter, req *http.Request) {
	log.WithFields(log.Fields{
		"req":  req,
		"upstream": u,
	}).Info("request", req.RequestURI,req.Host, req.Proto)

	scheme := req.URL.Scheme
	url, err := u.toUrl(scheme)
	if nil != err {
		log.WithFields(log.Fields{
			"where":  "reverse",
			"who":  "parsing",
		}).Error("Impossible to parse", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(url)

	proxy.ServeHTTP(w, req)
}