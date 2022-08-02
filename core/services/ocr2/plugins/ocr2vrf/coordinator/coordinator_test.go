package coordinator

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/chains/evm/logpoller"
	lp_mocks "github.com/smartcontractkit/chainlink/core/chains/evm/logpoller/mocks"
	evm_mocks "github.com/smartcontractkit/chainlink/core/chains/evm/mocks"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	dkg_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/ocr2vrf/generated/dkg"
	vrf_wrapper "github.com/smartcontractkit/chainlink/core/internal/gethwrappers/ocr2vrf/generated/vrf_beacon_coordinator"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ocr2vrf/coordinator/mocks"
)

func TestCoordinator_BeaconPeriod(t *testing.T) {
	t.Run("valid output", func(t *testing.T) {
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("IBeaconPeriodBlocks", mock.Anything).
			Return(big.NewInt(10), nil)
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		period, err := c.BeaconPeriod(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, uint16(10), period)
	})

	t.Run("invalid output", func(t *testing.T) {
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("IBeaconPeriodBlocks", mock.Anything).
			Return(nil, errors.New("rpc error"))
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		_, err := c.BeaconPeriod(context.TODO())
		assert.Error(t, err)
	})
}

func TestCoordinator_DKGVRFCommittees(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		// In this test the DKG and VRF committees have the same signers and
		// transmitters. This may (?) be different in practice.

		lp := &lp_mocks.LogPoller{}
		tp := newTopics()

		coordinatorAddress := newAddress(t)
		dkgAddress := newAddress(t)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, coordinatorAddress, 1).
			Return(&logpoller.Log{
				Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000a6fca200010576e704b4a519484d6239ef17f1f5b4a82e330b0daf827ed4dc2789971b0000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000002e000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000a8cbea12a06869d3ec432ab9682dab6c761d591000000000000000000000000f4f9db7bb1d16b7cdfb18ec68994c26964f5985300000000000000000000000022fb3f90c539457f00d8484438869135e604a65500000000000000000000000033cbcedccb11c9773ad78e214ba342e979255ab30000000000000000000000006ffaa96256fbc1012325cca88c79f725c33eed80000000000000000000000000000000000000000000000000000000000000000500000000000000000000000074103cf8b436465870b26aa9fa2f62ad62b22e3500000000000000000000000038a6cb196f805cc3041f6645a5a6cec27b64430d00000000000000000000000047d7095cfebf8285bdaa421bc8268d0db87d933c000000000000000000000000a8842be973800ff61d80d2d53fa62c3a685380eb0000000000000000000000003750e31321aee8c024751877070e8d5f704ce98700000000000000000000000000000000000000000000000000000000000000206f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf000000000000000000000000000000000000000000000000000000000000004220880d88ee16f1080c8afa0251880c8afa025208090dfc04a288090dfc04a30033a05010101010142206c5ca6f74b532222ac927dd3de235d46a943e372c0563393a33b01dcfd3f371c4220855114d25c2ef5e85fffe4f20a365672d8f2dba3b2ec82333f494168a2039c0442200266e835634db00977cbc1caa4db10e1676c1a4c0fcbc6ba7f09300f0d1831824220980cd91f7a73f20f4b0d51d00cd4e00373dc2beafbb299ca3c609757ab98c8304220eb6d36e2af8922085ff510bbe1eb8932a0e3295ca9f047fef25d90e69c52948f4a34313244334b6f6f574463364b7232644542684b59326b336e685057694676544565325331703978544532544b74344d7572716f684a34313244334b6f6f574b436e4367724b637743324a3577576a626e355435335068646b6b6f57454e534a39546537544b7836366f4a4a34313244334b6f6f575239616f675948786b357a38636b624c4c56346e426f7a777a747871664a7050586671336d4a7232796452474a34313244334b6f6f5744695444635565675637776b313133473366476a69616259756f54436f3157726f6f53656741343263556f544a34313244334b6f6f574e64687072586b5472665370354d5071736270467a70364167394a53787358694341434442676454424c656652820300050e416c74424e2d3132382047e282810e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c216a812d7616e47f0bd38fa4863f48fbcda6a38af4c58d2233dfa7cf79620947042d09f923e0a2f7a2270391e8b058d8bdb8f79fe082b7b627f025651c7290382fdff97c3181d15d162c146ce87ff752499d2acc2b26011439a12e29571a6f1e1defb1751c3be4258c493984fd9f0f6b4a26c539870b5f15bfed3d8ffac92499eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b127a112970e1adf615f823b2b2180754c2f0ee01f1b389e56df55ca09702cd0401b66ff71779d2dd67222503a85ab921b28c329cc1832800b192d0b0247c0776e1b9653dc00df48daa6364287c84c0382f5165e7269fef06d10bc67c1bba252305d1af0dc7bb0fe92558eb4c5f38c23163dee1cfb34a72020669dbdfe337c16f3307472616e736c61746f722066726f6d20416c74424e2d3132382047e2828120746f20416c74424e2d3132382047e282825880ade2046080c8afa0256880c8afa0257080ade204788094ebdc0382019e010a205034214e0bd4373f38e162cf9fc9133e2f3b71441faa4c3d1ac01c1877f1cd2712200e03e975b996f911abba2b79d2596c2150bc94510963c40a1137a03df6edacdb1a107dee1cdb894163813bb3da604c9c133c1a10bb33302eeafbd55d352e35dcc5d2b3311a10d2c658b6b93d74a02d467849b6fe75251a10fea5308cc1fea69e7246eafe7ca8a3a51a1048efe1ad873b6f025ac0243bdef715f8000000000000000000000000000000000000000000000000000000000000"),
			}, nil)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, dkgAddress, 1).
			Return(&logpoller.Log{
				Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000a6fca200010576e704b4a519484d6239ef17f1f5b4a82e330b0daf827ed4dc2789971b0000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000002e000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000a8cbea12a06869d3ec432ab9682dab6c761d591000000000000000000000000f4f9db7bb1d16b7cdfb18ec68994c26964f5985300000000000000000000000022fb3f90c539457f00d8484438869135e604a65500000000000000000000000033cbcedccb11c9773ad78e214ba342e979255ab30000000000000000000000006ffaa96256fbc1012325cca88c79f725c33eed80000000000000000000000000000000000000000000000000000000000000000500000000000000000000000074103cf8b436465870b26aa9fa2f62ad62b22e3500000000000000000000000038a6cb196f805cc3041f6645a5a6cec27b64430d00000000000000000000000047d7095cfebf8285bdaa421bc8268d0db87d933c000000000000000000000000a8842be973800ff61d80d2d53fa62c3a685380eb0000000000000000000000003750e31321aee8c024751877070e8d5f704ce98700000000000000000000000000000000000000000000000000000000000000206f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf000000000000000000000000000000000000000000000000000000000000004220880d88ee16f1080c8afa0251880c8afa025208090dfc04a288090dfc04a30033a05010101010142206c5ca6f74b532222ac927dd3de235d46a943e372c0563393a33b01dcfd3f371c4220855114d25c2ef5e85fffe4f20a365672d8f2dba3b2ec82333f494168a2039c0442200266e835634db00977cbc1caa4db10e1676c1a4c0fcbc6ba7f09300f0d1831824220980cd91f7a73f20f4b0d51d00cd4e00373dc2beafbb299ca3c609757ab98c8304220eb6d36e2af8922085ff510bbe1eb8932a0e3295ca9f047fef25d90e69c52948f4a34313244334b6f6f574463364b7232644542684b59326b336e685057694676544565325331703978544532544b74344d7572716f684a34313244334b6f6f574b436e4367724b637743324a3577576a626e355435335068646b6b6f57454e534a39546537544b7836366f4a4a34313244334b6f6f575239616f675948786b357a38636b624c4c56346e426f7a777a747871664a7050586671336d4a7232796452474a34313244334b6f6f5744695444635565675637776b313133473366476a69616259756f54436f3157726f6f53656741343263556f544a34313244334b6f6f574e64687072586b5472665370354d5071736270467a70364167394a53787358694341434442676454424c656652820300050e416c74424e2d3132382047e282810e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c216a812d7616e47f0bd38fa4863f48fbcda6a38af4c58d2233dfa7cf79620947042d09f923e0a2f7a2270391e8b058d8bdb8f79fe082b7b627f025651c7290382fdff97c3181d15d162c146ce87ff752499d2acc2b26011439a12e29571a6f1e1defb1751c3be4258c493984fd9f0f6b4a26c539870b5f15bfed3d8ffac92499eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b127a112970e1adf615f823b2b2180754c2f0ee01f1b389e56df55ca09702cd0401b66ff71779d2dd67222503a85ab921b28c329cc1832800b192d0b0247c0776e1b9653dc00df48daa6364287c84c0382f5165e7269fef06d10bc67c1bba252305d1af0dc7bb0fe92558eb4c5f38c23163dee1cfb34a72020669dbdfe337c16f3307472616e736c61746f722066726f6d20416c74424e2d3132382047e2828120746f20416c74424e2d3132382047e282825880ade2046080c8afa0256880c8afa0257080ade204788094ebdc0382019e010a205034214e0bd4373f38e162cf9fc9133e2f3b71441faa4c3d1ac01c1877f1cd2712200e03e975b996f911abba2b79d2596c2150bc94510963c40a1137a03df6edacdb1a107dee1cdb894163813bb3da604c9c133c1a10bb33302eeafbd55d352e35dcc5d2b3311a10d2c658b6b93d74a02d467849b6fe75251a10fea5308cc1fea69e7246eafe7ca8a3a51a1048efe1ad873b6f025ac0243bdef715f8000000000000000000000000000000000000000000000000000000000000"),
			}, nil)
		defer lp.AssertExpectations(t)

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
			topics:             tp,
			coordinatorAddress: coordinatorAddress,
			dkgAddress:         dkgAddress,
		}
		actualDKG, actualVRF, err := c.DKGVRFCommittees(context.TODO())
		assert.NoError(t, err)
		assert.ElementsMatch(t, expectedDKGVRF.Signers, actualDKG.Signers)
		assert.ElementsMatch(t, expectedDKGVRF.Transmitters, actualDKG.Transmitters)
		assert.ElementsMatch(t, expectedDKGVRF.Signers, actualVRF.Signers)
		assert.ElementsMatch(t, expectedDKGVRF.Transmitters, actualVRF.Transmitters)
	})

	t.Run("vrf log poll fails", func(t *testing.T) {
		lp := &lp_mocks.LogPoller{}
		tp := newTopics()

		coordinatorAddress := newAddress(t)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, coordinatorAddress, 1).
			Return(nil, errors.New("rpc error"))
		defer lp.AssertExpectations(t)

		c := &coordinator{
			lp:                 lp,
			topics:             tp,
			coordinatorAddress: coordinatorAddress,
		}

		_, _, err := c.DKGVRFCommittees(context.TODO())
		assert.Error(t, err)
	})

	t.Run("dkg log poll fails", func(t *testing.T) {
		lp := &lp_mocks.LogPoller{}
		tp := newTopics()
		coordinatorAddress := newAddress(t)
		dkgAddress := newAddress(t)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, coordinatorAddress, 1).
			Return(&logpoller.Log{
				Data: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000a6fca200010576e704b4a519484d6239ef17f1f5b4a82e330b0daf827ed4dc2789971b0000000000000000000000000000000000000000000000000000000000000032000000000000000000000000000000000000000000000000000000000000012000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000002e000000000000000000000000000000000000000000000000000000000000000050000000000000000000000000a8cbea12a06869d3ec432ab9682dab6c761d591000000000000000000000000f4f9db7bb1d16b7cdfb18ec68994c26964f5985300000000000000000000000022fb3f90c539457f00d8484438869135e604a65500000000000000000000000033cbcedccb11c9773ad78e214ba342e979255ab30000000000000000000000006ffaa96256fbc1012325cca88c79f725c33eed80000000000000000000000000000000000000000000000000000000000000000500000000000000000000000074103cf8b436465870b26aa9fa2f62ad62b22e3500000000000000000000000038a6cb196f805cc3041f6645a5a6cec27b64430d00000000000000000000000047d7095cfebf8285bdaa421bc8268d0db87d933c000000000000000000000000a8842be973800ff61d80d2d53fa62c3a685380eb0000000000000000000000003750e31321aee8c024751877070e8d5f704ce98700000000000000000000000000000000000000000000000000000000000000206f3b82406688b8ddb944c6f2e6d808f014c8fa8d568d639c25019568c715fbf000000000000000000000000000000000000000000000000000000000000004220880d88ee16f1080c8afa0251880c8afa025208090dfc04a288090dfc04a30033a05010101010142206c5ca6f74b532222ac927dd3de235d46a943e372c0563393a33b01dcfd3f371c4220855114d25c2ef5e85fffe4f20a365672d8f2dba3b2ec82333f494168a2039c0442200266e835634db00977cbc1caa4db10e1676c1a4c0fcbc6ba7f09300f0d1831824220980cd91f7a73f20f4b0d51d00cd4e00373dc2beafbb299ca3c609757ab98c8304220eb6d36e2af8922085ff510bbe1eb8932a0e3295ca9f047fef25d90e69c52948f4a34313244334b6f6f574463364b7232644542684b59326b336e685057694676544565325331703978544532544b74344d7572716f684a34313244334b6f6f574b436e4367724b637743324a3577576a626e355435335068646b6b6f57454e534a39546537544b7836366f4a4a34313244334b6f6f575239616f675948786b357a38636b624c4c56346e426f7a777a747871664a7050586671336d4a7232796452474a34313244334b6f6f5744695444635565675637776b313133473366476a69616259756f54436f3157726f6f53656741343263556f544a34313244334b6f6f574e64687072586b5472665370354d5071736270467a70364167394a53787358694341434442676454424c656652820300050e416c74424e2d3132382047e282810e86e8cf899ae9a1b43e023bbe8825b103659bb8d6d4e54f6a3cfae7b106069c216a812d7616e47f0bd38fa4863f48fbcda6a38af4c58d2233dfa7cf79620947042d09f923e0a2f7a2270391e8b058d8bdb8f79fe082b7b627f025651c7290382fdff97c3181d15d162c146ce87ff752499d2acc2b26011439a12e29571a6f1e1defb1751c3be4258c493984fd9f0f6b4a26c539870b5f15bfed3d8ffac92499eb62dbd2beb7c1524275a8019022f6ce6a7e86c9e65e3099452a2b96fc2432b127a112970e1adf615f823b2b2180754c2f0ee01f1b389e56df55ca09702cd0401b66ff71779d2dd67222503a85ab921b28c329cc1832800b192d0b0247c0776e1b9653dc00df48daa6364287c84c0382f5165e7269fef06d10bc67c1bba252305d1af0dc7bb0fe92558eb4c5f38c23163dee1cfb34a72020669dbdfe337c16f3307472616e736c61746f722066726f6d20416c74424e2d3132382047e2828120746f20416c74424e2d3132382047e282825880ade2046080c8afa0256880c8afa0257080ade204788094ebdc0382019e010a205034214e0bd4373f38e162cf9fc9133e2f3b71441faa4c3d1ac01c1877f1cd2712200e03e975b996f911abba2b79d2596c2150bc94510963c40a1137a03df6edacdb1a107dee1cdb894163813bb3da604c9c133c1a10bb33302eeafbd55d352e35dcc5d2b3311a10d2c658b6b93d74a02d467849b6fe75251a10fea5308cc1fea69e7246eafe7ca8a3a51a1048efe1ad873b6f025ac0243bdef715f8000000000000000000000000000000000000000000000000000000000000"),
			}, nil)
		lp.On("LatestLogByEventSigWithConfs", tp.configSetTopic, dkgAddress, 1).
			Return(nil, errors.New("rpc error"))
		defer lp.AssertExpectations(t)

		c := &coordinator{
			lp:                 lp,
			topics:             tp,
			coordinatorAddress: coordinatorAddress,
			dkgAddress:         dkgAddress,
		}
		_, _, err := c.DKGVRFCommittees(context.TODO())
		assert.Error(t, err)
	})
}

