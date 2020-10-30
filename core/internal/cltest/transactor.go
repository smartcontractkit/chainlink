package cltest

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

// OracleTransactor represents the identity of the oracle address used in cltest
var OracleTransactor *bind.TransactOpts

func init() {
	var err error
	k, err := keystore.DecryptKey([]byte(DefaultKeyJSON), "password")
	if err != nil {
		panic(err)
	}
	OracleTransactor = bind.NewKeyedTransactor(k.PrivateKey)
}
