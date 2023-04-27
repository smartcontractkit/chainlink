package coordinator

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/libocr/commontypes"
	ocr2Types "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"

	evmclimocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lp_mocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	dkg_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2vrf/coordinator/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/utils/mathutil"
)

func TestCoordinator_BeaconPeriod(t *testing.T) {
	t.Parallel()

	t.Run("valid output", func(t *testing.T) {
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("IBeaconPeriodBlocks", mock.Anything).
			Return(big.NewInt(10), nil)
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		period, err := c.BeaconPeriod(testutils.Context(t))
		assert.NoError(t, err)
		assert.Equal(t, uint16(10), period)
	})

	t.Run("invalid output", func(t *testing.T) {
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("IBeaconPeriodBlocks", mock.Anything).
			Return(nil, errors.New("rpc error"))
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		_, err := c.BeaconPeriod(testutils.Context(t))
		assert.Error(t, err)
	})
}

func TestCoordinator_DKGVRFCommittees(t *testing.T) {
	t.Parallel()
	evmClient := evmclimocks.NewClient(t)
	evmClient.On("ConfiguredChainID").Return(big.NewInt(1))

	t.Run("happy path", func(t *testing.T) {
		// In this test the DKG and VRF committees have the same signers and
		// transmitters. This may (?) be different in practice.

		lp := lp_mocks.NewLogPoller(t)
		tp := newTopics()

		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		dkgAddress := newAddress(t)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, beaconAddress, 10).
			Return(&logpoller.Log{
				Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000a6fca200010576e704b4a519484d6239ef17f1f5b4a82e330b0daf827ed4dc2789971b0000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000002e000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000a8cbea12a06869d3ec432ab9682dab6c761d591000000000000000000000000f4f9db7bb1d16b7cdfb18ec68994c26964f5985300000000000000000000000022fb3f90c539457f00d8484438869135e604a65500000000000000000000000033cbcedccb11c9773ad78e214ba342e979255ab30000000000000000000000006ffaa96256fbc1012325cca88c79f725c33eed80000000000000000000000000000000000000000000000000000000000000000500000000000000000000000074103cf8b436465870b26aa9fa2f62ad62b22e3500000000000000000000000038a6cb196f805cc3041f6645a5a6cec27b64430d00000000000000000000000047d7095cfebf8285bdaa421bc8268d0db87d933c000000000000000000000000a8842be973800ff61d80d2d53fa62c3a685380eb0000000000000000000000003750e31321aee8c024751877070e8d5f704ce98700000000000000000000000000000000000000000000000000000000000000206f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf000000000000000000000000000000000000000000000000000000000000004220880d88ee16f1080c8afa0251880c8afa025208090dfc04a288090dfc04a30033a05010101010142206c5ca6f74b532222ac927dd3de235d46a943e372c0563393a33b01dcfd3f371c4220855114d25c2ef5e85fffe4f20a365672d8f2dba3b2ec82333f494168a2039c0442200266e835634db00977cbc1caa4db10e1676c1a4c0fcbc6ba7f09300f0d1831824220980cd91f7a73f20f4b0d51d00cd4e00373dc2beafbb299ca3c609757ab98c8304220eb6d36e2af8922085ff510bbe1eb8932a0e3295ca9f047fef25d90e69c52948f4a34313244334b6f6f574463364b7232644542684b59326b336e685057694676544565325331703978544532544b74344d7572716f684a34313244334b6f6f574b436e4367724b637743324a3577576a626e355435335068646b6b6f57454e534a39546537544b7836366f4a4a34313244334b6f6f575239616f675948786b357a38636b624c4c56346e426f7a777a747871664a7050586671336d4a7232796452474a34313244334b6f6f5744695444635565675637776b313133473366476a69616259756f54436f3157726f6f53656741343263556f544a34313244334b6f6f574e64687072586b5472665370354d5071736270467a70364167394a53787358694341434442676454424c656652820300050e416c74424e2d3132382047e282810e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c216a812d7616e47f0bd38fa4863f48fbcda6a38af4c58d2233dfa7cf79620947042d09f923e0a2f7a2270391e8b058d8bdb8f79fe082b7b627f025651c7290382fdff97c3181d15d162c146ce87ff752499d2acc2b26011439a12e29571a6f1e1defb1751c3be4258c493984fd9f0f6b4a26c539870b5f15bfed3d8ffac92499eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b127a112970e1adf615f823b2b2180754c2f0ee01f1b389e56df55ca09702cd0401b66ff71779d2dd67222503a85ab921b28c329cc1832800b192d0b0247c0776e1b9653dc00df48daa6364287c84c0382f5165e7269fef06d10bc67c1bba252305d1af0dc7bb0fe92558eb4c5f38c23163dee1cfb34a72020669dbdfe337c16f3307472616e736c61746f722066726f6d20416c74424e2d3132382047e2828120746f20416c74424e2d3132382047e282825880ade2046080c8afa0256880c8afa0257080ade204788094ebdc0382019e010a205034214e0bd4373f38e162cf9fc9133e2f3b71441faa4c3d1ac01c1877f1cd2712200e03e975b996f911abba2b79d2596c2150bc94510963c40a1137a03df6edacdb1a107dee1cdb894163813bb3da604c9c133c1a10bb33302eeafbd55d352e35dcc5d2b3311a10d2c658b6b93d74a02d467849b6fe75251a10fea5308cc1fea69e7246eafe7ca8a3a51a1048efe1ad873b6f025ac0243bdef715f8000000000000000000000000000000000000000000000000000000000000"),
			}, nil)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, dkgAddress, 10).
			Return(&logpoller.Log{
				Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000a6fca200010576e704b4a519484d6239ef17f1f5b4a82e330b0daf827ed4dc2789971b0000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000002e000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000a8cbea12a06869d3ec432ab9682dab6c761d591000000000000000000000000f4f9db7bb1d16b7cdfb18ec68994c26964f5985300000000000000000000000022fb3f90c539457f00d8484438869135e604a65500000000000000000000000033cbcedccb11c9773ad78e214ba342e979255ab30000000000000000000000006ffaa96256fbc1012325cca88c79f725c33eed80000000000000000000000000000000000000000000000000000000000000000500000000000000000000000074103cf8b436465870b26aa9fa2f62ad62b22e3500000000000000000000000038a6cb196f805cc3041f6645a5a6cec27b64430d00000000000000000000000047d7095cfebf8285bdaa421bc8268d0db87d933c000000000000000000000000a8842be973800ff61d80d2d53fa62c3a685380eb0000000000000000000000003750e31321aee8c024751877070e8d5f704ce98700000000000000000000000000000000000000000000000000000000000000206f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf000000000000000000000000000000000000000000000000000000000000004220880d88ee16f1080c8afa0251880c8afa025208090dfc04a288090dfc04a30033a05010101010142206c5ca6f74b532222ac927dd3de235d46a943e372c0563393a33b01dcfd3f371c4220855114d25c2ef5e85fffe4f20a365672d8f2dba3b2ec82333f494168a2039c0442200266e835634db00977cbc1caa4db10e1676c1a4c0fcbc6ba7f09300f0d1831824220980cd91f7a73f20f4b0d51d00cd4e00373dc2beafbb299ca3c609757ab98c8304220eb6d36e2af8922085ff510bbe1eb8932a0e3295ca9f047fef25d90e69c52948f4a34313244334b6f6f574463364b7232644542684b59326b336e685057694676544565325331703978544532544b74344d7572716f684a34313244334b6f6f574b436e4367724b637743324a3577576a626e355435335068646b6b6f57454e534a39546537544b7836366f4a4a34313244334b6f6f575239616f675948786b357a38636b624c4c56346e426f7a777a747871664a7050586671336d4a7232796452474a34313244334b6f6f5744695444635565675637776b313133473366476a69616259756f54436f3157726f6f53656741343263556f544a34313244334b6f6f574e64687072586b5472665370354d5071736270467a70364167394a53787358694341434442676454424c656652820300050e416c74424e2d3132382047e282810e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c216a812d7616e47f0bd38fa4863f48fbcda6a38af4c58d2233dfa7cf79620947042d09f923e0a2f7a2270391e8b058d8bdb8f79fe082b7b627f025651c7290382fdff97c3181d15d162c146ce87ff752499d2acc2b26011439a12e29571a6f1e1defb1751c3be4258c493984fd9f0f6b4a26c539870b5f15bfed3d8ffac92499eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b127a112970e1adf615f823b2b2180754c2f0ee01f1b389e56df55ca09702cd0401b66ff71779d2dd67222503a85ab921b28c329cc1832800b192d0b0247c0776e1b9653dc00df48daa6364287c84c0382f5165e7269fef06d10bc67c1bba252305d1af0dc7bb0fe92558eb4c5f38c23163dee1cfb34a72020669dbdfe337c16f3307472616e736c61746f722066726f6d20416c74424e2d3132382047e2828120746f20416c74424e2d3132382047e282825880ade2046080c8afa0256880c8afa0257080ade204788094ebdc0382019e010a205034214e0bd4373f38e162cf9fc9133e2f3b71441faa4c3d1ac01c1877f1cd2712200e03e975b996f911abba2b79d2596c2150bc94510963c40a1137a03df6edacdb1a107dee1cdb894163813bb3da604c9c133c1a10bb33302eeafbd55d352e35dcc5d2b3311a10d2c658b6b93d74a02d467849b6fe75251a10fea5308cc1fea69e7246eafe7ca8a3a51a1048efe1ad873b6f025ac0243bdef715f8000000000000000000000000000000000000000000000000000000000000"),
			}, nil)

		expectedDKGVRF := ocr2vrftypes.OCRCommittee{
			Signers: []common.Address{
				common.HexToAddress("0x0A8cbEA12a06869d3EC432aB9682DAb6C761D591"),
				common.HexToAddress("0xF4f9db7BB1d16b7CDfb18Ec68994c26964F59853"),
				common.HexToAddress("0x22fB3F90C539457f00d8484438869135E604a655"),
				common.HexToAddress("0x33CbCedccb11c9773AD78e214Ba342E979255ab3"),
				common.HexToAddress("0x6ffaA96256fbC1012325cca88C79F725c33eED80"),
			},
			Transmitters: []common.Address{
				common.HexToAddress("0x74103Cf8b436465870b26aa9Fa2F62AD62b22E35"),
				common.HexToAddress("0x38A6Cb196f805cC3041F6645a5A6CEC27B64430D"),
				common.HexToAddress("0x47d7095CFEBF8285BdAa421Bc8268D0DB87D933C"),
				common.HexToAddress("0xa8842BE973800fF61D80d2d53fa62C3a685380eB"),
				common.HexToAddress("0x3750e31321aEE8c024751877070E8d5F704cE987"),
			},
		}

		c := &coordinator{
			lp:                 lp,
			lggr:               logger.TestLogger(t),
			topics:             tp,
			beaconAddress:      beaconAddress,
			coordinatorAddress: coordinatorAddress,
			dkgAddress:         dkgAddress,
			finalityDepth:      10,
			evmClient:          evmClient,
		}
		actualDKG, actualVRF, err := c.DKGVRFCommittees(testutils.Context(t))
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedDKGVRF.Signers, actualDKG.Signers)
		assert.ElementsMatch(t, expectedDKGVRF.Transmitters, actualDKG.Transmitters)
		assert.ElementsMatch(t, expectedDKGVRF.Signers, actualVRF.Signers)
		assert.ElementsMatch(t, expectedDKGVRF.Transmitters, actualVRF.Transmitters)
	})

	t.Run("vrf log poll fails", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		tp := newTopics()

		beaconAddress := newAddress(t)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, beaconAddress, 10).
			Return(nil, errors.New("rpc error"))

		c := &coordinator{
			lp:            lp,
			lggr:          logger.TestLogger(t),
			topics:        tp,
			beaconAddress: beaconAddress,
			finalityDepth: 10,
			evmClient:     evmClient,
		}

		_, _, err := c.DKGVRFCommittees(testutils.Context(t))
		assert.Error(t, err)
	})

	t.Run("dkg log poll fails", func(t *testing.T) {
		lp := lp_mocks.NewLogPoller(t)
		tp := newTopics()
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)
		dkgAddress := newAddress(t)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, beaconAddress, 10).
			Return(&logpoller.Log{
				Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000a6fca200010576e704b4a519484d6239ef17f1f5b4a82e330b0daf827ed4dc2789971b0000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000002e000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000a8cbea12a06869d3ec432ab9682dab6c761d591000000000000000000000000f4f9db7bb1d16b7cdfb18ec68994c26964f5985300000000000000000000000022fb3f90c539457f00d8484438869135e604a65500000000000000000000000033cbcedccb11c9773ad78e214ba342e979255ab30000000000000000000000006ffaa96256fbc1012325cca88c79f725c33eed80000000000000000000000000000000000000000000000000000000000000000500000000000000000000000074103cf8b436465870b26aa9fa2f62ad62b22e3500000000000000000000000038a6cb196f805cc3041f6645a5a6cec27b64430d00000000000000000000000047d7095cfebf8285bdaa421bc8268d0db87d933c000000000000000000000000a8842be973800ff61d80d2d53fa62c3a685380eb0000000000000000000000003750e31321aee8c024751877070e8d5f704ce98700000000000000000000000000000000000000000000000000000000000000206f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf000000000000000000000000000000000000000000000000000000000000004220880d88ee16f1080c8afa0251880c8afa025208090dfc04a288090dfc04a30033a05010101010142206c5ca6f74b532222ac927dd3de235d46a943e372c0563393a33b01dcfd3f371c4220855114d25c2ef5e85fffe4f20a365672d8f2dba3b2ec82333f494168a2039c0442200266e835634db00977cbc1caa4db10e1676c1a4c0fcbc6ba7f09300f0d1831824220980cd91f7a73f20f4b0d51d00cd4e00373dc2beafbb299ca3c609757ab98c8304220eb6d36e2af8922085ff510bbe1eb8932a0e3295ca9f047fef25d90e69c52948f4a34313244334b6f6f574463364b7232644542684b59326b336e685057694676544565325331703978544532544b74344d7572716f684a34313244334b6f6f574b436e4367724b637743324a3577576a626e355435335068646b6b6f57454e534a39546537544b7836366f4a4a34313244334b6f6f575239616f675948786b357a38636b624c4c56346e426f7a777a747871664a7050586671336d4a7232796452474a34313244334b6f6f5744695444635565675637776b313133473366476a69616259756f54436f3157726f6f53656741343263556f544a34313244334b6f6f574e64687072586b5472665370354d5071736270467a70364167394a53787358694341434442676454424c656652820300050e416c74424e2d3132382047e282810e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c216a812d7616e47f0bd38fa4863f48fbcda6a38af4c58d2233dfa7cf79620947042d09f923e0a2f7a2270391e8b058d8bdb8f79fe082b7b627f025651c7290382fdff97c3181d15d162c146ce87ff752499d2acc2b26011439a12e29571a6f1e1defb1751c3be4258c493984fd9f0f6b4a26c539870b5f15bfed3d8ffac92499eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b127a112970e1adf615f823b2b2180754c2f0ee01f1b389e56df55ca09702cd0401b66ff71779d2dd67222503a85ab921b28c329cc1832800b192d0b0247c0776e1b9653dc00df48daa6364287c84c0382f5165e7269fef06d10bc67c1bba252305d1af0dc7bb0fe92558eb4c5f38c23163dee1cfb34a72020669dbdfe337c16f3307472616e736c61746f722066726f6d20416c74424e2d3132382047e2828120746f20416c74424e2d3132382047e282825880ade2046080c8afa0256880c8afa0257080ade204788094ebdc0382019e010a205034214e0bd4373f38e162cf9fc9133e2f3b71441faa4c3d1ac01c1877f1cd2712200e03e975b996f911abba2b79d2596c2150bc94510963c40a1137a03df6edacdb1a107dee1cdb894163813bb3da604c9c133c1a10bb33302eeafbd55d352e35dcc5d2b3311a10d2c658b6b93d74a02d467849b6fe75251a10fea5308cc1fea69e7246eafe7ca8a3a51a1048efe1ad873b6f025ac0243bdef715f8000000000000000000000000000000000000000000000000000000000000"),
			}, nil)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, dkgAddress, 10).
			Return(nil, errors.New("rpc error"))

		c := &coordinator{
			lp:                 lp,
			topics:             tp,
			lggr:               logger.TestLogger(t),
			beaconAddress:      beaconAddress,
			coordinatorAddress: coordinatorAddress,
			dkgAddress:         dkgAddress,
			finalityDepth:      10,
			evmClient:          evmClient,
		}
		_, _, err := c.DKGVRFCommittees(testutils.Context(t))
		assert.Error(t, err)
	})
}