func TestCoordinator_ProvingKeyHash(t *testing.T) {
	t.Run("valid output", func(t *testing.T) {
		h := crypto.Keccak256Hash([]byte("hello world"))
		var expected [32]byte
		copy(expected[:], h.Bytes())
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("SProvingKeyHash", mock.Anything).
			Return(expected, nil)
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		provingKeyHash, err := c.ProvingKeyHash(context.TODO())
		assert.NoError(t, err)
		assert.ElementsMatch(t, expected[:], provingKeyHash[:])
	})

	t.Run("invalid output", func(t *testing.T) {
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("SProvingKeyHash", mock.Anything).
			Return([32]byte{}, errors.New("rpc error"))
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		_, err := c.ProvingKeyHash(context.TODO())
		assert.Error(t, err)
	})
}

func TestCoordinator_ReportBlocks(t *testing.T) {
	t.Run("happy path, beacon requests", func(t *testing.T) {
		coordinatorAddress := newAddress(t)

		// we only need the contract for unmarshaling raw log data,
		// so the backend can be safely set to nil.
		// in actual operation, the backend will be an evm client.
		coordinatorContract, err := vrf_wrapper.NewVRFBeaconCoordinator(coordinatorAddress, nil)
		require.NoError(t, err)

		latestHeadNumber := int64(200)
		evmClient := &evm_mocks.Client{}
		evmClient.On("HeadByNumber", mock.Anything, mock.Anything).
			Return(&evmtypes.Head{
				Number: latestHeadNumber,
			}, nil)
		defer evmClient.AssertExpectations(t)

		tp := newTopics()

		lookbackBlocks := int64(50)
		lp := &lp_mocks.LogPoller{}
		lp.On(
			"LogsWithSigs",
			latestHeadNumber-lookbackBlocks,
			latestHeadNumber,
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.newTransmissionTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191),
			newRandomnessRequestedLog(t, 3, 195, 192),
			newRandomnessRequestedLog(t, 3, 195, 193),
		}, nil)
		defer lp.AssertExpectations(t)

		lp.On("GetBlocks", []uint64{195}, mock.Anything).
			Return([]logpoller.LogPollerBlock{
				{
					BlockNumber: 195,
					BlockHash:   common.HexToHash("0x002"),
				},
			}, nil)

		c := &coordinator{
			coordinatorContract: coordinatorContract,
			coordinatorAddress:  coordinatorAddress,
			lp:                  lp,
			lookbackBlocks:      lookbackBlocks,
			lggr:                logger.TestLogger(t),
			topics:              tp,
			evmClient:           evmClient,
		}

		blocks, callbacks, err := c.ReportBlocks(
			context.TODO(),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 1)
		assert.Len(t, callbacks, 0)
	})

	t.Run("happy path, callback requests", func(t *testing.T) {
		coordinatorAddress := newAddress(t)

		// we only need the contract for unmarshaling raw log data,
		// so the backend can be safely set to nil.
		// in actual operation, the backend will be an evm client.
		coordinatorContract, err := vrf_wrapper.NewVRFBeaconCoordinator(coordinatorAddress, nil)
		require.NoError(t, err)

		latestHeadNumber := int64(200)
		evmClient := &evm_mocks.Client{}
		evmClient.On("HeadByNumber", mock.Anything, mock.Anything).
			Return(&evmtypes.Head{
				Number: latestHeadNumber,
			}, nil)
		defer evmClient.AssertExpectations(t)

		tp := newTopics()

		lookbackBlocks := int64(50)
		lp := &lp_mocks.LogPoller{}
		lp.On(
			"LogsWithSigs",
			latestHeadNumber-lookbackBlocks,
			latestHeadNumber,
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.newTransmissionTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3),
		}, nil)
		defer lp.AssertExpectations(t)

		lp.On("GetBlocks", []uint64{195}, mock.Anything).
			Return([]logpoller.LogPollerBlock{
				{
					BlockNumber: 195,
					BlockHash:   common.HexToHash("0x002"),
				},
			}, nil)

		c := &coordinator{
			coordinatorContract: coordinatorContract,
			coordinatorAddress:  coordinatorAddress,
			lp:                  lp,
			lookbackBlocks:      lookbackBlocks,
			lggr:                logger.TestLogger(t),
			topics:              tp,
			evmClient:           evmClient,
		}

		blocks, callbacks, err := c.ReportBlocks(
			context.TODO(),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 1)
		assert.Len(t, callbacks, 3)
	})

	t.Run("happy path, beacon requests, beacon fulfillments", func(t *testing.T) {
		coordinatorAddress := newAddress(t)

		// we only need the contract for unmarshaling raw log data,
		// so the backend can be safely set to nil.
		// in actual operation, the backend will be an evm client.
		coordinatorContract, err := vrf_wrapper.NewVRFBeaconCoordinator(coordinatorAddress, nil)
		require.NoError(t, err)

		latestHeadNumber := int64(200)
		evmClient := &evm_mocks.Client{}
		evmClient.On("HeadByNumber", mock.Anything, mock.Anything).
			Return(&evmtypes.Head{
				Number: latestHeadNumber,
			}, nil)
		defer evmClient.AssertExpectations(t)

		tp := newTopics()

		lookbackBlocks := int64(50)
		lp := &lp_mocks.LogPoller{}
		lp.On(
			"LogsWithSigs",
			latestHeadNumber-lookbackBlocks,
			latestHeadNumber,
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.newTransmissionTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessRequestedLog(t, 3, 195, 191),
			newRandomnessRequestedLog(t, 3, 195, 192),
			newRandomnessRequestedLog(t, 3, 195, 193),
			newNewTransmissionLog(t, []vrf_wrapper.VRFBeaconReportOutputServed{
				{
					Height:            195,
					ConfirmationDelay: big.NewInt(3),
				},
			}),
		}, nil)
		defer lp.AssertExpectations(t)

		var r []uint64
		lp.On("GetBlocks", r, mock.Anything).
			Return(nil, nil)

		c := &coordinator{
			coordinatorContract: coordinatorContract,
			coordinatorAddress:  coordinatorAddress,
			lp:                  lp,
			lookbackBlocks:      lookbackBlocks,
			lggr:                logger.TestLogger(t),
			topics:              tp,
			evmClient:           evmClient,
		}

		blocks, callbacks, err := c.ReportBlocks(
			context.TODO(),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 0)
		assert.Len(t, callbacks, 0)
	})

	t.Run("happy path, callback requests, callback fulfillments", func(t *testing.T) {
		coordinatorAddress := newAddress(t)

		// we only need the contract for unmarshaling raw log data,
		// so the backend can be safely set to nil.
		// in actual operation, the backend will be an evm client.
		coordinatorContract, err := vrf_wrapper.NewVRFBeaconCoordinator(coordinatorAddress, nil)
		require.NoError(t, err)

		latestHeadNumber := int64(200)
		evmClient := &evm_mocks.Client{}
		evmClient.On("HeadByNumber", mock.Anything, mock.Anything).
			Return(&evmtypes.Head{
				Number: latestHeadNumber,
			}, nil)
		defer evmClient.AssertExpectations(t)

		tp := newTopics()

		lookbackBlocks := int64(50)
		lp := &lp_mocks.LogPoller{}
		lp.On(
			"LogsWithSigs",
			latestHeadNumber-lookbackBlocks,
			latestHeadNumber,
			[]common.Hash{
				tp.randomnessRequestedTopic,
				tp.randomnessFulfillmentRequestedTopic,
				tp.randomWordsFulfilledTopic,
				tp.newTransmissionTopic,
			},
			coordinatorAddress,
			mock.Anything,
		).Return([]logpoller.Log{
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 191, 1),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 192, 2),
			newRandomnessFulfillmentRequestedLog(t, 3, 195, 193, 3),
			// Both RandomWordsFulfilled and NewTransmission events are emitted
			// when a VRF fulfillment happens on chain.
			newRandomWordsFulfilledLog(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, []byte{1, 1, 1}),
			newNewTransmissionLog(t, []vrf_wrapper.VRFBeaconReportOutputServed{
				{
					Height:            195,
					ConfirmationDelay: big.NewInt(3),
				},
			}),
		}, nil)
		defer lp.AssertExpectations(t)

		var r []uint64
		lp.On("GetBlocks", r, mock.Anything).
			Return(nil, nil)

		c := &coordinator{
			coordinatorContract: coordinatorContract,
			coordinatorAddress:  coordinatorAddress,
			lp:                  lp,
			lookbackBlocks:      lookbackBlocks,
			lggr:                logger.TestLogger(t),
			topics:              tp,
			evmClient:           evmClient,
		}

		blocks, callbacks, err := c.ReportBlocks(
			context.TODO(),
			0, // slotInterval: unused
			map[uint32]struct{}{3: {}},
			time.Duration(0),
			100, // maxBlocks: unused
			100, // maxCallbacks: unused
		)
		assert.NoError(t, err)
		assert.Len(t, blocks, 0)
		assert.Len(t, callbacks, 0)
	})
}

