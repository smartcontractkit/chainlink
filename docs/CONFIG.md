[//]: # (Documentation generated from docs/*.toml - DO NOT EDIT.)

This document describes the TOML format for configuration.

See also [SECRETS.md](SECRETS.md)

## Example

```toml
Log.Level = 'debug'

[[EVM]]
ChainID = '1' # Required

[[EVM.Nodes]]
Name = 'fake' # Required
WSURL = 'wss://foo.bar/ws'
HTTPURL = 'https://foo.bar' # Required
```

## Global
```toml
ExplorerURL = 'ws://explorer.url' # Example
InsecureFastScrypt = false # Default
RootDir = '~/.chainlink' # Default
ShutdownGracePeriod = '5s' # Default
```


### ExplorerURL
```toml
ExplorerURL = 'ws://explorer.url' # Example
```
ExplorerURL is the websocket URL used by the node to push stats. This variable is required to deliver telemetry.

### InsecureFastScrypt
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
InsecureFastScrypt = false # Default
```
InsecureFastScrypt causes all key stores to encrypt using "fast" scrypt params instead. This is insecure and only useful for local testing. DO NOT ENABLE THIS IN PRODUCTION.

### RootDir
```toml
RootDir = '~/.chainlink' # Default
```
RootDir is the Chainlink node's root directory. This is the default directory for logging, database backups, cookies, and other misc Chainlink node files. Chainlink nodes will always ensure this directory has 700 permissions because it might contain sensitive data.

### ShutdownGracePeriod
```toml
ShutdownGracePeriod = '5s' # Default
```
ShutdownGracePeriod is the maximum time allowed to shut down gracefully. If exceeded, the node will terminate immediately to avoid being SIGKILLed.

## Feature
```toml
[Feature]
FeedsManager = true # Default
LogPoller = false # Default
UICSAKeys = false # Default
```


### FeedsManager
```toml
FeedsManager = true # Default
```
FeedsManager enables the feeds manager service.

### LogPoller
```toml
LogPoller = false # Default
```
LogPoller enables the log poller, an experimental approach to processing logs, required if also using Evm.UseForwarders or OCR2.

### UICSAKeys
```toml
UICSAKeys = false # Default
```
UICSAKeys enables CSA Keys in the UI.

## Database
```toml
[Database]
DefaultIdleInTxSessionTimeout = '1h' # Default
DefaultLockTimeout = '15s' # Default
DefaultQueryTimeout = '10s' # Default
LogQueries = false # Default
MaxIdleConns = 10 # Default
MaxOpenConns = 20 # Default
MigrateOnStartup = true # Default
```


### DefaultIdleInTxSessionTimeout
```toml
DefaultIdleInTxSessionTimeout = '1h' # Default
```
DefaultIdleInTxSessionTimeout is the maximum time allowed for a transaction to be open and idle before timing out. See Postgres `idle_in_transaction_session_timeout` for more details.

### DefaultLockTimeout
```toml
DefaultLockTimeout = '15s' # Default
```
DefaultLockTimeout is the maximum time allowed to wait for database lock of any kind before timing out. See Postgres `lock_timeout` for more details.

### DefaultQueryTimeout
```toml
DefaultQueryTimeout = '10s' # Default
```
DefaultQueryTimeout is the maximum time allowed for standard queries before timing out.

### LogQueries
```toml
LogQueries = false # Default
```
LogQueries tells the Chainlink node to log database queries made using the default logger. SQL statements will be logged at `debug` level. Not all statements can be logged. The best way to get a true log of all SQL statements is to enable SQL statement logging on Postgres.

### MaxIdleConns
```toml
MaxIdleConns = 10 # Default
```
MaxIdleConns configures the maximum number of idle database connections that the Chainlink node will keep open. Think of this as the baseline number of database connections per Chainlink node instance. Increasing this number can help to improve performance under database-heavy workloads.

Postgres has connection limits, so you must use caution when increasing this value. If you are running several instances of a Chainlink node or another application on a single database server, you might run out of Postgres connection slots if you raise this value too high.

### MaxOpenConns
```toml
MaxOpenConns = 20 # Default
```
MaxOpenConns configures the maximum number of database connections that a Chainlink node will have open at any one time. Think of this as the maximum burst upper bound limit of database connections per Chainlink node instance. Increasing this number can help to improve performance under database-heavy workloads.

Postgres has connection limits, so you must use caution when increasing this value. If you are running several instances of a Chainlink node or another application on a single database server, you might run out of Postgres connection slots if you raise this value too high.

### MigrateOnStartup
```toml
MigrateOnStartup = true # Default
```
MigrateOnStartup controls whether a Chainlink node will attempt to automatically migrate the database on boot. If you want more control over your database migration process, set this variable to `false` and manually migrate the database using the CLI `migrate` command instead.

## Database.Backup
```toml
[Database.Backup]
Mode = 'none' # Default
Dir = 'test/backup/dir' # Example
OnVersionUpgrade = true # Default
Frequency = '1h' # Default
```
As a best practice, take regular database backups in case of accidental data loss. This best practice is especially important when you upgrade your Chainlink node to a new version. Chainlink nodes support automated database backups to make this process easier.

NOTE: Dumps can cause high load and massive database latencies, which will negatively impact the normal functioning of the Chainlink node. For this reason, it is recommended to set a `URL` and point it to a read replica if you enable automatic backups.

### Mode
```toml
Mode = 'none' # Default
```
Mode sets the type of automatic database backup, which can be one of _none_, `lite`, or `full`. If enabled, the Chainlink node will always dump a backup on every boot before running migrations. Additionally, it will automatically take database backups that overwrite the backup file for the given version at regular intervals if `Frequency` is set to a non-zero interval.

_none_ - Disables backups.
`lite` - Dumps small tables including configuration and keys that are essential for the node to function, which excludes historical data like job runs, transaction history, etc.
`full` - Dumps the entire database.

It will write to a file like `'Dir'/backup/cl_backup_<VERSION>.dump`. There is one backup dump file per version of the Chainlink node. If you upgrade the node, it will keep the backup taken right before the upgrade migration so you can restore to an older version if necessary.

### Dir
```toml
Dir = 'test/backup/dir' # Example
```
Dir sets the directory to use for saving the backup file. Use this if you want to save the backup file in a directory other than the default ROOT directory.

### OnVersionUpgrade
```toml
OnVersionUpgrade = true # Default
```
OnVersionUpgrade enables automatic backups of the database before running migrations, when you are upgrading to a new version.

### Frequency
```toml
Frequency = '1h' # Default
```
Frequency sets the interval for database dumps, if set to a positive duration and `Mode` is not _none_.

Set to `0` to disable periodic backups.

## Database.Listener
:warning: **_ADVANCED_**: _Do not change these settings unless you know what you are doing._
```toml
[Database.Listener]
MaxReconnectDuration = '10m' # Default
MinReconnectInterval = '1m' # Default
FallbackPollInterval = '30s' # Default
```
These settings control the postgres event listener.

### MaxReconnectDuration
```toml
MaxReconnectDuration = '10m' # Default
```
MaxReconnectDuration is the maximum duration to wait between reconnect attempts.

### MinReconnectInterval
```toml
MinReconnectInterval = '1m' # Default
```
MinReconnectInterval controls the duration to wait before trying to re-establish the database connection after connection loss. After each consecutive failure this interval is doubled, until MaxReconnectInterval is reached.  Successfully completing the connection establishment procedure resets the interval back to MinReconnectInterval.

### FallbackPollInterval
```toml
FallbackPollInterval = '30s' # Default
```
FallbackPollInterval controls how often clients should manually poll as a fallback in case the postgres event was missed/dropped.

## Database.Lock
:warning: **_ADVANCED_**: _Do not change these settings unless you know what you are doing._
```toml
[Database.Lock]
Enabled = true # Default
LeaseDuration = '10s' # Default
LeaseRefreshInterval = '1s' # Default
```
Ideally, you should use a container orchestration system like [Kubernetes](https://kubernetes.io/) to ensure that only one Chainlink node instance can ever use a specific Postgres database. However, some node operators do not have the technical capacity to do this. Common use cases run multiple Chainlink node instances in failover mode as recommended by our official documentation. The first instance takes a lock on the database and subsequent instances will wait trying to take this lock in case the first instance fails.

- If your nodes or applications hold locks open for several hours or days, Postgres is unable to complete internal cleanup tasks. The Postgres maintainers explicitly discourage holding locks open for long periods of time.

Because of the complications with advisory locks, Chainlink nodes with v2.0 and later only support `lease` locking mode. The `lease` locking mode works using the following process:

- Node A creates one row in the database with the client ID and updates it once per second.
- Node B spinlocks and checks periodically to see if the client ID is too old. If the client ID is not updated after a period of time, node B assumes that node A failed and takes over. Node B becomes the owner of the row and updates the client ID once per second.
- If node A comes back, it attempts to take out a lease, realizes that the database has been leased to another process, and exits the entire application immediately.

### Enabled
```toml
Enabled = true # Default
```
Enabled enables the database lock.

### LeaseDuration
```toml
LeaseDuration = '10s' # Default
```
LeaseDuration is how long the lease lock will last before expiring.

### LeaseRefreshInterval
```toml
LeaseRefreshInterval = '1s' # Default
```
LeaseRefreshInterval determines how often to refresh the lease lock. Also controls how often a standby node will check to see if it can grab the lease.

## TelemetryIngress
```toml
[TelemetryIngress]
UniConn = true # Default
Logging = false # Default
ServerPubKey = 'test-pub-key' # Example
URL = 'https://prom.test' # Example
BufferSize = 100 # Default
MaxBatchSize = 50 # Default
SendInterval = '500ms' # Default
SendTimeout = '10s' # Default
UseBatchSend = true # Default
```


### UniConn
```toml
UniConn = true # Default
```
UniConn toggles which ws connection style is used.

### Logging
```toml
Logging = false # Default
```
Logging toggles verbose logging of the raw telemetry messages being sent.

### ServerPubKey
```toml
ServerPubKey = 'test-pub-key' # Example
```
ServerPubKey is the public key of the telemetry server.

### URL
```toml
URL = 'https://prom.test' # Example
```
URL is where to send telemetry.

### BufferSize
```toml
BufferSize = 100 # Default
```
BufferSize is the number of telemetry messages to buffer before dropping new ones.

### MaxBatchSize
```toml
MaxBatchSize = 50 # Default
```
MaxBatchSize is the maximum number of messages to batch into one telemetry request.

### SendInterval
```toml
SendInterval = '500ms' # Default
```
SendInterval determines how often batched telemetry is sent to the ingress server.

### SendTimeout
```toml
SendTimeout = '10s' # Default
```
SendTimeout is the max duration to wait for the request to complete when sending batch telemetry.

### UseBatchSend
```toml
UseBatchSend = true # Default
```
UseBatchSend toggles sending telemetry to the ingress server using the batch client.

## AuditLogger
```toml
[AuditLogger]
Enabled = false # Default
ForwardToUrl = 'http://localhost:9898' # Example
JsonWrapperKey = 'event' # Example
Headers = ['Authorization: token', 'X-SomeOther-Header: value with spaces | and a bar+*'] # Example
```


### Enabled
```toml
Enabled = false # Default
```
Enabled determines if this logger should be configured at all

### ForwardToUrl
```toml
ForwardToUrl = 'http://localhost:9898' # Example
```
ForwardToUrl is where you want to forward logs to

### JsonWrapperKey
```toml
JsonWrapperKey = 'event' # Example
```
JsonWrapperKey if set wraps the map of data under another single key to make parsing easier

### Headers
```toml
Headers = ['Authorization: token', 'X-SomeOther-Header: value with spaces | and a bar+*'] # Example
```
Headers is the set of headers you wish to pass along with each request

## Log
```toml
[Log]
Level = 'info' # Default
JSONConsole = false # Default
UnixTS = false # Default
```


### Level
```toml
Level = 'info' # Default
```
Level determines both what is printed on the screen and what is written to the log file.

The available levels are:
- "debug": Useful for forensic debugging of issues.
- "info": High-level informational messages. (default)
- "warn": A mild error occurred that might require non-urgent action. Check these warnings semi-regularly to see if any of them require attention. These warnings usually happen due to factors outside of the control of the node operator. Examples: Unexpected responses from a remote API or misleading networking errors.
- "error": An unexpected error occurred during the regular operation of a well-maintained node. Node operators might need to take action to remedy this error. Check these regularly to see if any of them require attention. Examples: Use of deprecated configuration options or incorrectly configured settings that cause a job to fail.
- "crit": A critical error occurred. The node might be unable to function. Node operators should take immediate action to fix these errors. Examples: The node could not boot because a network socket could not be opened or the database became inaccessible.
- "panic": An exceptional error occurred that could not be handled. If the node is unresponsive, node operators should try to restart their nodes and notify the Chainlink team of a potential bug.
- "fatal": The node encountered an unrecoverable problem and had to exit.

### JSONConsole
```toml
JSONConsole = false # Default
```
JSONConsole enables JSON logging. Otherwise, the log is saved in a human-friendly console format.

### UnixTS
```toml
UnixTS = false # Default
```
UnixTS enables legacy unix timestamps.

Previous versions of Chainlink nodes wrote JSON logs with a unix timestamp. As of v1.1.0 and up, the default has changed to use ISO8601 timestamps for better readability.

## Log.File
```toml
[Log.File]
Dir = '/my/log/directory' # Example
MaxSize = '5120mb' # Default
MaxAgeDays = 0 # Default
MaxBackups = 1 # Default
```


### Dir
```toml
Dir = '/my/log/directory' # Example
```
Dir sets the log directory. By default, Chainlink nodes write log data to `$ROOT/log.jsonl`.

### MaxSize
```toml
MaxSize = '5120mb' # Default
```
MaxSize determines the log file's max size in megabytes before file rotation. Having this not set will disable logging to disk. If your disk doesn't have enough disk space, the logging will pause and the application will log errors until space is available again.

Values must have suffixes with a unit like: `5120mb` (5,120 megabytes). If no unit suffix is provided, the value defaults to `b` (bytes). The list of valid unit suffixes are:

- b (bytes)
- kb (kilobytes)
- mb (megabytes)
- gb (gigabytes)
- tb (terabytes)

### MaxAgeDays
```toml
MaxAgeDays = 0 # Default
```
MaxAgeDays determines the log file's max age in days before file rotation. Keeping this config with the default value will not remove log files based on age.

### MaxBackups
```toml
MaxBackups = 1 # Default
```
MaxBackups determines the maximum number of old log files to retain. Keeping this config with the default value retains all old log files. The `MaxAgeDays` variable can still cause them to get deleted.

## WebServer
```toml
[WebServer]
AllowOrigins = 'http://localhost:3000,http://localhost:6688' # Default
BridgeCacheTTL = '0s' # Default
BridgeResponseURL = 'https://my-chainlink-node.example.com:6688' # Example
HTTPWriteTimeout = '10s' # Default
HTTPPort = 6688 # Default
SecureCookies = true # Default
SessionTimeout = '15m' # Default
SessionReaperExpiration = '240h' # Default
```


### AllowOrigins
```toml
AllowOrigins = 'http://localhost:3000,http://localhost:6688' # Default
```
AllowOrigins controls the URLs Chainlink nodes emit in the `Allow-Origins` header of its API responses. The setting can be a comma-separated list with no spaces. You might experience CORS issues if this is not set correctly.

You should set this to the external URL that you use to access the Chainlink UI.

You can set `AllowOrigins = '*'` to allow the UI to work from any URL, but it is recommended for security reasons to make it explicit instead.

### BridgeCacheTTL
```toml
BridgeCacheTTL = '0s' # Default
```
BridgeCacheTTL controls the cache TTL for all bridge tasks to use old values in newer observations in case of intermittent failure. It's disabled by default.

### BridgeResponseURL
```toml
BridgeResponseURL = 'https://my-chainlink-node.example.com:6688' # Example
```
BridgeResponseURL defines the URL for bridges to send a response to. This _must_ be set when using async external adapters.

Usually this will be the same as the URL/IP and port you use to connect to the Chainlink UI.

### HTTPWriteTimeout
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
HTTPWriteTimeout = '10s' # Default
```
HTTPWriteTimeout controls how long the Chainlink node's API server can hold a socket open for writing a response to an HTTP request. Sometimes, this must be increased for pprof.

### HTTPPort
```toml
HTTPPort = 6688 # Default
```
HTTPPort is the port used for the Chainlink Node API, [CLI](/docs/configuration-variables/#cli-client), and GUI.

### SecureCookies
```toml
SecureCookies = true # Default
```
SecureCookies requires the use of secure cookies for authentication. Set to false to enable standard HTTP requests along with `TLSPort = 0`.

### SessionTimeout
```toml
SessionTimeout = '15m' # Default
```
SessionTimeout determines the amount of idle time to elapse before session cookies expire. This signs out GUI users from their sessions.

### SessionReaperExpiration
```toml
SessionReaperExpiration = '240h' # Default
```
SessionReaperExpiration represents how long an API session lasts before expiring and requiring a new login.

## WebServer.RateLimit
```toml
[WebServer.RateLimit]
Authenticated = 1000 # Default
AuthenticatedPeriod = '1m' # Default
Unauthenticated = 5 # Default
UnauthenticatedPeriod = '20s' # Default
```


### Authenticated
```toml
Authenticated = 1000 # Default
```
Authenticated defines the threshold to which authenticated requests get limited. More than this many authenticated requests per `AuthenticatedRateLimitPeriod` will be rejected.

### AuthenticatedPeriod
```toml
AuthenticatedPeriod = '1m' # Default
```
AuthenticatedPeriod defines the period to which authenticated requests get limited.

### Unauthenticated
```toml
Unauthenticated = 5 # Default
```
Unauthenticated defines the threshold to which authenticated requests get limited. More than this many unauthenticated requests per `UnAuthenticatedRateLimitPeriod` will be rejected.

### UnauthenticatedPeriod
```toml
UnauthenticatedPeriod = '20s' # Default
```
UnauthenticatedPeriod defines the period to which unauthenticated requests get limited.

## WebServer.MFA
```toml
[WebServer.MFA]
RPID = 'localhost' # Example
RPOrigin = 'http://localhost:6688/' # Example
```
The Operator UI frontend supports enabling Multi Factor Authentication via Webauthn per account. When enabled, logging in will require the account password and a hardware or OS security key such as Yubikey. To enroll, log in to the operator UI and click the circle purple profile button at the top right and then click **Register MFA Token**. Tap your hardware security key or use the OS public key management feature to enroll a key. Next time you log in, this key will be required to authenticate.

### RPID
```toml
RPID = 'localhost' # Example
```
RPID is the FQDN of where the Operator UI is served. When serving locally, the value should be `localhost`.

### RPOrigin
```toml
RPOrigin = 'http://localhost:6688/' # Example
```
RPOrigin is the origin URL where WebAuthn requests initiate, including scheme and port. When serving locally, the value should be `http://localhost:6688/`.

## WebServer.TLS
```toml
[WebServer.TLS]
CertPath = '~/.cl/certs' # Example
Host = 'tls-host' # Example
KeyPath = '/home/$USER/.chainlink/tls/server.key' # Example
HTTPSPort = 6689 # Default
ForceRedirect = false # Default
```
The TLS settings apply only if you want to enable TLS security on your Chainlink node.

### CertPath
```toml
CertPath = '~/.cl/certs' # Example
```
CertPath is the location of the TLS certificate file.

### Host
```toml
Host = 'tls-host' # Example
```
Host is the hostname configured for TLS to be used by the Chainlink node. This is useful if you configured a domain name specific for your Chainlink node.

### KeyPath
```toml
KeyPath = '/home/$USER/.chainlink/tls/server.key' # Example
```
KeyPath is the location of the TLS private key file.

### HTTPSPort
```toml
HTTPSPort = 6689 # Default
```
HTTPSPort is the port used for HTTPS connections. Set this to `0` to disable HTTPS. Disabling HTTPS also relieves Chainlink nodes of the requirement for a TLS certificate.

### ForceRedirect
```toml
ForceRedirect = false # Default
```
ForceRedirect forces TLS redirect for unencrypted connections.

## JobPipeline
```toml
[JobPipeline]
ExternalInitiatorsEnabled = false # Default
MaxRunDuration = '10m' # Default
MaxSuccessfulRuns = 10000 # Default
ReaperInterval = '1h' # Default
ReaperThreshold = '24h' # Default
ResultWriteQueueDepth = 100 # Default
```


### ExternalInitiatorsEnabled
```toml
ExternalInitiatorsEnabled = false # Default
```
ExternalInitiatorsEnabled enables the External Initiator feature. If disabled, `webhook` jobs can ONLY be initiated by a logged-in user. If enabled, `webhook` jobs can be initiated by a whitelisted external initiator.

### MaxRunDuration
```toml
MaxRunDuration = '10m' # Default
```
MaxRunDuration is the maximum time allowed for a single job run. If it takes longer, it will exit early and be marked errored. If set to zero, disables the time limit completely.

### MaxSuccessfulRuns
```toml
MaxSuccessfulRuns = 10000 # Default
```
MaxSuccessfulRuns caps the number of completed successful runs per pipeline
spec in the database. You can set it to zero as a performance optimisation;
this will avoid saving any successful run.

Note this is not a hard cap, it can drift slightly larger than this but not
by more than 5% or so.

### ReaperInterval
```toml
ReaperInterval = '1h' # Default
```
ReaperInterval controls how often the job pipeline reaper will run to delete completed jobs older than ReaperThreshold, in order to keep database size manageable.

Set to `0` to disable the periodic reaper.

### ReaperThreshold
```toml
ReaperThreshold = '24h' # Default
```
ReaperThreshold determines the age limit for job runs. Completed job runs older than this will be automatically purged from the database.

### ResultWriteQueueDepth
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
ResultWriteQueueDepth = 100 # Default
```
ResultWriteQueueDepth controls how many writes will be buffered before subsequent writes are dropped, for jobs that write results asynchronously for performance reasons, such as OCR.

## JobPipeline.HTTPRequest
```toml
[JobPipeline.HTTPRequest]
DefaultTimeout = '15s' # Default
MaxSize = '32768' # Default
```


### DefaultTimeout
```toml
DefaultTimeout = '15s' # Default
```
DefaultTimeout defines the default timeout for HTTP requests made by `http` and `bridge` adapters.

### MaxSize
```toml
MaxSize = '32768' # Default
```
MaxSize defines the maximum size for HTTP requests and responses made by `http` and `bridge` adapters.

## FluxMonitor
```toml
[FluxMonitor]
DefaultTransactionQueueDepth = 1 # Default
SimulateTransactions = false # Default
```


### DefaultTransactionQueueDepth
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DefaultTransactionQueueDepth = 1 # Default
```
DefaultTransactionQueueDepth controls the queue size for `DropOldestStrategy` in Flux Monitor. Set to 0 to use `SendEvery` strategy instead.

### SimulateTransactions
```toml
SimulateTransactions = false # Default
```
SimulateTransactions enables transaction simulation for Flux Monitor.

## OCR2
```toml
[OCR2]
Enabled = false # Default
ContractConfirmations = 3 # Default
BlockchainTimeout = '20s' # Default
ContractPollInterval = '1m' # Default
ContractSubscribeInterval = '2m' # Default
ContractTransmitterTransmitTimeout = '10s' # Default
DatabaseTimeout = '10s' # Default
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000' # Example
CaptureEATelemetry = false # Default
```


### Enabled
```toml
Enabled = false # Default
```
Enabled enables OCR2 jobs.

### ContractConfirmations
```toml
ContractConfirmations = 3 # Default
```
ContractConfirmations is the number of block confirmations to wait for before enacting an on-chain
configuration change. This value doesn't need to be very high (in
particular, it does not need to protect against malicious re-orgs).
Since configuration changes create some overhead, and mini-reorgs
are fairly common, recommended values are between two and ten.

Malicious re-orgs are not any more of concern here than they are in
blockchain applications in general: Since nodes check the contract for the
latest config every ContractConfigTrackerPollInterval.Seconds(), they will
come to a common view of the current config within any interval longer than
that, as long as the latest setConfig transaction in the longest chain is
stable. They will thus be able to continue reporting after the poll
interval, unless an adversary is able to repeatedly re-org the transaction
out during every poll interval, which would amount to the capability to
censor any transaction.

Note that 1 confirmation implies that the transaction/event has been mined in one block.
0 confirmations would imply that the event would be recognised before it has even been mined, which is not currently supported.
e.g.
Current block height: 42
Changed in block height: 43
Contract config confirmations: 1
STILL PENDING

Current block height: 43
Changed in block height: 43
Contract config confirmations: 1
CONFIRMED

### BlockchainTimeout
```toml
BlockchainTimeout = '20s' # Default
```
BlockchainTimeout is the timeout for blockchain queries (mediated through
ContractConfigTracker and ContractTransmitter).
(This is necessary because an oracle's operations are serialized, so
blocking forever on a chain interaction would break the oracle.)

### ContractPollInterval
```toml
ContractPollInterval = '1m' # Default
```
ContractPollInterval is the polling interval at which ContractConfigTracker is queried for# updated on-chain configurations. Recommended values are between
fifteen seconds and two minutes.

### ContractSubscribeInterval
```toml
ContractSubscribeInterval = '2m' # Default
```
ContractSubscribeInterval is the interval at which we try to establish a subscription on ContractConfigTracker
if one doesn't exist. Recommended values are between two and five minutes.

### ContractTransmitterTransmitTimeout
```toml
ContractTransmitterTransmitTimeout = '10s' # Default
```
ContractTransmitterTransmitTimeout is the timeout for ContractTransmitter.Transmit calls.

### DatabaseTimeout
```toml
DatabaseTimeout = '10s' # Default
```
DatabaseTimeout is the timeout for database interactions.
(This is necessary because an oracle's operations are serialized, so
blocking forever on an observation would break the oracle.)

### KeyBundleID
```toml
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000' # Example
```
KeyBundleID is a sha256 hexadecimal hash identifier.

### CaptureEATelemetry
```toml
CaptureEATelemetry = false # Default
```
CaptureEATelemetry toggles collecting extra information from External Adaptares

## OCR
```toml
[OCR]
Enabled = false # Default
ObservationTimeout = '5s' # Default
BlockchainTimeout = '20s' # Default
ContractPollInterval = '1m' # Default
ContractSubscribeInterval = '2m' # Default
DefaultTransactionQueueDepth = 1 # Default
KeyBundleID = 'acdd42797a8b921b2910497badc5000600000000000000000000000000000000' # Example
SimulateTransactions = false # Default
TransmitterAddress = '0xa0788FC17B1dEe36f057c42B6F373A34B014687e' # Example
CaptureEATelemetry = false # Default
```
This section applies only if you are running off-chain reporting jobs.

### Enabled
```toml
Enabled = false # Default
```
Enabled enables OCR jobs.

### ObservationTimeout
```toml
ObservationTimeout = '5s' # Default
```
ObservationTimeout is the timeout for making observations using the DataSource.Observe method.
(This is necessary because an oracle's operations are serialized, so
blocking forever on an observation would break the oracle.)

### BlockchainTimeout
```toml
BlockchainTimeout = '20s' # Default
```
BlockchainTimeout is the timeout for blockchain queries (mediated through
ContractConfigTracker and ContractTransmitter).
(This is necessary because an oracle's operations are serialized, so
blocking forever on a chain interaction would break the oracle.)

### ContractPollInterval
```toml
ContractPollInterval = '1m' # Default
```
ContractPollInterval is the polling interval at which ContractConfigTracker is queried for
updated on-chain configurations. Recommended values are between
fifteen seconds and two minutes.

### ContractSubscribeInterval
```toml
ContractSubscribeInterval = '2m' # Default
```
ContractSubscribeInterval is the interval at which we try to establish a subscription on ContractConfigTracker
if one doesn't exist. Recommended values are between two and five minutes.

### DefaultTransactionQueueDepth
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DefaultTransactionQueueDepth = 1 # Default
```
DefaultTransactionQueueDepth controls the queue size for `DropOldestStrategy` in OCR. Set to 0 to use `SendEvery` strategy instead.

### KeyBundleID
```toml
KeyBundleID = 'acdd42797a8b921b2910497badc5000600000000000000000000000000000000' # Example
```
KeyBundleID is the default key bundle ID to use for OCR jobs. If you have an OCR job that does not explicitly specify a key bundle ID, it will fall back to this value.

### SimulateTransactions
```toml
SimulateTransactions = false # Default
```
SimulateTransactions enables transaction simulation for OCR.

### TransmitterAddress
```toml
TransmitterAddress = '0xa0788FC17B1dEe36f057c42B6F373A34B014687e' # Example
```
TransmitterAddress is the default sending address to use for OCR. If you have an OCR job that does not explicitly specify a transmitter address, it will fall back to this value.

### CaptureEATelemetry
```toml
CaptureEATelemetry = false # Default
```
CaptureEATelemetry toggles collecting extra information from External Adaptares

## P2P
```toml
[P2P]
IncomingMessageBufferSize = 10 # Default
OutgoingMessageBufferSize = 10 # Default
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw' # Example
TraceLogging = false # Default
```
P2P supports multiple networking stack versions. You may configure `[P2P.V1]`, `[P2P.V2]`, or both to run simultaneously.
If both are configured, then for each link with another peer, V2 networking will be preferred. If V2 does not work, the link will
automatically fall back to V1. If V2 starts working again later, it will automatically be preferred again. This is useful
for migrating networks without downtime. Note that the two networking stacks _must not_ be configured to bind to the same IP/port.

All nodes in the OCR network should share the same networking stack.

### IncomingMessageBufferSize
```toml
IncomingMessageBufferSize = 10 # Default
```
IncomingMessageBufferSize is the per-remote number of incoming
messages to buffer. Any additional messages received on top of those
already in the queue will be dropped.

### OutgoingMessageBufferSize
```toml
OutgoingMessageBufferSize = 10 # Default
```
OutgoingMessageBufferSize is the per-remote number of outgoing
messages to buffer. Any additional messages send on top of those
already in the queue will displace the oldest.
NOTE: OutgoingMessageBufferSize should be comfortably smaller than remote's
IncomingMessageBufferSize to give the remote enough space to process
them all in case we regained connection and now send a bunch at once

### PeerID
```toml
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw' # Example
```
PeerID is the default peer ID to use for OCR jobs. If unspecified, uses the first available peer ID.

### TraceLogging
```toml
TraceLogging = false # Default
```
TraceLogging enables trace level logging.

## P2P.V1
```toml
[P2P.V1]
Enabled = true # Default
AnnounceIP = '1.2.3.4' # Example
AnnouncePort = 1337 # Example
BootstrapCheckInterval = '20s' # Default
DefaultBootstrapPeers = ['/dns4/example.com/tcp/1337/p2p/12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U', '/ip4/1.2.3.4/tcp/9999/p2p/12D3KooWLZ9uTC3MrvKfDpGju6RAQubiMDL7CuJcAgDRTYP7fh7R'] # Example
DHTAnnouncementCounterUserPrefix = 0 # Default
DHTLookupInterval = 10 # Default
ListenIP = '0.0.0.0' # Default
ListenPort = 1337 # Example
NewStreamTimeout = '10s' # Default
PeerstoreWriteInterval = '5m' # Default
```


### Enabled
```toml
Enabled = true # Default
```
Enabled enables P2P V1.

### AnnounceIP
```toml
AnnounceIP = '1.2.3.4' # Example
```
AnnounceIP should be set as the externally reachable IP address of the Chainlink node.

### AnnouncePort
```toml
AnnouncePort = 1337 # Example
```
AnnouncePort should be set as the externally reachable port of the Chainlink node.

### BootstrapCheckInterval
```toml
BootstrapCheckInterval = '20s' # Default
```
BootstrapCheckInterval is the interval at which nodes check connections to bootstrap nodes and reconnect if any of them is lost.
Setting this to a small value would allow newly joined bootstrap nodes to get more connectivity
more quickly, which helps to make bootstrap process faster. The cost of this operation is relatively
cheap. We set this to 1 minute during our test.

### DefaultBootstrapPeers
```toml
DefaultBootstrapPeers = ['/dns4/example.com/tcp/1337/p2p/12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U', '/ip4/1.2.3.4/tcp/9999/p2p/12D3KooWLZ9uTC3MrvKfDpGju6RAQubiMDL7CuJcAgDRTYP7fh7R'] # Example
```
DefaultBootstrapPeers is the default set of bootstrap peers.

### DHTAnnouncementCounterUserPrefix
```toml
DHTAnnouncementCounterUserPrefix = 0 # Default
```
DHTAnnouncementCounterUserPrefix can be used to restore the node's
ability to announce its IP/port on the P2P network after a database
rollback. Make sure to only increase this value, and *never* decrease it.
Don't use this variable unless you really know what you're doing, since you
could semi-permanently exclude your node from the P2P network by
misconfiguring it.

### DHTLookupInterval
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DHTLookupInterval = 10 # Default
```
DHTLookupInterval is the interval between which we do the expensive peer
lookup using DHT.

Every DHTLookupInterval failures to open a stream to a peer, we will
attempt to lookup its IP from DHT

### ListenIP
```toml
ListenIP = '0.0.0.0' # Default
```
ListenIP is the default IP address to bind to.

### ListenPort
```toml
ListenPort = 1337 # Example
```
ListenPort is the port to listen on. If left blank, the node randomly selects a different port each time it boots. It is highly recommended to set this to a static value to avoid network instability.

### NewStreamTimeout
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
NewStreamTimeout = '10s' # Default
```
NewStreamTimeout is the maximum length of time to wait to open a
stream before we give up.
We shouldn't hit this in practice since libp2p will give up fast if
it can't get a connection, but it is here anyway as a failsafe.
Set to 0 to disable any timeout on top of what libp2p gives us by default.

### PeerstoreWriteInterval
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
PeerstoreWriteInterval = '5m' # Default
```
PeerstoreWriteInterval controls how often the peerstore for the OCR V1 networking stack is persisted to the database.

## P2P.V2
```toml
[P2P.V2]
Enabled = false # Default
AnnounceAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
DefaultBootstrappers = ['12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U@1.2.3.4:9999', '12D3KooWM55u5Swtpw9r8aFLQHEtw7HR4t44GdNs654ej5gRs2Dh@example.com:1234'] # Example
DeltaDial = '15s' # Default
DeltaReconcile = '1m' # Default
ListenAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
```


### Enabled
```toml
Enabled = false # Default
```
Enabled enables P2P V2.
Note: V1.Enabled is true by default, so it must be set false in order to run V2 only.

### AnnounceAddresses
```toml
AnnounceAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
```
AnnounceAddresses is the addresses the peer will advertise on the network in `host:port` form as accepted by the TCP version of Go’s `net.Dial`.
The addresses should be reachable by other nodes on the network. When attempting to connect to another node,
a node will attempt to dial all of the other node’s AnnounceAddresses in round-robin fashion.

### DefaultBootstrappers
```toml
DefaultBootstrappers = ['12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U@1.2.3.4:9999', '12D3KooWM55u5Swtpw9r8aFLQHEtw7HR4t44GdNs654ej5gRs2Dh@example.com:1234'] # Example
```
DefaultBootstrappers is the default bootstrapper peers for libocr's v2 networking stack.

Oracle nodes typically only know each other’s PeerIDs, but not their hostnames, IP addresses, or ports.
DefaultBootstrappers are special nodes that help other nodes discover each other’s `AnnounceAddresses` so they can communicate.
Nodes continuously attempt to connect to bootstrappers configured in here. When a node wants to connect to another node
(which it knows only by PeerID, but not by address), it discovers the other node’s AnnounceAddresses from communications
received from its DefaultBootstrappers or other discovered nodes. To facilitate discovery,
nodes will regularly broadcast signed announcements containing their PeerID and AnnounceAddresses.

### DeltaDial
```toml
DeltaDial = '15s' # Default
```
DeltaDial controls how far apart Dial attempts are

### DeltaReconcile
```toml
DeltaReconcile = '1m' # Default
```
DeltaReconcile controls how often a Reconcile message is sent to every peer.

### ListenAddresses
```toml
ListenAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
```
ListenAddresses is the addresses the peer will listen to on the network in `host:port` form as accepted by `net.Listen()`,
but the host and port must be fully specified and cannot be empty. You can specify `0.0.0.0` (IPv4) or `::` (IPv6) to listen on all interfaces, but that is not recommended.

## Keeper
```toml
[Keeper]
DefaultTransactionQueueDepth = 1 # Default
GasPriceBufferPercent = 20 # Default
GasTipCapBufferPercent = 20 # Default
BaseFeeBufferPercent = 20 # Default
MaxGracePeriod = 100 # Default
TurnLookBack = 1_000 # Default
```


### DefaultTransactionQueueDepth
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DefaultTransactionQueueDepth = 1 # Default
```
DefaultTransactionQueueDepth controls the queue size for `DropOldestStrategy` in Keeper. Set to 0 to use `SendEvery` strategy instead.

### GasPriceBufferPercent
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
GasPriceBufferPercent = 20 # Default
```
GasPriceBufferPercent specifies the percentage to add to the gas price used for checking whether to perform an upkeep. Only applies in legacy mode (EIP-1559 off).

### GasTipCapBufferPercent
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
GasTipCapBufferPercent = 20 # Default
```
GasTipCapBufferPercent specifies the percentage to add to the gas price used for checking whether to perform an upkeep. Only applies in EIP-1559 mode.

### BaseFeeBufferPercent
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
BaseFeeBufferPercent = 20 # Default
```
BaseFeeBufferPercent specifies the percentage to add to the base fee used for checking whether to perform an upkeep. Applies only in EIP-1559 mode.

### MaxGracePeriod
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
MaxGracePeriod = 100 # Default
```
MaxGracePeriod is the maximum number of blocks that a keeper will wait after performing an upkeep before it resumes checking that upkeep

### TurnLookBack
```toml
TurnLookBack = 1_000 # Default
```
TurnLookBack is the number of blocks in the past to look back when getting a block for a turn.

## Keeper.Registry
```toml
[Keeper.Registry]
CheckGasOverhead = 200_000 # Default
PerformGasOverhead = 300_000 # Default
SyncInterval = '30m' # Default
MaxPerformDataSize = 5_000 # Default
SyncUpkeepQueueSize = 10 # Default
```


### CheckGasOverhead
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
CheckGasOverhead = 200_000 # Default
```
CheckGasOverhead is the amount of extra gas to provide checkUpkeep() calls to account for the gas consumed by the keeper registry.

### PerformGasOverhead
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
PerformGasOverhead = 300_000 # Default
```
PerformGasOverhead is the amount of extra gas to provide performUpkeep() calls to account for the gas consumed by the keeper registry

### SyncInterval
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
SyncInterval = '30m' # Default
```
SyncInterval is the interval in which the RegistrySynchronizer performs a full sync of the keeper registry contract it is tracking.

### MaxPerformDataSize
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
MaxPerformDataSize = 5_000 # Default
```
MaxPerformDataSize is the max size of perform data.

### SyncUpkeepQueueSize
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
SyncUpkeepQueueSize = 10 # Default
```
SyncUpkeepQueueSize represents the maximum number of upkeeps that can be synced in parallel.

## AutoPprof
```toml
[AutoPprof]
Enabled = false # Default
ProfileRoot = 'prof/root' # Example
PollInterval = '10s' # Default
GatherDuration = '10s' # Default
GatherTraceDuration = '5s' # Default
MaxProfileSize = '100mb' # Default
CPUProfileRate = 1 # Default
MemProfileRate = 1 # Default
BlockProfileRate = 1 # Default
MutexProfileFraction = 1 # Default
MemThreshold = '4gb' # Default
GoroutineThreshold = 5000 # Default
```
The Chainlink node is equipped with an internal "nurse" service that can perform automatic `pprof` profiling when the certain resource thresholds are exceeded, such as memory and goroutine count. These profiles are saved to disk to facilitate fine-grained debugging of performance-related issues. In general, if you notice that your node has begun to accumulate profiles, forward them to the Chainlink team.

To learn more about these profiles, read the [Profiling Go programs with pprof](https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/) guide.

### Enabled
```toml
Enabled = false # Default
```
Enabled enables the automatic profiling service.

### ProfileRoot
```toml
ProfileRoot = 'prof/root' # Example
```
ProfileRoot sets the location on disk where pprof profiles will be stored. Defaults to `RootDir`.

### PollInterval
```toml
PollInterval = '10s' # Default
```
PollInterval is the interval at which the node's resources are checked.

### GatherDuration
```toml
GatherDuration = '10s' # Default
```
GatherDuration is the duration for which profiles are gathered when profiling starts.

### GatherTraceDuration
```toml
GatherTraceDuration = '5s' # Default
```
GatherTraceDuration is the duration for which traces are gathered when profiling is kicked off. This is separately configurable because traces are significantly larger than other types of profiles.

### MaxProfileSize
```toml
MaxProfileSize = '100mb' # Default
```
MaxProfileSize is the maximum amount of disk space that profiles may consume before profiling is disabled.

### CPUProfileRate
```toml
CPUProfileRate = 1 # Default
```
CPUProfileRate sets the rate for CPU profiling. See https://pkg.go.dev/runtime#SetCPUProfileRate.

### MemProfileRate
```toml
MemProfileRate = 1 # Default
```
MemProfileRate sets the rate for memory profiling. See https://pkg.go.dev/runtime#pkg-variables.

### BlockProfileRate
```toml
BlockProfileRate = 1 # Default
```
BlockProfileRate sets the fraction of blocking events for goroutine profiling. See https://pkg.go.dev/runtime#SetBlockProfileRate.

### MutexProfileFraction
```toml
MutexProfileFraction = 1 # Default
```
MutexProfileFraction sets the fraction of contention events for mutex profiling. See https://pkg.go.dev/runtime#SetMutexProfileFraction.

### MemThreshold
```toml
MemThreshold = '4gb' # Default
```
MemThreshold sets the maximum amount of memory the node can actively consume before profiling begins.

### GoroutineThreshold
```toml
GoroutineThreshold = 5000 # Default
```
GoroutineThreshold is the maximum number of actively-running goroutines the node can spawn before profiling begins.

## Pyroscope
```toml
[Pyroscope]
ServerAddress = 'http://localhost:4040' # Example
Environment = 'mainnet' # Default
```


### ServerAddress
```toml
ServerAddress = 'http://localhost:4040' # Example
```
ServerAddress sets the address that will receive the profile logs. It enables the profiling service.

### Environment
```toml
Environment = 'mainnet' # Default
```
Environment sets the target environment tag in which profiles will be added to.

## Sentry
```toml
[Sentry]
Debug = false # Default
DSN = 'sentry-dsn' # Example
Environment = 'my-custom-env' # Example
Release = 'v1.2.3' # Example
```


### Debug
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
Debug = false # Default
```
Debug enables printing of Sentry SDK debug messages.

### DSN
```toml
DSN = 'sentry-dsn' # Example
```
DSN is the data source name where events will be sent. Sentry is completely disabled if this is left blank.

### Environment
```toml
Environment = 'my-custom-env' # Example
```
Environment overrides the Sentry environment to the given value. Otherwise autodetects between dev/prod.

### Release
```toml
Release = 'v1.2.3' # Example
```
Release overrides the Sentry release to the given value. Otherwise uses the compiled-in version number.

## Insecure
```toml
[Insecure]
DevWebServer = false # Default
OCRDevelopmentMode = false # Default
InfiniteDepthQueries = false # Default
DisableRateLimiting = false # Default
```
Insecure config family is only allowed in development builds.

### DevWebServer
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DevWebServer = false # Default
```
DevWebServer skips secure configuration for webserver AllowedHosts, SSL, etc.

### OCRDevelopmentMode
```toml
OCRDevelopmentMode = false # Default
```
OCRDevelopmentMode run OCR in development mode.

### InfiniteDepthQueries
```toml
InfiniteDepthQueries = false # Default
```
InfiniteDepthQueries skips graphql query depth limit checks.

### DisableRateLimiting
```toml
DisableRateLimiting = false # Default
```
DisableRateLimiting skips ratelimiting on asset requests.

## EVM
EVM defaults depend on ChainID:

<details><summary>Ethereum Mainnet (1)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x514910771AF9Ca656af840dff83E8264EcF986CA'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
OperatorFactoryAddress = '0x3E64Cd889482443324F91bFA9c84fE72A511f48A'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Ethereum Ropsten (3)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x20fE562d797A42Dcb3399062AE9546cd06f63280'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Ethereum Rinkeby (4)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x01BE23585060835E02B77ef475b0Cc51aA1e0709'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Ethereum Goerli (5)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x326C977E6efc84E512bB9C30f76E30c160eD06FB'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Optimism Mainnet (10)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimism'
FinalityDepth = 1
LinkContractAddress = '0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '15s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 10
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000
```

</p></details>

<details><summary>RSK Mainnet (30)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x14AdaE34beF7ca957Ce2dDe5ADD97ea050123827'
LogBackfillBatchSize = 100
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '50 mwei'
PriceMax = '50 gwei'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 mwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>RSK Testnet (31)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x8bBbd80981FE76d44854D8DF305e8985c19f0e78'
LogBackfillBatchSize = 100
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '50 mwei'
PriceMax = '50 gwei'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 mwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Ethereum Kovan (42)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0xa36085F69e2889c224210F603D836748e7dC0088'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
OperatorFactoryAddress = '0x8007e24251b1D2Fc518Eb843A701d9cD21fe0aA3'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>BSC Mainnet (56)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x404460C6A5EdE2D891e8297795264fDe62ADBB75'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 2

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '2s'
DatabaseTimeout = '2s'
ObservationGracePeriod = '500ms'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>OKX Testnet (65)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>OKX Mainnet (66)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Optimism Kovan (69)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimism'
FinalityDepth = 1
LinkContractAddress = '0x4911b761993b9c8c0d14Ba2d86902AF6B0074F5B'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '15s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 10
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000
```

</p></details>

<details><summary>xDai Mainnet (100)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'xdai'
FinalityDepth = 50
LinkContractAddress = '0xE2e73A1c69ecF83F464EFCE6A5be353a37cA09b2'
LogBackfillBatchSize = 100
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '1 gwei'
PriceMax = '500 gwei'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Heco Mainnet (128)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x404460C6A5EdE2D891e8297795264fDe62ADBB75'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 2

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '2s'
DatabaseTimeout = '2s'
ObservationGracePeriod = '500ms'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Polygon Mainnet (137)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 500
LinkContractAddress = '0xb0897686c545045aFc77CF20eC7A532E3120E0F1'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 5000
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '30 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '30 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '20 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Fantom Mainnet (250)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x6F43FF82CCA38001B6699a8AC47A2d0E66939407'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 2

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '15 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 3800000
```

</p></details>

<details><summary>Optimism Goerli (420)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
LinkContractAddress = '0xdc2CC710e42857672E7907CF474a69B63B93089f'
LogBackfillBatchSize = 100
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 wei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000
```

</p></details>

<details><summary>Metis Rinkeby (588)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'metis'
FinalityDepth = 1
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Klaytn Testnet (1001)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 1
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '750 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Metis Mainnet (1088)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'metis'
FinalityDepth = 1
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Simulated (1337)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 1
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '100'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '0s'
ResendAfterThreshold = '0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FixedPrice'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 10
MaxBufferSize = 100
SamplingInterval = '0s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Fantom Testnet (4002)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0xfaFedb041c0DD4fA2Dc0d87a6B0979Ee6FA7af5F'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 2

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '15 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 3800000
```

</p></details>

<details><summary>Klaytn Mainnet (8217)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 1
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '750 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Arbitrum Mainnet (42161)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
LinkContractAddress = '0xf97f4df75117a78c1A5a0DBb814Af92458539FB4'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'Arbitrum'
PriceDefault = '100 mwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 1000000000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Avalanche Fuji (43113)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 1
LinkContractAddress = '0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 2

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '25 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '25 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Avalanche Mainnet (43114)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 1
LinkContractAddress = '0x5947BB275c521040051D82396192181b413227A3'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 2

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '25 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '25 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Polygon Mumbai (80001)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 500
LinkContractAddress = '0x326C977E6efc84E512bB9C30f76E30c160eD06FB'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 5000
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '1 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '20 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Arbitrum Rinkeby (421611)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
LinkContractAddress = '0x615fBe6372676474d9e6933d310469c9b68e9726'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'Arbitrum'
PriceDefault = '100 mwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 1000000000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Arbitrum Goerli (421613)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
LinkContractAddress = '0xd14838A68E8AFBAdE5efb411d5871ea0011AFd28'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'Arbitrum'
PriceDefault = '100 mwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 1000000000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Ethereum Sepolia (11155111)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0xb227f007804c16546Bd054dfED2E7A1fD5437678'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Harmony Mainnet (1666600000)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x218532a12a389a4a92fC0C5Fb22901D1c19198aA'
LogBackfillBatchSize = 100
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>

<details><summary>Harmony Testnet (1666700000)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
LinkContractAddress = '0x8b12Ac23BFe11cAb03a634C1F117D64a7f2cFD3e'
LogBackfillBatchSize = 100
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1

[Transactions]
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 4
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5300000
```

</p></details>


### ChainID
```toml
ChainID = '1' # Example
```
ChainID is the EVM chain ID. Mandatory.

### Enabled
```toml
Enabled = true # Default
```
Enabled enables this chain.

### AutoCreateKey
```toml
AutoCreateKey = true # Default
```
AutoCreateKey, if set to true, will ensure that there is always at least one transmit key for the given chain.

### BlockBackfillDepth
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
BlockBackfillDepth = 10 # Default
```
BlockBackfillDepth specifies the number of blocks before the current HEAD that the log broadcaster will try to re-consume logs from.

### BlockBackfillSkip
```toml
BlockBackfillSkip = false # Default
```
BlockBackfillSkip enables skipping of very long backfills.

### ChainType
```toml
ChainType = 'Optimism' # Example
```
ChainType is automatically detected from chain ID. Set this to force a certain chain type regardless of chain ID.

### FinalityDepth
```toml
FinalityDepth = 50 # Default
```
FinalityDepth is the number of blocks after which an ethereum transaction is considered "final". Note that the default is automatically set based on chain ID so it should not be necessary to change this under normal operation.
BlocksConsideredFinal determines how deeply we look back to ensure that transactions are confirmed onto the longest chain
There is not a large performance penalty to setting this relatively high (on the order of hundreds)
It is practically limited by the number of heads we store in the database and should be less than this with a comfortable margin.
If a transaction is mined in a block more than this many blocks ago, and is reorged out, we will NOT retransmit this transaction and undefined behaviour can occur including gaps in the nonce sequence that require manual intervention to fix.
Therefore this number represents a number of blocks we consider large enough that no re-org this deep will ever feasibly happen.

Special cases:
`FinalityDepth`=0 would imply that transactions can be final even before they were mined into a block. This is not supported.
`FinalityDepth`=1 implies that transactions are final after we see them in one block.

Examples:

Transaction sending:
A transaction is sent at block height 42

`FinalityDepth` is set to 5
A re-org occurs at height 44 starting at block 41, transaction is marked for rebroadcast
A re-org occurs at height 46 starting at block 41, transaction is marked for rebroadcast
A re-org occurs at height 47 starting at block 41, transaction is NOT marked for rebroadcast

### FlagsContractAddress
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3' # Example
```
FlagsContractAddress can optionally point to a [Flags contract](../contracts/src/v0.8/Flags.sol). If set, the node will lookup that contract for each job that supports flags contracts (currently OCR and FM jobs are supported). If the job's contractAddress is set as hibernating in the FlagsContractAddress address, it overrides the standard update parameters (such as heartbeat/threshold).

### LinkContractAddress
```toml
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e' # Example
```
LinkContractAddress is the canonical ERC-677 LINK token contract address on the given chain. Note that this is usually autodetected from chain ID.

### LogBackfillBatchSize
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
LogBackfillBatchSize = 100 # Default
```
LogBackfillBatchSize sets the batch size for calling FilterLogs when we backfill missing logs.

### LogPollInterval
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
LogPollInterval = '15s' # Default
```
LogPollInterval works in conjunction with Feature.LogPoller. Controls how frequently the log poller polls for logs. Defaults to the block production rate.

### LogKeepBlocksDepth
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
LogKeepBlocksDepth = 100000 # Default
```
LogKeepBlocksDepth works in conjunction with Feature.LogPoller. Controls how many blocks the poller will keep, must be greater than FinalityDepth+1.

### MinContractPayment
```toml
MinContractPayment = '10000000000000 juels' # Default
```
MinContractPayment is the minimum payment in LINK required to execute a direct request job. This can be overridden on a per-job basis.

### MinIncomingConfirmations
```toml
MinIncomingConfirmations = 3 # Default
```
MinIncomingConfirmations is the minimum required confirmations before a log event will be consumed.

### NonceAutoSync
```toml
NonceAutoSync = true # Default
```
NonceAutoSync enables automatic nonce syncing on startup. Chainlink nodes will automatically try to sync its local nonce with the remote chain on startup and fast forward if necessary. This is almost always safe but can be disabled in exceptional cases by setting this value to false.

### NoNewHeadsThreshold
```toml
NoNewHeadsThreshold = '3m' # Default
```
NoNewHeadsThreshold controls how long to wait after receiving no new heads before `NodePool` marks rpc endpoints as
out-of-sync, and `HeadTracker` logs warnings.

Set to zero to disable out-of-sync checking.

### OperatorFactoryAddress
```toml
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418' # Example
```
OperatorFactoryAddress is the address of the canonical operator forwarder contract on the given chain. Note that this is usually autodetected from chain ID.

### RPCDefaultBatchSize
```toml
RPCDefaultBatchSize = 100 # Default
```
RPCDefaultBatchSize is the default batch size for batched RPC calls.

### RPCBlockQueryDelay
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
RPCBlockQueryDelay = 1 # Default
```
RPCBlockQueryDelay controls the number of blocks to trail behind head in the block history estimator and balance monitor.
For example, if this is set to 3, and we receive block 10, block history estimator will fetch block 7.

CAUTION: You might be tempted to set this to 0 to use the latest possible
block, but it is possible to receive a head BEFORE that block is actually
available from the connected node via RPC, due to race conditions in the code of the remote ETH node. In this case you will get false
"zero" blocks that are missing transactions.

## EVM.Transactions
```toml
[EVM.Transactions]
ForwardersEnabled = false # Default
MaxInFlight = 16 # Default
MaxQueued = 250 # Default
ReaperInterval = '1h' # Default
ReaperThreshold = '168h' # Default
ResendAfterThreshold = '1m' # Default
```


### ForwardersEnabled
```toml
ForwardersEnabled = false # Default
```
ForwardersEnabled enables or disables sending transactions through forwarder contracts.

### MaxInFlight
```toml
MaxInFlight = 16 # Default
```
MaxInFlight controls how many transactions are allowed to be "in-flight" i.e. broadcast but unconfirmed at any one time. You can consider this a form of transaction throttling.

The default is set conservatively at 16 because this is a pessimistic minimum that both geth and parity will hold without evicting local transactions. If your node is falling behind and you need higher throughput, you can increase this setting, but you MUST make sure that your ETH node is configured properly otherwise you can get nonce gapped and your node will get stuck.

0 value disables the limit. Use with caution.

### MaxQueued
```toml
MaxQueued = 250 # Default
```
MaxQueued is the maximum number of unbroadcast transactions per key that are allowed to be enqueued before jobs will start failing and rejecting send of any further transactions. This represents a sanity limit and generally indicates a problem with your ETH node (transactions are not getting mined).

Do NOT blindly increase this value thinking it will fix things if you start hitting this limit because transactions are not getting mined, you will instead only make things worse.

In deployments with very high burst rates, or on chains with large re-orgs, you _may_ consider increasing this.

0 value disables any limit on queue size. Use with caution.

### ReaperInterval
```toml
ReaperInterval = '1h' # Default
```
ReaperInterval controls how often the EthTx reaper will run.

### ReaperThreshold
```toml
ReaperThreshold = '168h' # Default
```
ReaperThreshold indicates how old an EthTx ought to be before it can be reaped.

### ResendAfterThreshold
```toml
ResendAfterThreshold = '1m' # Default
```
ResendAfterThreshold controls how long to wait before re-broadcasting a transaction that has not yet been confirmed.

## EVM.BalanceMonitor
```toml
[EVM.BalanceMonitor]
Enabled = true # Default
```


### Enabled
```toml
Enabled = true # Default
```
Enabled balance monitoring for all keys.

## EVM.GasEstimator
```toml
[EVM.GasEstimator]
Mode = 'BlockHistory' # Default
PriceDefault = '20 gwei' # Default
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether' # Default
PriceMin = '1 gwei' # Default
LimitDefault = 500_000 # Default
LimitMax = 500_000 # Default
LimitMultiplier = '1.0' # Default
LimitTransfer = 21_000 # Default
BumpMin = '5 gwei' # Default
BumpPercent = 20 # Default
BumpThreshold = 3 # Default
BumpTxDepth = 10 # Default
EIP1559DynamicFees = false # Default
FeeCapDefault = '100 gwei' # Default
TipCapDefault = '1 wei' # Default
TipCapMin = '1 wei' # Default
```


### Mode
```toml
Mode = 'BlockHistory' # Default
```
Mode controls what type of gas estimator is used.

- `FixedPrice` uses static configured values for gas price (can be set via API call).
- `BlockHistory` dynamically adjusts default gas price based on heuristics from mined blocks.
- `Optimism2`/`L2Suggested` is a special mode only for use with Optimism and Metis blockchains. This mode will use the gas price suggested by the rpc endpoint via `eth_gasPrice`.
- `Arbitrum` is a special mode only for use with Arbitrum blockchains. It uses the suggested gas price (up to `ETH_MAX_GAS_PRICE_WEI`, with `1000 gwei` default) as well as an estimated gas limit (up to `ETH_GAS_LIMIT_MAX`, with `1,000,000,000` default).

Chainlink nodes decide what gas price to use using an `Estimator`. It ships with several simple and battle-hardened built-in estimators that should work well for almost all use-cases. Note that estimators will change their behaviour slightly depending on if you are in EIP-1559 mode or not.

You can also use your own estimator for gas price by selecting the `FixedPrice` estimator and using the exposed API to set the price.

An important point to note is that the Chainlink node does _not_ ship with built-in support for go-ethereum's `estimateGas` call. This is for several reasons, including security and reliability. We have found empirically that it is not generally safe to rely on the remote ETH node's idea of what gas price should be.

### PriceDefault
```toml
PriceDefault = '20 gwei' # Default
```
PriceDefault is the default gas price to use when submitting transactions to the blockchain. Will be overridden by the built-in `BlockHistoryEstimator` if enabled, and might be increased if gas bumping is enabled.

(Only applies to legacy transactions)

Can be used with the `chainlink setgasprice` to be updated while the node is still running.

### PriceMax
```toml
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether' # Default
```
PriceMax is the maximum gas price. Chainlink nodes will never pay more than this for a transaction.
This applies to both legacy and EIP1559 transactions.
Note that it is impossible to disable the maximum limit. Setting this value to zero will prevent paying anything for any transaction (which can be useful in some rare cases).
Most chains by default have the maximum set to 2**256-1 Wei which is the maximum allowed gas price on EVM-compatible chains, and is so large it may as well be unlimited.

### PriceMin
```toml
PriceMin = '1 gwei' # Default
```
PriceMin is the minimum gas price. Chainlink nodes will never pay less than this for a transaction.

(Only applies to legacy transactions)

It is possible to force the Chainlink node to use a fixed gas price by setting a combination of these, e.g.

```toml
EIP1559DynamicFees = false
PriceMax = 100
PriceMin = 100
PriceDefault = 100
BumpThreshold = 0
Mode = 'FixedPrice'
```

### LimitDefault
```toml
LimitDefault = 500_000 # Default
```
LimitDefault sets default gas limit for outgoing transactions. This should not need to be changed in most cases.
Some job types, such as Keeper jobs, might set their own gas limit unrelated to this value.

### LimitMax
```toml
LimitMax = 500_000 # Default
```
LimitMax sets a maximum for _estimated_ gas limits. This currently only applies to `Arbitrum` `GasEstimatorMode`.

### LimitMultiplier
```toml
LimitMultiplier = '1.0' # Default
```
LimitMultiplier is the factor by which a transaction's GasLimit is multiplied before transmission. So if the value is 1.1, and the GasLimit for a transaction is 10, 10% will be added before transmission.

This factor is always applied, so includes Optimism L2 transactions which uses a default gas limit of 1 and is also applied to `LimitDefault`.

### LimitTransfer
```toml
LimitTransfer = 21_000 # Default
```
LimitTransfer is the gas limit used for an ordinary ETH transfer.

### BumpMin
```toml
BumpMin = '5 gwei' # Default
```
BumpMin is the minimum fixed amount of wei by which gas is bumped on each transaction attempt.

### BumpPercent
```toml
BumpPercent = 20 # Default
```
BumpPercent is the percentage by which to bump gas on a transaction that has exceeded `BumpThreshold`. The larger of `GasBumpPercent` and `GasBumpWei` is taken for gas bumps.

### BumpThreshold
```toml
BumpThreshold = 3 # Default
```
BumpThreshold is the number of blocks to wait for a transaction stuck in the mempool before automatically bumping the gas price. Set to 0 to disable gas bumping completely.

### BumpTxDepth
```toml
BumpTxDepth = 10 # Default
```
BumpTxDepth is the number of transactions to gas bump starting from oldest. Set to 0 for no limit (i.e. bump all).

### EIP1559DynamicFees
```toml
EIP1559DynamicFees = false # Default
```
EIP1559DynamicFees torces EIP-1559 transaction mode. Enabling EIP-1559 mode can help reduce gas costs on chains that support it. This is supported only on official Ethereum mainnet and testnets. It is not recommended to enable this setting on Polygon because the EIP-1559 fee market appears to be broken on all Polygon chains and EIP-1559 transactions are less likely to be included than legacy transactions.

#### Technical details

Chainlink nodes include experimental support for submitting transactions using type 0x2 (EIP-1559) envelope.

EIP-1559 mode is enabled by default on the Ethereum Mainnet, but can be enabled on a per-chain basis or globally.

This might help to save gas on spikes. Chainlink nodes should react faster on the upleg and avoid overpaying on the downleg. It might also be possible to set `EVM.GasEstimator.BlockHistory.BatchSize` to a smaller value such as 12 or even 6 because tip cap should be a more consistent indicator of inclusion time than total gas price. This would make Chainlink nodes more responsive and should reduce response time variance. Some experimentation is required to find optimum settings.

Set with caution, if you set this on a chain that does not actually support EIP-1559 your node will be broken.

In EIP-1559 mode, the total price for the transaction is the minimum of base fee + tip cap and fee cap. More information can be found on the [official EIP](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-1559.md).

Chainlink's implementation of EIP-1559 works as follows:

If you are using FixedPriceEstimator:
- With gas bumping disabled, it will submit all transactions with `feecap=PriceMax` and `tipcap=GasTipCapDefault`
- With gas bumping enabled, it will submit all transactions initially with `feecap=GasFeeCapDefault` and `tipcap=GasTipCapDefault`.

If you are using BlockHistoryEstimator (default for most chains):
- With gas bumping disabled, it will submit all transactions with `feecap=PriceMax` and `tipcap=<calculated using past blocks>`
- With gas bumping enabled (default for most chains) it will submit all transactions initially with `feecap = ( current block base fee * (1.125 ^ N) + tipcap )` where N is configurable by setting `EVM.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks` but defaults to `gas bump threshold+1` and `tipcap=<calculated using past blocks>`

Bumping works as follows:

- Increase tipcap by `max(tipcap * (1 + GasBumpPercent), tipcap + GasBumpWei)`
- Increase feecap by `max(feecap * (1 + GasBumpPercent), feecap + GasBumpWei)`

A quick note on terminology - Chainlink nodes use the same terms used internally by go-ethereum source code to describe various prices. This is not the same as the externally used terms. For reference:

- Base Fee Per Gas = BaseFeePerGas
- Max Fee Per Gas = FeeCap
- Max Priority Fee Per Gas = TipCap

In EIP-1559 mode, the following changes occur to how configuration works:

- All new transactions will be sent as type 0x2 transactions specifying a TipCap and FeeCap. Be aware that existing pending legacy transactions will continue to be gas bumped in legacy mode.
- `BlockHistoryEstimator` will apply its calculations (gas percentile etc) to the TipCap and this value will be used for new transactions (GasPrice will be ignored)
- `FixedPriceEstimator` will use `GasTipCapDefault` instead of `GasPriceDefault` for the tip cap
- `FixedPriceEstimator` will use `GasFeeCapDefault` instaed of `GasPriceDefault` for the fee cap
- `PriceMin` is ignored for new transactions and `GasTipCapMinimum` is used instead (default 0)
- `PriceMax` still represents that absolute upper limit that Chainlink will ever spend (total) on a single tx
- `Keeper.GasTipCapBufferPercent` is ignored in EIP-1559 mode and `Keeper.GasTipCapBufferPercent` is used instead

### FeeCapDefault
```toml
FeeCapDefault = '100 gwei' # Default
```
FeeCapDefault controls the fixed initial fee cap, if EIP1559 mode is enabled and `FixedPrice` gas estimator is used.

### TipCapDefault
```toml
TipCapDefault = '1 wei' # Default
```
TipCapDefault is the default gas tip to use when submitting transactions to the blockchain. Will be overridden by the built-in `BlockHistoryEstimator` if enabled, and might be increased if gas bumping is enabled.

(Only applies to EIP-1559 transactions)

### TipCapMin
```toml
TipCapMin = '1 wei' # Default
```
TipCapMinimum is the minimum gas tip to use when submitting transactions to the blockchain.

Only applies to EIP-1559 transactions)

