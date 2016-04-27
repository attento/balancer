package api

import (
	"github.com/gin-gonic/gin"
	"github.com/attento/balancer/core"
	"net/http"
	"fmt"
)

func apiServerPostUpstream(c *gin.Context) {

	var adr core.Address
	var sent bool
	if adr, sent = getAddressFromParam(c); sent {
		return
	}

	cnf := core.InMemoryRepository.Get()
	var ok  bool
	if _, ok = cnf.Server(adr); !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Server %s not found.", adr),
		})
		return
	}

	var bdUps core.Upstream
	var err error

	if err = c.BindJSON(&bdUps); err == nil {
		core.InMemoryRepository.AddUpstream(adr, &bdUps)
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": fmt.Sprintf("Json Format: %s.", err),
	})
}

func apiServerGetUpstream(c *gin.Context) {

	var adr core.Address
	var sent bool
	if adr, sent = getAddressFromParam(c); sent {
		return
	}

	var target string
	var port uint16
	if target, port, sent = getTargetPortFromParam(c); sent {
		return
	}

	cnf := core.InMemoryRepository.Get()
	var srv *core.Server
	var ok  bool
	var us *core.Upstream

	if srv, ok = cnf.Server(adr); !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Server %s not found.", adr),
		})
		return
	}

	if us, ok = srv.Upstream(target, port); !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Upstream %s-%d not foun on %s Server.", target, port, adr),
		})
		return
	}

	c.JSON(http.StatusOK, us)
}