func TestCoordinator_ReportWillBeTransmitted(t *testing.T) {
	c := &coordinator{}
	assert.NoError(t, c.ReportWillBeTransmitted(context.TODO(), ocr2vrftypes.AbstractReport{}))
}

func TestCoordinator_MarshalUnmarshal(t *testing.T) {
	coordinatorAddress := newAddress(t)
	vrfBeaconCoordinator, err := vrf_wrapper.NewVRFBeaconCoordinator(coordinatorAddress, nil)
	require.NoError(t, err)

	lg := newRandomnessRequestedLog(t, 3, 1500, 1450)
	rrIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	rr, ok := rrIface.(*vrf_wrapper.VRFBeaconCoordinatorRandomnessRequested)
	require.True(t, ok)
	assert.Equal(t, uint64(1500), rr.NextBeaconOutputHeight)
	assert.Equal(t, int64(3), rr.ConfDelay.Int64())

	lg = newRandomnessFulfillmentRequestedLog(t, 3, 1500, 1450, 1)
	rfrIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	rfr, ok := rfrIface.(*vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested)
	require.True(t, ok)
	assert.Equal(t, uint64(1500), rfr.NextBeaconOutputHeight)
	assert.Equal(t, int64(3), rfr.ConfDelay.Int64())
	assert.Equal(t, int64(1), rfr.Callback.RequestID.Int64())

	lg = newRandomWordsFulfilledLog(t, []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}, []byte{1, 1, 1})
	rwfIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	rwf, ok := rwfIface.(*vrf_wrapper.VRFBeaconCoordinatorRandomWordsFulfilled)
	require.True(t, ok)
	assert.Equal(t, []int64{1, 2, 3}, []int64{rwf.RequestIDs[0].Int64(), rwf.RequestIDs[1].Int64(), rwf.RequestIDs[2].Int64()})
	assert.Equal(t, []byte{1, 1, 1}, rwf.SuccessfulFulfillment)

	lg = newNewTransmissionLog(t, []vrf_wrapper.VRFBeaconReportOutputServed{
		{
			Height:            1500,
			ConfirmationDelay: big.NewInt(3),
		},
		{
			Height:            1505,
			ConfirmationDelay: big.NewInt(4),
		},
	})
	ntIface, err := vrfBeaconCoordinator.ParseLog(toGethLog(lg))
	require.NoError(t, err)
	nt, ok := ntIface.(*vrf_wrapper.VRFBeaconCoordinatorNewTransmission)
	require.True(t, ok)
	assert.Equal(t, uint64(1500), nt.OutputsServed[0].Height)
	assert.Equal(t, uint64(1505), nt.OutputsServed[1].Height)
	assert.Equal(t, int64(3), nt.OutputsServed[0].ConfirmationDelay.Int64())
	assert.Equal(t, int64(4), nt.OutputsServed[1].ConfirmationDelay.Int64())
}

