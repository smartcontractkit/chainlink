package chainlink_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	"github.com/smartcontractkit/chainlink-relay/pkg/loop"
	relayutils "github.com/smartcontractkit/chainlink-relay/pkg/utils"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink/v2/core/chains/cosmos"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/plugins"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	configtest "github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest/v2"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestCoreRelayerChainInteroperators(t *testing.T) {

	evmChainID1, evmChainID2 := utils.NewBig(big.NewInt(1)), utils.NewBig(big.NewInt(2))
	solanaChainID1, solanaChainID2 := "solana-id-1", "solana-id-2"
	starknetChainID1, starknetChainID2 := "starknet-id-1", "starknet-id-2"
	cosmosChainID1, cosmosChainID2 := "cosmos-id-1", "cosmos-id-2"

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {

		cfg := evmcfg.Defaults(evmChainID1)
		node1_1 := evmcfg.Node{
			Name:     ptr("Test node chain1:1"),
			WSURL:    models.MustParseURL("ws://localhost:8546"),
			HTTPURL:  models.MustParseURL("http://localhost:8546"),
			SendOnly: ptr(false),
			Order:    ptr(int32(15)),
		}
		node1_2 := evmcfg.Node{
			Name:     ptr("Test node chain1:2"),
			WSURL:    models.MustParseURL("ws://localhost:8547"),
			HTTPURL:  models.MustParseURL("http://localhost:8547"),
			SendOnly: ptr(false),
			Order:    ptr(int32(36)),
		}
		node2_1 := evmcfg.Node{
			Name:     ptr("Test node chain2:1"),
			WSURL:    models.MustParseURL("ws://localhost:8547"),
			HTTPURL:  models.MustParseURL("http://localhost:8547"),
			SendOnly: ptr(false),
			Order:    ptr(int32(11)),
		}
		c.EVM[0] = &evmcfg.EVMConfig{
			ChainID: evmChainID1,
			Enabled: ptr(true),
			Chain:   cfg,
			Nodes:   evmcfg.EVMNodes{&node1_1, &node1_2},
		}
		id2 := utils.NewBig(big.NewInt(2))
		c.EVM = append(c.EVM, &evmcfg.EVMConfig{
			ChainID: evmChainID2,
			Chain:   evmcfg.Defaults(id2),
			Enabled: ptr(true),
			Nodes:   evmcfg.EVMNodes{&node2_1},
		})

		c.Solana = solana.TOMLConfigs{
			&solana.TOMLConfig{
				ChainID: &solanaChainID1,
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes: []*solcfg.Node{{
					Name: ptr("solana chain 1 node 1"),
					URL:  ((*relayutils.URL)(models.MustParseURL("http://localhost:8547").URL())),
				}},
			},
			&solana.TOMLConfig{
				ChainID: &solanaChainID2,
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes: []*solcfg.Node{{
					Name: ptr("solana chain 2 node 1"),
					URL:  ((*relayutils.URL)(models.MustParseURL("http://localhost:8527").URL())),
				}},
			},
		}

		c.Starknet = stkcfg.TOMLConfigs{
			&stkcfg.TOMLConfig{
				ChainID: &starknetChainID1,
				Enabled: ptr(true),
				Chain:   stkcfg.Chain{},
				Nodes: []*stkcfg.Node{
					{
						Name: ptr("starknet chain 1 node 1"),
						URL:  ((*relayutils.URL)(models.MustParseURL("http://localhost:8547").URL())),
					},
					{
						Name: ptr("starknet chain 1 node 2"),
						URL:  ((*relayutils.URL)(models.MustParseURL("http://localhost:8548").URL())),
					},
					{
						Name: ptr("starknet chain 1 node 3"),
						URL:  ((*relayutils.URL)(models.MustParseURL("http://localhost:8549").URL())),
					},
				},
			},
			&stkcfg.TOMLConfig{
				ChainID: &starknetChainID2,
				Enabled: ptr(true),
				Chain:   stkcfg.Chain{},
				Nodes: []*stkcfg.Node{
					{
						Name: ptr("starknet chain 2 node 1"),
						URL:  ((*relayutils.URL)(models.MustParseURL("http://localhost:3547").URL())),
					},
				},
			},
		}

		c.Cosmos = cosmos.CosmosConfigs{
			&cosmos.CosmosConfig{
				ChainID: &cosmosChainID1,
				Enabled: ptr(true),
				Chain: coscfg.Chain{
					GasLimitMultiplier: ptr(decimal.RequireFromString("1.55555")),
					Bech32Prefix:       ptr("wasm"),
					GasToken:           ptr("cosm"),
				},
				Nodes: cosmos.CosmosNodes{
					&coscfg.Node{
						Name:          ptr("cosmos chain 1 node 1"),
						TendermintURL: (*relayutils.URL)(models.MustParseURL("http://localhost:9548").URL()),
					},
				},
			},
			&cosmos.CosmosConfig{
				ChainID: &cosmosChainID2,
				Enabled: ptr(true),
				Chain: coscfg.Chain{
					GasLimitMultiplier: ptr(decimal.RequireFromString("0.777")),
					Bech32Prefix:       ptr("wasm"),
					GasToken:           ptr("cosm"),
				},
				Nodes: cosmos.CosmosNodes{
					&coscfg.Node{
						Name:          ptr("cosmos chain 2 node 1"),
						TendermintURL: (*relayutils.URL)(models.MustParseURL("http://localhost:9598").URL()),
					},
				},
			},
		}
	})

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db, cfg.Database())

	lggr := logger.TestLogger(t)

	factory := chainlink.RelayerFactory{
		Logger:       lggr,
		LoopRegistry: plugins.NewLoopRegistry(lggr, nil),
		GRPCOpts:     loop.GRPCOpts{},
	}

	testctx := testutils.Context(t)

	tests := []struct {
		name                    string
		initFuncs               []chainlink.CoreRelayerChainInitFunc
		expectedRelayerNetworks map[relay.Network]struct{}

		expectedEVMChainCnt   int
		expectedEVMNodeCnt    int
		expectedEVMRelayerIds []relay.ID

		expectedSolanaChainCnt   int
		expectedSolanaNodeCnt    int
		expectedSolanaRelayerIds []relay.ID

		expectedStarknetChainCnt   int
		expectedStarknetNodeCnt    int
		expectedStarknetRelayerIds []relay.ID

		expectedCosmosChainCnt   int
		expectedCosmosNodeCnt    int
		expectedCosmosRelayerIds []relay.ID
	}{

		{name: "2 evm chains with 3 nodes",
			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitEVM(testctx, factory, chainlink.EVMFactoryConfig{
					ChainOpts: evm.ChainOpts{
						AppConfig:        cfg,
						EventBroadcaster: pg.NewNullEventBroadcaster(),
						MailMon:          &utils.MailboxMonitor{},
						DB:               db,
					},
					CSAETHKeystore: keyStore,
				}),
			},
			expectedEVMChainCnt: 2,
			expectedEVMNodeCnt:  3,
			expectedEVMRelayerIds: []relay.ID{
				{Network: relay.EVM, ChainID: relay.ChainID(evmChainID1.String())},
				{Network: relay.EVM, ChainID: relay.ChainID(evmChainID2.String())},
			},
			expectedRelayerNetworks: map[relay.Network]struct{}{relay.EVM: {}},
		},

		{name: "2 solana chain with 2 node",

			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitSolana(testctx, factory, chainlink.SolanaFactoryConfig{
					Keystore:    keyStore.Solana(),
					TOMLConfigs: cfg.SolanaConfigs()}),
			},
			expectedSolanaChainCnt: 2,
			expectedSolanaNodeCnt:  2,
			expectedSolanaRelayerIds: []relay.ID{
				{Network: relay.Solana, ChainID: relay.ChainID(solanaChainID1)},
				{Network: relay.Solana, ChainID: relay.ChainID(solanaChainID2)},
			},
			expectedRelayerNetworks: map[relay.Network]struct{}{relay.Solana: {}},
		},

		{name: "2 starknet chain with 4 nodes",

			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitStarknet(testctx, factory, chainlink.StarkNetFactoryConfig{
					Keystore:    keyStore.StarkNet(),
					TOMLConfigs: cfg.StarknetConfigs()}),
			},
			expectedStarknetChainCnt: 2,
			expectedStarknetNodeCnt:  4,
			expectedStarknetRelayerIds: []relay.ID{
				{Network: relay.StarkNet, ChainID: relay.ChainID(starknetChainID1)},
				{Network: relay.StarkNet, ChainID: relay.ChainID(starknetChainID2)},
			},
			expectedRelayerNetworks: map[relay.Network]struct{}{relay.StarkNet: {}},
		},

		{
			name: "2 cosmos chains with 2 nodes",
			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitCosmos(testctx, factory, chainlink.CosmosFactoryConfig{
					Keystore:         keyStore.Cosmos(),
					CosmosConfigs:    cfg.CosmosConfigs(),
					EventBroadcaster: pg.NewNullEventBroadcaster(),
					DB:               db,
					QConfig:          cfg.Database()}),
			},
			expectedCosmosChainCnt: 2,
			expectedCosmosNodeCnt:  2,
			expectedCosmosRelayerIds: []relay.ID{
				{Network: relay.Cosmos, ChainID: relay.ChainID(cosmosChainID1)},
				{Network: relay.Cosmos, ChainID: relay.ChainID(cosmosChainID2)},
			},
			expectedRelayerNetworks: map[relay.Network]struct{}{relay.Cosmos: {}},
		},

		{name: "all chains",

			initFuncs: []chainlink.CoreRelayerChainInitFunc{chainlink.InitSolana(testctx, factory, chainlink.SolanaFactoryConfig{
				Keystore:    keyStore.Solana(),
				TOMLConfigs: cfg.SolanaConfigs()}),
				chainlink.InitEVM(testctx, factory, chainlink.EVMFactoryConfig{
					ChainOpts: evm.ChainOpts{
						AppConfig:        cfg,
						EventBroadcaster: pg.NewNullEventBroadcaster(),
						MailMon:          &utils.MailboxMonitor{},
						DB:               db,
					},
					CSAETHKeystore: keyStore,
				}),
				chainlink.InitStarknet(testctx, factory, chainlink.StarkNetFactoryConfig{
					Keystore:    keyStore.StarkNet(),
					TOMLConfigs: cfg.StarknetConfigs()}),
				chainlink.InitCosmos(testctx, factory, chainlink.CosmosFactoryConfig{
					Keystore:         keyStore.Cosmos(),
					CosmosConfigs:    cfg.CosmosConfigs(),
					EventBroadcaster: pg.NewNullEventBroadcaster(),
					DB:               db,
					QConfig:          cfg.Database(),
				}),
			},
			expectedEVMChainCnt: 2,
			expectedEVMNodeCnt:  3,
			expectedEVMRelayerIds: []relay.ID{
				{Network: relay.EVM, ChainID: relay.ChainID(evmChainID1.String())},
				{Network: relay.EVM, ChainID: relay.ChainID(evmChainID2.String())},
			},

			expectedSolanaChainCnt: 2,
			expectedSolanaNodeCnt:  2,
			expectedSolanaRelayerIds: []relay.ID{
				{Network: relay.Solana, ChainID: relay.ChainID(solanaChainID1)},
				{Network: relay.Solana, ChainID: relay.ChainID(solanaChainID2)},
			},

			expectedStarknetChainCnt: 2,
			expectedStarknetNodeCnt:  4,
			expectedStarknetRelayerIds: []relay.ID{
				{Network: relay.StarkNet, ChainID: relay.ChainID(starknetChainID1)},
				{Network: relay.StarkNet, ChainID: relay.ChainID(starknetChainID2)},
			},

			expectedCosmosChainCnt: 2,
			expectedCosmosNodeCnt:  2,
			expectedCosmosRelayerIds: []relay.ID{
				{Network: relay.Cosmos, ChainID: relay.ChainID(cosmosChainID1)},
				{Network: relay.Cosmos, ChainID: relay.ChainID(cosmosChainID2)},
			},

			expectedRelayerNetworks: map[relay.Network]struct{}{relay.EVM: {}, relay.Cosmos: {}, relay.Solana: {}, relay.StarkNet: {}},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var cr *chainlink.CoreRelayerChainInteroperators
			{
				var err error
				cr, err = chainlink.NewCoreRelayerChainInteroperators(tt.initFuncs...)
				require.NoError(t, err)

				expectedChainCnt := tt.expectedEVMChainCnt + tt.expectedCosmosChainCnt + tt.expectedSolanaChainCnt + tt.expectedStarknetChainCnt
				allChainsStats, cnt, err := cr.ChainStatuses(testctx, 0, 0)
				assert.NoError(t, err)
				assert.Len(t, allChainsStats, expectedChainCnt)
				assert.Equal(t, cnt, len(allChainsStats))
				assert.Len(t, cr.Slice(), expectedChainCnt)

				// should be one relayer per chain and one service per relayer
				assert.Len(t, cr.Slice(), expectedChainCnt)
				assert.Len(t, cr.Services(), expectedChainCnt)

				expectedNodeCnt := tt.expectedEVMNodeCnt + tt.expectedCosmosNodeCnt + tt.expectedSolanaNodeCnt + tt.expectedStarknetNodeCnt
				allNodeStats, cnt, err := cr.NodeStatuses(testctx, 0, 0)
				assert.NoError(t, err)
				assert.Len(t, allNodeStats, expectedNodeCnt)
				assert.Equal(t, cnt, len(allNodeStats))
			}

			gotRelayerNetworks := make(map[relay.Network]struct{})
			for relayNetwork := range relay.SupportedRelays {
				var expectedChainCnt, expectedNodeCnt int
				switch relayNetwork {
				case relay.EVM:
					expectedChainCnt, expectedNodeCnt = tt.expectedEVMChainCnt, tt.expectedEVMNodeCnt
				case relay.Cosmos:
					expectedChainCnt, expectedNodeCnt = tt.expectedCosmosChainCnt, tt.expectedCosmosNodeCnt
				case relay.Solana:
					expectedChainCnt, expectedNodeCnt = tt.expectedSolanaChainCnt, tt.expectedSolanaNodeCnt
				case relay.StarkNet:
					expectedChainCnt, expectedNodeCnt = tt.expectedStarknetChainCnt, tt.expectedStarknetNodeCnt
				default:
					require.Fail(t, "untested relay network", relayNetwork)
				}

				interops := cr.List(chainlink.FilterRelayersByType(relayNetwork))
				assert.Len(t, cr.List(chainlink.FilterRelayersByType(relayNetwork)).Slice(), expectedChainCnt)
				if len(interops.Slice()) > 0 {
					gotRelayerNetworks[relayNetwork] = struct{}{}
				}

				// check legacy chains for those that haven't migrated fully to the loop relayer interface
				if relayNetwork == relay.EVM {
					_, wantEVM := tt.expectedRelayerNetworks[relay.EVM]
					if wantEVM {
						assert.Len(t, cr.LegacyEVMChains().Slice(), expectedChainCnt)
					} else {
						assert.Nil(t, cr.LegacyEVMChains())
					}
				}
				if relayNetwork == relay.Cosmos {
					_, wantCosmos := tt.expectedRelayerNetworks[relay.Cosmos]
					if wantCosmos {
						assert.Len(t, cr.LegacyCosmosChains().Slice(), expectedChainCnt)
					} else {
						assert.Nil(t, cr.LegacyCosmosChains())
					}
				}

				nodesStats, cnt, err := interops.NodeStatuses(testctx, 0, 0)
				assert.NoError(t, err)
				assert.Len(t, nodesStats, expectedNodeCnt)
				assert.Equal(t, cnt, len(nodesStats))

			}
			assert.EqualValues(t, gotRelayerNetworks, tt.expectedRelayerNetworks)

			allRelayerIds := [][]relay.ID{
				tt.expectedEVMRelayerIds,
				tt.expectedCosmosRelayerIds,
				tt.expectedSolanaRelayerIds,
				tt.expectedStarknetRelayerIds,
			}

			for _, chainSpecificRelayerIds := range allRelayerIds {
				for _, wantId := range chainSpecificRelayerIds {
					lr, err := cr.Get(wantId)
					assert.NotNil(t, lr)
					assert.NoError(t, err)
					stat, err := cr.ChainStatus(testctx, wantId)
					assert.NoError(t, err)
					assert.Equal(t, wantId.ChainID, stat.ID)
					// check legacy chains for evm and cosmos
					if wantId.Network == relay.EVM {
						c, err := cr.LegacyEVMChains().Get(wantId.ChainID)
						assert.NoError(t, err)
						assert.NotNil(t, c)
						assert.Equal(t, wantId.ChainID, c.ID().String())
					}
					if wantId.Network == relay.Cosmos {
						c, err := cr.LegacyCosmosChains().Get(wantId.ChainID)
						assert.NoError(t, err)
						assert.NotNil(t, c)
						assert.Equal(t, wantId.ChainID, c.ID())
					}
				}
			}

			expectedMissing := relay.ID{Network: relay.Cosmos, ChainID: "not a chain id"}
			unwanted, err := cr.Get(expectedMissing)
			assert.Nil(t, unwanted)
			assert.ErrorIs(t, err, chainlink.ErrNoSuchRelayer)

		})

	}

	t.Run("bad init func", func(t *testing.T) {
		t.Parallel()
		errBadFunc := errors.New("this is a bad func")
		badFunc := func() chainlink.CoreRelayerChainInitFunc {
			return func(op *chainlink.CoreRelayerChainInteroperators) error {
				return errBadFunc
			}
		}
		cr, err := chainlink.NewCoreRelayerChainInteroperators(badFunc())
		assert.Nil(t, cr)
		assert.ErrorIs(t, err, errBadFunc)
	})
}

func ptr[T any](t T) *T { return &t }
