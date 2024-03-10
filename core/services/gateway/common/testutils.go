package common

import (
	"crypto/ecdsa"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

type TestNode struct {
	Address    string
	PrivateKey *ecdsa.PrivateKey
}

func NewTestNodes(t *testing.T, n int) []TestNode {
	nodes := make([]TestNode, n)
	for i := 0; i < n; i++ {
		privateKey, err := crypto.GenerateKey()
		require.NoError(t, err)
		address := strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
		nodes[i] = TestNode{Address: address, PrivateKey: privateKey}
	}
	return nodes
}
