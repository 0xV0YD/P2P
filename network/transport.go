package network

// NetAddr represents a string address like "127.0.0.1:3000"
type NetAddr string

// RPC stands for "Remote Procedure Call".
// It holds the data coming from another node.
// From: Who sent this?
// Payload: The actual bytes (Block, Transaction, etc).

type RPC struct {
	From    NetAddr
	Payload []byte
}

// Transport is an interface that defines how we talk to the world.
// If we want to switch from TCP to UDP or Websockets later, we just implement this interface.
type Transport interface {
	// Consume returns a channel that provides incoming messages (RPCs).
	// We will "consume" messages from this channel to process them.
	// "<-chan" means this is a Read-Only channel.
	Consume() <-chan RPC

	// Connect allows us to dial another node.
	// We will implement this later.
	Connect(Transport) error

	// Addr returns the address this transport is listening on.
	Addr() NetAddr
}
