package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/gorilla/mux"
	"net/http/httptest"
	"net/http"
)

func TestShouldCreateAServer(t *testing.T) {

	s :=  &Server{}
	s.addUpstreamProperty("127.0.0.1", 8080, 1, 2)

	go listeners.onConfigCreatedNewAddress(":8000", s)

	time.Sleep(3*time.Second)

	assert.False(t, listeners[":8000"].Interrupted, listeners)
	listeners[":8000"].Stop(1*time.Second)
	time.Sleep(4*time.Second)
	assert.True(t, listeners[":8000"].Interrupted, listeners)
}

type o struct {}

func (s *o) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("YES"))
}


func TestNoFilterForServer(t *testing.T) {

	s :=  &Server{}
	s.addUpstreamProperty("127.0.0.1", 8080, 1, 2)
	s.putFilter(Filter{})

	req, err := http.NewRequest("GET", "/api/123123?jasidjkdsajk", nil)
	req.Host = "www.website.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req, 200, "YES")

	req2, err := http.NewRequest("GET", "/api/", nil)
	req2.Host = "www.liuggio.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req2, 200, "YES")

	req3, err := http.NewRequest("GET", "/123123123123132", nil)
	req3.Host = "www.123123123.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req3, 200, "YES")
}

func TestRouteHostForServer(t *testing.T) {

	s :=  &Server{}
	s.addUpstreamProperty("127.0.0.1", 8080, 1, 2)
	s.putFilter(Filter{Hosts: []string{"www.website.com"}})

	req, err := http.NewRequest("GET", "/web1", nil)
	req.Host = "www.website.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req, 200, "YES")

	req2, err := http.NewRequest("GET", "http://www.website.com/123123123123132", nil)
	assert.Nil(t, err, "Creating req failed!")

	assertServerRequest(t, s, req2, 200, "YES")

	req3, err := http.NewRequest("GET", "/123123123123132", nil)

	assert.Nil(t, err, "Creating req failed!")

	assertServerRequest(t, s, req3, 404, "404 page not found\n")
}


func TestRouteMultipleHostForServer(t *testing.T) {

	s :=  &Server{}
	s.addUpstreamProperty("127.0.0.1", 8080, 1, 2)
	s.putFilter(Filter{Hosts: []string{"www.website.com", "www.liuggio.com"}})

	req, err := http.NewRequest("GET", "/web1", nil)
	req.Host = "www.website.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req, 200, "YES")

	req2, err := http.NewRequest("GET", "/123123123123132", nil)
	req2.Host = "www.liuggio.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req2, 200, "YES")

	req3, err := http.NewRequest("GET", "/123123123123132", nil)
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req3, 404, "404 page not found\n")
}

func TestPrefixForServer(t *testing.T) {

	s :=  &Server{}
	s.addUpstreamProperty("127.0.0.1", 8080, 1, 2)
	s.putFilter(Filter{Hosts: []string{"www.website.com", "www.liuggio.com"}, PathPrefix: "/api"})

	req, err := http.NewRequest("GET", "/api/123123?jasidjkdsajk", nil)
	req.Host = "www.website.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req, 200, "YES")

	req2, err := http.NewRequest("GET", "/api/", nil)
	req2.Host = "www.liuggio.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req2, 200, "YES")

	req3, err := http.NewRequest("GET", "/123123123123132", nil)
	req3.Host = "www.liuggio.com"
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req3, 404, "404 page not found\n")
}

func TestOnlyPathPrefixForServer(t *testing.T) {

	s :=  &Server{}
	s.addUpstreamProperty("127.0.0.1", 8080, 1, 2)
	s.putFilter(Filter{PathPrefix: "/api"})

	req, err := http.NewRequest("GET", "/api/123123?jasidjkdsajk", nil)
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req, 200, "YES")

	req3, err := http.NewRequest("GET", "/123123123123132", nil)
	assert.Nil(t, err, "Creating req failed!")
	assertServerRequest(t, s, req3, 404, "404 page not found\n")
}

func assertServerRequest(t *testing.T, s *Server, req *http.Request, code int, body string) {

	rr := mux.NewRouter()
	a := &o{}
	// Routes consist of a path and a handler function.
	addFiltersOnRouter(rr, s, a)

	w := httptest.NewRecorder()
	rr.ServeHTTP(w, req)

	assert.Exactly(t, code, w.Code)
	assert.Exactly(t, body, w.Body.String())
}