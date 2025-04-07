package peer2peer

import (
	"log"
	"net"
)

type TCPPeer struct {
	net.Conn
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	Listener net.Listener
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.Listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		log.Printf("Failed to listen on %s: %v", t.ListenAddr, err)
		return err
	}

	go t.loopAndAccept()

	return nil
}

func (t *TCPTransport) loopAndAccept() error {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			log.Printf("Error: %+v\n", err)
			continue
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) error {
	peer := NewTCPPeer(conn, outbound)

	err := t.HandshakeFunc(peer)
	if err != nil {
		log.Printf("Handshake error: %v", err)
	}

	for {
		msg := &Message{}
		msg.From = conn.RemoteAddr()

		err = t.Decoder.Decode(conn, msg)
		if err != nil {
			log.Printf("Error: %+v\n", err)
			continue
		}
		log.Printf("message: %+v\n", msg)
	}
}
