package core

type Transaction struct {
	Data      []byte
	From      []byte
	Signature []byte
	Nonce     uint64
	hash      []byte
}

func NewTransaction(data []byte, from []byte, nonce uint64) *Transaction {
	return &Transaction{
		Data:  data,
		From:  from,
		Nonce: nonce,
	}
}

func (tx *Transaction) ID() []byte {
	return tx.hash
}
