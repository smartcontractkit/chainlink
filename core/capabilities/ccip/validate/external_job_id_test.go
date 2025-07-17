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
	require.Equal(t, id1, id2)
	require.Equal(t, id2, id3)
	require.Equal(t, id1, id3)

	// Validate UUID format
	_, err = uuid.Parse(id1)
	require.NoError(t, err)
}

func TestExternalJobID_Uniqueness(t *testing.T) {
	// Test that different inputs produce different outputs
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
			name: "different name",
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
			name: "different key ID",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "different-key-id",
			},
		},
		{
			name: "case sensitivity test",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "CCIP",
				CapabilityVersion:      "V1.0.0",
				P2PKeyID:               "TEST-P2P-KEY-ID",
			},
		},
	}

	generatedIDs := make(map[string]bool)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := validate.ExternalJobID(tt.specArgs)
			require.NoError(t, err)
			require.NotEmpty(t, id)

			// Validate UUID format
			_, err = uuid.Parse(id)
			require.NoError(t, err)

			// Check uniqueness
			require.False(t, generatedIDs[id], "Generated ID %s should be unique", id)
			generatedIDs[id] = true
		})
	}
}

func TestExternalJobID_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
		wantErr  bool
	}{
		{
			name: "special characters in name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip-v2@test",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "special characters in version",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0-beta.1+build.123",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "special characters in key ID",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id_123@domain.com",
			},
			wantErr: false,
		},
		{
			name: "unicode characters",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip-测试",
				CapabilityVersion:      "v1.0.0-测试",
				P2PKeyID:               "test-p2p-key-id-测试",
			},
			wantErr: false,
		},
		{
			name: "whitespace characters",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: " ccip ",
				CapabilityVersion:      " v1.0.0 ",
				P2PKeyID:               " test-p2p-key-id ",
			},
			wantErr: false,
		},
		{
			name: "newline characters",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip\n",
				CapabilityVersion:      "v1.0.0\n",
				P2PKeyID:               "test-p2p-key-id\n",
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
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
		wantErr  bool
	}{
		{
			name: "very long capability name",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: strings.Repeat("ccip", 1000),
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "very long version",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      strings.Repeat("v1.0.0", 1000),
				P2PKeyID:               "test-p2p-key-id",
			},
			wantErr: false,
		},
		{
			name: "very long key ID",
			specArgs: validate.SpecArgs{
				CapabilityLabelledName: "ccip",
				CapabilityVersion:      "v1.0.0",
				P2PKeyID:               strings.Repeat("test-p2p-key-id", 1000),
			},
			wantErr: false,
		},
		{
			name: "all fields very long",
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