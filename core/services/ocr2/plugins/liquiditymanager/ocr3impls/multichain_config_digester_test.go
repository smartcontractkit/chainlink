package ocr3impls_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/chains/evmutil"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/stretchr/testify/require"

	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/liquiditymanager/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/liquiditymanager/ocr3impls"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
)

func Test_MultichainConfigDigester(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		owner := testutils.MustNewSimTransactor(t)
		backend := backends.NewSimulatedBackend(core.GenesisAlloc{
			owner.From: core.GenesisAccount{
				Balance: assets.Ether(100).ToInt(),
			},
		}, 30e6)
		addr, _, _, err := no_op_ocr3.DeployNoOpOCR3(owner, backend)
		require.NoError(t, err, "failed to deploy contract")
		backend.Commit()
		wrapper, err := no_op_ocr3.NewNoOpOCR3(addr, backend)
		require.NoError(t, err, "failed to create wrapper")

		// masterChain needs to be 1337 since the on-chain config
		// digest is calculated using block.chainid, which is
		// always 1337 in the test.
		masterChain := commontypes.NewRelayID(relay.NetworkEVM, "1337")
		// rest of the chains don't matter for this test
		fChain1 := commontypes.NewRelayID(relay.NetworkEVM, "1338")
		fChain2 := commontypes.NewRelayID(relay.NetworkEVM, "1339")
		signers := []ocrtypes.OnchainPublicKey{
			testutils.NewAddress().Bytes(),
			testutils.NewAddress().Bytes(),
			testutils.NewAddress().Bytes(),
			testutils.NewAddress().Bytes(),
		}
		expectedSigners := func() (r []common.Address) {
			for _, signer := range signers {
				r = append(r, common.BytesToAddress(signer))
			}
			return
		}()
		masterTransmitters := randomOCRAccountArray(4)
		expectedTransmitters := func() (r []common.Address) {
			for _, transmitter := range masterTransmitters {
				r = append(r, common.HexToAddress(string(transmitter)))
			}
			return r
		}()
		fChain1Transmitters := randomOCRAccountArray(4)
		fChain2Transmitters := randomOCRAccountArray(4)
		_, err = wrapper.SetOCR3Config(owner, expectedSigners, expectedTransmitters, 1, []byte{}, 3, []byte{})
		require.NoError(t, err, "failed to set config")
		backend.Commit()
		iter, err := wrapper.FilterConfigSet(&bind.FilterOpts{
			Start: 1,
		})
		require.NoError(t, err, "failed to filter config set")
		var configDigest ocrtypes.ConfigDigest
		for iter.Next() {
			e := iter.Event
			require.Equal(t, uint8(1), e.F, "f doesn't match")
			require.Equal(t, uint64(1), e.ConfigCount, "config count doesn't match")
			require.Equal(t, []byte{}, e.OffchainConfig, "offchain config doesn't match")
			require.Equal(t, []byte{}, e.OnchainConfig, "onchain config doesn't match")
			require.Equal(t, uint64(3), e.OffchainConfigVersion, "offchain config version doesn't match")
			require.Equal(t, expectedSigners, e.Signers, "signers don't match")
			require.Equal(t, expectedTransmitters, e.Transmitters, "transmitters don't match")
			configDigest = e.ConfigDigest
		}
		require.NotEqual(t, ocrtypes.ConfigDigest{}, configDigest, "config digest must be nonzero")

		// follower chain config digests don't really matter for this test
		contractConfigs := map[commontypes.RelayID]ocrtypes.ContractConfig{
			masterChain: {
				ConfigDigest:          configDigest,
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
		require.Equal(t, configDigest, combined.ConfigDigest, "unexpected config digest")

		masterChainID, err := strconv.ParseUint(masterChain.ChainID, 10, 64)
		require.NoError(t, err)
		digester := ocr3impls.MultichainConfigDigester{
			MasterChainDigester: evmutil.EVMOffchainConfigDigester{
				ChainID:         masterChainID,
				ContractAddress: addr,
			},
		}
		offchainDigest, err := digester.ConfigDigest(combined)
		require.NoError(t, err)
		require.Equal(t, configDigest, offchainDigest, "offchain digest doesn't match onchain digest")
	})

	t.Run("signer is wrong length", func(t *testing.T) {
		masterChain := commontypes.NewRelayID(relay.NetworkEVM, "1337")
		masterChainID, err := strconv.ParseUint(masterChain.ChainID, 10, 64)
		require.NoError(t, err)
		digester := ocr3impls.MultichainConfigDigester{
			MasterChainDigester: evmutil.EVMOffchainConfigDigester{
				ChainID:         masterChainID,
				ContractAddress: testutils.NewAddress(),
			},
		}
		_, err = digester.ConfigDigest(ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{
				[]byte{1, 2, 3, 4}, // wrong length, must be 20 bytes
			},
		})
		require.Error(t, err)
		require.ErrorContains(t, err, "evm signer should be a 20 byte address, but got")
	})

	t.Run("num signers != num transmitters", func(t *testing.T) {
		masterChain := commontypes.NewRelayID(relay.NetworkEVM, "1337")
		masterChainID, err := strconv.ParseUint(masterChain.ChainID, 10, 64)
		require.NoError(t, err)
		digester := ocr3impls.MultichainConfigDigester{
			MasterChainDigester: evmutil.EVMOffchainConfigDigester{
				ChainID:         masterChainID,
				ContractAddress: testutils.NewAddress(),
			},
		}
		_, err = digester.ConfigDigest(ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
			},
			Transmitters: []ocrtypes.Account{
				ocrtypes.Account("1:2"),
				ocrtypes.Account("1:2"),
				ocrtypes.Account("1:2"),
				ocrtypes.Account("1:2"),
			},
		})
		require.Error(t, err)
		require.ErrorContains(t, err, "number of signers (3) does not match number of transmitters (4)")
	})

	t.Run("multi transmitter split fail", func(t *testing.T) {
		masterChain := commontypes.NewRelayID(relay.NetworkEVM, "1337")
		masterChainID, err := strconv.ParseUint(masterChain.ChainID, 10, 64)
		require.NoError(t, err)
		digester := ocr3impls.MultichainConfigDigester{
			MasterChainDigester: evmutil.EVMOffchainConfigDigester{
				ChainID:         masterChainID,
				ContractAddress: testutils.NewAddress(),
			},
		}
		_, err = digester.ConfigDigest(ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
			},
			Transmitters: []ocrtypes.Account{
				ocrtypes.Account("1:2:3:4"), // wrong split length, should be 2 not 4
				ocrtypes.Account("1:2:3:4"), // wrong split length, should be 2 not 4
				ocrtypes.Account("1:2:3:4"), // wrong split length, should be 2 not 4
				ocrtypes.Account("1:2:3:4"), // wrong split length, should be 2 not 4
			},
		})
		require.Error(t, err)
		require.ErrorContains(t, err, "unable to split multi-transmitter")
	})

	t.Run("wrong transmitter format", func(t *testing.T) {
		masterChain := commontypes.NewRelayID(relay.NetworkEVM, "1337")
		masterChainID, err := strconv.ParseUint(masterChain.ChainID, 10, 64)
		require.NoError(t, err)
		digester := ocr3impls.MultichainConfigDigester{
			MasterChainDigester: evmutil.EVMOffchainConfigDigester{
				ChainID:         masterChainID,
				ContractAddress: testutils.NewAddress(),
			},
		}
		_, err = digester.ConfigDigest(ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
			},
			Transmitters: []ocrtypes.Account{
				ocrtypes.Account("1337:blahblah"), // wrong transmitter format, must be 0x... address
				ocrtypes.Account("1337:blahblah"), // wrong transmitter format, must be 0x... address
				ocrtypes.Account("1337:blahblah"), // wrong transmitter format, must be 0x... address
				ocrtypes.Account("1337:blahblah"), // wrong transmitter format, must be 0x... address
			},
		})
		require.Error(t, err)
		require.ErrorContains(t, err, "evm transmitter should be a 42 character Ethereum address string, but got")
	})

	t.Run("multiple transmitters for master chain id", func(t *testing.T) {
		masterChain := commontypes.NewRelayID(relay.NetworkEVM, "1337")
		masterChainID, err := strconv.ParseUint(masterChain.ChainID, 10, 64)
		require.NoError(t, err)
		digester := ocr3impls.MultichainConfigDigester{
			MasterChainDigester: evmutil.EVMOffchainConfigDigester{
				ChainID:         masterChainID,
				ContractAddress: testutils.NewAddress(),
			},
		}
		_, err = digester.ConfigDigest(ocrtypes.ContractConfig{
			Signers: []ocrtypes.OnchainPublicKey{
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
				testutils.NewAddress().Bytes(),
			},
			Transmitters: []ocrtypes.Account{
				// chain id 1 shows up twice
				ocrtypes.Account(fmt.Sprintf("1337:%s,1337:%s", testutils.NewAddress().Hex(), testutils.NewAddress().Hex())),
				ocrtypes.Account(fmt.Sprintf("1:%s,1:%s", testutils.NewAddress().Hex(), testutils.NewAddress().Hex())),
				ocrtypes.Account(fmt.Sprintf("1:%s,1:%s", testutils.NewAddress().Hex(), testutils.NewAddress().Hex())),
				ocrtypes.Account(fmt.Sprintf("1:%s,1:%s", testutils.NewAddress().Hex(), testutils.NewAddress().Hex())),
			},
		})
		require.Error(t, err)
		require.ErrorContains(t, err, "same chain id appearing multiple times in parts")
	})
}
