package ccip

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/smartcontractkit/chainlink/integration-tests/testconfig/ccip"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
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
const simChainTestKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

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
	cribDeployOutput, err := cribEnv.GetConfig(simChainTestKey)
	require.NoError(t, err)
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)
	userOverrides.Validate(t, env)

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

	// initialize additional accounts on other chains
	transmitKeys, err := fundAdditionalKeys(lggr, *env, env.AllChainSelectors()[:*userOverrides.NumDestinationChains])
	// todo: fund keys on solana
	require.NoError(t, err)
	// todo: defer returning funds

	// Keep track of the block number for each chain so that event subscription can be done from that block.
	startBlocks := make(map[uint64]*uint64)
	state, err := ccipchangeset.LoadOnchainState(*env)
	require.NoError(t, err)

	finalSeqNrCommitChannels := make(map[uint64]chan finalSeqNrReport)
	finalSeqNrExecChannels := make(map[uint64]chan finalSeqNrReport)
	loadFinished := make(chan struct{})

	mm := NewMetricsManager(t, env.Logger, userOverrides, blockTimes)
	go mm.Start(ctx)

	// gunMap holds a destinationGun for every enabled destination chain
	gunMap := make(map[uint64]*DestinationGun)
	p := wasp.NewProfile()

	// potential source chains need a subscription
	for _, cs := range env.AllChainSelectors() {
		latesthdr, err := env.Chains[cs].Client.HeaderByNumber(ctx, nil)
		require.NoError(t, err)
		block := latesthdr.Number.Uint64()
		startBlocks[cs] = &block
		other := env.AllChainSelectorsExcluding([]uint64{cs})
		wg.Add(1)
		go subscribeTransmitEvents(
			ctx,
			lggr,
			state.Chains[cs].OnRamp,
			other,
			startBlocks[cs],
			cs,
			loadFinished,
			env.Chains[cs].Client,
			&wg,
			mm.InputChan,
			finalSeqNrCommitChannels,
			finalSeqNrExecChannels)
	}

	// confirmed dest chains need a subscription
	for ind := range *userOverrides.NumDestinationChains {
		cs := env.AllChainSelectors()[ind]

		evmSourceKeys := make(map[uint64]*bind.TransactOpts)
		solanaSourceKeys := make(map[uint64]*solana.PrivateKey)
		// todo: make solana source keys
		other := env.AllChainSelectorsExcluding([]uint64{cs})
		var mu sync.Mutex
		var wg2 sync.WaitGroup
		wg2.Add(len(other))
		for _, src := range other {
			go func(src uint64) {
				defer wg2.Done()
				mu.Lock()
				evmSourceKeys[src] = transmitKeys[src][ind]
				mu.Unlock()
				assert.NoError(t, prepareAccountToSendLink(
					t,
					state,
					*env,
					src,
					evmSourceKeys[src],
				))
			}(src)
		}
		wg2.Wait()

		gunMap[cs], err = NewDestinationGun(
			env.Logger,
			cs,
			*env,
			&state,
			state.Chains[cs].Receiver.Address(),
			userOverrides,
			evmSourceKeys,
			solanaSourceKeys,
			ind,
			mm.InputChan,
		)
		if err != nil {
			lggr.Errorw("Failed to initialize DestinationGun for", "chainSelector", cs, "error", err)
			t.Fatal(err)
		}

		finalSeqNrCommitChannels[cs] = make(chan finalSeqNrReport)
		finalSeqNrExecChannels[cs] = make(chan finalSeqNrReport)

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