func TestCoordinator_ProvingKeyHash(t *testing.T) {
	t.Parallel()

	t.Run("valid output", func(t *testing.T) {
		h := crypto.Keccak256Hash([]byte("hello world"))
		var expected [32]byte
		copy(expected[:], h.Bytes())
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("SProvingKeyHash", mock.Anything).
			Return(expected, nil)
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		provingKeyHash, err := c.ProvingKeyHash(testutils.Context(t))
		assert.NoError(t, err)
		assert.ElementsMatch(t, expected[:], provingKeyHash[:])
	})

	t.Run("invalid output", func(t *testing.T) {
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("SProvingKeyHash", mock.Anything).
			Return([32]byte{}, errors.New("rpc error"))
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		_, err := c.ProvingKeyHash(testutils.Context(t))
		assert.Error(t, err)
	})
}

func TestCoordinator_ReportBlocks(t *testing.T) {
	lggr := logger.TestLogger(t)
	proofG1X := big.NewInt(1)
	proofG1Y := big.NewInt(2)
	evmClient := evmclimocks.NewClient(t)
	evmClient.On("ConfiguredChainID").Return(big.NewInt(1))
	t.Run("happy path, beacon requests", func(t *testing.T) {
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)

		latestHeadNumber := uint64(200)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		tp := newTopics()

		lookbackBlocks := uint64(5)
		lp := getLogPoller(t, []uint64{195}, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191, 0, coordinatorAddress),
			newRandomnessRequestedLog(t, 3, 195, 192, 1, coordinatorAddress),
			newRandomnessRequestedLog(t, 3, 195, 193, 2, coordinatorAddress),
		}, nil).Once()

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 1)
		assert.Len(t, callbacks, 0)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, callback requests", func(t *testing.T) {
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)

		latestHeadNumber := uint64(200)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		tp := newTopics()

		lookbackBlocks := uint64(5)
		lp := getLogPoller(t, []uint64{195}, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3, 1000, coordinatorAddress),
		}, nil).Once()

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 1)
		assert.Len(t, callbacks, 3)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, beacon requests, beacon fulfillments", func(t *testing.T) {
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)

		latestHeadNumber := uint64(200)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		tp := newTopics()

		lookbackBlocks := uint64(5)
		lp := getLogPoller(t, []uint64{195}, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191, 0, coordinatorAddress),
			newRandomnessRequestedLog(t, 3, 195, 192, 1, coordinatorAddress),
			newRandomnessRequestedLog(t, 3, 195, 193, 2, coordinatorAddress),
			newOutputsServedLog(t, []vrf_coordinator.VRFBeaconTypesOutputServed{
				{
					Height:            195,
					ConfirmationDelay: big.NewInt(3),
					ProofG1X:          proofG1X,
					ProofG1Y:          proofG1Y,
				},
			}, coordinatorAddress),
		}, nil).Once()

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 0)
		assert.Len(t, callbacks, 0)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, callback requests, callback fulfillments", func(t *testing.T) {
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)

		latestHeadNumber := uint64(200)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		tp := newTopics()

		lookbackBlocks := uint64(5)
		lp := getLogPoller(t, []uint64{195}, latestHeadNumber, true, true, lookbackBlocks)
		// Both RandomWordsFulfilled and NewTransmission events are emitted
		// when a VRF fulfillment happens on chain.
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3, 1000, coordinatorAddress),
			// Regardless of success or failure, if the fulfillment has been tried once do not report again.
			newRandomWordsFulfilledLog(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, []byte{1, 0, 0}, coordinatorAddress),
			newOutputsServedLog(t, []vrf_coordinator.VRFBeaconTypesOutputServed{
				{
					Height:            195,
					ConfirmationDelay: big.NewInt(3),
					ProofG1X:          proofG1X,
					ProofG1Y:          proofG1Y,
				},
			}, coordinatorAddress),
		}, nil).Once()

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 0)
		assert.Len(t, callbacks, 0)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, only beacon fulfillment", func(t *testing.T) {
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)

		latestHeadNumber := uint64(200)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		tp := newTopics()

		lookbackBlocks := uint64(5)
		lp := getLogPoller(t, []uint64{}, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{newOutputsServedLog(t, []vrf_coordinator.VRFBeaconTypesOutputServed{
			{
				Height:            195,
				ConfirmationDelay: big.NewInt(3),
				ProofG1X:          proofG1X,
				ProofG1Y:          proofG1Y,
			},
		}, coordinatorAddress)}, nil).Once()

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 0)
		assert.Len(t, callbacks, 0)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, callback requests & callback fulfillments in-flight", func(t *testing.T) {
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)

		latestHeadNumber := uint64(200)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		tp := newTopics()

		lookbackBlocks := uint64(5)
		// Do not include latestHeadNumber in "GetBlocksRange" call for initial "ReportWillBeTransmitted."
		// Do not include recent blockhashes in range either.
		lp := getLogPoller(t, []uint64{195}, latestHeadNumber, false, false /* includeLatestHeadInRange */, 0)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		report := ocr2vrftypes.AbstractReport{
			RecentBlockHeight: 195,
			RecentBlockHash:   common.HexToHash("0x001"),
			Outputs: []ocr2vrftypes.AbstractVRFOutput{
				{
					BlockHeight:       195,
					ConfirmationDelay: 195,
					Callbacks: []ocr2vrftypes.AbstractCostedCallbackRequest{
						{
							RequestID:    1,
							BeaconHeight: 195,
						},
						{
							RequestID:    2,
							BeaconHeight: 195,
						},
						{
							RequestID:    3,
							BeaconHeight: 195,
						},
					},
				},
			},
		}

		err = c.ReportWillBeTransmitted(testutils.Context(t), report)
		require.NoError(t, err)

		// Include latestHeadNumber in "GetBlocksRange" call for "ReportBlocks" call.
		// Include recent blockhashes in range.
		lp = getLogPoller(t, []uint64{195}, latestHeadNumber, true, true /* includeLatestHeadInRange */, lookbackBlocks)
		c.lp = lp
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3, 1000, coordinatorAddress),
			newOutputsServedLog(t, []vrf_coordinator.VRFBeaconTypesOutputServed{
				{
					Height:            195,
					ConfirmationDelay: big.NewInt(3),
					ProofG1X:          proofG1X,
					ProofG1Y:          proofG1Y,
				},
			}, coordinatorAddress),
		}, nil).Once()

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 0)
		assert.Len(t, callbacks, 0)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, blocks requested hits batch gas limit", func(t *testing.T) {
		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		latestHeadNumber := uint64(400)
		lookbackBlocks := uint64(400)
		blockhashLookback := uint64(256)

		tp := newTopics()

		logs := []logpoller.Log{}
		requestedBlocks := []uint64{}

		// Populate 200 request blocks.
		for i := 0; i < 400; i += 2 {
			logs = append(logs, newRandomnessRequestedLog(t, 1, uint64(i), 0, int64(i), coordinatorAddress))
			requestedBlocks = append(requestedBlocks, uint64(i))
		}
		lp := getLogPoller(t, requestedBlocks, latestHeadNumber, true, true, blockhashLookback)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return(logs, nil)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        blockhashLookback,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{1: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)

		// Coordinator should allow 99 blocks, i.e 100 blocks - 1 block's worth of gas
		// for the coordinator overhead.
		assert.NoError(t, err)
		assert.Len(t, blocks, 99)
		assert.Len(t, callbacks, 0)
		assert.Equal(t, uint64(latestHeadNumber-blockhashLookback+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(blockhashLookback))
	})

	t.Run("happy path, last callback hits batch gas limit", func(t *testing.T) {
		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		latestHeadNumber := uint64(200)
		lookbackBlocks := uint64(5)

		tp := newTopics()

		requestedBlocks := []uint64{195}
		lp := getLogPoller(t, requestedBlocks, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191, 0, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1, 2_000_000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2, 2_900_000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3, 1, coordinatorAddress),
		}, nil)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)

		// Should allow the first two callbacks, which add up to 4_950_000 + 50_000 (1 block) = 5_000_000,
		// then reject the last callback for being out of gas.
		assert.NoError(t, err)
		assert.Len(t, blocks, 1)
		assert.Len(t, callbacks, 2)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, sandwiched callbacks hit batch gas limit", func(t *testing.T) {
		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		latestHeadNumber := uint64(200)
		lookbackBlocks := uint64(5)

		tp := newTopics()

		requestedBlocks := []uint64{195}
		lp := getLogPoller(t, requestedBlocks, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191, 0, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1, 10_000_000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3, 10_000_000, coordinatorAddress),
		}, nil)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)

		// Should allow the middle callback, with an acceptable gas allowance, to be processed.
		assert.NoError(t, err)
		assert.Len(t, blocks, 1)
		assert.Len(t, callbacks, 1)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("happy path, sandwiched callbacks with valid callback in next block hit batch gas limit", func(t *testing.T) {
		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		latestHeadNumber := uint64(200)
		lookbackBlocks := uint64(5)

		tp := newTopics()

		requestedBlocks := []uint64{195, 196}
		lp := getLogPoller(t, requestedBlocks, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191, 0, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1, 10_000_000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3, 10_000_000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 196, 194, 4, 1000, coordinatorAddress),
		}, nil)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)

		// Should allow the middle callback, with an acceptable gas allowance, to be processed,
		// then move to the next block and find a suitable callback. Also adds the block 196 for
		// that callback.
		assert.NoError(t, err)
		assert.Len(t, blocks, 2)
		assert.Len(t, callbacks, 2)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("correct blockhashes are retrieved with the maximum lookback", func(t *testing.T) {
		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		latestHeadNumber := uint64(1000)
		lookbackBlocks := uint64(256)

		tp := newTopics()

		requestedBlocks := []uint64{}
		lp := getLogPoller(t, requestedBlocks, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{}, nil)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		_, _, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)

		assert.NoError(t, err)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Equal(t, common.HexToHash(fmt.Sprintf("0x00%d", 1)), recentBlocks[0])
		assert.Equal(t, common.HexToHash(fmt.Sprintf("0x00%d", lookbackBlocks)), recentBlocks[len(recentBlocks)-1])
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("correct blockhashes are retrieved with a capped lookback (close to genesis block)", func(t *testing.T) {
		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		latestHeadNumber := uint64(100)
		lookbackBlocks := uint64(100)

		tp := newTopics()

		requestedBlocks := []uint64{}
		lp := getLogPoller(t, requestedBlocks, latestHeadNumber, true, true, lookbackBlocks)
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{}, nil)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		_, _, recentHeightStart, recentBlocks, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)

		assert.NoError(t, err)
		assert.Equal(t, uint64(latestHeadNumber-lookbackBlocks+1), recentHeightStart)
		assert.Equal(t, common.HexToHash(fmt.Sprintf("0x00%d", 1)), recentBlocks[0])
		assert.Equal(t, common.HexToHash(fmt.Sprintf("0x00%d", lookbackBlocks)), recentBlocks[len(recentBlocks)-1])
		assert.Len(t, recentBlocks, int(lookbackBlocks))
	})

	t.Run("logpoller GetBlocks returns error", func(tt *testing.T) {
		coordinatorAddress := newAddress(t)
		beaconAddress := newAddress(t)
		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		require.NoError(t, err)

		latestHeadNumber := uint64(200)
		lookbackBlocks := uint64(5)

		tp := newTopics()

		requestedBlocks := []uint64{195, 196}
		lp := lp_mocks.NewLogPoller(t)
		lp.On("LatestBlock", mock.Anything).
			Return(int64(latestHeadNumber), nil)

		lp.On("GetBlocksRange", mock.Anything, append(requestedBlocks, uint64(latestHeadNumber-lookbackBlocks+1), uint64(latestHeadNumber)), mock.Anything).
			Return(nil, errors.New("GetBlocks error"))
		lp.On(
			"LogsWithSigs",
			int64(latestHeadNumber-lookbackBlocks),
			int64(latestHeadNumber),
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.outputsServedTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191, 0, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1, 10_000_000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2, 1000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3, 10_000_000, coordinatorAddress),
			newRandomnessFulfillmentRequestedLog(t, 3, 196, 194, 4, 1000, coordinatorAddress),
		}, nil)

		c := &coordinator{
			onchainRouter:            onchainRouter,
			beaconAddress:            beaconAddress,
			coordinatorAddress:       coordinatorAddress,
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			topics:                   tp,
			evmClient:                evmClient,
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			blockhashLookback:        lookbackBlocks,
		}

		blocks, callbacks, _, _, err := c.ReportBlocks(
			testutils.Context(t),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)

		assert.Error(tt, err)
		assert.EqualError(tt, errors.Cause(err), "GetBlocks error")
		assert.Nil(tt, blocks)
		assert.Nil(tt, callbacks)
	})
}

