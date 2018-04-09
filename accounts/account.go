package accounts

import (
	"crypto/ecdsa"
	crand "crypto/rand"
	"encoding/hex"
	"io"
	"seth/common"
	"seth/crypto"
)

func newKey(rand io.Reader) (*ecdsa.PrivateKey, error) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
	if err != nil {
		return nil, err
	}
	return privateKeyECDSA, nil
}

//NewAccount new account return account address,publickey,privatekey
func NewAccount() (string, string, string) {
	key, err := newKey(crand.Reader)
	if err != nil {
		return "", "", ""
	}
	address := crypto.PubkeyToAddress(&key.PublicKey)
	publickey := hex.EncodeToString(crypto.FromECDSAPub(&key.PublicKey))
	Privatekey := hex.EncodeToString(crypto.FromECDSA(key))
	return address.Hex(), publickey, Privatekey
}

// NewRandomAccount new random account return address&privatekey
func NewRandomAccount() (common.Address, *crypto.PrivateKey) {
	key, err := newKey(crand.Reader)
	if err != nil {
		return common.Address{}, nil
	}
	return crypto.PubkeyToAddress(&key.PublicKey), (*crypto.PrivateKey)(key)

}
