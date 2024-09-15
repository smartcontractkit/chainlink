package ccipdeployment

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ocr3_config_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

const (
	NodeOperatorID         = 1
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"

	FirstBlockAge                           = 8 * time.Hour
	RemoteGasPriceBatchWriteFrequency       = 30 * time.Minute
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

func DeployCapReg(lggr logger.Logger, chains map[uint64]deployment.Chain, chainSel uint64) (deployment.AddressBook, common.Address, error) {
	ab := deployment.NewMemoryAddressBook()
	chain := chains[chainSel]
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
		return ab, common.Address{}, err
	}

	lggr.Infow("deployed capreg", "addr", capReg.Address)
	ccipConfig, err := deployContract(
		lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*ccip_config.CCIPConfig] {
			ccAddr, tx, cc, err2 := ccip_config.DeployCCIPConfig(
				chain.DeployerKey,
				chain.Client,
				capReg.Address,
			)
			return ContractDeploy[*ccip_config.CCIPConfig]{
				Address: ccAddr, Tv: deployment.NewTypeAndVersion(CCIPConfig, deployment.Version1_6_0_dev), Tx: tx, Err: err2, Contract: cc,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy ccip config", "err", err)
		return ab, common.Address{}, err
	}
	lggr.Infow("deployed ccip config", "addr", ccipConfig.Address)

	tx, err := capReg.Contract.AddCapabilities(chain.DeployerKey, []capabilities_registry.CapabilitiesRegistryCapability{
		{
			LabelledName:          CapabilityLabelledName,
			Version:               CapabilityVersion,
			CapabilityType:        2, // consensus. not used (?)
			ResponseType:          0, // report. not used (?)
			ConfigurationContract: ccipConfig.Address,
		},
	})
	if _, err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to add capabilities", "err", err)
		return ab, common.Address{}, err
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
		return ab, common.Address{}, err
	}
	return ab, capReg.Address, nil
}

func AddNodes(
	capReg *capabilities_registry.CapabilitiesRegistry,
	chain deployment.Chain,
	p2pIDs [][32]byte,
) error {
	var nodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, p2pID := range p2pIDs {
		nodeParam := capabilities_registry.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      NodeOperatorID,
			Signer:              p2pID, // Not used in tests
			P2pId:               p2pID,
			HashedCapabilityIds: [][32]byte{CCIPCapabilityID},
		}
		nodeParams = append(nodeParams, nodeParam)
	}
	tx, err := capReg.AddNodes(chain.DeployerKey, nodeParams)
	if err != nil {
		return err
	}
	_, err = chain.Confirm(tx)
	return err
}

