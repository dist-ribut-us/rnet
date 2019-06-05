package rnet

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type packetHandler struct {
	packet []byte
	addr   *Addr
}

func (ph *packetHandler) Receive(pck []byte, addr *Addr) {
	ph.packet = pck
	ph.addr = addr
}

func TestServer(t *testing.T) {
	port := Port(5556)
	addr := port.Addr()

	p := &packetHandler{}
	s, err := New(5555, p)
	assert.NoError(t, err)
	s2, err := RunNew(port, p)
	assert.NoError(t, err)

	s.Send([]byte{1, 2, 3}, addr)

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
		if len(p.packet) > 0 {
			break
		}
	}

	if len(p.packet) != 3 || p.packet[0] != 1 || p.packet[1] != 2 || p.packet[2] != 3 {
		t.Error("Incorrect Packet")
	}

	addrStr := p.addr.String()
	if l := len(addrStr); l < 4 || addrStr[l-4:] != "5555" {
		t.Error("Incorrect Address")
	}

	assert.Equal(t, Port(5555), s.GetPort())

	s.Close()
	s2.Close()
}

func TestStop(t *testing.T) {
	p := &packetHandler{}
	s, err := RunNew(5557, p)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond)
	if !s.running {
		t.Error("Server is not running")
	}
	s.Stop()
	time.Sleep(time.Millisecond)
	if s.running {
		t.Error("Server has not stopped")
	}
	s.Close()
}

func TestClose(t *testing.T) {
	p := &packetHandler{}
	s, err := RunNew(5558, p)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond)
	s.Close()
	time.Sleep(time.Millisecond)
	s2, err := RunNew(5558, p)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond)
	s2.Close()
}
