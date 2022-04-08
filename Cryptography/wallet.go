package cryptography

import (
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/keystore"
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