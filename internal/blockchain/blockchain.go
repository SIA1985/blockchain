package blockchain

import (
	"blockchain/internal/algorythms"
	"blockchain/internal/block"
	"fmt"
)

type Blockchain struct {
	blocks []*block.Block
}

func NewBlockchain() *Blockchain {
	return &Blockchain{[]*block.Block{block.NewGenesisBlock()}}
}

func (bc *Blockchain) AddBlock(data block.BlockData) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := block.NewBlock(data, prevBlock.Header.Hash)

	bc.blocks = append(bc.blocks, newBlock)
}

func (bc *Blockchain) ValidateBlocks() (result bool) {
	result = true

	for _, b := range bc.blocks {
		valid := algorythms.Validate(b.ToByteArrWithNonce(), b.Header.TargetBits)
		if valid {
			fmt.Printf("Block '%s' is valid\n", b.Data.Name)
		} else {
			fmt.Printf("Block '%s' is invalid\n", b.Data.Name)
		}
		result = result && valid
	}

	return
}
