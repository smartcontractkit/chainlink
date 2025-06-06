package example

import (
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

/**
DeployAndMintExampleChangeset demonstrates how to use Operations API to deploy and mint LINK token in EVM by using
a sequence of operations.
*/

var _ cldf.ChangeSetV2[SqDeployLinkInput] = DeployAndMintExampleChangeset{}

// SqDeployLinkInput must be JSON Serializable with no private fields
type SqDeployLinkInput struct {
	MintAmount *big.Int
	Amount     *big.Int
	To         common.Address
	ChainID    uint64
}

// SqDeployLinkOutput must be JSON Serializable with no private fields
type SqDeployLinkOutput struct {
	Address common.Address
}

type EthereumDeps struct {
	Auth  *bind.TransactOpts
	Chain cldf_evm.Chain
	AB    cldf.AddressBook
}

type DeployAndMintExampleChangeset struct{}

func (l DeployAndMintExampleChangeset) VerifyPreconditions(e cldf.Environment, config SqDeployLinkInput) error {
	// perform any preconditions checks here
	return nil
}

func (l DeployAndMintExampleChangeset) Apply(e cldf.Environment, config SqDeployLinkInput) (cldf.ChangesetOutput, error) {
	auth := e.BlockChains.EVMChains()[config.ChainID].DeployerKey
	ab := cldf.NewMemoryAddressBook()

	// build your custom dependencies needed in the sequence/operation
	deps := EthereumDeps{
		Auth:  auth,
		Chain: e.BlockChains.EVMChains()[config.ChainID],
		AB:    ab,
	}

	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, DeployAndMintSequence, deps, config)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReport.ExecutionReports,
	}, nil
}

// DeployAndMintSequence calls 3 operations in sequence:
// 1. DeployLinkOp: Deploys LINK token contract
// 2. GrantMintOp: Grants mint role to the same address
// 3. MintLinkOp: Mints some amount to the same address
// The output of the sequence is the address of the deployed LINK token contract.
var DeployAndMintSequence = operations.NewSequence(
	"deploy-mint-sequence",
	semver.MustParse("1.0.0"),
	"Deploy LINK token contract, grants mint and mints some amount to same address",
	func(b operations.Bundle, deps EthereumDeps, input SqDeployLinkInput) (SqDeployLinkOutput, error) {
		linkDeployReport, err := operations.ExecuteOperation(
			b, DeployLinkOp, deps, operations.EmptyInput{},
			operations.WithRetry[operations.EmptyInput, EthereumDeps](),
		)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		grantMintConfig := GrantMintRoleConfig{
			ContractAddress: linkDeployReport.Output,
			To:              deps.Auth.From,
		}
		_, err = operations.ExecuteOperation(
			b, GrantMintOp, deps, grantMintConfig,
			operations.WithRetry[GrantMintRoleConfig, EthereumDeps](),
		)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		mintConfig := MintLinkConfig{
			ContractAddress: linkDeployReport.Output,
			Amount:          input.MintAmount,
			To:              input.To,
		}
		_, err = operations.ExecuteOperation(
			b, MintLinkOp, deps, mintConfig,
			operations.WithRetry[MintLinkConfig, EthereumDeps](),
		)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		return SqDeployLinkOutput{Address: linkDeployReport.Output}, nil
	},
)
