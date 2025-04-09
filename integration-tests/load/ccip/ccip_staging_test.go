package ccip

import (
	"crypto/ecdsa"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/wasp"
	ccipchangeset "github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/crib"
	tc "github.com/smartcontractkit/chainlink/integration-tests/testconfig"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestStaging_CCIP_Load(t *testing.T) {
	lggr := logger.Test(t)

	sourceKey := simChainTestKey

	// get user defined configurations
	config, err := tc.GetConfig([]string{"Load"}, tc.CCIP)
	require.NoError(t, err)
	userOverrides := config.CCIP.Load

	// generate environment from crib-produced files
	cribEnv := crib.NewDevspaceEnvFromStateDir(lggr, *userOverrides.CribEnvDirectory)
	cribDeployOutput, err := cribEnv.GetConfig(sourceKey)
	require.NoError(t, err)
	env, err := crib.NewDeployEnvironmentFromCribOutput(lggr, cribDeployOutput)
	require.NoError(t, err)
	require.NotNil(t, env)
	userOverrides.Validate(t, env)
	state, err := ccipchangeset.LoadOnchainState(*env)
	require.NoError(t, err)

	// initialize additional accounts on other chains
	transmitKeys, err := fundAdditionalKeys(lggr, *env, env.AllChainSelectors()[:*userOverrides.NumDestinationChains])
	require.NoError(t, err)

	// gunMap holds a destinationGun for every enabled destination chain
	gunMap := make(map[uint64]*DestinationGun)
	p := wasp.NewProfile()
	for ind := range *userOverrides.NumDestinationChains {
		cs := env.AllChainSelectors()[ind]

		messageKeys := make(map[uint64]*bind.TransactOpts)
		other := env.AllChainSelectorsExcluding([]uint64{cs})
		var mu sync.Mutex
		var wg2 sync.WaitGroup
		wg2.Add(len(other))
		for _, src := range other {
			go func(src uint64) {
				defer wg2.Done()
				mu.Lock()
				messageKeys[src] = transmitKeys[src][ind]
				mu.Unlock()
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
			messageKeys,
			nil,
			ind,
			nil,
		)
		if err != nil {
			lggr.Errorw("Failed to initialize DestinationGun for", "chainSelector", cs, "error", err)
			t.Fatal(err)
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
	_, err = p.Run(true)
	require.NoError(t, err)

	lggr.Info("Load test complete, returning funds")
	// return funds to source address at the end of the test
	sourcePk, err := crypto.HexToECDSA(sourceKey)
	if err != nil {
		lggr.Errorw("could not return funds to source address")
	}
	// Derive the public key
	publicKey := sourcePk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		lggr.Errorw("could not return funds to source address")
	}

	// Get the address from the public key
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	err = reclaimFunds(lggr, *env, transmitKeys, common.HexToAddress(address))
	if err != nil {
		lggr.Errorw(err.Error())
	}
}