## EVM.GasEstimator.LimitJobType
```toml
[EVM.GasEstimator.LimitJobType]
OCR = 100_000 # Example
DR = 100_000 # Example
VRF = 100_000 # Example
FM = 100_000 # Example
Keeper = 100_000 # Example
```


### OCR
```toml
OCR = 100_000 # Example
```
OCR overrides LimitDefault for OCR jobs.

### DR
```toml
DR = 100_000 # Example
```
DR overrides LimitDefault for Direct Request jobs.

### VRF
```toml
VRF = 100_000 # Example
```
VRF overrides LimitDefault for VRF jobs.

### FM
```toml
FM = 100_000 # Example
```
FM overrides LimitDefault for Flux Monitor jobs.

### Keeper
```toml
Keeper = 100_000 # Example
```
Keeper overrides LimitDefault for Keeper jobs.

## EVM.GasEstimator.BlockHistory
```toml
[EVM.GasEstimator.BlockHistory]
BatchSize = 4 # Default
BlockHistorySize = 8 # Default
CheckInclusionBlocks = 12 # Default
CheckInclusionPercentile = 90 # Default
EIP1559FeeCapBufferBlocks = 13 # Example
TransactionPercentile = 60 # Default
```
These settings allow you to configure how your node calculates gas prices when using the block history estimator.
In most cases, leaving these values at their defaults should give good results.

