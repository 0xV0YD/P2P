package core

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"p2p/crypto"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	Signature *crypto.Signature
	Nonce     uint64
	hash      []byte
}

func NewTransaction(data []byte, from crypto.PublicKey, nonce uint64) *Transaction {
	return &Transaction{
		Data:  data,
		From:  from,
		Nonce: nonce,
	}
}

func (tx *Transaction) ID() []byte {
	return tx.hash
}

// Sign uses the private key to sign the transaction payload
func (tx *Transaction) Sign(privKey *crypto.PrivateKey) error {
	// 1. Calculate the hash of the transaction data (Data + From + Nonce)
	// We cannot use tx.Hash() yet because that might include the signature field depending on implementation.
	// Let's create a specific "Signable Data" payload.

	dataToSign, err := tx.EncodeBinary() // We will write this helper below
	if err != nil {
		return err
	}

	// 2. Sign the data
	sig := privKey.Sign(dataToSign)

	// 3. Attach signature to tx
	tx.Signature = sig
	return nil
}

// Verify checks if the signature matches the public key and data
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	// 1. Get the data that was supposed to be signed
	dataToVerify, err := tx.EncodeBinary()
	if err != nil {
		return err
	}

	// 2. Verify against the From field (which is the Public Key)
	if !tx.Signature.Verify(&tx.From, dataToVerify) {
		return fmt.Errorf("invalid transaction signature")
	}

	return nil
}

func (tx *Transaction) EncodeBinary() ([]byte, error) {
	var buf bytes.Buffer
	// We strictly encode the fields that define the transaction logic
	err := gob.NewEncoder(&buf).Encode(struct {
		Data  []byte
		From  []byte
		Nonce uint64
	}{
		Data:  tx.Data,
		From:  tx.From.Bytes(),
		Nonce: tx.Nonce,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
