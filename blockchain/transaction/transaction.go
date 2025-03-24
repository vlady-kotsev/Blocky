package transaction 

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

func CoinbaseTx(to, data string) (*Transaction, error) {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}
	txIn := TxInput{
		ID:  []byte{},
		Out: -1,
		Sig: data,
	}
	txOut := TxOutput{
		Value:  100,
		Pubkey: to,
	}
	tx := Transaction{
		ID:      nil,
		Inputs:  []TxInput{txIn},
		Outputs: []TxOutput{txOut},
	}
	err := tx.SetID()
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	encoder := gob.NewEncoder(&encoded)
	err := encoder.Encode(tx)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]

	return nil
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && tx.Inputs[0].Out == -1
}

