package v1_5

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccip"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
)

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
	CommitOffchainConfig     ccip.CommitOffchainConfig
	OCR2ConfigParams         confighelper.PublicConfig
}

func (c CommitOCR2ConfigParams) SetCommitOffChainCfg() error {
	cfgBytes, err := ccipconfig.EncodeOffchainConfig(
		testhelpers.NewCommitOffchainConfig(
			c.CommitOffchainConfig.GasPriceHeartBeat,
			c.CommitOffchainConfig.GasPriceDeviationPPB,
		))
}

func (c CommitOCR2ConfigParams) Validate(state changeset.CCIPOnChainState) error {
	if err := deployment.IsValidChainSelector(c.DestinationChainSelector); err != nil {
		return fmt.Errorf("invalid DestinationChainSelector: %w", err)
	}

	chain, exists := state.Chains[c.DestinationChainSelector]
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
	// TODO : add validation for rest of the configs
	return nil
}

type ExecuteOCR2ConfigParams struct {
	DestinationChainSelector uint64
	SourceChainSelector      uint64
	ExecOffchainConfig       ccip.ExecOffchainConfig
	ExecOnchainConfig        ccip.ExecOnchainConfig
	OCR2ConfigParams         types.OCRParameters
}

func (e ExecuteOCR2ConfigParams) Validate(state changeset.CCIPOnChainState) error {
	if err := e.OCR2ConfigParams.Validate(); err != nil {
		return err
	}
	if err := e.ExecOnchainConfig.Validate(); err != nil {
		return err
	}
	chain, exists := state.Chains[e.DestinationChainSelector]
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
	// TODO : add validation for rest of the configs
	return nil
}

type OCR2Config struct {
	CommitConfigs []CommitOCR2ConfigParams
	ExecConfigs   []ExecuteOCR2ConfigParams
}

func (o OCR2Config) Validate(state changeset.CCIPOnChainState) error {
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

func DeriveOCR2Config(
	env deployment.Environment,
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
				PeerID:            cfg.PeerID.String()[4:],
			}, ConfigEncryptionPublicKey: cfg.ConfigEncryptionPublicKey,
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
		ocrParams.F,
		ocrParams.OnchainConfig,
	)
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

// SetOCR2Config sets the OCR2 config on the chain for commit and offramp
// This is currently not suitable for prod environments it's only for testing
func SetOCR2Config(env deployment.Environment, c OCR2Config) (deployment.ChangesetOutput, error) {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load CCIP onchain state: %w", err)
	}
	if err := c.Validate(state); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid OCR2 config: %w", err)
	}
	for _, commit := range c.CommitConfigs {
		if err := setOCR2ConfigCommit(env, state, commit); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to set OCR2 config commit: %w", err)
		}
	}
	for _, exec := range c.ExecConfigs {
		if err := setOCR2ConfigExec(env, state, exec); err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to set OCR2 config exec: %w", err)
		}
	}
	return deployment.ChangesetOutput{}, nil
}

func setOCR2ConfigExec(env deployment.Environment, state changeset.CCIPOnChainState, c ExecuteOCR2ConfigParams) error {

	return nil
}

func setOCR2ConfigCommit(env deployment.Environment, state changeset.CCIPOnChainState, c CommitOCR2ConfigParams) error {
	return nil
}
