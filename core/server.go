package core

import (
	"net/http"
	"time"
	"gopkg.in/tylerb/graceful.v1"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"fmt"
)

var (
	listeners = Listeners{}
)

type Listeners map[Address]*graceful.Server

// Event
func (l Listeners) onConfigCreatedNewAddress(a Address, s *Server) {

	rr := mux.NewRouter()
	addFiltersOnRouter(rr, s, s)

	l[a] = &graceful.Server{
		Server: &http.Server{
			Addr: string(a),
			Handler: rr,
		},
	}
	err := l[a].ListenAndServe()
	if err != nil {
		log.Error(err)
	}
	// @todo event
	// server is not binded disable it!!
	// ugly remove from here we are in the domain...
	InMemoryRepository.RemoveServer(a)
}

// onConfigUpdateRemovedAddress
func (l Listeners) onConfigUpdateRemovedAddress(a Address) {
	log.Info("Dry Running ... stopping in ", time.Second, a)
	l[a].Stop(1 * time.Second)
}


func addHosts(rr *mux.Router, s *Server, r Reverser) {

	if len(s.Filter().Hosts) <= 0 {
		rr.MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			return true
		}).HandlerFunc(r.ServeHTTP)
		return
	}

	for _, v := range s.Filter().Hosts {
		 rr.Host(v).HandlerFunc(r.ServeHTTP)
	}
}

func addFiltersOnRouter(rr *mux.Router, s *Server, r Reverser) {

	var rs *mux.Router

	fmt.Errorf("pre-")
	if s.Filter().PathPrefix == "" {
		fmt.Errorf("vuoto-")
	}

	if "" != s.Filter().PathPrefix {
		rs = rr.PathPrefix(s.Filter().PathPrefix).Subrouter()
	}

	if rs != nil {
		addHosts(rs, s, r)
	}else {
		addHosts(rr, s, r)
	}

	return
}