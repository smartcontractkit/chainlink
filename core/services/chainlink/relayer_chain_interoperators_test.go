package chainlink_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	coscfg "github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/config"
	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"
	stkcfg "github.com/smartcontractkit/chainlink-starknet/relayer/pkg/chainlink/config"

	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/plugins"

	evmcfg "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
)

func TestCoreRelayerChainInteroperators(t *testing.T) {

	evmChainID1, evmChainID2 := ubig.New(big.NewInt(1)), ubig.New(big.NewInt(2))
	solanaChainID1, solanaChainID2 := "solana-id-1", "solana-id-2"
	starknetChainID1, starknetChainID2 := "starknet-id-1", "starknet-id-2"
	cosmosChainID1, cosmosChainID2 := "cosmos-id-1", "cosmos-id-2"

	cfg := configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {

		cfg := evmcfg.Defaults(evmChainID1)
		node1_1 := evmcfg.Node{
			Name:     ptr("Test node chain1:1"),
			WSURL:    commonconfig.MustParseURL("ws://localhost:8546"),
			HTTPURL:  commonconfig.MustParseURL("http://localhost:8546"),
			SendOnly: ptr(false),
			Order:    ptr(int32(15)),
		}
		node1_2 := evmcfg.Node{
			Name:     ptr("Test node chain1:2"),
			WSURL:    commonconfig.MustParseURL("ws://localhost:8547"),
			HTTPURL:  commonconfig.MustParseURL("http://localhost:8547"),
			SendOnly: ptr(false),
			Order:    ptr(int32(36)),
		}
		node2_1 := evmcfg.Node{
			Name:     ptr("Test node chain2:1"),
			WSURL:    commonconfig.MustParseURL("ws://localhost:8547"),
			HTTPURL:  commonconfig.MustParseURL("http://localhost:8547"),
			SendOnly: ptr(false),
			Order:    ptr(int32(11)),
		}
		c.EVM[0] = &evmcfg.EVMConfig{
			ChainID: evmChainID1,
			Enabled: ptr(true),
			Chain:   cfg,
			Nodes:   evmcfg.EVMNodes{&node1_1, &node1_2},
		}
		id2 := ubig.New(big.NewInt(2))
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
					URL:  ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:8547").URL())),
				}},
			},
			&solana.TOMLConfig{
				ChainID: &solanaChainID2,
				Enabled: ptr(true),
				Chain:   solcfg.Chain{},
				Nodes: []*solcfg.Node{{
					Name: ptr("solana chain 2 node 1"),
					URL:  ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:8527").URL())),
				}},
			},
		}

		c.Starknet = stkcfg.TOMLConfigs{
			&stkcfg.TOMLConfig{
				ChainID:   &starknetChainID1,
				Enabled:   ptr(true),
				Chain:     stkcfg.Chain{},
				FeederURL: commonconfig.MustParseURL("http://feeder.url"),
				Nodes: []*stkcfg.Node{
					{
						Name:   ptr("starknet chain 1 node 1"),
						URL:    ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:8547").URL())),
						APIKey: ptr("key"),
					},
					{
						Name:   ptr("starknet chain 1 node 2"),
						URL:    ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:8548").URL())),
						APIKey: ptr("key"),
					},
					{
						Name:   ptr("starknet chain 1 node 3"),
						URL:    ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:8549").URL())),
						APIKey: ptr("key"),
					},
				},
			},
			&stkcfg.TOMLConfig{
				ChainID:   &starknetChainID2,
				Enabled:   ptr(true),
				Chain:     stkcfg.Chain{},
				FeederURL: commonconfig.MustParseURL("http://feeder.url"),
				Nodes: []*stkcfg.Node{
					{
						Name:   ptr("starknet chain 2 node 1"),
						URL:    ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:3547").URL())),
						APIKey: ptr("key"),
					},
				},
			},
		}

		c.Cosmos = coscfg.TOMLConfigs{
			&coscfg.TOMLConfig{
				ChainID: &cosmosChainID1,
				Enabled: ptr(true),
				Chain: coscfg.Chain{
					GasLimitMultiplier: ptr(decimal.RequireFromString("1.55555")),
					Bech32Prefix:       ptr("wasm"),
					GasToken:           ptr("cosm"),
				},
				Nodes: coscfg.Nodes{
					&coscfg.Node{
						Name:          ptr("cosmos chain 1 node 1"),
						TendermintURL: (*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:9548").URL()),
					},
				},
			},
			&coscfg.TOMLConfig{
				ChainID: &cosmosChainID2,
				Enabled: ptr(true),
				Chain: coscfg.Chain{
					GasLimitMultiplier: ptr(decimal.RequireFromString("0.777")),
					Bech32Prefix:       ptr("wasm"),
					GasToken:           ptr("cosm"),
				},
				Nodes: coscfg.Nodes{
					&coscfg.Node{
						Name:          ptr("cosmos chain 2 node 1"),
						TendermintURL: (*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:9598").URL()),
					},
				},
			},
		}
	})

	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)

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
		expectedRelayerNetworks map[string]struct{}

		expectedEVMChainCnt   int
		expectedEVMNodeCnt    int
		expectedEVMRelayerIds []types.RelayID

		expectedSolanaChainCnt   int
		expectedSolanaNodeCnt    int
		expectedSolanaRelayerIds []types.RelayID

		expectedStarknetChainCnt   int
		expectedStarknetNodeCnt    int
		expectedStarknetRelayerIds []types.RelayID

		expectedCosmosChainCnt   int
		expectedCosmosNodeCnt    int
		expectedCosmosRelayerIds []types.RelayID
	}{

		{name: "2 evm chains with 3 nodes",
			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitEVM(testctx, factory, chainlink.EVMFactoryConfig{
					ChainOpts: legacyevm.ChainOpts{
						AppConfig: cfg,
						MailMon:   &mailbox.Monitor{},
						DS:        db,
					},
					CSAETHKeystore: keyStore,
				}),
			},
			expectedEVMChainCnt: 2,
			expectedEVMNodeCnt:  3,
			expectedEVMRelayerIds: []types.RelayID{
				{Network: types.NetworkEVM, ChainID: evmChainID1.String()},
				{Network: types.NetworkEVM, ChainID: evmChainID2.String()},
			},
			expectedRelayerNetworks: map[string]struct{}{types.NetworkEVM: {}},
		},

		{name: "2 solana chain with 2 node",

			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitSolana(testctx, factory, chainlink.SolanaFactoryConfig{
					Keystore:    keyStore.Solana(),
					TOMLConfigs: cfg.SolanaConfigs()}),
			},
			expectedSolanaChainCnt: 2,
			expectedSolanaNodeCnt:  2,
			expectedSolanaRelayerIds: []types.RelayID{
				{Network: types.NetworkSolana, ChainID: solanaChainID1},
				{Network: types.NetworkSolana, ChainID: solanaChainID2},
			},
			expectedRelayerNetworks: map[string]struct{}{types.NetworkSolana: {}},
		},

		{name: "2 starknet chain with 4 nodes",

			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitStarknet(testctx, factory, chainlink.StarkNetFactoryConfig{
					Keystore:    keyStore.StarkNet(),
					TOMLConfigs: cfg.StarknetConfigs()}),
			},
			expectedStarknetChainCnt: 2,
			expectedStarknetNodeCnt:  4,
			expectedStarknetRelayerIds: []types.RelayID{
				{Network: types.NetworkStarkNet, ChainID: starknetChainID1},
				{Network: types.NetworkStarkNet, ChainID: starknetChainID2},
			},
			expectedRelayerNetworks: map[string]struct{}{types.NetworkStarkNet: {}},
		},

		{
			name: "2 cosmos chains with 2 nodes",
			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitCosmos(testctx, factory, chainlink.CosmosFactoryConfig{
					Keystore:    keyStore.Cosmos(),
					TOMLConfigs: cfg.CosmosConfigs(),
					DS:          db,
				}),
			},
			expectedCosmosChainCnt: 2,
			expectedCosmosNodeCnt:  2,
			expectedCosmosRelayerIds: []types.RelayID{
				{Network: types.NetworkCosmos, ChainID: cosmosChainID1},
				{Network: types.NetworkCosmos, ChainID: cosmosChainID2},
			},
			expectedRelayerNetworks: map[string]struct{}{types.NetworkCosmos: {}},
		},

		{name: "all chains",

			initFuncs: []chainlink.CoreRelayerChainInitFunc{chainlink.InitSolana(testctx, factory, chainlink.SolanaFactoryConfig{
				Keystore:    keyStore.Solana(),
				TOMLConfigs: cfg.SolanaConfigs()}),
				chainlink.InitEVM(testctx, factory, chainlink.EVMFactoryConfig{
					ChainOpts: legacyevm.ChainOpts{
						AppConfig: cfg,

						MailMon: &mailbox.Monitor{},
						DS:      db,
					},
					CSAETHKeystore: keyStore,
				}),
				chainlink.InitStarknet(testctx, factory, chainlink.StarkNetFactoryConfig{
					Keystore:    keyStore.StarkNet(),
					TOMLConfigs: cfg.StarknetConfigs()}),
				chainlink.InitCosmos(testctx, factory, chainlink.CosmosFactoryConfig{
					Keystore:    keyStore.Cosmos(),
					TOMLConfigs: cfg.CosmosConfigs(),
					DS:          db,
				}),
			},
			expectedEVMChainCnt: 2,
			expectedEVMNodeCnt:  3,
			expectedEVMRelayerIds: []types.RelayID{
				{Network: types.NetworkEVM, ChainID: evmChainID1.String()},
				{Network: types.NetworkEVM, ChainID: evmChainID2.String()},
			},

			expectedSolanaChainCnt: 2,
			expectedSolanaNodeCnt:  2,
			expectedSolanaRelayerIds: []types.RelayID{
				{Network: types.NetworkSolana, ChainID: solanaChainID1},
				{Network: types.NetworkSolana, ChainID: solanaChainID2},
			},

			expectedStarknetChainCnt: 2,
			expectedStarknetNodeCnt:  4,
			expectedStarknetRelayerIds: []types.RelayID{
				{Network: types.NetworkStarkNet, ChainID: starknetChainID1},
				{Network: types.NetworkStarkNet, ChainID: starknetChainID2},
			},

			expectedCosmosChainCnt: 2,
			expectedCosmosNodeCnt:  2,
			expectedCosmosRelayerIds: []types.RelayID{
				{Network: types.NetworkCosmos, ChainID: cosmosChainID1},
				{Network: types.NetworkCosmos, ChainID: cosmosChainID2},
			},

			expectedRelayerNetworks: map[string]struct{}{types.NetworkEVM: {}, types.NetworkCosmos: {}, types.NetworkSolana: {}, types.NetworkStarkNet: {}},
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

			gotRelayerNetworks := make(map[string]struct{})
			for relayNetwork := range types.SupportedRelays {
				var expectedChainCnt, expectedNodeCnt int
				switch relayNetwork {
				case types.NetworkEVM:
					expectedChainCnt, expectedNodeCnt = tt.expectedEVMChainCnt, tt.expectedEVMNodeCnt
				case types.NetworkCosmos:
					expectedChainCnt, expectedNodeCnt = tt.expectedCosmosChainCnt, tt.expectedCosmosNodeCnt
				case types.NetworkSolana:
					expectedChainCnt, expectedNodeCnt = tt.expectedSolanaChainCnt, tt.expectedSolanaNodeCnt
				case types.NetworkStarkNet:
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
				if relayNetwork == types.NetworkEVM {
					_, wantEVM := tt.expectedRelayerNetworks[types.NetworkEVM]
					if wantEVM {
						assert.Len(t, cr.LegacyEVMChains().Slice(), expectedChainCnt)
					} else {
						assert.Nil(t, cr.LegacyEVMChains())
					}
				}
				if relayNetwork == types.NetworkCosmos {
					_, wantCosmos := tt.expectedRelayerNetworks[types.NetworkCosmos]
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

			allRelayerIds := [][]types.RelayID{
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
					if wantId.Network == types.NetworkEVM {
						c, err := cr.LegacyEVMChains().Get(wantId.ChainID)
						assert.NoError(t, err)
						assert.NotNil(t, c)
						assert.Equal(t, wantId.ChainID, c.ID().String())
					}
					if wantId.Network == types.NetworkCosmos {
						c, err := cr.LegacyCosmosChains().Get(wantId.ChainID)
						assert.NoError(t, err)
						assert.NotNil(t, c)
						assert.Equal(t, wantId.ChainID, c.ID())
					}
				}
			}

			expectedMissing := types.RelayID{Network: types.NetworkCosmos, ChainID: "not a chain id"}
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