### BatchSize
```toml
BatchSize = 4 # Default
```
BatchSize sets the maximum number of blocks to fetch in one batch in the block history estimator.
If the `BatchSize` variable is set to 0, it defaults to `EVM.RPCDefaultBatchSize`.

### BlockHistorySize
```toml
BlockHistorySize = 8 # Default
```
BlockHistorySize controls the number of past blocks to keep in memory to use as a basis for calculating a percentile gas price.

### CheckInclusionBlocks
```toml
CheckInclusionBlocks = 12 # Default
```
CheckInclusionBlocks is the number of recent blocks to use to detect if there is a transaction propagation/connectivity issue, and to prevent bumping in these cases.
This can help avoid the situation where RPC nodes are not propagating transactions for some non-price-related reason (e.g. go-ethereum bug, networking issue etc) and bumping gas would not help.

Set to zero to disable connectivity checking completely.

### CheckInclusionPercentile
```toml
CheckInclusionPercentile = 90 # Default
```
CheckInclusionPercentile controls the percentile that a transaction must have been higher than for all the blocks in the inclusion check window in order to register as a connectivity issue.

For example, if CheckInclusionBlocks=12 and CheckInclusionPercentile=90 then further bumping will be prevented for any transaction with any attempt that has a higher price than the 90th percentile for the most recent 12 blocks.

