package changeset

import (
	"fmt"
	"math/big"
	"reflect"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_remote"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

type RMNNopConfig struct {
	NodeIndex           uint64
	OffchainPublicKey   [32]byte
	EVMOnChainPublicKey common.Address
	PeerId              p2pkey.PeerID
}

func (c RMNNopConfig) ToRMNHomeNode() rmn_home.RMNHomeNode {
	return rmn_home.RMNHomeNode{
		PeerId:            c.PeerId,
		OffchainPublicKey: c.OffchainPublicKey,
	}
}

func (c RMNNopConfig) ToRMNRemoteSigner() rmn_remote.RMNRemoteSigner {
	return rmn_remote.RMNRemoteSigner{
		OnchainPublicKey: c.EVMOnChainPublicKey,
		NodeIndex:        c.NodeIndex,
	}
}

func (c RMNNopConfig) SetBit(bitmap *big.Int, value bool) {
	if value {
		bitmap.SetBit(bitmap, int(c.NodeIndex), 1)
	} else {
		bitmap.SetBit(bitmap, int(c.NodeIndex), 0)
	}
}

func getDeployer(e deployment.Environment, chain uint64, mcmConfig *MCMSConfig) *bind.TransactOpts {
	if mcmConfig == nil {
		return e.Chains[chain].DeployerKey
	}

	return deployment.SimTransactOpts()
}

type MCMSConfig struct {
	MinDelay time.Duration
}

type SetRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
	RMNStaticConfig   rmn_home.RMNHomeStaticConfig
	RMNDynamicConfig  rmn_home.RMNHomeDynamicConfig
	DigestToOverride  [32]byte
	MCMSConfig        *MCMSConfig
}

func (c SetRMNHomeCandidateConfig) Validate(state CCIPOnChainState) error {
	err := deployment.IsValidChainSelector(c.HomeChainSelector)
	if err != nil {
		return err
	}

	if len(c.RMNDynamicConfig.OffchainConfig) != 0 {
		return fmt.Errorf("RMNDynamicConfig.OffchainConfig must be empty")
	}
	if len(c.RMNStaticConfig.OffchainConfig) != 0 {
		return fmt.Errorf("RMNStaticConfig.OffchainConfig must be empty")
	}

	if len(c.RMNStaticConfig.Nodes) > 256 {
		return fmt.Errorf("RMNStaticConfig.Nodes must be less than 256")
	}

	var (
		peerIds            = make(map[[32]byte]struct{})
		offchainPublicKeys = make(map[[32]byte]struct{})
	)

	for _, node := range c.RMNStaticConfig.Nodes {
		if _, exists := peerIds[node.PeerId]; exists {
			return fmt.Errorf("peerId %x is duplicated", node.PeerId)
		}
		peerIds[node.PeerId] = struct{}{}

		if _, exists := offchainPublicKeys[node.OffchainPublicKey]; exists {
			return fmt.Errorf("offchainPublicKey %x is duplicated", node.OffchainPublicKey)
		}
		offchainPublicKeys[node.OffchainPublicKey] = struct{}{}
	}
	rmnHome := state.Chains[c.HomeChainSelector].RMNHome

	if rmnHome == nil {
		return fmt.Errorf("RMNHome not found for chain %d", c.HomeChainSelector)
	}

	currentDigest, err := rmnHome.GetCandidateDigest(nil)
	if err != nil {
		return fmt.Errorf("failed to get RMNHome candidate digest: %w", err)
	}

	if currentDigest != c.DigestToOverride {
		return fmt.Errorf("current digest (%x) does not match digest to override (%x)", currentDigest[:], c.DigestToOverride[:])
	}

	return nil
}

type PromoteRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
	DigestToPromote   [32]byte
	MCMSConfig        *MCMSConfig
}

func (c PromoteRMNHomeCandidateConfig) Validate(state CCIPOnChainState) error {
	err := deployment.IsValidChainSelector(c.HomeChainSelector)
	if err != nil {
		return err
	}

	rmnHome := state.Chains[c.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return fmt.Errorf("RMNHome not found for chain %d", c.HomeChainSelector)
	}

	currentCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	if err != nil {
		return fmt.Errorf("failed to get RMNHome candidate digest: %w", err)
	}

	if currentCandidateDigest != c.DigestToPromote {
		return fmt.Errorf("current digest (%x) does not match digest to promote (%x)", currentCandidateDigest[:], c.DigestToPromote[:])
	}

	return nil
}

