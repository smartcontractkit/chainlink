package ccip

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	// Third-party imports
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/gagliardetto/solana-go"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	solState "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/solana"

	selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/burnmint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/v0_1_1/lockrelease_token_pool"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"
)

var (
	CommonTestLabels = map[string]string{
		"branch": "ccip_load_1_6",
		"commit": "ccip_load_1_6",
	}
	wg sync.WaitGroup
)

// this key only works on simulated geth chains in crib
const (
	simChainTestKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	solTestKey      = "57qbvFjTChfNwQxqkFZwjHp7xYoPZa7f9ow6GA59msfCH1g6onSjKUTrrLp4w1nAwbwQuit8YgJJ2AwT9BSwownC"
	aptosTestKey    = "0x906b8a983b434318ca67b7eff7300f91b02744c84f87d243d2fbc3e528414366"
)

func runSafely(ops ...func()) {
	for _, op := range ops {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Recovered from panic: %v\n", r)
				}
			}()
			op()
		}()
	}
}

func SetProgramIDsSafe(state solState.CCIPChainState) {
	runSafely(
		func() {
			ccip_router.SetProgramID(state.Router)
		},
		func() {
			fee_quoter.SetProgramID(state.FeeQuoter)
		},
		func() {
			ccip_offramp.SetProgramID(state.OffRamp)
		},
		func() {
			for _, key := range state.LockReleaseTokenPools {
				lockrelease_token_pool.SetProgramID(key)
			}
		},
		func() {
			for _, key := range state.BurnMintTokenPools {
				burnmint_token_pool.SetProgramID(key)
			}
		},
	)
}

