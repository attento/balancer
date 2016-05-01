package api

import (
	"github.com/gin-gonic/gin"
	"github.com/attento/balancer/app/core"
	"net/http"
	"fmt"
)

func (a *Api) apiServerGetFilter(c *gin.Context) {

	var adr core.Address
	var sent bool
	if adr, sent = getAddressFromParam(c); sent {
		return
	}

	srv, sent := a.getConfigServerOr404(c, adr)
	if !sent {
		c.JSON(http.StatusOK, srv.Filter())
	}
}

func (a *Api) apiServerPostFilter(c *gin.Context) {

	var adr core.Address
	var sent bool
	if adr, sent = getAddressFromParam(c); sent {
		return
	}

	_, sent = a.getConfigServerOr404(c, adr)
	if sent {
		return
	}

	var filter core.Filter
	var err error
	if err = c.BindJSON(&filter); err == nil {
		a.app.PutConfigFilter(adr, filter)
		c.Data(http.StatusNoContent, gin.MIMEJSON, nil)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": fmt.Sprintf("Server %s not found.", err),
	})
}