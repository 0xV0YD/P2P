package main

import (
	"fmt"
	"p2p/core"
	"p2p/crypto"
)

func main() {
	// 1. Create Identities
	// Alice is the legitimate user
	alicePrivKey, _ := crypto.GeneratePrivateKey()
	alicePubKey := alicePrivKey.PublicKey()

	// Bob is the receiver (just an address for now)
	// bobPrivKey, _ := crypto.GeneratePrivateKey()

	fmt.Printf("Alice Public Key: %s\n", alicePubKey.Address())

	// 2. Create a Transaction
	msg := []byte("Send 100 Tokens to Bob")
	tx := core.NewTransaction(msg, *alicePubKey, 0)

	// 3. Alice Signs the Transaction
	// Without this step, the network should reject it
	if err := tx.Sign(alicePrivKey); err != nil {
		panic(err)
	}
	fmt.Println("Transaction Signed by Alice.")

	// 4. Verify the Transaction (Network Side)
	if err := tx.Verify(); err != nil {
		fmt.Printf("ERROR: Verification Failed: %s\n", err)
	} else {
		fmt.Println("SUCCESS: Transaction Signature Valid!")
	}

	fmt.Println("--- ATTACK SCENARIO ---")

	// 5. Eve tries to tamper with the data
	// Eve intercepts the transaction and changes the amount
	tx.Data = []byte("Send 100000 Tokens to Bob")

	fmt.Println("Eve changed the data to: Send 100000 Tokens to Bob")

	// 6. Verify again
	// The signature was created for "Send 100...", but the data is now "Send 100000..."
	// The math should fail.
	if err := tx.Verify(); err != nil {
		fmt.Printf("SUCCESS: Tampered transaction rejected! Error: %s\n", err)
	} else {
		fmt.Println("FAIL: The network accepted a hacked transaction.")
	}
}
