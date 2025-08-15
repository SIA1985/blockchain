package blockchain

import (
	"blockchain/internal/algorythms"
	"blockchain/internal/block"
	httpmap "blockchain/internal/httpMap"
	"blockchain/internal/transaction"
	"fmt"
	"slices"
)

const (
	/*hash(block) -> block*/
	BlocksFile = "blocks"

	/*tipKey -> block*/
	TipFile = "tip"
	tip     = "tipKey"

	/*txId -> array[TXOutput]*/
	UTXOFile = "utxo"

	SubsidyBase = 100
)

type Blockchain struct {
	tip       *block.Block
	myAddress []byte
}

func initStorage(address []byte) (err error) {
	genesis := block.NewGenesisBlock(transaction.NewCoinbaseTX(address, SubsidyBase, 0))

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

	ok, err = httpmap.CheckFiles([]string{BlocksFile, TipFile, UTXOFile})
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

func (bc *Blockchain) AddBlock(data block.BlockData, comission int64) (err error) {
	prevBlock := bc.tip
	newBlock := block.NewBlock(data, prevBlock.Header.Hash, prevBlock.Header.Height)

	data.Transactions = append(data.Transactions, transaction.NewCoinbaseTX(bc.myAddress, bc.getSubsidy(), comission))

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

	for b := range ForEach(bc) {
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

func (bc *Blockchain) FindUTXO(address []byte) (UTXO map[string][]transaction.TXOutput, err error) {
	publicKeyHash := algorythms.PublicKeyHash(address)

	/*todo: оптимизировать запросом всего*/
	txIds, err := httpmap.Keys(UTXOFile)
	if err != nil {
		return
	}

	var outsDeserialized string
	for _, txId := range txIds {
		outsDeserialized, err = httpmap.Load(UTXOFile, txId)
		if err != nil {
			return

		}

		var outs []transaction.TXOutput
		outs, err = transaction.TXOutArrayDesiralizeFromString(outsDeserialized)
		if err != nil {
			return
		}

		for _, out := range outs {
			if out.IsLockedWithKey(publicKeyHash) {
				UTXO[txId] = append(UTXO[txId], out)
			}
		}

	}

	return
}

func (bc Blockchain) UpdateUTXO(txs []*transaction.Transaction) (err error) {
	for _, tx := range txs {
		/*outs*/
		var outsString string
		outsString, err = transaction.TXOutArraySerializeToString(tx.VOut)
		if err != nil {
			return
		}

		err = httpmap.Store(UTXOFile, tx.TxId(), outsString)
		if err != nil {
			return
		}

		/*ins*/
		var value string
		var outs []transaction.TXOutput
		for _, in := range tx.VIn {
			value, err = httpmap.Load(UTXOFile, in.RefTxId())
			if err != nil {
				return
			}

			outs, err = transaction.TXOutArrayDesiralizeFromString(value)
			if err != nil {
				return
			}

			outs = slices.Delete(outs, int(in.VOut), int(in.VOut))

			value, err = transaction.TXOutArraySerializeToString(outs)
			if err != nil {
				return
			}

			if len(outs) == 0 {
				err = httpmap.Delete(UTXOFile, in.RefTxId())
			} else {
				err = httpmap.Store(UTXOFile, in.RefTxId(), value)
			}
			if err != nil {
				return
			}
		}
	}

	return
}

func (bc *Blockchain) FindOutputsToSpend(address []byte, amount int64) (unspentOutputs map[string][]int64, accumulated int64, err error) {
	publicKeyHash := algorythms.PublicKeyHash(address)

	unspentOutputs = make(map[string][]int64)
	accumulated = 0

	UTXO, err := bc.FindUTXO(address)
	if err != nil {
		return
	}

Work:
	for txId, outs := range UTXO {
		for outId, out := range outs {
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
