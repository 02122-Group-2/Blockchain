package cryptography

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"

	"crypto/elliptic"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func CreateNewWallet(password string) (string, error) {
	ks := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	newAcc, err := ks.NewAccount(password)
	
	if err != nil {
		return "", fmt.Errorf(err.Error())
	}
	return newAcc.Address.Hex(), nil
}

func GetAddress() string {
	ks := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	allAccs := ks.Accounts()
	return allAccs[len(allAccs)-1].Address.Hex()
}

func GetPrivateKey(password string) (*ecdsa.PrivateKey, error) {
	ks := keystore.NewKeyStore("./wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	allAccs := ks.Accounts()

	accountJson, err := ioutil.ReadFile(allAccs[len(allAccs)-1].URL.Path)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	privKey, err := keystore.DecryptKey(accountJson, password)

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	return privKey.PrivateKey, nil
}

// Given a password and a hashed transaction it will use the local wallet in order to sign the transaction.
func SignTransaction(password string, hashedTransaction [32]byte) ([]byte, error) {
	privKey, err := GetPrivateKey(password)
	if err != nil {
		return nil, err
	}
	signature, err := crypto.Sign(hashedTransaction[:], privKey)

	if (err != nil) {
		return nil, err
	}

	return signature, nil
}

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
