package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"fmt"
	"github.com/attento/balancer/app/core"
)

func (a *Api) apiServerGetUpstream(c *gin.Context) {

	var adr core.Address
	var us *core.Upstream
	var sent,ok bool
	var srv *core.Server

	adr, sent = getAddressFromParam(c)
	if sent {
		return
	}

	var target string
	var port uint16
	if target, port, sent = getTargetPortFromParam(c); sent {
		return
	}

	srv, sent = a.getConfigServerOr404(c, adr)
	if sent {
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



func (a *Api) apiServerPostUpstream(c *gin.Context) {

	var adr core.Address
	var sent bool
	var srv *core.Server
	var bdUps *core.Upstream
	var err error

	if adr, sent = getAddressFromParam(c); sent {
		return
	}
	srv, sent = a.getConfigServerOr404(c, adr)
	if sent {
		return
	}


	if err = c.BindJSON(&bdUps); err == nil {
		a.app.AddConfigUpstream(srv.Address(), bdUps)
		c.Data(http.StatusNoContent, gin.MIMEJSON, nil)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": fmt.Sprintf("Json Format: %s.", err),
	})
}
