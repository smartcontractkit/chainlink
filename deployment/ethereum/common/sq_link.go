package deployment_common

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/smartcontractkit/chainlink/deployment"
	deployment_ethereum "github.com/smartcontractkit/chainlink/deployment/ethereum/extension"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
)

type SqDeployLinkInput struct {
	MintAmount *big.Int
	Amount     *big.Int
	To         common.Address
	chainID    uint64
}

type SqDeployLinkOutput struct {
	Address common.Address
}

var DeployLinkSequence = deployment.NewSequence(
	"v1",
	"Deploy LINK token contract, grants mint and mints some amount to same address",
	func(e deployment.OpEnv, deps deployment_ethereum.EthereumDeps, input SqDeployLinkInput) (SqDeployLinkOutput, error) {

		linkDeployReport, err := deployment.ExecuteOp(e, DeployLinkOp, deps, deployment.EmptyInput{})
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		grantMintInput := GrantLinkInput{
			contractAddress: linkDeployReport.Output.Address,
			To:              deps.Auth.From,
		}

		_, err = deployment.ExecuteOp(e, GrantMintLinkOp, deps, grantMintInput)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		mintInput := MintLinkInput{
			contractAddress: linkDeployReport.Output.Address,
			To:              deps.Auth.From,
			Amount:          input.MintAmount,
		}

		_, err = deployment.ExecuteOp(e, MintLinkOp, deps, mintInput)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		transferInput := TransferLinkInput{
			contractAddress: linkDeployReport.Output.Address,
			To:              input.To,
			Amount:          input.Amount,
		}

		_, err = deployment.ExecuteOp(e, TransferLinkOp, deps, transferInput)
		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		// Sequences Executions will have by default every sub-operation report included in the reporter, no need to return them here, just the info relevant to the sequence
		return SqDeployLinkOutput{Address: linkDeployReport.Output.Address}, nil
	},
)

var DeployLinkWithAutobumpSequence = deployment.NewSequence(
	"v1",
	"Deploy LINK token contract, grants mint and mints some amount to same address using the autobump operation",
	func(e deployment.OpEnv, deps deployment_ethereum.EthereumDeps, input SqDeployLinkInput) (SqDeployLinkOutput, error) {

		// Use memory environment ethereum client to generate the data only
		backend := simulated.NewBackend(types.GenesisAlloc{})
		_, tx, _, _ := link_token.DeployLinkToken(deps.Auth, backend.Client())

		// This is a very explicit example, but Ethereum ops could be built using the Autobump operation
		// Or the send tx feature can simply be a dependency of the operation, and the bumper be abstracted there
		txReport, err := deployment.ExecuteOp(e, deployment_ethereum.SendTxWithGasBumpOp, deps, deployment_ethereum.GasBumpOpInput{
			deployment_ethereum.GasBump{
				RetryLimit:      10,
				RetryIntervalMs: 1000,
				BumpPercentage:  10,
			},
			struct {
				To    *common.Address
				Data  []byte
				Value *big.Int
			}{
				Data: tx.Data(),
			},
		})

		if err != nil {
			return SqDeployLinkOutput{}, err
		}

		return SqDeployLinkOutput{Address: txReport.Output.Address}, nil
	},
)
