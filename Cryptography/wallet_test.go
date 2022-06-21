package cryptography

import (
	"bytes"
	"testing"
)

var testingPassword = "testingPassword123"
var testingHashedTransaction = [32]byte{'h', 'e', 'j', 's', 'a', 'b', 'm', 'm', 'h', 'e', 'j', 's', 'a', 'b', 'm', 'm', 'h', 'e', 'j', 's', 'a', 'b', 'm', 'm', 'h', 'e', 'j', 's', 'a', 'b', 'm', 'm'}
var testingHashedTransaction2 = [32]byte{'k', 'e', 'j', 's', 'a', 'b', 'm', 'm', 'h', 'e', 'j', 's', 'a', 'b', 'm', 'm', 'h', 'e', 'j', 's', 'a', 'b', 'm', 'm', 'h', 'e', 'j', 's', 'a', 'b', 'm', 'm'}

// * Magnus, s204509
func TestCreateWallet(t *testing.T) {
	testAcc := "testAccount1"
	newAddr, _ := CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	testWallet.HardDelete()
	if newAddr != testWallet.Address {
		t.Errorf("new account address doesn't match the account address returned by GetAddress\n")
	}
}

// * Magnus, s204509
func TestCreateExistingWallet(t *testing.T) {
	testAcc := "testAccount1"
	newAddr1, _ := CreateNewWallet(testAcc, testingPassword)
	testWallet1, _ := AccessWallet(testAcc, testingPassword)

	newAddr2, _ := CreateNewWallet(testAcc, testingPassword)
	testWallet2, _ := AccessWallet(testAcc, testingPassword)

	testWallet1.HardDelete()
	testWallet2.HardDelete()
	if newAddr1 != testWallet1.Address {
		t.Errorf("new account address for account 1 doesn't match the account address returned by GetAddress\n")
	}
	if newAddr2 != testWallet2.Address {
		t.Errorf("new account address for account 2 doesn't match the account address returned by GetAddress\n")
	}
	if newAddr1 == newAddr2 {
		t.Errorf("The two accounts shouldn't have the same address\n")
	}
}

// * Emilie, s204471
func TestDeleteWallet(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	err := testWallet.Delete(testAcc, testingPassword)
	if err != nil {
		t.Log("Expected to delete account but didn't \n", err)
		t.Fail()
	}
}

// * Emilie, s204471
func TestDeleteWalletWithWrongPassword(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	err := testWallet.Delete(testAcc, testingPassword+"_wrong")

	testWallet.HardDelete()
	if err == nil {
		t.Log("Expected to not delete account when given wrong password, but did \n", err)
		t.Fail()
	}
}

// * Emilie, s204471
func TestAccessWalletWithWrongPassword(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)

	testWallet, _ := AccessWallet(testAcc, testingPassword)

	val, err := testWallet.GetPrivateKey(testingPassword + "_wrong")

	testWallet.HardDelete()
	if val != nil || err == nil {
		t.Errorf("expected to fail when given a wrong password, but didn't\n")
	}
}

// * Emilie, s204471
func TestAccessWalletWithWrongUsername(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)

	testWallet, err := AccessWallet(testAcc+"_wrong", testingPassword)

	testWallet.HardDelete()
	if err == nil {
		t.Errorf("expected to fail when given a wrong username, but didn't\n")
	}
}

// * Emilie, s204471
func TestAccessWalletWithCorrectPassword(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	val, err := testWallet.GetPrivateKey(testingPassword)
	testWallet.HardDelete()

	if val == nil || err != nil {

		t.Errorf("expected to be able to get private key with correct password, but couldn't\n")
	}
}

