package v1_6

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"

	ccipseq "github.com/smartcontractkit/chainlink/deployment/ccip/sequence/evm/v1_6"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/opsutil"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/evm"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var (
	_ cldf.ChangeSet[SetRMNRemoteOnRMNProxyConfig]  = SetRMNRemoteOnRMNProxyChangeset
	_ cldf.ChangeSet[SetRMNHomeCandidateConfig]     = SetRMNHomeCandidateConfigChangeset
	_ cldf.ChangeSet[PromoteRMNHomeCandidateConfig] = PromoteRMNHomeCandidateConfigChangeset
	_ cldf.ChangeSet[ccipseq.SetRMNRemoteConfig]    = SetRMNRemoteConfigChangeset
	_ cldf.ChangeSet[SetRMNHomeDynamicConfigConfig] = SetRMNHomeDynamicConfigChangeset
	_ cldf.ChangeSet[RevokeCandidateConfig]         = RevokeRMNHomeCandidateConfigChangeset
)

type SetRMNRemoteOnRMNProxyConfig struct {
	ChainSelectors []uint64
	MCMSConfig     *proposalutils.TimelockConfig
}

func (c SetRMNRemoteOnRMNProxyConfig) Validate(e cldf.Environment, state stateview.CCIPOnChainState) error {
	for _, chain := range c.ChainSelectors {
		err := cldf.IsValidChainSelector(chain)
		if err != nil {
			return err
		}
		chainState, exists := state.Chains[chain]
		if !exists {
			return fmt.Errorf("chain %d not found in state", chain)
		}
		if chainState.RMNRemote == nil {
			return fmt.Errorf("RMNRemote not found for chain %d", chain)
		}
		if chainState.RMNProxy == nil {
			return fmt.Errorf("RMNProxy not found for chain %d", chain)
		}

		chainEnv := e.BlockChains.EVMChains()[chain]
		if err := commoncs.ValidateOwnership(e.GetContext(), c.MCMSConfig != nil, chainEnv.DeployerKey.From, chainState.Timelock.Address(), chainState.RMNProxy); err != nil {
			return fmt.Errorf("failed to validate ownership of RMNProxy on %s: %w", chainEnv, err)
		}
	}
	return nil
}

func SetRMNRemoteOnRMNProxyChangeset(e cldf.Environment, cfg SetRMNRemoteOnRMNProxyConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	if err := cfg.Validate(e, state); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	timelocks := deployergroup.BuildTimelockAddressPerChain(e, state)

	inspectors := map[uint64]mcmssdk.Inspector{}
	timelockBatch := []mcmstypes.BatchOperation{}
	for _, sel := range cfg.ChainSelectors {
		chain, exists := e.BlockChains.EVMChains()[sel]
		if !exists {
			return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found", sel)
		}

		inspectors[sel], err = proposalutils.McmsInspectorForChain(e, sel)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get mcms inspector for chain %s: %w", chain.String(), err)
		}

		txOpts := chain.DeployerKey
		if cfg.MCMSConfig != nil {
			txOpts = cldf.SimTransactOpts()
		}
		batchOperation, err := setRMNRemoteOnRMNProxyOp(txOpts, chain, state.Chains[sel], cfg.MCMSConfig != nil)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set RMNRemote on RMNProxy for chain %s: %w", chain.String(), err)
		}

		if cfg.MCMSConfig != nil {
			timelockBatch = append(timelockBatch, batchOperation)
		}
	}
	// If we're not using MCMS, we can just return now as we've already confirmed the transactions
	if len(timelockBatch) == 0 {
		return cldf.ChangesetOutput{}, nil
	}
	mcmContract, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, cfg.MCMSConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	prop, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmContract,
		inspectors,
		timelockBatch,
		fmt.Sprintf("proposal to set RMNRemote on RMNProxy for chains %v", cfg.ChainSelectors),
		*cfg.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{
			*prop,
		},
	}, nil
}

func setRMNRemoteOnRMNProxyOp(
	txOpts *bind.TransactOpts, chain cldf_evm.Chain, chainState evm.CCIPChainState, mcmsEnabled bool,
) (mcmstypes.BatchOperation, error) {
	rmnProxy := chainState.RMNProxy
	rmnRemoteAddr := chainState.RMNRemote.Address()
	setRMNTx, err := rmnProxy.SetARM(txOpts, rmnRemoteAddr)

	// note: error check is handled below
	if !mcmsEnabled {
		_, err = cldf.ConfirmIfNoErrorWithABI(chain, setRMNTx, rmn_proxy_contract.RMNProxyABI, err)
		if err != nil {
			return mcmstypes.BatchOperation{}, fmt.Errorf("failed to confirm tx to set RMNRemote on RMNProxy  for chain %s: %w", chain.String(), cldf.MaybeDataErr(err))
		}
	} else if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to build call data/transaction to set RMNRemote on RMNProxy for chain %s: %w", chain.String(), err)
	}

	batchOperation, err := proposalutils.BatchOperationForChain(chain.Selector, rmnProxy.Address().Hex(),
		setRMNTx.Data(), big.NewInt(0), string(shared.RMN), []string{})
	if err != nil {
		return mcmstypes.BatchOperation{}, fmt.Errorf("failed to create batch operation for chain%s: %w", chain.String(), err)
	}

	return batchOperation, nil
}

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

