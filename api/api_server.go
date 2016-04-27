package api

import (
	"github.com/gin-gonic/gin"
	"github.com/attento/balancer/core"
	"net/http"
	"fmt"
	"strconv"
	"regexp"
)

type ApiV1server struct {
	Address   core.Address  		     `json:"address"`
	Filter    core.Filter		         `json:"filter"`
	Upstreams  map[string]*core.Upstream `json:"upstreams"`
}

func ApiV1serverNew(srv *core.Server) ApiV1server {
	return ApiV1server{
		srv.Address(),
		srv.Filter(),
		srv.Upstreams(),
	}
}

func apiServerGet(c *gin.Context) {

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

	c.JSON(http.StatusOK, ApiV1serverNew(srv))
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
			"message": fmt.Sprint("Error with \"upstream-port\"",upstream, err),
		})
		return target, port, true
	}
	matches := r.FindStringSubmatch(upstream)
	if len(matches) != 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprint("Error with \"upstream-port\"",upstream, matches),
		})
		return target, port, true
	}
	var port64 uint64
	port64, err = strconv.ParseUint(matches[2], 10, 16)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "error converting port to uint16 "+matches[2],
		})
		return target, port, true
	}

	return matches[1], uint16(port64), false
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