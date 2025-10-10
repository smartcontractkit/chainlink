package v1_5_1

import (
	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_1/token_pool_factory"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type DeployTokenPoolFactoryInput struct {
	ChainSelector      uint64
	TokenAdminRegistry common.Address
	// RegistryModule1_6Addresses indicates which registry module to use for each chain.
	// If the chain only has one 1.6.0 registry module, you do not need to specify it here.
	RegistryModule1_6Addresses common.Address
	RMNProxy                   common.Address
	Router                     common.Address
}

var (
	DeployTokenPoolFactoryOp = opsutil.NewEVMDeployOperation(
		"DeployTokenPoolFactory",
		semver.MustParse("1.0.0"),
		"Deploys TokenPoolFactory contract on the specified evm chain",
		shared.TokenPoolFactory,
		token_pool_factory.TokenPoolFactoryMetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_5_1,
			EVMBytecode:      common.FromHex(token_pool_factory.TokenPoolFactoryBin),
			ZkSyncVMBytecode: token_pool_factory.ZkBytecode,
		},
		func(input DeployTokenPoolFactoryInput) []any {
			return []any{
				input.TokenAdminRegistry,
				input.RegistryModule1_6Addresses,
				input.RMNProxy,
				input.Router,
			}
		},
	)
)
