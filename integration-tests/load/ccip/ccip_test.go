package ccip

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	solrpc "github.com/gagliardetto/solana-go/rpc"
	selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/burnmint_token_pool"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/ccip_router"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/solana/gobindings/lockrelease_token_pool"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	"github.com/smartcontractkit/chainlink-testing-framework/wasp"

	"github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"
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

func SetProgramIDsSafe(state changeset.SolCCIPChainState) {
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
			lockrelease_token_pool.SetProgramID(state.LockReleaseTokenPool)
		},
		func() {
			burnmint_token_pool.SetProgramID(state.BurnMintTokenPool)
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
	ctx, cancel := context.WithCancel(tests.Context(t))
	defer cancel()

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey, solTestKey)
	require.NoError(t, err)
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)
	userOverrides.Validate(t, env)

	destinationChains := env.AllChainSelectorsAllFamilies()[:*userOverrides.NumDestinationChains]
	// initialize the block time for each chain
	blockTimes := make(map[uint64]uint64)
	for _, cs := range env.AllChainSelectors() {
		// Get the first block
		block1, err := env.Chains[cs].Client.HeaderByNumber(context.Background(), big.NewInt(1))
		require.NoError(t, err)
		time1 := time.Unix(int64(block1.Time), 0) //nolint:gosec // G115

		// Get the second block
		block2, err := env.Chains[cs].Client.HeaderByNumber(context.Background(), big.NewInt(2))
		require.NoError(t, err)
		time2 := time.Unix(int64(block2.Time), 0) //nolint:gosec // G115

		blockTimeDiff := int64(time2.Sub(time1))
		blockNumberDiff := new(big.Int).Sub(block2.Number, block1.Number).Int64()
		blockTime := blockTimeDiff / blockNumberDiff / int64(time.Second)
		blockTimes[cs] = uint64(blockTime) //nolint:gosec // G115
		lggr.Infow("Chain block time", "chainSelector", cs, "blockTime", blockTime)
	}
	for _, cs := range env.AllChainSelectorsSolana() {
		blockTimes[cs] = 0
	}

	// initialize additional accounts on other chains
	evmSenders, err := fundAdditionalKeys(lggr, *env, destinationChains)
	solanaSenders := make(map[uint64][]solana.PrivateKey)
	for _, solSel := range env.AllChainSelectorsSolana() {
		solanaSenders[solSel] = make([]solana.PrivateKey, 0, len(destinationChains))
		for range len(destinationChains) {
			newPk, err := solana.NewRandomPrivateKey()
			require.NoError(t, err)
			solanaSenders[solSel] = append(solanaSenders[solSel], newPk)
		}
	}

	require.NoError(t, err)

	// Keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	state, err := ccipchangeset.LoadOnchainState(*env)
	require.NoError(t, err)

	for chainSel := range state.SolChains {
		SetProgramIDsSafe(state.SolChains[chainSel])
		err := prepSolAccount(ctx, t, lggr, env, solanaSenders[chainSel], chainSel, state.SolChains[chainSel].Router)
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

	// potential source chains need a subscription
	for _, cs := range env.AllChainSelectorsAllFamilies() {
		otherChains := env.AllChainSelectorsAllFamiliesExcluding([]uint64{cs})
		selectorFamily, err := selectors.GetSelectorFamily(cs)
		require.NoError(t, err)
		wg.Add(1)
		switch selectorFamily {
		case selectors.FamilyEVM:
			latesthdr, err := env.Chains[cs].Client.HeaderByNumber(ctx, nil)
			require.NoError(t, err)
			block := latesthdr.Number.Uint64()
			startBlocks[cs] = &block
			go subscribeTransmitEvents(
				ctx,
				lggr,
				state.Chains[cs].OnRamp,
				otherChains,
				startBlocks[cs],
				cs,
				loadFinished,
				env.Chains[cs].Client,
				&wg,
				mm.InputChan,
				finalSeqNrCommitChannels,
				finalSeqNrExecChannels)
		case selectors.FamilySolana:
			client := env.SolChains[cs].Client
			block, err := client.GetBlockHeight(ctx, solrpc.CommitmentConfirmed)
			require.NoError(t, err)
			startBlocks[cs] = &block
			go subscribeSolTransmitEvents(
				ctx,
				lggr,
				state.SolChains[cs].Router,
				otherChains,
				block,
				cs,
				loadFinished,
				env.SolChains[cs].Client,
				&wg,
				mm.InputChan,
				finalSeqNrCommitChannels,
				finalSeqNrExecChannels)
		}
	}

	// confirmed dest chains need a subscription
	for ind, cs := range destinationChains {
		evmSourceKeys := make(map[uint64]*bind.TransactOpts)
		solSourceKeys := make(map[uint64]*solana.PrivateKey)
		other := env.AllChainSelectorsAllFamiliesExcluding([]uint64{cs})
		var mu sync.Mutex
		var wg2 sync.WaitGroup
		for _, src := range other {
			selFamily, err := selectors.GetSelectorFamily(src)
			require.NoError(t, err)
			switch selFamily {
			case selectors.FamilyEVM:
				wg2.Add(1)
				go func(src uint64) {
					defer wg2.Done()
					mu.Lock()
					evmSourceKeys[src] = evmSenders[src][ind]
					mu.Unlock()
				}(src)
			case selectors.FamilySolana:
				solSourceKeys[src] = &solanaSenders[src][ind]
			}
		}
		wg2.Wait()

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
				evmSourceKeys,
				solSourceKeys,
				ind,
				mm.InputChan,
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
				other,
				startBlocks[cs],
				cs,
				env.Chains[cs].Client,
				finalSeqNrCommitChannels[cs],
				&wg,
				mm.InputChan)
			go subscribeExecutionEvents(
				ctx,
				lggr,
				state.Chains[cs].OffRamp,
				other,
				startBlocks[cs],
				cs,
				env.Chains[cs].Client,
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
				evmSourceKeys,
				solSourceKeys,
				ind,
				mm.InputChan,
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
				other,
				*startBlocks[cs],
				cs,
				env.SolChains[cs].Client,
				finalSeqNrCommitChannels[cs],
				&wg,
				mm.InputChan)

			go subscribeSolExecutionEvents(
				ctx,
				lggr,
				state.SolChains[cs].OffRamp,
				other,
				*startBlocks[cs],
				cs,
				env.SolChains[cs].Client,
				finalSeqNrCommitChannels[cs],
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
			config.CCIP.Load.GetLoadDuration(),
			config.CCIP.Load.GetRPCLatency(),
			config.CCIP.Load.GetRPCJitter(),
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

func prepareAccountToSendLink(
	t *testing.T,
	state ccipchangeset.CCIPOnChainState,
	e deployment.Environment,
	src uint64,
	srcAccount *bind.TransactOpts) error {
	lggr := logger.Test(t)
	srcDeployer := e.Chains[src].DeployerKey
	lggr.Infow("Setting up link token", "src", src)
	srcLink := state.Chains[src].LinkToken

	lggr.Infow("Granting mint and burn roles")
	tx, err := srcLink.GrantMintAndBurnRoles(srcDeployer, srcAccount.From)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	if err != nil {
		return err
	}

	lggr.Infow("Minting transfer amounts")
	//--------------------------------------------------------------------------------------------
	tx, err = srcLink.Mint(
		srcAccount,
		srcAccount.From,
		big.NewInt(20_000),
	)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	if err != nil {
		return err
	}

	//--------------------------------------------------------------------------------------------
	lggr.Infow("Approving routers")
	// Approve the router to spend the tokens and confirm the tx's
	// To prevent having to approve the router for every transfer, we approve a sufficiently large amount
	tx, err = srcLink.Approve(srcAccount, state.Chains[src].Router.Address(), math.MaxBig256)
	_, err = deployment.ConfirmIfNoError(e.Chains[src], tx, err)
	return err
}
