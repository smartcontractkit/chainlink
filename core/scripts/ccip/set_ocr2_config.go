package main

import (
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	ocrconfighelper "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/dione"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip/shared"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

type OCR2Params struct {
	// If an epoch (driven by a leader) fails to achieve progress (generate a
	// report) after DeltaProgress, we enter a new epoch. This parameter must be
	// chosen carefully. If the duration is too short, we may keep prematurely
	// switching epochs without ever achieving any progress, resulting in a
	// liveness failure!
	deltaProgress time.Duration
	// deltaResend determines how often Pacemaker newepoch messages should be
	// resent, allowing oracles that had crashed and are recovering to rejoin
	// the protocol more quickly. ~30s should be a reasonable default under most
	// circumstances.
	deltaResend time.Duration
	// deltaRound determines the minimal amount of time that should pass between
	// the start of report generation rounds. With OCR2 only (not OCR1!) you can
	// set this value very aggressively. Note that this only provides a lower
	// bound on the round interval; actual rounds might take longer.
	deltaRound time.Duration
	// Once the leader of a report generation round has collected sufficiently
	// many observations, it will wait for DeltaGrace to pass to allow slower
	// oracles to still contribute an observation before moving on to generating
	// the report. Consequently, rounds driven by correct leaders will always
	// take at least DeltaGrace.
	deltaGrace time.Duration

	query          time.Duration
	observation    time.Duration
	report         time.Duration
	shouldAccept   time.Duration
	shouldTransmit time.Duration
}

var (
	CommitOcr2Params = OCR2Params{
		deltaProgress:  2 * time.Minute,
		deltaResend:    5 * time.Second,
		deltaRound:     75 * time.Second,
		deltaGrace:     5 * time.Second,
		query:          100 * time.Millisecond, // commit does not use query
		observation:    35 * time.Second,
		report:         10 * time.Second,
		shouldAccept:   5 * time.Second,
		shouldTransmit: 10 * time.Second,
	}
	ExecOcr2Params = OCR2Params{
		deltaProgress:  100 * time.Second,
		deltaResend:    5 * time.Second,
		deltaRound:     40 * time.Second,
		deltaGrace:     5 * time.Second,
		query:          100 * time.Millisecond, // exec does not use query
		observation:    20 * time.Second,
		report:         8 * time.Second,
		shouldAccept:   5 * time.Second,
		shouldTransmit: 8 * time.Second,
	}
)

func (client *CCIPClient) SetOCR2Config(env dione.Environment) {
	verifierOCRConfig, err := client.Dest.CommitStore.LatestConfigDetails(&bind.CallOpts{})
	helpers.PanicErr(err)
	if verifierOCRConfig.BlockNumber != 0 {
		client.Dest.logger.Infof("CommitStore OCR config already found: %+v", verifierOCRConfig.ConfigDigest)
		client.Dest.logger.Infof("The new config will overwrite the current one.")
	}

	rampOCRConfig, err := client.Dest.OffRamp.LatestConfigDetails(&bind.CallOpts{})
	helpers.PanicErr(err)
	if rampOCRConfig.BlockNumber != 0 {
		client.Dest.logger.Infof("OffRamp OCR config already found: %+v", rampOCRConfig.ConfigDigest)
		client.Dest.logger.Infof("The new config will overwrite the current one.")
	}
	if client.Dest.Client.ChainId == 1337 || client.Source.Client.ChainId == 1337 {
		env = dione.Prod_Swift
	}
	don := dione.NewOfflineDON(env, client.Dest.logger)
	faults := len(don.Config.Nodes) / 3

	tx, err := client.setOCRConfig(client.Dest.CommitStore, client.getCommitStoreOffChainConfig(), client.getCommitStoreOnchainConfig(), CommitOcr2Params, faults, don.GenerateOracleIdentities(client.Dest.ChainId))
	helpers.PanicErr(err)
	client.Dest.logger.Infof("Config set on commitStore %s", helpers.ExplorerLink(int64(client.Dest.ChainId), tx.Hash()))

	tx, err = client.setOCRConfig(client.Dest.OffRamp, client.getOffRampOffChainConfig(), client.getOffRampOnchainConfig(), ExecOcr2Params, faults, don.GenerateOracleIdentities(client.Dest.ChainId))
	helpers.PanicErr(err)
	client.Dest.logger.Infof("Config set on offramp %s", helpers.ExplorerLink(int64(client.Dest.ChainId), tx.Hash()))
}

func (client *CCIPClient) setOCRConfig(ocrConf ocr2Configurer, pluginOffchainConfig []byte, onchainConfig []byte, ocr2Params OCR2Params, faults int, identities []ocrconfighelper.OracleIdentityExtra) (*types.Transaction, error) {
	// Simple transmission schedule of 1 node per stage.
	// sum(transmissionSchedule) should equal number of nodes.
	var transmissionSchedule []int
	for i := 0; i < len(identities); i++ {
		transmissionSchedule = append(transmissionSchedule, 1)
	}
	signers, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocrconfighelper.ContractSetConfigArgsForTests(
		ocr2Params.deltaProgress,
		ocr2Params.deltaResend,
		ocr2Params.deltaRound,
		ocr2Params.deltaGrace,
		client.Dest.TunableValues.InflightCacheExpiry.Duration(), // deltaStage
		3,
		transmissionSchedule,
		identities,
		pluginOffchainConfig,
		ocr2Params.query,
		ocr2Params.observation,
		ocr2Params.report,
		ocr2Params.shouldAccept,
		ocr2Params.shouldTransmit,
		faults,
		onchainConfig,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create args for ocr config tx")
	}
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	if err != nil {
		return nil, err
	}
	transmitterAddresses, err := evm.AccountToAddress(transmitters)
	if err != nil {
		return nil, err
	}

	tx, err := ocrConf.SetOCR2Config(
		client.Dest.Owner,
		signerAddresses,
		transmitterAddresses,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	if err != nil {
		return nil, err
	}
	err = shared.WaitForMined(client.Dest.logger, client.Dest.Client.Client, tx.Hash(), true)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (client *CCIPClient) getCommitStoreOnchainConfig() []byte {
	commitStoreOnchainConfig := ccipconfig.CommitOnchainConfig{
		PriceRegistry: client.Dest.PriceRegistry.Address(),
	}

	encodedCommitStoreOnchainConfig, err := abihelpers.EncodeAbiStruct(commitStoreOnchainConfig)
	helpers.PanicErr(err)

	return encodedCommitStoreOnchainConfig
}

func (client *CCIPClient) getCommitStoreOffChainConfig() []byte {
	if client.Source.TunableValues.FinalityDepth == 0 || client.Dest.TunableValues.FinalityDepth == 0 {
		panic("Please set the tunable chain values")
	}

	commitPluginConfig := ccipconfig.CommitOffchainConfig{
		SourceFinalityDepth:   client.Source.TunableValues.FinalityDepth,
		DestFinalityDepth:     client.Dest.TunableValues.FinalityDepth,
		FeeUpdateHeartBeat:    client.Dest.TunableValues.FeeUpdateHeartBeat,
		FeeUpdateDeviationPPB: client.Dest.TunableValues.FeeUpdateDeviationPPB,
		MaxGasPrice:           client.Dest.TunableValues.MaxGasPrice,
		InflightCacheExpiry:   client.Dest.TunableValues.InflightCacheExpiry,
	}

	encodedOffchainConfig, err := ccipconfig.EncodeOffchainConfig(commitPluginConfig)
	helpers.PanicErr(err)

	return encodedOffchainConfig
}

func (client *CCIPClient) getOffRampOnchainConfig() []byte {
	offRampOnchainConfig := ccipconfig.ExecOnchainConfig{
		PermissionLessExecutionThresholdSeconds: rhea.PERMISSIONLESS_EXEC_THRESHOLD_SEC,
		Router:                                  client.Dest.Router.Address(),
		PriceRegistry:                           client.Dest.PriceRegistry.Address(),
		MaxTokensLength:                         rhea.MAX_TOKEN_LENGTH,
		MaxDataSize:                             rhea.MAX_DATA_SIZE,
	}

	encodedOffRampOnchainConfig, err := abihelpers.EncodeAbiStruct(offRampOnchainConfig)
	helpers.PanicErr(err)

	return encodedOffRampOnchainConfig
}

func (client *CCIPClient) getOffRampOffChainConfig() []byte {
	if client.Source.TunableValues.FinalityDepth == 0 || client.Dest.TunableValues.FinalityDepth == 0 {
		panic("Please set the tunable chain values")
	}
	execPluginConfig := ccipconfig.ExecOffchainConfig{
		SourceFinalityDepth:         client.Source.TunableValues.FinalityDepth,
		DestFinalityDepth:           client.Dest.TunableValues.FinalityDepth,
		DestOptimisticConfirmations: client.Dest.TunableValues.OptimisticConfirmations,
		BatchGasLimit:               client.Dest.TunableValues.BatchGasLimit,
		RelativeBoostPerWaitHour:    client.Dest.TunableValues.RelativeBoostPerWaitHour,
		MaxGasPrice:                 client.Dest.TunableValues.MaxGasPrice,
		InflightCacheExpiry:         client.Dest.TunableValues.InflightCacheExpiry,
		RootSnoozeTime:              client.Dest.TunableValues.RootSnoozeTime,
	}

	encodedOffRampConfig, err := ccipconfig.EncodeOffchainConfig(execPluginConfig)
	helpers.PanicErr(err)

	return encodedOffRampConfig
}
