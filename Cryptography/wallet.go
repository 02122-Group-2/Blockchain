package cryptography

import (
	shared "blockchain/Shared"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"strings"

	"crypto/elliptic"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"os"
	"path/filepath"
)

// * All functions in file: Magnus, s204509

// Given a password, this function will create a new wallet in the ./wallet folder. It will not delete the old wallets.
func CreateNewWallet(username string, password string) (string, error) {
	hashedUsername := crypto.Keccak256Hash([]byte(username)).Hex()
	ks := keystore.NewKeyStore(filepath.Join(shared.LocalDirToWallets, hashedUsername), keystore.StandardScryptN, keystore.StandardScryptP)
	newAcc, err := ks.NewAccount(password)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return newAcc.Address.Hex(), nil
}

type Account struct {
	Username string           // Stores the username of the wallet
	Address  string           // Stores the address of the wallet
	Wallet   accounts.Account // Stores the actual wallet.
}

func AccessWallet(username string, password string) (Account, error) {
	hashedUsername := crypto.Keccak256Hash([]byte(username)).Hex()
	ks := keystore.NewKeyStore(filepath.Join(shared.LocalDirToWallets, hashedUsername), keystore.StandardScryptN, keystore.StandardScryptP)
	allAccs := ks.Accounts()
	if len(allAccs) == 0 {
		return Account{}, fmt.Errorf("not able to find the account")
	}

	wallet := allAccs[len(allAccs)-1]

	crypto.Keccak256Hash([]byte(username))

	account := Account{
		Username: hashedUsername,
		Address:  wallet.Address.Hex(),
		Wallet:   wallet,
	}

	err := account.tryLogIn(password)

	if err != nil {
		return Account{}, err
	}

	return account, nil
}

// Tries to log in with the password. This is done when accessing the wallet
func (account *Account) tryLogIn(password string) error {
	_, err := account.GetPrivateKey(password)
	return err
}

// Given a password, it gets the private key from the wallet.
func (account *Account) GetPrivateKey(password string) (*ecdsa.PrivateKey, error) {
	accountJson, err := ioutil.ReadFile(account.Wallet.URL.Path)
	if err != nil {
		return nil, err
	}

	privKey, err := keystore.DecryptKey(accountJson, password)

	if err != nil {
		return nil, err
	}

	return privKey.PrivateKey, nil
}

// Given a password and a hashed transaction it will use the wallet in order to sign the transaction.
func (account *Account) SignTransaction(password string, hashedTransaction [32]byte) ([]byte, error) {
	privKey, err := account.GetPrivateKey(password)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(hashedTransaction[:], privKey)

	if err != nil {
		return nil, err
	}

	return signature, nil
}

// Deletes the wallet - should require some verification before doing this
func (account *Account) Delete(username string, password string) error {

	hashedUsername := crypto.Keccak256Hash([]byte(username)).Hex()

	if strings.Compare(hashedUsername, account.Username) != 0 {
		return fmt.Errorf("Cannot delete wallet. Username is incorrect")
	}

	err := account.tryLogIn(password)
	if err != nil {
		return fmt.Errorf("Cannot delete wallet. Password is incorrect")
	}

	return os.Remove(account.Wallet.URL.Path)
}

//Hard delete the wallet - no required verification. Used for testing
func (account *Account) HardDelete() error {
	return os.Remove(account.Wallet.URL.Path)
}

// Given a signature and a signed transaction, it will return the public address of the signer
// This can be used with the transaction, which was hashed, to ensure the sender of the transaction is the one who signed it,
// as this would result in the value from this function is equal to the transaction sender.
func GetAddressFromSignedTransaction(signature []byte, hashedTransaction [32]byte) (string, error) {
	addr, err := crypto.SigToPub(hashedTransaction[:], signature)
	if err != nil {
		return "", err
	}

	addrBytes := elliptic.Marshal(crypto.S256(), addr.X, addr.Y)
	addrHash := crypto.Keccak256(addrBytes[1:])
	addrString := common.BytesToAddress(addrHash[:]).String()

	return addrString, nil
}
