package app

import (
	"github.com/gorilla/mux"
	"github.com/attento/balancer/app/core"
)

func (d *daemon) onNewConfigServer(s interface{}) error {
	return d.doOnNewConfigServer(s.(*core.Server))
}

func (d *daemon) onHttpServerStopped(a interface{}) error {
	return d.doOnHttpServerStopped(a.(core.Address))
}

func (d *daemon) onHttpServerStoppedWithError(a interface{}) error {
	return d.doOnHttpServerStopped(a.(core.Address))
}

func (d *daemon) doOnNewConfigServer(s *core.Server) error {

	rr := mux.NewRouter()
	reverser := NewReverse(s.Address(), d.repo)
	addFiltersOnRouter(rr, s, reverser)

	return d.httpServers.ListenAndServe(s.Address(), rr)
}

func addFiltersOnRouter(rr *mux.Router, s *core.Server, r Reverser) {

	var rs *mux.Router

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

func (d *daemon) doOnHttpServerStopped(a core.Address) error {

	d.repo.RemoveServer(a)

	return nil
}

func (d *daemon) doOnHttpServerStoppedWithError(a core.Address) error {

	// @todo disable server
	//d.repo.DisableServer(a)
	d.repo.RemoveServer(a)

	return nil
}

//
//// onConfigUpdateRemovedAddress
//func (l Listeners) onConfigUpdateRemovedAddress(a core.Address) {
//	log.Info("Dry Running ... stopping in ", time.Second, a)
//	d.httpServers.Stop(a, 1 * time.Second)
//	d.httpServers.events.Raise(core.EventListenerStopped, a)
//}
