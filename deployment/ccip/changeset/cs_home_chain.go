package changeset

import (
	"bytes"
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var _ deployment.ChangeSet[DeployHomeChainConfig] = DeployHomeChainChangeset

// DeployHomeChainChangeset is a separate changeset because it is a standalone deployment performed once in home chain for the entire CCIP deployment.
func DeployHomeChainChangeset(env deployment.Environment, cfg DeployHomeChainConfig) (deployment.ChangesetOutput, error) {
	err := cfg.Validate()
	if err != nil {
		return deployment.ChangesetOutput{}, errors.Wrapf(deployment.ErrInvalidConfig, "%v", err)
	}
	ab := deployment.NewMemoryAddressBook()
	// Note we also deploy the cap reg.
	_, err = deployHomeChain(env.Logger, env, ab, env.Chains[cfg.HomeChainSel], cfg.RMNStaticConfig, cfg.RMNDynamicConfig, cfg.NodeOperators, cfg.NodeP2PIDsPerNodeOpAdmin)
	if err != nil {
		env.Logger.Errorw("Failed to deploy cap reg", "err", err, "addresses", env.ExistingAddresses)
		return deployment.ChangesetOutput{
			AddressBook: ab,
		}, err
	}

	return deployment.ChangesetOutput{
		AddressBook: ab,
	}, nil
}

type DeployHomeChainConfig struct {
	HomeChainSel             uint64
	RMNStaticConfig          rmn_home.RMNHomeStaticConfig
	RMNDynamicConfig         rmn_home.RMNHomeDynamicConfig
	NodeOperators            []capabilities_registry.CapabilitiesRegistryNodeOperator
	NodeP2PIDsPerNodeOpAdmin map[string][][32]byte
}

func (c DeployHomeChainConfig) Validate() error {
	if c.HomeChainSel == 0 {
		return errors.New("home chain selector must be set")
	}
	if c.RMNDynamicConfig.OffchainConfig == nil {
		return errors.New("offchain config for RMNHomeDynamicConfig must be set")
	}
	if c.RMNStaticConfig.OffchainConfig == nil {
		return errors.New("offchain config for RMNHomeStaticConfig must be set")
	}
	if len(c.NodeOperators) == 0 {
		return errors.New("node operators must be set")
	}
	for _, nop := range c.NodeOperators {
		if nop.Admin == (common.Address{}) {
			return errors.New("node operator admin address must be set")
		}
		if nop.Name == "" {
			return errors.New("node operator name must be set")
		}
		if len(c.NodeP2PIDsPerNodeOpAdmin[nop.Name]) == 0 {
			return fmt.Errorf("node operator %s must have node p2p ids provided", nop.Name)
		}
	}

	return nil
}

// deployCapReg deploys the CapabilitiesRegistry contract if it is not already deployed
// and returns a deployment.ContractDeploy struct with the address and contract instance.
func deployCapReg(
	lggr logger.Logger,
	state CCIPOnChainState,
	ab deployment.AddressBook,
	chain deployment.Chain,
) (*deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry], error) {
	homeChainState, exists := state.Chains[chain.Selector]
	if exists {
		cr := homeChainState.CapabilityRegistry
		if cr != nil {
			lggr.Infow("Found CapabilitiesRegistry in chain state", "address", cr.Address().String())
			return &deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
				Address: cr.Address(), Contract: cr, Tv: deployment.NewTypeAndVersion(CapabilitiesRegistry, deployment.Version1_0_0),
			}, nil
		}
	}
	capReg, err := deployment.DeployContract(lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry] {
			crAddr, tx, cr, err2 := capabilities_registry.DeployCapabilitiesRegistry(
				chain.DeployerKey,
				chain.Client,
			)
			return deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
				Address: crAddr, Contract: cr, Tv: deployment.NewTypeAndVersion(CapabilitiesRegistry, deployment.Version1_0_0), Tx: tx, Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy capreg", "chain", chain.String(), "err", err)
		return nil, err
	}
	return capReg, nil
}

