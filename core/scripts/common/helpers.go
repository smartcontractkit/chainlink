package common

import (
	"crypto/ecdsa"
	"encoding/hex"
	"flag"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseArgs(flagSet *flag.FlagSet, args []string, requiredArgs ...string) {
	PanicErr(flagSet.Parse(args))
	seen := map[string]bool{}
	argValues := map[string]string{}
	flagSet.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
		argValues[f.Name] = f.Value.String()
	})
	for _, req := range requiredArgs {
		if !seen[req] {
			panic(fmt.Errorf("missing required -%s argument/flag", req))
		}
	}
}

// GetAccount returns the go-ethereum abstraction of an account on an EVM chain.
func GetAccount(privateKey string, chainID *big.Int) (*bind.TransactOpts, error) {
	b, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	d := new(big.Int).SetBytes(b)
	pubKeyX, pubKeyY := crypto.S256().ScalarBaseMult(d.Bytes())
	privKey := ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: crypto.S256(),
			X:     pubKeyX,
			Y:     pubKeyY,
		},
		D: d,
	}

	return bind.NewKeyedTransactorWithChainID(&privKey, chainID)
}
