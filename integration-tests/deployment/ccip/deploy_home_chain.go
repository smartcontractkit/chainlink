package ccipdeployment

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"

	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/rmn_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

const (
	NodeOperatorID         = 1
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"

	FirstBlockAge                           = 8 * time.Hour
	RemoteGasPriceBatchWriteFrequency       = 30 * time.Minute
	TokenPriceBatchWriteFrequency           = 30 * time.Minute
	BatchGasLimit                           = 6_500_000
	RelativeBoostPerWaitHour                = 1.5
	InflightCacheExpiry                     = 10 * time.Minute
	RootSnoozeTime                          = 30 * time.Minute
	BatchingStrategyID                      = 0
	DeltaProgress                           = 30 * time.Second
	DeltaResend                             = 10 * time.Second
	DeltaInitial                            = 20 * time.Second
	DeltaRound                              = 2 * time.Second
	DeltaGrace                              = 2 * time.Second
	DeltaCertifiedCommitRequest             = 10 * time.Second
	DeltaStage                              = 10 * time.Second
	Rmax                                    = 3
	MaxDurationQuery                        = 50 * time.Millisecond
	MaxDurationObservation                  = 5 * time.Second
	MaxDurationShouldAcceptAttestedReport   = 10 * time.Second
	MaxDurationShouldTransmitAcceptedReport = 10 * time.Second
)

var (
	CCIPCapabilityID = utils.Keccak256Fixed(MustABIEncode(`[{"type": "string"}, {"type": "string"}]`, CapabilityLabelledName, CapabilityVersion))
)

func MustABIEncode(abiString string, args ...interface{}) []byte {
	encoded, err := utils.ABIEncode(abiString, args...)
	if err != nil {
		panic(err)
	}
	return encoded
}

