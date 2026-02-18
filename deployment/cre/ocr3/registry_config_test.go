package ocr3

import (
	"encoding/base64"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateOCR3Config(t *testing.T) {
	tests := []struct {
		name         string
		signers      [][]byte
		transmitters [][]byte
		f            uint32
		wantErr      string
	}{
		{
			name:         "valid config with 4 nodes and f=1",
			signers:      make([][]byte, 4),
			transmitters: make([][]byte, 4),
			f:            1,
		},
		{
			name:         "valid config with 10 nodes and f=3",
			signers:      make([][]byte, 10),
			transmitters: make([][]byte, 10),
			f:            3,
		},
		{
			name:         "signers and transmitters count mismatch",
			signers:      make([][]byte, 4),
			transmitters: make([][]byte, 3),
			f:            1,
			wantErr:      "signers count (4) != transmitters count (3)",
		},
		{
			name:         "f must be positive",
			signers:      make([][]byte, 4),
			transmitters: make([][]byte, 4),
			f:            0,
			wantErr:      "f must be positive",
		},
		{
			name:         "not enough nodes for f (3f+1 rule)",
			signers:      make([][]byte, 3),
			transmitters: make([][]byte, 3),
			f:            1,
			wantErr:      "not enough nodes for f=1: need at least 4, got 3",
		},
		{
			name:         "exactly 3f nodes fails (need 3f+1)",
			signers:      make([][]byte, 6),
			transmitters: make([][]byte, 6),
			f:            2,
			wantErr:      "not enough nodes for f=2: need at least 7, got 6",
		},
		{
			name:    "empty transmitter",
			signers: [][]byte{{0x01}, {0x02}, {0x03}, {0x04}},
			transmitters: [][]byte{
				{0x01},
				{}, // empty
				{0x03},
				{0x04},
			},
			f:       1,
			wantErr: "transmitter 1 is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == "" || tt.wantErr == "signers count" || tt.wantErr == "f must be positive" {
				for i := range tt.transmitters {
					if len(tt.transmitters[i]) == 0 {
						tt.transmitters[i] = []byte{byte(i + 1)}
					}
				}
			}

			err := ValidateOCR3Config(tt.signers, tt.transmitters, tt.f)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOCR2OracleConfigToMap(t *testing.T) {
	config := &OCR2OracleConfig{
		Signers: [][]byte{
			{0x01, 0x02, 0x03},
			{0x04, 0x05, 0x06},
		},
		Transmitters: []common.Address{
			common.HexToAddress("0x1111111111111111111111111111111111111111"),
			common.HexToAddress("0x2222222222222222222222222222222222222222"),
		},
		F:                     2,
		OnchainConfig:         []byte{0xAA, 0xBB},
		OffchainConfigVersion: 3,
		OffchainConfig:        []byte{0xCC, 0xDD},
	}

	result := OCR2OracleConfigToMap(config, 5)

	signers, ok := result["signers"].([]string)
	require.True(t, ok)
	require.Len(t, signers, 2)
	decoded, err := base64.StdEncoding.DecodeString(signers[0])
	require.NoError(t, err)
	assert.Equal(t, []byte{0x01, 0x02, 0x03}, decoded)

	transmitters, ok := result["transmitters"].([]string)
	require.True(t, ok)
	require.Len(t, transmitters, 2)
	decoded, err = base64.StdEncoding.DecodeString(transmitters[0])
	require.NoError(t, err)
	assert.Equal(t, common.HexToAddress("0x1111111111111111111111111111111111111111").Bytes(), decoded)

	assert.Equal(t, uint32(2), result["f"])
	assert.Equal(t, uint64(3), result["offchainConfigVersion"])
	assert.Equal(t, uint64(5), result["configCount"])

	onchainConfig, ok := result["onchainConfig"].(string)
	require.True(t, ok)
	decoded, err = base64.StdEncoding.DecodeString(onchainConfig)
	require.NoError(t, err)
	assert.Equal(t, []byte{0xAA, 0xBB}, decoded)
}

func TestOCR2OracleConfigToMap_EmptyOnchainConfig(t *testing.T) {
	config := &OCR2OracleConfig{
		Signers:               [][]byte{{0x01}},
		Transmitters:          []common.Address{common.HexToAddress("0x1111111111111111111111111111111111111111")},
		F:                     1,
		OffchainConfigVersion: 1,
		OffchainConfig:        []byte{0xCC},
	}

	result := OCR2OracleConfigToMap(config, 1)

	_, hasOnchainConfig := result["onchainConfig"]
	assert.False(t, hasOnchainConfig, "onchainConfig should be omitted when empty")
}

func TestValidateOCR2OracleConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := &OCR2OracleConfig{
			Signers: [][]byte{{0x01}, {0x02}, {0x03}, {0x04}},
			Transmitters: []common.Address{
				common.HexToAddress("0x1111111111111111111111111111111111111111"),
				common.HexToAddress("0x2222222222222222222222222222222222222222"),
				common.HexToAddress("0x3333333333333333333333333333333333333333"),
				common.HexToAddress("0x4444444444444444444444444444444444444444"),
			},
			F: 1,
		}
		require.NoError(t, ValidateOCR2OracleConfig(config))
	})

	t.Run("zero address transmitter", func(t *testing.T) {
		config := &OCR2OracleConfig{
			Signers: [][]byte{{0x01}, {0x02}, {0x03}, {0x04}},
			Transmitters: []common.Address{
				common.HexToAddress("0x1111111111111111111111111111111111111111"),
				{}, // zero address
				common.HexToAddress("0x3333333333333333333333333333333333333333"),
				common.HexToAddress("0x4444444444444444444444444444444444444444"),
			},
			F: 1,
		}
		err := ValidateOCR2OracleConfig(config)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "zero address")
	})
}
