package rnet

import (
	"net"
	"time"
)

const MaxUdpPacketLength = 65507

type Server struct {
	conn          *net.UDPConn
	packetHandler PacketHandler
	localIP       string
	stop, running bool
}

type PacketHandler interface {
	Receive([]byte, *Addr)
}

// New creates a Server
// passing in ":0" for port will select any open port
func New(port string, packetHandler PacketHandler) (*Server, error) {
	laddr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		return nil, err
	}

	server := &Server{
		conn:          conn,
		packetHandler: packetHandler,
		localIP:       getLocalIP(),
		stop:          false,
		running:       false,
	}
	return server, nil
}

// RunNew is a wrapper around new that also calls Run in a Go routine if the
// server was created without error
func RunNew(port string, packetHandler PacketHandler) (*Server, error) {
	s, err := New(port, packetHandler)
	if err == nil {
		go s.Run()
	}
	return s, err
}

// Run is the servers listen loop. When it receives a message it will pass that
// message into the packetHandler
func (s *Server) Run() {
	if s.running || s.conn == nil {
		return
	}
	s.running = true
	buf := make([]byte, MaxUdpPacketLength)
	for {
		l, addr, err := s.conn.ReadFromUDP(buf)
		if s.stop {
			break
		}
		if err == nil {
			packet := make([]byte, l)
			copy(packet, buf[:l])
			go s.packetHandler.Receive(packet, &Addr{addr})
		}
	}
	s.running = false
}

// Returns true if the server is running and can receive messages
// Even if the server is not running, it can still send
func (s *Server) IsRunning() bool { return s.running }

// Returns true if the connection is open
// If the server is closed, it can neither send nor receive
func (s *Server) IsOpen() bool { return s.conn != nil }

// Stop will stop the server
func (s *Server) Stop() error {
	s.stop = true
	return s.conn.SetReadDeadline(time.Now()) // kill all reads
}

// Close will close the connection, freeing the port
func (s *Server) Close() error {
	err := s.Stop()
	if err != nil {
		return err
	}
	if err = s.conn.Close(); err != nil {
		s.conn = nil
	}
	return err
}

// Send will send a single packe (byte slice) to an address
// just a wrapper around WriteToUDP
func (s *Server) Send(packet []byte, addr *Addr) (int, error) {
	return s.conn.WriteToUDP(packet, addr.UDPAddr)
}

// SendAll sends a slice of packets (byte slices) to an address
// this will return the last error it encoutered, if it encountered any
func (s *Server) SendAll(packets [][]byte, addr *Addr) error {
	var last error
	for _, p := range packets {
		if _, e := s.Send(p, addr); e != nil {
			last = e
		}
		time.Sleep(time.Millisecond)
	}
	return last
}

// LocalIP is a getter for the localIP, which is set when the server starts
func (s *Server) LocalIP() string { return s.localIP }

// getLocalIp loops over the interface addresses to find one that is not a loopback
// address and uses that as it's local IP. It may not be fool proof and requires
// further investigation, but it does seem to work.
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, a := range addrs {
		var ip *net.IP
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			ip = &ipnet.IP
		} else if ipaddr, ok := a.(*net.IPAddr); ok && !ipaddr.IP.IsLoopback() {
			ip = &ipaddr.IP
		}
		if ip != nil {
			if ip.To4() != nil {
				addr := ip.String()
				if addr != "0.0.0.0" {
					return addr
				}
			}
		}
	}
	return ""
}
