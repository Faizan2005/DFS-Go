package peer2peer

import "net"

type Message struct {
	From    net.Addr
	Payload []byte
}
