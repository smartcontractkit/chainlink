# How to run Log Poller's tests

## Limitations
* currently they can only be run in Docker, not in Kubernetes
* when using `looped` runner it's not possible to directly control execution time
* WASP's `gun` implementation is imperfect in terms of generated load

## Configuration
Due to unfinished migration to TOML config tests use a mixed configuration approach:
* network, RPC endpoints, funding keys, etc need to be provided by env vars
* test-specific configuration can be provided by TOML file or via a `Config` struct (to which TOML is parsed anyway) additionally some of it can be overridden by env vars (for ease of use in CI)
** smoke tests use the programmatical approach
** load test uses the TOML approach

## Approximated test scenario
Different tests might have slightly modified scenarios, but generally they follow this pattern:
* start CL nodes
* setup OCR
* upload Automation Registry 2.1
* deploy UpKeep Consumers
* deploy test contracts
* register filters for test contracts
* make sure all CL nodes have filters registered
* emit test logs
* wait for log poller to finalise last block in which logs were emitted
** block number is determined either by finality tag or fixed finality depth depending on network configuration
* wait for all CL nodes to have expected log count
* compare logs that present in the EVM node with logs in CL nodes

All of the checks use fluent waits.

### Required env vars
* `CHAINLINK_IMAGE`
* `CHAINLINK_VERSION`
* `SELECTED_NETWORKS`

### Env vars required for live testnet tests
* `EVM_WS_URL` -- RPC websocket
* `EVM_HTTP_URL` -- RPC HTTP
* `EVM_KEYS` -- private keys used for funding

Since on live testnets we are using existing and canonical LINK contracts funding keys need to contain enough LINK to pay for the test. There's an automated check that fails during setup if there's not enough LINK. Approximately `9 LINK` is required for each UpKeep contract test uses to register a `LogTrigger`. Test contract emits 3 types of events and unless configured otherwise (programmatically!) all of them will be used, which means that due to Automation's limitation we need to register a separate `LogTrigger` for each event type for each contract. So if you want to test with 100 contracts, then you'd need to register 300 UpKeep contracts and thus your funding address needs to have at least 2700 LINK.

### Programmatical config
There are two load generators available:
* `looped` -- it's a simple generator that just loops over all contracts and emits events at random intervals
* `wasp` -- based on WASP load testing tool, it's more sophisticated and allows to control execution time

#### Looped config
```
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,                              # number of test contracts to deploy
			EventsPerTx:    4,                              # number of events to emit in a single transaction
			UseFinalityTag: false,                          # if set to true then Log Poller will use finality tag returned by chain, when determining last finalised block (won't work on a simulated network, it requires eth2)
		},
		LoopedConfig: &logpoller.LoopedConfig{
			ContractConfig: logpoller.ContractConfig{
				ExecutionCount: 100,                        # number of times each contract will be called
			},
			FuzzConfig: logpoller.FuzzConfig{
				MinEmitWaitTimeMs: 200,                     # minimum number of milliseconds to wait before emitting events
				MaxEmitWaitTimeMs: 500,                     # maximum number of milliseconds to wait before emitting events
			},
		},
	}

    eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {     # modify that function to emit only logs you want
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit
```

Remember that final number of events emitted will be `Contracts * EventsPerTx * ExecutionCount * len(eventToEmit)`. And that that last number by default is equal to `3` (that's because we want to emit different event types, not just one). You can change that by overriding `EventsToEmit` field.

#### WASP config
```
	cfg := logpoller.Config{
		General: &logpoller.General{
			Generator:      logpoller.GeneratorType_Looped,
			Contracts:      2,
			EventsPerTx:    4,
			UseFinalityTag: false,
		},
		Wasp: &logpoller.WaspConfig{
			Load: &logpoller.Load{
				RPS:                   10,                                              # requests per second
				LPS:                   0,                                               # logs per second 
				RateLimitUnitDuration: models.MustNewDuration(5 * time.Minutes),        # for how long the load should be limited (ramp-up period)
				Duration:              models.MustNewDuration(5 * time.Minutes),        # how long to generate the load for
				CallTimeout:           models.MustNewDuration(5 * time.Minutes),        # how long to wait for a single call to finish
			},
		},
	}

    eventsToEmit := []abi.Event{}
	for _, event := range logpoller.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit
```

Remember that you cannot specify both `RPS` and `LPS`. If you want to use `LPS` then omit `RPS` field. Also remember that depending on the events you decide to emit RPS might mean 1 request or might mean 3 requests (if you go with the default `EventsToEmit`).

For other nuances do check [gun.go][integration-tests/universal/log_poller/gun.go].

### TOML config
That config follows the same structure as programmatical config shown above.

Sample config: [config.toml](integration-tests/load/log_poller/config.toml)

Use this snippet instead of creating the `Config` struct programmatically:
```
	cfg, err := lp_helpers.ReadConfig(lp_helpers.DefaultConfigFilename)
	require.NoError(t, err)
```

And remember to add events you want emit:
```
	eventsToEmit := []abi.Event{}
	for _, event := range lp_helpers.EmitterABI.Events {
		eventsToEmit = append(eventsToEmit, event)
	}

	cfg.General.EventsToEmit = eventsToEmit
```

### Timeouts
Various checks inside the tests have hardcoded timeouts, which might not be suitable for your execution parameters, for example if you decided to emit 1M logs, then waiting for all of them to be indexed for `1m` might not be enough. Remember to adjust them accordingly.

Sample snippet:
```
	gom.Eventually(func(g gomega.Gomega) {
		logCountMatches, err := clNodesHaveExpectedLogCount(startBlock, endBlock, testEnv.EVMClient.GetChainID(), totalLogsEmitted, expectedFilters, l, coreLogger, testEnv.ClCluster)
		if err != nil {
			l.Warn().Err(err).Msg("Error checking if CL nodes have expected log count. Retrying...")
		}
		g.Expect(logCountMatches).To(gomega.BeTrue(), "Not all CL nodes have expected log count")
	}, "1m", "30s").Should(gomega.Succeed()) # 1m is the timeout for all nodes to have expected log count
```

## Tests
* [Load](integration-tests/load/log_poller/log_poller_test.go)
* [Smoke](integration-tests/smoke/log_poller/log_poller_test.go)

## Running tests
After setting all the environment variables you can run the test with:
```
# run in the root folder of chainlink repo
go test -v -test.timeout=2700s -run TestLogPollerReplay integration-tests/smoke/log_poller_test.go
```

Remember to adjust test timeout accordingly to match expected duration.


## Github Actions
If all of that seems too complicated use this [on-demand workflow](https://github.com/smartcontractkit/chainlink/actions/workflows/on-demand-log-poller.yml).

Execution time here is an approximation, so depending on network conditions it might be slightly longer or shorter.