func DeployCapReg(lggr logger.Logger, ab deployment.AddressBook, chain deployment.Chain) (*ContractDeploy[*capabilities_registry.CapabilitiesRegistry], error) {
	capReg, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*capabilities_registry.CapabilitiesRegistry] {
			crAddr, tx, cr, err2 := capabilities_registry.DeployCapabilitiesRegistry(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
				Address: crAddr, Contract: cr, Tv: deployment.NewTypeAndVersion(CapabilitiesRegistry, deployment.Version1_0_0), Tx: tx, Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy capreg", "err", err)
		return nil, err
	}

	lggr.Infow("deployed capreg", "addr", capReg.Address)
	ccipHome, err := deployContract(
		lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*ccip_home.CCIPHome] {
			ccAddr, tx, cc, err2 := ccip_home.DeployCCIPHome(
				chain.DeployerKey,
				chain.Client,
				capReg.Address,
			)
			return ContractDeploy[*ccip_home.CCIPHome]{
				Address: ccAddr, Tv: deployment.NewTypeAndVersion(CCIPHome, deployment.Version1_6_0_dev), Tx: tx, Err: err2, Contract: cc,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy CCIPHome", "err", err)
		return nil, err
	}
	lggr.Infow("deployed CCIPHome", "addr", ccipHome.Address)

	rmnHome, err := deployContract(
		lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*rmn_home.RMNHome] {
			rmnAddr, tx, rmn, err2 := rmn_home.DeployRMNHome(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*rmn_home.RMNHome]{
				Address: rmnAddr, Tv: deployment.NewTypeAndVersion(RMNHome, deployment.Version1_6_0_dev), Tx: tx, Err: err2, Contract: rmn,
			}
		},
	)
	if err != nil {
		lggr.Errorw("Failed to deploy RMNHome", "err", err)
		return nil, err
	}
	lggr.Infow("deployed RMNHome", "addr", rmnHome.Address)

	// TODO: properly configure RMNHome
	tx, err := rmnHome.Contract.SetCandidate(chain.DeployerKey, rmn_home.RMNHomeStaticConfig{
		Nodes:          []rmn_home.RMNHomeNode{},
		OffchainConfig: []byte("static config"),
	}, rmn_home.RMNHomeDynamicConfig{
		SourceChains:   []rmn_home.RMNHomeSourceChain{},
		OffchainConfig: []byte("dynamic config"),
	}, [32]byte{})
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

	if rmnActiveDigest != rmnCandidateDigest {
		lggr.Errorw("RMNHome active digest does not match previously candidate digest",
			"active", rmnActiveDigest, "candidate", rmnCandidateDigest)
		return nil, errors.New("RMNHome active digest does not match candidate digest")
	}

	tx, err = capReg.Contract.AddCapabilities(chain.DeployerKey, []capabilities_registry.CapabilitiesRegistryCapability{
		{
			LabelledName:          CapabilityLabelledName,
			Version:               CapabilityVersion,
			CapabilityType:        2, // consensus. not used (?)
			ResponseType:          0, // report. not used (?)
			ConfigurationContract: ccipHome.Address,
		},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to add capabilities", "err", err)
		return nil, err
	}
	// TODO: Just one for testing.
	tx, err = capReg.Contract.AddNodeOperators(chain.DeployerKey, []capabilities_registry.CapabilitiesRegistryNodeOperator{
		{
			Admin: chain.DeployerKey.From,
			Name:  "NodeOperator",
		},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to add node operators", "err", err)
		return nil, err
	}
	return capReg, nil
}

func AddNodes(
	lggr logger.Logger,
	capReg *capabilities_registry.CapabilitiesRegistry,
	chain deployment.Chain,
	p2pIDs [][32]byte,
) error {
	var nodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, p2pID := range p2pIDs {
		// if any p2pIDs are empty throw error
		if bytes.Equal(p2pID[:], make([]byte, 32)) {
			return errors.Wrapf(errors.New("empty p2pID"), "p2pID: %x selector: %d", p2pID, chain.Selector)
		}
		nodeParam := capabilities_registry.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      NodeOperatorID,
			Signer:              p2pID, // Not used in tests
			P2pId:               p2pID,
			EncryptionPublicKey: p2pID, // Not used in tests
			HashedCapabilityIds: [][32]byte{CCIPCapabilityID},
		}
		nodeParams = append(nodeParams, nodeParam)
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

func BuildAddDONArgs(
	lggr logger.Logger,
	offRamp *offramp.OffRamp,
	dest deployment.Chain,
	feedChainSel uint64,
	// Token address on Dest chain to aggregate address on feed chain
	tokenInfo map[ocrtypes.Account]pluginconfig.TokenInfo,
	nodes deployment.Nodes,
) (map[cctypes.PluginType]ccip_home.CCIPHomeOCR3Config, error) {
	p2pIDs := nodes.PeerIDs()
	// Get OCR3 Config from helper
	var schedule []int
	var oracles []confighelper2.OracleIdentityExtra
	for _, node := range nodes {
		schedule = append(schedule, 1)
		cfg := node.SelToOCRConfig[dest.Selector]
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  cfg.OnchainPublicKey,
				TransmitAccount:   cfg.TransmitAccount,
				OffchainPublicKey: cfg.OffchainPublicKey,
				PeerID:            cfg.PeerID.String()[4:],
			}, ConfigEncryptionPublicKey: cfg.ConfigEncryptionPublicKey,
		})
	}

	// Add DON on capability registry contract
	ocr3Configs := make(map[cctypes.PluginType]ccip_home.CCIPHomeOCR3Config)
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		var encodedOffchainConfig []byte
		var err2 error
		if pluginType == cctypes.PluginTypeCCIPCommit {
			encodedOffchainConfig, err2 = pluginconfig.EncodeCommitOffchainConfig(pluginconfig.CommitOffchainConfig{
				RemoteGasPriceBatchWriteFrequency:  *commonconfig.MustNewDuration(RemoteGasPriceBatchWriteFrequency),
				TokenPriceBatchWriteFrequency:      *commonconfig.MustNewDuration(TokenPriceBatchWriteFrequency),
				PriceFeedChainSelector:             ccipocr3.ChainSelector(feedChainSel),
				TokenInfo:                          tokenInfo,
				NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
				MaxReportTransmissionCheckAttempts: 5,
				MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
				SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
			})
		} else {
			encodedOffchainConfig, err2 = pluginconfig.EncodeExecuteOffchainConfig(pluginconfig.ExecuteOffchainConfig{
				BatchGasLimit:             BatchGasLimit,
				RelativeBoostPerWaitHour:  RelativeBoostPerWaitHour,
				MessageVisibilityInterval: *commonconfig.MustNewDuration(FirstBlockAge),
				InflightCacheExpiry:       *commonconfig.MustNewDuration(InflightCacheExpiry),
				RootSnoozeTime:            *commonconfig.MustNewDuration(RootSnoozeTime),
				BatchingStrategyID:        BatchingStrategyID,
			})
		}
		if err2 != nil {
			return nil, err2
		}
		signers, transmitters, configF, _, offchainConfigVersion, offchainConfig, err2 := ocr3confighelper.ContractSetConfigArgsForTests(
			DeltaProgress,
			DeltaResend,
			DeltaInitial,
			DeltaRound,
			DeltaGrace,
			DeltaCertifiedCommitRequest,
			DeltaStage,
			Rmax,
			schedule,
			oracles,
			encodedOffchainConfig,
			nil, // maxDurationInitialization
			MaxDurationQuery,
			MaxDurationObservation,
			MaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport,
			int(nodes.DefaultF()),
			[]byte{}, // empty OnChainConfig
		)
		if err2 != nil {
			return nil, err2
		}

		signersBytes := make([][]byte, len(signers))
		for i, signer := range signers {
			signersBytes[i] = signer
		}

		transmittersBytes := make([][]byte, len(transmitters))
		for i, transmitter := range transmitters {
			parsed, err2 := common.ParseHexOrString(string(transmitter))
			if err2 != nil {
				return nil, err2
			}
			transmittersBytes[i] = parsed
		}

		var ocrNodes []ccip_home.CCIPHomeOCR3Node
		for i := range nodes {
			ocrNodes = append(ocrNodes, ccip_home.CCIPHomeOCR3Node{
				P2pId:          p2pIDs[i],
				SignerKey:      signersBytes[i],
				TransmitterKey: transmittersBytes[i],
			})
		}

		_, ok := ocr3Configs[pluginType]
		if ok {
			return nil, fmt.Errorf("pluginType %s already exists in ocr3Configs", pluginType.String())
		}

		ocr3Configs[pluginType] = ccip_home.CCIPHomeOCR3Config{
			PluginType:            uint8(pluginType),
			ChainSelector:         dest.Selector,
			FRoleDON:              configF,
			OffchainConfigVersion: offchainConfigVersion,
			OfframpAddress:        offRamp.Address().Bytes(),
			Nodes:                 ocrNodes,
			OffchainConfig:        offchainConfig,
			// TODO: Deploy RMNHome and set address here
			RmnHomeAddress: common.BytesToAddress(randomBytes(20)).Bytes(),
		}
	}

	return ocr3Configs, nil
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func LatestCCIPDON(registry *capabilities_registry.CapabilitiesRegistry) (*capabilities_registry.CapabilitiesRegistryDONInfo, error) {
	dons, err := registry.GetDONs(nil)
	if err != nil {
		return nil, err
	}
	var ccipDON capabilities_registry.CapabilitiesRegistryDONInfo
	for _, don := range dons {
		if len(don.CapabilityConfigurations) == 1 &&
			don.CapabilityConfigurations[0].CapabilityId == CCIPCapabilityID &&
			don.Id > ccipDON.Id {
			ccipDON = don
		}
	}
	return &ccipDON, nil
}

