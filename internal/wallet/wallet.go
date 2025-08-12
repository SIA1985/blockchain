package wallet

import (
	"blockchain/internal/algorythms"
	httpmap "blockchain/internal/httpMap"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

	//todo:
	if err := httpmap.Store(walletFile, hex.EncodeToString(w.PublicKey), fmt.Sprintf("%x", w.PrivateKey)); err != nil {
		return err
	}

	return nil
}

func (w *Wallet) Load(address string) (err error) {
	if ok, err := httpmap.CheckFiles([]string{walletFile}); !ok && err != nil {
		return err
	}

	//todo:
	value, err := httpmap.Load(walletFile, address)
	if err != nil {
		return
	}

	w.PublicKey, err = hex.DecodeString(value)
	if err != nil {
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
