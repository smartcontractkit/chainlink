package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_offramp"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

var _ cldf.ChangeSet[OCR2Config] = SetOCR2ConfigForTestChangeset

type FinalOCR2Config struct {
	Signers               []common.Address
	Transmitters          []common.Address
	F                     uint8
	OnchainConfig         []byte
	OffchainConfigVersion uint64
	OffchainConfig        []byte
}

type CommitOCR2ConfigParams struct {
	DestinationChainSelector uint64
	SourceChainSelector      uint64
	OCR2ConfigParams         confighelper.PublicConfig
	GasPriceHeartBeat        config.Duration
	DAGasPriceDeviationPPB   uint32
	ExecGasPriceDeviationPPB uint32
	TokenPriceHeartBeat      config.Duration
	TokenPriceDeviationPPB   uint32
	InflightCacheExpiry      config.Duration
	PriceReportingDisabled   bool
}

func (c *CommitOCR2ConfigParams) PopulateOffChainAndOnChainCfg(priceReg common.Address) error {
	var err error
	c.OCR2ConfigParams.ReportingPluginConfig, err = testhelpers.NewCommitOffchainConfig(
		c.GasPriceHeartBeat,
		c.DAGasPriceDeviationPPB,
		c.ExecGasPriceDeviationPPB,
		c.TokenPriceHeartBeat,
		c.TokenPriceDeviationPPB,
		c.InflightCacheExpiry,
		c.PriceReportingDisabled,
	).Encode()
	if err != nil {
		return errors.Wrapf(err, "failed to encode offchain config for source chain %d and destination chain %d",
			c.SourceChainSelector, c.DestinationChainSelector)
	}
	c.OCR2ConfigParams.OnchainConfig, err = abihelpers.EncodeAbiStruct(testhelpers.NewCommitOnchainConfig(priceReg))
	if err != nil {
		return fmt.Errorf("failed to encode onchain config for source chain %d and destination chain %d: %w",
			c.SourceChainSelector, c.DestinationChainSelector, err)
	}
	return nil
}

func (c *CommitOCR2ConfigParams) Validate(state stateview.CCIPOnChainState) error {
	if err := cldf.IsValidChainSelector(c.DestinationChainSelector); err != nil {
		return fmt.Errorf("invalid DestinationChainSelector: %w", err)
	}
	if err := cldf.IsValidChainSelector(c.SourceChainSelector); err != nil {
		return fmt.Errorf("invalid SourceChainSelector: %w", err)
	}

	chain, exists := state.EVMChainState(c.DestinationChainSelector)
	if !exists {
		return fmt.Errorf("chain %d does not exist in state", c.DestinationChainSelector)
	}
	if chain.CommitStore == nil {
		return fmt.Errorf("chain %d does not have a commit store", c.DestinationChainSelector)
	}
	_, exists = chain.CommitStore[c.SourceChainSelector]
	if !exists {
		return fmt.Errorf("chain %d does not have a commit store for source chain %d", c.DestinationChainSelector, c.SourceChainSelector)
	}
	if chain.PriceRegistry == nil {
		return fmt.Errorf("chain %d does not have a price registry", c.DestinationChainSelector)
	}
	return nil
}

type ExecuteOCR2ConfigParams struct {
	DestinationChainSelector    uint64
	SourceChainSelector         uint64
	DestOptimisticConfirmations uint32
	BatchGasLimit               uint32
	RelativeBoostPerWaitHour    float64
	InflightCacheExpiry         config.Duration
	RootSnoozeTime              config.Duration
	BatchingStrategyID          uint32
	MessageVisibilityInterval   config.Duration
	ExecOnchainConfig           evm_2_evm_offramp.EVM2EVMOffRampDynamicConfig
	OCR2ConfigParams            confighelper.PublicConfig
}

