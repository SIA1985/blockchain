package block

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

type BlockHeader struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte
}

func (h *BlockHeader) toByteArr() []byte {
	t := []byte(strconv.FormatInt(h.Timestamp, 10))

	return bytes.Join([][]byte{t, h.PrevBlockHash}, []byte{})
}

type BlockData struct {
	Name string
}

func (d *BlockData) toByteArr() []byte {
	return []byte(d.Name)
}

type Block struct {
	Header BlockHeader
	Data   BlockData
}

func (b *Block) setHash() {
	byteBlock := append(b.Header.toByteArr(), b.Data.toByteArr()...)

	shaHash := sha256.Sum256(byteBlock)

	b.Header.Hash = shaHash[:]
}

func NewBlock(data BlockData, prevBlockHash []byte) *Block {
	header := BlockHeader{time.Now().Unix(), prevBlockHash, []byte{}}

	block := &Block{header, data}
	block.setHash()

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock(BlockData{"Genesis"}, []byte{})
}
