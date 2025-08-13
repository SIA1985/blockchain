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
	"math/big"
)

const (
	walletFile = "wallets"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
}

func (w Wallet) PublicKey() []byte {
	return append(w.PrivateKey.PublicKey.X.Bytes(), w.PrivateKey.PublicKey.Y.Bytes()...)
}

func (w Wallet) Address() []byte {
	publicKeyHash := algorythms.HashPublicKey(w.PublicKey())

	versionedPayload := append([]byte("01"), publicKeyHash...)
	checksum := algorythms.Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := fullPayload //todo: Base58Encode

	return address
}

type saveWallet struct {
	X *big.Int
	Y *big.Int
	D *big.Int
}

func (w Wallet) Save() error {
	if ok, err := httpmap.CheckFiles([]string{walletFile}); !ok && err != nil {
		return err
	}

	sw := saveWallet{w.PrivateKey.X, w.PrivateKey.Y, w.PrivateKey.D}

	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(&sw); err != nil {
		return err
	}

	if err := httpmap.Store(walletFile, hex.EncodeToString(w.Address()), hex.EncodeToString(result.Bytes())); err != nil {
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

	sw := saveWallet{}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err = decoder.Decode(&sw); err != nil {
		return
	}

	w.PrivateKey.D = sw.D
	w.PrivateKey.X = sw.X
	w.PrivateKey.Y = sw.Y
	w.PrivateKey.Curve = elliptic.P256()

	return nil
}

type Wallets map[string]*Wallet

func NewWallet() *Wallet {
	return &Wallet{newKeyPair()}
}

func newKeyPair() ecdsa.PrivateKey {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}

	return *private
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
