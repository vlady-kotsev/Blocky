package main

import (
	"fmt"
)

func main() {
	chain := InitBlockchain()

	chain.AddBlock("1")
	chain.AddBlock("2")
	chain.AddBlock("3")

	for _, block := range chain.blocks {
		fmt.Printf("Hash: %x, Data: %s, PrevHash: %x\n", block.Hash, block.Data, block.PrevHash)
	}
}