// NewSetRMNHomeCandidateConfigChangeset creates a changeset to set the RMNHome candidate config
// DigestToOverride is the digest of the current candidate config that the new config will override
// StaticConfig contains the list of nodes with their peerIDs (found in their rageproxy keystore) and offchain public keys (found in the RMN keystore)
// DynamicConfig contains the list of source chains with their chain selectors, f value and the bitmap of the nodes that are oberver for each source chain
// The bitmap is a 256 bit array where each bit represents a node. If the bit matching the index of the node in the static config is set it means that the node is an observer
func NewSetRMNHomeCandidateConfigChangeset(e deployment.Environment, config SetRMNHomeCandidateConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = config.Validate(state)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	homeChain, ok := e.Chains[config.HomeChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found", config.HomeChainSelector)
	}

	rmnHome := state.Chains[config.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("RMNHome not found for chain %s", homeChain.String())
	}

	deployer := getDeployer(e, config.HomeChainSelector, config.MCMSConfig)
	setCandidateTx, err := rmnHome.SetCandidate(deployer, config.RMNStaticConfig, config.RMNDynamicConfig, config.DigestToOverride)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("build RMNHome set candidate calldata for chain %s: %w", homeChain.String(), err)
	}

	if config.MCMSConfig == nil {
		chain := e.Chains[config.HomeChainSelector]
		_, err := chain.Confirm(setCandidateTx)

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm tx for chain %s: %w", homeChain.String(), deployment.MaybeDataErr(err))
		}

		return deployment.ChangesetOutput{}, nil
	}

	op := mcms.Operation{
		To:    rmnHome.Address(),
		Data:  setCandidateTx.Data(),
		Value: big.NewInt(0),
	}

	batches := []timelock.BatchChainOperation{
		{
			ChainIdentifier: mcms.ChainIdentifier(config.HomeChainSelector),
			Batch:           []mcms.Operation{op},
		},
	}

	timelocksPerChain := buildTimelockAddressPerChain(e, state)

	proposerMCMSes := buildProposerPerChain(e, state)

	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to set candidate config",
		config.MCMSConfig.MinDelay,
	)

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

func NewPromoteCandidateConfigChangeset(e deployment.Environment, config PromoteRMNHomeCandidateConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = config.Validate(state)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	homeChain, ok := e.Chains[config.HomeChainSelector]

	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found", config.HomeChainSelector)
	}

	rmnHome := state.Chains[config.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("RMNHome not found for chain %s", homeChain.String())
	}

	currentCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get RMNHome candidate digest for chain %s: %w", homeChain.String(), err)
	}

	currentActiveDigest, err := rmnHome.GetActiveDigest(nil)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get RMNHome active digest for chain %s: %w", homeChain.String(), err)
	}

	deployer := getDeployer(e, config.HomeChainSelector, config.MCMSConfig)
	promoteCandidateTx, err := rmnHome.PromoteCandidateAndRevokeActive(deployer, currentCandidateDigest, currentActiveDigest)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("get call data to promote RMNHome candidate digest for chain %s: %w", homeChain.String(), err)
	}

	if config.MCMSConfig == nil {
		chain := e.Chains[config.HomeChainSelector]
		_, err := chain.Confirm(promoteCandidateTx)

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm tx for chain %s: %w", homeChain.String(), deployment.MaybeDataErr(err))
		}

		return deployment.ChangesetOutput{}, nil
	}

	op := mcms.Operation{
		To:    rmnHome.Address(),
		Data:  promoteCandidateTx.Data(),
		Value: big.NewInt(0),
	}

	batches := []timelock.BatchChainOperation{
		{
			ChainIdentifier: mcms.ChainIdentifier(config.HomeChainSelector),
			Batch:           []mcms.Operation{op},
		},
	}

	timelocksPerChain := buildTimelockAddressPerChain(e, state)

	proposerMCMSes := buildProposerPerChain(e, state)

	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to promote candidate config",
		config.MCMSConfig.MinDelay,
	)

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal for chain %s: %w", homeChain.String(), err)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}

func buildTimelockPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*commonchangeset.TimelockExecutionContracts {
	timelocksPerChain := make(map[uint64]*commonchangeset.TimelockExecutionContracts)
	for _, chain := range e.Chains {
		timelocksPerChain[chain.Selector] = &commonchangeset.TimelockExecutionContracts{
			Timelock:  state.Chains[chain.Selector].Timelock,
			CallProxy: state.Chains[chain.Selector].CallProxy,
		}
	}
	return timelocksPerChain
}

func buildTimelockAddressPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]common.Address {
	timelocksPerChain := buildTimelockPerChain(e, state)
	timelockAddressPerChain := make(map[uint64]common.Address)
	for chain, timelock := range timelocksPerChain {
		timelockAddressPerChain[chain] = timelock.Timelock.Address()
	}
	return timelockAddressPerChain
}

func buildProposerPerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*gethwrappers.ManyChainMultiSig {
	proposerPerChain := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for _, chain := range e.Chains {
		proposerPerChain[chain.Selector] = state.Chains[chain.Selector].ProposerMcm
	}
	return proposerPerChain
}

func buildRMNRemotePerChain(e deployment.Environment, state CCIPOnChainState) map[uint64]*rmn_remote.RMNRemote {
	timelocksPerChain := make(map[uint64]*rmn_remote.RMNRemote)
	for _, chain := range e.Chains {
		timelocksPerChain[chain.Selector] = state.Chains[chain.Selector].RMNRemote
	}
	return timelocksPerChain
}

type SetRMNRemoteConfig struct {
	HomeChainSelector uint64
	Signers           []rmn_remote.RMNRemoteSigner
	F                 uint64
	MCMSConfig        *MCMSConfig
}

func (c SetRMNRemoteConfig) Validate() error {
	err := deployment.IsValidChainSelector(c.HomeChainSelector)
	if err != nil {
		return err
	}

	for i := 0; i < len(c.Signers)-1; i++ {
		if c.Signers[i].NodeIndex >= c.Signers[i+1].NodeIndex {
			return fmt.Errorf("signers must be in ascending order of nodeIndex")
		}
	}

	if len(c.Signers) < 2*int(c.F)+1 {
		return fmt.Errorf("signers count must greater than or equal to %d", 2*c.F+1)
	}

	return nil
}

func NewSetRMNRemoteConfigChangeset(e deployment.Environment, config SetRMNRemoteConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	lggr := e.Logger

	err = config.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}

	homeChain, ok := e.Chains[config.HomeChainSelector]

	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("chain %d not found", config.HomeChainSelector)
	}

	rmnHome := state.Chains[config.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("RMNHome not found for chain %s", homeChain.String())
	}

	activeConfig, err := rmnHome.GetActiveDigest(nil)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get RMNHome active digest for chain %s: %w", homeChain.String(), err)
	}

	rmnRemotePerChain := buildRMNRemotePerChain(e, state)
	batches := make([]timelock.BatchChainOperation, 0)
	for chain, remote := range rmnRemotePerChain {
		if remote == nil {
			continue
		}

		currentVersionConfig, err := remote.GetVersionedConfig(nil)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to get RMNRemote config for chain %s: %w", e.Chains[chain].String(), err)
		}

		newConfig := rmn_remote.RMNRemoteConfig{
			RmnHomeContractConfigDigest: activeConfig,
			Signers:                     config.Signers,
			F:                           config.F,
		}

		if reflect.DeepEqual(currentVersionConfig.Config, newConfig) {
			lggr.Infow("RMNRemote config already up to date", "chain", e.Chains[chain].String())
			continue
		}

		deployer := getDeployer(e, chain, config.MCMSConfig)
		tx, err := remote.SetConfig(deployer, newConfig)

		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("build call data to set RMNRemote config for chain %s: %w", e.Chains[chain].String(), err)
		}

		if config.MCMSConfig == nil {
			_, err := e.Chains[chain].Confirm(tx)

			if err != nil {
				return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm tx for chain %s: %w", e.Chains[chain].String(), deployment.MaybeDataErr(err))
			}
		}

		op := mcms.Operation{
			To:    remote.Address(),
			Data:  tx.Data(),
			Value: big.NewInt(0),
		}

		batch := timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(chain),
			Batch:           []mcms.Operation{op},
		}

		batches = append(batches, batch)
	}

	if config.MCMSConfig == nil {
		return deployment.ChangesetOutput{}, nil
	}

	timelocksPerChain := buildTimelockAddressPerChain(e, state)

	proposerMCMSes := buildProposerPerChain(e, state)

	prop, err := proposalutils.BuildProposalFromBatches(
		timelocksPerChain,
		proposerMCMSes,
		batches,
		"proposal to promote candidate config",
		config.MCMSConfig.MinDelay,
	)

	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal for chain %s: %w", homeChain.String(), err)
	}

	return deployment.ChangesetOutput{
		Proposals: []timelock.MCMSWithTimelockProposal{*prop},
	}, nil
}
