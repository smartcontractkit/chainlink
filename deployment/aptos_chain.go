package deployment

import (
	"github.com/aptos-labs/aptos-go-sdk"
)

// AptosChain represents an Aptos chain.
type AptosChain struct {
	Selector       uint64
	Client         aptos.AptosRpcClient
	DeployerSigner aptos.TransactionSigner
}
