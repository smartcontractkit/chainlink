# Changelog Chainlink Core

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Task definitions in v2 jobs (those with TOML specs) now support quoting strings with angle brackets (which DOT already permitted). This is particularly useful when defining JSON blobs to post to external adapters. For example:

    ``` 
    my_bridge [type=bridge name="my_bridge" requestData="{\\"hi\\": \\"hello\\"}"]
    ```
    ... can now be written as:
    ``` 
    my_bridge [type=bridge name="my_bridge" requestData=<{"hi": "hello"}>]
    ```
    Multiline strings are supported with this syntax as well:
    ``` 
    my_bridge [type=bridge
               name="my_bridge"
               requestData=<{
                   "hi": "hello",
                   "foo": "bar"
               }>]
    ```

- v2 jobs (those with TOML specs) now support variable interpolation in pipeline definitions. For example:

    ```
    fetch1    [type=bridge name="fetch"]
    parse1    [type=jsonparse path="foo,bar"]
    fetch2    [type=bridge name="fetch"]
    parse2    [type=jsonparse path="foo,bar"]
    medianize [type=median]
    submit    [type=bridge name="submit"
               requestData=<{
                              "result": $(medianize),
                              "fetchedData": [ $(parse1), $(parse2) ]
                            }>]

    fetch1 -> parse1 -> medianize
    fetch2 -> parse2 -> medianize
    medianize -> submit
    ```

    This syntax is supported by the following tasks/parameters:

    - `bridge`
        - `requestData`
    - `http`
        - `requestData`
    - `jsonparse`
        - `data` (falls back to the first input if unspecified)
    - `median`
        - `values` (falls back to the array of inputs if unspecified)
    - `multiply`
        - `input` (falls back to the first input if unspecified)
        - `times`

- Add `ETH_MAX_IN_FLIGHT_TRANSACTIONS` configuration option. This defaults to 16 and controls how many unconfirmed transactions may be in-flight at any given moment. This is set conservatively by default, node operators running many jobs on high throughput chains will probably need to increase this above the default to avoid lagging behind. However, before increasing this value, you MUST first ensure your ethereum node is configured not to ever evict local transactions that exceed this number otherwise your node may get permanently stuck. Set to 0 to disable the limit entirely (the old behaviour). Disabling this setting is not recommended.

Relevant settings for geth (and forks e.g. BSC)

```toml
[Eth.TxPool]
Locals = ["0xYourNodeAddress1", "0xYourNodeAddress2"]  # Add your node addresses here
NoLocals = false # Disabled by default but might as well make sure
Journal = "transactions.rlp" # Make sure you set a journal file
Rejournal = 3600000000000 # Default 1h, it might make sense to reduce this to e.g. 5m
PriceBump = 10 # Must be set less than or equal to chainlink's ETH_GAS_BUMP_PERCENT
AccountSlots = 16 # Highly recommended to increase this, must be greater than or equal to chainlink's ETH_MAX_IN_FLIGHT_TRANSACTIONS setting
GlobalSlots = 4096 # Increase this as necessary
AccountQueue = 64 # Increase this as necessary
GlobalQueue = 1024 # Increase this as necessary
Lifetime = 10800000000000 # Default 3h, this is probably ok, you might even consider reducing it

```

Relevant settings for parity/openethereum (and forks e.g. xDai)

NOTE: There is a bug in parity (and xDai) where occasionally local transactions are inexplicably culled. See: https://github.com/openethereum/parity-ethereum/issues/10228

Adjusting the settings below might help.

```toml
tx_queue_locals = ["0xYourNodeAddress1", "0xYourNodeAddress2"] # Add your node addresses here
tx_queue_size = 8192 # Increase this as necessary
tx_queue_per_sender = 16 # Highly recommended to increase this, must be greater than or equal to chainlink's ETH_MAX_IN_FLIGHT_TRANSACTIONS setting
tx_queue_mem_limit = 4 # In MB. Highly recommended to increase this or set to 0
tx_queue_no_early_reject = true # Recommended to set this
tx_queue_no_unfamiliar_locals = false # This is disabled by default but might as well make sure
```

### Fixed

- It is no longer required to set `DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS=true` to enable local fetches on bridge tasks. Please remove this if you had it set and no longer need it, since it introduces a slight security risk.

### Changed

- The v2 (TOML) `bridge` task's `includeInputAtKey` parameter is being deprecated in favor of variable interpolation. Please migrate your jobs to the new syntax as soon as possible.

- Chainlink no longers writes/reads eth key files to disk
- Add sensible default configuration settings for Fantom

