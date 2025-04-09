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
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	Listener net.Listener
	rpcch    chan RPC
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
		rpcch:            make(chan RPC, 1024),
	}
}

func (t *TCPTransport) Addr() string {
	return t.ListenAddr
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (p *TCPPeer) Close() error {
	return p.Conn.Close()
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

func (t *TCPTransport) loopAndAccept() {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			log.Printf("Error: %+v\n", err)
			continue
		}

		go t.handleConn(conn, false)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	peer := NewTCPPeer(conn, outbound)

	defer func() {
		log.Println("Dropping peer connection")
		conn.Close()
	}()

	err := t.HandshakeFunc(peer)
	if err != nil {
		log.Printf("Handshake error: %v", err)
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	for {
		msg := RPC{}
		msg.From = conn.RemoteAddr()

		err = t.Decoder.Decode(conn, &msg)
		if err != nil {
			log.Printf("Decode error from %v: %v", conn.RemoteAddr(), err)
			return
		}

		t.rpcch <- msg
	}
}
