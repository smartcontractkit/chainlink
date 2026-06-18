package pkg

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink/deployment/cre/forwarder/solana"
)

func GetOCR3CapabilityAddressRefKey(chainSel uint64, qualifier string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(
		chainSel,
		"OCR3Capability",
		semver.MustParse("1.0.0"),
		qualifier,
	)
}

func GetKeystoneForwarderCapabilityAddressRefKey(chainSel uint64, qualifier string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(
		chainSel,
		"KeystoneForwarder",
		semver.MustParse("1.0.0"),
		qualifier,
	)
}

// Solana CRE forwarder entries match cre/forwarder/solana.DeployForwarder (program + state).
const (
	SolanaForwarderProgramType datastore.ContractType = solana.ForwarderContract
	SolanaForwarderStateType   datastore.ContractType = solana.ForwarderState
)

func GetSolanaForwarderProgramRefKey(chainSel uint64, version *semver.Version, qualifier string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(chainSel, SolanaForwarderProgramType, version, qualifier)
}

func GetSolanaForwarderStateRefKey(chainSel uint64, version *semver.Version, qualifier string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(chainSel, SolanaForwarderStateType, version, qualifier)
}

func GetCapRegAddressRefKey(chainSel uint64, qualifier string, version string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(
		chainSel,
		"CapabilitiesRegistry",
		semver.MustParse(version),
		qualifier,
	)
}

func GetShardConfigAddressRefKey(chainSel uint64, qualifier string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(
		chainSel,
		"ShardConfig",
		semver.MustParse("1.0.0"),
		qualifier,
	)
}
