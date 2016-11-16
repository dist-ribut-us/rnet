package rnet

import (
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
