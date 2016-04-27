package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryRepositoryShouldPutAndGet(t *testing.T) {
	assertFilter := Filter{[]string{"www.google.com"}, [2]string{"http"}, ""}

	cnf := Create()
	cnf.PutFilter(":80", []string{"www.google.com"}, [2]string{"http"}, "")

	InMemoryRepository.Put(cnf)
	cnfg := InMemoryRepository.Get()
	assert.Exactly(t, cnfg[":80"].Filter(), assertFilter)
}