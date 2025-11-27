package core

import "time"

type BlockHeader struct {
	Version       uint32
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
			PrevBlockHash: prevHash,
			Timestamp:     time.Now().UnixNano(),
			Height:        height,
		},
		Transactions: txs,
	}
}
