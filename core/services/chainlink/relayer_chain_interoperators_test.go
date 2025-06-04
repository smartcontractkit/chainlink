package chainlink_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"

	solcfg "github.com/smartcontractkit/chainlink-solana/pkg/solana/config"

	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func TestCoreRelayerChainInteroperators(t *testing.T) {
	evmChainID1, evmChainID2 := ubig.New(big.NewInt(1)), ubig.New(big.NewInt(2))
	solanaChainID1, solanaChainID2 := "solana-id-1", "solana-id-2"

	newConfig := func() chainlink.GeneralConfig {
		return configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
			node1_1 := toml.Node{
				Name:     ptr("Test node chain1:1"),
				WSURL:    commonconfig.MustParseURL("ws://localhost:8546"),
				HTTPURL:  commonconfig.MustParseURL("http://localhost:8546"),
				SendOnly: ptr(false),
				Order:    ptr(int32(15)),
			}
			node1_2 := toml.Node{
				Name:     ptr("Test node chain1:2"),
				WSURL:    commonconfig.MustParseURL("ws://localhost:8547"),
				HTTPURL:  commonconfig.MustParseURL("http://localhost:8547"),
				SendOnly: ptr(false),
				Order:    ptr(int32(36)),
			}
			node2_1 := toml.Node{
				Name:     ptr("Test node chain2:1"),
				WSURL:    commonconfig.MustParseURL("ws://localhost:8547"),
				HTTPURL:  commonconfig.MustParseURL("http://localhost:8547"),
				SendOnly: ptr(false),
				Order:    ptr(int32(11)),
			}
			c.EVM[0] = &toml.EVMConfig{
				ChainID: evmChainID1,
				Enabled: ptr(true),
				Chain:   toml.Defaults(evmChainID1),
				Nodes:   toml.EVMNodes{&node1_1, &node1_2},
			}
			id2 := ubig.New(big.NewInt(2))
			c.EVM = append(c.EVM, &toml.EVMConfig{
				ChainID: evmChainID2,
				Chain:   toml.Defaults(id2),
				Enabled: ptr(true),
				Nodes:   toml.EVMNodes{&node2_1},
			})

			c.Solana = solcfg.TOMLConfigs{
				&solcfg.TOMLConfig{
					ChainID: &solanaChainID1,
					Enabled: ptr(true),
					Chain:   solcfg.NewDefault().Chain,
					Nodes: []*solcfg.Node{{
						Name: ptr("solana chain 1 node 1"),
						URL:  ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:8547").URL())),
					}},
				},
				&solcfg.TOMLConfig{
					ChainID: &solanaChainID2,
					Enabled: ptr(true),
					Chain:   solcfg.NewDefault().Chain,
					Nodes: []*solcfg.Node{{
						Name: ptr("solana chain 2 node 1"),
						URL:  ((*commonconfig.URL)(commonconfig.MustParseURL("http://localhost:8527").URL())),
					}},
				},
			}
		})
	}

	cfg := newConfig()
	db := pgtest.NewSqlxDB(t)
	keyStore := cltest.NewKeyStore(t, db)

	lggr := logger.TestLogger(t)

	factory := chainlink.RelayerFactory{
		Logger:               lggr,
		LoopRegistry:         plugins.NewTestLoopRegistry(lggr),
		GRPCOpts:             loop.GRPCOpts{},
		CapabilitiesRegistry: capabilities.NewRegistry(lggr),
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

		expectedStarknetChainCnt int
		expectedStarknetNodeCnt  int

		expectedDummyChainCnt   int
		expectedDummyNodeCnt    int
		expectedDummyRelayerIds []types.RelayID

		expectedCosmosChainCnt int
		expectedCosmosNodeCnt  int
	}{

		{name: "2 evm chains with 3 nodes",
			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitEVM(factory, chainlink.EVMFactoryConfig{
					ChainOpts: legacyevm.ChainOpts{
						ChainConfigs:   cfg.EVMConfigs(),
						DatabaseConfig: cfg.Database(),
						ListenerConfig: cfg.Database().Listener(),
						FeatureConfig:  cfg.Feature(),
						MailMon:        &mailbox.Monitor{},
						DS:             db,
					},
					EthKeystore: keyStore.Eth(),
					CSAKeystore: &keystore.CSASigner{CSA: keyStore.CSA()},
				}),
			},
			expectedEVMChainCnt: 2,
			expectedEVMNodeCnt:  3,
			expectedEVMRelayerIds: []types.RelayID{
				{Network: relay.NetworkEVM, ChainID: evmChainID1.String()},
				{Network: relay.NetworkEVM, ChainID: evmChainID2.String()},
			},
			expectedRelayerNetworks: map[string]struct{}{relay.NetworkEVM: {}},
		},

		{name: "2 solana chain with 2 node",
			initFuncs: []chainlink.CoreRelayerChainInitFunc{
				chainlink.InitSolana(factory, keyStore.Solana(), chainlink.SolanaFactoryConfig{
					TOMLConfigs: newConfig().SolanaConfigs()}),
			},
			expectedSolanaChainCnt: 2,
			expectedSolanaNodeCnt:  2,
			expectedSolanaRelayerIds: []types.RelayID{
				{Network: relay.NetworkSolana, ChainID: solanaChainID1},
				{Network: relay.NetworkSolana, ChainID: solanaChainID2},
			},
			expectedRelayerNetworks: map[string]struct{}{relay.NetworkSolana: {}},
		},

		{name: "all chains",
			initFuncs: []chainlink.CoreRelayerChainInitFunc{chainlink.InitSolana(factory, keyStore.Solana(), chainlink.SolanaFactoryConfig{
				TOMLConfigs: newConfig().SolanaConfigs()}),
				chainlink.InitEVM(factory, chainlink.EVMFactoryConfig{
					ChainOpts: legacyevm.ChainOpts{
						ChainConfigs:   cfg.EVMConfigs(),
						DatabaseConfig: cfg.Database(),
						ListenerConfig: cfg.Database().Listener(),
						FeatureConfig:  cfg.Feature(),

						MailMon: &mailbox.Monitor{},
						DS:      db,
					},
					EthKeystore: keyStore.Eth(),
					CSAKeystore: &keystore.CSASigner{CSA: keyStore.CSA()},
				}),
				chainlink.InitStarknet(factory, keyStore.StarkNet(), cfg.StarknetConfigs()),
				chainlink.InitCosmos(factory, keyStore.Cosmos(), cfg.CosmosConfigs()),
			},
			expectedEVMChainCnt: 2,
			expectedEVMNodeCnt:  3,
			expectedEVMRelayerIds: []types.RelayID{
				{Network: relay.NetworkEVM, ChainID: evmChainID1.String()},
				{Network: relay.NetworkEVM, ChainID: evmChainID2.String()},
			},

			expectedSolanaChainCnt: 2,
			expectedSolanaNodeCnt:  2,
			expectedSolanaRelayerIds: []types.RelayID{
				{Network: relay.NetworkSolana, ChainID: solanaChainID1},
				{Network: relay.NetworkSolana, ChainID: solanaChainID2},
			},

			expectedRelayerNetworks: map[string]struct{}{relay.NetworkEVM: {}, relay.NetworkCosmos: {}, relay.NetworkSolana: {}, relay.NetworkStarkNet: {}},
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
				assert.Len(t, allChainsStats, cnt)
				assert.Len(t, cr.Slice(), expectedChainCnt)

				// should be one relayer per chain and one service per relayer
				assert.Len(t, cr.Slice(), expectedChainCnt)
				assert.Len(t, cr.Services(), expectedChainCnt)

				expectedNodeCnt := tt.expectedEVMNodeCnt + tt.expectedCosmosNodeCnt + tt.expectedSolanaNodeCnt + tt.expectedStarknetNodeCnt
				allNodeStats, cnt, err := cr.NodeStatuses(testctx, 0, 0)
				assert.NoError(t, err)
				assert.Len(t, allNodeStats, expectedNodeCnt)
				assert.Len(t, allNodeStats, cnt)
			}

			gotRelayerNetworks := make(map[string]struct{})
			for relayNetwork := range relay.SupportedNetworks {
				var expectedChainCnt, expectedNodeCnt int
				switch relayNetwork {
				case relay.NetworkEVM:
					expectedChainCnt, expectedNodeCnt = tt.expectedEVMChainCnt, tt.expectedEVMNodeCnt
				case relay.NetworkCosmos:
					expectedChainCnt, expectedNodeCnt = tt.expectedCosmosChainCnt, tt.expectedCosmosNodeCnt
				case relay.NetworkSolana:
					expectedChainCnt, expectedNodeCnt = tt.expectedSolanaChainCnt, tt.expectedSolanaNodeCnt
				case relay.NetworkStarkNet:
					expectedChainCnt, expectedNodeCnt = tt.expectedStarknetChainCnt, tt.expectedStarknetNodeCnt
				case relay.NetworkDummy:
					expectedChainCnt, expectedNodeCnt = tt.expectedDummyChainCnt, tt.expectedDummyNodeCnt
				case relay.NetworkAptos:
					t.Skip("aptos doesn't need a CoreRelayerChainInteroperator")
				case relay.NetworkTron:
					t.Skip("tron doesn't need a CoreRelayerChainInteroperator")
				case relay.NetworkTON:
					t.Skip("ton doesn't need a CoreRelayerChainInteroperator")

				default:
					require.Fail(t, "untested relay network", relayNetwork)
				}

				interops := cr.List(chainlink.FilterRelayersByType(relayNetwork))
				assert.Len(t, cr.List(chainlink.FilterRelayersByType(relayNetwork)).Slice(), expectedChainCnt)
				if len(interops.Slice()) > 0 {
					gotRelayerNetworks[relayNetwork] = struct{}{}
				}

				// check legacy chains for those that haven't migrated fully to the loop relayer interface
				if relayNetwork == relay.NetworkEVM {
					_, wantEVM := tt.expectedRelayerNetworks[relay.NetworkEVM]
					if wantEVM {
						assert.Len(t, cr.LegacyEVMChains().Slice(), expectedChainCnt)
					} else {
						assert.Nil(t, cr.LegacyEVMChains())
					}
				}

				nodesStats, cnt, err := interops.NodeStatuses(testctx, 0, 0)
				assert.NoError(t, err)
				assert.Len(t, nodesStats, expectedNodeCnt)
				assert.Len(t, nodesStats, cnt)
			}
			assert.EqualValues(t, tt.expectedRelayerNetworks, gotRelayerNetworks)

			allRelayerIds := [][]types.RelayID{
				tt.expectedEVMRelayerIds,
				tt.expectedSolanaRelayerIds,
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
					if wantId.Network == relay.NetworkEVM {
						c, err := cr.LegacyEVMChains().Get(wantId.ChainID)
						assert.NoError(t, err)
						assert.NotNil(t, c)
						assert.Equal(t, wantId.ChainID, c.ID().String())
					}
				}
			}

			expectedMissing := types.RelayID{Network: relay.NetworkCosmos, ChainID: "not a chain id"}
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
