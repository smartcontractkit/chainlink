package evm_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"
)

func TestExtractNetwork(t *testing.T) {
	testCases := []struct {
		networkName  string
		expectedName string
		expectedErr  bool
	}{
		{
			networkName:  "ethereum-testnet-goerli",
			expectedName: "testnet",
			expectedErr:  false,
		},
		{
			networkName:  "ethereum-mainnet",
			expectedName: "mainnet",
			expectedErr:  false,
		},
		{
			networkName:  "polygon-devnet",
			expectedName: "devnet",
			expectedErr:  false,
		},
		{
			networkName:  "ethereum_test",
			expectedName: "",
			expectedErr:  true,
		},
		{
			networkName:  "ethereum",
			expectedName: "",
			expectedErr:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.networkName, func(t *testing.T) {
			networkName, err := chainselectors.ExtractNetworkEnvName(tc.networkName)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedName, networkName)
		})
	}
}
