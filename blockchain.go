package main

import "fmt"

type Blockchain struct {
	blocks []*Block
}

func (bc *Blockchain) AddBlock(data string) error {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	if len(bc.blocks) == 0 {
		return fmt.Errorf("Blockchain not instantiated")
	}
	newBlock := CreateBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
	return nil
}

func InitBlockchain() *Blockchain {
	chain := Blockchain{}
	zeroHash := [32]byte{}
	chain.blocks = append(chain.blocks, CreateBlock(GenesisData, zeroHash[:]))
	return &chain

}
