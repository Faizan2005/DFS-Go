package main

import (
	"log"
	"time"

	peer2peer "github.com/Faizan2005/DFS-Go/Peer2Peer"
)

func onPeer(peer peer2peer.Peer) error {
	//peer.Close()
	return nil
}

func main() {
	metaPath := "test_metadata.json"

	tcpOpts := peer2peer.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: peer2peer.NOPEHandshakeFunc,
		Decoder:       peer2peer.DefaultDecoder{},
		OnPeer:        onPeer,
	}
	opts := ServerOpts{
		pathTransform: CASPathTransformFunc,
		tcpTransport:  *peer2peer.NewTCPTransport(tcpOpts),
		metaData:      *NewMetadata(metaPath),
		
	}

	s := NewServer(opts)
	go func() {
		if err := s.Run(); err != nil {
			log.Println("Server error:", err)
		}
	}()

	time.Sleep(3 * time.Second)
	s.Stop()
}