- Rename `ETH_MAX_UNCONFIRMED_TRANSACTIONS` to `ETH_MAX_QUEUED_TRANSACTIONS`. It still performs the same function but the name was misleading and would have caused confusion with the new `ETH_MAX_IN_FLIGHT_TRANSACTIONS`.


## [0.10.7] - 2021-05-24

- If a CLI command is issued after the session has expired, and an api credentials file is found, auto login should now work.

- GasUpdater now works on RSK and xDai

- Offchain reporting jobs that have had a latest round requested can now be deleted from the UI without error

### Added

- Add `ETH_GAS_LIMIT_MULTIPLIER` configuration option, the gas limit is multiplied by this value before transmission. So a value of 1.1 will add 10% to the on chain gas limit when a transaction is submitted.

- Add `ETH_MIN_GAS_PRICE_WEI` configuration option. This defaults to 1Gwei on mainnet. Chainlink will never send a transaction at a price lower than this value.

- Add `chainlink node db migrate` for running database migrations. It's
  recommended to use this and set `MIGRATE_DATABASE=false` if you want to run
  the migrations separately outside of application startup.

### Changed

- Chainlink now automatically cleans up old eth_txes to reduce database size. By default, any eth_txes older than a week are pruned on a regular basis. It is recommended to use the default value, however the default can be overridden by setting the `ETH_TX_REAPER_THRESHOLD` env var e.g. `ETH_TX_REAPER_THRESHOLD=24h`. Reaper can be disabled entirely by setting `ETH_TX_REAPER_THRESHOLD=0`. The reaper will run on startup and again every hour (interval is configurable using `ETH_TX_REAPER_INTERVAL`).

- Heads corresponding to new blocks are now delivered in a sampled way, which is to improve 
  node performance on fast chains. The frequency is by default 1 second, and can be changed 
  by setting `ETH_HEAD_TRACKER_SAMPLING_INTERVAL` env var e.g. `ETH_HEAD_TRACKER_SAMPLING_INTERVAL=5s`.

- Database backups: default directory is now a subdirectory 'backup' of chainlink root dir, and can be changed 
  to any chosed directory by setting a new configuration value: `DATABASE_BACKUP_DIR`

## [0.10.6] - 2021-05-10

### Added

- Add `MockOracle.sol` for testing contracts

- Web job types can now be created from the operator UI as a new job. 

- See example web job spec below: 

```
type            = "web"
schemaVersion   = 1
jobID           = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F46"
observationSource = """
ds          [type=http method=GET url="http://example.com"];
ds_parse    [type=jsonparse path="data"];
ds -> ds_parse;
"""
```

- New CLI command to convert v1 flux monitor jobs (JSON) to 
v2 flux monitor jobs (TOML). Running it will archive the v1 
job and create a new v2 job. Example:
```
// Get v1 job ID:
chainlink job_specs list
// Migrate it to v2:
chainlink jobs migrate fe279ed9c36f4eef9dc1bdb7bef21264

// To undo the migration:
1. Archive the v2 job in the UI
2. Unarchive the v1 job manually in the db:
update job_specs set deleted_at = null where id = 'fe279ed9-c36f-4eef-9dc1-bdb7bef21264'
```

- Improved support for Optimism chain. Added a new boolean `OPTIMISM_GAS_FEES` configuration variable which makes a call to estimate gas before all transactions, suitable for use with Optimism's L2 chain. When this option is used `ETH_GAS_LIMIT_DEFAULT` is ignored.

- Chainlink now supports routing certain calls to the eth node over HTTP instead of websocket, when available. This has a number of advantages - HTTP is more robust and simpler than websockets, reducing complexity and allowing us to make large queries without running the risk of hitting websocket send limits. The HTTP url should point to the same node as the ETH_URL and can be specified with an env var like so: `ETH_HTTP_URL=https://my.ethereumnode.example/endpoint`.

Adding an HTTP endpoint is particularly recommended for BSC, which is hitting websocket limitations on certain queries due to its large block size.

- Support for legacy pipeline (V1 job specs) can now be turned off by setting `ENABLE_LEGACY_JOB_PIPELINE=false`. This can yield marginal performance improvements if you don't need to support the legacy JSON job spec format.

## [0.10.5] - 2021-04-26

### Added

- Add `MockOracle.sol` for testing contracts
- Cron jobs can now be created for the v2 job pipeline:
```
type            = "cron"
schemaVersion   = 1
schedule        = "*/10 * * * *"
observationSource   = """
ds          [type=http method=GET url="http://example.com"];
ds_parse    [type=jsonparse path="data"];
ds -> ds_parse;
"""
```

### Changed

