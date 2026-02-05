package main

import (
	"os"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func main() {
	// TODO: more keys to choose
	k, _ := p2pkey.NewV2()
	b, _ := k.ToEncryptedJSON("password", utils.ScryptParams{N: 2, P: 1})
	os.WriteFile("out.txt", b, 0o644)
}
