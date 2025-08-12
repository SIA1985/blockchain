package wallet

import (
	"blockchain/internal/algorythms"
	httpmap "blockchain/internal/httpMap"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
)

const (
	walletFile = "wallets"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w Wallet) GetAddress() []byte {
	publicKeyHash := algorythms.HashPublicKey(w.PublicKey)

	versionedPayload := append([]byte("01"), publicKeyHash...)
	checksum := algorythms.Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := fullPayload //todo: Base58Encode

	return address
}

func (w Wallet) Save() error {
	if ok, err := httpmap.CheckFiles([]string{walletFile}); !ok && err != nil {
		return err
	}

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(w); err != nil {
		return err
	}

	if err := httpmap.Store(walletFile, hex.EncodeToString(w.PublicKey), hex.EncodeToString(result.Bytes())); err != nil {
		return err
	}

	return nil
}

func (w *Wallet) Load(address string) (err error) {
	if ok, err := httpmap.CheckFiles([]string{walletFile}); !ok && err != nil {
		return err
	}

	value, err := httpmap.Load(walletFile, address)
	if err != nil {
		return
	}

	var data []byte
	data, err = hex.DecodeString(value)
	if err != nil {
		return
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(w); err != nil {
		return
	}

	return nil
}

type Wallets map[string]*Wallet

func NewWallet() *Wallet {
	private, public := newKeyPair()

	return &Wallet{private, public}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	return *private, append(private.PublicKey.X.Bytes(), private.Y.Bytes()...)
}

func Addresses() (addresses []string, err error) {
	if ok, err := httpmap.CheckFiles([]string{walletFile}); !ok && err != nil {
		return []string{}, err
	}

	if addresses, err = httpmap.Keys(walletFile); err != nil {
		return
	}

	return
}
