package main

import (
	"log"

	peer2peer "github.com/Faizan2005/DFS-Go/Peer2Peer"
)

func main() {
	opts := peer2peer.TCPTransportOpts{
		ListenAddr:    ":3000",
		HandshakeFunc: peer2peer.NOPHandshakeFunc,
		Decoder:       peer2peer.DefaultDecoder{},
	}

	tsp := peer2peer.NewTCPTransport(opts)
	err := tsp.ListenAndAccept()
	if err != nil {
		log.Printf("Error: %+v\n", err)
	}

	select {}
}