func TestCoordinator_ReportIsOnchain(t *testing.T) {
	t.Run("report is on-chain", func(t *testing.T) {
		tp := newTopics()
		coordinatorAddress := newAddress(t)

		epoch := uint32(20)
		round := uint8(3)
		epochAndRound := toEpochAndRoundUint40(epoch, round)
		enrTopic := common.BytesToHash(common.LeftPadBytes(epochAndRound.Bytes(), 32))
		lp := &lp_mocks.LogPoller{}
		lp.On("IndexedLogs", tp.newTransmissionTopic, coordinatorAddress, 2, []common.Hash{
			enrTopic,
		}, 1, mock.Anything).Return([]logpoller.Log{
			{
				BlockNumber: 195,
			},
		}, nil)

		c := &coordinator{
			lp:                 lp,
			lggr:               logger.TestLogger(t),
			coordinatorAddress: coordinatorAddress,
			topics:             tp,
		}

		present, err := c.ReportIsOnchain(context.TODO(), epoch, round)
		assert.NoError(t, err)
		assert.True(t, present)
	})

	t.Run("report is not on-chain", func(t *testing.T) {
		tp := newTopics()
		coordinatorAddress := newAddress(t)

		epoch := uint32(20)
		round := uint8(3)
		epochAndRound := toEpochAndRoundUint40(epoch, round)
		enrTopic := common.BytesToHash(common.LeftPadBytes(epochAndRound.Bytes(), 32))
		lp := &lp_mocks.LogPoller{}
		lp.On("IndexedLogs", tp.newTransmissionTopic, coordinatorAddress, 2, []common.Hash{
			enrTopic,
		}, 1, mock.Anything).Return([]logpoller.Log{}, nil)

		c := &coordinator{
			lp:                 lp,
			lggr:               logger.TestLogger(t),
			coordinatorAddress: coordinatorAddress,
			topics:             tp,
		}

		present, err := c.ReportIsOnchain(context.TODO(), epoch, round)
		assert.NoError(t, err)
		assert.False(t, present)
	})

}

