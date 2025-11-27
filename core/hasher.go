package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

// Hasher is an interface that defines how to hash a struct
type Hasher interface {
	Hash(input any) ([]byte, error)
}

// BlockHasher acts as a simplified hashing logic
type BlockHasher struct{}

// Implementation of the interface
func (bh BlockHasher) Hash(v any) ([]byte, error) {
	var buf bytes.Buffer
	// gob encodes the struct fields into bytes
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}

	hash := sha256.Sum256(buf.Bytes())
	return hash[:], nil
}

// --- Logic to Hash Specific Structs ---

// HashBlock calculates the hash of the Block Header
func (b *Block) Hash(hasher Hasher) ([]byte, error) {
	if b.hash != nil {
		return b.hash, nil
	}

	// Now this matches the interface signature Hash(input any)
	val, err := hasher.Hash(b.Header)
	if err != nil {
		return nil, err
	}

	b.hash = val
	return val, nil
}

// HashTransaction calculates the hash of the Transaction
func (tx *Transaction) Hash(hasher Hasher) ([]byte, error) {
	if tx.hash != nil {
		return tx.hash, nil
	}

	val, err := hasher.Hash(tx)
	if err != nil {
		return nil, err
	}

	tx.hash = val
	return val, nil
}
