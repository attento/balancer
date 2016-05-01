package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/attento/balancer/app/core"
	"net/http/httptest"
	"net/http"
 	"bytes"
)

func TestGetServerShouldResponse200(t *testing.T) {

	a,_,r,repo := createRepoAppAndRoutes()
	serverRoutes(r, a)
	repo.NewServer(":8484")
	repo.AddUpstream(":8484", &core.Upstream{"127.0.0.1", 80, 1, 2})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/:8484", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "{\"address\":\":8484\",\"filter\":{\"Hosts\":null,\"Schemes\":[\"\",\"\"],\"PathPrefix\":\"\"},\"upstreams\":[{\"Target\":\"127.0.0.1\",\"Port\":80,\"Priority\":1,\"Weight\":2}]}\n")
}

func TestShouldResponse404(t *testing.T) {

	a,_,r,repo := createRepoAppAndRoutes()
	serverRoutes(r, a)
	repo.NewServer(":8484")
	repo.AddUpstream(":8484", &core.Upstream{"127.0.0.1", 80, 1, 2})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 404)
}


func TestOnPostShouldResponse204(t *testing.T) {

	a,_,r,repo := createRepoAppAndRoutes()
	serverRoutes(r, a)

	var jsonStr = []byte(`{"address":":8484","filter":{"Hosts":["www.website.com"],"Schemes":["",""],"PathPrefix":""},"upstreams":[{"Target":"127.0.0.1","Port":80,"Priority":1,"Weight":2}]}`)
	req, _ := http.NewRequest("POST", "/server", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 204)
	assert.Equal(t, w.Body.String(), "")

	cnf := repo.Get()
	srv, ok := cnf.Server(":8484")
	assert.True(t, ok)
	assert.Exactly(t, srv.Filter().Hosts, []string{"www.website.com"})
}