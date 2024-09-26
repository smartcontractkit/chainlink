package feeds

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
)

type FeedConfig struct {
	MaximumGasPrice         uint32
	ReasonableGasPrice      uint32
	MicroLinkPerEth         uint32
	LinkGweiPerObservation  uint32
	LinkGweiPerTransmission uint32
	Link                    common.Address
	MinAnswer               *big.Int
	MaxAnswer               *big.Int
}

type Contracts interface {
	*offchainaggregator.OffchainAggregator
}

type ContractDeploy[C Contracts] struct {
	// We just keep all the deploy return values
	// since some will be empty if there's an error.
	Address  common.Address
	Contract C
	Tx       *types.Transaction
	Tv       deployment.TypeAndVersion
	Err      error
}

// TODO: pull up to general deployment pkg somehow
// without exposing all product specific contracts?
func deployContract[C Contracts](
	lggr logger.Logger,
	chain deployment.Chain,
	addressBook deployment.AddressBook,
	deploy func(chain deployment.Chain) ContractDeploy[C],
) (*ContractDeploy[C], error) {
	contractDeploy := deploy(chain)
	if contractDeploy.Err != nil {
		lggr.Errorw("Failed to deploy contract", "err", contractDeploy.Err)
		return nil, contractDeploy.Err
	}
	_, err := chain.Confirm(contractDeploy.Tx)
	if err != nil {
		lggr.Errorw("Failed to confirm deployment", "err", err)
		return nil, err
	}
	err = addressBook.Save(chain.Selector, contractDeploy.Address.String(), contractDeploy.Tv)
	if err != nil {
		lggr.Errorw("Failed to save contract address", "err", err)
		return nil, err
	}
	return &contractDeploy, nil
}

func DeployFeed(e deployment.Environment, chainSel uint64, cfg FeedConfig) (deployment.AddressBook, error) {
	ab := deployment.NewMemoryAddressBook()
	_, err := deployContract(e.Logger, e.Chains[chainSel], ab,
		func(chain deployment.Chain) ContractDeploy[*offchainaggregator.OffchainAggregator] {
			receiverAddr, tx, receiver, err2 := offchainaggregator.DeployOffchainAggregator(
				chain.DeployerKey,
				chain.Client,
				cfg.MaximumGasPrice,
				cfg.ReasonableGasPrice,
				cfg.MicroLinkPerEth,
				cfg.LinkGweiPerObservation,
				cfg.LinkGweiPerTransmission,
				cfg.Link,
				cfg.MinAnswer,
				cfg.MaxAnswer,
				common.HexToAddress("0x1"),
				common.HexToAddress("0x2"),
				8,
				"test",
			)
			return ContractDeploy[*offchainaggregator.OffchainAggregator]{
				receiverAddr, receiver, tx, deployment.NewTypeAndVersion("OffchainAggregator", deployment.Version1_0_0), err2,
			}
		})
	if err != nil {
		e.Logger.Errorw("Failed to deploy receiver", "err", err)
		return ab, err
	}
	return ab, nil
}
