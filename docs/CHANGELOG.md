# Changelog Chainlink Core

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- OCR bootstrap node now sends telemetry to the endpoint specified in the OCR job spec under `MonitoringEndpoint`.
- Adds "Account addresses" table to the `/keys` page.

### Changed

- Old jobs now allow duplicate job names. Also, if the name field is empty we no longer generate a name.

### Fixed

- Brings `/runs` tab back to the operator UI.
- Signs out a user from operator UI on authentication error.

### Changes

- Removes broken `ACCOUNT_ADDRESS` field from `/config` page.

### Changed

- Commands for creating/managing legacy jobs and OCR jobs have changed, to accomodate additional types of jobs using the new pipeline.

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