func TestCoordinator_ReportWillBeTransmitted(t *testing.T) {
	evmClient := evmclimocks.NewClient(t)
	evmClient.On("ConfiguredChainID").Return(big.NewInt(1))
	t.Run("happy path", func(t *testing.T) {
		lookbackBlocks := uint64(0)
		lp := getLogPoller(t, []uint64{199}, 200, false, false, 0)
		c := &coordinator{
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			evmClient:                evmClient,
		}
		assert.NoError(t, c.ReportWillBeTransmitted(testutils.Context(t), ocr2vrftypes.AbstractReport{
			RecentBlockHeight: 199,
			RecentBlockHash:   common.HexToHash("0x001"),
		}))
	})

	t.Run("re-org", func(t *testing.T) {
		lookbackBlocks := uint64(0)
		lp := getLogPoller(t, []uint64{199}, 200, false, false, 0)
		c := &coordinator{
			lp:                       lp,
			lggr:                     logger.TestLogger(t),
			toBeTransmittedBlocks:    NewBlockCache[blockInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			toBeTransmittedCallbacks: NewBlockCache[callbackInReport](time.Duration(int64(lookbackBlocks) * int64(time.Second))),
			coordinatorConfig:        newCoordinatorConfig(lookbackBlocks),
			evmClient:                evmClient,
		}
		assert.Error(t, c.ReportWillBeTransmitted(testutils.Context(t), ocr2vrftypes.AbstractReport{
			RecentBlockHeight: 199,
			RecentBlockHash:   common.HexToHash("0x009"),
		}))
	})
}

func TestCoordinator_MarshalUnmarshal(t *testing.T) {
	t.Parallel()
	proofG1X := big.NewInt(1)
	proofG1Y := big.NewInt(2)
	lggr := logger.TestLogger(t)
	evmClient := evmclimocks.NewClient(t)

	coordinatorAddress := newAddress(t)
	beaconAddress := newAddress(t)
	vrfBeaconCoordinator, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
	require.NoError(t, err)

	lg := newRandomnessRequestedLog(t, 3, 1500, 1450, 1, coordinatorAddress)
	rrIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	rr, ok := rrIface.(*vrf_coordinator.VRFCoordinatorRandomnessRequested)
	require.True(t, ok)
	assert.Equal(t, uint64(1500), rr.NextBeaconOutputHeight)
	assert.Equal(t, int64(3), rr.ConfDelay.Int64())

	lg = newRandomnessFulfillmentRequestedLog(t, 3, 1500, 1450, 1, 1000, coordinatorAddress)
	rfrIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	rfr, ok := rfrIface.(*vrf_coordinator.VRFCoordinatorRandomnessFulfillmentRequested)
	require.True(t, ok)
	assert.Equal(t, uint64(1500), rfr.NextBeaconOutputHeight)
	assert.Equal(t, int64(3), rfr.ConfDelay.Int64())
	assert.Equal(t, int64(1), rfr.RequestID.Int64())

	configDigest := common.BigToHash(big.NewInt(10))
	lg = newNewTransmissionLog(t, beaconAddress, configDigest)
	ntIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	nt, ok := ntIface.(*vrf_beacon.VRFBeaconNewTransmission)
	require.True(t, ok)
	assert.True(t, bytes.Equal(nt.ConfigDigest[:], configDigest[:]))
	assert.Equal(t, 0, nt.JuelsPerFeeCoin.Cmp(big.NewInt(1_000)))
	assert.Equal(t, 0, nt.EpochAndRound.Cmp(big.NewInt(1)))

	lg = newRandomWordsFulfilledLog(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, []byte{1, 1, 1}, coordinatorAddress)
	rwfIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	rwf, ok := rwfIface.(*vrf_coordinator.VRFCoordinatorRandomWordsFulfilled)
	require.True(t, ok)
	assert.Equal(t, []int64{1, 2, 3}, []int64{rwf.RequestIDs[0].Int64(), rwf.RequestIDs[1].Int64(), rwf.RequestIDs[2].Int64()})
	assert.Equal(t, []byte{1, 1, 1}, rwf.SuccessfulFulfillment)

	lg = newOutputsServedLog(t, []vrf_coordinator.VRFBeaconTypesOutputServed{
		{
			Height:            1500,
			ConfirmationDelay: big.NewInt(3),
			ProofG1X:          proofG1X,
			ProofG1Y:          proofG1Y,
		},
		{
			Height:            1505,
			ConfirmationDelay: big.NewInt(4),
			ProofG1X:          proofG1X,
			ProofG1Y:          proofG1Y,
		},
	}, coordinatorAddress)

	osIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	os, ok := osIface.(*vrf_coordinator.VRFCoordinatorOutputsServed)
	require.True(t, ok)
	assert.Equal(t, uint64(1500), os.OutputsServed[0].Height)
	assert.Equal(t, uint64(1505), os.OutputsServed[1].Height)
	assert.Equal(t, int64(3), os.OutputsServed[0].ConfirmationDelay.Int64())
	assert.Equal(t, int64(4), os.OutputsServed[1].ConfirmationDelay.Int64())
}

func TestCoordinator_ReportIsOnchain(t *testing.T) {
	evmClient := evmclimocks.NewClient(t)
	evmClient.On("ConfiguredChainID").Return(big.NewInt(1))

	t.Run("report is on-chain", func(t *testing.T) {
		tp := newTopics()
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)
		lggr := logger.TestLogger(t)

		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		assert.NoError(t, err)

		epoch := uint32(20)
		round := uint8(3)
		epochAndRound := toEpochAndRoundUint40(epoch, round)
		enrTopic := common.BytesToHash(common.LeftPadBytes(epochAndRound.Bytes(), 32))
		lp := lp_mocks.NewLogPoller(t)
		configDigest := common.BigToHash(big.NewInt(1337))
		log := newNewTransmissionLog(t, beaconAddress, configDigest)
		log.BlockNumber = 195
		lp.On("IndexedLogs", tp.newTransmissionTopic, beaconAddress, 2, []common.Hash{
			enrTopic,
		}, 1, mock.Anything).Return([]logpoller.Log{log}, nil)

		c := &coordinator{
			lp:            lp,
			onchainRouter: onchainRouter,
			lggr:          logger.TestLogger(t),
			beaconAddress: beaconAddress,
			topics:        tp,
			evmClient:     evmClient,
		}

		present, err := c.ReportIsOnchain(testutils.Context(t), epoch, round, configDigest)
		assert.NoError(t, err)
		assert.True(t, present)
	})

	t.Run("report is on-chain for old config digest", func(t *testing.T) {
		tp := newTopics()
		beaconAddress := newAddress(t)
		coordinatorAddress := newAddress(t)
		lggr := logger.TestLogger(t)

		onchainRouter, err := newRouter(lggr, beaconAddress, coordinatorAddress, evmClient)
		assert.NoError(t, err)

		epoch := uint32(20)
		round := uint8(3)
		epochAndRound := toEpochAndRoundUint40(epoch, round)
		enrTopic := common.BytesToHash(common.LeftPadBytes(epochAndRound.Bytes(), 32))
		lp := lp_mocks.NewLogPoller(t)
		oldConfigDigest := common.BigToHash(big.NewInt(1337))
		newConfigDigest := common.BigToHash(big.NewInt(8888))
		log := newNewTransmissionLog(t, beaconAddress, oldConfigDigest)
		log.BlockNumber = 195
		lp.On("IndexedLogs", tp.newTransmissionTopic, beaconAddress, 2, []common.Hash{
			enrTopic,
		}, 1, mock.Anything).Return([]logpoller.Log{log}, nil)

		c := &coordinator{
			lp:            lp,
			onchainRouter: onchainRouter,
			lggr:          logger.TestLogger(t),
			beaconAddress: beaconAddress,
			topics:        tp,
			evmClient:     evmClient,
		}

		present, err := c.ReportIsOnchain(testutils.Context(t), epoch, round, newConfigDigest)
		assert.NoError(t, err)
		assert.False(t, present)
	})

	t.Run("report is not on-chain", func(t *testing.T) {
		tp := newTopics()
		beaconAddress := newAddress(t)

		epoch := uint32(20)
		round := uint8(3)
		epochAndRound := toEpochAndRoundUint40(epoch, round)
		enrTopic := common.BytesToHash(common.LeftPadBytes(epochAndRound.Bytes(), 32))
		lp := lp_mocks.NewLogPoller(t)
		lp.On("IndexedLogs", tp.newTransmissionTopic, beaconAddress, 2, []common.Hash{
			enrTopic,
		}, 1, mock.Anything).Return([]logpoller.Log{}, nil)

		c := &coordinator{
			lp:            lp,
			lggr:          logger.TestLogger(t),
			beaconAddress: beaconAddress,
			topics:        tp,
			evmClient:     evmClient,
		}

		configDigest := common.BigToHash(big.NewInt(0))
		present, err := c.ReportIsOnchain(testutils.Context(t), epoch, round, configDigest)
		assert.NoError(t, err)
		assert.False(t, present)
	})

}

