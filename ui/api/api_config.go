package api

import (
	"net/http"

	"github.com/attento/balancer/app/core"
	"github.com/gin-gonic/gin"
)

type ApiV1Config struct {
	Servers []ApiV1server `json:"servers"`
}

func ApiV1ConfigNew(c core.Config) ApiV1Config {

	servs := make([]ApiV1server, len(c.Servers()))
	cnf := ApiV1Config{}
	i := 0
	for _, k := range c.Servers() {
		servs[i] = ApiV1serverNew(k)
		i++
	}

	cnf.Servers = servs

	return cnf
}

func (a *Api) apiConfigGet(c *gin.Context) {
	cnf := *a.app.Config()
	cc := ApiV1ConfigNew(cnf)
	c.JSON(http.StatusOK, cc)
}
