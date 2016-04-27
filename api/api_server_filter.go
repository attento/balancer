package api

import (
	"github.com/gin-gonic/gin"
	"github.com/attento/balancer/core"
	"net/http"
	"fmt"
)

func apiServerGetFilter(c *gin.Context) {

	var adr core.Address
	var sent bool
	if adr, sent = getAddressFromParam(c); sent {
		return
	}

	cnf := core.InMemoryRepository.Get()

	var srv *core.Server
	var ok  bool

	if srv, ok = cnf.Server(adr); !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Server %s not found.", adr),
		})
		return
	}

	c.JSON(http.StatusOK, srv.Filter())
}

func apiServerPostFilter(c *gin.Context) {

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

	var filter core.Filter
	var err error
	if err = c.BindJSON(&filter); err == nil {
		core.InMemoryRepository.PutFilter(adr, filter)
		c.Data(http.StatusNoContent, gin.MIMEJSON, nil)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": fmt.Sprintf("Server %s not found.", err),
	})
}