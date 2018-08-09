package models_test

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
)

func TestNewServiceAgreementFromRequest(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		input       string
		wantDigest  string
		wantPayment int64
	}{
		{"basic",
			`{"payment":1,"initiators":[{"type":"web"}],"tasks":[{"type":"httpget","url":"https://bitstamp.net/api/ticker/"},{"type":"jsonparse","path":["last"]},{"type":"ethbytes32"},{"type":"ethtx"}]}`,
			"0xa0911e6d17e4b992e41a12a2dae111382c42328a895983894e1bb6912213d385", 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreementRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			sa, err := models.NewServiceAgreementFromRequest(sar)
			assert.NoError(t, err)
			assert.Equal(t, test.wantDigest, sa.ID)
			assert.Equal(t, big.NewInt(test.wantPayment), sa.Encumbrance.Payment)
		})
	}
}

func TestEncumbrance_ABI(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		payment    *big.Int
		expiration *big.Int
		want       string
	}{
		{"basic", big.NewInt(1), big.NewInt(2), "00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002"},
		{"empty", nil, nil, "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			enc := models.Encumbrance{
				Payment:    test.payment,
				Expiration: test.expiration,
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
		wantPayment int64
	}{
		{"basic",
			`{"payment":1,"initiators":[{"type":"web"}],"tasks":[{"type":"httpget",` +
				`"url":"https://bitstamp.net/api/ticker/"},{"type":"jsonparse",` +
				`"path":["last"]},{"type":"ethbytes32"},{"type":"ethtx"}]}`,
			"0x57bf5be3447b9a3f8491b6538b01f828bcfcaf2d685ea90375ed4ec2943f4865", 1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var sar models.ServiceAgreementRequest
			assert.NoError(t, json.Unmarshal([]byte(test.input), &sar))

			assert.Equal(t, big.NewInt(test.wantPayment), sar.Encumbrance.Payment)
			assert.Equal(t, cltest.NormalizedJSONString([]byte(test.input)), sar.Normalized)
		})
	}
}