func (e *ExecuteOCR2ConfigParams) PopulateOffChainAndOnChainCfg(router, priceReg common.Address) error {
	var err error
	e.OCR2ConfigParams.ReportingPluginConfig, err = testhelpers.NewExecOffchainConfig(
		e.DestOptimisticConfirmations,
		e.BatchGasLimit,
		e.RelativeBoostPerWaitHour,
		e.InflightCacheExpiry,
		e.RootSnoozeTime,
		e.BatchingStrategyID,
	).Encode()
	if err != nil {
		return fmt.Errorf("failed to encode offchain config for exec plugin, source chain %d dest chain %d :%w",
			e.SourceChainSelector, e.DestinationChainSelector, err)
	}
	e.OCR2ConfigParams.OnchainConfig, err = abihelpers.EncodeAbiStruct(testhelpers.NewExecOnchainConfig(
		e.ExecOnchainConfig.PermissionLessExecutionThresholdSeconds,
		router,
		priceReg,
		e.ExecOnchainConfig.MaxNumberOfTokensPerMsg,
		e.ExecOnchainConfig.MaxDataBytes,
	))
	if err != nil {
		return fmt.Errorf("failed to encode onchain config for exec plugin, source chain %d dest chain %d :%w",
			e.SourceChainSelector, e.DestinationChainSelector, err)
	}
	return nil
}

func (e *ExecuteOCR2ConfigParams) Validate(state stateview.CCIPOnChainState) error {
	if err := cldf.IsValidChainSelector(e.SourceChainSelector); err != nil {
		return fmt.Errorf("invalid SourceChainSelector: %w", err)
	}
	if err := cldf.IsValidChainSelector(e.DestinationChainSelector); err != nil {
		return fmt.Errorf("invalid DestinationChainSelector: %w", err)
	}
	chain, exists := state.EVMChainState(e.DestinationChainSelector)
	if !exists {
		return fmt.Errorf("chain %d does not exist in state", e.DestinationChainSelector)
	}
	if chain.EVM2EVMOffRamp == nil {
		return fmt.Errorf("chain %d does not have an EVM2EVMOffRamp", e.DestinationChainSelector)
	}
	_, exists = chain.EVM2EVMOffRamp[e.SourceChainSelector]
	if !exists {
		return fmt.Errorf("chain %d does not have an EVM2EVMOffRamp for source chain %d", e.DestinationChainSelector, e.SourceChainSelector)
	}
	if chain.PriceRegistry == nil {
		return fmt.Errorf("chain %d does not have a price registry", e.DestinationChainSelector)
	}
	if chain.Router == nil {
		return fmt.Errorf("chain %d does not have a router", e.DestinationChainSelector)
	}
	return nil
}

type OCR2Config struct {
	CommitConfigs []CommitOCR2ConfigParams
	ExecConfigs   []ExecuteOCR2ConfigParams
}

func (o OCR2Config) Validate(state stateview.CCIPOnChainState) error {
	for _, c := range o.CommitConfigs {
		if err := c.Validate(state); err != nil {
			return err
		}
	}
	for _, e := range o.ExecConfigs {
		if err := e.Validate(state); err != nil {
			return err
		}
	}
	return nil
}

// SetOCR2ConfigForTestChangeset sets the OCR2 config on the chain for commit and offramp
// This is currently not suitable for prod environments it's only for testing
func SetOCR2ConfigForTestChangeset(env cldf.Environment, c OCR2Config) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load CCIP onchain state: %w", err)
	}
	if err := c.Validate(state); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid OCR2 config: %w", err)
	}
	for _, commit := range c.CommitConfigs {
		if err := commit.PopulateOffChainAndOnChainCfg(state.MustGetEVMChainState(commit.DestinationChainSelector).PriceRegistry.Address()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate offchain and onchain config for commit: %w", err)
		}
		finalCfg, err := deriveOCR2Config(env, commit.DestinationChainSelector, commit.OCR2ConfigParams)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to derive OCR2 config for commit: %w", err)
		}
		commitStore := state.MustGetEVMChainState(commit.DestinationChainSelector).CommitStore[commit.SourceChainSelector]
		chain := env.BlockChains.EVMChains()[commit.DestinationChainSelector]
		tx, err := commitStore.SetOCR2Config(
			chain.DeployerKey,
			finalCfg.Signers,
			finalCfg.Transmitters,
			finalCfg.F,
			finalCfg.OnchainConfig,
			finalCfg.OffchainConfigVersion,
			finalCfg.OffchainConfig,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set OCR2 config for commit store %s on chain %s: %w",
				commitStore.Address().String(), chain.String(), cldf.MaybeDataErr(err))
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm OCR2 for commit store %s config on chain %s: %w",
				commitStore.Address().String(), chain.String(), err)
		}
	}
	for _, exec := range c.ExecConfigs {
		if err := exec.PopulateOffChainAndOnChainCfg(
			state.MustGetEVMChainState(exec.DestinationChainSelector).Router.Address(),
			state.MustGetEVMChainState(exec.DestinationChainSelector).PriceRegistry.Address()); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to populate offchain and onchain config for offramp: %w", err)
		}
		finalCfg, err := deriveOCR2Config(env, exec.DestinationChainSelector, exec.OCR2ConfigParams)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to derive OCR2 config for offramp: %w", err)
		}
		offRamp := state.MustGetEVMChainState(exec.DestinationChainSelector).EVM2EVMOffRamp[exec.SourceChainSelector]
		chain := env.BlockChains.EVMChains()[exec.DestinationChainSelector]
		tx, err := offRamp.SetOCR2Config(
			chain.DeployerKey,
			finalCfg.Signers,
			finalCfg.Transmitters,
			finalCfg.F,
			finalCfg.OnchainConfig,
			finalCfg.OffchainConfigVersion,
			finalCfg.OffchainConfig,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set OCR2 config for offramp %s on chain %s: %w",
				offRamp.Address().String(), chain.String(), err)
		}
		_, err = chain.Confirm(tx)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to confirm OCR2 for offramp %s config on chain %s: %w",
				offRamp.Address().String(), chain.String(), err)
		}
	}
	return cldf.ChangesetOutput{}, nil
}