func deployHomeChain(
	lggr logger.Logger,
	e deployment.Environment,
	ab deployment.AddressBook,
	chain deployment.Chain,
	rmnHomeStatic rmn_home.RMNHomeStaticConfig,
	rmnHomeDynamic rmn_home.RMNHomeDynamicConfig,
	nodeOps []capabilities_registry.CapabilitiesRegistryNodeOperator,
	nodeP2PIDsPerNodeOpAdmin map[string][][32]byte,
) (*deployment.ContractDeploy[*capabilities_registry.CapabilitiesRegistry], error) {
	// load existing state
	state, err := LoadOnchainState(e)
	if err != nil {
		return nil, fmt.Errorf("failed to load onchain state: %w", err)
	}
	// Deploy CapabilitiesRegistry, CCIPHome, RMNHome
	capReg, err := deployCapReg(lggr, state, ab, chain)
	if err != nil {
		return nil, err
	}

	lggr.Infow("deployed/connected to capreg", "addr", capReg.Address)
	var ccipHomeAddr common.Address
	if state.Chains[chain.Selector].CCIPHome != nil {
		lggr.Infow("CCIPHome already deployed", "addr", state.Chains[chain.Selector].CCIPHome.Address().String())
		ccipHomeAddr = state.Chains[chain.Selector].CCIPHome.Address()
	} else {
		ccipHome, err := deployment.DeployContract(
			lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*ccip_home.CCIPHome] {
				ccAddr, tx, cc, err2 := ccip_home.DeployCCIPHome(
					chain.DeployerKey,
					chain.Client,
					capReg.Address,
				)
				return deployment.ContractDeploy[*ccip_home.CCIPHome]{
					Address: ccAddr, Tv: deployment.NewTypeAndVersion(CCIPHome, deployment.Version1_6_0), Tx: tx, Err: err2, Contract: cc,
				}
			})
		if err != nil {
			lggr.Errorw("Failed to deploy CCIPHome", "chain", chain.String(), "err", err)
			return nil, err
		}
		ccipHomeAddr = ccipHome.Address
	}
	rmnHome := state.Chains[chain.Selector].RMNHome
	if state.Chains[chain.Selector].RMNHome != nil {
		lggr.Infow("RMNHome already deployed", "addr", state.Chains[chain.Selector].RMNHome.Address().String())
	} else {
		rmnHomeContract, err := deployment.DeployContract(
			lggr, chain, ab,
			func(chain deployment.Chain) deployment.ContractDeploy[*rmn_home.RMNHome] {
				rmnAddr, tx, rmn, err2 := rmn_home.DeployRMNHome(
					chain.DeployerKey,
					chain.Client,
				)
				return deployment.ContractDeploy[*rmn_home.RMNHome]{
					Address: rmnAddr, Tv: deployment.NewTypeAndVersion(RMNHome, deployment.Version1_6_0), Tx: tx, Err: err2, Contract: rmn,
				}
			},
		)
		if err != nil {
			lggr.Errorw("Failed to deploy RMNHome", "chain", chain.String(), "err", err)
			return nil, err
		}
		rmnHome = rmnHomeContract.Contract
	}

	// considering the RMNHome is recently deployed, there is no digest to overwrite
	configs, err := rmnHome.GetAllConfigs(nil)
	if err != nil {
		return nil, err
	}
	setCandidate := false
	promoteCandidate := false

	// check if the candidate is already set and equal to static and dynamic configs
	if isRMNDynamicConfigEqual(rmnHomeDynamic, configs.CandidateConfig.DynamicConfig) &&
		isRMNStaticConfigEqual(rmnHomeStatic, configs.CandidateConfig.StaticConfig) {
		lggr.Infow("RMNHome candidate is already set and equal to given static and dynamic configs,skip setting candidate")
	} else {
		setCandidate = true
	}
	// check the active config is equal to the static and dynamic configs
	if isRMNDynamicConfigEqual(rmnHomeDynamic, configs.ActiveConfig.DynamicConfig) &&
		isRMNStaticConfigEqual(rmnHomeStatic, configs.ActiveConfig.StaticConfig) {
		lggr.Infow("RMNHome active is already set and equal to given static and dynamic configs," +
			"skip setting and promoting candidate")
		setCandidate = false
		promoteCandidate = false
	} else {
		promoteCandidate = true
	}

	if setCandidate {
		tx, err := rmnHome.SetCandidate(
			chain.DeployerKey, rmnHomeStatic, rmnHomeDynamic, configs.CandidateConfig.ConfigDigest)
		if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, rmn_home.RMNHomeABI, err); err != nil {
			lggr.Errorw("Failed to set candidate on RMNHome", "err", err)
			return nil, err
		}
		lggr.Infow("Set candidate on RMNHome", "chain", chain.String())
	}
	if promoteCandidate {
		rmnCandidateDigest, err := rmnHome.GetCandidateDigest(nil)
		if err != nil {
			lggr.Errorw("Failed to get RMNHome candidate digest", "chain", chain.String(), "err", err)
			return nil, err
		}

		tx, err := rmnHome.PromoteCandidateAndRevokeActive(chain.DeployerKey, rmnCandidateDigest, [32]byte{})
		if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, rmn_home.RMNHomeABI, err); err != nil {
			lggr.Errorw("Failed to promote candidate and revoke active on RMNHome", "chain", chain.String(), "err", err)
			return nil, err
		}

		rmnActiveDigest, err := rmnHome.GetActiveDigest(nil)
		if err != nil {
			lggr.Errorw("Failed to get RMNHome active digest", "chain", chain.String(), "err", err)
			return nil, err
		}
		lggr.Infow("Got rmn home active digest", "digest", rmnActiveDigest)

		if rmnActiveDigest != rmnCandidateDigest {
			lggr.Errorw("RMNHome active digest does not match previously candidate digest",
				"active", rmnActiveDigest, "candidate", rmnCandidateDigest)
			return nil, errors.New("RMNHome active digest does not match candidate digest")
		}
		lggr.Infow("Promoted candidate and revoked active on RMNHome", "chain", chain.String())
	}
	// check if ccip capability exists in cap reg
	capabilities, err := capReg.Contract.GetCapabilities(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get capabilities: %w", err)
	}
	capabilityToAdd := capabilities_registry.CapabilitiesRegistryCapability{
		LabelledName:          internal.CapabilityLabelledName,
		Version:               internal.CapabilityVersion,
		CapabilityType:        2, // consensus. not used (?)
		ResponseType:          0, // report. not used (?)
		ConfigurationContract: ccipHomeAddr,
	}
	addCapability := true
	for _, cap := range capabilities {
		if cap.LabelledName == capabilityToAdd.LabelledName && cap.Version == capabilityToAdd.Version {
			lggr.Infow("Capability already exists, skipping adding capability",
				"labelledName", cap.LabelledName, "version", cap.Version)
			addCapability = false
			break
		}
	}
	// Add the capability to the CapabilitiesRegistry contract only if it does not exist
	if addCapability {
		tx, err := capReg.Contract.AddCapabilities(
			chain.DeployerKey, []capabilities_registry.CapabilitiesRegistryCapability{
				capabilityToAdd,
			})
		if _, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, capabilities_registry.CapabilitiesRegistryABI, err); err != nil {
			lggr.Errorw("Failed to add capabilities", "chain", chain.String(), "err", err)
			return nil, err
		}
		lggr.Infow("Added capability to CapabilitiesRegistry",
			"labelledName", capabilityToAdd.LabelledName, "version", capabilityToAdd.Version)
	}

	existingNodeOps, err := capReg.Contract.GetNodeOperators(nil)
	if err != nil {
		return nil, err
	}
	nodeOpsMap := make(map[string]capabilities_registry.CapabilitiesRegistryNodeOperator)
	for _, nop := range nodeOps {
		nodeOpsMap[nop.Admin.String()] = nop
	}
	for _, existingNop := range existingNodeOps {
		if _, ok := nodeOpsMap[existingNop.Admin.String()]; ok {
			lggr.Infow("Node operator already exists", "admin", existingNop.Admin.String())
			delete(nodeOpsMap, existingNop.Admin.String())
		}
	}
	nodeOpsToAdd := make([]capabilities_registry.CapabilitiesRegistryNodeOperator, 0, len(nodeOpsMap))
	for _, nop := range nodeOpsMap {
		nodeOpsToAdd = append(nodeOpsToAdd, nop)
	}
	// Need to fetch nodeoperators ids to be able to add nodes for corresponding node operators
	p2pIDsByNodeOpID := make(map[uint32][][32]byte)
	if len(nodeOpsToAdd) > 0 {
		tx, err := capReg.Contract.AddNodeOperators(chain.DeployerKey, nodeOps)
		txBlockNum, err := deployment.ConfirmIfNoErrorWithABI(chain, tx, capabilities_registry.CapabilitiesRegistryABI, err)
		if err != nil {
			lggr.Errorw("Failed to add node operators", "chain", chain.String(), "err", err)
			return nil, err
		}
		addedEvent, err := capReg.Contract.FilterNodeOperatorAdded(&bind.FilterOpts{
			Start:   txBlockNum,
			Context: context.Background(),
		}, nil, nil)
		if err != nil {
			lggr.Errorw("Failed to filter NodeOperatorAdded event", "chain", chain.String(), "err", err)
			return capReg, err
		}

		for addedEvent.Next() {
			for nopName, p2pID := range nodeP2PIDsPerNodeOpAdmin {
				if addedEvent.Event.Name == nopName {
					lggr.Infow("Added node operator", "admin", addedEvent.Event.Admin, "name", addedEvent.Event.Name)
					p2pIDsByNodeOpID[addedEvent.Event.NodeOperatorId] = p2pID
				}
			}
		}
	} else {
		lggr.Infow("No new node operators to add")
		foundNopID := make(map[uint32]bool)
		for nopName, p2pID := range nodeP2PIDsPerNodeOpAdmin {
			// this is to find the node operator id for the given node operator name
			// node operator start from id 1, starting from 1 to len(existingNodeOps)
			totalNops := len(existingNodeOps)
			if totalNops >= math.MaxUint32 {
				return nil, errors.New("too many node operators")
			}
			for nopID := uint32(1); nopID <= uint32(totalNops); nopID++ {
				// if we already found the node operator id, skip
				if foundNopID[nopID] {
					continue
				}
				nodeOp, err := capReg.Contract.GetNodeOperator(nil, nopID)
				if err != nil {
					return capReg, fmt.Errorf("failed to get node operator %d: %w", nopID, err)
				}
				if nodeOp.Name == nopName {
					p2pIDsByNodeOpID[nopID] = p2pID
					foundNopID[nopID] = true
					break
				}
			}
		}
	}
	if len(p2pIDsByNodeOpID) != len(nodeP2PIDsPerNodeOpAdmin) {
		lggr.Errorw("Failed to add all node operators", "added", maps.Keys(p2pIDsByNodeOpID), "expected", maps.Keys(nodeP2PIDsPerNodeOpAdmin), "chain", chain.String())
		return capReg, errors.New("failed to add all node operators")
	}
	// Adds initial set of nodes to CR, who all have the CCIP capability
	if err := addNodes(lggr, capReg.Contract, chain, p2pIDsByNodeOpID); err != nil {
		return capReg, err
	}
	return capReg, nil
}