func SetupConfigInfo(chainSelector uint64, readers [][32]byte, fChain uint8, cfg []byte) ccip_config.CCIPConfigTypesChainConfigInfo {
	return ccip_config.CCIPConfigTypesChainConfigInfo{
		ChainSelector: chainSelector,
		ChainConfig: ccip_config.CCIPConfigTypesChainConfig{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}

func AddChainConfig(
	lggr logger.Logger,
	h deployment.Chain,
	ccipConfig *ccip_config.CCIPConfig,
	chainSelector uint64,
	p2pIDs [][32]byte,
) (ccip_config.CCIPConfigTypesChainConfigInfo, error) {
	// First Add ChainConfig that includes all p2pIDs as readers
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return ccip_config.CCIPConfigTypesChainConfigInfo{}, err
	}
	chainConfig := SetupConfigInfo(chainSelector, p2pIDs, uint8(len(p2pIDs)/3), encodedExtraChainConfig)
	tx, err := ccipConfig.ApplyChainConfigUpdates(h.DeployerKey, nil, []ccip_config.CCIPConfigTypesChainConfigInfo{
		chainConfig,
	})
	if _, err := deployment.ConfirmIfNoError(h, tx, err); err != nil {
		return ccip_config.CCIPConfigTypesChainConfigInfo{}, err
	}
	lggr.Infow("Applied chain config updates", "chainConfig", chainConfig)
	return chainConfig, nil
}

func BuildAddDONArgs(
	lggr logger.Logger,
	offRamp *offramp.OffRamp,
	dest deployment.Chain,
	nodes deployment.Nodes,
) ([]byte, error) {
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

	tabi, err := ocr3_config_encoder.IOCR3ConfigEncoderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	// Add DON on capability registry contract
	var ocr3Configs []ocr3_config_encoder.CCIPConfigTypesOCR3Config
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		var encodedOffchainConfig []byte
		var err2 error
		if pluginType == cctypes.PluginTypeCCIPCommit {
			encodedOffchainConfig, err2 = pluginconfig.EncodeCommitOffchainConfig(pluginconfig.CommitOffchainConfig{
				RemoteGasPriceBatchWriteFrequency: *commonconfig.MustNewDuration(RemoteGasPriceBatchWriteFrequency),
				// TODO: implement token price writes
				// TokenPriceBatchWriteFrequency:     *commonconfig.MustNewDuration(tokenPriceBatchWriteFrequency),
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

		ocr3Configs = append(ocr3Configs, ocr3_config_encoder.CCIPConfigTypesOCR3Config{
			PluginType:            uint8(pluginType),
			ChainSelector:         dest.Selector,
			F:                     configF,
			OffchainConfigVersion: offchainConfigVersion,
			OfframpAddress:        offRamp.Address().Bytes(),
			P2pIds:                p2pIDs,
			Signers:               signersBytes,
			Transmitters:          transmittersBytes,
			OffchainConfig:        offchainConfig,
		})
	}

	// TODO: Can just use utils.ABIEncode directly here.
	encodedCall, err := tabi.Pack("exposeOCR3Config", ocr3Configs)
	if err != nil {
		return nil, err
	}

	// Trim first four bytes to remove function selector.
	encodedConfigs := encodedCall[4:]
	return encodedConfigs, nil
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
	ccipConfig *ccip_config.CCIPConfig,
) ([]offramp.MultiOCR3BaseOCRConfigArgs, error) {
	var offrampOCR3Configs []offramp.MultiOCR3BaseOCRConfigArgs
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err2 := ccipConfig.GetOCRConfig(&bind.CallOpts{
			Context: context.Background(),
		}, donID, uint8(pluginType))
		if err2 != nil {
			return nil, err2
		}
		if len(ocrConfig) != 1 {
			return nil, errors.New("expected exactly one OCR3 config")
		}
		var signerAddresses []common.Address
		for _, signer := range ocrConfig[0].Config.Signers {
			signerAddresses = append(signerAddresses, common.BytesToAddress(signer))
		}

		var transmitterAddresses []common.Address
		for _, transmitter := range ocrConfig[0].Config.Transmitters {
			transmitterAddresses = append(transmitterAddresses, common.BytesToAddress(transmitter))
		}

		offrampOCR3Configs = append(offrampOCR3Configs, offramp.MultiOCR3BaseOCRConfigArgs{
			ConfigDigest:                   ocrConfig[0].ConfigDigest,
			OcrPluginType:                  uint8(pluginType),
			F:                              ocrConfig[0].Config.F,
			IsSignatureVerificationEnabled: pluginType == cctypes.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}
	return offrampOCR3Configs, nil
}

func AddDON(
	lggr logger.Logger,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipConfig *ccip_config.CCIPConfig,
	offRamp *offramp.OffRamp,
	dest deployment.Chain,
	home deployment.Chain,
	nodes deployment.Nodes,
) error {
	encodedConfigs, err := BuildAddDONArgs(lggr, offRamp, dest, nodes)
	if err != nil {
		return err
	}
	tx, err := capReg.AddDON(home.DeployerKey, nodes.PeerIDs(), []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: CCIPCapabilityID,
			Config:       encodedConfigs,
		},
	}, false, false, nodes.DefaultF())
	if _, err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return err
	}
	don, err := LatestCCIPDON(capReg)
	if err != nil {
		return err
	}
	offrampOCR3Configs, err := BuildSetOCR3ConfigArgs(don.Id, ccipConfig)
	if err != nil {
		return err
	}

	tx, err = offRamp.SetOCR3Configs(dest.DeployerKey, offrampOCR3Configs)
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
