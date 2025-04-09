package peer2peer

type Peer interface {
	Close() error
}

type Transport interface {
	Addr() string
	listenAndAccept() error
	Consume() <-chan RPC
}
