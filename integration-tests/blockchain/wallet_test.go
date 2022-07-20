package blockchain_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/integration-tests/blockchain"
	"github.com/stretchr/testify/require"
)

// Publicly available private key that is used as default in hardhat, geth, etc...
var key = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// Address of the key above
var address = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

func TestWallet(t *testing.T) {
	wallet, err := blockchain.NewEthereumWallet(key)
	require.NoError(t, err)

	require.Equal(t, address, wallet.Address(), "Address of key '%s' not as expected", key)
	require.Equal(t, key, wallet.PrivateKey(), "Private key not as expected")
	require.Equal(t, key, wallet.RawPrivateKey(), "Private key not as expected")
}
