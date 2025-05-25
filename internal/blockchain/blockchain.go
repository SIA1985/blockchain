package blockchain

import "blockchain/internal/block"

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
