package cryptography

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

var testingPassword = "testingPassword123"

func removeNewestWallet() {
	ks := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	allAccs := ks.Accounts()
	os.Remove(allAccs[len(allAccs)-1].URL.Path)
}

func TestCreateWallet(t *testing.T) {
	newAddr, _ := CreateNewWallet(testingPassword)
	
	if newAddr != GetAddress() {
		removeNewestWallet();
		t.Errorf("new account address doesn't match the account address returned by GetAddress\n")
	}

	removeNewestWallet();
}

func TestAccessWalletWithWrongPassword(t *testing.T) {
	CreateNewWallet(testingPassword)
	val, err := GetPrivateKey(testingPassword+"_wrong")
	removeNewestWallet()
	if val != nil || err == nil {
		t.Errorf("expected to fail when given a wrong password, but didn't\n")
	}
}

func TestAccessWalletWithCorrectPassword(t *testing.T) {
	CreateNewWallet(testingPassword)
	val, err := GetPrivateKey(testingPassword)
	if val == nil || err != nil {
		removeNewestWallet()
		t.Errorf("expected to be able to get private key with correct password, but couldn't\n")
	}
	removeNewestWallet()
}