// * Magnus, s204509
func TestSignTransaction(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	signature1, err := testWallet.SignTransaction(testingPassword, testingHashedTransaction)

	if err != nil {
		testWallet.HardDelete()
		t.Errorf("Failed to sign signature1\n")
	}

	signature2, err := testWallet.SignTransaction(testingPassword, testingHashedTransaction)
	testWallet.HardDelete()

	if err != nil {
		t.Errorf("Failed to sign signature2\n")
	}

	if bytes.Compare(signature1, signature2) != 0 {
		t.Errorf("Signature of same transaction is not identical\n")
	}
}

// * Emilie, s204471
func TestSignTransactionWithWrongPassword(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	_, err := testWallet.SignTransaction(testingPassword+"_wrong", testingHashedTransaction)
	testWallet.HardDelete()
	if err == nil {
		t.Log("Expected to fail signing transaction with wrong password but didn't\n")
		t.Fail()
	}
}

// * Emilie, s204471
func TestSignTransactionWhereSenderIsntTheSigner(t *testing.T) {
	//Sending account1
	testAcc := "testAccount1"
	sendAddr, _ := CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	//Signing account2
	testAcc2 := "testAccount2"
	CreateNewWallet(testAcc2, testingPassword)
	testWallet2, _ := AccessWallet(testAcc2, testingPassword)

	//Sign with account 2
	signature, _ := testWallet2.SignTransaction(testingPassword, testingHashedTransaction)

	//Retrive address
	retrievedAddr, _ := GetAddressFromSignedTransaction(signature, testingHashedTransaction)

	testWallet.HardDelete()
	testWallet2.HardDelete()

	if sendAddr == retrievedAddr {
		t.Log("Signing address and sending address match, but they shouldn't.\n")
		t.Fail()
	}

}

// * Magnus, s204509
func TestRetrieveAddressFromSignature(t *testing.T) {
	testAcc := "testAccount1"
	addr, _ := CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	signature, _ := testWallet.SignTransaction(testingPassword, testingHashedTransaction)
	retrievedAddr, err := GetAddressFromSignedTransaction(signature, testingHashedTransaction)

	testWallet.HardDelete()
	if err != nil {
		t.Errorf("failed to get address from signed transaction\n")
	}

	if addr != retrievedAddr {
		t.Errorf("wallet address and address retrieved from signed transaction doesn't match.\n")
	}
}

// * Magnus, s204509
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
		testWallet.HardDelete()
		testWallet2.HardDelete()
		t.Errorf("Failed to sign signature1\n")
	}

	signature2_account1, err := testWallet.SignTransaction(testingPassword, testingHashedTransaction)
	signature2_account2, err2 := testWallet2.SignTransaction(testingPassword, testingHashedTransaction)

	testWallet.HardDelete()
	testWallet2.HardDelete()
	if err != nil || err2 != nil {
		t.Errorf("Failed to sign signature2\n")
	}

	if bytes.Compare(signature1_account1, signature2_account1) != 0 {
		t.Errorf("Signature of same transaction for wallet one is not identical\n")
	}

	if bytes.Compare(signature1_account2, signature2_account2) != 0 {
		t.Errorf("Signature of same transaction for wallet two is not identical\n")
	}

	if bytes.Compare(signature1_account1, signature1_account2) == 0 || bytes.Compare(signature2_account1, signature2_account2) == 0 {
		t.Errorf("Signature for the two accounts are the same...")
	}
}

// * Emilie, s204471
func TestSignaturesFromSameWalletDifferForTransactions(t *testing.T) {
	testAcc := "testAccount1"
	CreateNewWallet(testAcc, testingPassword)
	testWallet, _ := AccessWallet(testAcc, testingPassword)

	//Sign transaction 1
	signature1, _ := testWallet.SignTransaction(testingPassword, testingHashedTransaction)

	//Sign transaction 2
	signature2, _ := testWallet.SignTransaction(testingPassword, testingHashedTransaction2)

	testWallet.HardDelete()
	if bytes.Compare(signature1, signature2) == 0 {
		t.Log("Signatures are equal for two different transactions but shouldn't be")
		t.Fail()
	}
}
