package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/liuggio/balancer/core"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"net/http"
)

func TestShouldResponse200(t *testing.T) {

	r := gin.New()
	routes(r)

	cnf := core.Create()
	cnf.AddUpstreamProperty(":8484", "127.0.0.1", 80, 1, 2)

	core.InMemoryRepository.Put(cnf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/:8484", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "{\"address\":\":8484\",\"filter\":{\"Hosts\":null,\"Schemes\":[\"\",\"\"],\"PathPrefix\":\"\"},\"upstreams\":{\"127.0.0.1:80\":{\"Target\":\"127.0.0.1\",\"Port\":80,\"Priority\":1,\"Weight\":2}}}\n")
}

func TestShouldResponse404(t *testing.T) {

	r := gin.New()
	routes(r)

	cnf := core.Create()
	cnf.AddUpstreamProperty(":8585", "127.0.0.1", 80, 1, 2)

	core.InMemoryRepository.Put(cnf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/123", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 404)
}