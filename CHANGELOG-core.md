# Changelog Chainlink Core

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

This release contains a number of features aimed at improving the node's reliability when putting transactions on-chain.

- An experimental new transaction manager is introduced that delivers reliability improvements compared to the old one, especially when faced with difficult network conditions or spiking gas prices. It also reduces load on the database and makes fewer calls to the eth node compared to the old tx manager.
- Along with the new transaction manager is a local client command for manually controlling the node nonce - `setnextnonce`. This should never be necessary under normal operation and is included only for use in emergencies.
- New prometheus metrics for the head tracker:
  - `head_tracker_heads_in_queue` - The number of heads currently waiting to be executed. You can think of this as the 'load' on the head tracker. Should rarely or never be more than 0.
  - `head_tracker_callback_execution_duration` - How long it took to execute all callbacks. If the average of this exceeds the time between blocks, your node could lag behind and delay transactions.

### Env var changes

- `ENABLE_BULLETPROOF_TX_MANAGER` - set this to true to enable the experimental new transaction manager. CAUTION: Once the node is booted with this toggled on, a configuration value is written that makes it permanent. Toggling it off again will have no effect. This is done because BulletproofTxManager requires tight control over the nonce sequence and if you revert to the old tx manager, the nonce can get out of sync. If you _must_ revert to the old transaction manager, at your own risk, you can execute the following SQL command: `DELETE FROM configurations WHERE name = 'ENABLE_BULLETPROOF_TX_MANAGER';`.
- `ETH_GAS_BUMP_PERCENT` default value has been increased from 10% to 20%
- `ETH_GAS_BUMP_THRESHOLD` default value has been decreased from 12 to 3
- `ETH_FINALITY_DEPTH` specifies how deep protection should be against re-orgs. The default is 50. It only applies if BulletproofTxManager is enabled. It is not recommended to change this setting.
- `EthHeadTrackerHistoryDepth` specifies how many heads the head tracker should keep in the database. The default is 100. It is not recommended to change this setting.

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
