package changeset

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"

	"github.com/smartcontractkit/chainlink/deployment"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	p2ptypes "github.com/smartcontractkit/chainlink/v2/core/services/p2p/types"
)

// DeployCapReg deploys the CapabilitiesRegistry contract if it is not already deployed
// and returns a deployment.ContractDeploy struct with the address and contract instance.
func DeployCapReg(
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
	capReg, err := DeployCapReg(lggr, state, ab, chain)
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
	if err := AddNodes(lggr, capReg.Contract, chain, p2pIDsByNodeOpId); err != nil {
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

func AddNodes(
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

func SetupConfigInfo(chainSelector uint64, readers [][32]byte, fChain uint8, cfg []byte) ccip_home.CCIPHomeChainConfigArgs {
	return ccip_home.CCIPHomeChainConfigArgs{
		ChainSelector: chainSelector,
		ChainConfig: ccip_home.CCIPHomeChainConfig{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}

func AddChainConfig(
	lggr logger.Logger,
	h deployment.Chain,
	ccipConfig *ccip_home.CCIPHome,
	chainSelector uint64,
	p2pIDs [][32]byte,
) (ccip_home.CCIPHomeChainConfigArgs, error) {
	// First Add ChainConfig that includes all p2pIDs as readers
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return ccip_home.CCIPHomeChainConfigArgs{}, err
	}
	chainConfig := SetupConfigInfo(chainSelector, p2pIDs, uint8(len(p2pIDs)/3), encodedExtraChainConfig)
	tx, err := ccipConfig.ApplyChainConfigUpdates(h.DeployerKey, nil, []ccip_home.CCIPHomeChainConfigArgs{
		chainConfig,
	})
	if _, err := deployment.ConfirmIfNoError(h, tx, err); err != nil {
		return ccip_home.CCIPHomeChainConfigArgs{}, err
	}
	lggr.Infow("Applied chain config updates", "chainConfig", chainConfig)
	return chainConfig, nil
}

// CreateDON creates one DON with 2 plugins (commit and exec)
// It first set a new candidate for the DON with the first plugin type and AddDON on capReg
// Then for subsequent operations it uses UpdateDON to promote the first plugin to the active deployment
// and to set candidate and promote it for the second plugin
func CreateDON(
	lggr logger.Logger,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	ocr3Configs map[cctypes.PluginType]ccip_home.CCIPHomeOCR3Config,
	home deployment.Chain,
	newChainSel uint64,
	nodes deployment.Nodes,
) error {
	commitConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPCommit]
	if !ok {
		return fmt.Errorf("missing commit plugin in ocr3Configs")
	}

	execConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPExec]
	if !ok {
		return fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	latestDon, err := internal.LatestCCIPDON(capReg)
	if err != nil {
		return err
	}

	donID := latestDon.Id + 1

	err = internal.SetupCommitDON(donID, commitConfig, capReg, home, nodes, ccipHome)
	if err != nil {
		return fmt.Errorf("setup commit don: %w", err)
	}

	// TODO: bug in contract causing this to not work as expected.
	err = internal.SetupExecDON(donID, execConfig, capReg, home, nodes, ccipHome)
	if err != nil {
		return fmt.Errorf("setup exec don: %w", err)
	}
	return ValidateCCIPHomeConfigSetUp(capReg, ccipHome, newChainSel)
}

// SetCandidateCommitPluginWithAddDonOps sets the candidate commit config by calling setCandidate on CCIPHome contract through the AddDON call on CapReg contract
// This should be done first before calling any other UpdateDON calls
// This proposes to set up OCR3 config for the commit plugin for the DON
func NewDonWithCandidateOp(
	donID uint32,
	pluginConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	nodes deployment.Nodes,
) (mcms.Operation, error) {
	encodedSetCandidateCall, err := internal.CCIPHomeABI.Pack(
		"setCandidate",
		donID,
		pluginConfig.PluginType,
		pluginConfig,
		[32]byte{},
	)
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("pack set candidate call: %w", err)
	}
	addDonTx, err := capReg.AddDON(deployment.SimTransactOpts(), nodes.PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: internal.CCIPCapabilityID,
			Config:       encodedSetCandidateCall,
		},
	}, false, false, nodes.DefaultF())
	if err != nil {
		return mcms.Operation{}, fmt.Errorf("could not generate add don tx w/ commit config: %w", err)
	}
	return mcms.Operation{
		To:    capReg.Address(),
		Data:  addDonTx.Data(),
		Value: big.NewInt(0),
	}, nil
}

// ValidateCCIPHomeConfigSetUp checks that the commit and exec active and candidate configs are set up correctly
func ValidateCCIPHomeConfigSetUp(
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	chainSel uint64,
) error {
	// fetch DONID
	donID, err := internal.DonIDForChain(capReg, ccipHome, chainSel)
	if err != nil {
		return fmt.Errorf("fetch don id for chain: %w", err)
	}
	// final sanity checks on configs.
	commitConfigs, err := ccipHome.GetAllConfigs(&bind.CallOpts{
		//Pending: true,
	}, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get all commit configs: %w", err)
	}
	commitActiveDigest, err := ccipHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get active commit digest: %w", err)
	}
	commitCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get commit candidate digest: %w", err)
	}
	if commitConfigs.ActiveConfig.ConfigDigest == [32]byte{} {
		return fmt.Errorf(
			"active config digest is empty for commit, expected nonempty, donID: %d, cfg: %+v, config digest from GetActiveDigest call: %x, config digest from GetCandidateDigest call: %x",
			donID, commitConfigs.ActiveConfig, commitActiveDigest, commitCandidateDigest)
	}
	if commitConfigs.CandidateConfig.ConfigDigest != [32]byte{} {
		return fmt.Errorf(
			"candidate config digest is nonempty for commit, expected empty, donID: %d, cfg: %+v, config digest from GetCandidateDigest call: %x, config digest from GetActiveDigest call: %x",
			donID, commitConfigs.CandidateConfig, commitCandidateDigest, commitActiveDigest)
	}

	execConfigs, err := ccipHome.GetAllConfigs(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return fmt.Errorf("get all exec configs: %w", err)
	}
	if execConfigs.ActiveConfig.ConfigDigest == [32]byte{} {
		return fmt.Errorf("active config digest is empty for exec, expected nonempty, cfg: %v", execConfigs.ActiveConfig)
	}
	if execConfigs.CandidateConfig.ConfigDigest != [32]byte{} {
		return fmt.Errorf("candidate config digest is nonempty for exec, expected empty, cfg: %v", execConfigs.CandidateConfig)
	}
	return nil
}

