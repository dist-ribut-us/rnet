package rnet

import (
	"net"
	"time"
)

// MaxUDPPacketLength is the max possible length after the UDP headers are
// removed
const MaxUDPPacketLength = 65507

// Server is a UDP server that can be used to send and receive UPD packets
type Server struct {
	conn          *net.UDPConn
	packetHandler PacketHandler
	stop, running bool
	port          Port
}

// PacketHandler is an interface for receiving packets from a UDP server
type PacketHandler interface {
	Receive([]byte, *Addr)
}

// New creates a Server passing in ":0" for port will select any open port. It
// is also possible to specify a full IP address for port, as long as the
// address is local, but generally only a port is specified.
func New(port Port, packetHandler PacketHandler) (*Server, error) {
	laddr, err := net.ResolveUDPAddr("udp", port.String())
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
		stop:          false,
		running:       false,
		port:          port,
	}
	return server, nil
}

// Port returns the port the server is listening on
func (s *Server) Port() Port {
	if s.port == 0 && s.conn != nil {
		addr := s.conn.LocalAddr()
		if udpaddr, ok := addr.(*net.UDPAddr); ok {
			s.port = Port(udpaddr.Port)
		}
	}
	return s.port
}

// RunNew is a wrapper around new that also calls Run in a Go routine if the
// server was created without error
func RunNew(port Port, packetHandler PacketHandler) (*Server, error) {
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
	buf := make([]byte, MaxUDPPacketLength)
	for {
		l, addr, err := s.conn.ReadFromUDP(buf)
		if s.stop {
			break
		}
		if err == nil {
			packet := make([]byte, l)
			copy(packet, buf[:l])
			go s.packetHandler.Receive(packet, &Addr{addr, nil})
		}
	}
	s.running = false
}

// IsRunning returns true if the server is running and can receive messages.
// Even if the server is not running, it can still send.
func (s *Server) IsRunning() bool { return s.running }

// IsOpen returns true if the connection is open. If the server is closed, it
// can neither send nor receive
func (s *Server) IsOpen() bool { return s.conn != nil }

// Stop will stop the server
func (s *Server) Stop() error {
	s.stop = true
	return s.conn.SetReadDeadline(time.Now()) // kill all reads
}

// Close will close the connection, freeing the port
func (s *Server) Close() error {
	if err := s.Stop(); err != nil {
		return err
	}
	if err := s.conn.Close(); err != nil {
		return err
	}
	s.conn = nil
	return nil
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
