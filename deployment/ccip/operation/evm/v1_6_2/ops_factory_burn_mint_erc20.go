package v1_6_2

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_2/factory_burn_mint_erc20"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	opsutil "github.com/smartcontractkit/chainlink/deployment/common/opsutils"
)

type DeployFactoryBurnMintERC20TokenInput struct {
	Name      string         `json:"name"`
	Symbol    string         `json:"symbol"`
	Decimals  uint8          `json:"decimals_"`
	MaxSupply *big.Int       `json:"maxSupply_"`
	PreMint   *big.Int       `json:"preMint"`
	NewOwner  common.Address `json:"newOwner"`
}

var (
	DeployFactoryBurnMintERC20TokenOp = opsutil.NewEVMDeployOperation(
		"FactoryBurnMintERC20",
		semver.MustParse("1.0.0"),
		"Deploys FactoryBurnMintERC20 1.6.2 contract on the specified evm chain",
		shared.FactoryBurnMintERC20Token,
		factory_burn_mint_erc20.FactoryBurnMintERC20MetaData,
		&opsutil.ContractOpts{
			Version:          &deployment.Version1_6_2,
			EVMBytecode:      common.FromHex(factory_burn_mint_erc20.FactoryBurnMintERC20Bin),
			ZkSyncVMBytecode: factory_burn_mint_erc20.ZkBytecode,
		},
		func(input DeployFactoryBurnMintERC20TokenInput) []any {
			return []any{
				input.Name,
				input.Symbol,
				input.Decimals,
				input.MaxSupply,
				input.PreMint,
				input.NewOwner,
			}
		},
	)
)
