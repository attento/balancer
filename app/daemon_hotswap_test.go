package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/attento/balancer/app/core"
	"github.com/attento/balancer/app/repository"
	"github.com/liuggio/events"
	"github.com/stretchr/testify/assert"
)

type MyMockedReverser struct {
	status int
}

func (r *MyMockedReverser) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Yes")
	w.WriteHeader(r.status)
}

func TestShouldHotSwapFilter(t *testing.T) {

	errorHappened := false
	errfn := func(err events.ListenerError) { errorHappened = true }

	e := events.NewWithErrListener(errfn)
	e.AddInMemoryEventRepo()

	d := New(repository.NewInMemoryConfigRepository(), e)
	d.httpServers.proxies[":8889"] = &MyMockedReverser{201}
	upstreams := []*core.Upstream{&core.Upstream{"127.0.0.1", 8000, 1, 1}}

	// asserting Event Stopped execution
	start := d.StartHttpServer(":8889", core.Filter{}, upstreams)
	assert.Nil(t, start, start)
	time.Sleep(2 * time.Second)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	d.httpServers.servers[":8889"].Handler.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Exactly(t, w.Body.String(), "Yes")

	d.PutConfigFilter(":8889", core.Filter{PathPrefix: "/api"})
	time.Sleep(2 * time.Second)

	req2, _ := http.NewRequest("GET", "/", nil)
	w404 := httptest.NewRecorder()
	d.httpServers.servers[":8889"].Handler.ServeHTTP(w404, req2)

	assert.Equal(t, w404.Code, 404)
	assert.Equal(t, w404.Body.String(), "404 page not found\n")

	req3, _ := http.NewRequest("GET", "/api", nil)
	w200 := httptest.NewRecorder()
	d.httpServers.servers[":8889"].Handler.ServeHTTP(w200, req3)

	assert.Equal(t, w200.Code, 200)
	assert.Exactly(t, w.Body.String(), "Yes")

	d.httpServers.Stop(":8889", 1*time.Second)
	e.Wait()
}
