package pkg

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func GetCapRegV2AddressRefKey(chainSel uint64, qualifier string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(
		chainSel,
		"CapabilitiesRegistry",
		semver.MustParse("2.0.0"),
		qualifier,
	)
}
