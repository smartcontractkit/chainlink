// Package confighelper provides helpers for converting between the gethwrappers/OCR2Aggregator.SetConfig
// event and types.ContractConfig
package confighelper

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr2config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// OracleIdentity is identical to the internal type in package config.
// We intentionally make a copy to make potential future internal modifications easier.
type OracleIdentity struct {
	OffchainPublicKey types.OffchainPublicKey
	// For EVM-chains, this an *address*.
	OnchainPublicKey types.OnchainPublicKey
	PeerID           string
	TransmitAccount  types.Account
}

// PublicConfig is identical to the internal type in package config.
// We intentionally make a copy to make potential future internal modifications easier.
type PublicConfig struct {
	DeltaProgress    time.Duration
	DeltaResend      time.Duration
	DeltaRound       time.Duration
	DeltaGrace       time.Duration
	DeltaStage       time.Duration
	RMax             uint8
	S                []int
	OracleIdentities []OracleIdentity

	ReportingPluginConfig []byte

	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationReport                       time.Duration
	MaxDurationShouldAcceptFinalizedReport  time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration

	F             int
	OnchainConfig []byte
	ConfigDigest  types.ConfigDigest
}

func (pc PublicConfig) N() int {
	return len(pc.OracleIdentities)
}

func PublicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	internalPublicConfig, err := ocr2config.PublicConfigFromContractConfig(skipResourceExhaustionChecks, change)
	if err != nil {
		return PublicConfig{}, err
	}
	identities := []OracleIdentity{}
	for _, internalIdentity := range internalPublicConfig.OracleIdentities {
		identities = append(identities, OracleIdentity{
			internalIdentity.OffchainPublicKey,
			internalIdentity.OnchainPublicKey,
			internalIdentity.PeerID,
			internalIdentity.TransmitAccount,
		})
	}
	return PublicConfig{
		internalPublicConfig.DeltaProgress,
		internalPublicConfig.DeltaResend,
		internalPublicConfig.DeltaRound,
		internalPublicConfig.DeltaGrace,
		internalPublicConfig.DeltaStage,
		internalPublicConfig.RMax,
		internalPublicConfig.S,
		identities,
		internalPublicConfig.ReportingPluginConfig,
		internalPublicConfig.MaxDurationQuery,
		internalPublicConfig.MaxDurationObservation,
		internalPublicConfig.MaxDurationReport,
		internalPublicConfig.MaxDurationShouldAcceptFinalizedReport,
		internalPublicConfig.MaxDurationShouldTransmitAcceptedReport,
		internalPublicConfig.F,
		internalPublicConfig.OnchainConfig,
		internalPublicConfig.ConfigDigest,
	}, nil
}

type OracleIdentityExtra struct {
	OracleIdentity
	ConfigEncryptionPublicKey types.ConfigEncryptionPublicKey
}

// ContractSetConfigArgsForIntegrationTest generates setConfig args for integration tests in core.
// Only use this for testing, *not* for production.
func ContractSetConfigArgsForEthereumIntegrationTest(
	oracles []OracleIdentityExtra,
	f int,
	alphaPPB uint64,
) (
	signers []common.Address,
	transmitters []common.Address,
	f_ uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	S := []int{}
	identities := []config.OracleIdentity{}
	sharedSecretEncryptionPublicKeys := []types.ConfigEncryptionPublicKey{}
	for _, oracle := range oracles {
		S = append(S, 1)
		identities = append(identities, config.OracleIdentity{
			oracle.OffchainPublicKey,
			oracle.OnchainPublicKey,
			oracle.PeerID,
			oracle.TransmitAccount,
		})
		sharedSecretEncryptionPublicKeys = append(sharedSecretEncryptionPublicKeys, oracle.ConfigEncryptionPublicKey)
	}
	sharedConfig := ocr2config.SharedConfig{
		ocr2config.PublicConfig{
			2 * time.Second,
			1 * time.Second,
			1 * time.Second,
			500 * time.Millisecond,
			2 * time.Second,
			3,
			S,
			identities,
			median.OffchainConfig{
				false,
				alphaPPB,
				false,
				alphaPPB,
				0,
			}.Encode(),
			50 * time.Millisecond,
			50 * time.Millisecond,
			50 * time.Millisecond,
			50 * time.Millisecond,
			50 * time.Millisecond,
			f,
			nil, // The median reporting plugin has an empty onchain config
			types.ConfigDigest{},
		},
		&[config.SharedSecretSize]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8},
	}
	setConfigArgs, err := ocr2config.XXXContractSetConfigArgsFromSharedConfigEthereum(sharedConfig, sharedSecretEncryptionPublicKeys)
	return setConfigArgs.Signers,
		setConfigArgs.Transmitters,
		setConfigArgs.F,
		setConfigArgs.OnchainConfig,
		setConfigArgs.OffchainConfigVersion,
		setConfigArgs.OffchainConfig,
		err
}

