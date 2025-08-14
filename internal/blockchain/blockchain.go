package blockchain

import (
	"blockchain/internal/algorythms"
	"blockchain/internal/block"
	httpmap "blockchain/internal/httpMap"
	"blockchain/internal/transaction"
	"encoding/hex"
	"fmt"
)

const (
	BlocksFile = "blocks"

	TipFile = "tip"
	tip     = "tipKey"

	SubsidyBase = 100
)

type Blockchain struct {
	tip       *block.Block
	myAddress []byte
}

func initStorage(address []byte) (err error) {
	genesis := block.NewGenesisBlock(transaction.NewCoinbaseTX(address, SubsidyBase))

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

func NewBlockchain(address []byte) (b *Blockchain, err error) {
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
		err = initStorage(address)
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

	b = &Blockchain{tipBlock, address}
	return
}

func (bc Blockchain) getSubsidy() int64 {

	return SubsidyBase >> (2 * int64((bc.tip.Header.Height / 3)))
}

func (bc *Blockchain) AddBlock(data block.BlockData) (err error) {
	prevBlock := bc.tip
	newBlock := block.NewBlock(data, prevBlock.Header.Hash, prevBlock.Header.Height)
	data.Transactions = append(data.Transactions, transaction.NewCoinbaseTX(bc.myAddress, bc.getSubsidy()))

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

	var blocks []*block.Block
	for b := range ForEach(bc) {
		blocks = append(blocks, b)
	}

	for _, b := range blocks {
		valid := algorythms.Validate(b.PrepareForValidate(), b.Header.TargetBits)
		if valid {
			fmt.Printf("%s - valid\n", b.StringHash())
		} else {
			fmt.Printf("%s - invalid\n", b.StringHash())
		}
		result = result && valid
	}

	return
}

func (bc *Blockchain) Iterator() BlockchainIterator {
	return BlockchainIterator{bc.tip.StringHash()}
}

func (bc *Blockchain) FindUnspentTx(address []byte) (unspent []transaction.Transaction) {
	publicKeyHash := algorythms.PublicKeyHash(address)
	spentTx := make(map[string][]int)

	for b := range ForEach(bc) {
		for _, tx := range b.Data.Transactions {
			txId := hex.EncodeToString(tx.Hash)

		outputs:
			for outIdx, out := range tx.VOut {
				if _, ok := spentTx[txId]; ok {
					for spentOut := range spentTx[txId] {
						if spentOut == outIdx {
							continue outputs
						}
					}
				}

				if out.IsLockedWithKey(publicKeyHash) {
					unspent = append(unspent, *tx)
				}
			}

			if tx.IsCoinbase() {
				continue
			}

			for _, in := range tx.VIn {
				if in.UsesKey(publicKeyHash) {
					inTxId := hex.EncodeToString(in.TxId)
					spentTx[inTxId] = append(spentTx[inTxId], int(in.VOut))
				}
			}
		}

		if len(b.Header.PrevBlockHash) == 0 {
			break
		}
	}

	return
}

func (bc *Blockchain) FindUTXO(address []byte) (UTXO []transaction.TXOutput) {
	publicKeyHash := algorythms.PublicKeyHash(address)

	for _, tx := range bc.FindUnspentTx(address) {
		for _, out := range tx.VOut {
			if out.IsLockedWithKey(publicKeyHash) {
				UTXO = append(UTXO, out)
			}
		}
	}

	return
}

func (bc *Blockchain) FindOutputsToSpend(address []byte, amount int64) (unspentOutputs map[string][]int64, accumulated int64) {
	publicKeyHash := algorythms.PublicKeyHash(address)

	unspentOutputs = make(map[string][]int64)
	accumulated = 0

Work:
	for _, tx := range bc.FindUnspentTx(address) {
		txId := hex.EncodeToString(tx.Hash)

		for outId, out := range tx.VOut {
			if out.IsLockedWithKey(publicKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txId] = append(unspentOutputs[txId], int64(outId))

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return
}
