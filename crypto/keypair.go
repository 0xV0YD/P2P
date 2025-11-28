package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

type PublicKey struct {
	key ed25519.PublicKey
}

type Signature struct {
	value []byte
}

func GeneratePrivateKey() (*PrivateKey, error) {
	seed := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		return nil, err
	}

	// ed25519 keys are generated from a seed
	key := ed25519.NewKeyFromSeed(seed)
	return &PrivateKey{key: key}, nil
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(msg []byte) *Signature {
	sig := ed25519.Sign(p.key, msg)
	return &Signature{value: sig}
}

func (p *PrivateKey) PublicKey() *PublicKey {
	// The last 32 bytes of the Ed25519 private key is the public key
	pub := make([]byte, 32)
	copy(pub, p.key[32:])

	return &PublicKey{key: pub}
}

// Bytes returns the raw bytes of the public key
func (p *PublicKey) Bytes() []byte {
	return p.key
}

// Address returns a string representation (hex) of the public key
// In Ethereum this would be 0x..., here we just hex the key
func (p *PublicKey) Address() string {
	return hex.EncodeToString(p.key)
}

// Verify checks if the signature is valid for the given message
func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.key, msg, s.value)
}

// Bytes returns the raw bytes of the signature
func (s *Signature) Bytes() []byte {
	return s.value
}
