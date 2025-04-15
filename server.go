package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
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

type Message struct {
	from    string
	payload any
}

type dataMessage struct {
	key  string
	data []byte
}

func (s *Server) StoreData(key string, w io.Reader) error {
	buff := new(bytes.Buffer)
	tee := io.TeeReader(w, buff)

	err := s.Store.WriteStream(key, tee)
	if err != nil {
		return err
	}

	p := &dataMessage{
		key:  key,
		data: buff.Bytes(),
	}

	err = s.Broadcast(*p)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Broadcast(d dataMessage) error {
	peerList := []io.Writer{}
	for _, peer := range s.peers {
		peerList = append(peerList, peer)
	}

	mw := io.MultiWriter(peerList...)

	msg := Message{
		from:    "self",
		payload: &d,
	}

	return gob.NewEncoder(mw).Encode(msg)
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
			var message Message
			err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&message)
			log.Print(err)
			fmt.Println(message)

			if err := s.handleMessage(&message); err != nil {
				return
			}

		case <-s.quitch:
			return
		}
	}
}

func (s *Server) handleMessage(msg *Message) error {
	switch m := msg.payload.(type) {
	case *dataMessage:
		fmt.Printf("recieved data: %+v\n", m)
	}

	return nil
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