### EIP1559FeeCapBufferBlocks
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
EIP1559FeeCapBufferBlocks = 13 # Example
```
EIP1559FeeCapBufferBlocks controls the buffer blocks to add to the current base fee when sending a transaction. By default, the gas bumping threshold + 1 block is used.

Only applies to EIP-1559 transactions)

### TransactionPercentile
```toml
TransactionPercentile = 60 # Default
```
TransactionPercentile specifies gas price to choose. E.g. if the block history contains four transactions with gas prices `[100, 200, 300, 400]` then picking 25 for this number will give a value of 200. If the calculated gas price is higher than `GasPriceDefault` then the higher price will be used as the base price for new transactions.

Must be in range 0-100.

Only has an effect if gas updater is enabled.

Think of this number as an indicator of how aggressive you want your node to price its transactions.

Setting this number higher will cause the Chainlink node to select higher gas prices.

Setting it lower will tend to set lower gas prices.

## EVM.HeadTracker
```toml
[EVM.HeadTracker]
HistoryDepth = 100 # Default
MaxBufferSize = 3 # Default
SamplingInterval = '1s' # Default
```
The head tracker continually listens for new heads from the chain.

In addition to these settings, it log warnings if `EVM.NoNewHeadsThreshold` is exceeded without any new blocks being emitted.

### HistoryDepth
```toml
HistoryDepth = 100 # Default
```
HistoryDepth tracks the top N block numbers to keep in the `heads` database table.
Note that this can easily result in MORE than N records since in the case of re-orgs we keep multiple heads for a particular block height.
This number should be at least as large as `FinalityDepth`.
There may be a small performance penalty to setting this to something very large (10,000+)

### MaxBufferSize
```toml
MaxBufferSize = 3 # Default
```
MaxBufferSize is the maximum number of heads that may be
buffered in front of the head tracker before older heads start to be
dropped. You may think of it as something like the maximum permittable "lag"
for the head tracker before we start dropping heads to keep up.

### SamplingInterval
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
SamplingInterval = '1s' # Default
```
SamplingInterval means that head tracker callbacks will at maximum be made once in every window of this duration. This is a performance optimisation for fast chains. Set to 0 to disable sampling entirely.

