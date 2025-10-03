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
