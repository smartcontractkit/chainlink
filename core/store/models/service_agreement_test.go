package models_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUnsignedServiceAgreementFromRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		wantDigest  string
		wantPayment *assets.Link
	}{
		{
			"basic",
			`{"payment":"1","initiators":[{"type":"web"}],"tasks":[` +
				`{"type":"httpget","url":"https://bitstamp.net/api/ticker/"},` +
				`{"type":"jsonparse","path":["last"]},` +
				`{"type":"ethbytes32"},{"type":"ethtx"}` +
				`],` +
				`"aggregator":"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",` +
				`"aggInitiateJobSelector":"0x12345678",` +
				`"aggFulfillSelector":"0x87654321"` +
				`}`,
			"0xad12826461f2259eac07e762d9f1d32dd6af2e4ed0797b08cb3d8a8c3c4dd61d",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.JobSpecRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			us, err := models.NewUnsignedServiceAgreementFromRequest(strings.NewReader(test.input))
			require.NoError(t, err)
			assert.Equal(t, test.wantDigest, us.ID.String())
			assert.Equal(t, test.wantPayment, us.Encumbrance.Payment)
			assert.Equal(t, cltest.NormalizedJSON(t, []byte(test.input)), us.RequestBody)
		})
	}
}

func TestBuildServiceAgreement(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		wantDigest  string
		wantPayment *assets.Link
	}{
		{
			"basic",
			`{"payment":"1","initiators":[{"type":"web"}],"tasks":[` +
				`{"type":"httpget","url":"https://bitstamp.net/api/ticker/"},` +
				`{"type":"jsonparse","path":["last"]},` +
				`{"type":"ethbytes32"},{"type":"ethtx"}` +
				`],` +
				`"aggregator":"0xDeaDbeefdEAdbeefdEadbEEFdeadbeEFdEaDbeeF",` +
				`"aggInitiateJobSelector":"0x12345678",` +
				`"aggFulfillSelector":"0x87654321"` +
				`}`,
			"0xad12826461f2259eac07e762d9f1d32dd6af2e4ed0797b08cb3d8a8c3c4dd61d",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.JobSpecRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			us, err := models.NewUnsignedServiceAgreementFromRequest(strings.NewReader(test.input))
			require.NoError(t, err)

			sa, err := models.BuildServiceAgreement(us, cltest.MockSigner{})
			require.NoError(t, err)
			assert.Equal(t, test.wantDigest, sa.ID)
			assert.Equal(t, test.wantPayment, sa.Encumbrance.Payment)
			assert.Equal(t, cltest.NormalizedJSON(t, []byte(test.input)), sa.RequestBody)
			assert.NotEqual(t, models.AnyTime{}, sa.CreatedAt)
			assert.NotEqual(t, "", sa.Signature.String())
		})
	}
}

