package rnet

import (
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
)

func TestAddrString(t *testing.T) {
	var a *Addr
	assert.Equal(t, "", a.String())

	addrStr := "127.0.0.1:1234"
	u, err := net.ResolveUDPAddr("udp", addrStr)
	assert.NoError(t, err)
	a = &Addr{u}
	assert.Equal(t, addrStr, a.String())
}
