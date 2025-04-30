package deployment

import (
	"github.com/aptos-labs/aptos-go-sdk"

	"github.com/smartcontractkit/chainlink-aptos/bindings/bind"

	ccip_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
)

// AptosChain represents an Aptos chain.
type AptosChain struct {
	Selector       uint64
	Client         aptos.AptosRpcClient
	DeployerSigner aptos.TransactionSigner
	URL            string

	Confirm func(txHash string, opts ...any) error
}

func (c AptosChain) GetOfframpDynamicConfig(ccipAddress aptos.AccountAddress) (module_offramp.DynamicConfig, error) {
	offrampBind := ccip_offramp.Bind(ccipAddress, c.Client)
	return offrampBind.Offramp().GetDynamicConfig(&bind.CallOpts{})
}