func TestCoordinator_ConfirmationDelays(t *testing.T) {
	t.Parallel()

	t.Run("valid output", func(t *testing.T) {
		expected := [8]uint32{1, 2, 3, 4, 5, 6, 7, 8}
		ret := [8]*big.Int{}
		for i, delay := range expected {
			ret[i] = big.NewInt(int64(delay))
		}
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("GetConfirmationDelays", mock.Anything).
			Return(ret, nil)
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		confDelays, err := c.ConfirmationDelays(testutils.Context(t))
		assert.NoError(t, err)
		assert.Equal(t, expected[:], confDelays[:])
	})

	t.Run("invalid output", func(t *testing.T) {
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("GetConfirmationDelays", mock.Anything).
			Return([8]*big.Int{}, errors.New("rpc error"))
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		_, err := c.ConfirmationDelays(testutils.Context(t))
		assert.Error(t, err)
	})
}

func TestCoordinator_getBlockCacheKey(t *testing.T) {
	t.Parallel()

	t.Run("calculates key correctly", func(t *testing.T) {
		hash := getBlockCacheKey(1, 11)
		assert.Equal(
			t,
			common.HexToHash("0x000000000000000000000000000000000000000000000001000000000000000b"),
			hash,
		)
	})
}

func TestCoordinator_KeyID(t *testing.T) {
	t.Parallel()

	t.Run("valid output", func(t *testing.T) {
		var keyIDBytes [32]byte
		keyIDBytes[0] = 1
		expected := dkg.KeyID(keyIDBytes)
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("SKeyID", mock.Anything).
			Return(keyIDBytes, nil)
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		keyID, err := c.KeyID(testutils.Context(t))
		assert.NoError(t, err)
		assert.Equal(t, expected[:], keyID[:])
	})

	t.Run("invalid output", func(t *testing.T) {
		var emptyBytes [32]byte
		onchainRouter := mocks.NewVRFBeaconCoordinator(t)
		onchainRouter.
			On("SKeyID", mock.Anything).
			Return(emptyBytes, errors.New("rpc error"))
		c := &coordinator{
			onchainRouter: onchainRouter,
		}
		_, err := c.KeyID(testutils.Context(t))
		assert.Error(t, err)
	})
}

