package tracking

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOverlyLongMetadataTruncation(t *testing.T) {
	metadata := map[string]any{
		"to_truncate": "abcde" + strings.Repeat("1234567890", 110), // overly long
		"note":        "keep this short",
		"debug":       "0123456789... etc",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{"to_truncate", "debug"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["debug"], truncated["debug"], "debug should not be truncated")
	require.Equal(t, metadata["note"], truncated["note"], "note should not be truncated")
	require.NotEqual(t, metadata["to_truncate"], truncated["to_truncate"], "to_truncate should be truncated")
}

func TestOnlyOneOveryLongMetadataTruncationNoPriority(t *testing.T) {
	metadata := map[string]any{
		"to_truncate":     "abcde" + strings.Repeat("1234567890", 110), // overly long
		"do_not_truncate": "abcde" + strings.Repeat("1234567890", 50),
		"note":            "keep this short",
		"debug":           "0123456789... etc",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["debug"], truncated["debug"], "debug should not be truncated")
	require.Equal(t, metadata["note"], truncated["note"], "note should not be truncated")
	require.NotEqual(t, metadata["to_truncate"], truncated["to_truncate"], "to_truncate should be truncated")
	require.Equal(t, metadata["do_not_truncate"], truncated["do_not_truncate"], "do_not_truncate should not be truncated")
}

func TestOnlyOneOveryLongMetadataTruncationNonExistingPriority(t *testing.T) {
	metadata := map[string]any{
		"to_truncate":     "abcde" + strings.Repeat("1234567890", 110), // overly long
		"do_not_truncate": "abcde" + strings.Repeat("1234567890", 50),
		"note":            "keep this short",
		"debug":           "0123456789... etc",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{"i_dont_exist"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["debug"], truncated["debug"], "debug should not be truncated")
	require.Equal(t, metadata["note"], truncated["note"], "note should not be truncated")
	require.NotEqual(t, metadata["to_truncate"], truncated["to_truncate"], "to_truncate should be truncated")
	require.Equal(t, metadata["do_not_truncate"], truncated["do_not_truncate"], "do_not_truncate should not be truncated")
}

func TestTogetherOveryLongMetadataTruncationNoPriority(t *testing.T) {
	metadata := map[string]any{
		"to_truncate":     "abcde" + strings.Repeat("1234567890", 60), // together both fields are overly long and one of them has to be truncated
		"do_not_truncate": "abcde" + strings.Repeat("1234567890", 60),
		"note":            "keep this short",
		"debug":           "0123456789... etc",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["debug"], truncated["debug"], "debug should not be truncated")
	require.Equal(t, metadata["note"], truncated["note"], "note should not be truncated")

	truncatedFields := []string{}
	if metadata["to_truncate"] != truncated["to_truncate"] {
		truncatedFields = append(truncatedFields, "to_truncate")
	}
	if metadata["do_not_truncate"] != truncated["do_not_truncate"] {
		truncatedFields = append(truncatedFields, "do_not_truncate")
	}

	require.Len(t, truncatedFields, 1, "only one field should be truncated")
}

func TestNoTruncation(t *testing.T) {
	// exactly 1024 bytes
	metadata := map[string]any{
		"all_good": "abcde" + strings.Repeat("x", 951),
		"note":     "keep this short",
		"debug":    "0123456789... etc",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{"all_good", "debug"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["debug"], truncated["debug"], "debug should not be truncated")
	require.Equal(t, metadata["note"], truncated["note"], "note should not be truncated")
	require.Equal(t, metadata["all_good"], truncated["all_good"], "all_good should not be truncated")
}

func TestNoTruncationAndNoPriority(t *testing.T) {
	// exactly 1024 bytes
	metadata := map[string]any{
		"all_good": "abcde" + strings.Repeat("x", 951),
		"note":     "keep this short",
		"debug":    "0123456789... etc",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["debug"], truncated["debug"], "debug should not be truncated")
	require.Equal(t, metadata["note"], truncated["note"], "note should not be truncated")
	require.Equal(t, metadata["all_good"], truncated["all_good"], "all_good should not be truncated")
}

func TestEmptyMetadata(t *testing.T) {
	metadata := map[string]any{}

	truncated, truncateErr := truncateMetadata(metadata, []string{"nonexistent"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "empty metadata should be valid")
	require.Len(t, metadata, len(truncated), "empty metadata should remain empty")
}

func TestNonStringValuesOnly(t *testing.T) {
	metadata := map[string]any{
		"number":  12345,
		"boolean": true,
		"array":   []int{1, 2, 3, 4, 5},
		"object":  map[string]int{"key": 123},
		"null":    nil,
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{"number", "boolean"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be valid")
	require.Equal(t, metadata["number"], truncated["number"], "number should not be truncated")
	require.Equal(t, metadata["boolean"], truncated["boolean"], "boolean should not be truncated")
	require.Equal(t, metadata["array"], truncated["array"], "array should not be truncated")
	require.Equal(t, metadata["object"], truncated["object"], "object should not be truncated")
	require.Equal(t, metadata["null"], truncated["null"], "null should not be truncated")
}

func TestMixedStringAndNonStringValues(t *testing.T) {
	metadata := map[string]any{
		"long_string": "abcde" + strings.Repeat("1234567890", 100), // overly long
		"number":      12345,
		"boolean":     true,
		"short_str":   "short",
		"array":       []int{1, 2, 3, 4, 5},
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["number"], truncated["number"], "number should not be truncated")
	require.Equal(t, metadata["boolean"], truncated["boolean"], "boolean should not be truncated")
	require.Equal(t, metadata["short_str"], truncated["short_str"], "short_str should not be truncated")
	require.Equal(t, metadata["array"], truncated["array"], "array should not be truncated")
	require.NotEqual(t, metadata["long_string"], truncated["long_string"], "long_string should be truncated")
	require.Contains(t, truncated["long_string"].(string), "(... truncated)", "truncated string should have suffix")
}

func TestMixedStringAndNonStringValuesNoTruncation(t *testing.T) {
	metadata := map[string]any{
		"not_long_string": "abcde" + strings.Repeat("1234567890", 10), // overly long
		"number":          12345,
		"boolean":         true,
		"short_str":       "short",
		"array":           []int{1, 2, 3, 4, 5},
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["number"], truncated["number"], "number should not be truncated")
	require.Equal(t, metadata["boolean"], truncated["boolean"], "boolean should not be truncated")
	require.Equal(t, metadata["short_str"], truncated["short_str"], "short_str should not be truncated")
	require.Equal(t, metadata["array"], truncated["array"], "array should not be truncated")
	require.Equal(t, metadata["long_string"], truncated["long_string"], "long_string should be truncated")
	require.Equal(t, metadata["not_long_string"], truncated["not_long_string"], "truncated string should have suffix")
}

func TestPriorityFieldsNeedTruncation(t *testing.T) {
	metadata := map[string]any{
		"priority1": "abcde" + strings.Repeat("1234567890", 50), // long priority field
		"priority2": "abcde" + strings.Repeat("abcdefghij", 50), // another long priority field
		"regular":   "abcde" + strings.Repeat("xyz", 50),        // regular long field
		"short":     "keep this short",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{"priority1", "priority2"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["short"], truncated["short"], "short should not be truncated")
	// Priority fields should be considered first, but algorithm only truncates as much as needed
	require.NotEqual(t, metadata["priority1"], truncated["priority1"], "priority1 should be truncated")
	// priority2 might not need truncation if priority1 truncation was enough
	require.Equal(t, metadata["regular"], truncated["regular"], "regular should not be truncated since priority fields are truncated first")
}

func TestEmptyStrings(t *testing.T) {
	metadata := map[string]any{
		"empty1": "",
		"empty2": "",
		"filled": "abcde" + strings.Repeat("1234567890", 100), // overly long
		"short":  "short",
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{"empty1", "filled"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["empty1"], truncated["empty1"], "empty1 should not be truncated")
	require.Equal(t, metadata["empty2"], truncated["empty2"], "empty2 should not be truncated")
	require.Equal(t, metadata["short"], truncated["short"], "short should not be truncated")
	require.NotEqual(t, metadata["filled"], truncated["filled"], "filled should be truncated")
}

func TestExtremelyLargeDataThatCannotFit(t *testing.T) {
	// Create metadata so large that even after truncation it won't fit
	metadata := map[string]any{
		"huge1": "abcde" + strings.Repeat("1234567890", 200),
		"huge2": "abcde" + strings.Repeat("abcdefghij", 200),
		"huge3": "abcde" + strings.Repeat("zyxwvutsrq", 200),
		"huge4": "abcde" + strings.Repeat("0987654321", 200),
		"huge5": "abcde" + strings.Repeat("qwertyuiop", 200),
	}

	truncated, truncateErr := truncateMetadata(metadata, []string{})

	// The function should return an error when it can't fit metadata within 1024 bytes
	if truncateErr != nil {
		require.Error(t, truncateErr)
		require.Contains(t, truncateErr.Error(), "unable to fit metadata within 1024 bytes")
	} else {
		// If no error, it should at least be within the limit
		b, err := json.Marshal(truncated)
		require.NoError(t, err)
		require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	}
}

func TestPriorityFieldsDoNotExist(t *testing.T) {
	metadata := map[string]any{
		"existing": "abcde" + strings.Repeat("1234567890", 100), // overly long
		"short":    "short",
	}

	// Priority fields that don't exist in metadata
	truncated, truncateErr := truncateMetadata(metadata, []string{"nonexistent1", "nonexistent2", "existing"})
	require.NoError(t, truncateErr)

	b, err := json.Marshal(truncated)
	require.NoError(t, err)

	require.LessOrEqual(t, len(b), 1024, "metadata should be truncated to less than 1024 bytes")
	require.Equal(t, metadata["short"], truncated["short"], "short should not be truncated")
	require.NotEqual(t, metadata["existing"], truncated["existing"], "existing should be truncated")
}
