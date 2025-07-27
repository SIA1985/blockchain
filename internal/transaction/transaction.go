package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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