func TestEncumbrance_ABI(t *testing.T) {
	t.Parallel()
	endAt := cltest.ParseISO8601(t, "2007-01-02T15:04:05.000Z")
	digest := cltest.NewHash()
	requestDigestHex := utils.RemoveHexPrefix(digest.Hex())

	tests := []struct {
		name                   string
		payment                *assets.Link
		expiration             int
		endAt                  models.AnyTime
		oracles                []ethkey.EIP55Address
		aggregator             string
		aggInitiateJobSelector string
		aggFulfillSelector     string
		want                   string
	}{
		{"basic", assets.NewLink(1), 2, models.AnyTime{}, nil,
			"0x0000000000000000000000000000000000000000000000000000000000000000", "0x00000000", "0x00000000",
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000001" + // Payment
				"0000000000000000000000000000000000000000000000000000000000000002" + // Expiration time
				"0000000000000000000000000000000000000000000000000000000000000000" + // EndAt
				requestDigestHex +
				"0000000000000000000000000000000000000000" + // Aggregator address
				"00000000" + "00000000", // Function selectors
		},
		{"basic dead beef payment", assets.NewLink(3735928559), 2, models.AnyTime{}, nil,
			"0x0000000000000000000000000000000000000000000000000000000000000000", "0x00000000", "0x00000000",
			"0x" +
				"00000000000000000000000000000000000000000000000000000000deadbeef" + // Payment
				"0000000000000000000000000000000000000000000000000000000000000002" + // Expiration time
				"0000000000000000000000000000000000000000000000000000000000000000" + // EndAt
				requestDigestHex +
				"0000000000000000000000000000000000000000" + // Aggregator address
				"00000000" + "00000000", // Function selectors
		},
		{"empty", assets.NewLink(0), 0, models.AnyTime{}, nil,
			"0x0000000000000000000000000000000000000000000000000000000000000000", "0x00000000", "0x00000000",
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000000" + // Payment
				"0000000000000000000000000000000000000000000000000000000000000000" + // Expiration time
				"0000000000000000000000000000000000000000000000000000000000000000" + // EndAt
				requestDigestHex +
				"0000000000000000000000000000000000000000" + // Aggregator address
				"00000000" + "00000000", // Function selectors
		},
		{"oracle address", assets.NewLink(0), 0, models.AnyTime{},
			[]ethkey.EIP55Address{ethkey.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")},
			"0x0000000000000000000000000000000000000000000000000000000000000000", "0x00000000", "0x00000000",
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000000" + // Payment
				"0000000000000000000000000000000000000000000000000000000000000000" + // Expiration time
				"0000000000000000000000000000000000000000000000000000000000000000" + // EndAt
				"000000000000000000000000a0788fc17b1dee36f057c42b6f373a34b014687e" + // Oracle address
				requestDigestHex +
				"0000000000000000000000000000000000000000" + // Aggregator address
				"00000000" + "00000000", // Function selectors
		},
		{"different endAt", assets.NewLink(0), 0, models.NewAnyTime(endAt),
			[]ethkey.EIP55Address{ethkey.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")},
			"0x0000000000000000000000000000000000000000000000000000000000000000", "0x00000000", "0x00000000",
			"0x" +
				"0000000000000000000000000000000000000000000000000000000000000000" + // Payment
				"0000000000000000000000000000000000000000000000000000000000000000" + // Expiration time
				"00000000000000000000000000000000000000000000000000000000459a7465" + // EndAt
				"000000000000000000000000a0788fc17b1dee36f057c42b6f373a34b014687e" + // Oracle address
				requestDigestHex +
				"0000000000000000000000000000000000000000" + // Aggregator address
				"00000000" + "00000000", // Function selectors
		},
		{name: "aggregator info", expiration: 0, endAt: models.NewAnyTime(endAt),
			oracles: []ethkey.EIP55Address{
				ethkey.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e"),
			},
			aggregator:             "0x3141592653589793238462643383279502884197",
			aggInitiateJobSelector: "0x12345678", aggFulfillSelector: "0x87654321",
			want: "0x" +
				"0000000000000000000000000000000000000000000000000000000000000000" + // Payment
				"0000000000000000000000000000000000000000000000000000000000000000" + // Expiration time
				"00000000000000000000000000000000000000000000000000000000459a7465" + // EndAt
				"000000000000000000000000a0788fc17b1dee36f057c42b6f373a34b014687e" + // Oracle address
				requestDigestHex +
				"3141592653589793238462643383279502884197" + // Aggregator address
				"12345678" + "87654321", // Function selectors
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fs := func(s string) models.FunctionSelector {
				return models.BytesToFunctionSelector(hexutil.MustDecode(s))
			}
			enc := models.Encumbrance{
				Payment:                test.payment,
				Expiration:             uint64(test.expiration),
				EndAt:                  test.endAt,
				Oracles:                test.oracles,
				Aggregator:             ethkey.EIP55Address(test.aggregator),
				AggInitiateJobSelector: fs(test.aggInitiateJobSelector),
				AggFulfillSelector:     fs(test.aggFulfillSelector),
			}

			ebytes, err := enc.ABI(digest)
			assert.NoError(t, err)
			assert.Equal(t, test.want, hexutil.Encode(ebytes))
		})
	}
}

func TestServiceAgreementRequest_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		wantDigest  string
		wantPayment *assets.Link
	}{
		{
			"basic",
			`{"encumbrance": {` +
				`"payment":"1",` +
				`"initiators":[{"type":"web"}],` +
				`"tasks":[` +
				`{"type":"HttpGet","params":{"get":"https://bitstamp.net/api/ticker/"}},` +
				`{"type":"JsonParse","params":{"path":["last"]}},` +
				`{"type":"EthBytes32","params":{"type":"ethtx"}}` +
				`],` +
				`"endAt":"2018-06-19T22:17:19Z"}` +
				`}`,
			"0x57bf5be3447b9a3f8491b6538b01f828bcfcaf2d685ea90375ed4ec2943f4865",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreement
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			assert.Equal(t, test.wantPayment, sar.Encumbrance.Payment)
		})
	}
}
