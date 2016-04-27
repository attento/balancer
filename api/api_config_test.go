package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/liuggio/balancer/core"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"net/http"
)

func TestConfigShouldResponse200(t *testing.T) {

	r := gin.New()
	routes(r)

	cnf := core.Create()
	cnf.AddUpstreamProperty(":8080", "127.0.0.1", 8080, 1, 2)

	core.InMemoryRepository.Put(cnf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "{\"servers\":[{\"address\":\":8080\",\"filter\":{\"Hosts\":null,\"Schemes\":[\"\",\"\"],\"PathPrefix\":\"\"},\"upstreams\":{\"127.0.0.1:8080\":{\"Target\":\"127.0.0.1\",\"Port\":8080,\"Priority\":1,\"Weight\":2}}}]}\n")
}

func TestConfigShouldNeverResponse404(t *testing.T) {
	r := gin.New()
	routes(r)
	cnf := core.Create()
	core.InMemoryRepository.Put(cnf)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 200)
}