package ocr3impls_test

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/ocr3impls"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func Test_TransmitterCombiner(t *testing.T) {
	masterChain := commontypes.NewRelayID(relay.NetworkEVM, "1")
	fChain1 := commontypes.NewRelayID(relay.NetworkEVM, "2")
	fChain2 := commontypes.NewRelayID(relay.NetworkEVM, "3")
	signers := []ocrtypes.OnchainPublicKey{
		testutils.NewAddress().Bytes(),
		testutils.NewAddress().Bytes(),
		testutils.NewAddress().Bytes(),
		testutils.NewAddress().Bytes(),
	}
	masterTransmitters := randomOCRAccountArray(4)
	fChain1Transmitters := randomOCRAccountArray(4)
	fChain2Transmitters := randomOCRAccountArray(4)

	t.Run("master chain config not provided", func(t *testing.T) {
		contractConfigs := map[commontypes.RelayID]ocrtypes.ContractConfig{
			fChain1: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          fChain1Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain2: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          fChain2Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
		}
		_, err := ocr3impls.TransmitterCombiner(masterChain, contractConfigs)
		require.Error(t, err, "expected error")
	})

	t.Run("mismatched signers and transmitters lengths", func(t *testing.T) {
		contractConfigs := map[commontypes.RelayID]ocrtypes.ContractConfig{
			masterChain: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          masterTransmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain1: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          fChain1Transmitters[:len(fChain1Transmitters)-1], // bad config here
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain2: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          fChain2Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
		}
		_, err := ocr3impls.TransmitterCombiner(masterChain, contractConfigs)
		require.Error(t, err, "expected error")
	})

	t.Run("signer not found on follower chain config", func(t *testing.T) {
		contractConfigs := map[commontypes.RelayID]ocrtypes.ContractConfig{
			masterChain: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          masterTransmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain1: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers[:len(signers)-1], // bad config here
				Transmitters:          fChain1Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain2: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          fChain2Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
		}
		_, err := ocr3impls.TransmitterCombiner(masterChain, contractConfigs)
		require.Error(t, err, "expected error")
	})

	t.Run("happy path", func(t *testing.T) {
		contractConfigs := map[commontypes.RelayID]ocrtypes.ContractConfig{
			masterChain: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          masterTransmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain1: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          fChain1Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain2: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          fChain2Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
		}
		combined, err := ocr3impls.TransmitterCombiner(masterChain, contractConfigs)
		require.NoError(t, err, "TransmitterCombiner should not error")
		// check sorted order of transmitters
		// due to lexicographic sorting of "chainID:address" strings
		// it should be sorted in increasing chain ID order
		expectedTransmitters := []ocrtypes.Account{
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[0], fChain1.ChainID, fChain1Transmitters[0], fChain2.ChainID, fChain2Transmitters[0])),
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[1], fChain1.ChainID, fChain1Transmitters[1], fChain2.ChainID, fChain2Transmitters[1])),
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[2], fChain1.ChainID, fChain1Transmitters[2], fChain2.ChainID, fChain2Transmitters[2])),
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[3], fChain1.ChainID, fChain1Transmitters[3], fChain2.ChainID, fChain2Transmitters[3])),
		}
		require.Equal(t, contractConfigs[masterChain].ConfigDigest, combined.ConfigDigest)
		require.Equal(t, contractConfigs[masterChain].ConfigCount, combined.ConfigCount)
		require.Equal(t, contractConfigs[masterChain].F, combined.F)
		require.Equal(t, contractConfigs[masterChain].OffchainConfig, combined.OffchainConfig)
		require.Equal(t, contractConfigs[masterChain].OffchainConfigVersion, combined.OffchainConfigVersion)
		require.Equal(t, contractConfigs[masterChain].OnchainConfig, combined.OnchainConfig)
		require.Equal(t, contractConfigs[masterChain].Signers, combined.Signers)
		require.Equal(t, expectedTransmitters, combined.Transmitters)
	})

	t.Run("different signer order on follower chains", func(t *testing.T) {
		fChain1Signers := shuffledSlice(signers)
		fChain2Signers := shuffledSlice(signers)
		contractConfigs := map[commontypes.RelayID]ocrtypes.ContractConfig{
			masterChain: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               signers,
				Transmitters:          masterTransmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain1: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               fChain1Signers,
				Transmitters:          fChain1Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
			fChain2: {
				ConfigDigest:          testutils.Random32Byte(),
				ConfigCount:           1,
				Signers:               fChain2Signers,
				Transmitters:          fChain2Transmitters,
				F:                     1,
				OnchainConfig:         []byte{},
				OffchainConfigVersion: 3,
				OffchainConfig:        []byte{},
			},
		}
		combined, err := ocr3impls.TransmitterCombiner(masterChain, contractConfigs)
		require.NoError(t, err, "TransmitterCombiner should not error")
		// check sorted order of transmitters
		// due to lexicographic sorting of "chainID:address" strings
		// it should be sorted in increasing chain ID order
		expectedTransmitters := []ocrtypes.Account{
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[0],
				fChain1.ChainID, fChain1Transmitters[getSignerIndex(t, signers[0], fChain1Signers)],
				fChain2.ChainID, fChain2Transmitters[getSignerIndex(t, signers[0], fChain2Signers)])),
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[1],
				fChain1.ChainID, fChain1Transmitters[getSignerIndex(t, signers[1], fChain1Signers)],
				fChain2.ChainID, fChain2Transmitters[getSignerIndex(t, signers[1], fChain2Signers)])),
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[2],
				fChain1.ChainID, fChain1Transmitters[getSignerIndex(t, signers[2], fChain1Signers)],
				fChain2.ChainID, fChain2Transmitters[getSignerIndex(t, signers[2], fChain2Signers)])),
			ocrtypes.Account(fmt.Sprintf("%s:%s,%s:%s,%s:%s",
				masterChain.ChainID, masterTransmitters[3],
				fChain1.ChainID, fChain1Transmitters[getSignerIndex(t, signers[3], fChain1Signers)],
				fChain2.ChainID, fChain2Transmitters[getSignerIndex(t, signers[3], fChain2Signers)])),
		}
		require.Equal(t, contractConfigs[masterChain].ConfigDigest, combined.ConfigDigest)
		require.Equal(t, contractConfigs[masterChain].ConfigCount, combined.ConfigCount)
		require.Equal(t, contractConfigs[masterChain].F, combined.F)
		require.Equal(t, contractConfigs[masterChain].OffchainConfig, combined.OffchainConfig)
		require.Equal(t, contractConfigs[masterChain].OffchainConfigVersion, combined.OffchainConfigVersion)
		require.Equal(t, contractConfigs[masterChain].OnchainConfig, combined.OnchainConfig)
		require.Equal(t, contractConfigs[masterChain].Signers, combined.Signers)
		require.Equal(t, expectedTransmitters, combined.Transmitters)
	})
}

