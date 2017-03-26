package rnet

import (
	"crypto/rand"
	"fmt"
	"net"
)

// Addr wraps net.UDPAddr and adds additional useful methods
type Addr struct {
	*net.UDPAddr
	Err error
}

// NewAddr creates an address from the primatives in net.UDPAddr
func NewAddr(ip []byte, port int, zone string) *Addr {
	return &Addr{
		UDPAddr: &net.UDPAddr{
			IP:   ip,
			Port: port,
			Zone: zone,
		},
	}
}

// String returns the address as IP:Port
func (a *Addr) String() string {
	if a == nil || a.UDPAddr == nil {
		return ""
	}
	return fmt.Sprintf("%s%s", a.IP, a.Port())
}

// Port returns the port of an address
func (a *Addr) Port() Port {
	if a == nil || a.UDPAddr == nil {
		return Port(0)
	}
	return Port(a.UDPAddr.Port)
}

// ResolveAddr takes a string and returns an Addr
func ResolveAddr(addr string) (*Addr, error) {
	udp, err := net.ResolveUDPAddr("udp", addr)
	return &Addr{udp, err}, err
}

// GetLocalIPs returns all local IP addresses that are not loopback addresses
func GetLocalIPs() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				addr := ipnet.IP.String()
				if addr != "0.0.0.0" {
					ips = append(ips, addr)
				}
			}
		} else if ipaddr, ok := a.(*net.IPAddr); ok && !ipaddr.IP.IsLoopback() {
			if ipaddr.IP.To4() != nil {
				addr := ipaddr.IP.String()
				if addr != "0.0.0.0" {
					ips = append(ips, addr)
				}
			}
		}
	}
	return ips
}

// Port is a convenience as sometimes a port needs to be a number and sometimes
// it needs to be a string.
type Port uint16

// String return the port as string starting with :
func (p Port) String() string { return fmt.Sprintf(":%d", p) }

// RawStr return the port as string
func (p Port) RawStr() string { return fmt.Sprintf("%d", p) }

// On returns a reference to the port on the given ip as an address
func (p Port) On(ip string) *Addr {
	a, _ := ResolveAddr(fmt.Sprintf("%s:%d", ip, p))
	return a
}

// Addr returns the port as an *Addr
func (p Port) Addr() *Addr {
	a, _ := ResolveAddr(fmt.Sprintf(":%d", p))
	return a
}

// Local returns the port as an *Addr on 127.0.0.1
func (p Port) Local() *Addr {
	a, _ := ResolveAddr(fmt.Sprintf("127.0.0.1:%d", p))
	return a
}

// RandomPort picks a random port number between 1000 and 65534 (inclusive)
func RandomPort() Port {
	var p uint16
	b := make([]byte, 2)
	for p < 1000 {
		rand.Read(b)
		p = uint16(b[0]) + uint16(b[1])<<8
	}
	return Port(p)
}
