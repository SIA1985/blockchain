package block

import (
	"blockchain/internal/algorythms"
	"bytes"
	"time"
)

type BlockHeader struct {
	Timestamp     int64
	PrevBlockHash []byte
	Hash          []byte

	/*nohasable*/
	Nonce      int64
	TargetBits uint64 /*todo: Присваивать после Proof-of-Work*/
}

func (h *BlockHeader) toByteArr() []byte {
	timestamp := algorythms.Int64ToByteArr(h.Timestamp)

	return bytes.Join([][]byte{timestamp, h.PrevBlockHash}, []byte{})
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

func (b *Block) ToByteArr() []byte {
	return bytes.Join([][]byte{b.Header.toByteArr(), b.Data.toByteArr()}, []byte{})
}

func (b *Block) ToByteArrWithNonce() []byte {
	return bytes.Join([][]byte{b.ToByteArr(), algorythms.Int64ToByteArr(b.Header.Nonce)}, []byte{})
}

func NewBlock(data BlockData, prevBlockHash []byte) *Block {
	header := BlockHeader{
		time.Now().Unix(),
		prevBlockHash,
		[]byte{},
		0,
		0}

	block := &Block{header, data}

	block.Header.Hash, block.Header.Nonce, block.Header.TargetBits = algorythms.ProofOfWork(block.ToByteArr())

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock(BlockData{"Genesis"}, []byte{})
}