func TestTopics_DKGConfigSet_VRFConfigSet(t *testing.T) {
	dkgConfigSetTopic := dkg_wrapper.DKGConfigSet{}.Topic()
	vrfConfigSetTopic := vrf_beacon.VRFBeaconConfigSet{}.Topic()
	assert.Equal(t, dkgConfigSetTopic, vrfConfigSetTopic, "config set topics of vrf and dkg must be equal")
}

func Test_UpdateConfiguration(t *testing.T) {
	t.Parallel()

	t.Run("valid binary", func(t *testing.T) {
		c := &coordinator{coordinatorConfig: newCoordinatorConfig(10), lggr: logger.TestLogger(t)}
		cacheEvictionWindowSeconds := int64(60)
		cacheEvictionWindow := time.Duration(cacheEvictionWindowSeconds * int64(time.Second))
		c.toBeTransmittedBlocks = NewBlockCache[blockInReport](cacheEvictionWindow)
		c.toBeTransmittedCallbacks = NewBlockCache[callbackInReport](cacheEvictionWindow)

		newCoordinatorConfig := &ocr2vrftypes.CoordinatorConfig{
			CacheEvictionWindowSeconds: 30,
			BatchGasLimit:              1_000_000,
			CoordinatorOverhead:        10_000,
			CallbackOverhead:           10_000,
			BlockGasOverhead:           10_000,
			LookbackBlocks:             1_000,
		}

		require.Equal(t, cacheEvictionWindow, c.toBeTransmittedBlocks.evictionWindow)
		require.Equal(t, cacheEvictionWindow, c.toBeTransmittedCallbacks.evictionWindow)

		expectedConfigDigest := ocr2Types.ConfigDigest(common.HexToHash("asd"))
		expectedOracleID := commontypes.OracleID(3)
		err := c.UpdateConfiguration(ocr2vrf.OffchainConfig(newCoordinatorConfig), expectedConfigDigest, expectedOracleID)
		newCacheEvictionWindow := time.Duration(newCoordinatorConfig.CacheEvictionWindowSeconds * int64(time.Second))
		require.NoError(t, err)
		require.Equal(t, newCoordinatorConfig.CacheEvictionWindowSeconds, c.coordinatorConfig.CacheEvictionWindowSeconds)
		require.Equal(t, newCoordinatorConfig.BatchGasLimit, c.coordinatorConfig.BatchGasLimit)
		require.Equal(t, newCoordinatorConfig.CoordinatorOverhead, c.coordinatorConfig.CoordinatorOverhead)
		require.Equal(t, newCoordinatorConfig.CallbackOverhead, c.coordinatorConfig.CallbackOverhead)
		require.Equal(t, newCoordinatorConfig.BlockGasOverhead, c.coordinatorConfig.BlockGasOverhead)
		require.Equal(t, newCoordinatorConfig.LookbackBlocks, c.coordinatorConfig.LookbackBlocks)
		require.Equal(t, newCacheEvictionWindow, c.toBeTransmittedBlocks.evictionWindow)
		require.Equal(t, newCacheEvictionWindow, c.toBeTransmittedCallbacks.evictionWindow)
		require.Equal(t, expectedConfigDigest, c.configDigest)
		require.Equal(t, expectedOracleID, c.oracleID)
	})

	t.Run("invalid binary", func(t *testing.T) {
		c := &coordinator{coordinatorConfig: newCoordinatorConfig(10), lggr: logger.TestLogger(t)}

		err := c.UpdateConfiguration([]byte{123}, ocr2Types.ConfigDigest{}, commontypes.OracleID(0))
		require.Error(t, err)
	})
}