func isEqualCapabilitiesRegistryNodeParams(a, b capabilities_registry.CapabilitiesRegistryNodeParams) bool {
	if len(a.HashedCapabilityIds) != len(b.HashedCapabilityIds) {
		return false
	}
	for i := range a.HashedCapabilityIds {
		if !bytes.Equal(a.HashedCapabilityIds[i][:], b.HashedCapabilityIds[i][:]) {
			return false
		}
	}
	return a.NodeOperatorId == b.NodeOperatorId &&
		bytes.Equal(a.Signer[:], b.Signer[:]) &&
		bytes.Equal(a.P2pId[:], b.P2pId[:]) &&
		bytes.Equal(a.EncryptionPublicKey[:], b.EncryptionPublicKey[:])
}

func addNodes(
	lggr logger.Logger,
	capReg *capabilities_registry.CapabilitiesRegistry,
	chain deployment.Chain,
	p2pIDsByNodeOpId map[uint32][][32]byte,
) error {
	var nodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	nodes, err := capReg.GetNodes(nil)
	if err != nil {
		return err
	}
	existingNodeParams := make(map[p2ptypes.PeerID]capabilities_registry.CapabilitiesRegistryNodeParams)
	for _, node := range nodes {
		existingNodeParams[node.P2pId] = capabilities_registry.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      node.NodeOperatorId,
			Signer:              node.Signer,
			P2pId:               node.P2pId,
			EncryptionPublicKey: node.EncryptionPublicKey,
			HashedCapabilityIds: node.HashedCapabilityIds,
		}
	}
	for nopID, p2pIDs := range p2pIDsByNodeOpId {
		for _, p2pID := range p2pIDs {
			// if any p2pIDs are empty throw error
			if bytes.Equal(p2pID[:], make([]byte, 32)) {
				return errors.Wrapf(errors.New("empty p2pID"), "p2pID: %x selector: %d", p2pID, chain.Selector)
			}
			nodeParam := capabilities_registry.CapabilitiesRegistryNodeParams{
				NodeOperatorId:      nopID,
				Signer:              p2pID, // Not used in tests
				P2pId:               p2pID,
				EncryptionPublicKey: p2pID, // Not used in tests
				HashedCapabilityIds: [][32]byte{internal.CCIPCapabilityID},
			}
			if existing, ok := existingNodeParams[p2pID]; ok {
				if isEqualCapabilitiesRegistryNodeParams(existing, nodeParam) {
					lggr.Infow("Node already exists", "p2pID", p2pID)
					continue
				}
			}

			nodeParams = append(nodeParams, nodeParam)
		}
	}
	if len(nodeParams) == 0 {
		lggr.Infow("No new nodes to add")
		return nil
	}
	tx, err := capReg.AddNodes(chain.DeployerKey, nodeParams)
	if err != nil {
		lggr.Errorw("Failed to add nodes", "chain", chain.String(),
			"err", deployment.DecodedErrFromABIIfDataErr(err, capabilities_registry.CapabilitiesRegistryABI))
		return err
	}
	_, err = chain.Confirm(tx)
	return err
}

