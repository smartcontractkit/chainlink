// alias for offchainreporting2plus/types
package types

import "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

type ConfigDigestPrefix = types.ConfigDigestPrefix

const (
	// Deprecated: Use equivalent offchainreporting2plus/types.ConfigDigestPrefixEVMSimple instead
	ConfigDigestPrefixEVM        ConfigDigestPrefix = types.ConfigDigestPrefixEVM //nolint:staticcheck
	ConfigDigestPrefixSolana     ConfigDigestPrefix = types.ConfigDigestPrefixSolana
	ConfigDigestPrefixStarknet   ConfigDigestPrefix = types.ConfigDigestPrefixStarknet
	ConfigDigestPrefixMercuryV02 ConfigDigestPrefix = types.ConfigDigestPrefixMercuryV02
	ConfigDigestPrefixOCR1       ConfigDigestPrefix = types.ConfigDigestPrefixOCR1
)

type ConfigDigest = types.ConfigDigest

func BytesToConfigDigest(b []byte) (ConfigDigest, error) {
	return types.BytesToConfigDigest(b)
}

type OffchainConfigDigester = types.OffchainConfigDigester

const MaxOracles = types.MaxOracles

type ConfigDatabase = types.ConfigDatabase

type Database = types.Database

type PendingTransmission = types.PendingTransmission

type PersistentState = types.PersistentState

const EnableDangerousDevelopmentMode = types.EnableDangerousDevelopmentMode

type LocalConfig = types.LocalConfig

type BinaryNetworkEndpointLimits = types.BinaryNetworkEndpointLimits

type BinaryNetworkEndpointFactory = types.BinaryNetworkEndpointFactory

type BootstrapperFactory = types.BootstrapperFactory

type Query = types.Query

type Observation = types.Observation

type AttributedObservation = types.AttributedObservation

type ReportTimestamp = types.ReportTimestamp

type ReportContext = types.ReportContext

type Report = types.Report

type AttributedOnchainSignature = types.AttributedOnchainSignature

type ReportingPluginFactory = types.ReportingPluginFactory

type ReportingPluginConfig = types.ReportingPluginConfig

type ReportingPlugin = types.ReportingPlugin

const (
	MaxMaxQueryLength       = types.MaxMaxQueryLength
	MaxMaxObservationLength = types.MaxMaxObservationLength
	MaxMaxReportLength      = types.MaxMaxReportLength
)

type ReportingPluginLimits = types.ReportingPluginLimits

type ReportingPluginInfo = types.ReportingPluginInfo

type Account = types.Account

type ContractTransmitter = types.ContractTransmitter

type ContractConfigTracker = types.ContractConfigTracker

type ContractConfig = types.ContractConfig

type OffchainPublicKey = types.OffchainPublicKey

type OnchainPublicKey = types.OnchainPublicKey

type ConfigEncryptionPublicKey = types.ConfigEncryptionPublicKey

type OffchainKeyring = types.OffchainKeyring

type OnchainKeyring = types.OnchainKeyring