func BuildSetOCR3ConfigArgs(
	donID uint32,
	ccipHome *ccip_home.CCIPHome,
	destSelector uint64,
) ([]offramp.MultiOCR3BaseOCRConfigArgs, error) {
	var offrampOCR3Configs []offramp.MultiOCR3BaseOCRConfigArgs
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err2 := ccipHome.GetAllConfigs(&bind.CallOpts{
			Context: context.Background(),
		}, donID, uint8(pluginType))
		if err2 != nil {
			return nil, err2
		}

		fmt.Printf("pluginType: %s, destSelector: %d, donID: %d, activeConfig digest: %x, candidateConfig digest: %x\n",
			pluginType.String(), destSelector, donID, ocrConfig.ActiveConfig.ConfigDigest, ocrConfig.CandidateConfig.ConfigDigest)

		// we expect only an active config and no candidate config.
		if ocrConfig.ActiveConfig.ConfigDigest == [32]byte{} || ocrConfig.CandidateConfig.ConfigDigest != [32]byte{} {
			return nil, fmt.Errorf("invalid OCR3 config state, expected active config and no candidate config, donID: %d", donID)
		}

		activeConfig := ocrConfig.ActiveConfig
		var signerAddresses []common.Address
		var transmitterAddresses []common.Address
		for _, node := range activeConfig.Config.Nodes {
			signerAddresses = append(signerAddresses, common.BytesToAddress(node.SignerKey))
			transmitterAddresses = append(transmitterAddresses, common.BytesToAddress(node.TransmitterKey))
		}

		offrampOCR3Configs = append(offrampOCR3Configs, offramp.MultiOCR3BaseOCRConfigArgs{
			ConfigDigest:                   activeConfig.ConfigDigest,
			OcrPluginType:                  uint8(pluginType),
			F:                              activeConfig.Config.FRoleDON,
			IsSignatureVerificationEnabled: pluginType == cctypes.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}
	return offrampOCR3Configs, nil
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
	nodes deployment.Nodes,
) ([]mcms.Operation, error) {
	commitConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPCommit]
	if !ok {
		return nil, fmt.Errorf("missing commit plugin in ocr3Configs")
	}

	execConfig, ok := ocr3Configs[cctypes.PluginTypeCCIPExec]
	if !ok {
		return nil, fmt.Errorf("missing exec plugin in ocr3Configs")
	}

	tabi, err := ccip_home.CCIPHomeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	latestDon, err := LatestCCIPDON(capReg)
	if err != nil {
		return nil, err
	}

	donID := latestDon.Id + 1
	mcmsOps := []mcms.Operation{}

	err = setupCommitDON(tabi, donID, commitConfig, capReg, home, nodes, ccipHome)
	if err != nil {
		return nil, fmt.Errorf("setup commit don: %w", err)
	}

	// TODO: bug in contract causing this to not work as expected.
	err = setupExecDON(tabi, donID, execConfig, capReg, home, nodes, ccipHome)
	if err != nil {
		return nil, fmt.Errorf("setup exec don: %w", err)
	}

	// final sanity checks on configs.
	commitConfigs, err := ccipHome.GetAllConfigs(&bind.CallOpts{
		Pending: true,
	}, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return nil, fmt.Errorf("get all commit configs: %w", err)
	}
	commitActiveDigest, err := ccipHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return nil, fmt.Errorf("get active commit digest: %w", err)
	}
	commitCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return nil, fmt.Errorf("get commit candidate digest: %w", err)
	}
	if commitConfigs.ActiveConfig.ConfigDigest == [32]byte{} {
		return nil, fmt.Errorf(
			"active config digest is empty for commit, expected nonempty, donID: %d, cfg: %+v, config digest from GetActiveDigest call: %x, config digest from GetCandidateDigest call: %x",
			donID, commitConfigs.ActiveConfig, commitActiveDigest, commitCandidateDigest)
	}
	if commitConfigs.CandidateConfig.ConfigDigest != [32]byte{} {
		return nil, fmt.Errorf(
			"candidate config digest is nonempty for commit, expected empty, donID: %d, cfg: %+v, config digest from GetCandidateDigest call: %x, config digest from GetActiveDigest call: %x",
			donID, commitConfigs.CandidateConfig, commitCandidateDigest, commitActiveDigest)
	}

	execConfigs, err := ccipHome.GetAllConfigs(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return nil, fmt.Errorf("get all exec configs: %w", err)
	}
	if execConfigs.ActiveConfig.ConfigDigest == [32]byte{} {
		return nil, fmt.Errorf("active config digest is empty for exec, expected nonempty, cfg: %v", execConfigs.ActiveConfig)
	}
	if execConfigs.CandidateConfig.ConfigDigest != [32]byte{} {
		return nil, fmt.Errorf("candidate config digest is nonempty for exec, expected empty, cfg: %v", execConfigs.CandidateConfig)
	}

	return mcmsOps, nil
}