func TestCoordinator_ConfirmationDelays(t *testing.T) {
	t.Run("valid output", func(t *testing.T) {
		expected := [8]uint32{1, 2, 3, 4, 5, 6, 7, 8}
		ret := [8]*big.Int{}
		for i, delay := range expected {
			ret[i] = big.NewInt(int64(delay))
		}
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("GetConfirmationDelays", mock.Anything).
			Return(ret, nil)
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		confDelays, err := c.ConfirmationDelays(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expected[:], confDelays[:])
	})

	t.Run("invalid output", func(t *testing.T) {
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("GetConfirmationDelays", mock.Anything).
			Return([8]*big.Int{}, errors.New("rpc error"))
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		_, err := c.ConfirmationDelays(context.TODO())
		assert.Error(t, err)
	})
}

func TestCoordinator_KeyID(t *testing.T) {
	t.Run("valid output", func(t *testing.T) {
		var keyIDBytes [32]byte
		keyIDBytes[0] = 1
		expected := dkg.KeyID(keyIDBytes)
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("SKeyID", mock.Anything).
			Return(keyIDBytes, nil)
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		keyID, err := c.KeyID(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, expected[:], keyID[:])
	})

	t.Run("invalid output", func(t *testing.T) {
		var emptyBytes [32]byte
		coordinatorContract := &mocks.VRFBeaconCoordinator{}
		coordinatorContract.
			On("SKeyID", mock.Anything).
			Return(emptyBytes, errors.New("rpc error"))
		defer coordinatorContract.AssertExpectations(t)
		c := &coordinator{
			coordinatorContract: coordinatorContract,
		}
		_, err := c.KeyID(context.TODO())
		assert.Error(t, err)
	})
}

