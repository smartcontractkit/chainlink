# Changelog Chainlink Core

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<!-- unreleased -->
## [dev]

...

<!-- unreleasedstop -->

## 1.13.0 - 2023-03-16

### Added

- Support for sending OCR2 job specs to the feeds manager
- Log poller filters now saved in db, restored on node startup to guard against missing logs during periods where services are temporarily unable to start

### Updated

- TOML env var `CL_CONFIG` always processed as the last configuration, with the effect of being the final override 
of any values provided via configuration files.

### Changed

- The config option `FeatureFeedsManager`/`FEATURE_FEEDS_MANAGER` is now true by default.

### Removed

- Terra is no longer supported

## 1.12.0 - 2023-02-15

### Added

- Prometheus gauge `mailbox_load_percent` for percent of "`Mailbox`" capacity used.
- New config option, `JobPipeline.MaxSuccessfulRuns` caps the total number of
  saved completed runs per job. This is done in response to the `pipeline_runs`
  table potentially becoming large, which can cause performance degradation.
  The default is set to 10,000. You can set it to 0 to disable run saving
  entirely. **NOTE**: This can only be configured via TOML and not with an
  environment variable.
- Prometheus gauge vector `feeds_job_proposal_count` to track counts of job proposals partitioned by proposal status.
- Support for variable expression for the `minConfirmations` parameter on the `ethtx` task.

### Updated

- Removed `KEEPER_TURN_FLAG_ENABLED` as all networks/nodes have switched this to `true` now. The variable should be completely removed my NOPs.
- Removed `Keeper.UpkeepCheckGasPriceEnabled` config (`KEEPER_CHECK_UPKEEP_GAS_PRICE_FEATURE_ENABLED` in old env var configuration) as this feature is deprecated now. The variable should be completely removed by NOPs.

### Fixed

- Fixed (SQLSTATE 42P18) error on Job Runs page, when attempting to view specific older or infrequenty run jobs
- The `config dump` subcommand was fixed to dump the correct config data. 
  - The `P2P.V1.Enabled` config logic incorrectly matched V2, by only setting explicit true values so that otherwise the default is used. The `V1.Enabled` default value is actually true already, and is now updated to only set explicit false values.
  - The `[EVM.Transactions]` config fields `MaxQueued` & `MaxInFlight` will now correctly match `ETH_MAX_QUEUED_TRANSACTIONS` & `ETH_MAX_IN_FLIGHT_TRANSACTIONS`.

## 1.11.0 - 2022-12-12

### Added

- New `EVM.NodePool.SelectionMode` `TotalDifficulty` to use the node with the greatest total difficulty.
- Add the following prometheus metrics (labelled by bridge name) for monitoring external adapter queries:
    - `bridge_latency_seconds`
    - `bridge_errors_total`
    - `bridge_cache_hits_total`
    - `bridge_cache_errors_total`
- `EVM.NodePool.SyncThreshold` to ensure that live nodes do not lag too far behind.

> ```toml
> SyncThreshold = 5 # Default
> ```
> 
> SyncThreshold controls how far a node may lag behind the best node before being marked out-of-sync.
Depending on `SelectionMode`, this represents a difference in the number of blocks (`HighestHead`, `RoundRobin`), or total difficulty (`TotalDifficulty`).
>
> Set to 0 to disable this check.

#### TOML Configuration (experimental)

Chainlink now supports static configuration via TOML files as an alternative to the existing combination of environment variables and persisted database configurations.

This is currently _experimental_, but in the future (with `v2.0.0`), it will become *mandatory* as the only supported configuration method. Avoid using TOML for configuration unless running on a test network for this release.

##### How to use

TOML configuration can be enabled by simply using the new `-config <filename>` flag or `CL_CONFIG` environment variable.
Multiple files can be used (`-c configA.toml -c configB.toml`), and will be applied in order with duplicated fields overriding any earlier values.

Existing nodes can automatically generate their equivalent TOML configuration via the `config dump` subcommand.
Secrets must be configured manually and passed via `-secrets <filename>` or equivalent environment variables.

Format details: [CONFIG.md](../docs/CONFIG.md) • [SECRETS.md](../docs/SECRETS.md)

**Note:** You _cannot_ mix legacy environment variables with TOML configuration. Leaving any legacy env vars set will fail validation and prevent boot.

##### Examples

Dump your current configuration as TOML.
```bash
chainlink config dump > config.toml
```

Inspect your full effective configuration, and ensure it is valid. This includes defaults.
```bash
chainlink --config config.toml --secrets secrets.toml config validate
```

Run the node.
```bash
chainlink -c config.toml -s secrets.toml node start
```

#### Bridge caching
##### BridgeCacheTTL

- Default: 0s

When set to `d` units of time, this variable enables using cached bridge responses that are at most `d` units old. Caching is disabled by default.

Example `BridgeCacheTTL=10s`, `BridgeCacheTTL=1m`

### Fixed

- Fixed a minor bug whereby Chainlink would not always resend all pending transactions when using multiple keys

### Updated

- `NODE_NO_NEW_HEADS_THRESHOLD=0` no longer requires `NODE_SELECTION_MODE=RoundRobin`. 

## 1.10.0 - 2022-11-15

### Added

#### New optional external logger added
##### AUDIT_LOGGER_FORWARD_TO_URL

- Default: _none_

When set, this environment variable configures and enables an optional HTTP logger which is used specifically to send audit log events. Audit logs events are emitted when specific actions are performed by any of the users through the node's API. The value of this variable should be a full URL. Log items will be sent via POST

There are audit log implemented for the following events:
  - Auth & Sessions (new session, login success, login failed, 2FA enrolled, 2FA failed, password reset, password reset failed, etc.)
  - CRUD actions for all resources (add/create/delete resources such as bridges, nodes, keys)
  - Sensitive actions (keys exported/imported, config changed, log level changed, environment dumped)

A full list of audit log enum types can be found in the source within the `audit` package (`audit_types.go`).

The following `AUDIT_LOGGER_*` environment variables below configure this optional audit log HTTP forwarder.

##### AUDIT_LOGGER_HEADERS

- Default: _none_

An optional list of HTTP headers to be added for every optional audit log event. If the above `AUDIT_LOGGER_FORWARD_TO_URL` is set, audit log events will be POSTed to that URL, and will include headers specified in this environment variable. One example use case is auth for example: ```AUDIT_LOGGER_HEADERS="Authorization||{{token}}"```.

Header keys and values are delimited on ||, and multiple headers can be added with a forward slash delimiter ('\\'). An example of multiple key value pairs:
```AUDIT_LOGGER_HEADERS="Authorization||{{token}}\Some-Other-Header||{{token2}}"```

##### AUDIT_LOGGER_JSON_WRAPPER_KEY

- Default: _none_

When the audit log HTTP forwarder is enabled, if there is a value set for this optional environment variable then the POST body will be wrapped in a dictionary in a field specified by the value of set variable. This is to help enable specific logging service integrations that may require the event JSON in a special shape. For example: `AUDIT_LOGGER_JSON_WRAPPER_KEY=event` will create the POST body:
```
{
  "event": {
    "eventID":  EVENT_ID_ENUM,
    "data": ...
  }
}
```

#### Automatic connectivity detection; Chainlink will no longer bump excessively if the network is broken

This feature only applies on EVM chains when using BlockHistoryEstimator (the most common case).

Chainlink will now try to automatically detect if there is a transaction propagation/connectivity issue and prevent bumping in these cases. This can help avoid the situation where RPC nodes are not propagating transactions for some reason (e.g. go-ethereum bug, networking issue etc) and Chainlink responds in a suboptimal way by bumping transactions to a very high price in an effort to get them mined. This can lead to unnecessary expense when the connectivity issue is resolved and the transactions are finally propagated into the mempool.

This feature is enabled by default with fairly conservative settings: if a transaction has been priced above the 90th percentile of the past 12 blocks, but still wants to bump due to not being mined, a connectivity/propagation issue is assumed and all further bumping will be prevented for this transaction. In this situation, Chainlink will start firing the `block_history_estimator_connectivity_failure_count` prometheus counter and logging at critical level until the transaction is mined.

The default settings should work fine for most users. For advanced users, the values can be tweaked by changing `BLOCK_HISTORY_ESTIMATOR_CHECK_INCLUSION_BLOCKS` and `BLOCK_HISTORY_ESTIMATOR_CHECK_INCLUSION_PERCENTILE`.

To disable connectivity checking completely, set `BLOCK_HISTORY_ESTIMATOR_CHECK_INCLUSION_BLOCKS=0`.

### Changed

- The default maximum gas price on most networks is now effectively unlimited.
  - Chainlink will bump as high as necessary to get a transaction included. The connectivity checker is relied on to prevent excessive bumping when there is a connectivity failure.
  - If you want to change this, you can manually set `ETH_MAX_GAS_PRICE_WEI`.

- EVMChainID field will be auto-added with default chain id to job specs of newly created OCR jobs, if not explicitly included.
  - Old OCR jobs missing EVMChainID will continue to run on any chain ETH_CHAIN_ID is set to (or first chain if unset), which may be changed after a restart.
  - Newly created OCR jobs will only run on a single fixed chain, unaffected by changes to ETH_CHAIN_ID after the job is added.
  - It should no longer be possible to end up with multiple OCR jobs for a single contract running on the same chain; only one job per contract per chain is allowed
  - If there are any existing duplicate jobs (per contract per chain), all but the job with the latest creation date will be pruned during upgrade.

