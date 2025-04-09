package peer2peer

type HandshakeFunc func(Peer) error

func NOPEHandshakeFunc(Peer) error { return nil }
