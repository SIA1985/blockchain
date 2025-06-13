package block

import (
	"blockchain/internal/algorythms"
	"bytes"
	"encoding/gob"
	"encoding/hex"
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

func (h *BlockHeader) prepareForPOW() []byte {
	timestamp := algorythms.Int64ToByteArr(h.Timestamp)

	return bytes.Join([][]byte{timestamp, h.PrevBlockHash}, []byte{})
}

type BlockData struct {
	Name string
}

func (d *BlockData) prepareForPOW() []byte {
	return []byte(d.Name)
}

type Block struct {
	Header BlockHeader
	Data   BlockData
}

func (b *Block) PrepareForPOW() []byte {
	return bytes.Join([][]byte{b.Header.prepareForPOW(), b.Data.prepareForPOW()}, []byte{})
}

func (b *Block) PrepareForValidate() []byte {
	return bytes.Join([][]byte{b.PrepareForPOW(), algorythms.Int64ToByteArr(b.Header.Nonce)}, []byte{})
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)

	return result.Bytes(), err
}

func (b *Block) StringSerialize() (value string, err error) {
	data, err := b.Serialize()
	if err != nil {
		return
	}

	value = hex.EncodeToString(data)
	return
}

func DeserializeBlock(data []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)

	return &block, err
}

func StringDeserializeBlock(value string) (*Block, error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return nil, err
	}

	return DeserializeBlock(data)
}

func (b *Block) StringHash() string {
	return hex.EncodeToString(b.Header.Hash)
}

func NewBlock(data BlockData, prevBlockHash []byte) *Block {
	header := BlockHeader{
		time.Now().Unix(),
		prevBlockHash,
		[]byte{},
		0,
		0}

	block := &Block{header, data}

	block.Header.Hash, block.Header.Nonce, block.Header.TargetBits = algorythms.ProofOfWork(block.PrepareForPOW())

	return block
}

func NewGenesisBlock() *Block {
	return NewBlock(BlockData{"Genesis"}, []byte{})
}
