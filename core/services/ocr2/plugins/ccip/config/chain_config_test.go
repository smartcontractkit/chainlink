package config

import (
	"math/big"
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/common/chains/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestGetChainFromSpec(t *testing.T) {
	testChainID := int64(1337)

	tests := []struct {
		name           string
		spec           *job.OCR2OracleSpec
		expectedErr    bool
		expectedErrMsg string
	}{
		{
			name: "success",
			spec: &job.OCR2OracleSpec{
				RelayConfig: job.JSONConfig{
					"chainID": float64(testChainID),
				},
			},
			expectedErr: false,
		},
		{
			name:           "missing_chain_ID",
			spec:           &job.OCR2OracleSpec{},
			expectedErr:    true,
			expectedErrMsg: "chainID must be provided in relay config",
		},
	}

	mockChain := mocks.NewChain(t)
	mockChain.On("ID").Return(big.NewInt(testChainID)).Maybe()

	mockChainSet := mocks.NewLegacyChainContainer(t)
	mockChainSet.On("Get", strconv.FormatInt(testChainID, 10)).Return(mockChain, nil).Maybe()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chain, chainID, err := GetChainFromSpec(test.spec, mockChainSet)
			if test.expectedErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), test.expectedErrMsg)
			} else {
				require.NoError(t, err)
				require.Equal(t, mockChain, chain)
				require.Equal(t, testChainID, chainID)
			}
		})
	}
}

func TestGetChainByChainSelector_success(t *testing.T) {
	mockChain := mocks.NewChain(t)
	mockChain.On("ID").Return(big.NewInt(11155111))

	mockChainSet := mocks.NewLegacyChainContainer(t)
	mockChainSet.On("Get", "11155111").Return(mockChain, nil)

	// Ethereum Sepolia chain selector.
	chain, chainID, err := GetChainByChainSelector(mockChainSet, uint64(16015286601757825753))
	require.NoError(t, err)
	require.Equal(t, mockChain, chain)
	require.Equal(t, int64(11155111), chainID)
}

func TestGetChainByChainSelector_selectorNotFound(t *testing.T) {
	mockChainSet := mocks.NewLegacyChainContainer(t)

	_, _, err := GetChainByChainSelector(mockChainSet, uint64(444000444))
	require.Error(t, err)
}

func TestGetChainById_notFound(t *testing.T) {
	mockChainSet := mocks.NewLegacyChainContainer(t)
	mockChainSet.On("Get", "444").Return(nil, errors.New("test")).Maybe()

	_, _, err := GetChainByChainID(mockChainSet, uint64(444))
	require.Error(t, err)
	require.Contains(t, err.Error(), "chain not found in chainset")
}

func TestResolveChainNames(t *testing.T) {
	tests := []struct {
		name                    string
		sourceChainId           int64
		destChainId             int64
		expectedSourceChainName string
		expectedDestChainName   string
		expectedErr             bool
	}{
		{
			name:                    "success",
			sourceChainId:           1,
			destChainId:             10,
			expectedSourceChainName: "ethereum-mainnet",
			expectedDestChainName:   "ethereum-mainnet-optimism-1",
		},
		{
			name:          "source chain not found",
			sourceChainId: 901278309182,
			destChainId:   10,
			expectedErr:   true,
		},
		{
			name:          "dest chain not found",
			sourceChainId: 1,
			destChainId:   901278309182,
			expectedErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sourceChainName, destChainName, err := ResolveChainNames(test.sourceChainId, test.destChainId)
			if test.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.expectedSourceChainName, sourceChainName)
				assert.Equal(t, test.expectedDestChainName, destChainName)
			}
		})
	}
}
