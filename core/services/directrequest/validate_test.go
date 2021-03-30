package directrequest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatedDirectRequestSpec(t *testing.T) {
	toml := `
type                = "directrequest"
schemaVersion       = 1
name                = "example eth request event spec"
contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
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
	assert.NotZero(t, s.DirectRequestSpec.OnChainJobSpecID[:])
	assert.Equal(t, time.Time{}, s.DirectRequestSpec.CreatedAt)
	assert.Equal(t, time.Time{}, s.DirectRequestSpec.UpdatedAt)
}
