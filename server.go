package main

import (
	"fmt"
	"log"
	"sync"

	peer2peer "github.com/Faizan2005/DFS-Go/Peer2Peer"
)

type ServerOpts struct {
	storageRoot    string
	pathTransform  pathTransform
	tcpTransport   peer2peer.TCPTransport
	metaData       Metadata
	bootstrapNodes []string
}

type Server struct {
	peerLock sync.Mutex
	peers    map[string]peer2peer.Peer

	serverOpts ServerOpts
	Store      *Store
	quitch     chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	StoreOpts := StructOpts{
		PathTransformFunc: opts.pathTransform,
		Metadata:          &opts.metaData,
		Root:              opts.storageRoot,
	}

	return &Server{
		peers:      map[string]peer2peer.Peer{},
		serverOpts: opts,
		Store:      NewStore(StoreOpts),
		quitch:     make(chan struct{}),
	}
}

func (s *Server) Run() error {
	err := s.serverOpts.tcpTransport.ListenAndAccept()
	if err != nil {
		return err
	}

	if len(s.serverOpts.bootstrapNodes) != 0 {
		err := s.BootstrapNetwork()
		if err != nil {
			return err
		}
	}

	s.loop()
	return nil
}

func (s *Server) loop() {
	defer func() {
		s.serverOpts.tcpTransport.Close()
		log.Println("file server closed due to user quit action")
	}()

	for {
		select {
		case msg := <-s.serverOpts.tcpTransport.Consume():
			fmt.Println(msg)

		case <-s.quitch:
			return
		}
	}
}

func (s *Server) Stop() {
	close(s.quitch)
}

func (s *Server) BootstrapNetwork() error {
	for _, addr := range s.serverOpts.bootstrapNodes {

		if addr == "" {
			continue
		}
		go func() {
			fmt.Println("attempting to connect with remote: ", addr)
			if err := s.serverOpts.tcpTransport.Dial(addr); err != nil {
				log.Println("Dial error:", err)
			}
		}()
	}

	return nil
}

func (s *Server) OnPeer(p peer2peer.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	addr := p.RemoteAddr()

	s.peers[addr.String()] = p
	log.Printf("[OnPeer] Connected with remote peer: %s\n", addr.String())
	return nil
}
