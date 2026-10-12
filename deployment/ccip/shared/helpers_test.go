package shared

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
)

func TestResolveFeeQuoterAddressAndVersion(t *testing.T) {
	version := semver.MustParse("2.0.0")
	older := semver.MustParse("1.6.0")
	refs := []datastore.AddressRef{
		{ChainSelector: 1, Type: datastore.ContractType("FeeQuoter"), Version: older, Address: "0x0000000000000000000000000000000000000001"},
		{ChainSelector: 1, Type: datastore.ContractType("FeeQuoter"), Version: version, Address: "0x0000000000000000000000000000000000000002"},
		{ChainSelector: 1, Type: datastore.ContractType("FeeQuoter"), Version: version, Address: "0x0000000000000000000000000000000000000003", Labels: datastore.NewLabelSet(SupersededLabel)},
	}

	address, gotVersion, err := ResolveFeeQuoterAddressAndVersion(refs, 1)
	require.NoError(t, err)
	require.Equal(t, "0x0000000000000000000000000000000000000002", address.Hex())
	require.Equal(t, *version, gotVersion)
}

func TestResolveFeeQuoterAddressAndVersionRejectsAmbiguity(t *testing.T) {
	version := semver.MustParse("2.0.0")
	refs := []datastore.AddressRef{
		{ChainSelector: 1, Type: datastore.ContractType("FeeQuoter"), Version: version, Address: "0x0000000000000000000000000000000000000002"},
		{ChainSelector: 1, Type: datastore.ContractType("FeeQuoter"), Version: version, Address: "0x0000000000000000000000000000000000000003"},
	}

	_, _, err := ResolveFeeQuoterAddressAndVersion(refs, 1)
	require.ErrorContains(t, err, "ambiguous fee quoter")
}
