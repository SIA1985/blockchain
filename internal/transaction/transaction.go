package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
)

const (
	GenesisCoinbaseData = "Pust vse budet kak reshit monolog tvoei dushi!"
)

type Transaction struct {
	Hash []byte
	VOut []TXOutput
	VIn  []TXInput
}

func (t *Transaction) SetHash() (err error) {
	var buffer bytes.Buffer

	err = gob.NewEncoder(&buffer).Encode(t)
	if err != nil {
		return err
	}

	hash := sha256.Sum256(buffer.Bytes())
	t.Hash = hash[:]

	return
}

func (t *Transaction) IsCoinbase() bool {
	return len(t.VIn) == 1 && len(t.VIn[0].TxId) == 0 && t.VIn[0].VOut == -1
}

type TXOutput struct {
	Value           int64
	ScriptPublicKey string
}

func (out *TXOutput) CanBeUnlockedWith(data string) bool {
	return out.ScriptPublicKey == data
}

type TXInput struct {
	TxId            []byte
	VOut            int64
	ScriptSignature string
}

func (in TXInput) CanUnlockOutputWith(data string) bool {
	return in.ScriptSignature == data
}

/*todo: динамическая, от блоков*/
func getSubsidy() int64 {
	return 1
}

func NewCoinbaseTX(to, data string) *Transaction {
	txin := TXInput{[]byte{}, -1, data}
	txout := TXOutput{getSubsidy(), to}

	tx := Transaction{nil, []TXOutput{txout}, []TXInput{txin}}
	tx.SetHash()

	return &tx
}

func NewUTXOTransaction(from, to string, amount int64, unspentOutputs map[string][]int64, accumulated int64) (tx *Transaction) {
	var TXin []TXInput
	var TXout []TXOutput

	if accumulated < amount {
		panic("Not enough money!")
	}

	/*inputs*/
	for txId, outs := range unspentOutputs {
		txId, _ := hex.DecodeString(txId)

		for _, out := range outs {
			TXin = append(TXin, TXInput{txId, out, from})
		}
	}

	/*outputs*/
	TXout = append(TXout, TXOutput{amount, to})
	if accumulated > amount {
		TXout = append(TXout, TXOutput{accumulated - amount, from})
	}

	tx = &Transaction{nil, TXout, TXin}
	tx.SetHash()

	return
}
