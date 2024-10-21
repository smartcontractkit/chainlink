package onchain

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	bmerc "github.com/smartcontractkit/chainlink/e2e/capabilities/components/gethwrappers"
	"math/big"
)

type Input struct {
	URL string  `toml:"url" validate:"required"`
	Out *Output `toml:"out"`
}

type Output struct {
	Addresses []common.Address `toml:"addresses"`
}

func NewProductOnChainDeployment(in *Input) (*Output, error) {
	if in.Out != nil && framework.UseCache() {
		return in.Out, nil
	}

	// deploy your contracts here, example

	c, err := seth.NewClientBuilder().
		WithRpcUrl(in.URL).
		// TODO: use secrets here from AWSecretsManager later, for now only deployment for Geth/Anvil
		WithPrivateKeys([]string{blockchain.DefaultAnvilPrivateKey}).
		Build()
	if err != nil {
		return nil, err
	}

	addr, tx, _, err := bmerc.DeployBurnMintERC677(
		c.NewTXOpts(),
		c.Client,
		"TestToken",
		"TestToken",
		18,
		big.NewInt(100),
	)
	_, err = bind.WaitMined(context.Background(), c.Client, tx)
	if err != nil {
		return nil, err
	}

	out := &Output{
		// save all the addresses to output, so it can be cached
		Addresses: []common.Address{addr},
	}
	in.Out = out
	return out, nil
}
