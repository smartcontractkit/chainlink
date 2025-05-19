package deployment

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKMSToEthSigConversion(t *testing.T) {
	kmsSigBytes, err := hex.DecodeString("304402206168865941bafcae3a8cf8b26edbb5693d62222b2e54d962c1aabbeaddf33b6802205edc7f597d2bf2d1eaa14fc514a6202bafcffe52b13ae3fec00674d92a874b73")
	require.NoError(t, err)
	ecdsaPublicKeyBytes, err := hex.DecodeString("04a735e9e3cb526f83be23b03f1f5ae7788a8654e3f0fcfb4f978290de07ebd47da30eeb72e904fdd4a81b46e320908ff4345e119148f89c1f04674c14a506e24b")
	require.NoError(t, err)
	txHashBytes, err := hex.DecodeString("a2f037301e90f58c084fe4bec2eef14b26e620d6b6cb46051037d03b29ab7d9a")
	require.NoError(t, err)
	expectedEthSignBytes, err := hex.DecodeString("6168865941bafcae3a8cf8b26edbb5693d62222b2e54d962c1aabbeaddf33b685edc7f597d2bf2d1eaa14fc514a6202bafcffe52b13ae3fec00674d92a874b7300")
	require.NoError(t, err)

	actualEthSig, err := kmsToEthSig(
		kmsSigBytes,
		ecdsaPublicKeyBytes,
		txHashBytes,
	)
	require.NoError(t, err)
	require.Equal(t, expectedEthSignBytes, actualEthSig)
}
