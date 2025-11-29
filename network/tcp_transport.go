package network

import (
	"fmt"
	"io"
	"net"
	"sync"
)

// TCPPeer represents a connection to another node.
// Currently, it just wraps the raw TCP connection.
type TCPPeer struct {
	conn net.Conn // The underlying TCP connection

	// If true, WE called THEM. If false, THEY called US.
	outbound bool
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener

	// rpcCh is the channel where we will push incoming messages.
	// The "Consume" method will return this channel to the caller.
	rpcCh chan RPC

	// peers is a map of all active connections.
	// Key: The network address of the peer. Value: The peer struct.
	peers map[net.Addr]*TCPPeer

	// mu is a Mutex (Lock).
	// Because multiple peers might connect at the exact same time,
	// we need to lock the map before adding/removing to prevent crashes.
	mu      sync.RWMutex
	decoder Decoder
}

// NewTCPTransport is the constructor.
func NewTCPTransport(addr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: addr,
		rpcCh:         make(chan RPC),
		peers:         make(map[net.Addr]*TCPPeer),
		decoder:       GOBDecoder{},
	}
}

// Consume implements the Transport interface.
// It just gives the channel to the caller (read-only access).
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) listenAndAccept() error {
	var err error
	// open the port (Start listening)
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	// Start a background thread (Goroutine) to accept incoming connections.
	// We use `go` because we don't want to block the main program here.
	go t.startAcceptLoop()

	return nil
}

// startAcceptLoop runs forever, waiting for new connections.
func (t *TCPTransport) startAcceptLoop() {
	for {
		// 1. Wait here until someone connects
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept Error: %s\n", err)
			continue
		}

		// 2. Someone connected! Let's handle them in a NEW Goroutine.
		// Why new Goroutine? So we can go back to waiting for the next person immediately.
		go t.handleConn(conn, false)
	}
}

// handleConn is the logic for a single connection.
func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	fmt.Printf("New connection from: %+v\n", conn.RemoteAddr())

	// Create the Peer struct
	peer := &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}

	// Add to our map (Thread-Safe)
	t.mu.Lock()
	t.peers[conn.RemoteAddr()] = peer
	t.mu.Unlock()

	// --- NEW: The Read Loop ---
	// We defer closing the connection. If the loop breaks (error), we hang up.
	defer func() {
		fmt.Printf("Dropping peer: %s\n", conn.RemoteAddr())
		conn.Close()
	}()

	for {
		rpc := RPC{} // Create an empty envelope

		// 1. Read from the wire into the envelope
		err := t.decoder.Decode(conn, &rpc)

		if err == io.EOF {
			// This means the other side closed the connection
			return
		}
		if err != nil {
			fmt.Printf("TCP Read Error: %s\n", err)
			continue
		}

		// 2. Tag the message so we know who sent it
		rpc.From = NetAddr(conn.RemoteAddr().String())

		// 3. Send it to the main channel for the Server to process
		// This will block here until the Server reads it!
		t.rpcCh <- rpc
	}
}

// The Logic Flow
// Main program calls ListenAndAccept().
// Socket opens.
// startAcceptLoop begins running in the background.
// Wait... (It pauses at t.listener.Accept() until a user connects).
// User connects! Accept returns a conn object.
// handleConn starts in another background thread to manage that specific user.
// startAcceptLoop immediately loops back to Accept to wait for the next user.
