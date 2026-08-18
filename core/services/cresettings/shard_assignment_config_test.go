package cresettings

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseShardAssignmentConfig_Empty(t *testing.T) {
	t.Parallel()
	cfg, err := ParseShardAssignmentConfig("")
	require.NoError(t, err)
	assert.Empty(t, cfg.StaticDefaultAssignment)
	assert.Empty(t, cfg.PerOwnerAssignment)
	assert.Empty(t, cfg.PerOrgAssignment)
	assert.Empty(t, cfg.HashedOwnerAssignment)
	assert.False(t, cfg.HashedDefaultAssignment)
}

func TestParseShardAssignmentConfig_Valid(t *testing.T) {
	t.Parallel()
	raw := `
static_default_assignment = [0, 1]
disabled_shards = []
hashed_default_assignment = false
hashed_owner_assignment = ["0x3c44cdddb6a900fa8b69783f73a91f60e8c1c348"]

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1, 2]
  "0x70997970c51812dc3a010c7d01b50e0d17dc79c8" = [0, 2]
`
	cfg, err := ParseShardAssignmentConfig(raw)
	require.NoError(t, err)
	assert.Equal(t, []uint32{0, 1}, cfg.StaticDefaultAssignment)
	assert.Len(t, cfg.PerOwnerAssignment, 2)
	assert.Equal(t, []uint32{1, 2}, cfg.PerOwnerAssignment["f39fd6e51aad88f6f4ce6ab8827279cfffb92266"])
	assert.Equal(t, []uint32{0, 2}, cfg.PerOwnerAssignment["70997970c51812dc3a010c7d01b50e0d17dc79c8"])
	assert.True(t, cfg.HashedOwnerAssignment["3c44cdddb6a900fa8b69783f73a91f60e8c1c348"])
	assert.False(t, cfg.HashedDefaultAssignment)
}

func TestParseShardAssignmentConfig_OwnerNormalization(t *testing.T) {
	t.Parallel()
	raw := `
[per_owner_assignment]
  "0xF39FD6E51AAD88F6F4CE6AB8827279CFFFB92266" = [1]
  "0x70997970C51812DC3A010C7D01B50E0D17DC79C8" = [0]
`
	cfg, err := ParseShardAssignmentConfig(raw)
	require.NoError(t, err)
	assert.Equal(t, []uint32{1}, cfg.PerOwnerAssignment["f39fd6e51aad88f6f4ce6ab8827279cfffb92266"])
	assert.Equal(t, []uint32{0}, cfg.PerOwnerAssignment["70997970c51812dc3a010c7d01b50e0d17dc79c8"])
}

func TestParseShardAssignmentConfig_PerOrgAssignment(t *testing.T) {
	t.Parallel()
	raw := `
static_default_assignment = [0]

[per_org_assignment]
  org_4wfzBbaifbL32SMN = [1, 2]
  org_mE8piSzea082U52J = [0, 2]
`
	cfg, err := ParseShardAssignmentConfig(raw)
	require.NoError(t, err)
	assert.Len(t, cfg.PerOrgAssignment, 2)
	assert.Equal(t, []uint32{1, 2}, cfg.PerOrgAssignment["org_4wfzBbaifbL32SMN"])
	assert.Equal(t, []uint32{0, 2}, cfg.PerOrgAssignment["org_mE8piSzea082U52J"])
}

func TestParseShardAssignmentConfig_PerOrgAndPerOwner(t *testing.T) {
	t.Parallel()
	raw := `
static_default_assignment = [0]

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]

[per_org_assignment]
  org_4wfzBbaifbL32SMN = [2]
`
	cfg, err := ParseShardAssignmentConfig(raw)
	require.NoError(t, err)
	assert.Len(t, cfg.PerOwnerAssignment, 1)
	assert.Len(t, cfg.PerOrgAssignment, 1)
	assert.Equal(t, []uint32{1}, cfg.PerOwnerAssignment["f39fd6e51aad88f6f4ce6ab8827279cfffb92266"])
	assert.Equal(t, []uint32{2}, cfg.PerOrgAssignment["org_4wfzBbaifbL32SMN"])
}

func TestParseShardAssignmentConfig_EmptyPerOrgShardList(t *testing.T) {
	t.Parallel()
	raw := `
[per_org_assignment]
  org_4wfzBbaifbL32SMN = []
`
	_, err := ParseShardAssignmentConfig(raw)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestParseShardAssignmentConfig_HashedDefaultTrue(t *testing.T) {
	t.Parallel()
	raw := `
hashed_default_assignment = true
static_default_assignment = [0]
`
	cfg, err := ParseShardAssignmentConfig(raw)
	require.NoError(t, err)
	assert.True(t, cfg.HashedDefaultAssignment)
}

func TestParseShardAssignmentConfig_InvalidOwnerAddress(t *testing.T) {
	t.Parallel()
	raw := `
[per_owner_assignment]
  "0xinvalid" = [1]
`
	_, err := ParseShardAssignmentConfig(raw)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid owner address")
}

func TestParseShardAssignmentConfig_DuplicateOwnerInBothMaps(t *testing.T) {
	t.Parallel()
	raw := `
[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]
`
	_, err := ParseShardAssignmentConfig(raw)
	require.NoError(t, err)

	raw = `
hashed_owner_assignment = ["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"]

[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = [1]
`
	_, err = ParseShardAssignmentConfig(raw)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "appears in both")
}

func TestParseShardAssignmentConfig_NegativeShardID(t *testing.T) {
	t.Parallel()
	raw := `
static_default_assignment = [-1]
`
	_, err := ParseShardAssignmentConfig(raw)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must be in range")
}

func TestParseShardAssignmentConfig_EmptyPerOwnerShardList(t *testing.T) {
	t.Parallel()
	raw := `
[per_owner_assignment]
  "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266" = []
`
	_, err := ParseShardAssignmentConfig(raw)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "must not be empty")
}

func TestNormalizeOwnerHex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input    string
		expected string
		hasError bool
	}{
		{"0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "f39fd6e51aad88f6f4ce6ab8827279cfffb92266", false},
		{"0XF39FD6E51AAD88F6F4CE6AB8827279CFFFB92266", "f39fd6e51aad88f6f4ce6ab8827279cfffb92266", false},
		{"F39FD6E51AAD88F6F4CE6AB8827279CFFFB92266", "f39fd6e51aad88f6f4ce6ab8827279cfffb92266", false},
		{"0xshort", "", true},
		{"0xgggggggggggggggggggggggggggggggggggggggg", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			result, err := normalizeOwnerHex(tt.input)
			if tt.hasError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