func AddDON(
	lggr logger.Logger,
	ocrSecrets deployment.OCRSecrets,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	rmnHomeAddress common.Address,
	offRamp *offramp.OffRamp,
	feedChainSel uint64,
	// Token address on Dest chain to aggregate address on feed chain
	tokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
	dest deployment.Chain,
	home deployment.Chain,
	nodes deployment.Nodes,
	tokenConfigs []pluginconfig.TokenDataObserverConfig,
) error {
	ocrConfigs, err := internal.BuildOCR3ConfigForCCIPHome(
		ocrSecrets, offRamp, dest, feedChainSel, tokenInfo, nodes, rmnHomeAddress, tokenConfigs)
	if err != nil {
		return err
	}
	err = CreateDON(lggr, capReg, ccipHome, ocrConfigs, home, dest.Selector, nodes)
	if err != nil {
		return err
	}
	don, err := internal.LatestCCIPDON(capReg)
	if err != nil {
		return err
	}
	lggr.Infow("Added DON", "donID", don.Id)

	offrampOCR3Configs, err := internal.BuildSetOCR3ConfigArgs(don.Id, ccipHome, dest.Selector)
	if err != nil {
		return err
	}
	lggr.Infow("Setting OCR3 Configs",
		"offrampOCR3Configs", offrampOCR3Configs,
		"configDigestCommit", hex.EncodeToString(offrampOCR3Configs[cctypes.PluginTypeCCIPCommit].ConfigDigest[:]),
		"configDigestExec", hex.EncodeToString(offrampOCR3Configs[cctypes.PluginTypeCCIPExec].ConfigDigest[:]),
		"chainSelector", dest.Selector,
	)

	tx, err := offRamp.SetOCR3Configs(dest.DeployerKey, offrampOCR3Configs)
	if _, err := deployment.ConfirmIfNoError(dest, tx, err); err != nil {
		return err
	}

	mapOfframpOCR3Configs := make(map[cctypes.PluginType]offramp.MultiOCR3BaseOCRConfigArgs)
	for _, config := range offrampOCR3Configs {
		mapOfframpOCR3Configs[cctypes.PluginType(config.OcrPluginType)] = config
	}

	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err := offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: context.Background(),
		}, uint8(pluginType))
		if err != nil {
			return err
		}
		// TODO: assertions to be done as part of full state
		// resprentation validation CCIP-3047
		if mapOfframpOCR3Configs[pluginType].ConfigDigest != ocrConfig.ConfigInfo.ConfigDigest {
			return fmt.Errorf("%s OCR3 config digest mismatch", pluginType.String())
		}
		if mapOfframpOCR3Configs[pluginType].F != ocrConfig.ConfigInfo.F {
			return fmt.Errorf("%s OCR3 config F mismatch", pluginType.String())
		}
		if mapOfframpOCR3Configs[pluginType].IsSignatureVerificationEnabled != ocrConfig.ConfigInfo.IsSignatureVerificationEnabled {
			return fmt.Errorf("%s OCR3 config signature verification mismatch", pluginType.String())
		}
		if pluginType == cctypes.PluginTypeCCIPCommit {
			// only commit will set signers, exec doesn't need them.
			for i, signer := range mapOfframpOCR3Configs[pluginType].Signers {
				if !bytes.Equal(signer.Bytes(), ocrConfig.Signers[i].Bytes()) {
					return fmt.Errorf("%s OCR3 config signer mismatch", pluginType.String())
				}
			}
		}
		for i, transmitter := range mapOfframpOCR3Configs[pluginType].Transmitters {
			if !bytes.Equal(transmitter.Bytes(), ocrConfig.Transmitters[i].Bytes()) {
				return fmt.Errorf("%s OCR3 config transmitter mismatch", pluginType.String())
			}
		}
	}

	return nil
}

func ApplyChainConfigUpdatesOp(
	e deployment.Environment,
	state CCIPOnChainState,
	homeChainSel uint64,
	chains []uint64,
) (mcms.Operation, error) {
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil {
		return mcms.Operation{}, err
	}
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return mcms.Operation{}, err
	}
	var chainConfigUpdates []ccip_home.CCIPHomeChainConfigArgs
	for _, chainSel := range chains {
		chainConfig := SetupConfigInfo(chainSel, nodes.NonBootstraps().PeerIDs(),
			nodes.DefaultF(), encodedExtraChainConfig)
		chainConfigUpdates = append(chainConfigUpdates, chainConfig)
	}

	addChain, err := state.Chains[homeChainSel].CCIPHome.ApplyChainConfigUpdates(
		deployment.SimTransactOpts(),
		nil,
		chainConfigUpdates,
	)
	if err != nil {
		return mcms.Operation{}, err
	}
	return mcms.Operation{
		To:    state.Chains[homeChainSel].CCIPHome.Address(),
		Data:  addChain.Data(),
		Value: big.NewInt(0),
	}, nil
}
