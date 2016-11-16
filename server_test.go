package rnet

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type PH struct {
	packet []byte
	addr   *Addr
}

func (ph *PH) Receive(pck []byte, addr *Addr) {
	ph.packet = pck
	ph.addr = addr
}

func TestServer(t *testing.T) {
	p := &PH{}
	s, err := New(":5555", p)
	assert.NoError(t, err)
	s2, err := RunNew(":5556", p)
	assert.NoError(t, err)

	addr, err := ResolveAddr(":5556")
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

	s.Close()
	s2.Close()
}

func TestStop(t *testing.T) {
	p := &PH{}
	s, err := RunNew(":5557", p)
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
	p := &PH{}
	s, err := RunNew(":5558", p)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond)
	s.Close()
	time.Sleep(time.Millisecond)
	s2, err := RunNew(":5558", p)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond)
	s2.Close()
}
