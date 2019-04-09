package models_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/smartcontractkit/chainlink/core/store/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/tools/cltest"
	"github.com/stretchr/testify/assert"
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
				`{"type":"ethbytes32"},{"type":"ethtx"}]}`,
			"0x4080e87b11b47454e49e19de88af26d9c80628cff774780f4fb4260c12a7c8de",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.JobSpecRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			us, err := models.NewUnsignedServiceAgreementFromRequest(strings.NewReader(test.input))
			assert.NoError(t, err)
			assert.Equal(t, test.wantDigest, us.ID.String())
			assert.Equal(t, test.wantPayment, us.Encumbrance.Payment)
			assert.Equal(t, cltest.NormalizedJSON([]byte(test.input)), us.RequestBody)
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
				`{"type":"ethbytes32"},{"type":"ethtx"}]}`,
			"0x4080e87b11b47454e49e19de88af26d9c80628cff774780f4fb4260c12a7c8de",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.JobSpecRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			us, err := models.NewUnsignedServiceAgreementFromRequest(strings.NewReader(test.input))
			assert.NoError(t, err)

			sa, err := models.BuildServiceAgreement(us, cltest.MockSigner{})
			assert.NoError(t, err)
			assert.Equal(t, test.wantDigest, sa.ID)
			assert.Equal(t, test.wantPayment, sa.Encumbrance.Payment)
			assert.Equal(t, cltest.NormalizedJSON([]byte(test.input)), sa.RequestBody)
			assert.NotEqual(t, models.AnyTime{}, sa.CreatedAt)
			assert.NotEqual(t, "", sa.Signature.String())
		})
	}
}

func TestEncumbrance_ABI(t *testing.T) {
	t.Parallel()
	endAt, _ := time.Parse("2006-01-02T15:04:05.000Z", "2007-01-02T15:04:05.000Z")

	tests := []struct {
		name       string
		payment    *assets.Link
		expiration int
		endAt      models.AnyTime
		oracles    []models.EIP55Address
		want       string
	}{
		{"basic", assets.NewLink(1), 2, models.AnyTime{}, []models.EIP55Address{},
			"0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002886e0900"},
		{"basic dead beef payment", assets.NewLink(3735928559), 2, models.AnyTime{}, []models.EIP55Address{},
			"0x00000000000000000000000000000000000000000000000000000000deadbeef0000000000000000000000000000000000000000000000000000000000000002886e0900"},
		{"empty", nil, 0, models.AnyTime{}, []models.EIP55Address{}, "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000886e0900"},
		{"oracle address", nil, 0, models.AnyTime{},
			[]models.EIP55Address{models.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")},
			"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000886e0900000000000000000000000000a0788fc17b1dee36f057c42b6f373a34b014687e"},
		{"oracle address", nil, 0, models.AnyTime{},
			[]models.EIP55Address{models.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")},
			"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000886e0900000000000000000000000000a0788fc17b1dee36f057c42b6f373a34b014687e"},
		{"different endAt", nil, 0, models.NewAnyTime(endAt),
			[]models.EIP55Address{models.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")},
			"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000459a7465000000000000000000000000a0788fc17b1dee36f057c42b6f373a34b014687e"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enc := models.Encumbrance{
				Payment:    test.payment,
				Expiration: uint64(test.expiration),
				EndAt:      test.endAt,
				Oracles:    test.oracles,
			}

			ebytes, err := enc.ABI()
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