## EVM.KeySpecific
```toml
[[EVM.KeySpecific]]
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
GasEstimator.PriceMax = '79 gwei' # Example
```


### Key
```toml
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
```
Key is the account to apply these settings to

### PriceMax
```toml
GasEstimator.PriceMax = '79 gwei' # Example
```
GasEstimator.PriceMax overrides the maximum gas price for this key. See EVM.GasEstimator.PriceMax.

## EVM.NodePool
```toml
[EVM.NodePool]
PollFailureThreshold = 5 # Default
PollInterval = '10s' # Default
SelectionMode = 'HighestHead' # Default
SyncThreshold = 5 # Default
```
The node pool manages multiple RPC endpoints.

In addition to these settings, `EVM.NoNewHeadsThreshold` controls how long to wait after receiving no new heads before marking the node as out-of-sync.

### PollFailureThreshold
```toml
PollFailureThreshold = 5 # Default
```
PollFailureThreshold indicates how many consecutive polls must fail in order to mark a node as unreachable.

Set to zero to disable poll checking.

### PollInterval
```toml
PollInterval = '10s' # Default
```
PollInterval controls how often to poll the node to check for liveness.

Set to zero to disable poll checking.

### SelectionMode
```toml
SelectionMode = 'HighestHead' # Default
```
SelectionMode controls node selection strategy:
- HighestHead: use the node with the highest head number
- RoundRobin: rotate through nodes, per-request
- TotalDifficulty: use the node with the greatest total difficulty

