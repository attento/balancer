package app

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/attento/balancer/app/core"
	"github.com/liuggio/events"
	"time"
	"github.com/attento/balancer/app/repository"
	"net"
)

func TestWhenAServerHasBeenCreated(t *testing.T) {

	errorHappened := false
	errfn := func(err events.ListenerError) {errorHappened = true;}

	e := events.NewWithErrListener(errfn)
	e.AddInMemoryEventRepo()

	d := &daemon{
		httpServers: NewHttpServers(e),
		repo: repository.NewInMemoryConfigRepository(),
		e: e,
	}

	upstreams := []*core.Upstream{&core.Upstream{"127.0.0.1", 90, 1, 1}}

	d.StartHttpServer(":8080", core.Filter{}, upstreams)

	assert.True(t, true)
	assert.True(t, e.GetEventRepo().Contains(core.EventConfigServerCreated))
}

func TestFunctionalAServerHasBeenCreatedAndStartStopHttpServer(t *testing.T) {

	errorHappened := false
	errfn := func(err events.ListenerError) {errorHappened = true;}

	e := events.NewWithErrListener(errfn)
	e.AddInMemoryEventRepo()

	d := New(repository.NewInMemoryConfigRepository(), e)

	upstreams := []*core.Upstream{&core.Upstream{"127.0.0.1", 90, 1, 1}}

	start := d.StartHttpServer(":8080", core.Filter{}, upstreams)
	assert.Nil(t, start, start)
	time.Sleep(2*time.Second)
	d.httpServers.Stop(":8080", 1*time.Second)
	e.Wait()
	assert.True(t, true)
	assert.True(t, e.GetEventRepo().Contains(core.EventConfigServerCreated))
	assert.True(t, e.GetEventRepo().Contains(core.EventHttpServerStarted))
	assert.True(t, e.GetEventRepo().Contains(core.EventHttpServerStopped))

	assert.False(t, e.GetEventRepo().Contains(core.EventHttpServerStoppedWithError))
}

func TestFunctionalAServerHasBeenCreatedAndStartOnBindPortHttpServer(t *testing.T) {

	errorHappened := false
	errfn := func(err events.ListenerError) {errorHappened = true;}

	e := events.NewWithErrListener(errfn)
	e.AddInMemoryEventRepo()

	d := New(repository.NewInMemoryConfigRepository(), e)

	upstreams := []*core.Upstream{&core.Upstream{"127.0.0.1", 90, 1, 1}}

	// binding port 8889
	l, err := net.Listen("tcp", ":8080")
	assert.Nil(t, err)

	// asserting Event Stopped execution
	start := d.StartHttpServer(":8889", core.Filter{}, upstreams)
	assert.Nil(t, start, start)
	time.Sleep(2*time.Second)
	l.Close()
	d.httpServers.Stop(":8080", 1*time.Second)
	e.Wait()
	assert.True(t, true)
	assert.True(t, e.GetEventRepo().Contains(core.EventConfigServerCreated))
	assert.True(t, e.GetEventRepo().Contains(core.EventHttpServerStarted))
	assert.True(t, e.GetEventRepo().Contains(core.EventHttpServerStoppedWithError))

}