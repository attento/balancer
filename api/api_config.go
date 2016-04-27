package api

import (
	"github.com/gin-gonic/gin"
	"github.com/attento/balancer/core"
	"net/http"
)

type ApiV1Config struct {
	Servers []ApiV1server `json:"servers"`
}

func ApiV1ConfigNew(c core.Config) ApiV1Config {

	servs := make([]ApiV1server, len(c))
	cnf := ApiV1Config{}
	i := 0
	for _, k := range c {
		servs[i] = ApiV1serverNew(k)
		i++
	}

	cnf.Servers = servs

	return cnf
}

func apiConfigGet(c *gin.Context) {

	cnf := core.InMemoryRepository.Get()
	cc := ApiV1ConfigNew(cnf)

	c.JSON(http.StatusOK, cc)
}
