package cryptography

import (
	"bytes"
	"testing"
)

var testingPassword = "testingPassword123"
var testingHashedTransaction = [32]byte{'h', 'e', 'j', 's', 'a','b','m','m','h', 'e', 'j', 's', 'a','b','m','m','h', 'e', 'j', 's', 'a','b','m','m','h', 'e', 'j', 's', 'a','b','m','m'}

func TestCreateWallet(t *testing.T) {
	testAcc := "testAccount1"
	newAddr, _ := CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	testWallet.Delete()
	if newAddr != testWallet.Address {
		t.Errorf("new account address doesn't match the account address returned by GetAddress\n")
	}
}

func TestAccessWalletWithWrongPassword(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)

	testWallet, _ := AccessWallet(testAcc, testingPassword)

	val, err := testWallet.GetPrivateKey(testingPassword+"_wrong")

	testWallet.Delete()
	if val != nil || err == nil {
		t.Errorf("expected to fail when given a wrong password, but didn't\n")
	}
}

func TestAccessWalletWithCorrectPassword(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	val, err := testWallet.GetPrivateKey(testingPassword)
	testWallet.Delete()

	if val == nil || err != nil {
		
		t.Errorf("expected to be able to get private key with correct password, but couldn't\n")
	}
}


func TestSignTransaction(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	signature1, err := testWallet.SignTransaction(testingPassword, testingHashedTransaction)

	if err != nil {
		testWallet.Delete()
		t.Errorf("Failed to sign signature1\n")
	}

	signature2, err := testWallet.SignTransaction(testingPassword, testingHashedTransaction)
	testWallet.Delete()

	if err != nil {
		t.Errorf("Failed to sign signature2\n")
	}

	if bytes.Compare(signature1, signature2) != 0  {
		t.Errorf("Signature of same transaction is not identical\n")
	}
}

func TestRetrieveAddressFromSignature(t *testing.T) {
	testAcc := "testAccount1"
	addr, _ := CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	signature, _ := testWallet.SignTransaction(testingPassword, testingHashedTransaction)
	retrievedAddr, err := GetAddressFromSignedTransaction(signature, testingHashedTransaction)

	testWallet.Delete()
	if err != nil {
		t.Errorf("failed to get address from signed transaction\n")
	}

	if addr != retrievedAddr {
		t.Errorf("wallet address and address retrieved from signed transaction doesn't match.\n")
	}
}

func TestSignTwoTransactionWithTwoAccounts(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	testAcc2 := "testAccount2"
	CreateNewWallet(testAcc2, testingPassword)
	testWallet2, _ := AccessWallet(testAcc2, testingPassword)

	signature1_account1, err1 := testWallet.SignTransaction(testingPassword, testingHashedTransaction)
	signature1_account2, err2 := testWallet2.SignTransaction(testingPassword, testingHashedTransaction)

	if err1 != nil || err2 != nil {
		testWallet.Delete()
		testWallet2.Delete()
		t.Errorf("Failed to sign signature1\n")
	}

	signature2_account1, err := testWallet.SignTransaction(testingPassword, testingHashedTransaction)
	signature2_account2, err2 := testWallet2.SignTransaction(testingPassword, testingHashedTransaction)

	testWallet.Delete()
	testWallet2.Delete()
	if err != nil || err2 != nil {
		t.Errorf("Failed to sign signature2\n")
	}

	if bytes.Compare(signature1_account1, signature2_account1) != 0  {
		t.Errorf("Signature of same transaction for wallet one is not identical\n")
	}

	if bytes.Compare(signature1_account2, signature2_account2) != 0  {
		t.Errorf("Signature of same transaction for wallet two is not identical\n")
	}

	if bytes.Compare(signature1_account1, signature1_account2) == 0 || bytes.Compare(signature2_account1, signature2_account2) == 0 {
		t.Errorf("Signature for the two accounts are the same...")
	}
}




