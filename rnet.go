package rnet

import (
	"fmt"
	"net"
)

// Addr wraps net.UDPAddr and adds additional useful methods
type Addr struct {
	*net.UDPAddr
}

func (a *Addr) String() string {
	if a == nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", a.IP.String(), a.Port)
}

func ResolveAddr(addr string) (*Addr, error) {
	udp, err := net.ResolveUDPAddr("udp", addr)
	return &Addr{udp}, err
}

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
