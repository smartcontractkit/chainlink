package v0_5_0

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/types"
	verifier "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

var SetConfigChangeset = deployment.CreateChangeSet(setConfigLogic, setConfigPrecondition)

type SetConfigConfig struct {
	ConfigsByChain map[uint64][]SetConfig
	MCMSConfig     *changeset.MCMSConfig
}

type SetConfig struct {
	VerifierAddress            common.Address
	ConfigDigest               [32]byte
	Signers                    []common.Address
	F                          uint8
	RecipientAddressesAndProps []verifier.CommonAddressAndWeight
}

type VerifierState struct {
	Verifier *verifier.Verifier
}

func (cfg SetConfigConfig) Validate() error {
	if len(cfg.ConfigsByChain) == 0 {
		return errors.New("ConfigsByChain cannot be empty")
	}
	return nil
}

func setConfigPrecondition(_ deployment.Environment, cc SetConfigConfig) error {
	if err := cc.Validate(); err != nil {
		return fmt.Errorf("invalid SetConfig config: %w", err)
	}
	return nil
}

func setConfigLogic(e deployment.Environment, cfg SetConfigConfig) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	allBatches := []mcmstypes.BatchOperation{}

	var timelockAddressesPerChain map[uint64]string
	var proposerAddressPerChain map[uint64]string
	var inspectorPerChain map[uint64]mcmssdk.Inspector
	if cfg.MCMSConfig != nil {
		timelockAddressesPerChain = make(map[uint64]string)
		proposerAddressPerChain = make(map[uint64]string)
		inspectorPerChain = make(map[uint64]mcmssdk.Inspector)
	}

	for chainSelector, configs := range cfg.ConfigsByChain {
		chain, chainExists := e.Chains[chainSelector]
		if !chainExists {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found in environment", chainSelector)
		}

		batch := mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chainSelector),
			Transactions:  []mcmstypes.Transaction{},
		}

		if cfg.MCMSConfig != nil {
			chainState := state.Chains[chainSelector]
			timelockAddressesPerChain[chainSelector] = chainState.Timelock.Address().String()
			proposerAddressPerChain[chainSelector] = chainState.ProposerMcm.Address().String()
			inspectorPerChain[chainSelector] = evm.NewInspector(chain.Client)
		}

		opts := getTransactOptsSetConfig(e, chainSelector, cfg.MCMSConfig)
		for _, config := range configs {
			confState, err := maybeLoadVerifier(e, chainSelector, config.VerifierAddress.String())
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			tx, err := setConfigOrBuildTx(e, confState.Verifier, config, opts, chain, cfg.MCMSConfig)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}

			op := evm.NewTransaction(
				config.VerifierAddress,
				tx.Data(),
				big.NewInt(0),
				string(types.Verifier),
				[]string{},
			)
			batch.Transactions = append(batch.Transactions, op)
		}

		allBatches = append(allBatches, batch)
	}

	if cfg.MCMSConfig != nil {
		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			e,
			timelockAddressesPerChain,
			proposerAddressPerChain,
			inspectorPerChain,
			allBatches,
			"SetConfig proposal",
			cfg.MCMSConfig.MinDelay,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}

		return deployment.ChangesetOutput{
			MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
		}, nil
	}

	return deployment.ChangesetOutput{}, nil
}

func setConfigOrBuildTx(
	e deployment.Environment,
	verifierContract *verifier.Verifier,
	cfg SetConfig,
	opts *bind.TransactOpts,
	chain deployment.Chain,
	mcmsConfig *changeset.MCMSConfig,
) (*ethTypes.Transaction, error) {
	tx, err := verifierContract.SetConfig(
		opts,
		cfg.ConfigDigest,               // bytes32 configDigest
		cfg.Signers,                    // address[] signers
		cfg.F,                          // uint8 f
		cfg.RecipientAddressesAndProps, // CommonAddressAndWeight[]
	)
	if err != nil {
		return nil, fmt.Errorf("error packing setConfig tx data: %w", err)
	}

	if mcmsConfig == nil {
		if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
			e.Logger.Errorw("Failed to confirm setConfig tx", "chain", chain.String(), "err", err)
			return nil, err
		}
	}
	return tx, nil
}

func maybeLoadVerifier(e deployment.Environment, chainSel uint64, contractAddr string) (*VerifierState, error) {
	chain, ok := e.Chains[chainSel]
	if !ok {
		return nil, fmt.Errorf("chain %d not found", chainSel)
	}

	addresses, err := e.ExistingAddresses.AddressesForChain(chainSel)
	if err != nil {
		return nil, err
	}

	tv, found := addresses[contractAddr]
	if !found {
		return nil, fmt.Errorf("unable to find Verifier contract on chain %s (selector %d)", chain.Name(), chain.Selector)
	}
	if tv.Type != types.Verifier || tv.Version != deployment.Version0_5_0 {
		return nil, fmt.Errorf("unexpected contract type %s for Verifier on chain %s (selector %d)", tv, chain.Name(), chain.Selector)
	}

	conf, err := verifier.NewVerifier(common.HexToAddress(contractAddr), chain.Client)
	if err != nil {
		return nil, err
	}

	return &VerifierState{Verifier: conf}, nil
}

func getTransactOptsSetConfig(e deployment.Environment, chainSel uint64, mcmsConfig *changeset.MCMSConfig) *bind.TransactOpts {
	if mcmsConfig == nil {
		return e.Chains[chainSel].DeployerKey
	}

	return deployment.SimTransactOpts()
}
