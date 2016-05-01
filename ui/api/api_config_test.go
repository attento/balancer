package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/attento/balancer/app/core"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"net/http"
	"time"
	"github.com/attento/balancer/app/repository"
	"github.com/attento/balancer/app"
)

func createRepoAppAndRoutes() (*Api, app.DaemonInterface, *gin.Engine, core.ConfigRepository) {
	repo := repository.NewInMemoryConfigRepository()
	d := app.NewStandardRepo(repo)
	a := NewWithApp(d, 1*time.Second)
	r := gin.New()

	return a,d,r,repo
}

func TestConfigShouldResponse200(t *testing.T) {

	a,_,r,repo := createRepoAppAndRoutes()
	configRoutes(r, a)

	repo.AddUpstream(":8080", &core.Upstream{"127.0.0.1", 8080, 1, 2})


	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "{\"servers\":[{\"address\":\":8080\",\"filter\":{\"Hosts\":null,\"Schemes\":[\"\",\"\"],\"PathPrefix\":\"\"},\"upstreams\":[{\"Target\":\"127.0.0.1\",\"Port\":8080,\"Priority\":1,\"Weight\":2}]}]}\n")
}

func TestConfigShouldNeverResponse404(t *testing.T) {


	a,_,r,_ := createRepoAppAndRoutes()
	configRoutes(r, a)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 200)
}