package mercury

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_OffchainConfig(t *testing.T) {
	t.Run("decoding", func(t *testing.T) {
		t.Run("with number type for USD fee", func(t *testing.T) {
			json := `
{
	"expirationWindow": 42,
	"baseUSDFee": 123.456
}
`
			c, err := DecodeOffchainConfig([]byte(json))
			require.NoError(t, err)

			assert.Equal(t, decimal.NewFromFloat32(123.456), c.BaseUSDFee)

			json = `
{
	"expirationWindow": 42,
	"baseUSDFee": 123
}
`
			c, err = DecodeOffchainConfig([]byte(json))
			require.NoError(t, err)

			assert.Equal(t, decimal.NewFromInt32(123), c.BaseUSDFee)

			json = `
{
	"expirationWindow": 42,
	"baseUSDFee": 0.12
}
`
			c, err = DecodeOffchainConfig([]byte(json))
			require.NoError(t, err)

			assert.Equal(t, decimal.NewFromFloat32(0.12), c.BaseUSDFee)
		})
		t.Run("with string type for USD fee", func(t *testing.T) {
			json := `
{
	"expirationWindow": 42,
	"baseUSDFee": "123.456"
}
`
			c, err := DecodeOffchainConfig([]byte(json))
			require.NoError(t, err)

			assert.Equal(t, decimal.NewFromFloat32(123.456), c.BaseUSDFee)

			json = `
{
	"expirationWindow": 42,
	"baseUSDFee": "123"
}
`
			c, err = DecodeOffchainConfig([]byte(json))
			require.NoError(t, err)

			assert.Equal(t, decimal.NewFromInt32(123), c.BaseUSDFee)

			json = `
{
	"expirationWindow": 42,
	"baseUSDFee": "0.12"
}
`
			c, err = DecodeOffchainConfig([]byte(json))
			require.NoError(t, err)

			assert.Equal(t, decimal.NewFromFloat32(0.12), c.BaseUSDFee)
		})
	})
	t.Run("serialize/deserialize", func(t *testing.T) {
		c := OffchainConfig{32, decimal.NewFromFloat32(1.23)}

		serialized, err := c.Encode()
		require.NoError(t, err)

		deserialized, err := DecodeOffchainConfig(serialized)
		require.NoError(t, err)

		assert.Equal(t, c, deserialized)
	})
}
