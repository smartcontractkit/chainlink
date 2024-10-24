package onchain

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	bmerc "github.com/smartcontractkit/chainlink/e2e/capabilities/components/gethwrappers"
	"math/big"
	"os"
)

type Input struct {
	URL string  `toml:"url"`
	Out *Output `toml:"out"`
}

type Output struct {
	UseCache  bool             `toml:"use_cache"`
	Addresses []common.Address `toml:"addresses"`
}

func NewProductOnChainDeployment(in *Input) (*Output, error) {
	if in.Out.UseCache {
		return in.Out, nil
	}

	// deploy your contracts here, example

	c, err := seth.NewClientBuilder().
		WithRpcUrl(in.URL).
		WithPrivateKeys([]string{os.Getenv("PRIVATE_KEY")}).
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
		UseCache: true,
		// save all the addresses to output, so it can be cached
		Addresses: []common.Address{addr},
	}
	in.Out = out
	return out, nil
}