func setupExecDON(
	tabi *abi.ABI,
	donID uint32,
	execConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	home deployment.Chain,
	nodes deployment.Nodes,
	ccipHome *ccip_home.CCIPHome,
) error {
	encodedSetCandidateCall, err := tabi.Pack(
		"setCandidate",
		donID,
		execConfig.PluginType,
		execConfig,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack set candidate call: %w", err)
	}

	// set candidate call
	tx, err := capReg.UpdateDON(
		home.DeployerKey,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedSetCandidateCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return fmt.Errorf("update don w/ exec config: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return fmt.Errorf("confirm update don w/ exec config: %w", err)
	}

	execCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, execConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get commit candidate digest: %w", err)
	}

	if execCandidateDigest == [32]byte{} {
		return fmt.Errorf("candidate digest is empty, expected nonempty")
	}

	// promote candidate call
	encodedPromotionCall, err := tabi.Pack(
		"promoteCandidateAndRevokeActive",
		donID,
		execConfig.PluginType,
		execCandidateDigest,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack promotion call: %w", err)
	}

	tx, err = capReg.UpdateDON(
		home.DeployerKey,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedPromotionCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return fmt.Errorf("update don w/ exec config: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return fmt.Errorf("confirm update don w/ exec config: %w", err)
	}

	// check that candidate digest is empty.
	execCandidateDigest, err = ccipHome.GetCandidateDigest(nil, donID, execConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get exec candidate digest 2nd time: %w", err)
	}

	if execCandidateDigest != [32]byte{} {
		return fmt.Errorf("candidate digest is nonempty after promotion, expected empty")
	}

	// check that active digest is non-empty.
	execActiveDigest, err := ccipHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return fmt.Errorf("get active exec digest: %w", err)
	}

	if execActiveDigest == [32]byte{} {
		return fmt.Errorf("active exec digest is empty, expected nonempty")
	}

	execConfigs, err := ccipHome.GetAllConfigs(nil, donID, uint8(cctypes.PluginTypeCCIPExec))
	if err != nil {
		return fmt.Errorf("get all exec configs 2nd time: %w", err)
	}

	// print the above info
	fmt.Printf("completed exec DON creation and promotion: donID: %d execCandidateDigest: %x, execActiveDigest: %x, execCandidateDigestFromGetAllConfigs: %x, execActiveDigestFromGetAllConfigs: %x\n",
		donID, execCandidateDigest, execActiveDigest, execConfigs.CandidateConfig.ConfigDigest, execConfigs.ActiveConfig.ConfigDigest)

	return nil
}

