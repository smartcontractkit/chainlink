// alias for offchainreporting2plus/confighelper
package confighelper

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

type OracleIdentity = confighelper.OracleIdentity

type PublicConfig = confighelper.PublicConfig

func PublicConfigFromContractConfig(skipResourceExhaustionChecks bool, change types.ContractConfig) (PublicConfig, error) {
	return confighelper.PublicConfigFromContractConfig(skipResourceExhaustionChecks, change)
}

type OracleIdentityExtra = confighelper.OracleIdentityExtra

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
	return confighelper.ContractSetConfigArgsForEthereumIntegrationTest(oracles, f, alphaPPB)
}

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
	return confighelper.ContractSetConfigArgsForTestsWithAuxiliaryArgs(
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
		auxiliaryArgs,
	)
}

type AuxiliaryArgs = confighelper.AuxiliaryArgs

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
	return confighelper.ContractSetConfigArgsForTests(
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
	)
}
