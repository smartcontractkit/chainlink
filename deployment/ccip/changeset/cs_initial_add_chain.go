package changeset

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/imdario/mergo"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/merklemulti"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/internal"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_home"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

var _ deployment.ChangeSet[NewChainsConfig] = ConfigureNewChains

// ConfigureNewChains enables new chains as destination(s) for CCIP
// It performs the following steps per chain:
// - addChainConfig + AddDON (candidate->primary promotion i.e. init) on the home chain
// - SetOCR3Config on the remote chain
// ConfigureNewChains assumes that the home chain is already enabled and all CCIP contracts are already deployed.
func ConfigureNewChains(env deployment.Environment, c NewChainsConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid NewChainsConfig: %w", err)
	}
	err := configureChain(env, c)
	if err != nil {
		env.Logger.Errorw("Failed to configure chain", "err", err)
		return deployment.ChangesetOutput{}, deployment.MaybeDataErr(err)
	}
	return deployment.ChangesetOutput{
		Proposals:   []timelock.MCMSWithTimelockProposal{},
		AddressBook: nil,
		JobSpecs:    nil,
	}, nil
}

type CCIPOCRParams struct {
	OCRParameters types.OCRParameters
	// Note contains pointers to Arb feeds for prices
	CommitOffChainConfig pluginconfig.CommitOffchainConfig
	// Note ontains USDC config
	ExecuteOffChainConfig pluginconfig.ExecuteOffchainConfig
}

// Override overrides non-empty dst CCIPOCRParams attributes with non-empty src CCIPOCRParams attributes values
// and returns the updated CCIPOCRParams.
func (c CCIPOCRParams) Override(overrides CCIPOCRParams) (CCIPOCRParams, error) {
	err := mergo.Merge(&c, &overrides, mergo.WithOverride)
	if err != nil {
		return CCIPOCRParams{}, err
	}
	return c, nil
}

func (c CCIPOCRParams) Validate() error {
	if err := c.OCRParameters.Validate(); err != nil {
		return fmt.Errorf("invalid OCR parameters: %w", err)
	}
	if err := c.CommitOffChainConfig.Validate(); err != nil {
		return fmt.Errorf("invalid commit off-chain config: %w", err)
	}
	if err := c.ExecuteOffChainConfig.Validate(); err != nil {
		return fmt.Errorf("invalid execute off-chain config: %w", err)
	}
	return nil
}

type NewChainsConfig struct {
	// Common to all chains
	HomeChainSel uint64
	FeedChainSel uint64
	// Per chain config
	ChainConfigByChain map[uint64]CCIPOCRParams
}

func (c NewChainsConfig) Chains() []uint64 {
	chains := make([]uint64, 0, len(c.ChainConfigByChain))
	for chain := range c.ChainConfigByChain {
		chains = append(chains, chain)
	}
	return chains
}

func (c NewChainsConfig) Validate() error {
	if err := deployment.IsValidChainSelector(c.HomeChainSel); err != nil {
		return fmt.Errorf("invalid home chain selector: %d - %w", c.HomeChainSel, err)
	}
	if err := deployment.IsValidChainSelector(c.FeedChainSel); err != nil {
		return fmt.Errorf("invalid feed chain selector: %d - %w", c.FeedChainSel, err)
	}
	// Validate chain config
	for chain, cfg := range c.ChainConfigByChain {
		if err := cfg.Validate(); err != nil {
			return fmt.Errorf("invalid OCR params for chain %d: %w", chain, err)
		}
		if cfg.CommitOffChainConfig.PriceFeedChainSelector != ccipocr3.ChainSelector(c.FeedChainSel) {
			return fmt.Errorf("chain %d has invalid feed chain selector", chain)
		}
	}
	return nil
}

