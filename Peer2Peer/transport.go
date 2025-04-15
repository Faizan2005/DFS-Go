package peer2peer

import "net"

type Peer interface {
	net.Conn
	// Close() error
	// RemoteAddr() net.Addr
	Send([]byte) error
}

type Transport interface {
	Addr() string
	listenAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
