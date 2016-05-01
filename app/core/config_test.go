package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestShouldCreateAServerAndAFilter(t *testing.T) {

	assertFilter := Filter{[]string{"www.google.com"}, [2]string{"http"}, ""}

	cnf := Create()
	cnf.PutFilterProperties(":8080", []string{"www.google.com"}, [2]string{"http"}, "")

	assert.Exactly(t, cnf.servers[":8080"].Filter(), assertFilter)
}

func TestShouldCreateAndRemoveAServerAndAnUpstream(t *testing.T) {

	upstreamAssertion := &Upstream{"127.0.0.1", 80, 1, 2}
	upstreamAssertion2 := &Upstream{"127.0.0.2", 80, 1, 2}

	cnf := Create()
	cnf.AddUpstreamProperty(":8080", "127.0.0.1", 80, 1, 2)
	cnf.AddUpstreamProperty(":8080", "127.0.0.2", 80, 1, 2)

	u, _ := cnf.servers[":8080"].Upstream("127.0.0.1", 80)
	assert.Exactly(t, u, upstreamAssertion)

	u2, _ := cnf.servers[":8080"].Upstream("127.0.0.2", 80)
	assert.Exactly(t, u2, upstreamAssertion2)

	cnf.RemoveUpstream(":8080", "127.0.0.2", 80)
	_, ok := cnf.servers[":8080"].Upstream("127.0.0.2", 80)
	assert.False(t, ok)
}