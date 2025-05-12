package v1_6

import (
	"fmt"
	"math/big"

	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	ccipchangesethelpers "github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	soltokens "github.com/smartcontractkit/chainlink-ccip/chains/solana/utils/tokens"
)

type ApproveTokenEVMConfig struct {
	ChainSelector		uint64
	TokenAddress		string
	RouterAddress		string
	Amount			*big.Int
}

func ApproveTokenTransferEVMChangeset(e deployment.Environment, cfg ApproveTokenEVMConfig) (cldf.ChangesetOutput, error) {
	err := ccipchangesethelpers.ApproveToken(e, cfg.ChainSelector, cfg.TokenAddress,cfg.RouterAddress,cfg.Amount)

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to approve token transfer: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func ApproveTokenTransferSolChangeset(e deployment.Environment, cfg ApproveTokenEVMConfig) (cldf.ChangesetOutput, error) {
	ixApprove, err := soltokens.TokenApproveChecked(1e9, 9, tokenProgram, deployerWSOL, wSOL, billingSignerPDA, deployer.PublicKey(), []solana.PublicKey{})

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create approve instruction: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}
