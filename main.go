package main

import (
	"log"

	peer2peer "github.com/Faizan2005/DFS-Go/Peer2Peer"
)

func onPeer(peer peer2peer.Peer) error {
	peer.Close()
	return nil
}

func main() {
	opts := peer2peer.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: peer2peer.NOPEHandshakeFunc,
		Decoder:       peer2peer.DefaultDecoder{},
		OnPeer:        onPeer, //func(peer2peer.Peer) error { return fmt.Errorf("Failed the OnPeer func") }
	}

	tsp := peer2peer.NewTCPTransport(opts)
	err := tsp.ListenAndAccept()
	if err != nil {
		log.Printf("Error: %+v\n", err)
	}

	go func() {
		for {
			msg := <-tsp.Consume()
			log.Printf("Message: %+v\n", msg)
		}
	}()

	select {}
}
