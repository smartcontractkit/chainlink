package pkg

import (
	"github.com/Masterminds/semver/v3"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func GetOCR3CapabilityV2AddressRefKey(chainSel uint64, qualifier string) datastore.AddressRefKey {
	return datastore.NewAddressRefKey(
		chainSel,
		"OCR3Capability",
		semver.MustParse("2.0.0"),
		qualifier,
	)
}