func newCoordinatorConfig(lookbackBlocks uint64) *ocr2vrftypes.CoordinatorConfig {
	return &ocr2vrftypes.CoordinatorConfig{
		CacheEvictionWindowSeconds: 60,
		BatchGasLimit:              5_000_000,
		CoordinatorOverhead:        50_000,
		CallbackOverhead:           50_000,
		BlockGasOverhead:           50_000,
		LookbackBlocks:             lookbackBlocks,
	}
}

func newRandomnessRequestedLog(
	t *testing.T,
	confDelay int64,
	nextBeaconOutputHeight uint64,
	requestBlock uint64,
	requestID int64,
	coordinatorAddress common.Address,
) logpoller.Log {
	//event RandomnessRequested(
	//    RequestID indexed requestID,
	//    address indexed requester,
	//    uint64 nextBeaconOutputHeight,
	//    ConfirmationDelay confDelay,
	//    uint64 subID,
	//    uint16 numWords
	//);
	e := vrf_coordinator.VRFCoordinatorRandomnessRequested{
		RequestID:              big.NewInt(requestID),
		Requester:              common.HexToAddress("0x1234567890"),
		ConfDelay:              big.NewInt(confDelay),
		NextBeaconOutputHeight: nextBeaconOutputHeight,
		NumWords:               1,
		SubID:                  big.NewInt(1),
		Raw: types.Log{
			BlockNumber: requestBlock,
		},
	}
	var unindexed abi.Arguments
	for _, a := range vrfCoordinatorABI.Events[randomnessRequestedEvent].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(e.NextBeaconOutputHeight, e.ConfDelay, e.SubID, e.NumWords)
	require.NoError(t, err)

	requestIDType, err := abi.NewType("uint64", "", nil)
	require.NoError(t, err)

	requesterType, err := abi.NewType("address", "", nil)
	require.NoError(t, err)

	requestIDArg := abi.Arguments{abi.Argument{
		Name:    "requestID",
		Type:    requestIDType,
		Indexed: true,
	}}
	requesterArg := abi.Arguments{abi.Argument{
		Name:    "requester",
		Type:    requesterType,
		Indexed: true,
	}}

	topic1, err := requestIDArg.Pack(e.RequestID.Uint64())
	require.NoError(t, err)
	topic2, err := requesterArg.Pack(e.Requester)
	require.NoError(t, err)

	topic0 := vrfCoordinatorABI.Events[randomnessRequestedEvent].ID
	lg := logpoller.Log{
		Address: coordinatorAddress,
		Data:    nonIndexedData,
		Topics: [][]byte{
			// first topic is the event signature
			topic0.Bytes(),
			// second topic is requestID since it's indexed
			topic1,
			// third topic is requester since it's indexed
			topic2,
		},
		BlockNumber: int64(requestBlock),
		EventSig:    topic0,
	}
	return lg
}

