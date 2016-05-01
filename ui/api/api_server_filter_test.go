package api

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/attento/balancer/app/core"
	"net/http/httptest"
	"net/http"
	"bytes"
)

func TestOnFilterShouldResponse204(t *testing.T) {

	a,_,r,repo := createRepoAppAndRoutes()
	serverFilterRoutes(r, a)

	f := core.Filter{Schemes: [2]string{"http"}}
	repo.PutFilter(core.Address(":8383"), f)

	w := httptest.NewRecorder()

	var jsonStr = []byte(`{"Hosts":null,"Schemes":["",""],"PathPrefix":""}`)
	req, _ := http.NewRequest("POST", "/server/:8383/filter", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Code, 204)
	assert.Equal(t, w.Body.String(), "")
}

