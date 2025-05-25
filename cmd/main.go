package main

import (
	"blockchain/internal/block"
	"blockchain/internal/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()

	bc.AddBlock(block.BlockData{Name: "Block 1"})
	bc.AddBlock(block.BlockData{Name: "Block 2"})
}
