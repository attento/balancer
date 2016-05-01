package api

import (
	"github.com/gin-gonic/gin"
	"github.com/attento/balancer/app"
	"time"
	log "github.com/Sirupsen/logrus"
)

type Api struct {
	app app.DaemonInterface
	drainingTime time.Duration
}

func New(d time.Duration) *Api {
	return &Api{app.NewStandard(), d}
}

func NewWithApp(a app.DaemonInterface, d time.Duration) *Api {
	return &Api{a, d}
}


func (a *Api) DrainingTime() time.Duration {
	return a.drainingTime
}


func configRoutes(r *gin.Engine, a *Api) {
	r.GET("/", a.apiConfigGet)
	r.HEAD("/", a.apiConfigGet)
}

func serverRoutes(r *gin.Engine, a *Api) {

	r.POST("/server", a.apiServerPost)
	r.GET("/server/:address", a.apiServerGet)
	r.HEAD("/server/:address", a.apiServerGet)
	r.DELETE("/server/:address", a.apiServerDelete)
}

func serverUpstreamRoutes(r *gin.Engine, a *Api) {
	r.GET("/server/:address/upstream/:upstream", a.apiServerGetUpstream)
	r.HEAD("/server/:address/upstream/:upstream", a.apiServerGetUpstream)
	r.POST("/server/:address/upstream", a.apiServerPostUpstream)
	r.PUT("/server/:address/upstream", a.apiServerPostUpstream)
}

func serverFilterRoutes(r *gin.Engine, a *Api) {
	r.GET("/server/:address/filter", a.apiServerGetFilter)
	r.HEAD("/server/:address/filter", a.apiServerGetFilter)
	r.POST("/server/:address/filter", a.apiServerPostFilter)
	r.PUT("/server/:address/filter", a.apiServerPostFilter)
}

func Run(addr string, isDebug bool) {

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}

	a := New(3*time.Second)

	r := gin.Default()

	configRoutes(r, a)
	serverRoutes(r, a)
	serverUpstreamRoutes(r, a)
	serverFilterRoutes(r, a)
	log.Info("Started API Server on:", addr)
	log.Error(r.Run(addr))
}