func newRandomnessFulfillmentRequestedLog(
	t *testing.T,
	confDelay int64,
	nextBeaconOutputHeight uint64,
	requestBlock uint64,
	requestID int64,
	gasAllowance uint32,
	coordinatorAddress common.Address,
) logpoller.Log {
	//event RandomnessFulfillmentRequested(
	//    RequestID indexed requestID,
	//    address indexed requester,
	//    uint64 nextBeaconOutputHeight,
	//    ConfirmationDelay confDelay,
	//    uint64 subID,
	//    uint16 numWords,
	//    uint32 gasAllowance,
	//    uint256 gasPrice,
	//    uint256 weiPerUnitLink,
	//    bytes arguments
	//);
	e := vrf_coordinator.VRFCoordinatorRandomnessFulfillmentRequested{
		ConfDelay:              big.NewInt(confDelay),
		NextBeaconOutputHeight: nextBeaconOutputHeight,
		RequestID:              big.NewInt(1),
		NumWords:               1,
		GasAllowance:           gasAllowance,
		GasPrice:               big.NewInt(0),
		WeiPerUnitLink:         big.NewInt(0),
		SubID:                  big.NewInt(1),
		Requester:              common.HexToAddress("0x1234567890"),
		Raw: types.Log{
			BlockNumber: requestBlock,
		},
	}
	var unindexed abi.Arguments
	for _, a := range vrfCoordinatorABI.Events[randomnessFulfillmentRequestedEvent].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(e.NextBeaconOutputHeight, e.ConfDelay, e.SubID, e.NumWords,
		e.GasAllowance, e.GasPrice, e.WeiPerUnitLink, e.Arguments)
	require.NoError(t, err)

	requestIDType, err := abi.NewType("uint64", "", nil)
	require.NoError(t, err)

	requesterType, err := abi.NewType("address", "", nil)
	require.NoError(t, err)

	requestIDArg := abi.Arguments{abi.Argument{
		Name:    "requestID",
		Type:    requestIDType,
		Indexed: true,
	}}
	requesterArg := abi.Arguments{abi.Argument{
		Name:    "requester",
		Type:    requesterType,
		Indexed: true,
	}}

	topic0 := vrfCoordinatorABI.Events[randomnessFulfillmentRequestedEvent].ID
	topic1, err := requestIDArg.Pack(e.RequestID.Uint64())
	require.NoError(t, err)
	topic2, err := requesterArg.Pack(e.Requester)
	require.NoError(t, err)
	return logpoller.Log{
		Address:  coordinatorAddress,
		Data:     nonIndexedData,
		EventSig: topic0,
		Topics: [][]byte{
			topic0.Bytes(),
			topic1,
			topic2,
		},
		BlockNumber: int64(requestBlock),
	}
}

func newRandomWordsFulfilledLog(
	t *testing.T,
	requestIDs []*big.Int,
	successfulFulfillment []byte,
	coordinatorAddress common.Address,
) logpoller.Log {
	//event RandomWordsFulfilled(
	//  RequestID[] requestIDs,
	//  bytes successfulFulfillment,
	//  bytes[] truncatedErrorData
	//);
	e := vrf_coordinator.VRFCoordinatorRandomWordsFulfilled{
		RequestIDs:            requestIDs,
		SuccessfulFulfillment: successfulFulfillment,
	}
	packed, err := vrfCoordinatorABI.Events[randomWordsFulfilledEvent].Inputs.Pack(
		e.RequestIDs, e.SuccessfulFulfillment, e.TruncatedErrorData)
	require.NoError(t, err)
	topic0 := vrfCoordinatorABI.Events[randomWordsFulfilledEvent].ID
	return logpoller.Log{
		Address:  coordinatorAddress,
		Data:     packed,
		EventSig: topic0,
		Topics:   [][]byte{topic0.Bytes()},
	}
}