### SyncThreshold
```toml
SyncThreshold = 5 # Default
```
SyncThreshold controls how far a node may lag behind the best node before being marked out-of-sync.
Depending on `SelectionMode`, this represents a difference in the number of blocks (`HighestHead`, `RoundRobin`), or total difficulty (`TotalDifficulty`).

Set to 0 to disable this check.

## EVM.OCR
```toml
[EVM.OCR]
ContractConfirmations = 4 # Default
ContractTransmitterTransmitTimeout = '10s' # Default
DatabaseTimeout = '10s' # Default
ObservationGracePeriod = '1s' # Default
```


### ContractConfirmations
```toml
ContractConfirmations = 4 # Default
```
ContractConfirmations sets `OCR.ContractConfirmations` for this EVM chain.

### ContractTransmitterTransmitTimeout
```toml
ContractTransmitterTransmitTimeout = '10s' # Default
```
ContractTransmitterTransmitTimeout sets `OCR.ContractTransmitterTransmitTimeout` for this EVM chain.

### DatabaseTimeout
```toml
DatabaseTimeout = '10s' # Default
```
DatabaseTimeout sets `OCR.DatabaseTimeout` for this EVM chain.

### ObservationGracePeriod
```toml
ObservationGracePeriod = '1s' # Default
```
ObservationGracePeriod sets `OCR.ObservationGracePeriod` for this EVM chain.

