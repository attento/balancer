package core

import (
	"net/http"
	"time"
	"gopkg.in/tylerb/graceful.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

var (
	listeners = Listeners{}
)

type Listeners map[Address]*graceful.Server

// Event
func (l Listeners) onConfigCreatedNewAddress(a Address, s *Server) {

	rr := mux.NewRouter()
	var r *mux.Route

	for _, v := range s.Filter().Schemes {
		if r != nil {
			r = r.Host(v)
			continue
		}
		r = rr.Schemes(v)
	}

	for _, v := range s.Filter().Hosts {
		if r != nil {
			r = r.Host(v)
			continue
		}
		r = rr.Host(v)
	}

	if "" != s.Filter().PathPrefix {
		if r != nil {
			r = r.PathPrefix(s.Filter().PathPrefix)
		} else {
			r = rr.PathPrefix(s.Filter().PathPrefix)
		}
	}

	if r != nil {
		r = rr.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return true
		})
	}

	r.HandlerFunc(s.ServeHTTP)

	l[a] = &graceful.Server{

		Server: &http.Server{
			Addr: string(a),
			Handler: r.GetHandler(),
		},
	}
	err := l[a].ListenAndServe()
	if err != nil {
		log.Error(err)
	}
}

// onConfigUpdateRemovedAddress
func (l Listeners) onConfigUpdateRemovedAddress(a Address) {
	log.Info("Dry Running ... stopping in 5 seconds")
	l[a].Stop(5 * time.Second)
}