- Default for `JOB_PIPELINE_REAPER_THRESHOLD` has been reduced from 1 week to 1 day to save database space. This variable controls how long past job run history for OCR is kept. To keep the old behaviour, you can set `JOB_PIPELINE_REAPER_THRESHOLD=168h`
- Removed support for the env var `JOB_PIPELINE_PARALLELISM`. 
- OCR jobs no longer show `TaskRuns` in success cases. This reduces
DB load and significantly improves the performance of archiving OCR jobs.
- Archiving OCR jobs should be 5-10x faster.

### Fixed

- Added `GAS_UPDATER_BATCH_SIZE` option to workaround `websocket: read limit exceeded` issues on BSC

- Basic support for Optimism chain: node no longer gets stuck with 'nonce too low' error if connection is lost

## [0.10.4] - 2021-04-05

### Added

- VRF Jobs now support an optional `coordinatorAddress` field that, when present, will tell the node to check the fulfillment status of any VRF request before attempting the fulfillment transaction. This will assist in the effort to run multiple nodes with one VRF key.

- Experimental: Add `DATABASE_BACKUP_MODE`, `DATABASE_BACKUP_FREQUENCY` and `DATABASE_BACKUP_URL` configuration variables

    - It's now possible to configure database backups: on node start and separately, to be run at given frequency. `DATABASE_BACKUP_MODE` enables the initial backup on node start (with one of the values: `none`, `lite`, `full` where `lite` excludes
    potentially large tables related to job runs, among others). Additionally, if `DATABASE_BACKUP_FREQUENCY` variable is set to a duration of
    at least '1m', it enables periodic backups.
    - `DATABASE_BACKUP_URL` can be optionally set to point to e.g. a database replica, in order to avoid excessive load on the main one. Example settings:
        1. `DATABASE_BACKUP_MODE="full"` and `DATABASE_BACKUP_FREQUENCY` not set, will run a full back only at the start of the node.
        2. `DATABASE_BACKUP_MODE="lite"` and `DATABASE_BACKUP_FREQUENCY="1h"` will lead to a partial backup on node start and then again a partial backup every one hour.

