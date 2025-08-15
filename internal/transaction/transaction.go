package transaction

import (
	"blockchain/internal/algorythms"
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

func (t Transaction) TxId() string {
	return hex.EncodeToString(t.Hash)
}

type TXOutput struct {
	Value int64

	//todo: свой скриптовый язык
	// ScriptPublicKey string

	PublicKeyHash []byte
}

func NewTXOutput(value int64, address []byte) (out *TXOutput) {
	out = &TXOutput{value, []byte{}}
	out.Lock(address)

	return
}

func (out *TXOutput) Lock(address []byte) {
	publicKeyHash := address //todo: Base58Decode
	publicKeyHash = algorythms.PublicKeyHash(address)
	out.PublicKeyHash = publicKeyHash
}

func (out TXOutput) IsLockedWithKey(publicKeyHash []byte) bool {
	return bytes.Equal(out.PublicKeyHash, publicKeyHash)
}

type TXInput struct {
	TxId []byte /*Id транзакции, на которую ссылаемся*/
	VOut int64  /*Индекс выхода транзакции с Id=TxId*/

	//todo: свой скриптовый язык
	// ScriptSignature string

	Signature []byte
	PublicKey []byte
}

func (in TXInput) RefTxId() string {
	return hex.EncodeToString(in.TxId)
}
func NewTXInput(txId []byte, vOut int64) *TXInput {
	return &TXInput{txId, vOut, []byte{}, []byte{}}
}

func (in TXInput) UsesKey(publicKeyHash []byte) bool {
	lockingHash := algorythms.HashPublicKey(in.PublicKey)

	return bytes.Equal(lockingHash, publicKeyHash)
}

func NewCoinbaseTX(to []byte, subsidy int64, commision int64) *Transaction {
	txin := *NewTXInput([]byte{}, -1)
	txout := *NewTXOutput(subsidy, to)
	txCommision := *NewTXOutput(commision, to)

	tx := Transaction{nil, []TXOutput{txout, txCommision}, []TXInput{txin}}
	tx.SetHash()

	return &tx
}

func NewUTXOTransaction(from, to []byte, amount int64, unspentOutputs map[string][]int64, accumulated int64) (tx *Transaction) {
	var TXin []TXInput
	var TXout []TXOutput

	if accumulated < amount {
		panic("Not enough money!")
	}

	/*inputs*/
	for txId, outs := range unspentOutputs {
		txId, _ := hex.DecodeString(txId)

		for _, out := range outs {
			TXin = append(TXin, *NewTXInput(txId, out))
		}
	}

	/*outputs*/
	TXout = append(TXout, *NewTXOutput(amount, to))
	if accumulated > amount {
		TXout = append(TXout, *NewTXOutput(accumulated-amount, from))
	}

	tx = &Transaction{nil, TXout, TXin}
	tx.SetHash()

	return
}

func TXOutArraySerializeToString(outs []TXOutput) (value string, err error) {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err = encoder.Encode(outs)

	value = hex.EncodeToString(result.Bytes())

	return
}

func TXOutArrayDesiralizeFromString(value string) (outs []TXOutput, err error) {
	data, err := hex.DecodeString(value)
	if err != nil {
		return
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err = decoder.Decode(&outs)
	if err != nil {
		return
	}

	return
}