func getDeployer(e cldf.Environment, chain uint64, mcmConfig *proposalutils.TimelockConfig) *bind.TransactOpts {
	if mcmConfig == nil {
		return e.BlockChains.EVMChains()[chain].DeployerKey
	}

	return cldf.SimTransactOpts()
}

type SetRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
	RMNStaticConfig   rmn_home.RMNHomeStaticConfig
	RMNDynamicConfig  rmn_home.RMNHomeDynamicConfig
	DigestToOverride  [32]byte
	MCMSConfig        *proposalutils.TimelockConfig
}

func (c SetRMNHomeCandidateConfig) Validate(state stateview.CCIPOnChainState) error {
	err := cldf.IsValidChainSelector(c.HomeChainSelector)
	if err != nil {
		return err
	}

	if len(c.RMNDynamicConfig.OffchainConfig) != 0 {
		return errors.New("RMNDynamicConfig.OffchainConfig must be empty")
	}
	if len(c.RMNStaticConfig.OffchainConfig) != 0 {
		return errors.New("RMNStaticConfig.OffchainConfig must be empty")
	}

	if len(c.RMNStaticConfig.Nodes) > 256 {
		return errors.New("RMNStaticConfig.Nodes must be less than 256")
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

	homeChain, ok := state.Chains[c.HomeChainSelector]
	if !ok {
		return fmt.Errorf("chain %d not found", c.HomeChainSelector)
	}

	rmnHome := homeChain.RMNHome
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

func isRMNStaticConfigEqual(a, b rmn_home.RMNHomeStaticConfig) bool {
	if len(a.Nodes) != len(b.Nodes) {
		return false
	}
	nodesByPeerID := make(map[p2pkey.PeerID]rmn_home.RMNHomeNode)
	for i := range a.Nodes {
		nodesByPeerID[a.Nodes[i].PeerId] = a.Nodes[i]
	}
	for i := range b.Nodes {
		node, exists := nodesByPeerID[b.Nodes[i].PeerId]
		if !exists {
			return false
		}
		if !bytes.Equal(node.OffchainPublicKey[:], b.Nodes[i].OffchainPublicKey[:]) {
			return false
		}
	}

	return bytes.Equal(a.OffchainConfig, b.OffchainConfig)
}

func isRMNDynamicConfigEqual(a, b rmn_home.RMNHomeDynamicConfig) bool {
	if len(a.SourceChains) != len(b.SourceChains) {
		return false
	}
	sourceChainBySelector := make(map[uint64]rmn_home.RMNHomeSourceChain)
	for i := range a.SourceChains {
		sourceChainBySelector[a.SourceChains[i].ChainSelector] = a.SourceChains[i]
	}
	for i := range b.SourceChains {
		sourceChain, exists := sourceChainBySelector[b.SourceChains[i].ChainSelector]
		if !exists {
			return false
		}
		if sourceChain.FObserve != b.SourceChains[i].FObserve {
			return false
		}
		if sourceChain.ObserverNodesBitmap.Cmp(b.SourceChains[i].ObserverNodesBitmap) != 0 {
			return false
		}
	}
	return bytes.Equal(a.OffchainConfig, b.OffchainConfig)
}

type PromoteRMNHomeCandidateConfig struct {
	HomeChainSelector uint64
	DigestToPromote   [32]byte
	MCMSConfig        *proposalutils.TimelockConfig
}

func (c PromoteRMNHomeCandidateConfig) Validate(state stateview.CCIPOnChainState) error {
	err := cldf.IsValidChainSelector(c.HomeChainSelector)
	if err != nil {
		return err
	}

	homeChain, ok := state.Chains[c.HomeChainSelector]

	if !ok {
		return fmt.Errorf("chain %d not found", c.HomeChainSelector)
	}

	rmnHome := homeChain.RMNHome
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

// SetRMNHomeCandidateConfigChangeset creates a changeset to set the RMNHome candidate config
// DigestToOverride is the digest of the current candidate config that the new config will override
// StaticConfig contains the list of nodes with their peerIDs (found in their rageproxy keystore) and offchain public keys (found in the RMN keystore)
// DynamicConfig contains the list of source chains with their chain selectors, f value and the bitmap of the nodes that are oberver for each source chain
// The bitmap is a 256 bit array where each bit represents a node. If the bit matching the index of the node in the static config is set it means that the node is an observer
func SetRMNHomeCandidateConfigChangeset(e cldf.Environment, config SetRMNHomeCandidateConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = config.Validate(state)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	homeChain, ok := e.BlockChains.EVMChains()[config.HomeChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found", config.HomeChainSelector)
	}

	rmnHome := state.Chains[config.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("RMNHome not found for chain %s", homeChain.String())
	}

	deployer := getDeployer(e, config.HomeChainSelector, config.MCMSConfig)
	setCandidateTx, err := rmnHome.SetCandidate(deployer, config.RMNStaticConfig, config.RMNDynamicConfig, config.DigestToOverride)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("build RMNHome set candidate calldata for chain %s: %w", homeChain.String(), err)
	}

	if config.MCMSConfig == nil {
		chain := e.BlockChains.EVMChains()[config.HomeChainSelector]
		_, err := chain.Confirm(setCandidateTx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm tx for chain %s: %w", homeChain.String(), cldf.MaybeDataErr(err))
		}

		return cldf.ChangesetOutput{}, nil
	}

	operation, err := proposalutils.BatchOperationForChain(homeChain.Selector, rmnHome.Address().Hex(),
		setCandidateTx.Data(), big.NewInt(0), string(shared.RMN), []string{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create batch operation for chain %s: %w", homeChain.String(), err)
	}

	timelocks := deployergroup.BuildTimelockAddressPerChain(e, state)
	mcmContract, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, config.MCMSConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	inspectors, err := proposalutils.McmsInspectors(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get mcms inspector for chain %s: %w", homeChain.String(), err)
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmContract,
		inspectors,
		[]mcmstypes.BatchOperation{operation},
		"proposal to set candidate config",
		*config.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal for chain %s: %w", homeChain.String(), err)
	}

	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
}

func PromoteRMNHomeCandidateConfigChangeset(e cldf.Environment, config PromoteRMNHomeCandidateConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	err = config.Validate(state)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	homeChain, ok := e.BlockChains.EVMChains()[config.HomeChainSelector]

	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found", config.HomeChainSelector)
	}

	rmnHome := state.Chains[config.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("RMNHome not found for chain %s", homeChain.String())
	}

	currentCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get RMNHome candidate digest for chain %s: %w", homeChain.String(), err)
	}

	currentActiveDigest, err := rmnHome.GetActiveDigest(nil)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get RMNHome active digest for chain %s: %w", homeChain.String(), err)
	}

	deployer := getDeployer(e, config.HomeChainSelector, config.MCMSConfig)
	promoteCandidateTx, err := rmnHome.PromoteCandidateAndRevokeActive(deployer, currentCandidateDigest, currentActiveDigest)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("get call data to promote RMNHome candidate digest for chain %s: %w", homeChain.String(), err)
	}

	if config.MCMSConfig == nil {
		chain := e.BlockChains.EVMChains()[config.HomeChainSelector]
		_, err := chain.Confirm(promoteCandidateTx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm tx for chain %s: %w", homeChain.String(), cldf.MaybeDataErr(err))
		}

		return cldf.ChangesetOutput{}, nil
	}

	operation, err := proposalutils.BatchOperationForChain(homeChain.Selector, rmnHome.Address().Hex(),
		promoteCandidateTx.Data(), big.NewInt(0), string(shared.RMN), []string{})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create batch operation for chain %s: %w", homeChain.String(), err)
	}

	timelocks := deployergroup.BuildTimelockAddressPerChain(e, state)
	mcmContract, err := deployergroup.BuildMcmAddressesPerChainByAction(e, state, config.MCMSConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	inspectors := map[uint64]mcmssdk.Inspector{}
	inspectors[config.HomeChainSelector], err = proposalutils.McmsInspectorForChain(e, config.HomeChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get mcms inspector for chain %s: %w", homeChain.String(), err)
	}

	proposal, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		mcmContract,
		inspectors,
		[]mcmstypes.BatchOperation{operation},
		"proposal to promote candidate config",
		*config.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build proposal for chain %s: %w", homeChain.String(), err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal},
	}, nil
}

