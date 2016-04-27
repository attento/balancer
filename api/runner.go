package api

import (
	"github.com/gin-gonic/gin"
)

func routes(r *gin.Engine)  {

	r.POST("/server", apiServerPost)
	r.GET("/server/:address", apiServerGet)
	r.HEAD("/server/:address", apiServerGet)
	r.DELETE("/server/:address", apiServerDelete)

	r.GET("/server/:address/upstream/:upstream", apiServerGetUpstream)
	r.HEAD("/server/:address/upstream/:upstream", apiServerGetUpstream)
	r.POST("/server/:address/upstream", apiServerPostUpstream)
	r.PUT("/server/:address/upstream", apiServerPostUpstream)

	r.GET("/server/:address/filter", apiServerGetFilter)
	r.HEAD("/server/:address/filter", apiServerGetFilter)
	r.POST("/server/:address/filter", apiServerPostFilter)
	r.PUT("/server/:address/filter", apiServerPostFilter)

	r.GET("/", apiConfigGet)
	r.HEAD("/", apiConfigGet)
}

func Run(addr string, isDebug bool) {

	if !isDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	routes(r)

	r.Run(addr)
}