// DefaultOCRParams returns the default OCR parameters for a chain,
// except for a few values which must be parameterized (passed as arguments).
func DefaultOCRParams(
	feedChainSel uint64,
	tokenInfo map[ccipocr3.UnknownEncodedAddress]pluginconfig.TokenInfo,
	tokenDataObservers []pluginconfig.TokenDataObserverConfig,
) CCIPOCRParams {
	return CCIPOCRParams{
		OCRParameters: types.OCRParameters{
			DeltaProgress:                           internal.DeltaProgress,
			DeltaResend:                             internal.DeltaResend,
			DeltaInitial:                            internal.DeltaInitial,
			DeltaRound:                              internal.DeltaRound,
			DeltaGrace:                              internal.DeltaGrace,
			DeltaCertifiedCommitRequest:             internal.DeltaCertifiedCommitRequest,
			DeltaStage:                              internal.DeltaStage,
			Rmax:                                    internal.Rmax,
			MaxDurationQuery:                        internal.MaxDurationQuery,
			MaxDurationObservation:                  internal.MaxDurationObservation,
			MaxDurationShouldAcceptAttestedReport:   internal.MaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport: internal.MaxDurationShouldTransmitAcceptedReport,
		},
		ExecuteOffChainConfig: pluginconfig.ExecuteOffchainConfig{
			BatchGasLimit:             internal.BatchGasLimit,
			RelativeBoostPerWaitHour:  internal.RelativeBoostPerWaitHour,
			InflightCacheExpiry:       *config.MustNewDuration(internal.InflightCacheExpiry),
			RootSnoozeTime:            *config.MustNewDuration(internal.RootSnoozeTime),
			MessageVisibilityInterval: *config.MustNewDuration(internal.FirstBlockAge),
			BatchingStrategyID:        internal.BatchingStrategyID,
			TokenDataObservers:        tokenDataObservers,
		},
		CommitOffChainConfig: pluginconfig.CommitOffchainConfig{
			RemoteGasPriceBatchWriteFrequency:  *config.MustNewDuration(internal.RemoteGasPriceBatchWriteFrequency),
			TokenPriceBatchWriteFrequency:      *config.MustNewDuration(internal.TokenPriceBatchWriteFrequency),
			TokenInfo:                          tokenInfo,
			PriceFeedChainSelector:             ccipocr3.ChainSelector(feedChainSel),
			NewMsgScanBatchSize:                merklemulti.MaxNumberTreeLeaves,
			MaxReportTransmissionCheckAttempts: 5,
			RMNEnabled:                         os.Getenv("ENABLE_RMN") == "true", // only enabled in manual test
			RMNSignaturesTimeout:               30 * time.Minute,
			MaxMerkleTreeSize:                  merklemulti.MaxNumberTreeLeaves,
			SignObservationPrefix:              "chainlink ccip 1.6 rmn observation",
		},
	}
}

// configureChain assumes the all the Home chain contracts and CCIP contracts are deployed
// It does -
// 1. addChainConfig for each chain in CCIPHome
// 2. Registers the nodes with the capability registry
// 3. SetOCR3Config on the remote chain
func configureChain(
	e deployment.Environment,
	c NewChainsConfig,
) error {
	if e.OCRSecrets.IsEmpty() {
		return fmt.Errorf("OCR secrets are empty")
	}
	nodes, err := deployment.NodeInfo(e.NodeIDs, e.Offchain)
	if err != nil || len(nodes) == 0 {
		e.Logger.Errorw("Failed to get node info", "err", err)
		return err
	}
	existingState, err := LoadOnchainState(e)
	if err != nil {
		e.Logger.Errorw("Failed to load existing onchain state", "err")
		return err
	}
	homeChain := e.Chains[c.HomeChainSel]
	capReg := existingState.Chains[c.HomeChainSel].CapabilityRegistry
	if capReg == nil {
		e.Logger.Errorw("Failed to get capability registry", "chain", homeChain.String())
		return fmt.Errorf("capability registry not found")
	}
	ccipHome := existingState.Chains[c.HomeChainSel].CCIPHome
	if ccipHome == nil {
		e.Logger.Errorw("Failed to get ccip home", "chain", homeChain.String(), "err", err)
		return fmt.Errorf("ccip home not found")
	}
	rmnHome := existingState.Chains[c.HomeChainSel].RMNHome
	if rmnHome == nil {
		e.Logger.Errorw("Failed to get rmn home", "chain", homeChain.String(), "err", err)
		return fmt.Errorf("rmn home not found")
	}

	for chainSel, chainConfig := range c.ChainConfigByChain {
		chain, _ := e.Chains[chainSel]
		chainState, ok := existingState.Chains[chain.Selector]
		if !ok {
			return fmt.Errorf("chain state not found for chain %d", chain.Selector)
		}
		if chainState.OffRamp == nil {
			return fmt.Errorf("off ramp not found for chain %d", chain.Selector)
		}
		_, err = addChainConfig(
			e.Logger,
			e.Chains[c.HomeChainSel],
			ccipHome,
			chain.Selector,
			nodes.NonBootstraps().PeerIDs())
		if err != nil {
			return err
		}
		// For each chain, we create a DON on the home chain (2 OCR instances)
		if err := addDON(
			e.Logger,
			e.OCRSecrets,
			capReg,
			ccipHome,
			rmnHome.Address(),
			chainState.OffRamp,
			chain,
			e.Chains[c.HomeChainSel],
			nodes.NonBootstraps(),
			chainConfig,
		); err != nil {
			e.Logger.Errorw("Failed to add DON", "err", err)
			return err
		}
	}

	return nil
}

func setupConfigInfo(chainSelector uint64, readers [][32]byte, fChain uint8, cfg []byte) ccip_home.CCIPHomeChainConfigArgs {
	return ccip_home.CCIPHomeChainConfigArgs{
		ChainSelector: chainSelector,
		ChainConfig: ccip_home.CCIPHomeChainConfig{
			Readers: readers,
			FChain:  fChain,
			Config:  cfg,
		},
	}
}

