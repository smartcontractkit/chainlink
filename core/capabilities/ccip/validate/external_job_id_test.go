package validate_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"
	"github.com/stretchr/testify/require"
)

func TestExternalJobID_Generation(t *testing.T) {
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
		wantErr  bool
	}{
		{
			name: "valid basic spec args",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "different capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip-v2",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "different version",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v2.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "different p2p key ID",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "different-key-id",
			},
			wantErr: false,
		},
		{
			name: "empty capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false, // Should still generate valid UUID
		},
		{
			name: "empty version",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false, // Should still generate valid UUID
		},
		{
			name: "empty p2p key ID",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "",
			},
			wantErr: false, // Should still generate valid UUID
		},
		{
			name: "all empty fields",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "",
				CapabilityVersion:      "",
				P2PKeyID:               "",
			},
			wantErr: false, // Should still generate valid UUID
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate.ExternalJobID(tt.specArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, got)

				// Validate that the result is a valid UUID
				_, err := uuid.Parse(got)
				require.NoError(t, err, "Generated ID should be a valid UUID")
			}
		})
	}
}

func TestExternalJobID_Deterministic(t *testing.T) {
	// Test that the same input always produces the same output
	specArgs := validate.SpecArgs{
		CapabilityLabelledName: "ccip",
		CapabilityVersion:      "v1.0.0",
		P2PKeyID:               "test-p2p-key-id",
	}

	// Generate ID multiple times
	id1, err := validate.ExternalJobID(specArgs)
	require.NoError(t, err)

	id2, err := validate.ExternalJobID(specArgs)
	require.NoError(t, err)

	id3, err := validate.ExternalJobID(specArgs)
	require.NoError(t, err)

	// All IDs should be identical
	require.Equal(t, id1, id2, "Same input should produce same ID")
	require.Equal(t, id2, id3, "Same input should produce same ID")
	require.Equal(t, id1, id3, "Same input should produce same ID")
}

func TestExternalJobID_Uniqueness(t *testing.T) {
	// Test that different inputs produce different IDs
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
	}{
		{
			name: "base case",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
		},
		{
			name: "different capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip-v2",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
		},
		{
			name: "different version",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v2.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
		},
		{
			name: "different p2p key ID",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "different-key-id",
			},
		},
	}

	generatedIDs := make(map[string]bool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := validate.ExternalJobID(tt.specArgs)
			require.NoError(t, err)
			require.NotEmpty(t, id)

			// Check that this ID hasn't been generated before
			require.False(t, generatedIDs[id], "ID %s should be unique, but was already generated", id)
			generatedIDs[id] = true
		})
	}
}

func TestExternalJobID_CaseSensitivity(t *testing.T) {
	// Test that case differences produce different IDs
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
	}{
		{
			name: "lowercase capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
		},
		{
			name: "uppercase capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "CCIP",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
		},
		{
			name: "mixed case capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "CcIp",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
		},
	}

	generatedIDs := make([]string, 0, len(tests))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := validate.ExternalJobID(tt.specArgs)
			require.NoError(t, err)
			require.NotEmpty(t, id)
			generatedIDs = append(generatedIDs, id)
		})
	}

	// Verify all IDs are different
	for i := 0; i < len(generatedIDs); i++ {
		for j := i + 1; j < len(generatedIDs); j++ {
			require.NotEqual(t, generatedIDs[i], generatedIDs[j],
				"IDs should be different for case-sensitive inputs: %s vs %s", generatedIDs[i], generatedIDs[j])
		}
	}
}

func TestExternalJobID_SpecialCharacters(t *testing.T) {
	// Test handling of special characters
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
		wantErr  bool
	}{
		{
			name: "special characters in capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip-@#$%^&*()",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "unicode characters",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip-测试-🚀",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "whitespace characters",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip with spaces",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test p2p key id",
			},
			wantErr: false,
		},
		{
			name: "newline and tab characters",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip\n\t",
				CapabilityVersion:      "v1.0.0\n",
				P2PKeyID:               "test\tp2p\nkey\rid",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := validate.ExternalJobID(tt.specArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, id)

				// Validate UUID format
				_, err = uuid.Parse(id)
				require.NoError(t, err)
			}
		})
	}
}

func TestExternalJobID_LargeInputs(t *testing.T) {
	// Test handling of very large inputs
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
		wantErr  bool
	}{
		{
			name: "large capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: strings.Repeat("ccip", 1000),
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "large version",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      strings.Repeat("v1.0.0", 500),
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "large p2p key ID",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               strings.Repeat("test-p2p-key-id", 500),
			},
			wantErr: false,
		},
		{
			name: "all large fields",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: strings.Repeat("ccip", 500),
				CapabilityVersion:      strings.Repeat("v1.0.0", 500),
				P2PKeyID:               strings.Repeat("test-p2p-key-id", 500),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := validate.ExternalJobID(tt.specArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, id)

				// Validate UUID format
				_, err = uuid.Parse(id)
				require.NoError(t, err)
			}
		})
	}
}

func TestExternalJobID_UUIDVersion4Properties(t *testing.T) {
	// Test that generated UUIDs have correct version 4 properties
	specArgs := validate.SpecArgs{
		CapabilityLabelledName: "ccip",
		CapabilityVersion:      "v1.0.0",
		P2PKeyID:               "test-p2p-key-id",
	}

	id, err := validate.ExternalJobID(specArgs)
	require.NoError(t, err)

	// Parse the UUID
	parsedUUID, err := uuid.Parse(id)
	require.NoError(t, err)

	// Check that it's version 4
	require.Equal(t, uuid.Version(4), parsedUUID.Version(), "UUID should be version 4")

	// Check that variant is correct (should be 10 in binary, which is 2 in decimal)
	require.Equal(t, uuid.RFC4122, parsedUUID.Variant(), "UUID should have RFC4122 variant")

	// Verify the UUID string format (should be 36 characters with hyphens)
	require.Len(t, id, 36, "UUID string should be 36 characters long")
	require.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, id, "UUID should match version 4 format")
}