func deriveOCR2Config(
	env cldf.Environment,
	chainSel uint64,
	ocrParams confighelper.PublicConfig,
) (FinalOCR2Config, error) {
	nodeInfo, err := deployment.NodeInfo(env.NodeIDs, env.Offchain)
	if err != nil {
		return FinalOCR2Config{}, fmt.Errorf("failed to get node info: %w", err)
	}
	nodes := nodeInfo.NonBootstraps()
	// Get OCR3 Config from helper
	var schedule []int
	var oracles []confighelper.OracleIdentityExtra
	for _, node := range nodes {
		schedule = append(schedule, 1)
		cfg, exists := node.OCRConfigForChainSelector(chainSel)
		if !exists {
			return FinalOCR2Config{}, fmt.Errorf("no OCR config for chain %d", chainSel)
		}
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  cfg.OnchainPublicKey,
				TransmitAccount:   cfg.TransmitAccount,
				OffchainPublicKey: cfg.OffchainPublicKey,
				PeerID:            cfg.PeerID.Raw(),
			},
			ConfigEncryptionPublicKey: cfg.ConfigEncryptionPublicKey,
		})
	}

	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		ocrParams.DeltaProgress,
		ocrParams.DeltaResend,
		ocrParams.DeltaRound,
		ocrParams.DeltaGrace,
		ocrParams.DeltaStage,
		ocrParams.RMax,
		schedule,
		oracles,
		ocrParams.ReportingPluginConfig,
		nil,
		ocrParams.MaxDurationQuery,
		ocrParams.MaxDurationObservation,
		ocrParams.MaxDurationReport,
		ocrParams.MaxDurationShouldAcceptFinalizedReport,
		ocrParams.MaxDurationShouldTransmitAcceptedReport,
		int(nodes.DefaultF()),
		ocrParams.OnchainConfig,
	)
	if err != nil {
		return FinalOCR2Config{}, fmt.Errorf("failed to derive OCR2 config: %w", err)
	}
	var signersAddresses []common.Address
	for _, signer := range signers {
		if len(signer) != 20 {
			return FinalOCR2Config{}, fmt.Errorf("address is not 20 bytes %s", signer)
		}
		signersAddresses = append(signersAddresses, common.BytesToAddress(signer))
	}
	var transmittersAddresses []common.Address
	for _, transmitter := range transmitters {
		bytes, err := hexutil.Decode(string(transmitter))
		if err != nil {
			return FinalOCR2Config{}, errors.Wrap(err, fmt.Sprintf("given address is not valid %s", transmitter))
		}
		if len(bytes) != 20 {
			return FinalOCR2Config{}, errors.Errorf("address is not 20 bytes %s", transmitter)
		}
		transmittersAddresses = append(transmittersAddresses, common.BytesToAddress(bytes))
	}
	return FinalOCR2Config{
		Signers:               signersAddresses,
		Transmitters:          transmittersAddresses,
		F:                     threshold,
		OnchainConfig:         onchainConfig,
		OffchainConfigVersion: offchainConfigVersion,
		OffchainConfig:        offchainConfig,
	}, nil
}