type RemoveDONsConfig struct {
	HomeChainSel uint64
	DonIDs       []uint32
	MCMS         *MCMSConfig
}

func (c RemoveDONsConfig) Validate(homeChain CCIPChainState) error {
	if err := deployment.IsValidChainSelector(c.HomeChainSel); err != nil {
		return fmt.Errorf("home chain selector must be set %w", err)
	}
	if len(c.DonIDs) == 0 {
		return errors.New("don ids must be set")
	}
	// Cap reg must exist
	if homeChain.CapabilityRegistry == nil {
		return errors.New("cap reg does not exist")
	}
	if homeChain.CCIPHome == nil {
		return errors.New("ccip home does not exist")
	}
	if err := internal.DONIdExists(homeChain.CapabilityRegistry, c.DonIDs); err != nil {
		return err
	}
	return nil
}

// RemoveDONs removes DONs from the CapabilitiesRegistry contract.
// TODO: Could likely be moved to common, but needs
// a common state struct first.
func RemoveDONs(e deployment.Environment, cfg RemoveDONsConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(e)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	homeChain, ok := e.Chains[cfg.HomeChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("home chain %d not found", cfg.HomeChainSel)
	}
	homeChainState := state.Chains[cfg.HomeChainSel]
	if err := cfg.Validate(homeChainState); err != nil {
		return deployment.ChangesetOutput{}, err
	}
	txOpts := homeChain.DeployerKey
	if cfg.MCMS != nil {
		txOpts = deployment.SimTransactOpts()
	}

	tx, err := homeChainState.CapabilityRegistry.RemoveDONs(txOpts, cfg.DonIDs)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	if cfg.MCMS == nil {
		_, err = homeChain.Confirm(tx)
		if err != nil {
			return deployment.ChangesetOutput{}, err
		}
		e.Logger.Infof("Removed dons using deployer key tx %s", tx.Hash().String())
		return deployment.ChangesetOutput{}, nil
	}
	p, err := proposalutils.BuildProposalFromBatches(
		map[uint64]common.Address{
			cfg.HomeChainSel: homeChainState.Timelock.Address(),
		},
		map[uint64]*gethwrappers.ManyChainMultiSig{
			cfg.HomeChainSel: homeChainState.ProposerMcm,
		},
		[]timelock.BatchChainOperation{
			{
				ChainIdentifier: mcms.ChainIdentifier(cfg.HomeChainSel),
				Batch: []mcms.Operation{
					{
						To:    homeChainState.CapabilityRegistry.Address(),
						Data:  tx.Data(),
						Value: big.NewInt(0),
					},
				},
			},
		},
		"Remove DONs",
		cfg.MCMS.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	e.Logger.Infof("Created proposal to remove dons")
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}
