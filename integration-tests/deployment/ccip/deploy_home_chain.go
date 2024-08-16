package ccipdeployment

import (
	"bytes"
	"context"
	"errors"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	deployment2 "github.com/smartcontractkit/ccip/integration-tests/deployment"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ocr3_config_encoder"
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

func DeployCapReg(lggr logger.Logger, chains map[uint64]deployment.Chain, chainSel uint64) (deployment.AddressBook, error) {
	ab := deployment.NewMemoryAddressBook()
	chain := chains[chainSel]
	capReg, err := deployContract(lggr, chain, ab,
		func(chain deployment.Chain) ContractDeploy[*capabilities_registry.CapabilitiesRegistry] {
			crAddr, tx, cr, err2 := capabilities_registry.DeployCapabilitiesRegistry(
				chain.DeployerKey,
				chain.Client,
			)
			return ContractDeploy[*capabilities_registry.CapabilitiesRegistry]{
				Address: crAddr, Contract: cr, TvStr: CapabilitiesRegistry_1_0_0, Tx: tx, Err: err2,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy capreg", "err", err)
		return ab, err
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
				Address: ccAddr, TvStr: CCIPConfig_1_6_0, Tx: tx, Err: err2, Contract: cc,
			}
		})
	if err != nil {
		lggr.Errorw("Failed to deploy ccip config", "err", err)
		return ab, err
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
	if err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to add capabilities", "err", err)
		return ab, err
	}
	// TODO: Just one for testing.
	tx, err = capReg.Contract.AddNodeOperators(chain.DeployerKey, []capabilities_registry.CapabilitiesRegistryNodeOperator{
		{
			Admin: chain.DeployerKey.From,
			Name:  "NodeOperator",
		},
	})
	if err := deployment.ConfirmIfNoError(chain, tx, err); err != nil {
		lggr.Errorw("Failed to add node operators", "err", err)
		return ab, err
	}
	return ab, nil
}

func sortP2PIDS(p2pIDs [][32]byte) {
	sort.Slice(p2pIDs, func(i, j int) bool {
		return bytes.Compare(p2pIDs[i][:], p2pIDs[j][:]) < 0
	})
}

func AddNodes(
	capReg *capabilities_registry.CapabilitiesRegistry,
	chain deployment.Chain,
	p2pIDs [][32]byte,
	capabilityIDs [][32]byte,
) error {
	// Need to sort, otherwise _checkIsValidUniqueSubset onChain will fail
	sortP2PIDS(p2pIDs)
	var nodeParams []capabilities_registry.CapabilitiesRegistryNodeParams
	for _, p2pID := range p2pIDs {
		nodeParam := capabilities_registry.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      NodeOperatorID,
			Signer:              p2pID, // Not used in tests
			P2pId:               p2pID,
			HashedCapabilityIds: capabilityIDs,
		}
		nodeParams = append(nodeParams, nodeParam)
	}
	tx, err := capReg.AddNodes(chain.DeployerKey, nodeParams)
	if err != nil {
		return err
	}
	return chain.Confirm(tx.Hash())
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
	f uint8,
) (ccip_config.CCIPConfigTypesChainConfigInfo, error) {
	// Need to sort, otherwise _checkIsValidUniqueSubset onChain will fail
	sortP2PIDS(p2pIDs)
	// First Add ChainConfig that includes all p2pIDs as readers
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		FinalityDepth:           10,
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return ccip_config.CCIPConfigTypesChainConfigInfo{}, err
	}
	chainConfig := SetupConfigInfo(chainSelector, p2pIDs, f, encodedExtraChainConfig)
	inputConfig := []ccip_config.CCIPConfigTypesChainConfigInfo{
		chainConfig,
	}
	tx, err := ccipConfig.ApplyChainConfigUpdates(h.DeployerKey, nil, inputConfig)
	if err != nil {
		return ccip_config.CCIPConfigTypesChainConfigInfo{}, err
	}
	if err := h.Confirm(tx.Hash()); err != nil {
		return ccip_config.CCIPConfigTypesChainConfigInfo{}, err
	}
	return chainConfig, nil
}

