package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

// ERC20 transfer method selector (transfer(address,uint256))
var erc20TransferSelector = []byte{0xa9, 0x05, 0x9c, 0xbb}

// encodeERC20TransferCalldata returns ABI-encoded transfer(to, amount) for standard ERC20
func encodeERC20TransferCalldata(tr types.ERC20Transfer) ([]byte, error) {
	uint256Ty, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, err
	}
	addressTy, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, err
	}
	args := abi.Arguments{
		{Type: addressTy},
		{Type: uint256Ty},
	}
	packed, err := args.Pack(
		common.HexToAddress(tr.Payee),
		tr.Amount,
	)
	if err != nil {
		return nil, err
	}
	return append(erc20TransferSelector, packed...), nil
}

var TransferERC20Changeset cldf.ChangeSetV2[types.TransferERC20Config] = transferERC20Changeset{}

type transferERC20Changeset struct{}

func (t transferERC20Changeset) VerifyPreconditions(e cldf.Environment, cfg types.TransferERC20Config) error {
	return ValidateTransferERC20Config(e, cfg)
}

func (t transferERC20Changeset) Apply(e cldf.Environment, cfg types.TransferERC20Config) (cldf.ChangesetOutput, error) {
	lggr := e.Logger

	lggr.Infow("Starting ERC20 transfer from timelock",
		"chain", cfg.ChainSelector,
		"timelock_id", cfg.TimelockIdentifier,
		"transfers", len(cfg.Transfers),
		"description", cfg.Description)

	evmChains := e.BlockChains.EVMChains()
	chain, exists := evmChains[cfg.ChainSelector]
	if !exists {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", cfg.ChainSelector)
	}

	timelockAddr, err := GetContractAddressWithQualifier(e.DataStore, cfg.ChainSelector, commontypes.RBACTimelock, cfg.TimelockIdentifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("timelock not found for chain %d: %w", cfg.ChainSelector, err)
	}

	proposerAddr, err := GetContractAddressWithQualifier(e.DataStore, cfg.ChainSelector, commontypes.ProposerManyChainMultisig, cfg.TimelockIdentifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("proposer not found for chain %d: %w", cfg.ChainSelector, err)
	}

	var transactions []mcmstypes.Transaction
	for i, tr := range cfg.Transfers {
		data, err := encodeERC20TransferCalldata(tr)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("transfer %d: failed to encode calldata: %w", i, err)
		}

		tx, err := proposalutils.TransactionForChain(
			cfg.ChainSelector,
			tr.Token,
			data,
			nil,
			"ERC20Transfer",
			[]string{"vault", "erc20-transfer"},
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("transfer %d: %w", i, err)
		}
		transactions = append(transactions, tx)
	}

	batches := []mcmstypes.BatchOperation{{
		ChainSelector: mcmstypes.ChainSelector(cfg.ChainSelector),
		Transactions:  transactions,
	}}

	description := cfg.Description
	if description == "" {
		description = "Vault ERC20 Transfer"
	}

	mcmsConfig := proposalutils.TimelockConfig{MinDelay: 0}
	if cfg.MCMSConfig != nil {
		mcmsConfig = *cfg.MCMSConfig
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		map[uint64]string{cfg.ChainSelector: timelockAddr},
		map[uint64]string{cfg.ChainSelector: proposerAddr},
		map[uint64]mcmsevmsdk.Inspector{cfg.ChainSelector: mcmsevmsdk.NewInspector(chain.Client)},
		batches,
		description,
		mcmsConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}

	lggr.Infow("ERC20 transfer proposal built",
		"chain", cfg.ChainSelector,
		"transfers", len(cfg.Transfers))

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
	}, nil
}
