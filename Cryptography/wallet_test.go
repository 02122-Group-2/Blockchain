package cryptography

import (
	"bytes"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/keystore"
)

var testingPassword = "testingPassword123"
var testingHashedTransaction = [32]byte{'h', 'e', 'j', 's', 'a','b','m','m','h', 'e', 'j', 's', 'a','b','m','m','h', 'e', 'j', 's', 'a','b','m','m','h', 'e', 'j', 's', 'a','b','m','m'}

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


func TestSignTransaction(t *testing.T) {
	CreateNewWallet(testingPassword)
	signature1, err := SignTransaction(testingPassword, testingHashedTransaction)

	if err != nil {
		removeNewestWallet()
		t.Errorf("Failed to sign signature1\n")
	}

	signature2, err := SignTransaction(testingPassword, testingHashedTransaction)

	if err != nil {
		removeNewestWallet()
		t.Errorf("Failed to sign signature2\n")
	}

	if bytes.Compare(signature1, signature2) != 0  {
		removeNewestWallet()
		t.Errorf("Signature of same transaction is not identical\n")
	}


	removeNewestWallet()
}

func TestRetrieveAddressFromSignature(t *testing.T) {
	addr, _ := CreateNewWallet(testingPassword)
	signature, _ := SignTransaction(testingPassword, testingHashedTransaction)
	retrievedAddr, err := GetAddressFromSignedTransaction(signature, testingHashedTransaction)

	if err != nil {
		removeNewestWallet()
		t.Errorf("failed to get address from signed transaction\n")
	}

	if addr != retrievedAddr {
		removeNewestWallet()
		t.Errorf("wallet address and address retrieved from signed transaction doesn't match.\n")
	}


	removeNewestWallet()
}




