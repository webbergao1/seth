package types

import (
	"fmt"
	"math/big"
	"seth/accounts"
	"testing"
)

func Test_Transaction_SignVerify(t *testing.T) {
	fromaddress, fromprivatekey := accounts.NewRandomAccount()
	toaddress, _ := accounts.NewRandomAccount()

	signer := NewSethSigner(big.NewInt(1))
	tx := NewTransaction(toaddress, big.NewInt(10), 0)
	err := tx.Sign(signer, fromprivatekey)
	if err != nil {
		t.Fatalf("Failed to sign a tx!")
	}
	if tx.ChainID().Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("Error: the chain id of tx is not equal to signer!")
	}
	if tx.Verify(signer) != true {
		t.Fatalf("Failed: verify the transaction content failed!")
	}

	address, err := tx.Sender(signer)
	address, err = tx.Sender(signer)
	if err != nil {
		t.Fatalf("Failed: get the transaction sender failed!")
	}
	fmt.Println(address.Hex())
	fmt.Println(fromaddress.Hex())
	if address.Hex() != fromaddress.Hex() {
		t.Fatalf("Failed: the transaction sender verify failed!")
	}
}

func Test_Transaction_Hash(t *testing.T) {
	_, fromprivatekey := accounts.NewRandomAccount()
	toaddress, _ := accounts.NewRandomAccount()
	signer := NewSethSigner(big.NewInt(1))
	tx := NewTransaction(toaddress, big.NewInt(10), 0)
	hashBeforeSign := tx.Hash()
	err := tx.Sign(signer, fromprivatekey)
	if err != nil {
		t.Fatalf("Failed to sign a tx!")
	}
	hashAfterSign := tx.Hash()
	if hashBeforeSign.Equal(hashAfterSign) != true {
		t.Fatalf("Failed two hash of the transaction not equal!")
	}

}
