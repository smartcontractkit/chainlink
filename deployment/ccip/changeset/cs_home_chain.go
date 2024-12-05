package changeset

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

var _ deployment.ChangeSet[DeployHomeChainConfig] = DeployHomeChain

// DeployHomeChain is a separate changeset because it is a standalone deployment performed once in home chain for the entire CCIP deployment.
func DeployHomeChain(env deployment.Environment, cfg DeployHomeChainConfig) (deployment.ChangesetOutput, error) {
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
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: ab,
		JobSpecs:    nil,
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
		return fmt.Errorf("home chain selector must be set")
	}
	if c.RMNDynamicConfig.OffchainConfig == nil {
		return fmt.Errorf("offchain config for RMNHomeDynamicConfig must be set")
	}
	if c.RMNStaticConfig.OffchainConfig == nil {
		return fmt.Errorf("offchain config for RMNHomeStaticConfig must be set")
	}
	if len(c.NodeOperators) == 0 {
		return fmt.Errorf("node operators must be set")
	}
	for _, nop := range c.NodeOperators {
		if nop.Admin == (common.Address{}) {
			return fmt.Errorf("node operator admin address must be set")
		}
		if nop.Name == "" {
			return fmt.Errorf("node operator name must be set")
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
		lggr.Errorw("Failed to deploy capreg", "err", err)
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
	ccipHome, err := deployment.DeployContract(
		lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*ccip_home.CCIPHome] {
			ccAddr, tx, cc, err2 := ccip_home.DeployCCIPHome(
				chain.DeployerKey,
				chain.Client,
				capReg.Address,
			)
			return deployment.ContractDeploy[*ccip_home.CCIPHome]{
				Address: ccAddr, Tv: deployment.NewTypeAndVersion(CCIPHome, deployment.Version1_6_0_dev), Tx: tx, Err: err2, Contract: cc,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy CCIPHome", "err", err)
		return nil, err
	}
	lggr.Infow("deployed CCIPHome", "addr", ccipHome.Address)

	rmnHome, err := deployment.DeployContract(
		lggr, chain, ab,
		func(chain deployment.Chain) deployment.ContractDeploy[*rmn_home.RMNHome] {
			rmnAddr, tx, rmn, err2 := rmn_home.DeployRMNHome(
				chain.DeployerKey,
				chain.Client,
			)
			return deployment.ContractDeploy[*rmn_home.RMNHome]{
				Address: rmnAddr, Tv: deployment.NewTypeAndVersion(RMNHome, deployment.Version1_6_0_dev), Tx: tx, Err: err2, Contract: rmn,
			}
		},
	)
	if err != nil {
		lggr.Errorw("Failed to deploy RMNHome", "err", err)
		return nil, err
	}
	lggr.Infow("deployed RMNHome", "addr", rmnHome.Address)

	// considering the RMNHome is recently deployed, there is no digest to overwrite
	tx, err := rmnHome.Contract.SetCandidate(chain.DeployerKey, rmnHomeStatic, rmnHomeDynamic, [32]byte{})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to set candidate on RMNHome", "err", err)
		return nil, err
	}

	rmnCandidateDigest, err := rmnHome.Contract.GetCandidateDigest(nil)
	if err != nil {
		lggr.Errorw("Failed to get RMNHome candidate digest", "err", err)
		return nil, err
	}

	tx, err = rmnHome.Contract.PromoteCandidateAndRevokeActive(chain.DeployerKey, rmnCandidateDigest, [32]byte{})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to promote candidate and revoke active on RMNHome", "err", err)
		return nil, err
	}

	rmnActiveDigest, err := rmnHome.Contract.GetActiveDigest(nil)
	if err != nil {
		lggr.Errorw("Failed to get RMNHome active digest", "err", err)
		return nil, err
	}
	lggr.Infow("Got rmn home active digest", "digest", rmnActiveDigest)

	if rmnActiveDigest != rmnCandidateDigest {
		lggr.Errorw("RMNHome active digest does not match previously candidate digest",
			"active", rmnActiveDigest, "candidate", rmnCandidateDigest)
		return nil, errors.New("RMNHome active digest does not match candidate digest")
	}

	tx, err = capReg.Contract.AddCapabilities(chain.DeployerKey, []capabilities_registry.CapabilitiesRegistryCapability{
		{
			LabelledName:          internal.CapabilityLabelledName,
			Version:               internal.CapabilityVersion,
			CapabilityType:        2, // consensus. not used (?)
			ResponseType:          0, // report. not used (?)
			ConfigurationContract: ccipHome.Address,
		},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to add capabilities", "err", err)
		return nil, err
	}

	tx, err = capReg.Contract.AddNodeOperators(chain.DeployerKey, nodeOps)
	txBlockNum, err := deployment.ConfirmIfNoError(chain, tx, err)
	if err != nil {
		lggr.Errorw("Failed to add node operators", "err", err)
		return nil, err
	}
	addedEvent, err := capReg.Contract.FilterNodeOperatorAdded(&bind.FilterOpts{
		Start:   txBlockNum,
		Context: context.Background(),
	}, nil, nil)
	if err != nil {
		lggr.Errorw("Failed to filter NodeOperatorAdded event", "err", err)
		return capReg, err
	}
	// Need to fetch nodeoperators ids to be able to add nodes for corresponding node operators
	p2pIDsByNodeOpId := make(map[uint32][][32]byte)
	for addedEvent.Next() {
		for nopName, p2pId := range nodeP2PIDsPerNodeOpAdmin {
			if addedEvent.Event.Name == nopName {
				lggr.Infow("Added node operator", "admin", addedEvent.Event.Admin, "name", addedEvent.Event.Name)
				p2pIDsByNodeOpId[addedEvent.Event.NodeOperatorId] = p2pId
			}
		}
	}
	if len(p2pIDsByNodeOpId) != len(nodeP2PIDsPerNodeOpAdmin) {
		lggr.Errorw("Failed to add all node operators", "added", maps.Keys(p2pIDsByNodeOpId), "expected", maps.Keys(nodeP2PIDsPerNodeOpAdmin))
		return capReg, errors.New("failed to add all node operators")
	}
	// Adds initial set of nodes to CR, who all have the CCIP capability
	if err := addNodes(lggr, capReg.Contract, chain, p2pIDsByNodeOpId); err != nil {
		return capReg, err
	}
	return capReg, nil
}

func isEqualCapabilitiesRegistryNodeParams(a, b capabilities_registry.CapabilitiesRegistryNodeParams) (bool, error) {
	aBytes, err := json.Marshal(a)
	if err != nil {
		return false, err
	}
	bBytes, err := json.Marshal(b)
	if err != nil {
		return false, err
	}
	return bytes.Equal(aBytes, bBytes), nil
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
				if isEqual, err := isEqualCapabilitiesRegistryNodeParams(existing, nodeParam); err != nil && isEqual {
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
		lggr.Errorw("Failed to add nodes", "err", deployment.MaybeDataErr(err))
		return err
	}
	_, err = chain.Confirm(tx)
	return err
}