### Fixed

- Fixed minor bug where Chainlink would attempt (and fail) to estimate a tip cap higher than the maximum configured gas price in EIP1559 mode. It now caps the tipcap to the max instead of erroring.
- Fixed bug whereby it was impossible to remove eth keys that had extant transactions. Now, removing an eth key will drop all associated data automatically including past transactions.

## 1.9.0 - 2022-10-12

### Added

- Added `length` and `lessthan` tasks (pipeline).
- Added `gasUnlimited` parameter to `ethcall` task. 
- `/keys` page in Operator UI now exposes several admin commands, namely:
  - "abandon" to abandon all current txes
  - enable/disable a key for a given chain
  - manually set the nonce for a key
  See [this PR](https://github.com/smartcontractkit/chainlink/pull/7406) for a screenshot example.

## 1.8.1 - 2022-09-29

### Added

-  New `GAS_ESTIMATOR_MODE` for Arbitrum to support Nitro's multi-dimensional gas model, with dynamic gas pricing and limits.
   -  NOTE: It is recommended to remove `GAS_ESTIMATOR_MODE` as an env var if you have it set in order to use the new default.
   -  This new, default estimator for Arbitrum networks uses the suggested gas price (up to `ETH_MAX_GAS_PRICE_WEI`, with `1000 gwei` default) as well as an estimated gas limit (up to `ETH_GAS_LIMIT_MAX`, with `1,000,000,000` default).
- `ETH_GAS_LIMIT_MAX` to put a maximum on the gas limit returned by the `Arbitrum` estimator.

### Changed

- EIP1559 is now enabled by default on Goerli network

## 1.8.0 - 2022-09-01

### Added

- Added `hexencode` and `base64encode` tasks (pipeline).
- `forwardingAllowed` per job attribute to allow forwarding txs submitted by the job.
- Keypath now supports paths with any depth, instead of limiting it to 2
- `Arbitrum` chains are no longer restricted to only `FixedPrice` `GAS_ESTIMATOR_MODE`
- Updated `Arbitrum Rinkeby & Mainnet & Mainnet` configurationss for Nitro
- Add `Arbitrum Goerli` configuration
- It is now possible to use the same key across multiple chains.
- `NODE_SELECTION_MODE` (`EVM.NodePool.SelectionMode`) controls node picking strategy. Supported values: `HighestHead` (default) and `RoundRobin`:
  - `RoundRobin` mode simply iterates among available alive nodes. This was the default behavior prior to this release.
  - `HighestHead` mode picks a node having the highest reported head number among other alive nodes. When several nodes have the same latest head number, the strategy sticks to the last used node.
  For chains having `NODE_NO_NEW_HEADS_THRESHOLD=0` (such as Arbitrum, Optimism), the implementation will fall back to `RoundRobin` mode.
- New `keys eth chain` command
  - This can also be accessed at `/v2/keys/evm/chain`.
  - Usage examples:
    - Manually (re)set a nonce:
      - `chainlink keys eth chain --address "0xEXAMPLE" --evmChainID 99 --setNextNonce 42`
    - Enable a key for a particular chain:
      - `chainlink keys eth chain --address "0xEXAMPLE" --evmChainID 99 --enable`
    - Disable a key for a particular chain:
      - `chainlink keys eth chain --address "0xEXAMPLE" --evmChainID 99 --disable`
    - Abandon all currently pending transactions (use with caution!):
      - `chainlink evm keys chain --address "0xEXAMPLE" --evmChainID 99 --abandon`
  - Commands can be combined e.g.
    - Reset nonce and abandon all currently pending transaction:
      - `chainlink evm keys chain --address "0xEXAMPLE" --evmChainID 99 --setNextNonce 42 --abandon`

### Changed

- The `setnextnonce` local client command has been removed, and replaced by a more general key/chain client command.
- `chainlink admin users update` command is replaced with `chainlink admin users chrole` (only the role can be changed for a user)

## 1.7.1 - 2022-08-22

### Added

- `Arbitrum Nitro` client error support

## 1.7.0 - 2022-08-08

### Added

- `p2pv2Bootstrappers` has been added as a new optional property of OCR1 job specs; default may still be specified with P2PV2_BOOTSTRAPPERS config param
- Added official support for Sepolia chain
- Added `hexdecode` and `base64decode` tasks (pipeline).
- Added support for Besu execution client (note that while Chainlink supports Besu, Besu itself [has](https://github.com/hyperledger/besu/issues/4212) [multiple](https://github.com/hyperledger/besu/issues/4192) [bugs](https://github.com/hyperledger/besu/issues/4114) that make it unreliable).
- Added the functionality to allow the root admin CLI user (and any additional admin users created) to create and assign tiers of role based access to new users. These new API users will be able to log in to the Operator UI independently, and can each have specific roles tied to their account. There are four roles: `admin`, `edit`, `run`, and `view`.
  - User management can be configured through the use of the new admin CLI command `chainlink admin users`. Be sure to run `chainlink adamin login`. For example, a readonly user can be created with: `chainlink admin users create --email=operator-ui-read-only@test.com --role=view`.
  - Updated documentation repo with a break down of actions to required role level
- Added per job spec and per job type gas limit control. The following rule of precedence is applied:

1. task-specific parameter `gasLimit` overrides anything else when specified (e.g. `ethtx` task has such a parameter).
2. job-spec attribute `gasLimit` has the scope of the current job spec only.
3. job-type limits `ETH_GAS_LIMIT_*_JOB_TYPE` affect any jobs of the corresponding type:

```
ETH_GAS_LIMIT_OCR_JOB_TYPE    # EVM.GasEstimator.LimitOCRJobType
ETH_GAS_LIMIT_DR_JOB_TYPE     # EVM.GasEstimator.LimitDRJobType
ETH_GAS_LIMIT_VRF_JOB_TYPE    # EVM.GasEstimator.LimitVRFJobType
ETH_GAS_LIMIT_FM_JOB_TYPE     # EVM.GasEstimator.LimitFMJobType
ETH_GAS_LIMIT_KEEPER_JOB_TYPE # EVM.GasEstimator.LimitKeeperJobType
```

4. global `ETH_GAS_LIMIT_DEFAULT` (`EVM.GasEstimator.LimitDefault`) value is the last resort.

### Fixed

- Addressed a very rare bug where using multiple nodes with differently configured RPC tx fee caps could cause missed transaction. Reminder to everyone to ensure that your RPC nodes have no caps (for more information see the [performance and tuning guide](https://docs.chain.link/docs/evm-performance-configuration/)).
- Improved handling of unknown transaction error types, making Chainlink more robust in certain cases on unsupported chains/RPC clients

## [1.6.0] - 2022-07-20

### Changed

- After feedback from users, password complexity requirements have been simplified. These are the new, simplified requirements for any kind of password used with Chainlink:
1. Must be 16 characters or more
2. Must not contain leading or trailing whitespace
3. User passwords must not contain the user's API email

- Simplified the Keepers job spec by removing the observation source from the required parameters.

## [1.5.1] - 2022-06-27

### Fixed

- Fix rare out-of-sync to invalid-chain-id transaction
- Fix key-specific max gas limits for gas estimator and ensure we do not bump gas beyond key-specific limits
- Fix EVM_FINALITY_DEPTH => ETH_FINALITY_DEPTH

## [1.5.0] - 2022-06-21

### Changed

- Chainlink will now log a warning if the postgres database password is missing or too insecure. Passwords should conform to the following rules:
```
Must be longer than 12 characters
Must comprise at least 3 of:
	lowercase characters
	uppercase characters
	numbers
	symbols
Must not comprise:
	More than three identical consecutive characters
	Leading or trailing whitespace (note that a trailing newline in the password file, if present, will be ignored)
```
For backward compatibility all insecure passwords will continue to work, however in a future version of Chainlink insecure passwords will prevent application boot. To bypass this check at your own risk, you may set `SKIP_DATABASE_PASSWORD_COMPLEXITY_CHECK=true`.

- `MIN_OUTGOING_CONFIRMATIONS` has been removed and no longer has any effect. `ETH_FINALITY_DEPTH` is now used as the default for `ethtx` confirmations instead. You may override this on a per-task basis by setting `minConfirmations` in the task definition e.g. `foo [type=ethtx minConfirmations=42 ...]`. NOTE: This may have a minor impact on performance on very high throughput chains. If you don't care about reporting task status in the UI, it is recommended to set `minConfirmations=0` in your job specs. For more details, see the [relevant section of the performance tuning guide](https://www.notion.so/chainlink/EVM-performance-configuration-handbook-a36b9f84dcac4569ba68772aa0c1368c#e9998c2f722540b597301a640f53cfd4).

- The following ENV variables have been deprecated, and will be removed in a future release: `INSECURE_SKIP_VERIFY`, `CLIENT_NODE_URL`, `ADMIN_CREDENTIALS_FILE`. These vars only applied to Chainlink when running in client mode and have been replaced by command line args, notably: `--insecure-skip-verify`, `--remote-node-url URL` and `--admin-credentials-file FILE` respectively. More information can be found by running `./chainlink --help`.

- The `Optimism2` `GAS_ESTIMATOR_MODE` has been renamed to `L2Suggested`. The old name is still supported for now.

- The `p2pBootstrapPeers` property on OCR2 job specs has been renamed to `p2pv2Bootstrappers`.

### Added
- Added `ETH_USE_FORWARDERS` config option to enable transactions forwarding contracts.
- In job pipeline (direct request) the three new block variables are exposed:
  - `$(jobRun.blockReceiptsRoot)` : the root of the receipts trie of the block (hash)
  - `$(jobRun.blockTransactionsRoot)` : the root of the transaction trie of the block (hash)
  - `$(jobRun.blockStateRoot)` : the root of the final state trie of the block (hash)
- `ethtx` tasks can now be configured to error if the transaction reverts on-chain. You must set `failOnRevert=true` on the task to enable this behavior, like so:

`foo [type=ethtx failOnRevert=true ...]`

So the `ethtx` task now works as follows:

If minConfirmations == 0, task always succeeds and nil is passed as output
If minConfirmations > 0, the receipt is passed through as output
If minConfirmations > 0 and failOnRevert=true then the ethtx task will error on revert

If `minConfirmations` is not set on the task, the chain default will be used which is usually 12 and always greater than 0.

- `http` task now allows specification of request headers. Use like so: `foo [type=http headers="[\\"X-Header-1\\", \\"value1\\", \\"X-Header-2\\", \\"value2\\"]"]`.


### Fixed
- Fixed `max_unconfirmed_age` metric. Previously this would incorrectly report the max time since the last rebroadcast, capping the upper limit to the EthResender interval. This now reports the correct value of total time elapsed since the _first_ broadcast.
- Correctly handle the case where bumped gas would exceed the RPC node's configured maximum on Fantom (note that node operators should check their Fantom RPC node configuration and remove the fee cap if there is one)
- Fixed handling of Metis internal fee change

### Removed

- The `Optimism` OVM 1.0 `GAS_ESTIMATOR_MODE` has been removed.

## [1.4.1] - 2022-05-11

### Fixed

- Ensure failed EthSubscribe didn't register a (*rpc.ClientSubscription)(nil) which would lead to a panic on Unsubscribe
- Fixes parsing of float values on job specs

## [1.4.0] - 2022-05-02

### Added

- JSON parse tasks (v2) now support a custom `separator` parameter to substitute for the default `,`.
- Log slow SQL queries
- Fantom and avalanche block explorer urls
- Display `requestTimeout` in job UI
- Keeper upkeep order is shuffled

### Fixed

- `LOG_FILE_MAX_SIZE` handling
- Improved websocket subscription management (fixes issues with multiple-primary-node failover from 1.3.x)
- VRFv2 fixes and enhancements
- UI support for `minContractPaymentLinkJuels`

## [1.3.0] - 2022-04-18

### Added

- Added support for Keeper registry v1.2 in keeper jobs
- Added disk rotating logs. Chainlink will now always log to disk at debug level. The default output directory for debug logs is Chainlink's root directory (ROOT_DIR) but can be configured by setting LOG_FILE_DIR. This makes it easier for node operators to report useful debugging information to Chainlink's team, since all the debug logs are conveniently located in one directory. Regular logging to STDOUT still works as before and respects the LOG_LEVEL env var. If you want to log in disk at a particular level, you can pipe STDOUT to disk. This automatic debug-logs-to-disk feature is enabled by default, and will remain enabled as long as the `LOG_FILE_MAX_SIZE` ENV var is set to a value greater than zero. The amount of disk space required for this feature to work can be calculated with the following formula: `LOG_FILE_MAX_SIZE` * (`LOG_FILE_MAX_BACKUPS` + 1). If your disk doesn't have enough disk space, the logging will pause and the application will log Errors until space is available again. New environment variables related to this feature:
  - `LOG_FILE_MAX_SIZE` (default: 5120mb) - this env var allows you to override the log file's max size (in megabytes) before file rotation.
  - `LOG_FILE_MAX_AGE` (default: 0) - if `LOG_FILE_MAX_SIZE` is set, this env var allows you to override the log file's max age (in days) before file rotation. Keeping this config with the default value means not to remove old log files.
  - `LOG_FILE_MAX_BACKUPS` (default: 1) - if `LOG_FILE_MAX_SIZE` is set, this env var allows you to override the max amount of old log files to retain. Keeping this config with the default value means to retain 1 old log file at most (though `LOG_FILE_MAX_AGE` may still cause them to get deleted). If this is set to 0, the node will retain all old log files instead.
- Added support for the `force` flag on `chainlink blocks replay`. If set to true, already consumed logs that would otherwise be skipped will be rebroadcasted.
- Added version compatibility check when using CLI to login to a remote node. flag `bypass-version-check` skips this check.
- Interrim solution to set multiple nodes/chains from ENV. This gives the ability to specify multiple RPCs that the Chainlink node will constantly monitor for health and sync status, detecting dead nodes and out of sync nodes, with automatic failover. This is a temporary stand-in until configuration is overhauled and will be removed in future in favor of a config file. Set as such: `EVM_NODES='{...}'` where the var is a JSON array containing the node specifications. This is not compatible with using any other way to specify node via env (e.g. `ETH_URL`, `ETH_SECONDARY_URL`, `ETH_CHAIN_ID` etc). **WARNING**: Setting this environment variable will COMPLETELY ERASE your `evm_nodes` table on every boot and repopulate from the given data, nullifying any runtime modifications. Make sure to carefully read the [EVM performance configuration guide](https://chainlink.notion.site/EVM-performance-configuration-handbook-a36b9f84dcac4569ba68772aa0c1368c) for best practices here.

For example:

```bash
export EVM_NODES='
[
	{
		"name": "primary_1",
		"evmChainId": "137",
		"wsUrl": "wss://endpoint-1.example.com/ws",
    "httpUrl": "http://endpoint-1.example.com/",
		"sendOnly": false
	},
	{
		"name": "primary_2",
		"evmChainId": "137",
		"wsUrl": "ws://endpoint-2.example.com/ws",
    "httpUrl": "http://endpoint-2.example.com/",
		"sendOnly": false
	},
	{
		"name": "primary_3",
		"evmChainId": "137",
		"wsUrl": "wss://endpoint-3.example.com/ws",
    "httpUrl": "http://endpoint-3.example.com/",
		"sendOnly": false
	},
	{
		"name": "sendonly_1",
		"evmChainId": "137",
		"httpUrl": "http://endpoint-4.example.com/",
		"sendOnly": true
	},
  {
		"name": "sendonly_2",
		"evmChainId": "137",
		"httpUrl": "http://endpoint-5.example.com/",
		"sendOnly": true
	}
]
'
```

### Changed

- Changed default locking mode to "dual". Bugs in lease locking have been ironed out and this paves the way to making "lease" the default in the future. It is recommended to set `DATABASE_LOCKING_MODE=lease`, default is set to "dual" only for backwards compatibility.
- EIP-1559 is now enabled by default on mainnet. To disable (go back to legacy mode) set `EVM_EIP1559_DYNAMIC_FEES=false`. The default settings should work well, but if you wish to tune your gas controls, see the [documentation](https://docs.chain.link/docs/configuration-variables/#evm-gas-controls).

Note that EIP-1559 can be manually enabled on other chains by setting `EVM_EIP1559_DYNAMIC_FEES=true` but we only support it for official Ethereum mainnet and testnets. It is _not_ recommended enabling this setting on Polygon since during our testing process we found that the EIP-1559 fee market appears to be broken on all Polygon chains and EIP-1559 transactions are actually less likely to get included than legacy transactions.

See issue: https://github.com/maticnetwork/bor/issues/347

- The pipeline task runs have changed persistence protocol (database), which will result in inability to decode some existing task runs. All new runs should be working with no issues.

### Removed

- `LOG_TO_DISK` ENV var.

## [1.2.1] - 2022-03-17

This release hotfixes issues from moving a new CI/CD system. Feature-wise the functionality is the same as `v1.2.0`.

### Fixed

- Fixed CI/CD issue where environment variables were not being passed into the underlying build

## [1.2.0] - 2022-03-02

### Added

- Added support for the Nethermind Ethereum client.
- Added support for batch sending telemetry to the ingress server to improve performance.
- Added v2 P2P networking support (alpha)

New ENV vars:

- `ADVISORY_LOCK_CHECK_INTERVAL` (default: 1s) - when advisory locking mode is enabled, this controls how often Chainlink checks to make sure it still holds the advisory lock. It is recommended to leave this at the default.
- `ADVISORY_LOCK_ID` (default: 1027321974924625846) - when advisory locking mode is enabled, the application advisory lock ID can be changed using this env var. All instances of Chainlink that might run on a particular database must share the same advisory lock ID. It is recommended to leave this at the default.
- `LOG_FILE_DIR` (default: chainlink root directory) - if `LOG_FILE_MAX_SIZE` is set, this env var allows you to override the output directory for logging.
- `SHUTDOWN_GRACE_PERIOD` (default: 5s) - when node is shutting down gracefully and exceeded this grace period, it terminates immediately (trying to close DB connection) to avoid being SIGKILLed.
- `SOLANA_ENABLED` (default: false) - set to true to enable Solana support
- `TERRA_ENABLED` (default: false) - set to true to enable Terra support
- `BLOCK_HISTORY_ESTIMATOR_EIP1559_FEE_CAP_BUFFER_BLOCKS` - if EIP1559 mode is enabled, this optional env var controls the buffer blocks to add to the current base fee when sending a transaction. By default, the gas bumping threshold + 1 block is used. It is not recommended to change this unless you know what you are doing.
- `TELEMETRY_INGRESS_BUFFER_SIZE` (default: 100) - the number of telemetry messages to buffer before dropping new ones
- `TELEMETRY_INGRESS_MAX_BATCH_SIZE` (default: 50) - the maximum number of messages to batch into one telemetry request
- `TELEMETRY_INGRESS_SEND_INTERVAL` (default: 500ms) - the cadence on which batched telemetry is sent to the ingress server
- `TELEMETRY_INGRESS_SEND_TIMEOUT` (default: 10s) - the max duration to wait for the request to complete when sending batch telemetry
- `TELEMETRY_INGRESS_USE_BATCH_SEND` (default: true) - toggles sending telemetry using the batch client to the ingress server
- `NODE_NO_NEW_HEADS_THRESHOLD` (default: 3m) - RPC node will be marked out-of-sync if it does not receive a new block for this length of time. Set to 0 to disable head monitoring for liveness checking,
- `NODE_POLL_FAILURE_THRESHOLD` (default: 5) - number of consecutive failed polls before an RPC node is marked dead. Set to 0 to disable poll liveness checking.
- `NODE_POLL_INTERVAL` (default: 10s) - how often to poll. Set to 0 to disable all polling.

#### Bootstrap job

Added a new `bootstrap` job type. This job removes the need for every job to implement their own bootstrapping logic.
OCR2 jobs with `isBootstrapPeer=true` are automatically migrated to the new format.
The spec parameters are similar to a basic OCR2 job, an example would be:

```
type            = "bootstrap"
name            = "bootstrap"
relay           = "evm"
schemaVersion	= 1
contractID      = "0xAb5801a7D398351b8bE11C439e05C5B3259aeC9B"
[relayConfig]
chainID	        = 4
```

#### EVM node hot failover and liveness checking

Chainlink now supports hot failover and liveness checking for EVM nodes. This completely supercedes and replaces the Fiews failover proxy and should remove the need for any kind of failover proxy between Chainlink and its RPC nodes.

In order to use this feature, you'll need to set multiple primary RPC nodes.

### Removed

- `deleteuser` CLI command.

### Changed

`EVM_DISABLED` has been deprecated and replaced by `EVM_ENABLED` for consistency with other feature flags.
`ETH_DISABLED` has been deprecated and replaced by `EVM_RPC_ENABLED` for consistency, and because this was confusingly named. In most cases you want to set `EVM_ENABLED=false` and not `EVM_RPC_ENABLED=false`.

Log colorization is now disabled by default because it causes issues when piped to text files. To re-enable log colorization, set `LOG_COLOR=true`.

#### Polygon/matic defaults changed

Due to increasingly hostile network conditions on Polygon we have had to increase a number of default limits. This is to work around numerous and very deep re-orgs, high mempool pressure and a failure by the network to propagate transactions properly. These new limits are likely to increase load on both your Chainlink node and database, so please be sure to monitor CPU and memory usage on both and make sure they are adequately specced to handle the additional load.

## [1.1.1] - 2022-02-14

### Added

- `BLOCK_HISTORY_ESTIMATOR_EIP1559_FEE_CAP_BUFFER_BLOCKS` - if EIP1559 mode is enabled, this optional env var controls the buffer blocks to add to the current base fee when sending a transaction. By default, the gas bumping threshold + 1 block is used. It is not recommended to change this unless you know what you are doing.
- `EVM_GAS_FEE_CAP_DEFAULT` - if EIP1559 mode is enabled, and FixedPrice gas estimator is used, this env var controls the fixed initial fee cap.
- Allow dumping pprof even when not in dev mode, useful for debugging (go to /v2/debug/pprof as a logged in user)

### Fixed

- Update timeout so we don’t exit early on very large log broadcaster backfills

#### EIP-1559 Fixes

Fixed issues with EIP-1559 related to gas bumping. Due to [go-ethereum's implementation](https://github.com/ethereum/go-ethereum/blob/bff330335b94af3643ac2fb809793f77de3069d4/core/tx_list.go#L298) which introduces additional restrictions on top of the EIP-1559 spec, we must bump the FeeCap at least 10% each time in order for the gas bump to be accepted.

The new EIP-1559 implementation works as follows:

If you are using FixedPriceEstimator:
- With gas bumping disabled, it will submit all transactions with `feecap=ETH_MAX_GAS_PRICE_WEI` and `tipcap=EVM_GAS_TIP_CAP_DEFAULT`
- With gas bumping enabled, it will submit all transactions initially with `feecap=EVM_GAS_FEE_CAP_DEFAULT` and `tipcap=EVM_GAS_TIP_CAP_DEFAULT`.

If you are using BlockHistoryEstimator (default for most chains):
- With gas bumping disabled, it will submit all transactions with `feecap=ETH_MAX_GAS_PRICE_WEI` and `tipcap=<calculated using past blocks>`
- With gas bumping enabled (default for most chains) it will submit all transactions initially with `feecap = ( current block base fee * (1.125 ^ N) + tipcap )` where N is configurable by setting BLOCK_HISTORY_ESTIMATOR_EIP1559_FEE_CAP_BUFFER_BLOCKS but defaults to `gas bump threshold+1` and `tipcap=<calculated using past blocks>`

Bumping works as follows:

- Increase tipcap by `max(tipcap * (1 + ETH_GAS_BUMP_PERCENT), tipcap + ETH_GAS_BUMP_WEI)`
- Increase feecap by `max(feecap * (1 + ETH_GAS_BUMP_PERCENT), feecap + ETH_GAS_BUMP_WEI)`

## [1.1.0] - 2022-01-25

### Added

- Added support for Sentry error reporting. Set `SENTRY_DSN` at run-time to enable reporting.
- Added Prometheus counters: `log_warn_count`, `log_error_count`, `log_critical_count`, `log_panic_count` and `log_fatal_count` representing the corresponding number of warning/error/critical/panic/fatal messages in the log.
- The new prometheus metric `tx_manager_tx_attempt_count` is a Prometheus Gauge that should represent the total number of Transactions attempts that awaiting confirmation for this node.
- The new prometheus metric `version` that displays the node software version (tag) as well as the corresponding commit hash.
- CLI command `keys eth list` is updated to display key specific max gas prices.
- CLI command `keys eth create` now supports optional `maxGasPriceGWei` parameter.
- CLI command `keys eth update` is added to update key specific parameters like `maxGasPriceGWei`.
- Add partial support for Moonriver chain
- For OCR jobs, `databaseTimeout`, `observationGracePeriod` and `contractTransmitterTransmitTimeout` can be specified to override chain-specific default values.

Two new log levels have been added.

- `[crit]`: *Critical* level logs are more severe than `[error]` and require quick action from the node operator.
- `[debug] [trace]`: *Trace* level logs contain extra `[debug]` information for development, and must be compiled in via `-tags trace`.

#### [Beta] Multichain support added

As a beta feature, Chainlink now supports connecting to multiple different EVM chains simultaneously.

This means that one node can run jobs on Goerli, Kovan, BSC and Mainnet (for example). Note that you can still have as many eth keys as you like, but each eth key is pegged to one chain only.

Extensive efforts have been made to make migration for existing nops as seamless as possible. Generally speaking, you should not have to make any changes when upgrading your existing node to this version. All your jobs will continue to run as before.

The overall summary of changes is such:

##### Chains/Ethereum Nodes

EVM chains are now represented as a first class object within the chainlink node. You can create/delete/list them using the CLI or API.

At least one primary node is required in order for a chain to connect. You may additionally specify zero or more send-only nodes for a chain. It is recommended to use the CLI/API or GUI to add nodes to chain.

###### Creation

```bash
chainlink chains evm create -id 42 # creates an evm chain with chain ID 42 (see: https://chainlist.org/)
chainlink nodes create -chain-id 42 -name 'my-primary-kovan-full-node' -type primary -ws-url ws://node.example/ws -http-url http://node.example/rpc # http-url is optional but recommended for primaries
chainlink nodes create -chain-id 42 -name 'my-send-only-backup-kovan-node' -type sendonly -http-url http://some-public-node.example/rpc
```

###### Listing

```bash
chainlink chains evm list
chainlink nodes list
```

###### Deletion

```bash
chainlink nodes delete 'my-send-only-backup-kovan-node'
chainlink chains evm delete 42
```

###### Legacy eth ENV vars

The old way of specifying chains using environment variables is still supported but discouraged. It works as follows:

If you specify `ETH_URL` then the values of `ETH_URL`, `ETH_CHAIN_ID`, `ETH_HTTP_URL` and `ETH_SECONDARY_URLS` will be used to create/update chains and nodes representing these values in the database. If an existing chain/node is found it will be overwritten. This behavior is used mainly to ease the process of upgrading, and on subsequent runs (once your old settings have been written to the database) it is recommended to unset these ENV vars and use the API commands exclusively to administer chains and nodes.

##### Jobs/tasks

By default, all jobs/tasks will continue to use the default chain (specified by `ETH_CHAIN_ID`). However, the following jobs now allow an additional `evmChainID` key in their TOML:

- VRF
- DirectRequest
- Keeper
- OCR
- Fluxmonitor

You can pin individual jobs to a particular chain by specifying the `evmChainID` explicitly. Here is an example job to demonstrate:

```toml
type            = "keeper"
evmChainID      = 3
schemaVersion   = 1
name            = "example keeper spec"
contractAddress = "0x9E40733cC9df84636505f4e6Db28DCa0dC5D1bba"
externalJobID   = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F49"
fromAddress     = "0xa8037A20989AFcBC51798de9762b351D63ff462e"
```

The above keeper job will _always_ run on chain ID 3 (Ropsten) regardless of the `ETH_CHAIN_ID` setting. If no chain matching this ID has been added to the chainlink node, the job cannot be created (you must create the chain first).

In addition, you can also specify `evmChainID` on certain pipeline tasks. This allows for cross-chain requests, for example:

```toml
type                = "directrequest"
schemaVersion       = 1
evmChainID          = 42
name                = "example cross chain spec"
contractAddress     = "0x613a38AC1659769640aaE063C651F48E0250454C"
externalJobID       = "0EEC7E1D-D0D2-476C-A1A8-72DFB6633F90"
observationSource   = """
    decode_log   [type=ethabidecodelog ... ]
    ...
    submit [type=ethtx to="0x613a38AC1659769640aaE063C651F48E0250454C" data="$(encode_tx)" minConfirmations="2" evmChainID="3"]
    decode_log-> ... ->submit;
"""
```

In the example above (which excludes irrelevant pipeline steps for brevity) a log can be read from the chain with ID 42 (Kovan) and a transaction emitted on chain with ID 3 (Ropsten).

Tasks that support the `evmChainID` parameter are as follows:

- `ethcall`
- `estimategaslimit`
- `ethtx`

###### Defaults

If the job- or task-specific `evmChainID` is _not_ given, the job/task will simply use the default as specified by the `ETH_CHAIN_ID` env variable.

Generally speaking, the default config values for each chain are good enough. But in some cases it is necessary to be able to override the defaults on a per-chain basis.

This used to be done via environment variables e.g. `MINIMUM_CONTRACT_PAYMENT_LINK_JUELS`.

These still work, but if set they will override that value for _all_ chains. This may not always be what you want. Consider a node that runs both Matic and Mainnet. You may want to set a higher value for `MINIMUM_CONTRACT_PAYMENT` on Mainnet, due to the more expensive gas costs. However, setting `MINIMUM_CONTRACT_PAYMENT_LINK_JUELS` using env variables will set that value for _all_ chains including matic.

To help you work around this, Chainlink now supports setting per-chain configuration options.

**Examples**

To set initial configuration when creating a chain, pass in the full json string as an optional parameter at the end:

`chainlink evm chains create -id 42 '{"BlockHistoryEstimatorBlockDelay": "100"}'`

To set configuration on an existing chain, specify key values pairs as such:

`chainlink evm chains configure -id 42 BlockHistoryEstimatorBlockDelay=100 GasEstimatorMode=FixedPrice`

The full list of chain-specific configuration options can be found by looking at the `ChainCfg` struct in `core/chains/evm/types/types.go`.

#### Async support in external adapters

External Adapters making async callbacks can now error job runs. This required a slight change to format, the correct way to callback from an asynchronous EA is using the following JSON:

SUCCESS CASE:

```json
{
    "value": < any valid json object >
}
```

ERROR CASE:


```json
{
    "error": "some error string"
}
```

This only applies to EAs using the `X-Chainlink-Pending` header to signal that the result will be POSTed back to the Chainlink node sometime 'later'. Regular synchronous calls to EAs work just as they always have done.

(NOTE: Official documentation for EAs needs to be updated)

#### New optional VRF v2 field: `requestedConfsDelay`

Added a new optional field for VRF v2 jobs called `requestedConfsDelay`, which configures a
number of blocks to wait in addition to the request specified `requestConfirmations` before servicing
the randomness request, i.e. the Chainlink node will wait `max(nodeMinConfs, requestConfirmations + requestedConfsDelay)`
blocks before servicing the request.

It can be used in the following way:

```toml
type = "vrf"
externalJobID = "123e4567-e89b-12d3-a456-426655440001"
schemaVersion = 1
name = "vrf-v2-secondary"
coordinatorAddress = "0xABA5eDc1a551E55b1A570c0e1f1055e5BE11eca7"
requestedConfsDelay = 10
# ... rest of job spec ...
```

Use of this field requires a database migration.

#### New locking mode: 'lease'

Chainlink now supports a new environment variable `DATABASE_LOCKING_MODE`. It can be set to one of the following values:

- `dual` (the default - uses both locking types for backwards and forwards compatibility)
- `advisorylock` (advisory lock only)
- `lease` (lease lock only)
- `none` (no locking at all - useful for advanced deployment environments when you can be sure that only one instance of chainlink will ever be running)

The database lock ensures that only one instance of Chainlink can be run on the database at a time. Running multiple instances of Chainlink on a single database at the same time would likely to lead to strange errors and possibly even data integrity failures and should not be allowed.

Ideally, node operators would be using a container orchestration system (e.g. Kubernetes) that ensures that only one instance of Chainlink ever runs on a particular postgres database.

However, we are aware that many node operators do not have the technical capacity to do this. So a common use case is to run multiple Chainlink instances in failover mode (as recommended by our official documentation, although this will be changing in future). The first instance will take some kind of lock on the database and subsequent instances will wait trying to take this lock in case the first instance disappears or dies.

Traditionally Chainlink has used an advisory lock to manage this. However, advisory locks come with several problems, notably:
- Postgres does not really like it when you hold locks open for a very long time (hours/days). It hampers certain internal cleanup tasks and is explicitly discouraged by the postgres maintainers.
- The advisory lock can silently disappear on postgres upgrade, meaning that a new instance can take over even while the old one is still running.
- Advisory locks do not play nicely with pooling tools such as pgbouncer.
- If the application crashes, the advisory lock can be left hanging around for a while (sometimes hours) and can require manual intervention to remove it before another instance of Chainlink will allow itself to boot.

For this reason, we have introduced a new locking mode, `lease`, which is likely to become the default in the future. `lease`-mode works as follows:
- Have one row in a database which is updated periodically with the client ID.
- CL node A will run a background process on start that updates this e.g. once per second.
- CL node B will spinlock, checking periodically to see if the update got too old. If it goes more than a set period without updating, it assumes that node A is dead and takes over. Now CL node B is the owner of the row, and it updates this every second.
- If CL node A comes back somehow, it will go to take out a lease and realise that the database has been leased to another process, so it will exit the entire application immediately.

The default is set to `dual` which used both advisory locking AND lease locking, for backwards compatibility. However, it is recommended that node operators who know what they are doing, or explicitly want to stop using the advisory locking mode set `DATABASE_LOCKING_MODE=lease` in their env.

Lease locking can be configured using the following ENV vars:

`LEASE_LOCK_REFRESH_INTERVAL` (default 1s)
`LEASE_LOCK_DURATION` (default 30s)

It is recommended to leave these set to the default values.

#### Duplicate Job Configuration

When duplicating a job, the new job's configuration settings that have not been overridden by the user can still reflect the chainlink node configuration.

#### Nurse (automatic pprof profiler)

Added new automatic pprof profiling service. Profiling is triggered when the node exceeds certain resource thresholds (currently, memory and goroutine count). The following environment variables have been added to allow configuring this service:

- `AUTO_PPROF_ENABLED`: Set to `true` to enable the automatic profiling service. Defaults to `false`.
- `AUTO_PPROF_PROFILE_ROOT`: The location on disk where pprof profiles will be stored. Defaults to `$CHAINLINK_ROOT`.
- `AUTO_PPROF_POLL_INTERVAL`: The interval at which the node's resources are checked. Defaults to `10s`.
- `AUTO_PPROF_GATHER_DURATION`: The duration for which profiles are gathered when profiling is kicked off. Defaults to `10s`.
- `AUTO_PPROF_GATHER_TRACE_DURATION`: The duration for which traces are gathered when profiling is kicked off. This is separately configurable because traces are significantly larger than other types of profiles. Defaults to `5s`.
- `AUTO_PPROF_MAX_PROFILE_SIZE`: The maximum amount of disk space that profiles may consume before profiling is disabled. Defaults to `100mb`.
- `AUTO_PPROF_CPU_PROFILE_RATE`: See https://pkg.go.dev/runtime#SetCPUProfileRate. Defaults to `1`.
- `AUTO_PPROF_MEM_PROFILE_RATE`: See https://pkg.go.dev/runtime#pkg-variables. Defaults to `1`.
- `AUTO_PPROF_BLOCK_PROFILE_RATE`: See https://pkg.go.dev/runtime#SetBlockProfileRate. Defaults to `1`.
- `AUTO_PPROF_MUTEX_PROFILE_FRACTION`: See https://pkg.go.dev/runtime#SetMutexProfileFraction. Defaults to `1`.
- `AUTO_PPROF_MEM_THRESHOLD`: The maximum amount of memory the node can actively consume before profiling begins. Defaults to `4gb`.
- `AUTO_PPROF_GOROUTINE_THRESHOLD`: The maximum number of actively-running goroutines the node can spawn before profiling begins. Defaults to `5000`.

**Adventurous node operators are encouraged to read [this guide on how to analyze pprof profiles](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/).**

#### `merge` task type

A new task type has been added, called `merge`. It can be used to merge two maps/JSON values together. Merge direction is from right to left such that `right` will clobber values of `left`. If no `left` is provided, it uses the input of the previous task. Example usage as such:


```
decode_log   [type=ethabidecodelog ...]
merge        [type=merge right=<{"foo": 42}>];

decode_log -> merge;
```

Or, to reverse merge direction:

```
decode_log   [type=ethabidecodelog ...]
merge        [type=merge left=<{"foo": 42}> right="$(decode_log)"];

decode_log -> merge;
```

#### Enhanced ABI encoding support

The `ethabiencode2` task supports ABI encoding using the abi specification generated by `solc`. e.g:

    {
        "name": "call",
        "inputs": [
          {
            "name": "value",
            "type": "tuple",
            "components": [
              {
                "name": "first",
                "type": "bytes32"
              },
              {
                "name": "last",
                "type": "bool"
              }
            ]
          }
        ]
    }

This would allow for calling of a function `call` with a tuple containing two values, the first a `bytes32` and the second a `bool`. You can supply a named map or an array.

#### Transaction Simulation (Gas Savings)

Chainlink now supports transaction simulation for certain types of job. When this is enabled, transactions will be simulated using `eth_call` before initial send. If the transaction reverted, the tx is marked as errored without being broadcast, potentially avoiding an expensive on-chain revert.

This can add a tiny bit of latency (upper bound 2s, generally much shorter under good conditions) and will add marginally more load to the eth client, since it adds an extra call for every transaction sent. However, it may help to save gas in some cases especially during periods of high demand by avoiding unnecessary reverts (due to outdated round etc.).

This option is EXPERIMENTAL and disabled by default.

To enable for FM or OCR:

`FM_SIMULATE_TRANSACTIONs=true`
`OCR_SIMULATE_TRANSACTIONS=true`

To enable in the pipeline, use the `simulate=true` option like so:

```
submit [type=ethtx to="0xDeadDeadDeadDeadDeadDeadDeadDead" data="0xDead" simulate=true]
```

Use at your own risk.

#### Misc

Chainlink now supports more than one primary eth node per chain. Requests are round-robined between available primaries.

Add CRUD functionality for EVM Chains and Nodes through Operator UI.

Non-fatal errors to a pipeline run are preserved including any run that succeeds but has more than one fatal error.

Chainlink now supports configuring max gas price on a per-key basis (allows implementation of keeper "lanes").

The Operator UI now supports login MFA with hardware security keys. `MFA_RPID` and `MFA_RPORIGIN` environment variables have been added to the config and are required if using the new MFA feature.

Keys and Configuration navigation links have been moved into a settings dropdown to make space for multichain navigation links.

#### Full EIP1559 Support (Gas Savings)

Chainlink now includes experimental support for submitting transactions using type 0x2 (EIP-1559) envelope.

EIP-1559 mode is off by default but can be enabled on a per-chain basis or globally.

This may help to save gas on spikes: Chainlink ought to react faster on the upleg and avoid overpaying on the downleg. It may also be possible to set `BLOCK_HISTORY_ESTIMATOR_BATCH_SIZE` to a smaller value e.g. 12 or even 6 because tip cap ought to be a more consistent indicator of inclusion time than total gas price. This would make Chainlink more responsive and ought to reduce response time variance. Some experimentation will be needed here to find optimum settings.

To enable globally, set `EVM_EIP1559_DYNAMIC_FEES=true`. Set with caution, if you set this on a chain that does not actually support EIP-1559 your node will be broken.

In EIP-1559 mode, the total price for the transaction is the minimum of base fee + tip cap and fee cap. More information can be found on the [official EIP](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1559.md).

Chainlink's implementation of this is to set a large fee cap and modify the tip cap to control confirmation speed of transactions. So, when in EIP1559 mode, the tip cap takes the place of gas price roughly speaking, with the varying base price remaining a constant (we always pay it).

A quick note on terminology - Chainlink uses the same terms used internally by go-ethereum source code to describe various prices. This is not the same as the externally used terms. For reference:

Base Fee Per Gas = BaseFeePerGas
Max Fee Per Gas = FeeCap
Max Priority Fee Per Gas = TipCap

In EIP-1559 mode, the following changes occur to how configuration works:

- All new transactions will be sent as type 0x2 transactions specifying a TipCap and FeeCap (NOTE: existing pending legacy transactions will continue to be gas bumped in legacy mode)
- BlockHistoryEstimator will apply its calculations (gas percentile etc.) to the TipCap and this value will be used for new transactions (GasPrice will be ignored)
- FixedPriceEstimator will use `EVM_GAS_TIP_CAP_DEFAULT` instead of `ETH_GAS_PRICE_DEFAULT`
- `ETH_GAS_PRICE_DEFAULT` is ignored for new transactions and `EVM_GAS_TIP_CAP_DEFAULT` is used instead (default 20GWei)
- `ETH_MIN_GAS_PRICE_WEI` is ignored for new transactions and `EVM_GAS_TIP_CAP_MINIMUM` is used instead (default 0)
- `ETH_MAX_GAS_PRICE_WEI` controls the FeeCap
- `KEEPER_GAS_PRICE_BUFFER_PERCENT` is ignored in EIP-1559 mode and `KEEPER_TIP_CAP_BUFFER_PERCENT` is used instead

The default tip cap is configurable per-chain but can be specified for all chains using `EVM_GAS_TIP_CAP_DEFAULT`. The fee cap is derived from `ETH_MAX_GAS_PRICE_WEI`.

When using the FixedPriceEstimator, the default gas tip will be used for all transactions.

When using the BlockHistoryEstimator, Chainlink will calculate the tip cap based on transactions already included (in the same way it calculates gas price in legacy mode).

Enabling EIP1559 mode might lead to marginally faster transaction inclusion and make the node more responsive to sharp rises/falls in gas price, keeping response times more consistent.

In addition, `ethcall` tasks now accept `gasTipCap` and `gasFeeCap` parameters in addition to `gasPrice`. This is required for Keeper jobs, i.e.:

```
check_upkeep_tx          [type=ethcall
                          failEarly=true
                          extractRevertReason=true
                          contract="$(jobSpec.contractAddress)"
                          gas="$(jobSpec.checkUpkeepGasLimit)"
                          gasPrice="$(jobSpec.gasPrice)"
                          gasTipCap="$(jobSpec.gasTipCap)"
                          gasFeeCap="$(jobSpec.gasFeeCap)"
                          data="$(encode_check_upkeep_tx)"]
```


NOTE: AccessLists are part of the 0x2 transaction type spec and Chainlink also implements support for these internally. This is not currently exposed in any way, if there is demand for this it ought to be straightforward enough to do so.

Avalanche AP4 defaults have been added (you can remove manually set ENV vars controlling gas pricing).

#### New env vars

`CHAIN_TYPE` - Configure the type of chain (if not standard). `Arbitrum`, `ExChain`, `Optimism`, or `XDai`. Replaces `LAYER_2_TYPE`. NOTE: This is a global override, to set on a per-chain basis you must use the CLI/API or GUI to change the chain-specific config for that chain (`ChainType`).

`BLOCK_EMISSION_IDLE_WARNING_THRESHOLD` - Controls global override for the time after which node will start logging warnings if no heads are received.

`ETH_DEFAULT_BATCH_SIZE` - Controls the default number of items per batch when making batched RPC calls. It is unlikely that you will need to change this from the default value.

NOTE: `ETH_URL` used to default to "ws://localhost:8546" and `ETH_CHAIN_ID` used to default to 1. These defaults have now been removed. The env vars are no longer required, since node configuration is now done via CLI/API/GUI and stored in the database.

### Removed

- `belt/` and `evm-test-helpers/` removed from the codebase.

#### Deprecated env vars

`LAYER_2_TYPE` - Use `CHAIN_TYPE` instead.

#### Removed env vars

`FEATURE_CRON_V2`, `FEATURE_FLUX_MONITOR_V2`, `FEATURE_WEBHOOK_V2` - all V2 job types are now enabled by default.

### Fixed

- Fixed a regression whereby the BlockHistoryEstimator would use a bumped value on old gas price even if the new current price was larger than the bumped value.
- Fixed a bug where creating lots of jobs very quickly in parallel would cause the node to hang
- Propagating `evmChainID` parameter in job specs supporting this parameter.

Fixed `LOG_LEVEL` behavior in respect to the corresponding UI setting: Operator can override `LOG_LEVEL` until the node is restarted.

### Changed

- The default `GAS_ESTIMATOR_MODE` for Optimism chains has been changed to `Optimism2`.
- Default minimum payment on mainnet has been reduced from 1 LINK to 0.1 LINK.
- Logging timestamp output has been changed from unix to ISO8601 to aid in readability. To keep the old unix format, you may set `LOG_UNIX_TS=true`
- Added WebAuthn support for the Operator UI and corresponding support in the Go backend

#### Log to Disk

This feature has been disabled by default, turn on with LOG_TO_DISK. For most production uses this is not desirable.

## [1.0.1] - 2021-11-23

### Added

- Improved error reporting
- Panic and recovery improvements

### Fixed

- Resolved config conversion errors for ETH_FINALITY_DEPTH, ETH_HEAD_TRACKER_HISTORY, and ETH_GAS_LIMIT_MULTIPLIER
- Proper handling for "nonce too low" errors on Avalanche

## [1.0.0] - 2021-10-19

### Added

- `chainlink node db status` will now display a table of applied and pending migrations.
- Add support for OKEx/ExChain.

### Changed

**Legacy job pipeline (JSON specs) are no longer supported**

This version will refuse to migrate the database if job specs are still present. You must manually delete or migrate all V1 job specs before upgrading.

For more information on migrating, see [the docs](https://docs.chain.link/chainlink-nodes/).

This release will DROP legacy job tables so please take a backup before upgrading.

#### KeyStore changes

* We no longer support "soft deleting", or archiving keys. From now on, keys can only be hard-deleted.
* Eth keys can no longer be imported directly to the database. If you with to import an eth key, you _must_ start the node first and import through the remote client.

#### New env vars

`LAYER_2_TYPE` - For layer 2 chains only. Configure the type of chain, either `Arbitrum` or `Optimism`.

#### Misc

- Head sampling can now be optionally disabled by setting `ETH_HEAD_TRACKER_SAMPLING_INTERVAL = "0s"` - this will result in every new head being delivered to running jobs,
  regardless of the head frequency from the chain.
- When creating new FluxMonitor jobs, the validation logic now checks that only one of: drumbeat ticker or idle timer is enabled.
- Added a new Prometheus metric: `uptime_seconds` which measures the number of seconds the node has been running. It can be helpful in detecting potential crashes.

### Fixed

Fixed a regression whereby the BlockHistoryEstimator would use a bumped value on old gas price even if the new current price was larger than the bumped value.

## [0.10.15] - 2021-10-14

**It is highly recommended upgrading to this version before upgrading to any newer versions to avoid any complications.**

### Fixed

- Prevent release from clobbering databases that have previously been upgraded

## [0.10.14] - 2021-09-06

### Added

- FMv2 spec now contains DrumbeatRandomDelay parameter that can be used to introduce variation between round of submits of different oracles, if drumbeat ticker is enabled.

- OCR Hibernation

#### Requesters/MinContractPaymentLinkJuels

V2 direct request specs now support two additional keys:

- "requesters" key which allows whitelisting requesters
- "minContractPaymentLinkJuels" key which allows to specify a job-specific minimum contract payment.

For example:

```toml
type                        = "directrequest"
schemaVersion               = 1
requesters                  = ["0xaaaa1F8ee20f5565510B84f9353F1E333E753B7a", "0xbbbb70F0e81C6F3430dfdC9fa02fB22BdD818C4e"] # optional
minContractPaymentLinkJuels = "100000000000000" # optional
name                        = "example eth request event spec with requesters"
contractAddress             = "..."
externalJobID               = "..."
observationSource           = """
...
"""
```

## [0.10.13] - 2021-08-25

### Fixed

- Resolved exiting Hibernation bug on FMv2

## [0.10.12] - 2021-08-16

### Fixed

- Resolved FMv2 stalling in Hibernation mode
- Resolved rare issue when the Gas Estimator fails on start
- Resolved the handling of nil values for gas price

## [0.10.11] - 2021-08-09

A new configuration variable, `BLOCK_BACKFILL_SKIP`, can be optionally set to "true" in order to strongly limit the depth of the log backfill.
This is useful if the node has been offline for a longer time and after startup should not be concerned with older events from the chain.

Three new configuration variables are added for the new telemetry ingress service support. `TELEMETRY_INGRESS_URL` sets the URL to connect to for telemetry ingress, `TELEMETRY_INGRESS_SERVER_PUB_KEY` sets the public key of the telemetry ingress server, and `TELEMETRY_INGRESS_LOGGING` toggles verbose logging of the raw telemetry messages being sent.

* Fixes the logging configuration form not displaying the current values
* Updates the design of the configuration cards to be easier on the eyes
* View Coordinator Service Authentication keys in the Operator UI. This is hidden
  behind a feature flag until usage is enabled.
* Adds support for the new telemetry ingress service.

### Changed

**The legacy job pipeline (JSON specs) has been officially deprecated and support for these jobs will be dropped in an upcoming release.**

Any node operators still running jobs with JSON specs should migrate their jobs to TOML format instead.

The format for V2 Webhook job specs has changed. They now allow specifying 0 or more external initiators. Example below:

```toml
type            = "webhook"
schemaVersion   = 1
externalInitiators = [
    { name = "foo-ei", spec = '{"foo": 42}' },
    { name = "bar-ei", spec = '{"bar": 42}' }
]
observationSource   = """
ds          [type=http method=GET url="https://chain.link/ETH-USD"];
ds_parse    [type=jsonparse path="data,price"];
ds_multiply [type=multiply times=100];
ds -> ds_parse -> ds_multiply;
"""
```

These external initiators will be notified with the given spec after the job is created, and also at deletion time.

Only the External Initiators listed in the toml spec may trigger a run for that job. Logged-in users can always trigger a run for any job.

#### Migrating Jobs

- OCR
All OCR jobs are already using v2 pipeline by default - no need to do anything here.

- Flux Monitor v1
We have created a tool to help you automigrate flux monitor specs in JSON format to the new TOML format. You can migrate a job like this:

```
chainlink jobs migrate <job id>
```

This can be automated by using the API like so:

```
POST http://yournode.example/v2/migrate/<job id>
```

- VRF v1
Automigration is not supported for VRF jobs. They must be manually converted into v2 format.

- Ethlog/Runlog/Cron/web
All other job types must also be manually converted into v2 format.

#### Technical details

Why are we doing this?

To give some background, the legacy job pipeline has been around since before Chainlink went to mainnet and is getting quite long in the tooth. The code is brittle and difficult to understand and maintain. For a while now we have been developing a v2 job pipeline in parallel which uses the TOML format. The new job pipeline is simpler, more performant and more powerful. Every job that can be represented in the legacy pipeline should be able to be represented in the v2 pipeline - if it can't be, that's a bug, so please let us know ASAP.

The v2 pipeline has now been extensively tested in production and proved itself reliable. So, we made the decision to drop V1 support entirely in favour of focusing developer effort on new features like native multichain support, EIP1559-compatible fees, further gas saving measures and support for more blockchains. By dropping support for the old pipeline, we can deliver these features faster and better support our community.

#### KeyStore changes

* Key export files are changing format and will not be compatible between versions. Ex, a key exported in 0.10.12, will not be importable by a node running 1.0.0, and vice-versa.
* We no longer support "soft deleting", or archiving keys. From now on, keys can only be hard-deleted.
* Eth keys can no longer be imported directly to the database. If you with to import an eth key, you _must_ start the node first and import through the remote client.

## [0.10.10] - 2021-07-19

### Changed

This update will truncate `pipeline_runs`, `pipeline_task_runs`, `flux_monitor_round_stats_v2` DB tables as a part of the migration.

#### Gas Estimation

Gas estimation has been revamped and full support for Optimism has been added.

The following env vars have been deprecated, and will be removed in a future release:

```
GAS_UPDATER_ENABLED
GAS_UPDATER_BATCH_SIZE
GAS_UPDATER_BLOCK_DELAY
GAS_UPDATER_BLOCK_HISTORY_SIZE
GAS_UPDATER_TRANSACTION_PERCENTILE
```

If you are using any of the env vars above, please switch to using the following instead:

```
GAS_ESTIMATOR_MODE
BLOCK_HISTORY_ESTIMATOR_BATCH_SIZE
BLOCK_HISTORY_ESTIMATOR_BLOCK_DELAY
BLOCK_HISTORY_ESTIMATOR_BLOCK_HISTORY_SIZE
BLOCK_HISTORY_ESTIMATOR_TRANSACTION_PERCENTILE
```

Valid values for `GAS_ESTIMATOR_MODE` are as follows:

`GAS_ESTIMATOR_MODE=BlockHistory` (equivalent to `GAS_UPDATER_ENABLED=true`)
`GAS_ESTIMATOR_MODE=FixedPrice` (equivalent to `GAS_UPDATER_ENABLED=false`)
`GAS_ESTIMATOR_MODE=Optimism` (new)

New gas estimator modes may be added in the future.

In addition, a minor annoyance has been fixed whereby previously if you enabled the gas updater, it would overwrite the locally stored value for gas price and continue to use this even if it was disabled after a reboot. This will no longer happen: BlockHistory mode will not clobber the locally stored value for fixed gas price, which can still be adjusted via remote API call or using `chainlink config setgasprice XXX`. In order to use this manually fixed gas price, you must enable FixedPrice estimator mode.

### Added

Added support for latest version of libocr with the V2 networking stack. New env vars to configure this are:

```
P2P_NETWORKING_STACK
P2PV2_ANNOUNCE_ADDRESSES
P2PV2_BOOTSTRAPPERS
P2PV2_DELTA_DIAL
P2PV2_DELTA_RECONCILE
P2PV2_LISTEN_ADDRESSES
```

All of these are currently optional, by default OCR will continue to use the existing V1 stack. The new env vars will be used internally for OCR testing.

### Fixed

- Fix inability to create jobs with a cron schedule.

## [0.10.9] - 2021-07-05

### Changed

#### Transaction Strategies

FMv2, Keeper and OCR jobs now use a new strategy for sending transactions. By default, if multiple transactions are queued up, only the latest one will be sent. This should greatly reduce the number of stale rounds and reverted transactions, and help node operators to save significant gas especially during times of high congestion or when catching up on a deep backlog.

Defaults should work well, but it can be controlled if necessary using the following new env vars:

`FM_DEFAULT_TRANSACTION_QUEUE_DEPTH`
`KEEPER_DEFAULT_TRANSACTION_QUEUE_DEPTH`
`OCR_DEFAULT_TRANSACTION_QUEUE_DEPTH`

Setting to 0 will disable (the old behaviour). Setting to 1 (the default) will keep only the latest transaction queued up at any given time. Setting to 2, 3 etc. will allow this many transactions to be queued before starting to drop older items.

Note that it has no effect on FMv1 jobs. Node operators will need to upgrade to FMv2 to take advantage of this feature.

## [0.10.8] - 2021-06-21

### Fixed

- The HTTP adapter would remove a trailing slash on a subdirectory when specifying an extended path, so for instance `http://example.com/subdir/` with a param of `?query=` extended path would produce the URL `http://example.com/subdir?query=`, but should now produce: `http://example.com/subdir/?query=`.

- Matic autoconfig is now enabled for mainnet. Matic nops should remove any custom tweaks they have been running with. In addition, we have better default configs for Optimism, Arbitrum and RSK.

- It is no longer required to set `DEFAULT_HTTP_ALLOW_UNRESTRICTED_NETWORK_ACCESS=true` to enable local fetches on bridge or http tasks. If the URL for the http task is specified as a variable, then set the AllowUnrestrictedNetworkAccess option for this task. Please remove this if you had it set and no longer need it, since it introduces a slight security risk.

- Chainlink can now run with ETH_DISABLED=true without spewing errors everywhere

- Removed prometheus metrics that were no longer valid after recent changes to head tracking:
  `head_tracker_heads_in_queue`, `head_tracker_callback_execution_duration`,
  `head_tracker_callback_execution_duration_hist`, `head_tracker_num_heads_dropped`

### Added

- MINIMUM_CONTRACT_PAYMENT_LINK_JUELS replaces MINIMUM_CONTRACT_PAYMENT, which will be deprecated in a future release.

- INSECURE_SKIP_VERIFY configuration variable disables verification of the Chainlink SSL certificates when using the CLI.

- JSON parse tasks (v2) now permit an empty `path` parameter.

- Eth->eth transfer gas limit is no longer hardcoded at 21000 and can now be adjusted using `ETH_GAS_LIMIT_TRANSFER`

- HTTP and Bridge tasks (v2 pipeline) now log the request parameters (including the body) upon making the request when `LOG_LEVEL=debug`.

- Webhook v2 jobs now support two new parameters, `externalInitiatorName` and `externalInitiatorSpec`. The v2 version of the following v1 spec:
    ```
    {
      "initiators": [
        {
          "type": "external",
          "params": {
            "name": "substrate",
            "body": {
              "endpoint": "substrate",
              "feed_id": 0,
              "account_id": "0x7c522c8273973e7bcf4a5dbfcc745dba4a3ab08c1e410167d7b1bdf9cb924f6c",
              "fluxmonitor": {
                "requestData": {
                  "data": { "from": "DOT", "to": "USD" }
                },
                "feeds": [{ "url": "http://adapter1:8080" }],
                "threshold": 0.5,
                "absoluteThreshold": 0,
                "precision": 8,
                "pollTimer": { "period": "30s" },
                "idleTimer": { "duration": "1m" }
              }
            }
          }
        }
      ],
      "tasks": [
        {
          "type": "substrate-adapter1",
          "params": { "multiply": 1e8 }
        }
      ]
    }
    ```
    is:
    ```
    type            = "webhook"
    schemaVersion   = 1
    jobID           = "0EEC7E1D-D0D2-475C-A1A8-72DFB6633F46"
    externalInitiatorName = "substrate"
    externalInitiatorSpec = """
        {
          "endpoint": "substrate",
          "feed_id": 0,
          "account_id": "0x7c522c8273973e7bcf4a5dbfcc745dba4a3ab08c1e410167d7b1bdf9cb924f6c",
          "fluxmonitor": {
            "requestData": {
              "data": { "from": "DOT", "to": "USD" }
            },
            "feeds": [{ "url": "http://adapter1:8080" }],
            "threshold": 0.5,
            "absoluteThreshold": 0,
            "precision": 8,
            "pollTimer": { "period": "30s" },
            "idleTimer": { "duration": "1m" }
          }
        }
    """
    observationSource   = """
        submit [type=bridge name="substrate-adapter1" requestData=<{ "multiply": 1e8 }>]
    """
    ```


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

- Keeper jobs now support prometheus metrics, they are considered a pipeline with a single `keeper` task type. Example:
```
pipeline_run_errors{job_id="1",job_name="example keeper spec"} 1
pipeline_run_total_time_to_completion{job_id="1",job_name="example keeper spec"} 8.470456e+06
pipeline_task_execution_time{job_id="1",job_name="example keeper spec",task_type="keeper"} 8.470456e+06
pipeline_tasks_total_finished{job_id="1",job_name="example keeper spec",status="completed",task_type="keeper"} 1
```

### Changed

- The v2 (TOML) `bridge` task's `includeInputAtKey` parameter is being deprecated in favor of variable interpolation. Please migrate your jobs to the new syntax as soon as possible.

- Chainlink no longer writes/reads eth key files to disk

- Add sensible default configuration settings for Fantom

- Rename `ETH_MAX_UNCONFIRMED_TRANSACTIONS` to `ETH_MAX_QUEUED_TRANSACTIONS`. It still performs the same function but the name was misleading and would have caused confusion with the new `ETH_MAX_IN_FLIGHT_TRANSACTIONS`.

- The VRF keys are now managed remotely through the node only. Example commands:
```
// Starting a node with a vrf key
chainlink node start -p path/to/passwordfile -vp path/to/vrfpasswordfile

// Remotely managing the vrf keys
chainlink keys vrf create // Creates a key with path/to/vrfpasswordfile
chainlink keys vrf list // Lists all keys on the node
chainlink keys vrf delete // Lists all keys on the node

// Archives (soft deletes) vrf key with compressed pub key 0x788..
chainlink keys vrf delete 0x78845e23b6b22c47e4c81426fdf6fc4087c4c6a6443eba90eb92cf4d11c32d3e00

// Hard deletes vrf key with compressed pub key 0x788..
chainlink keys vrf delete 0x78845e23b6b22c47e4c81426fdf6fc4087c4c6a6443eba90eb92cf4d11c32d3e00 --hard

// Exports 0x788.. key to file 0x788_exported_key on disk encrypted with path/to/vrfpasswordfile
// Note you can re-encrypt it with a different password if you like when exporting.
chainlink keys vrf export 0x78845e23b6b22c47e4c81426fdf6fc4087c4c6a6443eba90eb92cf4d11c32d3e00 -p path/to/vrfpasswordfile -o 0x788_exported_key

// Import key material in 0x788_exported_key using path/to/vrfpasswordfile to decrypt.
// Will be re-encrypted with the nodes vrf password file i.e. "-vp"
chainlink keys vrf import -p path/to/vrfpasswordfile 0x788_exported_key
```



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
  to any chosen directory by setting a new configuration value: `DATABASE_BACKUP_DIR`

## [0.10.6] - 2021-05-10

### Added

- Add `MockOracle.sol` for testing contracts

- Web job types can now be created from the operator UI as a new job.

- See example web job spec below:

```
type            = "webhook"
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
update initiators set deleted_at = null where job_spec_id = 'fe279ed9-c36f-4eef-9dc1-bdb7bef21264'
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
- Fixed an issue where expired session tokens in operator UI would cause a large number of requests to be sent to the node, resulting in a temporary rate-limit and 429 errors.
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
- OCR jobs no longer require defining v1 bootstrap peers unless `P2P_NETWORKING_STACK=V1`

#### BREAKING CHANGES

- Commands for creating/managing legacy jobs and OCR jobs have changed, to reduce confusion and accommodate additional types of jobs using the new pipeline.
- If `P2P_NETWORKING_STACK=V1V2`, then `P2PV2_BOOTSTRAPPERS` must also be set

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
- As a side effect, we now no longer handle the case where an external wallet used the chainlink ethereum private key to send a transaction. This use-case was already explicitly unsupported, but we made a best-effort attempt to handle it. We now make no attempt at all to handle it and doing this WILL result in your node not sending the data that it expected to be sent for the nonces that were used by an external wallet.
- Operator UI now shows booleans correctly

### Changed

- ETH_MAX_GAS_PRICE_WEI now 1500Gwei by default

## [0.8.18] - 2020-10-01

### Fixed

- Prometheus gas_updater_set_gas_price metric now only shows last gas price instead of every block since restart

## [0.8.17] - 2020-09-28

### Added

- Add new env variable ETH_SECONDARY_URL. Default is unset. You may optionally set this to a http(s) ethereum RPC client URL. If set, transactions will also be broadcast to this secondary ethereum node. This allows transaction broadcasting to be more robust in the face of primary ethereum node bugs or failures.
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

IMPORTANT: The selection mechanism for keys has changed. When an ethtx task spec is not pinned to a particular key by defining `fromAddress` or `fromAddresses`, the node will now cycle through all available keys in round-robin fashion. This is a change from the previous behaviour where nodes would only pick the earliest created key.

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

- `ethtx` tasks now support a new parameter, `minRequiredOutgoingConfirmations` which allows you to tune how many confirmations are required before moving on from an `ethtx` task on a per-task basis (only works with BulletproofTxManager). If it is not supplied, the default of `MIN_OUTGOING_CONFIRMATIONS` is used (same as the old behaviour).

### Changed

- HeadTracker now automatically backfills missing heads up to `ETH_FINALITY_DEPTH`
- The strategy for gas bumping has been changed to produce a potentially higher gas cost in exchange for the transaction getting through faster.

### Breaking changes

- `admin withdraw` command has been removed. This was only ever useful to withdraw LINK if the Oracle contract was owned by the Chainlink node address. It is no longer recommended having the Oracle owner be the chainlink node address.
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
- `ETH_FINALITY_DEPTH` specifies how deep protection should be against re-orgs. The default is 50. It only applies if BulletproofTxManager is enabled. It is not recommended changing this setting.
- `EthHeadTrackerHistoryDepth` specifies how many heads the head tracker should keep in the database. The default is 100. It is not recommended changing this setting.
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