- Added periodic resending of eth transactions. This means that we no longer rely exclusively on gas bumping to resend unconfirmed transactions that got "lost" for whatever reason. This has two advantages:
    1. Chainlink no longer relies on gas bumping settings to ensure our transactions always end up in the mempool
    2. Chainlink will continue to resend existing transactions even in the event that heads are delayed. This is especially useful on chains like Arbitrum which have very long wait times between heads.
    - Periodic resending can be controlled using the `ETH_TX_RESEND_AFTER_THRESHOLD` env var (default 30s). Unconfirmed transactions will be resent periodically at this interval. It is recommended to leave this at the default setting, but it can be set to any [valid duration](https://golang.org/pkg/time/#ParseDuration) or to 0 to disable periodic resending.

- Logging can now be configured in the Operator UI.

- Tuned defaults for certain Eth-compatible chains

- Chainlink node now uses different sets of default values depending on the given Chain ID. Tuned configs are built-in for the following chains:
    - Ethereum Mainnet and test chains
    - Polygon (Matic)
    - BSC
    - HECO

- If you have manually set ENV vars specific to these chains, you may want to remove those and allow the node to use its configured defaults instead.

- New prometheus metric "tx_manager_num_tx_reverted" which counts the number of reverted transactions on chain.

### Fixed

- Under certain circumstances a poorly configured Explorer could delay Chainlink node startup by up to 45 seconds.

- Chainlink node now automatically sets the correct nonce on startup if you are restoring from a previous backup (manual setnextnonce is no longer necessary).

- Flux monitor jobs should now work correctly with [outlier-detection](https://github.com/smartcontractkit/external-adapters-js/tree/develop/composite/outlier-detection) and [market-closure](https://github.com/smartcontractkit/external-adapters-js/tree/develop/composite/market-closure) external adapters.

- Performance improvements to OCR job adds. Removed the pipeline_task_specs table
and added a new column `dot_id` to the pipeline_task_runs table which links a pipeline_task_run
to a dotID in the pipeline_spec.dot_dag_source.

- Fixed bug where node will occasionally submit an invalid OCR transmission which reverts with "address not authorized to sign". 

- Fixed bug where a node will sometimes double submit on runlog jobs causing reverted transactions on-chain


## [0.10.3] - 2021-03-22

### Added

- Add `STATS_PUSHER_LOGGING` to toggle stats pusher raw message logging (DEBUG
  level).

- Add `ADMIN_CREDENTIALS_FILE` configuration variable

This variable defaults to `$ROOT/apicredentials` and when defined / the
file exists, any command using the CLI that requires authentication will use it
to automatically log in.

- Add `ETH_MAX_UNCONFIRMED_TRANSACTIONS` configuration variable

Chainlink node now has a maximum number of unconfirmed transactions that
may be in flight at any one time (per key).

If this limit is reached, further attempts to send transactions will fail
and the relevant job will be marked as failed.

Jobs will continue to fail until at least one transaction is confirmed
and the queue size is reduced. This is introduced as a sanity limit to
prevent unbounded sending of transactions e.g. in the case that the eth
node is failing to broadcast to the network.

The default is set to 500 which considered high enough that it should
never be reached under normal operation. This limit can be changed
by setting the `ETH_MAX_UNCONFIRMED_TRANSACTIONS` environment variable.

- Support requestNewRound in libocr

requestNewRound enables dedicated requesters to request a fresh report to
be sent to the contract right away regardless of heartbeat or deviation.

- New prometheus metric:

```
Name: "head_tracker_eth_connection_errors",
Help: "The total number of eth node connection errors",
```

- Gas bumping can now be disabled by setting `ETH_GAS_BUMP_THRESHOLD=0`

- Support for arbitrum

### Fixed

- Improved handling of the case where we exceed the configured TX fee cap in geth.

Node will now fatally error jobs if the total transaction costs exceeds the
configured cap (default 1 Eth). Also, it will no longer continue to bump gas on
transactions that started hitting this limit and instead continue to resubmit
at the highest price that worked.

Node operators should check their geth nodes and remove this cap if configured,
you can do this by running your geth node with `--rpc.gascap=0
--rpc.txfeecap=0` or setting these values in your config toml.

- Make head backfill asynchronous. This should eliminate some harmless but
  annoying errors related to backfilling heads, logged on startup and
  occasionally during normal operation on fast chains like Kovan.

- Improvements to the GasUpdater

Various efficiency and correctness improvements have been made to the
GasUpdater. It places less load on the ethereum node and now features re-org
detection.

Most notably, GasUpdater no longer takes a 24 block delay to "warm up" on
application start and instead loads all relevant block history immediately.
This means that the application gas price will always be updated correctly
after reboot before the first transaction is ever sent, eliminating the previous
scenario where the node could send underpriced or overpriced transactions for a
period after a reboot, until the gas updater caught up.

### Changed

- Bump `ORM_MAX_OPEN_CONNS` default from 10 to 20
- Bump `ORM_MAX_IDLE_CONNS` default from 5 to 10

Each Chainlink node will now use a maximum of 23 database connections (up from previous max of 13). Make sure your postgres database is tuned accordingly, especially if you are running multiple Chainlink nodes on a single database. If you find yourself hitting connection limits, you can consider reducing `ORM_MAX_OPEN_CONNS` but this may result in degraded performance.

- The global env var `JOB_PIPELINE_MAX_TASK_DURATION` is no longer supported
for OCR jobs.

## [0.10.2] - 2021-02-26

### Fixed

- Add contexts so that database queries timeout when necessary.
- Use manual updates instead of gorm update associations.

## [0.10.1] - 2021-02-25

### Fixed

- Prevent autosaving Task Spec on when Task Runs are saved to lower database load.

## [0.10.0] - 2021-02-22

### Fixed

- Fix a case where archiving jobs could try to delete it from the external initiator even if the job was not an EI job.
- Improved performance of the transaction manager by fetching receipts in
  batches. This should help prevent the node from getting stuck when processing
  large numbers of OCR jobs.
- Fixed a fluxmonitor job bug where submitting a value outside the acceptable range would stall the job
  permanently. Now a job spec error will be thrown if the polled answer is outside the
  acceptable range and no ethtx will be submitted. As additional protection, we also now
  check the receipts of the ethtx's and if they were reverted, we mark the ethtx task as failed.

### Breaking

- Squashed migrations into a single 1_initial migration. If you were running a version
  older than 0.9.10, you need to upgrade to 0.9.10 first before upgrading to the next
  version so that the migrations are run.

### Added

- A new Operator UI feature that visualize JSON and TOML job spec tasks on a 'New Job' page.

## [0.9.10] - 2021-01-30

### Fixed

- Fixed a UI bug with fluxmonitor jobs where initiator params were bunched up.
- Improved performance of OCR jobs to reduce database load. OCR jobs now run with unlimited parallelism and are not affected by `JOB_PIPELINE_PARALLELISM`.

### Added

- A new env var `JOB_PIPELINE_MAX_RUN_DURATION` has been added which controls maximum duration of the total run.

## [0.9.9] - 2021-01-18

### Added

- New CLI commands for key management:
  - `chainlink keys eth import`
  - `chainlink keys eth export`
  - `chainlink keys eth delete`
- All keys other than VRF keys now share the same password. If you have OCR, P2P, and ETH keys encrypted with different passwords, re-insert them into your DB encrypted with the same password prior to upgrading.

### Fixed

- Fixed reading of function selector values in DB.
- Support for bignums encoded in CBOR
- Silence spurious `Job spawner ORM attempted to claim locally-claimed job` warnings
- OCR now drops transmissions instead of queueing them if the node is out of Ether
- Fixed a long-standing issue where standby nodes would hold transactions open forever while waiting for a lock. This was preventing postgres from running necessary cleanup operations, resulting in bad database performance. Any node operators running standby failover chainlink nodes should see major database performance improvements with this release and may be able to reduce the size of their database instances.
- Fixed an issue where expired session tokens in operator UI would cause a large number of reqeusts to be sent to the node, resulting in a temporary rate-limit and 429 errors.
- Fixed issue whereby http client could leave too many open file descriptors

### Changed

- Key-related API endpoints have changed. All key-related commands are now namespaced under `/v2/keys/...`, and are standardized across key types.
- All key deletion commands now perform a soft-delete (i.e. archive) by default. A special CLI flag or query string parameter must be provided to hard-delete a key.
- Node now supports multiple OCR jobs sharing the same peer ID. If you have more than one key in your database, you must now specify `P2P_PEER_ID` to indicate which key to use.
- `DATABASE_TIMEOUT` is now set to 0 by default, so that nodes will wait forever for a lock. If you already have `DATABASE_TIMEOUT=0` set explicitly in your env (most node operators) then you don't need to do anything. If you didn't have it set, and you want to keep the old default behaviour where a node exits shortly if it can't get a lock, you can manually set `DATABASE_TIMEOUT=500ms` in your env.
- OCR bootstrap node no longer sends telemetry to the endpoint specified in the OCR job spec under `MonitoringEndpoint`.

## [0.9.8] - 2020-12-17

### Fixed

- An issue where the node would emit warnings on startup for fluxmonitor contracts

## [0.9.7] - 2020-12-14

### Added

- OCR bootstrap node now sends telemetry to the endpoint specified in the OCR job spec under `MonitoringEndpoint`.
- Adds "Account addresses" table to the `/keys` page.

### Changed

- Old jobs now allow duplicate job names. Also, if the name field is empty we no longer generate a name.
- Removes broken `ACCOUNT_ADDRESS` field from `/config` page.

### Fixed

- Brings `/runs` tab back to the operator UI.
- Signs out a user from operator UI on authentication error.

#### BREAKING CHANGES

- Commands for creating/managing legacy jobs and OCR jobs have changed, to reduce confusion and accomodate additional types of jobs using the new pipeline.

#### V1 jobs

`jobs archive` => `job_specs archive`
`jobs create` => `job_specs create`
`jobs list` => `job_specs list`
`jobs show` => `job_specs show`

#### V2 jobs (currently only applies to OCR)

`jobs createocr` => `jobs create`
`jobs deletev2` => `jobs delete`
`jobs run` => `jobs run`

## [0.9.6] - 2020-11-23

- OCR pipeline specs can now be configured on a per-task basis to allow unrestricted network access for http tasks. Example like so:

```
ds1          [type=http method=GET url="http://example.com" allowunrestrictednetworkaccess="true"];
ds1_parse    [type=jsonparse path="USD" lax="true"];
ds1_multiply [type=multiply times=100];
ds1 -> ds1_parse -> ds1_multiply;
```

- New prometheus metrics as follows:

```
Name: "pipeline_run_errors",
Help: "Number of errors for each pipeline spec",

Name: "pipeline_run_total_time_to_completion",
Help: "How long each pipeline run took to finish (from the moment it was created)",

Name: "pipeline_tasks_total_finished",
Help: "The total number of pipline tasks which have finished",

Name: "pipeline_task_execution_time",
Help: "How long each pipeline task took to execute",

Name: "pipeline_task_http_fetch_time",
Help: "Time taken to fully execute the HTTP request",

Name: "pipeline_task_http_response_body_size",
Help: "Size (in bytes) of the HTTP response body",

Name: "pipeline_runs_queued",
Help: "The total number of pipline runs that are awaiting execution",

Name: "pipeline_task_runs_queued",
Help: "The total number of pipline task runs that are awaiting execution",
```

### Changed

Numerous key-related UX improvements:

- All key-related commands have been consolidated under the `chainlink keys` subcommand:
  - `chainlink createextrakey` => `chainlink keys eth create`
  - `chainlink admin info` => `chainlink keys eth list`
  - `chainlink node p2p [create|list|delete]` => `chainlink keys p2p [create|list|delete]`
  - `chainlink node ocr [create|list|delete]` => `chainlink keys ocr [create|list|delete]`
  - `chainlink node vrf [create|list|delete]` => `chainlink keys vrf [create|list|delete]`
- Deleting OCR key bundles and P2P key bundles now archives them (i.e., soft delete) so that they can be recovered if needed. If you want to hard delete a key, pass the new `--hard` flag to the command, e.g. `chainlink keys p2p delete --hard 6`.
- Output from ETH/OCR/P2P/VRF key CLI commands now renders consistently.
- Deleting an OCR/P2P/VRF key now requires confirmation from the user. To skip confirmation (e.g. in shell scripts), pass `--yes` or `-y`.
- The `--ocrpassword` flag has been removed. OCR/P2P keys now share the same password at the ETH key (i.e., the password specified with the `--password` flag).

Misc:

- Two new env variables are added `P2P_ANNOUNCE_IP` and `P2P_ANNOUNCE_PORT` which allow node operators to override locally detected values for the chainlink node's externally reachable IP/port.
- `OCR_LISTEN_IP` and `OCR_LISTEN_PORT` have been renamed to `P2P_LISTEN_IP` and `P2P_LISTEN_PORT` for consistency.
- Support for adding a job with the same name as one that was deleted.

### Fixed

- Fixed an issue where the HTTP adapter would send an empty body on retries.
- Changed the default `JOB_PIPELINE_REAPER_THRESHOLD` value from `7d` to `168h` (hours are the highest time unit allowed by `time.Duration`).

## [0.9.5] - 2020-11-12

### Changed

- Updated from Go 1.15.4 to 1.15.5.

## [0.9.4] - 2020-11-04

### Fixed

- Hotfix to fix an issue with httpget adapter

## [0.9.3] - 2020-11-02

### Added

- Add new subcommand `node hard-reset` which is used to remove all state for unstarted and pending job runs from the database.

### Changed

- Chainlink now requires Postgres >= 11.x. Previously this was a recommendation, this is now a hard requirement. Migrations will fail if run on an older version of Postgres.
- Database improvements that greatly reduced the number of open Postgres connections
- Operator UI /jobs page is now searchable
- Jobs now accept a name field in the jobspecs

## [0.9.2] - 2020-10-15

### Added

- Bulletproof transaction manager enabled by default
- Fluxmonitor support enabled by default

### Fixed

- Improve transaction manager architecture to be more compatible with `ETH_SECONDARY_URL` option (i.e. concurrent transaction submission to multiple different eth nodes). This also comes with some minor performance improvements in the tx manager and more correct handling of some extremely rare edge cases.
- As a side-effect, we now no longer handle the case where an external wallet used the chainlink ethereum private key to send a transaction. This use-case was already explicitly unsupported, but we made a best-effort attempt to handle it. We now make no attempt at all to handle it and doing this WILL result in your node not sending the data that it expected to be sent for the nonces that were used by an external wallet.
- Operator UI now shows booleans correctly

### Changed

- ETH_MAX_GAS_PRICE_WEI now 1500Gwei by default

## [0.8.18] - 2020-10-01

### Fixed

- Prometheus gas_updater_set_gas_price metric now only shows last gas price instead of every block since restart

## [0.8.17] - 2020-09-28

### Added

- Add new env variable ETH_SECONDARY_URL. Default is unset. You may optionally set this to an http(s) ethereum RPC client URL. If set, transactions will also be broadcast to this secondary ethereum node. This allows transaction broadcasting to be more robust in the face of primary ethereum node bugs or failures.
- Remove configuration option ORACLE_CONTRACT_ADDRESS, it had no effect
- Add configuration option OPERATOR_CONTRACT_ADDRESS, it filters the contract addresses the node should listen to for Run Logs
- At startup, the chainlink node will create a new funding address. This will initially be used to pay for cancelling stuck transactions.

### Fixed

- Gas bumper no longer hits database constraint error if ETH_MAX_GAS_PRICE_WEI is reached (this was actually mostly harmless, but the errors were annoying)

### Changes

- ETH_MAX_GAS_PRICE_WEI now defaults to 1500 gwei

## [0.8.16] - 2020-09-18

### Added

- The chainlink node now will bump a limited configurable number of transactions at once. This is configured with the ETH_GAS_BUMP_TX_DEPTH variable which is 10 by default. Set to 0 to disable (the old behaviour).

### Fixed

- ETH_DISABLED flag works again

## [0.8.15] - 2020-09-14

### Added

- Chainlink header images to the following `README.md` files: root, core,
  evm-contracts, and evm-test-helpers.
- Database migrations: new log_consumptions records will contain the number of the associated block.
  This migration will allow future version of chainlink to automatically clean up unneeded log_consumption records.
  This migration should execute very fast.
- External Adapters for the Flux Monitor will now receive the Flux Monitor round state info as the meta payload.
- Reduce frequency of balance checking.

### Fixed

Previously when the node was overloaded with heads there was a minor possibility it could get backed up with a very large head queue, and become unstable. Now, we drop heads instead in this case and noisily emit an error. This means the node should more gracefully handle overload conditions, although this is still dangerous and node operators should deal with it immediately to avoid missing jobs.

A new environment variable is introduced to configure this, called `ETH_HEAD_TRACKER_MAX_BUFFER_SIZE`. It is recommended to leave this set to the default of "3".

A new prometheus metric is also introduced to track dropped heads, called `head_tracker_num_heads_dropped`. You may wish to set an alert on a rule such as `increase(chainlink_dropped_heads[5m]) > 0`.

## [0.8.14] - 2020-09-02

## Changed

- Fix for gas bumper
- Fix for broadcast-transactions function

## [0.8.13] - 2020-08-31

## Changed

- Fix for gas bumper
- Fix for broadcast-transactions function

## [0.8.13] - 2020-08-31

### Changed

- Performance improvements when using BulletproofTxManager.

## [0.8.12] - 2020-08-10

### Fixed

- Added a workaround for Infura users who are seeing "error getting balance: header not found".
  This behaviour is due to Infura announcing it has a block, but when we request our balance in this block, the eth node doesn't have the block in memory. The workaround is to add a configurable lag time on balance update requests. The default is set to 1 but this is configurable via a new environment variable `ETH_BALANCE_MONITOR_BLOCK_DELAY`.

## [0.8.11] - 2020-07-27

### Added

- Job specs now support pinning to multiple keys using the new `fromAddresses` field in the ethtx task spec.

### Changed

- Using `fromAddress` in ethtx task specs has been deprecated. Please use `fromAddresses` instead.

### Breaking changes

- Support for RunLogTopic0original and RunLogTopic20190123withFullfillmentParams logs has been dropped. This should not affect any users since these logs predate Chainlink's mainnet launch and have never been used on mainnet.

IMPORTANT: The selection mechanism for keys has changed. When an ethtx task spec is not pinned to a particular key by defining `fromAddress` or `fromAddresses`, the node will now cycle through all available keys in round robin fashion. This is a change from the previous behaviour where nodes would only pick the earliest created key.

This is done to allow increases in throughput when a node operator has multiple whitelisted addresses for their oracle.

If your node has multiple keys, you will need to take one of the three following actions:

1. Make sure all keys are valid for all job specs
2. Pin job specs to a valid subset of key(s) using `fromAddresses`
3. Delete the key(s) you don't want to use

If your node only has one key, no action is required.

## [0.8.10] - 2020-07-14

### Fixed

- Incorrect sequence on keys table in some edge cases

## [0.8.9] - 2020-07-13

### Added

- Added a check on sensitive file ownership that gives a warning if certain files are not owned by the user running chainlink
- Added mechanism to asynchronously communicate when a job spec has an ethereum interaction error (or any async error) with a UI screen
- Gas Bumper now bumps based on the current gas price instead of the gas price of the original transaction

### Fixed

- Support for multiple node addresses

## [0.8.8] - 2020-06-29

### Added

- `ethtx` tasks now support a new parameter, `minRequiredOutgoingConfirmations` which allows you to tune how many confirmations are required before moving on from an `ethtx` task on a per task basis (only works with BulletproofTxManager). If it is not supplied, the default of `MIN_OUTGOING_CONFIRMATIONS` is used (same as the old behaviour).

### Changed

- HeadTracker now automatically backfills missing heads up to `ETH_FINALITY_DEPTH`
- The strategy for gas bumping has been changed to produce a potentially higher gas cost in exchange for the transaction getting through faster.

### Breaking changes

- `admin withdraw` command has been removed. This was only ever useful to withdraw LINK if the Oracle contract was owned by the Chainlink node address. It is no longer recommended to have the Oracle owner be the chainlink node address.
- Fixed `txs create` to send the amount in Eth not in Wei (as per the documentation)

## [0.8.7] - 2020-06-15

### Added

This release contains a number of features aimed at improving the node's reliability when putting transactions on-chain.

- An experimental new transaction manager is introduced that delivers reliability improvements compared to the old one, especially when faced with difficult network conditions or spiking gas prices. It also reduces load on the database and makes fewer calls to the eth node compared to the old tx manager.
- Along with the new transaction manager is a local client command for manually controlling the node nonce - `setnextnonce`. This should never be necessary under normal operation and is included only for use in emergencies.
- New prometheus metrics for the head tracker:
  - `head_tracker_heads_in_queue` - The number of heads currently waiting to be executed. You can think of this as the 'load' on the head tracker. Should rarely or never be more than 0.
  - `head_tracker_callback_execution_duration` - How long it took to execute all callbacks. If the average of this exceeds the time between blocks, your node could lag behind and delay transactions.
- Nodes transmit their build info to Explorer for better debugging/tracking.

### Env var changes

- `ENABLE_BULLETPROOF_TX_MANAGER` - set this to true to enable the experimental new transaction manager
- `ETH_GAS_BUMP_PERCENT` default value has been increased from 10% to 20%
- `ETH_GAS_BUMP_THRESHOLD` default value has been decreased from 12 to 3
- `ETH_FINALITY_DEPTH` specifies how deep protection should be against re-orgs. The default is 50. It only applies if BulletproofTxManager is enabled. It is not recommended to change this setting.
- `EthHeadTrackerHistoryDepth` specifies how many heads the head tracker should keep in the database. The default is 100. It is not recommended to change this setting.
- Update README.md with links to mockery, jq, and gencodec as they are required to run `go generate ./...`

## [0.8.6] - 2020-06-08

### Added

- The node now logs the eth client RPC calls
- More reliable Ethereum block header tracking
- Limit the amount of an HTTP response body that the node will read
- Make Aggregator contract interface viewable
- More resilient handling of chain reorganizations

## [0.8.5] - 2020-06-01

### Added

- The chainlink node can now be configured to backfill logs from `n` blocks after a
  connection to the ethereum client is reset. This value is specified with an environment
  variable `BLOCK_BACKFILL_DEPTH`.
- The chainlink node now sets file permissions on sensitive files on startup (tls, .api, .env, .password and secret)
- AggregatorInterface now has description and version fields.

### Changed

- Solidity: Renamed the previous `AggregatorInterface.sol` to
  `HistoricAggregatorInterface.sol`. Users are encouraged to use the new methods
  introduced on the `AggregatorInterface`(`getRoundData` and `latestRoundData`),
  as they return metadata to indicate freshness of the data in a single
  cross-contract call.
- Solidity: Marked `HistoricAggregatorInterface` methods (`latestAnswer`,
  `latestRound`, `latestTimestamp`, `getAnswer`, `getTimestamp`) as deprecated
  on `FluxAggregator`, `WhitelistedAggregator`, `AggregatorProxy`,
  `WhitelistedAggregatorProxy`.
- Updated the solidity compiler version for v0.6 from 0.6.2 to 0.6.6.
- AccessControlledAggregatorProxy checks an external contract for users to be able to
  read functions.

### Fixed

- Fluxmonitor jobs now respect the `minPayment` field on job specs and won't poll if the contract
  does not have sufficient funding. This allows certain jobs to require a larger payment
  than `MINIMUM_CONTRACT_PAYMENT`.

## [0.8.4] - 2020-05-18

### Added

- Fluxmonitor initiators may now optionally include an `absoluteThreshold`
  parameter. To trigger a new on-chain report, the absolute difference in the feed
  value must change by at least the `absoluteThreshold` value. If it is
  unspecified or zero, fluxmonitor behavior is unchanged.
- Database Migrations: Add created_at and updated_at to all tables allowing for
  better historical insights. This migration may take a minute or two on large
  databases.

### Fixed

- Fix incorrect permissions on some files written by the node
  Prevent a case where duplicate ethereum keys could be added
  Improve robustness and reliability of ethtx transaction logic

## [0.8.3] - 2020-05-04

### Added

- Added Changelog.
- Database Migrations: There a number of database migrations included in this
  release as part of our ongoing effort to make the node even more reliable and
  stable, and build a firm foundation for future development.

### Changed

- New cron strings MUST now include time zone. If you want your jobs to run in
  UTC for example: `CRON_TZ=UTC * * * * *`. Previously, jobs specified without a
  time zone would run in the server's native time zone, which in most cases is UTC
  but this was never guaranteed.

### Fixed

- Fix crash in experimental gas updater when run on Kovan network

## [0.8.2] - 2020-04-20

## [0.8.1] - 2020-04-08

## [0.8.0] - 2020-04-06
