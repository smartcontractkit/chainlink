package operations

import (
	"encoding/json"
	"math"
	"sync"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type argument struct {
	def   Definition
	input any
}

func Test_constructUniqueHashFrom(t *testing.T) {
	t.Parallel()

	type Input struct {
		A int
		B int
	}

	definition := Definition{
		ID:          "plus1",
		Version:     semver.MustParse("1.0.0"),
		Description: "plus 1",
	}
	altDefinition := Definition{
		ID:          "plus2",
		Version:     semver.MustParse("1.0.0"),
		Description: "plus 2",
	}

	tests := []struct {
		name           string
		left           argument
		right          argument
		wantSame       bool
		wantErr        string
		skipCacheCheck bool
	}{
		{
			name: "Same def and input should have same hash (struct)",
			left: argument{
				def:   definition,
				input: Input{A: 1, B: 2},
			},
			right: argument{
				def:   definition,
				input: Input{A: 1, B: 2},
			},
			wantSame: true,
		},
		{
			name: "Same def but different input should have different hash",
			left: argument{
				def:   definition,
				input: Input{A: 1, B: 2},
			},
			right: argument{
				def:   definition,
				input: Input{A: 2, B: 1},
			},
			wantSame: false,
		},
		{
			name: "Same def but different literal input should have different hash",
			left: argument{
				def:   definition,
				input: 2,
			},
			right: argument{
				def:   definition,
				input: 1,
			},
			wantSame: false,
		},
		{
			name: "Different def same input should have different hash",
			left: argument{
				def:   definition,
				input: Input{A: 1, B: 2},
			},
			right: argument{
				def:   altDefinition,
				input: Input{A: 1, B: 2},
			},
			wantSame: false,
		},
		{
			name: "Invalid input should error",
			left: argument{
				def:   definition,
				input: math.NaN(),
			},
			right: argument{
				def:   definition,
				input: math.NaN(),
			},
			wantErr:        "json: unsupported value: NaN",
			skipCacheCheck: true,
		},
		{
			name: "different order of map keys should have the same hash",
			left: argument{
				def: definition,
				input: map[string]any{
					"a": 1,
					"b": 2,
					"c": 3,
				},
			},
			right: argument{
				def: definition,
				input: map[string]any{
					"b": 2,
					"c": 3,
					"a": 1,
				},
			},
			wantSame: true,
		},
		{
			name: "different order of nested struct in a slice should have the same hash",
			left: argument{
				def: definition,
				input: []map[string]any{
					{
						"a": 1,
						"b": 2,
					},
					{
						"b": 2,
						"a": 1,
					},
				}},
			right: argument{
				def: definition,
				input: []map[string]any{
					{
						"b": 2,
						"a": 1,
					},
					{
						"a": 1,
						"b": 2,
					},
				},
			},
			wantSame: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cache := &sync.Map{}
			hash1, err1 := constructUniqueHashFrom(cache, tt.left.def, tt.left.input)
			hash2, err2 := constructUniqueHashFrom(&sync.Map{}, tt.right.def, tt.right.input)

			if tt.wantErr != "" {
				require.Error(t, err1)
				require.Error(t, err2)
				require.ErrorContains(t, err1, tt.wantErr)
				require.ErrorContains(t, err2, tt.wantErr)
				return
			}

			require.NoError(t, err1)
			require.NoError(t, err2)

			if tt.wantSame {
				assert.Equal(t, hash1, hash2)
			} else {
				assert.NotEqual(t, hash1, hash2)
			}

			if !tt.skipCacheCheck {
				// Verify cache behavior
				defBytes1, err := json.Marshal(tt.left.def)
				require.NoError(t, err)
				inputBytes1, err := json.Marshal(tt.left.input)
				require.NoError(t, err)
				cacheKey1 := string(append(defBytes1, inputBytes1...))

				cached1, ok := cache.Load(cacheKey1)
				require.True(t, ok)
				assert.Equal(t, hash1, cached1)

				// Second call should use cache
				cachedHash1, err := constructUniqueHashFrom(cache, tt.left.def, tt.left.input)
				require.NoError(t, err)
				assert.Equal(t, hash1, cachedHash1)
			}
		})
	}
}
