package deployment

import (
	"github.com/aptos-labs/aptos-go-sdk"
)

// AptChain represents an Aptos chain.
type AptChain struct {
	Selector       uint64
	Client         aptos.AptosRpcClient
	DeployerSigner aptos.TransactionSigner
}
