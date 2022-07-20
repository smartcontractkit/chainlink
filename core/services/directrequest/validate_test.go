package directrequest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatedDirectRequestSpec(t *testing.T) {
	t.Parallel()

	toml := `
type                = "directrequest"
schemaVersion       = 1
name                = "example eth request event spec"
contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID       = "A5AC14E8-7629-4726-B1F1-1AE053FC829E"
observationSource   = """
    ds1          [type=http method=GET url="example.com" allowunrestrictednetworkaccess="true"];
    ds1_parse    [type=jsonparse path="USD"];
    ds1_multiply [type=multiply times=100];
    ds1 -> ds1_parse -> ds1_multiply;
"""
`

	s, err := ValidatedDirectRequestSpec(toml)
	require.NoError(t, err)

	assert.Equal(t, int32(0), s.ID)
	assert.Equal(t, "0x613a38AC1659769640aaE063C651F48E0250454C", s.DirectRequestSpec.ContractAddress.Hex())
	assert.Equal(t, "0x6135616331346538373632393437323662316631316165303533666338323965", s.ExternalIDEncodeStringToTopic().String())
	assert.Equal(t, "0xa5ac14e876294726b1f11ae053fc829e00000000000000000000000000000000", s.ExternalIDEncodeBytesToTopic().String())
	assert.NotZero(t, s.ExternalJobID.Bytes()[:])
	assert.Equal(t, time.Time{}, s.DirectRequestSpec.CreatedAt)
	assert.Equal(t, time.Time{}, s.DirectRequestSpec.UpdatedAt)
}

func TestValidatedDirectRequestSpec_MinIncomingConfirmations(t *testing.T) {
	t.Parallel()

	t.Run("no minIncomingConfirmations specified", func(t *testing.T) {
		t.Parallel()

		toml := `
		type                = "directrequest"
		schemaVersion       = 1
		name                = "example eth request event spec"
		observationSource   = """
		"""
		`

		s, err := ValidatedDirectRequestSpec(toml)
		require.NoError(t, err)

		assert.False(t, s.DirectRequestSpec.MinIncomingConfirmations.Valid)
	})

	t.Run("minIncomingConfirmations set to 100", func(t *testing.T) {
		t.Parallel()

		toml := `
		type                = "directrequest"
		schemaVersion       = 1
		name                = "example eth request event spec"
		minIncomingConfirmations = 100
		observationSource   = """
		"""
		`

		s, err := ValidatedDirectRequestSpec(toml)
		require.NoError(t, err)

		assert.True(t, s.DirectRequestSpec.MinIncomingConfirmations.Valid)
		assert.Equal(t, uint32(100), s.DirectRequestSpec.MinIncomingConfirmations.Uint32)
	})
}
