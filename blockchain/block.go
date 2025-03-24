package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"vlady-kotsev/blocky/blockchain/transaction"
)

type Block struct {
	Hash         []byte
	Transactions []*transaction.Transaction
	PrevHash     []byte
	Nonce        uint64
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.ID)
	}

	txBytes := bytes.Join(txHashes, []byte{})
	txHash := sha256.Sum256(txBytes)

	return txHash[:]
}

func (b *Block) SerializeBlock() ([]byte, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func DeserializeBlock(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}
