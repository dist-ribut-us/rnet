package rnet

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddrString(t *testing.T) {
	var a *Addr
	assert.Equal(t, "", a.String())

	addrStr := "127.0.0.1:1234"
	a, err := ResolveAddr(addrStr)
	assert.NoError(t, err)
	assert.Equal(t, addrStr, a.String())
}

func TestIncrementer(t *testing.T) {
	pi := NewPortIncrementer(5555)
	assert.Equal(t, Port(5556), pi.Next())
	assert.Equal(t, Port(5557), pi.Next())
}