func setupCommitDON(
	tabi *abi.ABI,
	donID uint32,
	commitConfig ccip_home.CCIPHomeOCR3Config,
	capReg *capabilities_registry.CapabilitiesRegistry,
	home deployment.Chain,
	nodes deployment.Nodes,
	ccipHome *ccip_home.CCIPHome,
) error {
	encodedSetCandidateCall, err := tabi.Pack(
		"setCandidate",
		donID,
		commitConfig.PluginType,
		commitConfig,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack set candidate call: %w", err)
	}

	tx, err := capReg.AddDON(home.DeployerKey, nodes.PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: CCIPCapabilityID,
			Config:       encodedSetCandidateCall,
		},
	}, false, false, nodes.DefaultF())
	if err != nil {
		return fmt.Errorf("add don w/ commit config: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return fmt.Errorf("confirm add don w/ commit config: %w", err)
	}

	commitCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, commitConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get commit candidate digest: %w", err)
	}

	if commitCandidateDigest == [32]byte{} {
		return fmt.Errorf("candidate digest is empty, expected nonempty")
	}
	fmt.Printf("commit candidate digest after setCandidate: %x\n", commitCandidateDigest)

	encodedPromotionCall, err := tabi.Pack(
		"promoteCandidateAndRevokeActive",
		donID,
		commitConfig.PluginType,
		commitCandidateDigest,
		[32]byte{},
	)
	if err != nil {
		return fmt.Errorf("pack promotion call: %w", err)
	}

	tx, err = capReg.UpdateDON(
		home.DeployerKey,
		donID,
		nodes.PeerIDs(),
		[]capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
			{
				CapabilityId: CCIPCapabilityID,
				Config:       encodedPromotionCall,
			},
		},
		false,
		nodes.DefaultF(),
	)
	if err != nil {
		return fmt.Errorf("update don w/ commit config: %w", err)
	}

	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return fmt.Errorf("confirm update don w/ commit config: %w", err)
	}

	// check that candidate digest is empty.
	commitCandidateDigest, err = ccipHome.GetCandidateDigest(nil, donID, commitConfig.PluginType)
	if err != nil {
		return fmt.Errorf("get commit candidate digest 2nd time: %w", err)
	}

	if commitCandidateDigest != [32]byte{} {
		return fmt.Errorf("candidate digest is nonempty after promotion, expected empty")
	}

	// check that active digest is non-empty.
	commitActiveDigest, err := ccipHome.GetActiveDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get active commit digest: %w", err)
	}

	if commitActiveDigest == [32]byte{} {
		return fmt.Errorf("active commit digest is empty, expected nonempty")
	}

	commitConfigs, err := ccipHome.GetAllConfigs(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get all commit configs 2nd time: %w", err)
	}

	// print the above information
	fmt.Printf("completed commit DON creation and promotion: donID: %d, commitCandidateDigest: %x, commitActiveDigest: %x, commitCandidateDigestFromGetAllConfigs: %x, commitActiveDigestFromGetAllConfigs: %x\n",
		donID, commitCandidateDigest, commitActiveDigest, commitConfigs.CandidateConfig.ConfigDigest, commitConfigs.ActiveConfig.ConfigDigest)

	return nil
}

