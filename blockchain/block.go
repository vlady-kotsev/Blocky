package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
	Nonce    uint64
}

func CreateGenesis() *Block {
	zeroHash := [32]byte{}
	data := []byte(GenesisData)
	bytes := bytes.Join([][]byte{
		data,
		zeroHash[:],
	}, []byte{})

	hash := sha256.Sum256(bytes)
	return &Block{
		Data:     data,
		PrevHash: zeroHash[:],
		Hash:     hash[:],
		Nonce:    0,
	}
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
