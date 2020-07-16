package cltest

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

// OracleTransactor representsthe identity of the oracle address used in cltest
var OracleTransactor *bind.TransactOpts
var privateKey3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea *ecdsa.PrivateKey

func init() {
	var err error
	privateKey3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea, err =
		crypto.HexToECDSA(
			// Extracted from (to abuse notation) keystore.DecryptKey(
			//   []byte(key3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea), "password").D
			"40052315eab136be2379e23d4c7290e74a6212a018a13b1eafae4152df5d4daa")
	if err != nil {
		panic(err)
	}
	OracleTransactor = bind.NewKeyedTransactor(
		privateKey3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea)
}
