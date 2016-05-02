package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/attento/balancer/app/core"
	"github.com/stretchr/testify/assert"
)

func TestOnUpstreamShouldResponse200(t *testing.T) {

	a, _, r, repo := createRepoAppAndRoutes()
	serverUpstreamRoutes(r, a)
	repo.NewServer(":8484")
	repo.AddUpstream(":8484", &core.Upstream{"127.0.0.1", 80, 1, 2})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/:8484/upstream/127.0.0.1-80", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "{\"Target\":\"127.0.0.1\",\"Port\":80,\"Priority\":1,\"Weight\":2}\n")
}

func TestOnUpstreamPostShouldResponse204(t *testing.T) {

	a, _, r, repo := createRepoAppAndRoutes()
	serverUpstreamRoutes(r, a)

	repo.NewServer(":8686")
	w := httptest.NewRecorder()

	var jsonStr = []byte(`{"Target":"127.0.0.1","Port":80,"Priority":1,"Weight":2}`)
	req, _ := http.NewRequest("POST", "/server/:8686/upstream", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 204)
	assert.Equal(t, w.Body.String(), "")
}

func TestOnUpstreamShouldResponse404(t *testing.T) {

	a, _, r, _ := createRepoAppAndRoutes()
	serverUpstreamRoutes(r, a)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/:8989/upstream/127.0.0.1-81", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 404)
}
