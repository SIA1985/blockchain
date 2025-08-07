package wallet

import (
	"blockchain/internal/algorythms"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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
