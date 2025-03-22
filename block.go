package main

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Hash     []byte
	Data     []byte
	PrevHash []byte
}

func (b *Block) DeriveHash() {
	bytes := bytes.Join([][]byte{b.Data, b.PrevHash}, []byte{})
	hash := sha256.Sum256(bytes)
	b.Hash = hash[:]
}

func CreateBlock(data string, prevHash []byte) *Block {
	block := Block{
		Data:     []byte(data),
		PrevHash: prevHash,
	}
	block.DeriveHash()
	return &block
}