func Test_SplitMultiTransmitter(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		multiTransmitter := "1:a1,2:a2,3:a3,4:a4"
		ret, err := ocr3impls.SplitMultiTransmitter(ocrtypes.Account(multiTransmitter))
		require.NoError(t, err)
		require.Contains(t, ret, commontypes.NewRelayID(relay.NetworkEVM, "1"))
		require.Contains(t, ret, commontypes.NewRelayID(relay.NetworkEVM, "2"))
		require.Contains(t, ret, commontypes.NewRelayID(relay.NetworkEVM, "3"))
		require.Contains(t, ret, commontypes.NewRelayID(relay.NetworkEVM, "4"))
		require.Equal(t, ocrtypes.Account("a1"), ret[commontypes.NewRelayID(relay.NetworkEVM, "1")])
		require.Equal(t, ocrtypes.Account("a2"), ret[commontypes.NewRelayID(relay.NetworkEVM, "2")])
		require.Equal(t, ocrtypes.Account("a3"), ret[commontypes.NewRelayID(relay.NetworkEVM, "3")])
		require.Equal(t, ocrtypes.Account("a4"), ret[commontypes.NewRelayID(relay.NetworkEVM, "4")])
	})

	t.Run("num parts is not 2", func(t *testing.T) {
		multiTransmitter := "1:lol:0:a1,2:rofl:0:a2,3:hex:0:a3,4:dex:0:a4"
		_, err := ocr3impls.SplitMultiTransmitter(ocrtypes.Account(multiTransmitter))
		require.Error(t, err)
		require.ErrorContains(t, err, "split on ':' must contain exactly 2 parts, got:")
	})

	t.Run("same chain id appearing multiple times", func(t *testing.T) {
		// chain id 1 appears twice
		multiTransmitter := "1:a1,1:a2,3:a3,4:a4"
		_, err := ocr3impls.SplitMultiTransmitter(ocrtypes.Account(multiTransmitter))
		require.Error(t, err)
		require.ErrorContains(t, err, "same chain id appearing multiple times in parts")
	})
}

func randomOCRAccountArray(n int) []ocrtypes.Account {
	addresses := make([]ocrtypes.Account, n)
	for i := 0; i < n; i++ {
		addresses[i] = ocrtypes.Account(testutils.NewAddress().Hex())
	}
	return addresses
}

func shuffledSlice[T any](s []T) []T {
	cpy := make([]T, len(s))
	copy(cpy, s)
	rand.Shuffle(len(s), func(i, j int) {
		cpy[i], cpy[j] = cpy[j], cpy[i]
	})
	return cpy
}

func getSignerIndex(t *testing.T, masterSigner ocrtypes.OnchainPublicKey, followerSigners []ocrtypes.OnchainPublicKey) int {
	idx := slices.IndexFunc(followerSigners, func(opk ocrtypes.OnchainPublicKey) bool {
		return hexutil.Encode(masterSigner) == hexutil.Encode(opk)
	})
	require.NotEqual(t, -1, idx, "expected something other than -1")
	return idx
}
