package core

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestShouldCreateAServer(t *testing.T) {

	s :=  &Server{}
	s.addUpstreamProperty("127.0.0.1", 80, 1, 2)

	go listeners.onConfigCreatedNewAddress(":8000", s)

	time.Sleep(3*time.Second)

	assert.False(t, listeners[":8000"].Interrupted, listeners)
	listeners[":8000"].Stop(1*time.Second)
	time.Sleep(4*time.Second)
	assert.True(t, listeners[":8000"].Interrupted, listeners)
}