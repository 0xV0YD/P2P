package network

import (
	"encoding/gob"
	"io"
)

// Decoder is an interface that reads from a generic reader (like a TCP connection)
// and fills a generic struct (RPC) with data.
type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GOBDecoder struct{}

func (g GOBDecoder) Decode(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}

// DefaultGobDecoder is a helper function we might use later
func DefaultGobDecoder(r io.Reader, msg *RPC) error {
	return gob.NewDecoder(r).Decode(msg)
}
