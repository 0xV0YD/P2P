package main

import (
	"encoding/hex"
	"fmt"
	"p2p/core" // Replace with your actual module name
)

func main() {
	// 1. Create the Genesis Block (Block 0)
	genesisTx := core.NewTransaction([]byte("Genesis Tx"), nil, 0)
	genesisBlock := core.NewBlock(nil, 0, []*core.Transaction{genesisTx})

	hasher := core.BlockHasher{}
	hash1, _ := genesisBlock.Hash(hasher)

	fmt.Printf("Block 0 Hash: %s\n", hex.EncodeToString(hash1))

	// 2. Create Block 1 (Linked to Block 0)
	tx2 := core.NewTransaction([]byte("Bob sends 5 coins to Alice"), nil, 1)
	block1 := core.NewBlock(hash1, 1, []*core.Transaction{tx2})

	hash2, _ := block1.Hash(hasher)
	fmt.Printf("Block 1 Hash: %s\n", hex.EncodeToString(hash2))

	// 3. Verify the Link
	fmt.Printf("Block 1 PrevHash: %s\n", hex.EncodeToString(block1.Header.PrevBlockHash))

	// Quick integrity check
	if hex.EncodeToString(block1.Header.PrevBlockHash) == hex.EncodeToString(hash1) {
		fmt.Println("SUCCESS: The chain is linked correctly!")
	} else {
		fmt.Println("FAIL: Broken chain.")
	}
}
