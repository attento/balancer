package repository

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/attento/balancer/app/core"
)

func TestInMemoryRepositoryShouldPutAndGet(t *testing.T) {
	assertFilter := core.Filter{[]string{"www.google.com"}, [2]string{"http"}, ""}

	cnf := core.Create()
	cnf.PutFilterProperties(":8080", []string{"www.google.com"}, [2]string{"http"}, "")

	repo := NewInMemoryConfigRepository()
	repo.Put(cnf)
	cnfg := repo.Get()
	s, _ := cnfg.Server(":8080")
	assert.Exactly(t, s.Filter(), assertFilter)
}