func AddDON(
	lggr logger.Logger,
	ccipCapabilityID [32]byte,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipConfig *ccip_config.CCIPConfig,
	offRamp *evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp,
	dest deployment.Chain,
	home deployment.Chain,
	f uint8,
	bootstrapP2PID [32]byte,
	p2pIDs [][32]byte,
	nodes []deployment2.Node,
) error {
	sortP2PIDS(p2pIDs)
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
		return err
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
			return err2
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
			int(f),
			[]byte{}, // empty OnChainConfig
		)
		if err2 != nil {
			return err2
		}

		signersBytes := make([][]byte, len(signers))
		for i, signer := range signers {
			signersBytes[i] = signer
		}

		transmittersBytes := make([][]byte, len(transmitters))
		for i, transmitter := range transmitters {
			parsed, err2 := common.ParseHexOrString(string(transmitter))
			if err != nil {
				return err2
			}
			transmittersBytes[i] = parsed
		}

		ocr3Configs = append(ocr3Configs, ocr3_config_encoder.CCIPConfigTypesOCR3Config{
			PluginType:            uint8(pluginType),
			ChainSelector:         dest.Selector,
			F:                     configF,
			OffchainConfigVersion: offchainConfigVersion,
			OfframpAddress:        offRamp.Address().Bytes(),
			BootstrapP2PIds:       [][32]byte{bootstrapP2PID},
			P2pIds:                p2pIDs,
			Signers:               signersBytes,
			Transmitters:          transmittersBytes,
			OffchainConfig:        offchainConfig,
		})
	}

	encodedCall, err := tabi.Pack("exposeOCR3Config", ocr3Configs)
	if err != nil {
		return err
	}

	// Trim first four bytes to remove function selector.
	encodedConfigs := encodedCall[4:]

	// commit so that we have an empty block to filter events from
	// TODO: required?
	//h.backend.Commit()

	tx, err := capReg.AddDON(home.DeployerKey, p2pIDs, []capabilities_registry.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: ccipCapabilityID,
			Config:       encodedConfigs,
		},
	}, false, false, f)
	if err := deployment.ConfirmIfNoError(home, tx, err); err != nil {
		return err
	}

	latestBlock, err := home.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return err
	}
	endBlock := latestBlock.Number.Uint64()
	iter, err := capReg.FilterConfigSet(&bind.FilterOpts{
		Start: endBlock - 1,
		End:   &endBlock,
	})
	if err != nil {
		return err
	}
	var donID uint32
	for iter.Next() {
		donID = iter.Event.DonId
		break
	}
	if donID == 0 {
		return errors.New("failed to get donID")
	}

	var signerAddresses []common.Address
	for _, oracle := range oracles {
		signerAddresses = append(signerAddresses, common.BytesToAddress(oracle.OnchainPublicKey))
	}

	var transmitterAddresses []common.Address
	for _, oracle := range oracles {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(oracle.TransmitAccount)))
	}

	// get the config digest from the ccip config contract and set config on the offramp.
	var offrampOCR3Configs []evm_2_evm_multi_offramp.MultiOCR3BaseOCRConfigArgs
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err2 := ccipConfig.GetOCRConfig(&bind.CallOpts{
			Context: context.Background(),
		}, donID, uint8(pluginType))
		if err2 != nil {
			return err2
		}
		if len(ocrConfig) != 1 {
			return errors.New("expected exactly one OCR3 config")
		}

		offrampOCR3Configs = append(offrampOCR3Configs, evm_2_evm_multi_offramp.MultiOCR3BaseOCRConfigArgs{
			ConfigDigest:                   ocrConfig[0].ConfigDigest,
			OcrPluginType:                  uint8(pluginType),
			F:                              f,
			IsSignatureVerificationEnabled: pluginType == cctypes.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}

	//uni.backend.Commit()

	tx, err = offRamp.SetOCR3Configs(dest.DeployerKey, offrampOCR3Configs)
	if err := deployment.ConfirmIfNoError(dest, tx, err); err != nil {
		return err
	}

	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		_, err = offRamp.LatestConfigDetails(&bind.CallOpts{
			Context: context.Background(),
		}, uint8(pluginType))
		if err != nil {
			//return err
			return deployment.MaybeDataErr(err)
		}
		// TODO: assertions
		//require.Equalf(t, offrampOCR3Configs[pluginType].ConfigDigest, ocrConfig.ConfigInfo.ConfigDigest, "%s OCR3 config digest mismatch", pluginType.String())
		//require.Equalf(t, offrampOCR3Configs[pluginType].F, ocrConfig.ConfigInfo.F, "%s OCR3 config F mismatch", pluginType.String())
		//require.Equalf(t, offrampOCR3Configs[pluginType].IsSignatureVerificationEnabled, ocrConfig.ConfigInfo.IsSignatureVerificationEnabled, "%s OCR3 config signature verification mismatch", pluginType.String())
		//if pluginType == cctypes.PluginTypeCCIPCommit {
		//	// only commit will set signers, exec doesn't need them.
		//	require.Equalf(t, offrampOCR3Configs[pluginType].Signers, ocrConfig.Signers, "%s OCR3 config signers mismatch", pluginType.String())
		//}
		//require.Equalf(t, offrampOCR3Configs[pluginType].TransmittersByEVMChainID, ocrConfig.TransmittersByEVMChainID, "%s OCR3 config transmitters mismatch", pluginType.String())
	}

	lggr.Infof("set ocr3 config on the offramp, signers: %+v, transmitters: %+v", signerAddresses, transmitterAddresses)
	return nil
}
