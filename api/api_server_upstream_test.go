package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/liuggio/balancer/core"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"net/http"
	"bytes"
)


func TestOnUpstreamShouldResponse200(t *testing.T) {

	r := gin.New()
	routes(r)

	cnf := core.Create()
	cnf.AddUpstreamProperty(":8080", "127.0.0.1", 80, 1, 2)

	core.InMemoryRepository.Put(cnf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/:8080/upstream/127.0.0.1-80", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "{\"Target\":\"127.0.0.1\",\"Port\":80,\"Priority\":1,\"Weight\":2}\n")
}


func TestOnUpstreamPostShouldResponse204(t *testing.T) {

	r := gin.New()
	routes(r)

	cnf := core.Create()
	cnf.AddUpstreamProperty(":8686", "127.0.0.1", 80, 1, 2)
	core.InMemoryRepository.Put(cnf)

	w := httptest.NewRecorder()

	var jsonStr = []byte(`{"Target":"127.0.0.1","Port":80,"Priority":1,"Weight":2}`)
	req, _ := http.NewRequest("POST", "/server/:8686/upstream", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 204)
	assert.Equal(t, w.Body.String(), "null\n")
}

func TestOnUpstreamShouldResponse404(t *testing.T) {

	r := gin.New()
	routes(r)

	cnf := core.Create()
	cnf.AddUpstreamProperty(":8989", "127.0.0.1", 80, 1, 2)

	core.InMemoryRepository.Put(cnf)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/server/:8989/upstream/127.0.0.1-81", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 404)
}