package blockchain

import (
	"blockchain/internal/algorythms"
	"blockchain/internal/block"
	httpmap "blockchain/internal/httpMap"
	"fmt"
)

const (
	BlocksFile = "blocks"

	TipFile = "tip"
	tip     = "tipKey"
)

type Blockchain struct {
	tip *block.Block
}

func initStorage() (err error) {
	genesis := block.NewGenesisBlock()

	value, err := genesis.StringSerialize()
	if err != nil {
		return
	}

	err = httpmap.Store(BlocksFile, genesis.StringHash(), value)
	if err != nil {
		return
	}

	err = httpmap.Store(TipFile, tip, value)
	if err != nil {
		return
	}

	return
}

func NewBlockchain() (b *Blockchain, err error) {
	var ok bool

	ok, err = httpmap.CheckFiles([]string{BlocksFile, TipFile})
	if err != nil && !ok {
		return
	}

	ok, err = httpmap.CheckKeys(TipFile, []string{tip})
	if err != nil {
		return
	}

	if !ok {
		err = initStorage()
		if err != nil {
			return
		}
	}

	data, err := httpmap.Load(TipFile, tip)
	if err != nil {
		return
	}

	tipBlock, err := block.StringDeserializeBlock(data)
	if err != nil {
		return
	}

	b = &Blockchain{tipBlock}
	return
}

func (bc *Blockchain) AddBlock(data block.BlockData) (err error) {
	prevBlock := bc.tip
	newBlock := block.NewBlock(data, prevBlock.Header.Hash)

	value, err := newBlock.StringSerialize()
	if err != nil {
		return
	}

	err = httpmap.Store(BlocksFile, newBlock.StringHash(), value)
	if err != nil {
		return
	}

	err = httpmap.Store(TipFile, tip, value)
	if err != nil {
		return
	}

	bc.tip = newBlock
	return
}

func (bc *Blockchain) ValidateBlocks() (result bool) {
	result = true

	keys, err := httpmap.Keys(BlocksFile)
	if err != nil {
		result = false
		return
	}

	var blocks []*block.Block
	/*! Тут собираем блоки в произвольном порядке !*/
	for _, key := range keys {
		value, err := httpmap.Load(BlocksFile, key)
		if err != nil {
			result = false
			return
		}

		b, err := block.StringDeserializeBlock(value)
		if err != nil {
			result = false
			return
		}

		blocks = append(blocks, b)
	}

	for _, b := range blocks {
		valid := algorythms.Validate(b.PrepareForValidate(), b.Header.TargetBits)
		if valid {
			fmt.Printf("Block '%s' is valid\n", b.Data.Name)
		} else {
			fmt.Printf("Block '%s' is invalid\n", b.Data.Name)
		}
		result = result && valid
	}

	return
}
