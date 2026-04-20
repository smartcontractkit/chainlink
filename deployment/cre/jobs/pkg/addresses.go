package pkg

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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
