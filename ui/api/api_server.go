package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/attento/balancer/app/core"
	"github.com/gin-gonic/gin"
)

type ApiV1server struct {
	Address   core.Address     `json:"address"`
	Filter    core.Filter      `json:"filter"`
	Upstreams []*core.Upstream `json:"upstreams"`
}

func ApiV1serverNew(srv *core.Server) ApiV1server {

	return ApiV1server{
		srv.Address,
		srv.Filter,
		core.ConvertConfigUpstreamFromMap(srv.Upstreams),
	}
}

func (a *Api) apiServerGet(c *gin.Context) {

	var adr core.Address
	var sent bool
	if adr, sent = getAddressFromParam(c); sent {
		return
	}

	srv, sent := a.getConfigServerOr404(c, adr)
	if !sent {
		c.JSON(http.StatusOK, ApiV1serverNew(srv))
	}
}

func (a *Api) apiServerPost(c *gin.Context) {

	var v1Server ApiV1server
	var err error
	if err = c.BindJSON(&v1Server); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("Json Format: %s.", err),
		})
	}

	a.app.StartHttpServer(v1Server.Address, v1Server.Filter, v1Server.Upstreams)
	c.Data(http.StatusNoContent, gin.MIMEJSON, nil)
	return
}

func (a *Api) apiServerDelete(c *gin.Context) {

	var adr core.Address
	var sent bool
	if adr, sent = getAddressFromParam(c); sent {
		return
	}
	srv, sent := a.getConfigServerOr404(c, adr)
	if !sent {
		a.app.StopHttpServer(srv.Address, a.DrainingTime())
		c.Data(http.StatusNoContent, gin.MIMEJSON, nil)
	}

	return
}

func getTargetPortFromParam(c *gin.Context) (target string, port uint16, sent bool) {

	upstream := c.Param("upstream")
	if "" == upstream {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "expecting upstream",
		})
		return target, port, true
	}

	r, err := regexp.Compile(`^(?P<Target>.+)-(?P<Port>\d+)$`) //65535
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprint("Error with \"upstream-port\"", upstream, err),
		})
		return target, port, true
	}
	matches := r.FindStringSubmatch(upstream)
	if len(matches) != 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprint("Error with \"upstream-port\"", upstream, matches),
		})
		return target, port, true
	}
	var port64 uint64
	port64, err = strconv.ParseUint(matches[2], 10, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error converting port to uint16 " + matches[2],
		})
		return target, port, true
	}

	return matches[1], uint16(port64), false
}

func (a *Api) getConfigServerOr404(c *gin.Context, adr core.Address) (srv *core.Server, sent bool) {
	sent = true
	srv, ok, err := a.app.ConfigServer(adr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Error on %s", adr),
		})
		log.Error("error 500 on", err, adr)
		return
	}

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Server %s not found.", adr),
		})
		return
	}
	sent = false
	return
}

func getAddressFromParam(c *gin.Context) (adr core.Address, sent bool) {

	adr = core.Address(c.Param("address"))
	if "" == adr {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "expecting address",
		})
		return adr, true
	}

	return adr, false
}