func AddDON(
	lggr logger.Logger,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	offRamp *offramp.OffRamp,
	feedChainSel uint64,
	// Token address on Dest chain to aggregate address on feed chain
	tokenInfo map[ocrtypes.Account]pluginconfig.TokenInfo,
	dest deployment.Chain,
	home deployment.Chain,
	nodes deployment.Nodes,
) error {
	ocrConfigs, err := BuildAddDONArgs(lggr, offRamp, dest, feedChainSel, tokenInfo, nodes)
	if err != nil {
		return err
	}
	_, err = CreateDON(lggr, capReg, ccipHome, ocrConfigs, home, nodes)
	if err != nil {
		return err
	}
	don, err := LatestCCIPDON(capReg)
	if err != nil {
		return err
	}
	lggr.Infow("Added DON", "donID", don.Id)

	offrampOCR3Configs, err := BuildSetOCR3ConfigArgs(don.Id, ccipHome, dest.Selector)
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

	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err := offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: context.Background(),
		}, uint8(pluginType))
		if err != nil {
			return err
		}
		// TODO: assertions to be done as part of full state
		// resprentation validation CCIP-3047
		if offrampOCR3Configs[pluginType].ConfigDigest != ocrConfig.ConfigInfo.ConfigDigest {
			return fmt.Errorf("%s OCR3 config digest mismatch", pluginType.String())
		}
		if offrampOCR3Configs[pluginType].F != ocrConfig.ConfigInfo.F {
			return fmt.Errorf("%s OCR3 config F mismatch", pluginType.String())
		}
		if offrampOCR3Configs[pluginType].IsSignatureVerificationEnabled != ocrConfig.ConfigInfo.IsSignatureVerificationEnabled {
			return fmt.Errorf("%s OCR3 config signature verification mismatch", pluginType.String())
		}
		if pluginType == cctypes.PluginTypeCCIPCommit {
			// only commit will set signers, exec doesn't need them.
			for i, signer := range offrampOCR3Configs[pluginType].Signers {
				if !bytes.Equal(signer.Bytes(), ocrConfig.Signers[i].Bytes()) {
					return fmt.Errorf("%s OCR3 config signer mismatch", pluginType.String())
				}
			}
		}
		for i, transmitter := range offrampOCR3Configs[pluginType].Transmitters {
			if !bytes.Equal(transmitter.Bytes(), ocrConfig.Transmitters[i].Bytes()) {
				return fmt.Errorf("%s OCR3 config transmitter mismatch", pluginType.String())
			}
		}
	}

	return nil
}
