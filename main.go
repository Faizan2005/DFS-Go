package main

import (
	"bytes"
	"log"
	"time"

	//"time"

	peer2peer "github.com/Faizan2005/DFS-Go/Peer2Peer"
)

func main() {
	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")

	go func() {
		if err := s1.Run(); err != nil {
			log.Println("Server s1 error:", err)
		}
	}()
	time.Sleep(1 * time.Second)

	go s2.Run()
	time.Sleep(1 * time.Second)

	data := bytes.NewReader([]byte("This is some data"))
	s2.StoreData("this is a secret key", data)

	select {}
}

func makeServer(listenAddr string, node ...string) *Server {
	metaPath := "test_metadata.json"

	tcpOpts := peer2peer.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: peer2peer.NOPEHandshakeFunc,
		Decoder:       peer2peer.DefaultDecoder{},
	}
	tcpTransport := peer2peer.NewTCPTransport(tcpOpts)

	s := &Server{} // create server first to use its OnPeer
	tcpTransport.OnPeer = s.OnPeer

	opts := ServerOpts{
		pathTransform:  CASPathTransformFunc,
		tcpTransport:   *tcpTransport,
		metaData:       *NewMetadata(metaPath),
		bootstrapNodes: node,
	}

	*s = *NewServer(opts)
	return s
}