func addChainConfig(
	lggr logger.Logger,
	h deployment.Chain,
	ccipConfig *ccip_home.CCIPHome,
	chainSelector uint64,
	p2pIDs [][32]byte,
) (ccip_home.CCIPHomeChainConfigArgs, error) {
	// First Add CCIPOCRParams that includes all p2pIDs as readers
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		OptimisticConfirmations: 1,
	})
	if err != nil {
		return ccip_home.CCIPHomeChainConfigArgs{}, err
	}
	chainConfig := setupConfigInfo(chainSelector, p2pIDs, uint8(len(p2pIDs)/3), encodedExtraChainConfig)
	tx, err := ccipConfig.ApplyChainConfigUpdates(h.DeployerKey, nil, []ccip_home.CCIPHomeChainConfigArgs{
		chainConfig,
	})
	if _, err := deployment.ConfirmIfNoError(h, tx, err); err != nil {
		return ccip_home.CCIPHomeChainConfigArgs{}, err
	}
	lggr.Infow("Applied chain config updates", "homeChain", h.String(), "addedChain", chainSelector, "chainConfig", chainConfig)
	return chainConfig, nil
}

// createDON creates one DON with 2 plugins (commit and exec)
// It first set a new candidate for the DON with the first plugin type and AddDON on capReg
// Then for subsequent operations it uses UpdateDON to promote the first plugin to the active deployment
// and to set candidate and promote it for the second plugin
func createDON(
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

	err = internal.SetupCommitDON(lggr, donID, commitConfig, capReg, home, nodes, ccipHome)
	if err != nil {
		return fmt.Errorf("setup commit don: %w", err)
	}

	// TODO: bug in contract causing this to not work as expected.
	err = internal.SetupExecDON(lggr, donID, execConfig, capReg, home, nodes, ccipHome)
	if err != nil {
		return fmt.Errorf("setup exec don: %w", err)
	}
	return ValidateCCIPHomeConfigSetUp(lggr, capReg, ccipHome, newChainSel)
}

func addDON(
	lggr logger.Logger,
	ocrSecrets deployment.OCRSecrets,
	capReg *capabilities_registry.CapabilitiesRegistry,
	ccipHome *ccip_home.CCIPHome,
	rmnHomeAddress common.Address,
	offRamp *offramp.OffRamp,
	dest deployment.Chain,
	home deployment.Chain,
	nodes deployment.Nodes,
	ocrParams CCIPOCRParams,
) error {
	ocrConfigs, err := internal.BuildOCR3ConfigForCCIPHome(
		ocrSecrets, offRamp, dest, nodes, rmnHomeAddress, ocrParams.OCRParameters, ocrParams.CommitOffChainConfig, ocrParams.ExecuteOffChainConfig)
	if err != nil {
		return err
	}
	err = createDON(lggr, capReg, ccipHome, ocrConfigs, home, dest.Selector, nodes)
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
		lggr.Debugw("Fetched OCR3 Configs",
			"MultiOCR3BaseOCRConfig.F", ocrConfig.ConfigInfo.F,
			"MultiOCR3BaseOCRConfig.N", ocrConfig.ConfigInfo.N,
			"MultiOCR3BaseOCRConfig.IsSignatureVerificationEnabled", ocrConfig.ConfigInfo.IsSignatureVerificationEnabled,
			"Signers", ocrConfig.Signers,
			"Transmitters", ocrConfig.Transmitters,
			"configDigest", hex.EncodeToString(ocrConfig.ConfigInfo.ConfigDigest[:]),
			"chain", dest.String(),
		)
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

// ValidateCCIPHomeConfigSetUp checks that the commit and exec active and candidate configs are set up correctly
func ValidateCCIPHomeConfigSetUp(
	lggr logger.Logger,
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
	lggr.Debugw("Fetched active commit digest", "commitActiveDigest", hex.EncodeToString(commitActiveDigest[:]))
	commitCandidateDigest, err := ccipHome.GetCandidateDigest(nil, donID, uint8(cctypes.PluginTypeCCIPCommit))
	if err != nil {
		return fmt.Errorf("get commit candidate digest: %w", err)
	}
	lggr.Debugw("Fetched candidate commit digest", "commitCandidateDigest", hex.EncodeToString(commitCandidateDigest[:]))
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
	lggr.Debugw("Fetched exec configs",
		"ActiveConfig.ConfigDigest", hex.EncodeToString(execConfigs.ActiveConfig.ConfigDigest[:]),
		"CandidateConfig.ConfigDigest", hex.EncodeToString(execConfigs.CandidateConfig.ConfigDigest[:]),
	)
	if execConfigs.ActiveConfig.ConfigDigest == [32]byte{} {
		return fmt.Errorf("active config digest is empty for exec, expected nonempty, cfg: %v", execConfigs.ActiveConfig)
	}
	if execConfigs.CandidateConfig.ConfigDigest != [32]byte{} {
		return fmt.Errorf("candidate config digest is nonempty for exec, expected empty, cfg: %v", execConfigs.CandidateConfig)
	}
	return nil
}