func TestTopics_DKGConfigSet_VRFConfigSet(t *testing.T) {
	dkgConfigSetTopic := dkg_wrapper.DKGConfigSet{}.Topic()
	vrfConfigSetTopic := vrf_wrapper.VRFBeaconCoordinatorConfigSet{}.Topic()
	assert.Equal(t, dkgConfigSetTopic, vrfConfigSetTopic, "config set topics of vrf and dkg must be equal")
}
func newRandomnessRequestedLog(
	t *testing.T,
	confDelay int64,
	nextBeaconOutputHeight uint64,
	requestBlock uint64,
) logpoller.Log {
	//event RandomnessRequested(
	//  uint64 indexed nextBeaconOutputHeight,
	//  ConfirmationDelay confDelay
	//);
	confDelayType, err := abi.NewType("uint24", "", nil)
	require.NoError(t, err)
	unindexedArgs := abi.Arguments{
		{
			Name: "confDelay",
			Type: confDelayType,
		},
	}
	logData, err := unindexedArgs.Pack(big.NewInt(confDelay))
	require.NoError(t, err)
	nextBeaconOutputHeightType, err := abi.NewType("uint64", "", nil)
	require.NoError(t, err)
	indexedArgs := abi.Arguments{abi.Argument{
		Name: "nextBeaconOutputHeight",
		Type: nextBeaconOutputHeightType,
	}}
	topicData, err := indexedArgs.Pack(nextBeaconOutputHeight)
	require.NoError(t, err)
	topic0 := vrfABI.Events[randomnessRequestedEvent].ID.Bytes()
	lg := logpoller.Log{
		Data: logData,
		Topics: [][]byte{
			// first topic is the event signature
			topic0,
			// second topic is nextBeaconOutputHeight since it's indexed
			topicData,
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
) logpoller.Log {
	//event RandomnessFulfillmentRequested(
	//  uint64 nextBeaconOutputHeight,
	//  ConfirmationDelay confDelay,
	//  uint64 subID,
	//  Callback callback
	//);
	e := vrf_wrapper.VRFBeaconCoordinatorRandomnessFulfillmentRequested{
		ConfDelay:              big.NewInt(confDelay),
		NextBeaconOutputHeight: nextBeaconOutputHeight,
		Callback: vrf_wrapper.VRFBeaconTypesCallback{
			RequestID:    big.NewInt(requestID),
			NumWords:     1,
			GasAllowance: big.NewInt(1000),
		},
		SubID: 1,
		Raw: types.Log{
			BlockNumber: requestBlock,
		},
	}
	packed, err := vrfABI.Events[randomnessFulfillmentRequestedEvent].Inputs.Pack(
		e.NextBeaconOutputHeight, e.ConfDelay, e.SubID, e.Callback)
	require.NoError(t, err)
	topic0 := vrfABI.Events[randomnessFulfillmentRequestedEvent].ID.Bytes()
	return logpoller.Log{
		Data:     packed,
		EventSig: topic0,
		Topics: [][]byte{
			topic0,
		},
		BlockNumber: int64(requestBlock),
	}
}

func newRandomWordsFulfilledLog(
	t *testing.T,
	requestIDs []*big.Int,
	successfulFulfillment []byte,
) logpoller.Log {
	//event RandomWordsFulfilled(
	//  RequestID[] requestIDs,
	//  bytes successfulFulfillment,
	//  bytes[] truncatedErrorData
	//);
	e := vrf_wrapper.VRFBeaconCoordinatorRandomWordsFulfilled{
		RequestIDs:            requestIDs,
		SuccessfulFulfillment: successfulFulfillment,
	}
	packed, err := vrfABI.Events[randomWordsFulfilledEvent].Inputs.Pack(
		e.RequestIDs, e.SuccessfulFulfillment, e.TruncatedErrorData)
	require.NoError(t, err)
	topic0 := vrfABI.Events[randomWordsFulfilledEvent].ID.Bytes()
	return logpoller.Log{
		Data:     packed,
		EventSig: topic0,
		Topics:   [][]byte{topic0},
	}
}

func newNewTransmissionLog(
	t *testing.T,
	outputsServed []vrf_wrapper.VRFBeaconReportOutputServed,
) logpoller.Log {
	//event NewTransmission(
	//  uint32 indexed aggregatorRoundId,
	//  uint40 indexed epochAndRound,
	//  address transmitter,
	//  uint192 juelsPerFeeCoin,
	//  bytes32 configDigest,
	//  OutputServed[] outputsServed
	//);
	e := vrf_wrapper.VRFBeaconCoordinatorNewTransmission{
		AggregatorRoundId: 1,
		OutputsServed:     outputsServed,
		JuelsPerFeeCoin:   big.NewInt(0),
		EpochAndRound:     big.NewInt(1),
		ConfigDigest:      crypto.Keccak256Hash([]byte("hello world")),
		Transmitter:       newAddress(t),
	}
	var unindexed abi.Arguments
	for _, a := range vrfABI.Events[newTransmissionEvent].Inputs {
		if !a.Indexed {
			unindexed = append(unindexed, a)
		}
	}
	nonIndexedData, err := unindexed.Pack(
		e.Transmitter, e.JuelsPerFeeCoin, e.ConfigDigest, e.OutputsServed)
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

	topic0 := vrfABI.Events[newTransmissionEvent].ID.Bytes()
	return logpoller.Log{
		Data: nonIndexedData,
		Topics: [][]byte{
			topic0,
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
