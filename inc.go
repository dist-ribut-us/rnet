package rnet

import (
	"sync/atomic"
)

// PortIncrementer provides a threadsafe way to return an incrementing port
// value.
type PortIncrementer struct {
	port int32
}

// NewPortIncrementer returns a PortIncrementer. The first value returned will
// be one greater than start.
func NewPortIncrementer(start uint16) *PortIncrementer {
	return &PortIncrementer{
		port: int32(start),
	}
}

// Next returnst the next port
func (p *PortIncrementer) Next() Port {
	return Port(atomic.AddInt32(&(p.port), 1))
}