## EVM.Nodes
```toml
[[EVM.Nodes]]
Name = 'foo' # Example
WSURL = 'wss://web.socket/test' # Example
HTTPURL = 'https://foo.web' # Example
SendOnly = false # Default
```


### Name
```toml
Name = 'foo' # Example
```
Name is a unique (per-chain) identifier for this node.

### WSURL
```toml
WSURL = 'wss://web.socket/test' # Example
```
WSURL is the WS(S) endpoint for this node. Required for primary nodes.

### HTTPURL
```toml
HTTPURL = 'https://foo.web' # Example
```
HTTPURL is the HTTP(S) endpoint for this node. Required for all nodes.

### SendOnly
```toml
SendOnly = false # Default
```
SendOnly limits usage to sending transaction broadcasts only. With this enabled, only HTTPURL is required, and WSURL is not used.

## EVM.OCR2.Automation
```toml
[EVM.OCR2.Automation]
GasLimit = 5300000 # Default
```


### GasLimit
```toml
GasLimit = 5300000 # Default
```
GasLimit controls the gas limit for transmit transactions from ocr2automation job.

## Cosmos
```toml
[[Cosmos]]
ChainID = 'Malaga-420' # Example
Enabled = true # Default
BlockRate = '6s' # Default
BlocksUntilTxTimeout = 30 # Default
ConfirmPollPeriod = '1s' # Default
FallbackGasPriceUAtom = '0.015' # Default
FCDURL = 'http://cosmos.com' # Example
GasLimitMultiplier = '1.5' # Default
MaxMsgsPerBatch = 100 # Default
OCR2CachePollPeriod = '4s' # Default
OCR2CacheTTL = '1m' # Default
TxMsgTimeout = '10m' # Default
```