func newOutputsServedLog(
	t *testing.T,
	outputsServed []vrf_coordinator.VRFBeaconTypesOutputServed,
	beaconAddress common.Address,
) logpoller.Log {
	// event OutputsServed(
	//     uint64 recentBlockHeight,
	//     address transmitter,
	//     uint192 juelsPerFeeCoin,
	//     OutputServed[] outputsServed
	// );
	e := vrf_coordinator.VRFCoordinatorOutputsServed{
		RecentBlockHeight: 0,
		// AggregatorRoundId: 1,
		OutputsServed:      outputsServed,
		JuelsPerFeeCoin:    big.NewInt(0),
		ReasonableGasPrice: 0,
		// EpochAndRound:     big.NewInt(1),
		// ConfigDigest:      crypto.Keccak256Hash([]byte("hello world")),
		Transmitter: newAddress(t),
	}
	var unindexed abi.Arguments
	for _, a := range vrfCoordinatorABI.Events[outputsServedEvent].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(e.RecentBlockHeight, e.Transmitter, e.JuelsPerFeeCoin, e.ReasonableGasPrice, e.OutputsServed)
	require.NoError(t, err)

	topic0 := vrfCoordinatorABI.Events[outputsServedEvent].ID
	return logpoller.Log{
		Address: beaconAddress,
		Data:    nonIndexedData,
		Topics: [][]byte{
			topic0.Bytes(),
		},
		EventSig: topic0,
	}
}

func newNewTransmissionLog(
	t *testing.T,
	beaconAddress common.Address,
	configDigest [32]byte,
) logpoller.Log {
	// event NewTransmission(
	//     uint32 indexed aggregatorRoundId,
	//     uint40 indexed epochAndRound,
	//     address transmitter,
	//     uint192 juelsPerFeeCoin,
	//     bytes32 configDigest
	// );
	e := vrf_beacon.VRFBeaconNewTransmission{
		AggregatorRoundId:  1,
		JuelsPerFeeCoin:    big.NewInt(1_000),
		ReasonableGasPrice: 1_000,
		EpochAndRound:      big.NewInt(1),
		ConfigDigest:       configDigest,
		Transmitter:        newAddress(t),
	}
	var unindexed abi.Arguments
	for _, a := range vrfBeaconABI.Events[newTransmissionEvent].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(
		e.Transmitter, e.JuelsPerFeeCoin, e.ReasonableGasPrice, e.ConfigDigest)
	require.NoError(t, err)

	// aggregatorRoundId is indexed
	aggregatorRoundIDType, err := abi.NewType("uint32", "", nil)
	require.NoError(t, err)
	indexedArgs := abi.Arguments{
		{
			Name: "aggregatorRoundId",
			Type: aggregatorRoundIDType,
		},
	}
	aggregatorPacked, err := indexedArgs.Pack(e.AggregatorRoundId)
	require.NoError(t, err)

	// epochAndRound is indexed
	epochAndRoundType, err := abi.NewType("uint40", "", nil)
	require.NoError(t, err)
	indexedArgs = abi.Arguments{
		{
			Name: "epochAndRound",
			Type: epochAndRoundType,
		},
	}
	epochAndRoundPacked, err := indexedArgs.Pack(e.EpochAndRound)
	require.NoError(t, err)

	topic0 := vrfBeaconABI.Events[newTransmissionEvent].ID
	return logpoller.Log{
		Address: beaconAddress,
		Data:    nonIndexedData,
		Topics: [][]byte{
			topic0.Bytes(),
			aggregatorPacked,
			epochAndRoundPacked,
		},
		EventSig: topic0,
	}
}

func newAddress(t *testing.T) common.Address {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	return common.HexToAddress(hexutil.Encode(b))
}

func getLogPoller(
	t *testing.T,
	requestedBlocks []uint64,
	latestHeadNumber uint64,
	needsLatestBlock bool,
	includeLatestHeadInRange bool,
	blockhashLookback uint64,
) *lp_mocks.LogPoller {
	lp := lp_mocks.NewLogPoller(t)
	if needsLatestBlock {
		lp.On("LatestBlock", mock.Anything).
			Return(int64(latestHeadNumber), nil)
	}
	var logPollerBlocks []logpoller.LogPollerBlock

	// If provided, ajust the blockhash range such that it starts at the most recent head.
	if includeLatestHeadInRange {
		requestedBlocks = append(requestedBlocks, latestHeadNumber)
	}

	// If provided, adjust the blockhash range such that it includes all recent available blockhashes.
	if blockhashLookback != 0 {
		requestedBlocks = append(requestedBlocks, latestHeadNumber-blockhashLookback+1)
	}

	// Sort the blocks to match the coordinator's calls.
	sort.Slice(requestedBlocks, func(a, b int) bool {
		return requestedBlocks[a] < requestedBlocks[b]
	})

	// Fill range of blocks based on requestedBlocks
	// example: requestedBlocks [195, 197] -> [{BlockNumber: 195, BlockHash: 0x001}, {BlockNumber: 196, BlockHash: 0x002}, {BlockNumber: 197, BlockHash: 0x003}]
	minRequestedBlock := mathutil.Min(requestedBlocks[0], requestedBlocks[1:]...)
	maxRequestedBlock := mathutil.Max(requestedBlocks[0], requestedBlocks[1:]...)
	for i := minRequestedBlock; i <= maxRequestedBlock; i++ {
		logPollerBlocks = append(logPollerBlocks, logpoller.LogPollerBlock{
			BlockNumber: int64(i),
			BlockHash:   common.HexToHash(fmt.Sprintf("0x00%d", i-minRequestedBlock+1)),
		})
	}

	lp.On("GetBlocksRange", mock.Anything, requestedBlocks, mock.Anything).
		Return(logPollerBlocks, nil)

	return lp
}

func TestFilterNamesFromSpec(t *testing.T) {
	beaconAddress := newAddress(t)
	coordinatorAddress := newAddress(t)
	dkgAddress := newAddress(t)

	spec := &job.OCR2OracleSpec{
		ContractID: beaconAddress.String(),
		PluginType: job.OCR2VRF,
		PluginConfig: job.JSONConfig{
			"VRFCoordinatorAddress": coordinatorAddress.String(),
			"DKGContractAddress":    dkgAddress.String(),
		},
	}

	names, err := FilterNamesFromSpec(spec)
	require.NoError(t, err)

	assert.Len(t, names, 1)
	assert.Equal(t, logpoller.FilterName("VRF Coordinator", beaconAddress, coordinatorAddress, dkgAddress), names[0])

	spec = &job.OCR2OracleSpec{
		PluginType:   job.OCR2VRF,
		ContractID:   beaconAddress.String(),
		PluginConfig: nil, // missing coordinator & dkg addresses
	}
	_, err = FilterNamesFromSpec(spec)
	require.ErrorContains(t, err, "is not a valid EIP55 formatted address")
}
