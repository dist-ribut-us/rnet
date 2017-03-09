package rnet

import (
	"crypto/rand"
	"fmt"
	"net"
)

// Addr wraps net.UDPAddr and adds additional useful methods
type Addr struct {
	*net.UDPAddr
}

// String returns the address as IP:Port
func (a *Addr) String() string {
	if a == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", a.IP.String(), a.Port)
}

// ResolveAddr takes a string and returns an Addr
func ResolveAddr(addr string) (*Addr, error) {
	udp, err := net.ResolveUDPAddr("udp", addr)
	return &Addr{udp}, err
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

func (p Port) On(ip string) (*Addr, error) {
	return ResolveAddr(fmt.Sprintf("%s%s", ip, p))
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