// step 1: setup
// Parse the test config
// step 2: subscribe
// Create event subscribers in src and dest
// step 3: load
// Use wasp to initiate load
// step 4: teardown
// wait for ccip to finish, push remaining data
func TestCCIPLoad_RPS(t *testing.T) {
	lggr := logger.Test(t)
	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(crib.DeployerKeys{
		EVMKey:   simChainTestKey,
		SolKey:   solTestKey,
		AptosKey: aptosTestKey,
	})
	require.NoError(t, err)
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)
	userOverrides.Validate(t, env)

	ctx := env.GetContext()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	destinationChains := env.BlockChains.ListChainSelectors()[:*userOverrides.NumDestinationChains]
	evmChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilyEVM))
	solChains := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilySolana))

	// initialize the block time for each chain
	blockTimes := make(map[uint64]uint64)
	for _, cs := range evmChains {
		// Get the first block
		block1, err := env.BlockChains.EVMChains()[cs].Client.HeaderByNumber(context.Background(), big.NewInt(1))
		require.NoError(t, err)
		time1 := time.Unix(int64(block1.Time), 0) //nolint:gosec // G115

		// Get the second block
		block2, err := env.BlockChains.EVMChains()[cs].Client.HeaderByNumber(context.Background(), big.NewInt(2))
		require.NoError(t, err)
		time2 := time.Unix(int64(block2.Time), 0) //nolint:gosec // G115

		blockTimeDiff := int64(time2.Sub(time1))
		blockNumberDiff := new(big.Int).Sub(block2.Number, block1.Number).Int64()
		blockTime := blockTimeDiff / blockNumberDiff / int64(time.Second)
		blockTimes[cs] = uint64(blockTime) //nolint:gosec // G115
		lggr.Infow("Chain block time", "chainSelector", cs, "blockTime", blockTime)
	}
	for _, cs := range solChains {
		blockTimes[cs] = 0
	}

	// initialize additional accounts on EVM, we need more accounts to avoid nonce issues
	// Solana doesn't have a nonce concept so we just use a single account for all chains
	evmSenders, err := fundAdditionalKeys(lggr, *env, destinationChains)
	require.NoError(t, err)

	// Keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	state, err := stateview.LoadOnchainState(*env)
	require.NoError(t, err)

	for chainSel := range state.SolChains {
		SetProgramIDsSafe(state.SolChains[chainSel])
		err := prepSolAccount(
			ctx,
			t,
			lggr,
			env,
			state,
			chainSel)
		require.NoError(t, err)
	}

	finalSeqNrCommitChannels := make(map[uint64]chan finalSeqNrReport)
	finalSeqNrExecChannels := make(map[uint64]chan finalSeqNrReport)
	loadFinished := make(chan struct{})

	mm := NewMetricsManager(t, env.Logger, userOverrides, blockTimes)
	go mm.Start(ctx)

	// gunMap holds a destinationGun for every enabled destination chain
	gunMap := make(map[uint64]*DestinationGun)
	p := wasp.NewProfile()

	// Discover lanes from deployed state
	laneConfig := &crib.LaneConfiguration{}
	err = laneConfig.DiscoverLanesFromDeployedState(*env, &state)
	require.NoError(t, err)
	laneConfig.LogLaneConfigInfo(lggr)

	// potential source chains need a subscription
	for _, cs := range env.BlockChains.ListChainSelectors() {
		destChains := laneConfig.GetDestinationChainsForSource(cs)
		selectorFamily, err := selectors.GetSelectorFamily(cs)
		require.NoError(t, err)
		wg.Add(1)
		switch selectorFamily {
		case selectors.FamilyEVM:
			latesthdr, err := env.BlockChains.EVMChains()[cs].Client.HeaderByNumber(ctx, nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[cs] = &block
			go subscribeTransmitEvents(
				ctx,
				lggr,
				state.Chains[cs].OnRamp,
				destChains,
				startBlocks[cs],
				cs,
				loadFinished,
				env.BlockChains.EVMChains()[cs].Client,
				&wg,
				mm.InputChan,
				finalSeqNrCommitChannels,
				finalSeqNrExecChannels)
		case selectors.FamilySolana:
			client := env.BlockChains.SolanaChains()[cs].Client
			block, err := client.GetBlockHeight(ctx, solrpc.CommitmentConfirmed)
			require.NoError(t, err)
			startBlocks[cs] = &block
			go subscribeSolTransmitEvents(
				ctx,
				lggr,
				state.SolChains[cs].Router,
				destChains,
				block,
				cs,
				loadFinished,
				client,
				&wg,
				mm.InputChan,
				finalSeqNrCommitChannels,
				finalSeqNrExecChannels)
		}
	}

	evmSourceKeys := make(map[uint64]map[uint64]*bind.TransactOpts)
	solSourceKeys := make(map[uint64]*solana.PrivateKey)
	var mu sync.Mutex

	for ind, cs := range destinationChains {
		srcChains := laneConfig.GetSourceChainsForDestination(cs)

		// Initialize the map for this destination
		evmSourceKeys[cs] = make(map[uint64]*bind.TransactOpts)

		for _, src := range srcChains {
			selFamily, err := selectors.GetSelectorFamily(src)
			if err != nil {
				lggr.Errorw("Failed to get selector family", "chainSelector", src, "error", err)
				continue
			}
			mu.Lock()
			switch selFamily {
			case selectors.FamilyEVM:
				// Check if we have enough senders for this source chain
				if ind < len(evmSenders[src]) {
					evmSourceKeys[cs][src] = evmSenders[src][ind]
				} else {
					lggr.Errorw("Not enough EVM senders for source chain",
						"sourceChain", src,
						"destinationChain", cs,
						"requiredIndex", ind,
						"availableSenders", len(evmSenders[src]))
				}
			case selectors.FamilySolana:
				if _, exists := solSourceKeys[src]; !exists {
					solSourceKeys[src] = env.BlockChains.SolanaChains()[src].DeployerKey
				}
			}
			mu.Unlock()
		}
	}

	// confirmed dest chains need a subscription
	for _, cs := range destinationChains {
		srcChains := laneConfig.GetSourceChainsForDestination(cs)

		g := new(errgroup.Group)
		for _, src := range srcChains {
			g.Go(func() error {
				selFamily, err := selectors.GetSelectorFamily(src)
				require.NoError(t, err)
				switch selFamily {
				case selectors.FamilyEVM:
					return prepareAccountToSendLink(
						lggr,
						state,
						*env,
						src,
						evmSourceKeys[cs][src],
					)
				default:
					return nil
				}
			})
		}
		require.NoError(t, g.Wait())

		finalSeqNrCommitChannels[cs] = make(chan finalSeqNrReport)
		finalSeqNrExecChannels[cs] = make(chan finalSeqNrReport)

		selectorFamily, err := selectors.GetSelectorFamily(cs)
		require.NoError(t, err)
		switch selectorFamily {
		case selectors.FamilyEVM:
			gunMap[cs], err = NewDestinationGun(
				env.Logger,
				cs,
				*env,
				&state,
				state.Chains[cs].Receiver.Address().Bytes(),
				userOverrides,
				evmSourceKeys[cs],
				solSourceKeys,
				mm.InputChan,
				srcChains,
			)
			if err != nil {
				lggr.Errorw("Failed to initialize DestinationGun for", "chainSelector", cs, "error", err)
				t.Fatal(err)
			}
			wg.Add(2)
			go subscribeCommitEvents(
				ctx,
				lggr,
				state.Chains[cs].OffRamp,
				srcChains,
				startBlocks[cs],
				cs,
				env.BlockChains.EVMChains()[cs].Client,
				finalSeqNrCommitChannels[cs],
				&wg,
				mm.InputChan)
			go subscribeExecutionEvents(
				ctx,
				lggr,
				state.Chains[cs].OffRamp,
				srcChains,
				startBlocks[cs],
				cs,
				env.BlockChains.EVMChains()[cs].Client,
				finalSeqNrExecChannels[cs],
				&wg,
				mm.InputChan)

			// error watchers
			go subscribeSkippedIncorrectNonce(
				ctx,
				cs,
				state.Chains[cs].NonceManager,
				lggr)

			go subscribeAlreadyExecuted(
				ctx,
				cs,
				state.Chains[cs].OffRamp,
				lggr)
		case selectors.FamilySolana:

			gunMap[cs], err = NewDestinationGun(
				env.Logger,
				cs,
				*env,
				&state,
				state.SolChains[cs].Receiver.Bytes(),
				userOverrides,
				evmSourceKeys[cs],
				solSourceKeys,
				mm.InputChan,
				srcChains,
			)
			if err != nil {
				lggr.Errorw("Failed to initialize DestinationGun for", "chainSelector", cs, "error", err)
				t.Fatal(err)
			}
			wg.Add(2)
			go subscribeSolCommitEvents(
				ctx,
				lggr,
				state.SolChains[cs].OffRamp,
				srcChains,
				*startBlocks[cs],
				cs,
				env.BlockChains.SolanaChains()[cs].Client,
				finalSeqNrCommitChannels[cs],
				&wg,
				mm.InputChan)

			go subscribeSolExecutionEvents(
				ctx,
				lggr,
				state.SolChains[cs].OffRamp,
				srcChains,
				*startBlocks[cs],
				cs,
				env.BlockChains.SolanaChains()[cs].Client,
				finalSeqNrExecChannels[cs],
				&wg,
				mm.InputChan)
		}
	}

	requestFrequency, err := time.ParseDuration(*userOverrides.RequestFrequency)
	require.NoError(t, err)

	for _, gun := range gunMap {
		p.Add(wasp.NewGenerator(&wasp.Config{
			T:           t,
			GenName:     "ccipLoad",
			LoadType:    wasp.RPS,
			CallTimeout: userOverrides.GetLoadDuration(),
			// 1 request per second for n seconds
			Schedule: wasp.Plain(1, userOverrides.GetLoadDuration()),
			// limit requests to 1 per duration
			RateLimitUnitDuration: requestFrequency,
			// will need to be divided by number of chains
			// this schedule is per generator
			// in this example, it would be 1 request per 5seconds per generator (dest chain)
			// so if there are 3 generators, it would be 3 requests per 5 seconds over the network
			Gun:        gun,
			Labels:     CommonTestLabels,
			LokiConfig: wasp.NewEnvLokiConfig(),
			// use the same loki client using `NewLokiClient` with the same config for sending events
		}))
	}

	switch config.CCIP.Load.ChaosMode {
	case ccip.ChaosModeTypeRPCLatency:
		go runRealisticRPCLatencySuite(t,
			config.CCIP.Load.GetLoadDuration()+userOverrides.GetTimeoutDuration(),
			config.CCIP.Load.GetRPCLatency(),
			config.CCIP.Load.GetRPCJitter(),
			len(evmChains),
		)
	case ccip.ChaosModeTypeFull:
		go runFullChaosSuite(t)
	case ccip.ChaosModeNone:
	}

	_, err = p.Run(true)
	require.NoError(t, err)
	// wait some duration so that transmits can happen
	go func() {
		time.Sleep(tickerDuration)
		close(loadFinished)
	}()

	// after load is finished, wait for a "timeout duration" before considering that messages are timed out
	timeout := userOverrides.GetTimeoutDuration()
	if timeout != 0 {
		testTimer := time.NewTimer(timeout)
		go func() {
			<-testTimer.C
			cancel()
			t.Fail()
		}()
	}

	wg.Wait()
	lggr.Infow("closed event subscribers")
}