// ContractSetConfigArgsForTestsWithAuxiliaryArgs generates setConfig args from
// the relevant parameters. Only use this for testing, *not* for production.
func ContractSetConfigArgsForTestsWithAuxiliaryArgs(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaStage time.Duration,
	rMax uint8,
	s []int,
	oracles []OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationQuery time.Duration,
	maxDurationObservation time.Duration,
	maxDurationReport time.Duration,
	maxDurationShouldAcceptFinalizedReport time.Duration,
	maxDurationShouldTransmitAcceptedReport time.Duration,

	f int,
	onchainConfig []byte,
	auxiliaryArgs AuxiliaryArgs,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	identities := []config.OracleIdentity{}
	configEncryptionPublicKeys := []types.ConfigEncryptionPublicKey{}
	for _, oracle := range oracles {
		identities = append(identities, config.OracleIdentity{
			oracle.OffchainPublicKey,
			oracle.OnchainPublicKey,
			oracle.PeerID,
			oracle.TransmitAccount,
		})
		configEncryptionPublicKeys = append(configEncryptionPublicKeys, oracle.ConfigEncryptionPublicKey)
	}

	sharedSecret := [config.SharedSecretSize]byte{}
	if _, err := io.ReadFull(auxiliaryArgs.rng(), sharedSecret[:]); err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}

	sharedConfig := ocr2config.SharedConfig{
		ocr2config.PublicConfig{
			deltaProgress,
			deltaResend,
			deltaRound,
			deltaGrace,
			deltaStage,
			rMax,
			s,
			identities,
			reportingPluginConfig,
			maxDurationQuery,
			maxDurationObservation,
			maxDurationReport,
			maxDurationShouldAcceptFinalizedReport,
			maxDurationShouldTransmitAcceptedReport,
			f,
			onchainConfig,
			types.ConfigDigest{},
		},
		&sharedSecret,
	}
	return ocr2config.XXXContractSetConfigArgsFromSharedConfig(sharedConfig, configEncryptionPublicKeys)
}

// AuxiliaryArgs provides keyword-style extra configuration for calls to
// ContractSetConfigArgsForTests
type AuxiliaryArgs struct {
	RNG io.Reader
}

func (a AuxiliaryArgs) rng() io.Reader {
	if a.RNG == nil {
		return rand.Reader
	}
	return a.RNG
}

// ContractSetConfigArgsForTests generates setConfig args from the relevant
// parameters. Only use this for testing, *not* for production.
func ContractSetConfigArgsForTests(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaStage time.Duration,
	rMax uint8,
	s []int,
	oracles []OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationQuery time.Duration,
	maxDurationObservation time.Duration,
	maxDurationReport time.Duration,
	maxDurationShouldAcceptFinalizedReport time.Duration,
	maxDurationShouldTransmitAcceptedReport time.Duration,

	f int,
	onchainConfig []byte,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	return ContractSetConfigArgsForTestsWithAuxiliaryArgs(
		deltaProgress,
		deltaResend,
		deltaRound,
		deltaGrace,
		deltaStage,
		rMax,
		s,
		oracles,
		reportingPluginConfig,
		maxDurationQuery,
		maxDurationObservation,
		maxDurationReport,
		maxDurationShouldAcceptFinalizedReport,
		maxDurationShouldTransmitAcceptedReport,
		f,
		onchainConfig,
		AuxiliaryArgs{},
	)
}

// Deprecated: Use corresponding function in ocr3confighelper
func ContractSetConfigArgsForTestsMercuryV02(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaInitial time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaCertifiedCommitRequest time.Duration,
	deltaStage time.Duration,
	rMax uint8,
	s []int,
	oracles []OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationObservation time.Duration,
	f int,
	onchainConfig []byte,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	return ContractSetConfigArgsForTestsOCR3(
		deltaProgress,
		deltaResend,
		deltaInitial,
		deltaRound,
		deltaGrace,
		deltaCertifiedCommitRequest,
		deltaStage,
		uint64(rMax),
		s,
		oracles,
		reportingPluginConfig,
		0,
		maxDurationObservation,
		0,
		0,
		f,
		onchainConfig,
	)
}

// Deprecated: Use corresponding function in ocr3confighelper
func ContractSetConfigArgsForTestsOCR3(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaInitial time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaCertifiedCommitRequest time.Duration,
	deltaStage time.Duration,
	rMax uint64,
	s []int,
	oracles []OracleIdentityExtra,
	reportingPluginConfig []byte,
	maxDurationQuery time.Duration,
	maxDurationObservation time.Duration,
	maxDurationShouldAcceptAttestedReport time.Duration,
	maxDurationShouldTransmitAcceptedReport time.Duration,
	f int,
	onchainConfig []byte,
) (
	signers []types.OnchainPublicKey,
	transmitters []types.Account,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	identities := []config.OracleIdentity{}
	configEncryptionPublicKeys := []types.ConfigEncryptionPublicKey{}
	for _, oracle := range oracles {
		identities = append(identities, config.OracleIdentity{
			oracle.OffchainPublicKey,
			oracle.OnchainPublicKey,
			oracle.PeerID,
			oracle.TransmitAccount,
		})
		configEncryptionPublicKeys = append(configEncryptionPublicKeys, oracle.ConfigEncryptionPublicKey)
	}

	sharedSecret := [config.SharedSecretSize]byte{}
	if _, err := io.ReadFull(rand.Reader, sharedSecret[:]); err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}

	sharedConfig := ocr3config.SharedConfig{
		ocr3config.PublicConfig{
			deltaProgress,
			deltaResend,
			deltaInitial,
			deltaRound,
			deltaGrace,
			deltaCertifiedCommitRequest,
			deltaStage,
			rMax,
			s,
			identities,
			reportingPluginConfig,
			maxDurationQuery,
			maxDurationObservation,
			maxDurationShouldAcceptAttestedReport,
			maxDurationShouldTransmitAcceptedReport,
			f,
			onchainConfig,
			types.ConfigDigest{},
		},
		&sharedSecret,
	}
	return ocr3config.XXXContractSetConfigArgsFromSharedConfig(sharedConfig, configEncryptionPublicKeys)
}
