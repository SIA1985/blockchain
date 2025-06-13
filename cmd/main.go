package main

import (
	"blockchain/internal/block"
	"blockchain/internal/blockchain"
)

func main() {
	var err error = nil

	bc, err := blockchain.NewBlockchain()
	if err != nil {
		panic(err)
	}

	bc.AddBlock(block.BlockData{Name: "Block 1"})
	bc.AddBlock(block.BlockData{Name: "Block 2"})

	bc.ValidateBlocks()
}
