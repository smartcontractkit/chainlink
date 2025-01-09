package deployment

import (
	"fmt"
	"strconv"

	"github.com/gagliardetto/solana-go"
	solRpc "github.com/gagliardetto/solana-go/rpc"

	solCommomUtil "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/common"
)

// SolChain represents a Solana chain.
type SolChain struct {
	// Selectors used as canonical chain identifier.
	Selector uint64
	// RPC cient
	Client *solRpc.Client
	// TODO: raw private key for now, need to replace with a more secure way
	DeployerKey *solana.PrivateKey
	Confirm     func(instructions []solana.Instruction, opts ...solCommomUtil.TxModifier) error
}

func (c SolChain) String() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	return fmt.Sprintf("%s (%d)", chainInfo.ChainName, chainInfo.ChainSelector)
}

func (c SolChain) Name() string {
	chainInfo, err := ChainInfo(c.Selector)
	if err != nil {
		// we should never get here, if the selector is invalid it should not be in the environment
		panic(err)
	}
	if chainInfo.ChainName == "" {
		return strconv.FormatUint(c.Selector, 10)
	}
	return chainInfo.ChainName
}