func BuildRMNRemotePerChain(e cldf.Environment, state stateview.CCIPOnChainState) map[uint64]*rmn_remote.RMNRemote {
	timelocksPerChain := make(map[uint64]*rmn_remote.RMNRemote)
	for _, chain := range e.BlockChains.EVMChains() {
		timelocksPerChain[chain.Selector] = state.Chains[chain.Selector].RMNRemote
	}
	return timelocksPerChain
}

type SetRMNHomeDynamicConfigConfig struct {
	HomeChainSelector uint64
	RMNDynamicConfig  rmn_home.RMNHomeDynamicConfig
	ActiveDigest      [32]byte
	MCMS              *proposalutils.TimelockConfig
}

func (c SetRMNHomeDynamicConfigConfig) Validate(e cldf.Environment) error {
	err := cldf.IsValidChainSelector(c.HomeChainSelector)
	if err != nil {
		return err
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	rmnHome := state.Chains[c.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return fmt.Errorf("RMNHome not found for chain %s", e.BlockChains.EVMChains()[c.HomeChainSelector].String())
	}

	currentDigest, err := rmnHome.GetActiveDigest(nil)
	if err != nil {
		return fmt.Errorf("failed to get RMNHome candidate digest for chain %s: %w", e.BlockChains.EVMChains()[c.HomeChainSelector].String(), err)
	}

	if currentDigest != c.ActiveDigest {
		return fmt.Errorf("onchain active digest (%x) does not match provided digest (%x)", currentDigest[:], c.ActiveDigest[:])
	}

	if len(c.RMNDynamicConfig.OffchainConfig) != 0 {
		return errors.New("RMNDynamicConfig.OffchainConfig must be empty")
	}

	return nil
}

func SetRMNHomeDynamicConfigChangeset(e cldf.Environment, cfg SetRMNHomeDynamicConfigConfig) (cldf.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("set RMNHome dynamic config")

	chain, exists := e.BlockChains.EVMChains()[cfg.HomeChainSelector]
	if !exists {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found", cfg.HomeChainSelector)
	}

	rmnHome := state.Chains[cfg.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("RMNHome not found for chain %s", chain.String())
	}

	deployer, err := deployerGroup.GetDeployer(cfg.HomeChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	_, err = rmnHome.SetDynamicConfig(deployer, cfg.RMNDynamicConfig, cfg.ActiveDigest)

	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to set RMNHome dynamic config for chain %s: %w", chain.String(), err)
	}

	return deployerGroup.Enact()
}

type RevokeCandidateConfig struct {
	HomeChainSelector uint64
	CandidateDigest   [32]byte
	MCMS              *proposalutils.TimelockConfig
}

func (c RevokeCandidateConfig) Validate(e cldf.Environment) error {
	err := cldf.IsValidChainSelector(c.HomeChainSelector)
	if err != nil {
		return err
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	rmnHome := state.Chains[c.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return fmt.Errorf("RMNHome not found for chain %s", e.BlockChains.EVMChains()[c.HomeChainSelector].String())
	}

	currentDigest, err := rmnHome.GetCandidateDigest(nil)
	if err != nil {
		return fmt.Errorf("failed to get RMNHome candidate digest for chain %s: %w", e.BlockChains.EVMChains()[c.HomeChainSelector].String(), err)
	}

	if currentDigest != c.CandidateDigest {
		return fmt.Errorf("onchain candidate digest (%x) does not match provided digest (%x)", currentDigest[:], c.CandidateDigest[:])
	}

	return nil
}

func RevokeRMNHomeCandidateConfigChangeset(e cldf.Environment, cfg RevokeCandidateConfig) (cldf.ChangesetOutput, error) {
	err := cfg.Validate(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(e, state, cfg.MCMS).WithDeploymentContext("revoke candidate config")

	chain, exists := e.BlockChains.EVMChains()[cfg.HomeChainSelector]
	if !exists {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain %d not found", cfg.HomeChainSelector)
	}

	rmnHome := state.Chains[cfg.HomeChainSelector].RMNHome
	if rmnHome == nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("RMNHome not found for chain %s", chain.String())
	}

	deployer, err := deployerGroup.GetDeployer(cfg.HomeChainSelector)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	_, err = rmnHome.RevokeCandidate(deployer, cfg.CandidateDigest)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to revoke candidate config for chain %s: %w", chain.String(), err)
	}

	return deployerGroup.Enact()
}

func SetRMNRemoteConfigChangeset(e cldf.Environment, config ccipseq.SetRMNRemoteConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deps := opsutil.ConfigureDependencies{
		Env:          e,
		CurrentState: state,
	}
	seqReport, err := operations.ExecuteSequence(e.OperationsBundle, ccipseq.SetRMNRemoteConfigSequence, deps, config)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute SetRMNRemoteConfig sequence: %w", err)
	}

	return seqReport.Output.ToChangesetOutput(nil), nil
}
