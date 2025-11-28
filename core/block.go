package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"time"
)

type BlockHeader struct {
	Version       uint32
	DataHash      []byte
	PrevBlockHash []byte // The link to the previous block
	Timestamp     int64  // When the block was created
	Height        uint32 // Block number (0, 1, 2...)
	Nonce         uint64 // Used later for consensus voting logic
}
type Block struct {
	Header       *BlockHeader
	Transactions []*Transaction
	hash         []byte
}

func NewBlock(prevHash []byte, height uint32, txs []*Transaction) *Block {
	return &Block{
		Header: &BlockHeader{
			Version:       1,
			DataHash:      calculateDataHash(txs),
			PrevBlockHash: prevHash,
			Timestamp:     time.Now().UnixNano(),
			Height:        height,
		},
		Transactions: txs,
	}
}

// calculateDataHash creates a SHA256 hash of all transaction data combined
func calculateDataHash(txs []*Transaction) []byte {
	// In a real blockchain, this would be a Merkle Tree.
	// For now, we simply concatenate all tx data and hash it.

	buf := &bytes.Buffer{}

	for _, tx := range txs {
		// We need to encode the tx to bytes to hash it
		if err := gob.NewEncoder(buf).Encode(tx); err != nil {
			// In production handle error, here we panic for simplicity
			panic(err)
		}
	}

	hash := sha256.Sum256(buf.Bytes())
	return hash[:]
}
