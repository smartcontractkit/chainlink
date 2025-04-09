package internal

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRaw_nonprintable(t *testing.T) {
	bytes, err := hex.DecodeString("f0d07ab448018b2754475f9a3b580218b0675a1456aad96ad607c7bbd7d9237b")
	require.NoError(t, err)
	r := NewRaw(bytes)

	exp := "<Raw Private Key>"

	assert.Equal(t, exp, fmt.Sprint(r))

	assert.Equal(t, exp, fmt.Sprintf("%v", r))

	assert.Equal(t, exp, fmt.Sprintf("%#v", r))

	assert.Equal(t, exp, fmt.Sprintf("%s", r)) //nolint:gosimple // S1025 deliberately testing formatting verbs

	got, err := json.Marshal(r) //nolint:staticcheck // SA9005 deliberately testing marshalling

	if assert.NoError(t, err) {
		assert.Equal(t, `{}`, string(got))
	}
}
