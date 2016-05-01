package app

import (
	"github.com/gorilla/mux"
	"github.com/attento/balancer/app/core"
	log "github.com/Sirupsen/logrus"
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
	rev, err := d.httpServers.createProxyIfNotExists(s.Address(), d.repo)
	if err != nil {
		return err
	}
	addFiltersOnRouter(rr, s, rev)

	return d.httpServers.ListenAndServe(s.Address(), rr)
}

func (d *daemon) onConfigServerChangedFilter(s interface{}) error {
	return d.doOnConfigServerChangedFilter(s.(*core.Server))
}

func (d *daemon) doOnConfigServerChangedFilter(s *core.Server) error {

	rr := mux.NewRouter()

	rev, err := d.httpServers.createProxyIfNotExists(s.Address(), d.repo)
	if err != nil {
		return err
	}

	addFiltersOnRouter(rr, s, rev)
	log.Info("changed routes")

	return d.httpServers.ChangeRoutes(s.Address(), rr)
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