### ChainID
```toml
ChainID = 'Malaga-420' # Example
```
ChainID is the Cosmos chain ID. Mandatory.

### Enabled
```toml
Enabled = true # Default
```
Enabled enables this chain.

### BlockRate
```toml
BlockRate = '6s' # Default
```
BlockRate is the average time between blocks.

### BlocksUntilTxTimeout
```toml
BlocksUntilTxTimeout = 30 # Default
```
BlocksUntilTxTimeout is the number of blocks to wait before giving up on the tx getting confirmed.

### ConfirmPollPeriod
```toml
ConfirmPollPeriod = '1s' # Default
```
ConfirmPollPeriod sets how often check for tx confirmation.

### FallbackGasPriceUAtom
```toml
FallbackGasPriceUAtom = '0.015' # Default
```
FallbackGasPriceUAtom sets a fallback gas price to use when the estimator is not available.

### FCDURL
```toml
FCDURL = 'http://cosmos.com' # Example
```
FCDURL sets the FCD (Full Client Daemon) URL.

### GasLimitMultiplier
```toml
GasLimitMultiplier = '1.5' # Default
```
GasLimitMultiplier scales the estimated gas limit.

### MaxMsgsPerBatch
```toml
MaxMsgsPerBatch = 100 # Default
```
MaxMsgsPerBatch limits the numbers of mesages per transaction batch.

### OCR2CachePollPeriod
```toml
OCR2CachePollPeriod = '4s' # Default
```
OCR2CachePollPeriod is the rate to poll for the OCR2 state cache.

### OCR2CacheTTL
```toml
OCR2CacheTTL = '1m' # Default
```
OCR2CacheTTL is the stale OCR2 cache deadline.

### TxMsgTimeout
```toml
TxMsgTimeout = '10m' # Default
```
TxMsgTimeout is the maximum age for resending transaction before they expire.

## Cosmos.Nodes
```toml
[[Cosmos.Nodes]]
Name = 'primary' # Example
TendermintURL = 'http://tender.mint' # Example
```


### Name
```toml
Name = 'primary' # Example
```
Name is a unique (per-chain) identifier for this node.

### TendermintURL
```toml
TendermintURL = 'http://tender.mint' # Example
```
TendermintURL is the HTTP(S) tendermint endpoint for this node.

## Solana
```toml
[[Solana]]
ChainID = 'mainnet' # Example
Enabled = false # Default
BalancePollPeriod = '5s' # Default
ConfirmPollPeriod = '500ms' # Default
OCR2CachePollPeriod = '1s' # Default
OCR2CacheTTL = '1m' # Default
TxTimeout = '1m' # Default
TxRetryTimeout = '10s' # Default
TxConfirmTimeout = '30s' # Default
SkipPreflight = true # Default
Commitment = 'confirmed' # Default
MaxRetries = 0 # Default
FeeEstimatorMode = 'fixed' # Default
ComputeUnitPriceMax = 1000 # Default
ComputeUnitPriceMin = 0 # Default
ComputeUnitPriceDefault = 0 # Default
FeeBumpPeriod = '3s' # Default
```


### ChainID
```toml
ChainID = 'mainnet' # Example
```
ChainID is the Solana chain ID. Must be one of: mainnet, testnet, devnet, localnet. Mandatory.

### Enabled
```toml
Enabled = false # Default
```
Enabled enables this chain.

### BalancePollPeriod
```toml
BalancePollPeriod = '5s' # Default
```
BalancePollPeriod is the rate to poll for SOL balance and update Prometheus metrics.

### ConfirmPollPeriod
```toml
ConfirmPollPeriod = '500ms' # Default
```
ConfirmPollPeriod is the rate to poll for signature confirmation.

### OCR2CachePollPeriod
```toml
OCR2CachePollPeriod = '1s' # Default
```
OCR2CachePollPeriod is the rate to poll for the OCR2 state cache.

### OCR2CacheTTL
```toml
OCR2CacheTTL = '1m' # Default
```
OCR2CacheTTL is the stale OCR2 cache deadline.

### TxTimeout
```toml
TxTimeout = '1m' # Default
```
TxTimeout is the timeout for sending txes to an RPC endpoint.

### TxRetryTimeout
```toml
TxRetryTimeout = '10s' # Default
```
TxRetryTimeout is the duration for tx manager to attempt rebroadcasting to RPC, before giving up.

### TxConfirmTimeout
```toml
TxConfirmTimeout = '30s' # Default
```
TxConfirmTimeout is the duration to wait when confirming a tx signature, before discarding as unconfirmed.

### SkipPreflight
```toml
SkipPreflight = true # Default
```
SkipPreflight enables or disables preflight checks when sending txs.

### Commitment
```toml
Commitment = 'confirmed' # Default
```
Commitment is the confirmation level for solana state and transactions. ([documentation](https://docs.solana.com/developing/clients/jsonrpc-api#configuring-state-commitment))

### MaxRetries
```toml
MaxRetries = 0 # Default
```
MaxRetries is the maximum number of times the RPC node will automatically rebroadcast a tx.
The default is 0 for custom txm rebroadcasting method, set to -1 to use the RPC node's default retry strategy.

### FeeEstimatorMode
```toml
FeeEstimatorMode = 'fixed' # Default
```
FeeEstimatorMode is the method used to determine the base fee

### ComputeUnitPriceMax
```toml
ComputeUnitPriceMax = 1000 # Default
```
ComputeUnitPriceMax is the maximum price per compute unit that a transaction can be bumped to

### ComputeUnitPriceMin
```toml
ComputeUnitPriceMin = 0 # Default
```
ComputeUnitPriceMin is the minimum price per compute unit that transaction can have

### ComputeUnitPriceDefault
```toml
ComputeUnitPriceDefault = 0 # Default
```
ComputeUnitPriceDefault is the default price per compute unit price, and the starting base fee when FeeEstimatorMode = 'fixed'

### FeeBumpPeriod
```toml
FeeBumpPeriod = '3s' # Default
```
FeeBumpPeriod is the amount of time before a tx is retried with a fee bump

## Solana.Nodes
```toml
[[Solana.Nodes]]
Name = 'primary' # Example
URL = 'http://solana.web' # Example
```


### Name
```toml
Name = 'primary' # Example
```
Name is a unique (per-chain) identifier for this node.

### URL
```toml
URL = 'http://solana.web' # Example
```
URL is the HTTP(S) endpoint for this node.

## Starknet
```toml
[[Starknet]]
ChainID = 'foobar' # Example
Enabled = true # Default
OCR2CachePollPeriod = '5s' # Default
OCR2CacheTTL = '1m' # Default
RequestTimeout = '10s' # Default
TxTimeout = '1m' # Default
TxSendFrequency = '5s' # Default
TxMaxBatchSize = 100 # Default
```


### ChainID
```toml
ChainID = 'foobar' # Example
```
ChainID is the Starknet chain ID.

### Enabled
```toml
Enabled = true # Default
```
Enabled enables this chain.

### OCR2CachePollPeriod
```toml
OCR2CachePollPeriod = '5s' # Default
```
OCR2CachePollPeriod is the rate to poll for the OCR2 state cache.

### OCR2CacheTTL
```toml
OCR2CacheTTL = '1m' # Default
```
OCR2CacheTTL is the stale OCR2 cache deadline.

### RequestTimeout
```toml
RequestTimeout = '10s' # Default
```
RequestTimeout is the RPC client timeout.

### TxTimeout
```toml
TxTimeout = '1m' # Default
```
TxTimeout is the timeout for sending txes to an RPC endpoint.

### TxSendFrequency
```toml
TxSendFrequency = '5s' # Default
```
TxSendFrequency is how often to broadcast batches of txes.

### TxMaxBatchSize
```toml
TxMaxBatchSize = 100 # Default
```
TxMaxBatchSize limits the size of tx batches.

## Starknet.Nodes
```toml
[[Starknet.Nodes]]
Name = 'primary' # Example
URL = 'http://stark.node' # Example
```


### Name
```toml
Name = 'primary' # Example
```
Name is a unique (per-chain) identifier for this node.

### URL
```toml
URL = 'http://stark.node' # Example
```
URL is the base HTTP(S) endpoint for this node.

