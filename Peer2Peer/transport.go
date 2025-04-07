package peer2peer

import "net"

type Peer interface {
	net.Conn
}

type Transport interface {
	Addr() string
	listenAndAccept() error
}
