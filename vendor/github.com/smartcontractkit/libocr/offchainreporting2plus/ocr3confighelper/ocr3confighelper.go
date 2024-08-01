package ocr3confighelper

import (
	"crypto/rand"
	"io"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/internal/config/ocr3config"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// PublicConfig is identical to the internal type in package config.
// We intentionally make a copy to make potential future internal modifications easier.
type PublicConfig struct {
	DeltaProgress               time.Duration
	DeltaResend                 time.Duration
	DeltaInitial                time.Duration
	DeltaRound                  time.Duration
	DeltaGrace                  time.Duration
	DeltaCertifiedCommitRequest time.Duration
	DeltaStage                  time.Duration
	RMax                        uint64
	S                           []int
	OracleIdentities            []confighelper.OracleIdentity

	ReportingPluginConfig []byte

	MaxDurationQuery                        time.Duration
	MaxDurationObservation                  time.Duration
	MaxDurationShouldAcceptAttestedReport   time.Duration
	MaxDurationShouldTransmitAcceptedReport time.Duration

	F             int
	OnchainConfig []byte
	ConfigDigest  types.ConfigDigest
}

func (pc PublicConfig) N() int {
	return len(pc.OracleIdentities)
}

func PublicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	internalPublicConfig, err := ocr3config.PublicConfigFromContractConfig(skipResourceExhaustionChecks, change)
	if err != nil {
		return PublicConfig{}, err
	}
	identities := []confighelper.OracleIdentity{}
	for _, internalIdentity := range internalPublicConfig.OracleIdentities {
		identities = append(identities, confighelper.OracleIdentity{
			internalIdentity.OffchainPublicKey,
			internalIdentity.OnchainPublicKey,
			internalIdentity.PeerID,
			internalIdentity.TransmitAccount,
		})
	}
	return PublicConfig{
		internalPublicConfig.DeltaProgress,
		internalPublicConfig.DeltaResend,
		internalPublicConfig.DeltaInitial,
		internalPublicConfig.DeltaRound,
		internalPublicConfig.DeltaGrace,
		internalPublicConfig.DeltaCertifiedCommitRequest,
		internalPublicConfig.DeltaStage,
		internalPublicConfig.RMax,
		internalPublicConfig.S,
		identities,
		internalPublicConfig.ReportingPluginConfig,
		internalPublicConfig.MaxDurationQuery,
		internalPublicConfig.MaxDurationObservation,
		internalPublicConfig.MaxDurationShouldAcceptAttestedReport,
		internalPublicConfig.MaxDurationShouldTransmitAcceptedReport,
		internalPublicConfig.F,
		internalPublicConfig.OnchainConfig,
		internalPublicConfig.ConfigDigest,
	}, nil
}

// ContractSetConfigArgsForTestsWithAuxiliaryArgsMercuryV02 generates setConfig
// args for mercury v0.2. Only use this for testing, *not* for production.
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
	oracles []confighelper.OracleIdentityExtra,
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
	return ContractSetConfigArgsForTests(
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

// ContractSetConfigArgsForTestsOCR3 generates setConfig args for OCR3. Only use
// this for testing, *not* for production.
func ContractSetConfigArgsForTests(
	deltaProgress time.Duration,
	deltaResend time.Duration,
	deltaInitial time.Duration,
	deltaRound time.Duration,
	deltaGrace time.Duration,
	deltaCertifiedCommitRequest time.Duration,
	deltaStage time.Duration,
	rMax uint64,
	s []int,
	oracles []confighelper.OracleIdentityExtra,
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
