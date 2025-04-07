package peer2peer

import (
	"encoding/gob"
	"io"
	"log"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(w io.Reader, msg *Message) error {
	return gob.NewDecoder(w).Decode(msg)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(w io.Reader, msg *Message) error {
	buff := make([]byte, 1024)

	n, err := w.Read(buff)
	if err != nil {
		log.Printf("Error: %+v\n", err)
	}

	msg.Payload = buff[:n]

	return nil
}
