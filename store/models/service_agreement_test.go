package models_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
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
			"0xc7106c5877b5bd321e5aac3842cd6ae68faf21e7e6ee45556b13f7b386104381",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreementRequest
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
			"0xc7106c5877b5bd321e5aac3842cd6ae68faf21e7e6ee45556b13f7b386104381",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreementRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			us, err := models.NewUnsignedServiceAgreementFromRequest(strings.NewReader(test.input))
			assert.NoError(t, err)

			sa, err := models.BuildServiceAgreement(us, cltest.MockSigner{})
			assert.NoError(t, err)
			assert.Equal(t, test.wantDigest, sa.ID)
			assert.Equal(t, test.wantPayment, sa.Encumbrance.Payment)
			assert.Equal(t, cltest.NormalizedJSON([]byte(test.input)), sa.RequestBody)
			assert.NotEqual(t, models.Time{}, sa.CreatedAt)
			assert.NotEqual(t, "", sa.Signature.String())
		})
	}
}

func TestEncumbrance_ABI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		payment    *assets.Link
		expiration int
		oracles    []models.EIP55Address
		want       string
	}{
		{"basic", assets.NewLink(1), 2, []models.EIP55Address{}, "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002"},
		{"basic", assets.NewLink(3735928559), 2, []models.EIP55Address{}, "0x00000000000000000000000000000000000000000000000000000000deadbeef0000000000000000000000000000000000000000000000000000000000000002"},
		{"empty", nil, 0, []models.EIP55Address{}, "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
		{"oracle address", nil, 0, []models.EIP55Address{models.EIP55Address("0xa0788FC17B1dEe36f057c42B6F373A34B014687e")}, "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a0788fc17b1dee36f057c42b6f373a34b014687e"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enc := models.Encumbrance{
				Payment:    test.payment,
				Expiration: uint64(test.expiration),
				Oracles:    test.oracles,
			}

			ebytes, err := enc.ABI()
			assert.NoError(t, err)
			assert.Equal(t, test.want, common.ToHex(ebytes))
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
			`{"payment":"1","initiators":[{"type":"web"}],"tasks":[{"type":"httpget",` +
				`"url":"https://bitstamp.net/api/ticker/"},{"type":"jsonparse",` +
				`"path":["last"]},{"type":"ethbytes32"},{"type":"ethtx"}]}`,
			"0x57bf5be3447b9a3f8491b6538b01f828bcfcaf2d685ea90375ed4ec2943f4865",
			assets.NewLink(1),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreementRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			assert.Equal(t, test.wantPayment, sar.Payment)
		})
	}
}
