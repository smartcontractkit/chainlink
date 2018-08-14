package models_test

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/assets"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestNewServiceAgreementFromRequest(t *testing.T) {
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

			sa, err := models.NewServiceAgreementFromRequest(sar)
			assert.NoError(t, err)
			assert.Equal(t, test.wantDigest, sa.ID)
			assert.Equal(t, test.wantPayment, sa.Encumbrance.Payment)
			assert.NotEqual(t, models.Time{}, sa.CreatedAt)
		})
	}
}

func TestEncumbrance_ABI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		payment    *assets.Link
		expiration int
		want       string
	}{
		{"basic", assets.NewLink(1), 2, "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002"},
		{"basic", assets.NewLink(3735928559), 2, "00000000000000000000000000000000000000000000000000000000deadbeef0000000000000000000000000000000000000000000000000000000000000002"},
		{"empty", nil, 0, "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enc := models.Encumbrance{
				Payment:    test.payment,
				Expiration: uint64(test.expiration),
			}

			assert.Equal(t, test.want, enc.ABI())
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

			assert.Equal(t, test.wantPayment, sar.Encumbrance.Payment)
			assert.Equal(t, cltest.NormalizedJSON([]byte(test.input)), sar.NormalizedBody)
		})
	}
}
