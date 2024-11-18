package changeset

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/config"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/internal"
)

var _ deployment.ChangeSet[map[uint64]MCMSWithTimelockConfig] = DeployMCMSWithTimelock

type MCMSWithTimelockConfig struct {
	Canceller         config.Config
	Bypasser          config.Config
	Proposer          config.Config
	TimelockExecutors []common.Address
	TimelockMinDelay  *big.Int
}

func DeployMCMSWithTimelock(e deployment.Environment, cfgByChain map[uint64]MCMSWithTimelockConfig) (deployment.ChangesetOutput, error) {
	newAddresses := deployment.NewMemoryAddressBook()
	err := internal.DeployMCMSWithTimelockContractsBatch(
		e.Logger, e.Chains, newAddresses, cfgByChain,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{AddressBook: newAddresses}, nil
}
