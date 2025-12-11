package stringutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
)

func TestStringUtils_ToInt64(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "unit")

	want := int64(12)

	got, err := ToInt64("12")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestStringUtils_FromInt64(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "unit")

	want := "12"

	got := FromInt64(int64(12))

	assert.Equal(t, want, got)
}

func TestStringUtils_ToInt32(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "unit")

	want := int32(32)

	got, err := ToInt32("32")
	require.NoError(t, err)

	assert.Equal(t, want, got)
}

func TestStringUtils_FromInt32(t *testing.T) {
	t.Parallel()
	tests.BelongsToCISuite(t, "unit")

	want := "32"

	got := FromInt32(int32(32))

	assert.Equal(t, want, got)
}
