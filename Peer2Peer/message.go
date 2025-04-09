package peer2peer

import "net"

type RPC struct {
	From    net.Addr
	Payload []byte
}
