[//]: # (Documentation generated from docs.toml - DO NOT EDIT.)

## Table of contents

- [Global](#Global)
- [Feature](#Feature)
- [Database](#Database)
	- [Backup](#Database-Backup)
	- [Listener](#Database-Listener)
	- [Lock](#Database-Lock)
- [TelemetryIngress](#TelemetryIngress)
- [Log](#Log)
- [WebServer](#WebServer)
	- [RateLimit](#WebServer-RateLimit)
	- [MFA](#WebServer-MFA)
	- [TLS](#WebServer-TLS)
- [JobPipeline](#JobPipeline)
- [FluxMonitor](#FluxMonitor)
- [OCR2](#OCR2)
- [OCR](#OCR)
- [P2P](#P2P)
	- [V1](#P2P-V1)
	- [V2](#P2P-V2)
- [Keeper](#Keeper)
- [AutoPprof](#AutoPprof)
- [Sentry](#Sentry)
- [EVM](#EVM)
	- [BalanceMonitor](#EVM-BalanceMonitor)
	- [GasEstimator](#EVM-GasEstimator)
		- [BlockHistory](#EVM-GasEstimator-BlockHistory)
	- [HeadTracker](#EVM-HeadTracker)
	- [KeySpecific](#EVM-KeySpecific)
	- [NodePool](#EVM-NodePool)
	- [OCR](#EVM-OCR)
	- [Nodes](#EVM-Nodes)
- [Solana](#Solana)
	- [Nodes](#Solana-Nodes)
- [Terra](#Terra)
	- [Nodes](#Terra-Nodes)

## Global<a id='Global'></a>
```toml
ExplorerURL = 'ws://explorer.url' # Example
InsecureFastScrypt = false # Default
RootDir = '~/.chainlink' # Default
ShutdownGracePeriod = '5s' # Default
```


### ExplorerURL<a id='ExplorerURL'></a>
```toml
ExplorerURL = 'ws://explorer.url' # Example
```
ExplorerURL is the websocket URL for the node to push stats to.

### InsecureFastScrypt<a id='InsecureFastScrypt'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
InsecureFastScrypt = false # Default
```
InsecureFastScrypt causes all key stores to encrypt using "fast" scrypt params instead. This is insecure and only useful for local testing. DO NOT ENABLE THIS IN PRODUCTION.

### RootDir<a id='RootDir'></a>
```toml
RootDir = '~/.chainlink' # Default
```
RootDir is the Chainlink node's root directory. This is the default directory for logging, database backups, cookies, and other misc Chainlink node files. Chainlink nodes will always ensure this directory has 700 permissions because it might contain sensitive data.

### ShutdownGracePeriod<a id='ShutdownGracePeriod'></a>
```toml
ShutdownGracePeriod = '5s' # Default
```
ShutdownGracePeriod is the maximum time allowed to shut down gracefully. If exceeded, the node will terminate immediately to avoid being SIGKILLed.

## Feature<a id='Feature'></a>
```toml
[Feature]
FeedsManager = false # Default
LogPoller = false # Default
UICSA = false # Default
```


### FeedsManager<a id='Feature-FeedsManager'></a>
```toml
FeedsManager = false # Default
```
FeedsManager enables the experimental feeds manager service.

### LogPoller<a id='Feature-LogPoller'></a>
```toml
LogPoller = false # Default
```
LogPoller enables the log poller, an experimental approach to processing logs, required if also using Evm.UseForwarders or OCR2.

### UICSA<a id='Feature-UICSA'></a>
```toml
UICSA = false # Default
```
TODO

## Database<a id='Database'></a>
```toml
[Database]
DefaultIdleInTxSessionTimeout = '1h' # Default
DefaultLockTimeout = '15s' # Default
DefaultQueryTimeout = '10s' # Default
MigrateOnStartup = true # Default
ORMMaxIdleConns = 10 # Default
ORMMaxOpenConns = 20 # Default
```


### DefaultIdleInTxSessionTimeout<a id='Database-DefaultIdleInTxSessionTimeout'></a>
```toml
DefaultIdleInTxSessionTimeout = '1h' # Default
```
DefaultIdleInTxSessionTimeout is the maximum time allowed for queries to idle in transaction before timing out.

### DefaultLockTimeout<a id='Database-DefaultLockTimeout'></a>
```toml
DefaultLockTimeout = '15s' # Default
```
DefaultLockTimeout is the maximum time allowed for a query stuck waiting to take a lock before timing out.

### DefaultQueryTimeout<a id='Database-DefaultQueryTimeout'></a>
```toml
DefaultQueryTimeout = '10s' # Default
```
DefaultQueryTimeout is the maximum time allowed for standard queries before timing out.

### MigrateOnStartup<a id='Database-MigrateOnStartup'></a>
```toml
MigrateOnStartup = true # Default
```
MigrateOnStartup controls whether a Chainlink node will attempt to automatically migrate the database on boot. If you want more control over your database migration process, set this variable to `false` and manually migrate the database using the CLI `migrate` command instead.

### ORMMaxIdleConns<a id='Database-ORMMaxIdleConns'></a>
```toml
ORMMaxIdleConns = 10 # Default
```
ORMMaxIdleConns configures the maximum number of idle database connections that the Chainlink node will keep open. Think of this as the baseline number of database connections per Chainlink node instance. Increasing this number can help to improve performance under database-heavy workloads.

Postgres has connection limits, so you must use cation when increasing this value. If you are running several instances of a Chainlink node or another application on a single database server, you might run out of Postgres connection slots if you raise this value too high.

### ORMMaxOpenConns<a id='Database-ORMMaxOpenConns'></a>
```toml
ORMMaxOpenConns = 20 # Default
```
ORMMaxOpenConns configures the maximum number of database connections that a Chainlink node will have open at any one time. Think of this as the maximum burst upper bound limit of database connections per Chainlink node instance. Increasing this number can help to improve performance under database-heavy workloads.

Postgres has connection limits, so you must use cation when increasing this value. If you are running several instances of a Chainlink node or another application on a single database server, you might run out of Postgres connection slots if you raise this value too high.

## Database.Backup<a id='Database-Backup'></a>
```toml
[Database.Backup]
Mode = 'none' # Default
Dir = 'test/backup/dir' # Example
OnVersionUpgrade = true # Default
URL = 'http://test.back.up/fake' # Example
Frequency = '1h' # Default
```
As a best practice, take regular database backups in case of accidental data loss. This best practice is especially important when you upgrade your Chainlink node to a new version. Chainlink nodes support automated database backups to make this process easier.

NOTE: Dumps can cause high load and massive database latencies, which will negatively impact the normal functioning of the Chainlink node. For this reason, it is recommended to set a `URL` and point it to a read replica if you enable automatic backups.

### Mode<a id='Database-Backup-Mode'></a>
```toml
Mode = 'none' # Default
```
Mode sets the type of automatic database backup, which can be one of _none_, `lite`, or `full`. If enabled, the Chainlink node will always dump a backup on every boot before running migrations. Additionally, it will automatically take database backups that overwrite the backup file for the given version at regular intervals if `Frequency` is set to a non-zero interval.

_none_ - Disables backups.
`lite` - Dumps small tables including configuration and keys that are essential for the node to function, which excludes historical data like job runs, transaction history, etc.
`full` - Dumps the entire database.

It will write to a file like `$ROOT/backup/cl_backup_<VERSION>.dump`. There is one backup dump file per version of the Chainlink node. If you upgrade the node, it will keep the backup taken right before the upgrade migration so you can restore to an older version if necessary.

### Dir<a id='Database-Backup-Dir'></a>
```toml
Dir = 'test/backup/dir' # Example
```
Dir sets the directory to use for saving the backup file. Use this if you want to save the backup file in a directory other than the default ROOT directory.

### OnVersionUpgrade<a id='Database-Backup-OnVersionUpgrade'></a>
```toml
OnVersionUpgrade = true # Default
```
OnVersionUpgrade enables automatic backups of the database before running migrations, when you are upgrading to a new version.

### URL<a id='Database-Backup-URL'></a>
```toml
URL = 'http://test.back.up/fake' # Example
```
URL, if specified, is an alternative for the automatic database backup to use instead of the main database url.

It is recommended to set this value to a _read replica_ if you have one to avoid excessive load on the main database.

### Frequency<a id='Database-Backup-Frequency'></a>
```toml
Frequency = '1h' # Default
```
Frequency sets the interval for database dumps, if set to a positive duration and `Mode` is not _none_.

Set to `0` to disable periodic backups.

## Database.Listener<a id='Database-Listener'></a>
:warning: **_ADVANCED_**: _Do not change these settings unless you know what you are doing._
```toml
[Database.Listener]
MaxReconnectDuration = '10m' # Default
MinReconnectInterval = '1m' # Default
FallbackPollInterval = '30s' # Default
```
These settings control the postgres event listener.

### MaxReconnectDuration<a id='Database-Listener-MaxReconnectDuration'></a>
```toml
MaxReconnectDuration = '10m' # Default
```
MaxReconnectDuration is the maximum duration to wait between reconnect attempts.

### MinReconnectInterval<a id='Database-Listener-MinReconnectInterval'></a>
```toml
MinReconnectInterval = '1m' # Default
```
MinReconnectInterval controls the duration to wait before trying to re-establish the database connection after connection loss. After each consecutive failure this interval is doubled, until MaxReconnectInterval is reached.  Successfully completing the connection establishment procedure resets the interval back to MinReconnectInterval.

### FallbackPollInterval<a id='Database-Listener-FallbackPollInterval'></a>
```toml
FallbackPollInterval = '30s' # Default
```
FallbackPollInterval controls how often clients should manually poll as a fallback in case the postgres event was missed/dropped.

## Database.Lock<a id='Database-Lock'></a>
:warning: **_ADVANCED_**: _Do not change these settings unless you know what you are doing._
```toml
[Database.Lock]
LeaseDuration = '10s' # Default
LeaseRefreshInterval = '1s' # Default
```
Ideally, you should use a container orchestration system like [Kubernetes](https://kubernetes.io/) to ensure that only one Chainlink node instance can ever use a specific Postgres database. However, some node operators do not have the technical capacity to do this. Common use cases run multiple Chainlink node instances in failover mode as recommended by our official documentation. The first instance takes a lock on the database and subsequent instances will wait trying to take this lock in case the first instance fails.

- If your nodes or applications hold locks open for several hours or days, Postgres is unable to complete internal cleanup tasks. The Postgres maintainers explicitly discourage holding locks open for long periods of time.

Because of the complications with advisory locks, Chainlink nodes with v2.0 and later only support `lease` locking mode. The `lease` locking mode works using the following process:

- Node A creates one row in the database with the client ID and updates it once per second.
- Node B spinlocks and checks periodically to see if the client ID is too old. If the client ID is not updated after a period of time, node B assumes that node A failed and takes over. Node B becomes the owner of the row and updates the client ID once per second.
- If node A comes back, it attempts to take out a lease, realizes that the database has been leased to another process, and exits the entire application immediately.

### LeaseDuration<a id='Database-Lock-LeaseDuration'></a>
```toml
LeaseDuration = '10s' # Default
```
LeaseDuration is how long the lease lock will last before expiring.

This setting applies only if `Mode` is set to enable lease locking.

### LeaseRefreshInterval<a id='Database-Lock-LeaseRefreshInterval'></a>
```toml
LeaseRefreshInterval = '1s' # Default
```
LeaseRefreshInterval determines how often to refresh the lease lock. Also controls how often a standby node will check to see if it can grab the lease.

This setting applies only if Mode is set to enable lease locking.

## TelemetryIngress<a id='TelemetryIngress'></a>
```toml
[TelemetryIngress]
UniConn = true # Default
Logging = false # Default
ServerPubKey = 'test-pub-key' # Default
URL = 'https://prom.test' # Example
BufferSize = 100 # Default
MaxBatchSize = 50 # Default
SendInterval = '500ms' # Default
SendTimeout = '10s' # Default
UseBatchSend = true # Default
```


### UniConn<a id='TelemetryIngress-UniConn'></a>
```toml
UniConn = true # Default
```
UniConn toggles which ws connection style is used.

### Logging<a id='TelemetryIngress-Logging'></a>
```toml
Logging = false # Default
```
Logging toggles verbose logging of the raw telemetry messages being sent.

### ServerPubKey<a id='TelemetryIngress-ServerPubKey'></a>
```toml
ServerPubKey = 'test-pub-key' # Default
```
ServerPubKey is the public key of the telemetry server.

### URL<a id='TelemetryIngress-URL'></a>
```toml
URL = 'https://prom.test' # Example
```
URL is where to send telemetry.

### BufferSize<a id='TelemetryIngress-BufferSize'></a>
```toml
BufferSize = 100 # Default
```
BufferSize is the number of telemetry messages to buffer before dropping new ones.

### MaxBatchSize<a id='TelemetryIngress-MaxBatchSize'></a>
```toml
MaxBatchSize = 50 # Default
```
MaxBatchSize is the maximum number of messages to batch into one telemetry request.

### SendInterval<a id='TelemetryIngress-SendInterval'></a>
```toml
SendInterval = '500ms' # Default
```
SendInterval determines how often batched telemetry is sent to the ingress server.

### SendTimeout<a id='TelemetryIngress-SendTimeout'></a>
```toml
SendTimeout = '10s' # Default
```
SendTimeout is the max duration to wait for the request to complete when sending batch telemetry.

### UseBatchSend<a id='TelemetryIngress-UseBatchSend'></a>
```toml
UseBatchSend = true # Default
```
UseBatchSend toggles sending telemetry to the ingress server using the batch client.

## Log<a id='Log'></a>
```toml
[Log]
DatabaseQueries = false # Default
JSONConsole = false # Default
FileDir = '/my/log/directory' # Example
FileMaxSize = '5120mb' # Default
FileMaxAgeDays = 0 # Default
FileMaxBackups = 1 # Default
UnixTS = false # Default
```


### DatabaseQueries<a id='Log-DatabaseQueries'></a>
```toml
DatabaseQueries = false # Default
```
DatabaseQueries tells the Chainlink node to log database queries made using the default logger. SQL statements will be logged at `debug` level. Not all statements can be logged. The best way to get a true log of all SQL statements is to enable SQL statement logging on Postgres.

### JSONConsole<a id='Log-JSONConsole'></a>
```toml
JSONConsole = false # Default
```
JSONConsole enables JSON logging. Otherwise, the log is saved in a human-friendly console format.

### FileDir<a id='Log-FileDir'></a>
```toml
FileDir = '/my/log/directory' # Example
```
FileDir sets the log directory. By default, Chainlink nodes write log data to `$ROOT/log.jsonl`.

### FileMaxSize<a id='Log-FileMaxSize'></a>
```toml
FileMaxSize = '5120mb' # Default
```
FileMaxSize determines the log file's max size in megabytes before file rotation. Having this not set will disable logging to disk. If your disk doesn't have enough disk space, the logging will pause and the application will log errors until space is available again.

Values must have suffixes with a unit like: `5120mb` (5,120 megabytes). If no unit suffix is provided, the value defaults to `b` (bytes). The list of valid unit suffixes are:

- b (bytes)
- kb (kilobytes)
- mb (megabytes)
- gb (gigabytes)
- tb (terabytes)

### FileMaxAgeDays<a id='Log-FileMaxAgeDays'></a>
```toml
FileMaxAgeDays = 0 # Default
```
FileMaxAgeDays determines the log file's max age in days before file rotation. Keeping this config with the default value will not remove log files based on age.

### FileMaxBackups<a id='Log-FileMaxBackups'></a>
```toml
FileMaxBackups = 1 # Default
```
FileMaxBackups determines the maximum number of old log files to retain. Keeping this config with the default value retains all old log files. The `FileMaxAgeDays` variable can still cause them to get deleted.

### UnixTS<a id='Log-UnixTS'></a>
```toml
UnixTS = false # Default
```
UnixTS enables legacy unix timestamps.

Previous versions of Chainlink nodes wrote JSON logs with a unix timestamp. As of v1.1.0 and up, the default has changed to use ISO8601 timestamps for better readability.

## WebServer<a id='WebServer'></a>
```toml
[WebServer]
AllowOrigins = 'http://localhost:3000,http://localhost:6688' # Default
BridgeResponseURL = 'https://my-chainlink-node.example.com:6688' # Example
HTTPWriteTimeout = '10s' # Default
HTTPPort = 6688 # Default
SecureCookies = true # Default
SessionTimeout = '15m' # Default
SessionReaperExpiration = '240h' # Default
```


### AllowOrigins<a id='WebServer-AllowOrigins'></a>
```toml
AllowOrigins = 'http://localhost:3000,http://localhost:6688' # Default
```
AllowOrigins controls the URLs Chainlink nodes emit in the `Allow-Origins` header of its API responses. The setting can be a comma-separated list with no spaces. You might experience CORS issues if this is not set correctly.

You should set this to the external URL that you use to access the Chainlink UI.

You can set `AllowOrigins = '*'` to allow the UI to work from any URL, but it is recommended for security reasons to make it explicit instead.

### BridgeResponseURL<a id='WebServer-BridgeResponseURL'></a>
```toml
BridgeResponseURL = 'https://my-chainlink-node.example.com:6688' # Example
```
BridgeResponseURL defines the URL for bridges to send a response to. This _must_ be set when using async external adapters.

Usually this will be the same as the URL/IP and port you use to connect to the Chainlink UI.

### HTTPWriteTimeout<a id='WebServer-HTTPWriteTimeout'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
HTTPWriteTimeout = '10s' # Default
```
HTTPWriteTimeout controls how long the Chainlink node's API server can hold a socket open for writing a response to an HTTP request. Sometimes, this must be increased for pprof.

### HTTPPort<a id='WebServer-HTTPPort'></a>
```toml
HTTPPort = 6688 # Default
```
HTTPPort is the port used for the Chainlink Node API, [CLI](/docs/configuration-variables/#cli-client), and GUI.

### SecureCookies<a id='WebServer-SecureCookies'></a>
```toml
SecureCookies = true # Default
```
SecureCookies requires the use of secure cookies for authentication. Set to false to enable standard HTTP requests along with `TLSPort = 0`.

### SessionTimeout<a id='WebServer-SessionTimeout'></a>
```toml
SessionTimeout = '15m' # Default
```
SessionTimeout determines the amount of idle time to elapse before session cookies expire. This signs out GUI users from their sessions.

### SessionReaperExpiration<a id='WebServer-SessionReaperExpiration'></a>
```toml
SessionReaperExpiration = '240h' # Default
```
SessionReaperExpiration represents how long an API session lasts before expiring and requiring a new login.

## WebServer.RateLimit<a id='WebServer-RateLimit'></a>
```toml
[WebServer.RateLimit]
Authenticated = 42 # Default
AuthenticatedPeriod = '1m' # Default
Unauthenticated = 5 # Default
UnauthenticatedPeriod = '20s' # Default
```


### Authenticated<a id='WebServer-RateLimit-Authenticated'></a>
```toml
Authenticated = 42 # Default
```
Authenticated defines the threshold to which authenticated requests get limited. More than this many authenticated requests per `AuthenticatedRateLimitPeriod` will be rejected.

### AuthenticatedPeriod<a id='WebServer-RateLimit-AuthenticatedPeriod'></a>
```toml
AuthenticatedPeriod = '1m' # Default
```
AuthenticatedPeriod defines the period to which authenticated requests get limited.

### Unauthenticated<a id='WebServer-RateLimit-Unauthenticated'></a>
```toml
Unauthenticated = 5 # Default
```
Unauthenticated defines the threshold to which authenticated requests get limited. More than this many unauthenticated requests per `UnAuthenticatedRateLimitPeriod` will be rejected.

### UnauthenticatedPeriod<a id='WebServer-RateLimit-UnauthenticatedPeriod'></a>
```toml
UnauthenticatedPeriod = '20s' # Default
```
UnauthenticatedPeriod defines the period to which unauthenticated requests get limited.

## WebServer.MFA<a id='WebServer-MFA'></a>
```toml
[WebServer.MFA]
RPID = 'localhost' # Example
RPOrigin = 'http://localhost:6688/' # Example
```
The Operator UI frontend supports enabling Multi Factor Authentication via Webauthn per account. When enabled, logging in will require the account password and a hardware or OS security key such as Yubikey. To enroll, log in to the operator UI and click the circle purple profile button at the top right and then click **Register MFA Token**. Tap your hardware security key or use the OS public key management feature to enroll a key. Next time you log in, this key will be required to authenticate.

### RPID<a id='WebServer-MFA-RPID'></a>
```toml
RPID = 'localhost' # Example
```
RPID is the FQDN of where the Operator UI is served. When serving locally, the value should be `localhost`.

### RPOrigin<a id='WebServer-MFA-RPOrigin'></a>
```toml
RPOrigin = 'http://localhost:6688/' # Example
```
RPOrigin is the origin URL where WebAuthn requests initiate, including scheme and port. When serving locally, the value should be `http://localhost:6688/`.

## WebServer.TLS<a id='WebServer-TLS'></a>
```toml
[WebServer.TLS]
CertPath = '/home/$USER/.chainlink/tls/server.crt' # Example
Host = 'tls-host' # Example
KeyPath = '/home/$USER/.chainlink/tls/server.key' # Example
HTTPSPort = 6689 # Default
ForceRedirect = false # Default
```
The TLS settings apply only if you want to enable TLS security on your Chainlink node.

### CertPath<a id='WebServer-TLS-CertPath'></a>
```toml
CertPath = '/home/$USER/.chainlink/tls/server.crt' # Example
```
CertPath is the location of the TLS certificate file.

### Host<a id='WebServer-TLS-Host'></a>
```toml
Host = 'tls-host' # Example
```
Host is the hostname configured for TLS to be used by the Chainlink node. This is useful if you configured a domain name specific for your Chainlink node.

### KeyPath<a id='WebServer-TLS-KeyPath'></a>
```toml
KeyPath = '/home/$USER/.chainlink/tls/server.key' # Example
```
KeyPath is the location of the TLS private key file.

### HTTPSPort<a id='WebServer-TLS-HTTPSPort'></a>
```toml
HTTPSPort = 6689 # Default
```
HTTPSPort is the port used for HTTPS connections. Set this to `0` to disable HTTPS. Disabling HTTPS also relieves Chainlink nodes of the requirement for a TLS certificate.

### ForceRedirect<a id='WebServer-TLS-ForceRedirect'></a>
```toml
ForceRedirect = false # Default
```
ForceRedirect forces TLS redirect for unencrypted connections.

## JobPipeline<a id='JobPipeline'></a>
```toml
[JobPipeline]
HTTPRequestMaxSize = '32768' # Default
DefaultHTTPRequestTimeout = '15s' # Default
ExternalInitiatorsEnabled = false # Default
MaxRunDuration = '10m' # Default
ReaperInterval = '1h' # Default
ReaperThreshold = '24h' # Default
ResultWriteQueueDepth = 100 # Default
```


### HTTPRequestMaxSize<a id='JobPipeline-HTTPRequestMaxSize'></a>
```toml
HTTPRequestMaxSize = '32768' # Default
```
HTTPRequestMaxSize defines the maximum size for HTTP requests and responses made by `http` and `bridge` adapters.

### DefaultHTTPRequestTimeout<a id='JobPipeline-DefaultHTTPRequestTimeout'></a>
```toml
DefaultHTTPRequestTimeout = '15s' # Default
```
DefaultHTTPRequestTimeout defines the default timeout for HTTP requests made by `http` and `bridge` adapters.

### ExternalInitiatorsEnabled<a id='JobPipeline-ExternalInitiatorsEnabled'></a>
```toml
ExternalInitiatorsEnabled = false # Default
```
ExternalInitiatorsEnabled enables the External Initiator feature. If disabled, `webhook` jobs can ONLY be initiated by a logged-in user. If enabled, `webhook` jobs can be initiated by a whitelisted external initiator.

### MaxRunDuration<a id='JobPipeline-MaxRunDuration'></a>
```toml
MaxRunDuration = '10m' # Default
```
MaxRunDuration is the maximum time allowed for a single job run. If it takes longer, it will exit early and be marked errored. If set to zero, disables the time limit completely.

### ReaperInterval<a id='JobPipeline-ReaperInterval'></a>
```toml
ReaperInterval = '1h' # Default
```
ReaperInterval controls how often the job pipeline reaper will run to delete completed jobs older than ReaperThreshold, in order to keep database size manageable.

Set to `0` to disable the periodic reaper.

### ReaperThreshold<a id='JobPipeline-ReaperThreshold'></a>
```toml
ReaperThreshold = '24h' # Default
```
ReaperThreshold determines the age limit for job runs. Completed job runs older than this will be automatically purged from the database.

### ResultWriteQueueDepth<a id='JobPipeline-ResultWriteQueueDepth'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
ResultWriteQueueDepth = 100 # Default
```
ResultWriteQueueDepth controls how many writes will be buffered before subsequent writes are dropped, for jobs that write results asynchronously for performance reasons, such as OCR.

## FluxMonitor<a id='FluxMonitor'></a>
```toml
[FluxMonitor]
DefaultTransactionQueueDepth = 1 # Default
SimulateTransactions = false # Default
```


### DefaultTransactionQueueDepth<a id='FluxMonitor-DefaultTransactionQueueDepth'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DefaultTransactionQueueDepth = 1 # Default
```
DefaultTransactionQueueDepth controls the queue size for `DropOldestStrategy` in Flux Monitor. Set to 0 to use `SendEvery` strategy instead.

### SimulateTransactions<a id='FluxMonitor-SimulateTransactions'></a>
```toml
SimulateTransactions = false # Default
```
SimulateTransactions enables transaction simulation for Flux Monitor.

## OCR2<a id='OCR2'></a>
```toml
[OCR2]
Enabled = true # Default
ContractConfirmations = 3 # Default
BlockchainTimeout = '20s' # Default
ContractPollInterval = '1m' # Default
ContractSubscribeInterval = '2m' # Default
ContractTransmitterTransmitTimeout = '10s' # Default
DatabaseTimeout = '10s' # Default
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000' # Example
```


### Enabled<a id='OCR2-Enabled'></a>
```toml
Enabled = true # Default
```
Enabled enables OCR2 jobs.

### ContractConfirmations<a id='OCR2-ContractConfirmations'></a>
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

### BlockchainTimeout<a id='OCR2-BlockchainTimeout'></a>
```toml
BlockchainTimeout = '20s' # Default
```
BlockchainTimeout is the timeout for blockchain queries (mediated through
ContractConfigTracker and ContractTransmitter).
(This is necessary because an oracle's operations are serialized, so
blocking forever on a chain interaction would break the oracle.)

### ContractPollInterval<a id='OCR2-ContractPollInterval'></a>
```toml
ContractPollInterval = '1m' # Default
```
ContractPollInterval is the polling interval at which ContractConfigTracker is queried for# updated on-chain configurations. Recommended values are between
fifteen seconds and two minutes.

### ContractSubscribeInterval<a id='OCR2-ContractSubscribeInterval'></a>
```toml
ContractSubscribeInterval = '2m' # Default
```
ContractSubscribeInterval is the interval at which we try to establish a subscription on ContractConfigTracker
if one doesn't exist. Recommended values are between two and five minutes.

### ContractTransmitterTransmitTimeout<a id='OCR2-ContractTransmitterTransmitTimeout'></a>
```toml
ContractTransmitterTransmitTimeout = '10s' # Default
```
ContractTransmitterTransmitTimeout is the timeout for ContractTransmitter.Transmit calls.

### DatabaseTimeout<a id='OCR2-DatabaseTimeout'></a>
```toml
DatabaseTimeout = '10s' # Default
```
DatabaseTimeout is the timeout for database interactions.
(This is necessary because an oracle's operations are serialized, so
blocking forever on an observation would break the oracle.)

### KeyBundleID<a id='OCR2-KeyBundleID'></a>
```toml
KeyBundleID = '7a5f66bbe6594259325bf2b4f5b1a9c900000000000000000000000000000000' # Example
```
KeyBundleID is a sha256 hexadecimal hash identifier.

## OCR<a id='OCR'></a>
```toml
[OCR]
Enabled = true # Default
ObservationTimeout = '5s' # Default
BlockchainTimeout = '20s' # Default
ContractPollInterval = '1m' # Default
ContractSubscribeInterval = '2m' # Default
DefaultTransactionQueueDepth = 1 # Default
KeyBundleID = 'acdd42797a8b921b2910497badc5000600000000000000000000000000000000' # Example
SimulateTransactions = false # Default
TransmitterAddress = '0xa0788FC17B1dEe36f057c42B6F373A34B014687e' # Example
```
This section applies only if you are running off-chain reporting jobs.

### Enabled<a id='OCR-Enabled'></a>
```toml
Enabled = true # Default
```
Enabled enables OCR jobs.

### ObservationTimeout<a id='OCR-ObservationTimeout'></a>
```toml
ObservationTimeout = '5s' # Default
```
ObservationTimeout is the timeout for making observations using the DataSource.Observe method.
(This is necessary because an oracle's operations are serialized, so
blocking forever on an observation would break the oracle.)

### BlockchainTimeout<a id='OCR-BlockchainTimeout'></a>
```toml
BlockchainTimeout = '20s' # Default
```
BlockchainTimeout is the timeout for blockchain queries (mediated through
ContractConfigTracker and ContractTransmitter).
(This is necessary because an oracle's operations are serialized, so
blocking forever on a chain interaction would break the oracle.)

### ContractPollInterval<a id='OCR-ContractPollInterval'></a>
```toml
ContractPollInterval = '1m' # Default
```
ContractPollInterval is the polling interval at which ContractConfigTracker is queried for
updated on-chain configurations. Recommended values are between
fifteen seconds and two minutes.

### ContractSubscribeInterval<a id='OCR-ContractSubscribeInterval'></a>
```toml
ContractSubscribeInterval = '2m' # Default
```
ContractSubscribeInterval is the interval at which we try to establish a subscription on ContractConfigTracker
if one doesn't exist. Recommended values are between two and five minutes.

### DefaultTransactionQueueDepth<a id='OCR-DefaultTransactionQueueDepth'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DefaultTransactionQueueDepth = 1 # Default
```
DefaultTransactionQueueDepth controls the queue size for `DropOldestStrategy` in OCR. Set to 0 to use `SendEvery` strategy instead.

### KeyBundleID<a id='OCR-KeyBundleID'></a>
```toml
KeyBundleID = 'acdd42797a8b921b2910497badc5000600000000000000000000000000000000' # Example
```
KeyBundleID is the default key bundle ID to use for OCR jobs. If you have an OCR job that does not explicitly specify a key bundle ID, it will fall back to this value.

### SimulateTransactions<a id='OCR-SimulateTransactions'></a>
```toml
SimulateTransactions = false # Default
```
SimulateTransactions enables transaction simulation for OCR.

### TransmitterAddress<a id='OCR-TransmitterAddress'></a>
```toml
TransmitterAddress = '0xa0788FC17B1dEe36f057c42B6F373A34B014687e' # Example
```
TransmitterAddress is the default sending address to use for OCR. If you have an OCR job that does not explicitly specify a transmitter address, it will fall back to this value.

## P2P<a id='P2P'></a>
```toml
[P2P]
IncomingMessageBufferSize = 10 # Default
OutgoingMessageBufferSize = 10 # Default
TraceLogging = false # Default
```
P2P supports multiple networking stack versions. You may configure `[P2P.V1]`, `[P2P.V2]`, or both to run simultaneously.
If both are configured, then for each link with another peer, V2 networking will be preferred. If V2 does not work, the link will
automatically fall back to V1. If V2 starts working again later, it will automatically be preferred again. This is useful
for migrating networks without downtime. Note that the two networking stacks _must not_ be configured to bind to the same IP/port.

All nodes in the OCR network should share the same networking stack.

### IncomingMessageBufferSize<a id='P2P-IncomingMessageBufferSize'></a>
```toml
IncomingMessageBufferSize = 10 # Default
```
IncomingMessageBufferSize is the per-remote number of incoming
messages to buffer. Any additional messages received on top of those
already in the queue will be dropped.

### OutgoingMessageBufferSize<a id='P2P-OutgoingMessageBufferSize'></a>
```toml
OutgoingMessageBufferSize = 10 # Default
```
OutgoingMessageBufferSize is the per-remote number of outgoing
messages to buffer. Any additional messages send on top of those
already in the queue will displace the oldest.
NOTE: OutgoingMessageBufferSize should be comfortably smaller than remote's
IncomingMessageBufferSize to give the remote enough space to process
them all in case we regained connection and now send a bunch at once

### TraceLogging<a id='P2P-TraceLogging'></a>
```toml
TraceLogging = false # Default
```
TraceLogging enables trace level logging.

## P2P.V1<a id='P2P-V1'></a>
```toml
[P2P.V1]
AnnounceIP = '1.2.3.4' # Example
AnnouncePort = 1337 # Example
BootstrapCheckInterval = '20s' # Default
DefaultBootstrapPeers = ['/dns4/example.com/tcp/1337/p2p/12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U', '/ip4/1.2.3.4/tcp/9999/p2p/12D3KooWLZ9uTC3MrvKfDpGju6RAQubiMDL7CuJcAgDRTYP7fh7R'] # Example
DHTAnnouncementCounterUserPrefix = 0 # Default
DHTLookupInterval = 10 # Default
ListenIP = '0.0.0.0' # Default
ListenPort = 1337 # Example
NewStreamTimeout = '10s' # Default
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw' # Example
PeerstoreWriteInterval = '5m' # Default
```


### AnnounceIP<a id='P2P-V1-AnnounceIP'></a>
```toml
AnnounceIP = '1.2.3.4' # Example
```
AnnounceIP should be set as the externally reachable IP address of the Chainlink node.

### AnnouncePort<a id='P2P-V1-AnnouncePort'></a>
```toml
AnnouncePort = 1337 # Example
```
AnnouncePort should be set as the externally reachable port of the Chainlink node.

### BootstrapCheckInterval<a id='P2P-V1-BootstrapCheckInterval'></a>
```toml
BootstrapCheckInterval = '20s' # Default
```
BootstrapCheckInterval is the interval at which nodes check connections to bootstrap nodes and reconnect if any of them is lost.
Setting this to a small value would allow newly joined bootstrap nodes to get more connectivityBootstrapCheckInterval = '20s' # Default
more quickly, which helps to make bootstrap process faster. The cost of this operation is relatively# DefaultBootstrapPeers is the default set of bootstrap peers.
cheap. We set this to 1 minute during our test.DefaultBootstrapPeers = ['/dns4/example.com/tcp/1337/p2p/12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U', '/ip4/1.2.3.4/tcp/9999/p2p/12D3KooWLZ9uTC3MrvKfDpGju6RAQubiMDL7CuJcAgDRTYP7fh7R'] # Example

### DefaultBootstrapPeers<a id='P2P-V1-DefaultBootstrapPeers'></a>
```toml
DefaultBootstrapPeers = ['/dns4/example.com/tcp/1337/p2p/12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U', '/ip4/1.2.3.4/tcp/9999/p2p/12D3KooWLZ9uTC3MrvKfDpGju6RAQubiMDL7CuJcAgDRTYP7fh7R'] # Example
```
DefaultBootstrapPeers is the default set of bootstrap peers.

### DHTAnnouncementCounterUserPrefix<a id='P2P-V1-DHTAnnouncementCounterUserPrefix'></a>
```toml
DHTAnnouncementCounterUserPrefix = 0 # Default
```
DHTAnnouncementCounterUserPrefix can be used to restore the node's
ability to announce its IP/port on the P2P network after a database
rollback. Make sure to only increase this value, and *never* decrease it.
Don't use this variable unless you really know what you're doing, since you
could semi-permanently exclude your node from the P2P network by
misconfiguring it.

### DHTLookupInterval<a id='P2P-V1-DHTLookupInterval'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DHTLookupInterval = 10 # Default
```
DHTLookupInterval is the interval between which we do the expensive peer
lookup using DHT.

Every DHTLookupInterval failures to open a stream to a peer, we will
attempt to lookup its IP from DHT

### ListenIP<a id='P2P-V1-ListenIP'></a>
```toml
ListenIP = '0.0.0.0' # Default
```
ListenIP is the default IP address to bind to.

### ListenPort<a id='P2P-V1-ListenPort'></a>
```toml
ListenPort = 1337 # Example
```
ListenPort is the port to listen on. If left blank, the node randomly selects a different port each time it boots. It is highly recommended to set this to a static value to avoid network instability.

### NewStreamTimeout<a id='P2P-V1-NewStreamTimeout'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
NewStreamTimeout = '10s' # Default
```
NewStreamTimeout is the maximum length of time to wait to open a
stream before we give up.
We shouldn't hit this in practice since libp2p will give up fast if
it can't get a connection, but it is here anyway as a failsafe.
Set to 0 to disable any timeout on top of what libp2p gives us by default.

### PeerID<a id='P2P-V1-PeerID'></a>
```toml
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw' # Example
```
PeerID is the default peer ID to use for OCR jobs. If unspecified, uses the first available peer ID.

### PeerstoreWriteInterval<a id='P2P-V1-PeerstoreWriteInterval'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
PeerstoreWriteInterval = '5m' # Default
```
PeerstoreWriteInterval controls how often the peerstore for the OCR V1 networking stack is persisted to the database.

## P2P.V2<a id='P2P-V2'></a>
```toml
[P2P.V2]
AnnounceAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
DefaultBootstrappers = ['12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U@1.2.3.4:9999', '12D3KooWM55u5Swtpw9r8aFLQHEtw7HR4t44GdNs654ej5gRs2Dh@example.com:1234'] # Example
DeltaDial = '15s' # Default
DeltaReconcile = '1m' # Default
ListenAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
```


### AnnounceAddresses<a id='P2P-V2-AnnounceAddresses'></a>
```toml
AnnounceAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
```
AnnounceAddresses is the addresses the peer will advertise on the network in host:port form as accepted by net.Dial. The addresses should be reachable by peers of interest.

### DefaultBootstrappers<a id='P2P-V2-DefaultBootstrappers'></a>
```toml
DefaultBootstrappers = ['12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U@1.2.3.4:9999', '12D3KooWM55u5Swtpw9r8aFLQHEtw7HR4t44GdNs654ej5gRs2Dh@example.com:1234'] # Example
```
DefaultBootstrappers is the default bootstrapper peers for libocr's v2 networking stack.

### DeltaDial<a id='P2P-V2-DeltaDial'></a>
```toml
DeltaDial = '15s' # Default
```
DeltaDial controls how far apart Dial attempts are

### DeltaReconcile<a id='P2P-V2-DeltaReconcile'></a>
```toml
DeltaReconcile = '1m' # Default
```
DeltaReconcile controls how often a Reconcile message is sent to every peer.

### ListenAddresses<a id='P2P-V2-ListenAddresses'></a>
```toml
ListenAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
```
ListenAddresses is the addresses the peer will listen to on the network in `host:port` form as accepted by `net.Listen()`, but the host and port must be fully specified and cannot be empty. You can specify `0.0.0.0` (IPv4) or `::` (IPv6) to listen on all interfaces, but that is not recommended.

## Keeper<a id='Keeper'></a>
```toml
[Keeper]
DefaultTransactionQueueDepth = 1 # Default
GasPriceBufferPercent = 20 # Default
GasTipCapBufferPercent = 20 # Default
BaseFeeBufferPercent = 20 # Default
MaximumGracePeriod = 100 # Default
RegistryCheckGasOverhead = '200000' # Default
RegistryPerformGasOverhead = '150000' # Default
RegistrySyncInterval = '30m' # Default
RegistrySyncUpkeepQueueSize = 10 # Default
TurnLookBack = 1000 # Default
TurnFlagEnabled = false # Default
UpkeepCheckGasPriceEnabled = false # Default
```


### DefaultTransactionQueueDepth<a id='Keeper-DefaultTransactionQueueDepth'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DefaultTransactionQueueDepth = 1 # Default
```
DefaultTransactionQueueDepth controls the queue size for `DropOldestStrategy` in Keeper. Set to 0 to use `SendEvery` strategy instead.

### GasPriceBufferPercent<a id='Keeper-GasPriceBufferPercent'></a>
```toml
GasPriceBufferPercent = 20 # Default
```
GasPriceBufferPercent specifies the percentage to add to the gas price used for checking whether to perform an upkeep. Only applies in legacy mode (EIP-1559 off).

### GasTipCapBufferPercent<a id='Keeper-GasTipCapBufferPercent'></a>
```toml
GasTipCapBufferPercent = 20 # Default
```
GasTipCapBufferPercent specifies the percentage to add to the gas price used for checking whether to perform an upkeep. Only applies in EIP-1559 mode.

### BaseFeeBufferPercent<a id='Keeper-BaseFeeBufferPercent'></a>
```toml
BaseFeeBufferPercent = 20 # Default
```
BaseFeeBufferPercent specifies the percentage to add to the base fee used for checking whether to perform an upkeep. Applies only in EIP-1559 mode.

### MaximumGracePeriod<a id='Keeper-MaximumGracePeriod'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
MaximumGracePeriod = 100 # Default
```
MaximumGracePeriod is the maximum number of blocks that a keeper will wait after performing an upkeep before it resumes checking that upkeep

### RegistryCheckGasOverhead<a id='Keeper-RegistryCheckGasOverhead'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
RegistryCheckGasOverhead = '200000' # Default
```
RegistryCheckGasOverhead is the amount of extra gas to provide checkUpkeep() calls to account for the gas consumed by the keeper registry.

### RegistryPerformGasOverhead<a id='Keeper-RegistryPerformGasOverhead'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
RegistryPerformGasOverhead = '150000' # Default
```
RegistryPerformGasOverhead is the amount of extra gas to provide performUpkeep() calls to account for the gas consumed by the keeper registry

### RegistrySyncInterval<a id='Keeper-RegistrySyncInterval'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
RegistrySyncInterval = '30m' # Default
```
RegistrySyncInterval is the interval in which the RegistrySynchronizer performs a full sync of the keeper registry contract it is tracking.

### RegistrySyncUpkeepQueueSize<a id='Keeper-RegistrySyncUpkeepQueueSize'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
RegistrySyncUpkeepQueueSize = 10 # Default
```
RegistrySyncUpkeepQueueSize represents the maximum number of upkeeps that can be synced in parallel.

### TurnLookBack<a id='Keeper-TurnLookBack'></a>
```toml
TurnLookBack = 1000 # Default
```
TurnLookBack is the number of blocks in the past to look back when getting a block for a turn.

### TurnFlagEnabled<a id='Keeper-TurnFlagEnabled'></a>
```toml
TurnFlagEnabled = false # Default
```
TurnFlagEnabled enables a new algorithm for how keepers take turns.

### UpkeepCheckGasPriceEnabled<a id='Keeper-UpkeepCheckGasPriceEnabled'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
UpkeepCheckGasPriceEnabled = false # Default
```
UpkeepCheckGasPriceEnabled includes gas price in calls to `checkUpkeep()` when set to `true`.

## AutoPprof<a id='AutoPprof'></a>
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

### Enabled<a id='AutoPprof-Enabled'></a>
```toml
Enabled = false # Default
```
Enabled enables the automatic profiling service.

### ProfileRoot<a id='AutoPprof-ProfileRoot'></a>
```toml
ProfileRoot = 'prof/root' # Example
```
ProfileRoot sets the location on disk where pprof profiles will be stored. Defaults to `RootDir`.

### PollInterval<a id='AutoPprof-PollInterval'></a>
```toml
PollInterval = '10s' # Default
```
PollInterval is the interval at which the node's resources are checked.

### GatherDuration<a id='AutoPprof-GatherDuration'></a>
```toml
GatherDuration = '10s' # Default
```
GatherDuration is the duration for which profiles are gathered when profiling starts.

### GatherTraceDuration<a id='AutoPprof-GatherTraceDuration'></a>
```toml
GatherTraceDuration = '5s' # Default
```
GatherTraceDuration is the duration for which traces are gathered when profiling is kicked off. This is separately configurable because traces are significantly larger than other types of profiles.

### MaxProfileSize<a id='AutoPprof-MaxProfileSize'></a>
```toml
MaxProfileSize = '100mb' # Default
```
MaxProfileSize is the maximum amount of disk space that profiles may consume before profiling is disabled.

### CPUProfileRate<a id='AutoPprof-CPUProfileRate'></a>
```toml
CPUProfileRate = 1 # Default
```
CPUProfileRate sets the rate for CPU profiling. See https://pkg.go.dev/runtime#SetCPUProfileRate.

### MemProfileRate<a id='AutoPprof-MemProfileRate'></a>
```toml
MemProfileRate = 1 # Default
```
MemProfileRate sets the rate for memory profiling. See https://pkg.go.dev/runtime#pkg-variables.

### BlockProfileRate<a id='AutoPprof-BlockProfileRate'></a>
```toml
BlockProfileRate = 1 # Default
```
BlockProfileRate sets the fraction of blocking events for goroutine profiling. See https://pkg.go.dev/runtime#SetBlockProfileRate.

### MutexProfileFraction<a id='AutoPprof-MutexProfileFraction'></a>
```toml
MutexProfileFraction = 1 # Default
```
MutexProfileFraction sets the fraction of contention events for mutex profiling. See https://pkg.go.dev/runtime#SetMutexProfileFraction.

### MemThreshold<a id='AutoPprof-MemThreshold'></a>
```toml
MemThreshold = '4gb' # Default
```
MemThreshold sets the maximum amount of memory the node can actively consume before profiling begins.

### GoroutineThreshold<a id='AutoPprof-GoroutineThreshold'></a>
```toml
GoroutineThreshold = 5000 # Default
```
GoroutineThreshold is the maximum number of actively-running goroutines the node can spawn before profiling begins.

## Sentry<a id='Sentry'></a>
```toml
[Sentry]
Debug = false # Default
DSN = 'sentry-dsn' # Example
Environment = 'prod' # Default
Release = 'v1.2.3' # Example
```


### Debug<a id='Sentry-Debug'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
Debug = false # Default
```
Debug enables printing of Sentry SDK debug messages.

### DSN<a id='Sentry-DSN'></a>
```toml
DSN = 'sentry-dsn' # Example
```
DSN is the data source name where events will be sent. Sentry is completely disabled if this is left blank.

### Environment<a id='Sentry-Environment'></a>
```toml
Environment = 'prod' # Default
```
Environment overrides the Sentry environment to the given value. Otherwise autodetects between dev/prod.

### Release<a id='Sentry-Release'></a>
```toml
Release = 'v1.2.3' # Example
```
Release overrides the Sentry release to the given value. Otherwise uses the compiled-in version number.

## EVM<a id='EVM'></a>
EVM defaults depend on ChainID:

<details><summary>Ethereum Mainnet (1)<a id='EVM-1'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x514910771AF9Ca656af840dff83E8264EcF986CA'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.1 link'
NonceAutoSync = true
OperatorFactoryAddress = '0x3E64Cd889482443324F91bFA9c84fE72A511f48A'
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 4
TransactionPercentile = 50


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Ethereum Ropsten (3)<a id='EVM-3'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x20fE562d797A42Dcb3399062AE9546cd06f63280'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.1 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 4
TransactionPercentile = 50


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Ethereum Rinkeby (4)<a id='EVM-4'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x01BE23585060835E02B77ef475b0Cc51aA1e0709'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.1 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 4
TransactionPercentile = 50


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Ethereum Goerli (5)<a id='EVM-5'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x326C977E6efc84E512bB9C30f76E30c160eD06FB'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.1 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 4
TransactionPercentile = 50


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Optimism Mainnet (10)<a id='EVM-10'></a></summary><p>

```toml
ChainType = 'optimism'
FinalityDepth = 1
LinkContractAddress = '0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '15s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 0

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '0'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 0
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '0s'
HistoryDepth = 10
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>RSK Mainnet (30)<a id='EVM-30'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x14AdaE34beF7ca957Ce2dDe5ADD97ea050123827'
LogBackfillBatchSize = 100
LogPollInterval = '30s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '50 mwei'
PriceMax = '50 gwei'
PriceMin = '0'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 mwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>RSK Testnet (31)<a id='EVM-31'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x8bBbd80981FE76d44854D8DF305e8985c19f0e78'
LogBackfillBatchSize = 100
LogPollInterval = '30s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '50 mwei'
PriceMax = '50 gwei'
PriceMin = '0'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 mwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Ethereum Kovan (42)<a id='EVM-42'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0xa36085F69e2889c224210F603D836748e7dC0088'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.1 link'
NonceAutoSync = true
OperatorFactoryAddress = '0x8007e24251b1D2Fc518Eb843A701d9cD21fe0aA3'
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 4
TransactionPercentile = 50


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>BSC Mainnet (56)<a id='EVM-56'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x404460C6A5EdE2D891e8297795264fDe62ADBB75'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 2

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 2
BlockHistorySize = 24
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '2s'
DatabaseTimeout = '2s'
ObservationGracePeriod = '500ms'
```

</p></details>

<details><summary>OKX Testnet (65)<a id='EVM-65'></a></summary><p>

```toml
FinalityDepth = 50
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>OKX Mainnet (66)<a id='EVM-66'></a></summary><p>

```toml
FinalityDepth = 50
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Optimism Kovan (69)<a id='EVM-69'></a></summary><p>

```toml
ChainType = 'optimism'
FinalityDepth = 1
LinkContractAddress = '0x4911b761993b9c8c0d14Ba2d86902AF6B0074F5B'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '15s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 0

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '0'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 0
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '30m0s'
HistoryDepth = 10
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>xDai Mainnet (100)<a id='EVM-100'></a></summary><p>

```toml
ChainType = 'xdai'
FinalityDepth = 50
LinkContractAddress = '0xE2e73A1c69ecF83F464EFCE6A5be353a37cA09b2'
LogBackfillBatchSize = 100
LogPollInterval = '5s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '1 gwei'
PriceMax = '500 gwei'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '1m0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '3m0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Heco Mainnet (128)<a id='EVM-128'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x404460C6A5EdE2D891e8297795264fDe62ADBB75'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 2

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 2
BlockHistorySize = 24
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '2s'
DatabaseTimeout = '2s'
ObservationGracePeriod = '500ms'
```

</p></details>

<details><summary>Polygon Mainnet (137)<a id='EVM-137'></a></summary><p>

```toml
FinalityDepth = 500
LinkContractAddress = '0xb0897686c545045aFc77CF20eC7A532E3120E0F1'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 5000
MinIncomingConfirmations = 5
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 13

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '30 gwei'
PriceMax = '200 micro'
PriceMin = '30 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '20 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 10
BlockHistorySize = 24
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Fantom Mainnet (250)<a id='EVM-250'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x6F43FF82CCA38001B6699a8AC47A2d0E66939407'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 2

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '15 gwei'
PriceMax = '200 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 2
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Metis Rinkeby (588)<a id='EVM-588'></a></summary><p>

```toml
ChainType = 'metis'
FinalityDepth = 1
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 0

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '0'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 0
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Metis Mainnet (1088)<a id='EVM-1088'></a></summary><p>

```toml
ChainType = 'metis'
FinalityDepth = 1
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 0

[GasEstimator]
Mode = 'L2Suggested'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '0'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 0
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Fantom Testnet (4002)<a id='EVM-4002'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0xfaFedb041c0DD4fA2Dc0d87a6B0979Ee6FA7af5F'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 2

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '15 gwei'
PriceMax = '200 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 2
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Arbitrum Mainnet (42161)<a id='EVM-42161'></a></summary><p>

```toml
ChainType = 'arbitrum'
FinalityDepth = 50
LinkContractAddress = '0xf97f4df75117a78c1A5a0DBb814Af92458539FB4'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'FixedPrice'
PriceDefault = '1 micro'
PriceMax = '1 micro'
PriceMin = '1 micro'
LimitDefault = 7000000
LimitMultiplier = '1'
LimitTransfer = 800000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 0
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Avalanche Fuji (43113)<a id='EVM-43113'></a></summary><p>

```toml
FinalityDepth = 1
LinkContractAddress = '0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '25 gwei'
PriceMax = '1 micro'
PriceMin = '25 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 2
BlockHistorySize = 24
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Avalanche Mainnet (43114)<a id='EVM-43114'></a></summary><p>

```toml
FinalityDepth = 1
LinkContractAddress = '0x5947BB275c521040051D82396192181b413227A3'
LogBackfillBatchSize = 100
LogPollInterval = '3s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '25 gwei'
PriceMax = '1 micro'
PriceMin = '25 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 2
BlockHistorySize = 24
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Polygon Mumbai (80001)<a id='EVM-80001'></a></summary><p>

```toml
FinalityDepth = 500
LinkContractAddress = '0x326C977E6efc84E512bB9C30f76E30c160eD06FB'
LogBackfillBatchSize = 100
LogPollInterval = '1s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 5000
MinIncomingConfirmations = 5
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 13

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '1 gwei'
PriceMax = '200 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '20 gwei'
BumpPercent = 20
BumpThreshold = 5
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 10
BlockHistorySize = 24
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Arbitrum Rinkeby (421611)<a id='EVM-421611'></a></summary><p>

```toml
ChainType = 'arbitrum'
FinalityDepth = 50
LinkContractAddress = '0x615fBe6372676474d9e6933d310469c9b68e9726'
LogBackfillBatchSize = 100
LogPollInterval = '15s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 3
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'FixedPrice'
PriceDefault = '1 micro'
PriceMax = '1 micro'
PriceMin = '1 micro'
LimitDefault = 7000000
LimitMultiplier = '1'
LimitTransfer = 800000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 0
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '0s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '0s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Harmony Mainnet (1666600000)<a id='EVM-1666600000'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x218532a12a389a4a92fC0C5Fb22901D1c19198aA'
LogBackfillBatchSize = 100
LogPollInterval = '2s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>

<details><summary>Harmony Testnet (1666700000)<a id='EVM-1666700000'></a></summary><p>

```toml
FinalityDepth = 50
LinkContractAddress = '0x8b12Ac23BFe11cAb03a634C1F117D64a7f2cFD3e'
LogBackfillBatchSize = 100
LogPollInterval = '2s'
MaxInFlightTransactions = 16
MaxQueuedTransactions = 250
MinIncomingConfirmations = 1
MinimumContractPayment = '0.00001 link'
NonceAutoSync = true
RPCDefaultBatchSize = 100
TxReaperInterval = '1h0m0s'
TxReaperThreshold = '168h0m0s'
TxResendAfterThreshold = '1m0s'
UseForwarders = false

[BalanceMonitor]
Enabled = true
BlockDelay = 1

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
BumpTxDepth = 10
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMinimum = '1 wei'
[GasEstimator.BlockHistory]
BatchSize = 4
BlockDelay = 1
BlockHistorySize = 8
TransactionPercentile = 60


[HeadTracker]
BlockEmissionIdleWarningThreshold = '15s'
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'

[NodePool]
NoNewHeadsThreshold = '30s'
PollFailureThreshold = 5
PollInterval = '10s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
ObservationGracePeriod = '1s'
```

</p></details>


### ChainID<a id='EVM-ChainID'></a>
```toml
ChainID = '1' # Example
```
ChainID is the EVM chain ID. Mandatory.

### Enabled<a id='EVM-Enabled'></a>
```toml
Enabled = true # Default
```
Enabled enables this chain.

### BlockBackfillDepth<a id='EVM-BlockBackfillDepth'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
BlockBackfillDepth = 10 # Default
```
BlockBackfillDepth specifies the number of blocks before the current HEAD that the log broadcaster will try to re-consume logs from.

### BlockBackfillSkip<a id='EVM-BlockBackfillSkip'></a>
```toml
BlockBackfillSkip = false # Default
```
BlockBackfillSkip enables skipping of very long backfills.

### ChainType<a id='EVM-ChainType'></a>
```toml
ChainType = 'Optimism' # Example
```
ChainType is automatically detected from chain ID. Set this to force a certain chain type regardless of chain ID.

### FinalityDepth<a id='EVM-FinalityDepth'></a>
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

### FlagsContractAddress<a id='EVM-FlagsContractAddress'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
FlagsContractAddress = '0xae4E781a6218A8031764928E88d457937A954fC3' # Example
```
FlagsContractAddress can optionally point to a [Flags contract](../contracts/src/v0.8/Flags.sol). If set, the node will lookup that contract for each job that supports flags contracts (currently OCR and FM jobs are supported). If the job's contractAddress is set as hibernating in the FlagsContractAddress address, it overrides the standard update parameters (such as heartbeat/threshold).

### LinkContractAddress<a id='EVM-LinkContractAddress'></a>
```toml
LinkContractAddress = '0x538aAaB4ea120b2bC2fe5D296852D948F07D849e' # Example
```
LinkContractAddress is the canonical ERC-677 LINK token contract address on the given chain. Note that this is usually autodetected from chain ID.

### LogBackfillBatchSize<a id='EVM-LogBackfillBatchSize'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
LogBackfillBatchSize = 100 # Default
```
LogBackfillBatchSize sets the batch size for calling FilterLogs when we backfill missing logs.

### LogPollInterval<a id='EVM-LogPollInterval'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
LogPollInterval = '15s' # Default
```
LogPollInterval works in conjunction with Feature.LogPoller. Controls how frequently the log poller polls for logs. Defaults to the block production rate.

### MaxInFlightTransactions<a id='EVM-MaxInFlightTransactions'></a>
```toml
MaxInFlightTransactions = 16 # Default
```
MaxInFlightTransactions controls how many transactions are allowed to be "in-flight" i.e. broadcast but unconfirmed at any one time. You can consider this a form of transaction throttling.

The default is set conservatively at 16 because this is a pessimistic minimum that both geth and parity will hold without evicting local transactions. If your node is falling behind and you need higher throughput, you can increase this setting, but you MUST make sure that your ETH node is configured properly otherwise you can get nonce gapped and your node will get stuck.

0 value disables the limit. Use with caution.

### MaxQueuedTransactions<a id='EVM-MaxQueuedTransactions'></a>
```toml
MaxQueuedTransactions = 250 # Default
```
MaxQueuedTransactions is the maximum number of unbroadcast transactions per key that are allowed to be enqueued before jobs will start failing and rejecting send of any further transactions. This represents a sanity limit and generally indicates a problem with your ETH node (transactions are not getting mined).

Do NOT blindly increase this value thinking it will fix things if you start hitting this limit because transactions are not getting mined, you will instead only make things worse.

In deployments with very high burst rates, or on chains with large re-orgs, you _may_ consider increasing this.

0 value disables any limit on queue size. Use with caution.

### MinIncomingConfirmations<a id='EVM-MinIncomingConfirmations'></a>
```toml
MinIncomingConfirmations = 3 # Default
```
MinIncomingConfirmations is the minimum required confirmations before a log event will be consumed.

### MinimumContractPayment<a id='EVM-MinimumContractPayment'></a>
```toml
MinimumContractPayment = '10000000000000 juels' # Default
```
MinimumContractPayment is the minimum payment in LINK required to execute a direct request job. This can be overridden on a per-job basis.

### NonceAutoSync<a id='EVM-NonceAutoSync'></a>
```toml
NonceAutoSync = true # Default
```
NonceAutoSync enables automatic nonce syncing on startup. Chainlink nodes will automatically try to sync its local nonce with the remote chain on startup and fast forward if necessary. This is almost always safe but can be disabled in exceptional cases by setting this value to false.

### OperatorFactoryAddress<a id='EVM-OperatorFactoryAddress'></a>
```toml
OperatorFactoryAddress = '0xa5B85635Be42F21f94F28034B7DA440EeFF0F418' # Example
```
OperatorFactoryAddress is the address of the canonical operator forwarder contract on the given chain. Note that this is usually autodetected from chain ID.

### RPCDefaultBatchSize<a id='EVM-RPCDefaultBatchSize'></a>
```toml
RPCDefaultBatchSize = 100 # Default
```
RPCDefaultBatchSize is the default batch size for batched RPC calls.

### TxReaperInterval<a id='EVM-TxReaperInterval'></a>
```toml
TxReaperInterval = '1h' # Default
```
TxReaperInterval controls how often the EthTx reaper will run.

### TxReaperThreshold<a id='EVM-TxReaperThreshold'></a>
```toml
TxReaperThreshold = '168h' # Default
```
TxReaperThreshold indicates how old an EthTx ought to be before it can be reaped.

### TxResendAfterThreshold<a id='EVM-TxResendAfterThreshold'></a>
```toml
TxResendAfterThreshold = '1m' # Default
```
TxResendAfterThreshold controls how long to wait before re-broadcasting a transaction that has not yet been confirmed.

### UseForwarders<a id='EVM-UseForwarders'></a>
```toml
UseForwarders = false # Default
```
UseForwarders enables or disables sending transactions through forwarder contracts.

## EVM.BalanceMonitor<a id='EVM-BalanceMonitor'></a>
```toml
[EVM.BalanceMonitor]
Enabled = true # Default
BlockDelay = 1 # Default
```


### Enabled<a id='EVM-BalanceMonitor-Enabled'></a>
```toml
Enabled = true # Default
```
Enabled balance monitoring for all keys.

### BlockDelay<a id='EVM-BalanceMonitor-BlockDelay'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
BlockDelay = 1 # Default
```
BlockDelay is the number of blocks that the balance monitor trails behind head. This is required when load balancing
across multiple nodes announce a new head, then route a request to a different node which does not have this head yet.

## EVM.GasEstimator<a id='EVM-GasEstimator'></a>
```toml
[EVM.GasEstimator]
Mode = 'BlockHistory' # Default
PriceDefault = '20 gwei' # Default
PriceMax = '100 micro' # Default
PriceMin = '1 gwei' # Default
LimitDefault = 500_000 # Default
LimitOCRJobType = 100_000 # Example
LimitDRJobType = 100_000 # Example
LimitVRFJobType = 100_000 # Example
LimitFMJobType = 100_000 # Example
LimitKeeperJobType = 100_000 # Example
LimitMultiplier = '1.0' # Default
LimitTransfer = 21_000 # Default
BumpMin = '5 gwei' # Default
BumpPercent = 20 # Default
BumpThreshold = 3 # Default
BumpTxDepth = 10 # Default
EIP1559DynamicFees = false # Default
FeeCapDefault = '100 gwei' # Default
TipCapDefault = '1 wei' # Default
TipCapMinimum = '1 wei' # Default
```


### Mode<a id='EVM-GasEstimator-Mode'></a>
```toml
Mode = 'BlockHistory' # Default
```
Mode controls what type of gas estimator is used.

- `FixedPrice` uses static configured values for gas price (can be set via API call).
- `BlockHistory` dynamically adjusts default gas price based on heuristics from mined blocks.
- `L2Suggested`

Chainlink nodes decide what gas price to use using an `Estimator`. It ships with several simple and battle-hardened built-in estimators that should work well for almost all use-cases. Note that estimators will change their behaviour slightly depending on if you are in EIP-1559 mode or not.

You can also use your own estimator for gas price by selecting the `FixedPrice` estimator and using the exposed API to set the price.

An important point to note is that the Chainlink node does _not_ ship with built-in support for go-ethereum's `estimateGas` call. This is for several reasons, including security and reliability. We have found empirically that it is not generally safe to rely on the remote ETH node's idea of what gas price should be.

### PriceDefault<a id='EVM-GasEstimator-PriceDefault'></a>
```toml
PriceDefault = '20 gwei' # Default
```
PriceDefault is the default gas price to use when submitting transactions to the blockchain. Will be overridden by the built-in `BlockHistoryEstimator` if enabled, and might be increased if gas bumping is enabled.

(Only applies to legacy transactions)

Can be used with the `chainlink setgasprice` to be updated while the node is still running.

### PriceMax<a id='EVM-GasEstimator-PriceMax'></a>
```toml
PriceMax = '100 micro' # Default
```
PriceMax is the maximum gas price. Chainlink nodes will never pay more than this for a transaction.

### PriceMin<a id='EVM-GasEstimator-PriceMin'></a>
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

### LimitDefault<a id='EVM-GasEstimator-LimitDefault'></a>
```toml
LimitDefault = 500_000 # Default
```
LimitDefault sets default gas limit for outgoing transactions. This should not need to be changed in most cases.
Some job types, such as Keeper jobs, might set their own gas limit unrelated to this value.

### LimitOCRJobType<a id='EVM-GasEstimator-LimitOCRJobType'></a>
```toml
LimitOCRJobType = 100_000 # Example
```
LimitOCRJobType overrides LimitDefault for OCR jobs.

### LimitDRJobType<a id='EVM-GasEstimator-LimitDRJobType'></a>
```toml
LimitDRJobType = 100_000 # Example
```
LimitDRJobType overrides LimitDefault for Direct Request jobs.

### LimitVRFJobType<a id='EVM-GasEstimator-LimitVRFJobType'></a>
```toml
LimitVRFJobType = 100_000 # Example
```
LimitVRFJobType overrides LimitDefault for VRF jobs.

### LimitFMJobType<a id='EVM-GasEstimator-LimitFMJobType'></a>
```toml
LimitFMJobType = 100_000 # Example
```
LimitFMJobType overrides LimitDefault for Flux Monitor jobs.

### LimitKeeperJobType<a id='EVM-GasEstimator-LimitKeeperJobType'></a>
```toml
LimitKeeperJobType = 100_000 # Example
```
LimitKeeperJobType overrides LimitDefault for Keeper jobs.

### LimitMultiplier<a id='EVM-GasEstimator-LimitMultiplier'></a>
```toml
LimitMultiplier = '1.0' # Default
```
LimitMultiplier is the factor by which a transaction's GasLimit is multiplied before transmission. So if the value is 1.1, and the GasLimit for a transaction is 10, 10% will be added before transmission.

This factor is always applied, so includes Optimism L2 transactions which uses a default gas limit of 1 and is also applied to `LimitDefault`.

### LimitTransfer<a id='EVM-GasEstimator-LimitTransfer'></a>
```toml
LimitTransfer = 21_000 # Default
```
LimitTransfer is the gas limit used for an ordinary ETH transfer.

### BumpMin<a id='EVM-GasEstimator-BumpMin'></a>
```toml
BumpMin = '5 gwei' # Default
```
BumpMin is the minimum fixed amount of wei by which gas is bumped on each transaction attempt.

### BumpPercent<a id='EVM-GasEstimator-BumpPercent'></a>
```toml
BumpPercent = 20 # Default
```
BumpPercent is the percentage by which to bump gas on a transaction that has exceeded `BumpThreshold`. The larger of `GasBumpPercent` and `GasBumpWei` is taken for gas bumps.

### BumpThreshold<a id='EVM-GasEstimator-BumpThreshold'></a>
```toml
BumpThreshold = 3 # Default
```
BumpThreshold is the number of blocks to wait for a transaction stuck in the mempool before automatically bumping the gas price. Set to 0 to disable gas bumping completely.

### BumpTxDepth<a id='EVM-GasEstimator-BumpTxDepth'></a>
```toml
BumpTxDepth = 10 # Default
```
BumpTxDepth is the number of transactions to gas bump starting from oldest. Set to 0 for no limit (i.e. bump all).

### EIP1559DynamicFees<a id='EVM-GasEstimator-EIP1559DynamicFees'></a>
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
- With gas bumping disabled, it will submit all transactions with `feecap=MaxGasPriceWei` and `tipcap=GasTipCapDefault`
- With gas bumping enabled, it will submit all transactions initially with `feecap=GasFeeCapDefault` and `tipcap=GasTipCapDefault`.

If you are using BlockHistoryEstimator (default for most chains):
- With gas bumping disabled, it will submit all transactions with `feecap=MaxGasPriceWei` and `tipcap=<calculated using past blocks>`
- With gas bumping enabled (default for most chains) it will submit all transactions initially with `feecap=current block base fee * (1.125 ^ N)` where N is configurable by setting `EVM.GasEstimator.BlockHistory.EIP1559FeeCapBufferBlocks` but defaults to `gas bump threshold+1` and `tipcap=<calculated using past blocks>`

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
- `PriceMinWei` is ignored for new transactions and `GasTipCapMinimum` is used instead (default 0)
- `PriceMaxWei` still represents that absolute upper limit that Chainlink will ever spend (total) on a single tx
- `Keeper.GasTipCapBufferPercent` is ignored in EIP-1559 mode and `Keeper.GasTipCapBufferPercent` is used instead

### FeeCapDefault<a id='EVM-GasEstimator-FeeCapDefault'></a>
```toml
FeeCapDefault = '100 gwei' # Default
```
FeeCapDefault controls the fixed initial fee cap, if EIP1559 mode is enabled and `FixedPrice` gas estimator is used.

### TipCapDefault<a id='EVM-GasEstimator-TipCapDefault'></a>
```toml
TipCapDefault = '1 wei' # Default
```
TipCapDefault is the default gas tip to use when submitting transactions to the blockchain. Will be overridden by the built-in `BlockHistoryEstimator` if enabled, and might be increased if gas bumping is enabled.

(Only applies to EIP-1559 transactions)

### TipCapMinimum<a id='EVM-GasEstimator-TipCapMinimum'></a>
```toml
TipCapMinimum = '1 wei' # Default
```
TipCapMinimum is the minimum gas tip to use when submitting transactions to the blockchain.

Only applies to EIP-1559 transactions)

## EVM.GasEstimator.BlockHistory<a id='EVM-GasEstimator-BlockHistory'></a>
```toml
[EVM.GasEstimator.BlockHistory]
BatchSize = 4 # Default
BlockDelay = 1 # Default
BlockHistorySize = 8 # Default
EIP1559FeeCapBufferBlocks = 13 # Example
TransactionPercentile = 60 # Default
```
These settings allow you to configure how your node calculates gas prices when using the block history estimator.
In most cases, leaving these values at their defaults should give good results.

### BatchSize<a id='EVM-GasEstimator-BlockHistory-BatchSize'></a>
```toml
BatchSize = 4 # Default
```
BatchSize sets the maximum number of blocks to fetch in one batch in the block history estimator.
If the `BatchSize` variable is set to 0, it defaults to `EVM.RPCDefaultBatchSize`.

### BlockDelay<a id='EVM-GasEstimator-BlockHistory-BlockDelay'></a>
```toml
BlockDelay = 1 # Default
```
BlockDelay controls the number of blocks that the block history estimator trails behind head.
For example, if this is set to 3, and we receive block 10, block history estimator will fetch block 7.

CAUTION: You might be tempted to set this to 0 to use the latest possible
block, but it is possible to receive a head BEFORE that block is actually
available from the connected node via RPC, due to race conditions in the code of the remote ETH node. In this case you will get false
"zero" blocks that are missing transactions.

### BlockHistorySize<a id='EVM-GasEstimator-BlockHistory-BlockHistorySize'></a>
```toml
BlockHistorySize = 8 # Default
```
BlockHistorySize controls the number of past blocks to keep in memory to use as a basis for calculating a percentile gas price.

### EIP1559FeeCapBufferBlocks<a id='EVM-GasEstimator-BlockHistory-EIP1559FeeCapBufferBlocks'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
EIP1559FeeCapBufferBlocks = 13 # Example
```
EIP1559FeeCapBufferBlocks controls the buffer blocks to add to the current base fee when sending a transaction. By default, the gas bumping threshold + 1 block is used.

Only applies to EIP-1559 transactions)

### TransactionPercentile<a id='EVM-GasEstimator-BlockHistory-TransactionPercentile'></a>
```toml
TransactionPercentile = 60 # Default
```
TransactionPercentile specifies gas price to choose. E.g. if the block history contains four transactions with gas prices `[100, 200, 300, 400]` then picking 25 for this number will give a value of 200. If the calculated gas price is higher than `GasPriceDefault` then the higher price will be used as the base price for new transactions.

Must be in range 0-100.

Only has an effect if gas updater is enabled.

Think of this number as an indicator of how aggressive you want your node to price its transactions.

Setting this number higher will cause the Chainlink node to select higher gas prices.

Setting it lower will tend to set lower gas prices.

## EVM.HeadTracker<a id='EVM-HeadTracker'></a>
```toml
[EVM.HeadTracker]
BlockEmissionIdleWarningThreshold = '1m' # Default
HistoryDepth = 100 # Default
MaxBufferSize = 3 # Default
SamplingInterval = '1s' # Default
```


### BlockEmissionIdleWarningThreshold<a id='EVM-HeadTracker-BlockEmissionIdleWarningThreshold'></a>
```toml
BlockEmissionIdleWarningThreshold = '1m' # Default
```
BlockEmissionIdleWarningThreshold will cause Chainlink to log warnings if this duration is exceeded without any new blocks being emitted.

### HistoryDepth<a id='EVM-HeadTracker-HistoryDepth'></a>
```toml
HistoryDepth = 100 # Default
```
HistoryDepth tracks the top N block numbers to keep in the `heads` database table.
Note that this can easily result in MORE than N records since in the case of re-orgs we keep multiple heads for a particular block height.
This number should be at least as large as `FinalityDepth`.
There may be a small performance penalty to setting this to something very large (10,000+)

### MaxBufferSize<a id='EVM-HeadTracker-MaxBufferSize'></a>
```toml
MaxBufferSize = 3 # Default
```
MaxBufferSize is the maximum number of heads that may be
buffered in front of the head tracker before older heads start to be
dropped. You may think of it as something like the maximum permittable "lag"
for the head tracker before we start dropping heads to keep up.

### SamplingInterval<a id='EVM-HeadTracker-SamplingInterval'></a>
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
SamplingInterval = '1s' # Default
```
SamplingInterval means that head tracker callbacks will at maximum be made once in every window of this duration. This is a performance optimisation for fast chains. Set to 0 to disable sampling entirely.

## EVM.KeySpecific<a id='EVM-KeySpecific'></a>
```toml
[[EVM.KeySpecific]]
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
GasEstimator.PriceMax = '79 gwei' # Example
```


### Key<a id='EVM-KeySpecific-Key'></a>
```toml
Key = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
```
Key is the account to apply these settings to

### PriceMax<a id='EVM-KeySpecific-GasEstimator-PriceMax'></a>
```toml
GasEstimator.PriceMax = '79 gwei' # Example
```
GasEstimator.PriceMax overrides the maximum gas price for this key. See EVM.GasEstimator.PriceMaxWei.

## EVM.NodePool<a id='EVM-NodePool'></a>
```toml
[EVM.NodePool]
NoNewHeadsThreshold = '3m' # Default
PollFailureThreshold = 3 # Default
PollInterval = '10s' # Default
```


### NoNewHeadsThreshold<a id='EVM-NodePool-NoNewHeadsThreshold'></a>
```toml
NoNewHeadsThreshold = '3m' # Default
```
NoNewHeadsThreshold controls how long to wait after receiving no new heads before marking the node as out-of-sync.

Set to zero to disable out-of-sync checking.

### PollFailureThreshold<a id='EVM-NodePool-PollFailureThreshold'></a>
```toml
PollFailureThreshold = 3 # Default
```
PollFailureThreshold indicates how many consecutive polls must fail in order to mark a node as unreachable.

Set to zero to disable poll checking.

### PollInterval<a id='EVM-NodePool-PollInterval'></a>
```toml
PollInterval = '10s' # Default
```
PollInterval controls how often to poll the node to check for liveness.

Set to zero to disable poll checking.

## EVM.OCR<a id='EVM-OCR'></a>
```toml
[EVM.OCR]
ContractConfirmations = 4 # Default
ContractTransmitterTransmitTimeout = '10s' # Default
DatabaseTimeout = '10s' # Default
ObservationGracePeriod = '1s' # Default
ObservationTimeout = '1m' # Example
```


### ContractConfirmations<a id='EVM-OCR-ContractConfirmations'></a>
```toml
ContractConfirmations = 4 # Default
```
ContractConfirmations sets `OCR.ContractConfirmations` for this EVM chain.

### ContractTransmitterTransmitTimeout<a id='EVM-OCR-ContractTransmitterTransmitTimeout'></a>
```toml
ContractTransmitterTransmitTimeout = '10s' # Default
```
ContractTransmitterTransmitTimeout sets `OCR.ContractTransmitterTransmitTimeout` for this EVM chain.

### DatabaseTimeout<a id='EVM-OCR-DatabaseTimeout'></a>
```toml
DatabaseTimeout = '10s' # Default
```
DatabaseTimeout sets `OCR.DatabaseTimeout` for this EVM chain.

### ObservationGracePeriod<a id='EVM-OCR-ObservationGracePeriod'></a>
```toml
ObservationGracePeriod = '1s' # Default
```
ObservationGracePeriod sets `OCR.ObservationGracePeriod` for this EVM chain.

### ObservationTimeout<a id='EVM-OCR-ObservationTimeout'></a>
```toml
ObservationTimeout = '1m' # Example
```
ObservationTimeout sets `OCR.ObservationTimeout` for this EVM chain.

## EVM.Nodes<a id='EVM-Nodes'></a>
```toml
[[EVM.Nodes]]
Name = 'foo' # Example
WSURL = 'wss://web.socket/test' # Example
HTTPURL = 'https://foo.web' # Example
SendOnly = false # Default
```


### Name<a id='EVM-Nodes-Name'></a>
```toml
Name = 'foo' # Example
```
Name is a unique (per-chain) identifier for this node.

### WSURL<a id='EVM-Nodes-WSURL'></a>
```toml
WSURL = 'wss://web.socket/test' # Example
```
WSURL is the WS(S) endpoint for this node. Required for primary nodes.

### HTTPURL<a id='EVM-Nodes-HTTPURL'></a>
```toml
HTTPURL = 'https://foo.web' # Example
```
HTTPURL is the HTTP(S) endpoint for this node. Recommended for primary nodes. Required for `SendOnly`.

### SendOnly<a id='EVM-Nodes-SendOnly'></a>
```toml
SendOnly = false # Default
```
SendOnly limits usage to sending transaction broadcasts only. With this enabled, only HTTPURL is required, and WSURL is not used.

## Solana<a id='Solana'></a>
```toml
[[Solana]]
ChainID = 'mainnet' # Example
Enabled = false # Default
BalancePollPeriod = '5s' # Default
ConfirmPollPeriod = '500ms' # Default
OCR2CachePollPeriod = '1s' # Default
OCR2CacheTTL = '1m' # Default
TxTimeout = '1h' # Default
TxRetryTimeout = '10s' # Default
TxConfirmTimeout = '30s' # Default
SkipPreflight = true # Default
Commitment = 'confirmed' # Default
MaxRetries = 0 # Default
```


### ChainID<a id='Solana-ChainID'></a>
```toml
ChainID = 'mainnet' # Example
```
ChainID is the Solana chain ID. Must be one of: mainnet, testnet, devnet, localnet. Mandatory.

### Enabled<a id='Solana-Enabled'></a>
```toml
Enabled = false # Default
```
Enabled enables this chain.

### BalancePollPeriod<a id='Solana-BalancePollPeriod'></a>
```toml
BalancePollPeriod = '5s' # Default
```
BalancePollPeriod is the rate to poll for SOL balance and update Prometheus metrics.

### ConfirmPollPeriod<a id='Solana-ConfirmPollPeriod'></a>
```toml
ConfirmPollPeriod = '500ms' # Default
```
ConfirmPollPeriod is the rate to poll for signature confirmation.

### OCR2CachePollPeriod<a id='Solana-OCR2CachePollPeriod'></a>
```toml
OCR2CachePollPeriod = '1s' # Default
```
OCR2CachePollPeriod is the rate to poll for the OCR2 state cache.

### OCR2CacheTTL<a id='Solana-OCR2CacheTTL'></a>
```toml
OCR2CacheTTL = '1m' # Default
```
OCR2CacheTTL is the stale OCR2 cache deadline.

### TxTimeout<a id='Solana-TxTimeout'></a>
```toml
TxTimeout = '1h' # Default
```
TxTimeout is the timeout for sending txes to an RPC endpoint.

### TxRetryTimeout<a id='Solana-TxRetryTimeout'></a>
```toml
TxRetryTimeout = '10s' # Default
```
TxRetryTimeout is the duration for tx manager to attempt rebroadcasting to RPC, before giving up.

### TxConfirmTimeout<a id='Solana-TxConfirmTimeout'></a>
```toml
TxConfirmTimeout = '30s' # Default
```
TxConfirmTimeout is the duration to wait when confirming a tx signature, before discarding as unconfirmed.

### SkipPreflight<a id='Solana-SkipPreflight'></a>
```toml
SkipPreflight = true # Default
```
SkipPreflight enables or disables preflight checks when sending txs.

### Commitment<a id='Solana-Commitment'></a>
```toml
Commitment = 'confirmed' # Default
```
Commitment is the confirmation level for solana state and transactions. ([documentation](https://docs.solana.com/developing/clients/jsonrpc-api#configuring-state-commitment))

### MaxRetries<a id='Solana-MaxRetries'></a>
```toml
MaxRetries = 0 # Default
```
MaxRetries is the maximum number of times the RPC node will automatically rebroadcast a tx.
The default is 0 for custom txm rebroadcasting method, set to -1 to use the RPC node's default retry strategy.

## Solana.Nodes<a id='Solana-Nodes'></a>
```toml
[[Solana.Nodes]]
Name = 'primary' # Example
URL = 'http://solana.web' # Example
```


### Name<a id='Solana-Nodes-Name'></a>
```toml
Name = 'primary' # Example
```
Name is a unique (per-chain) identifier for this node.

### URL<a id='Solana-Nodes-URL'></a>
```toml
URL = 'http://solana.web' # Example
```
URL is the HTTP(S) endpoint for this node.

## Terra<a id='Terra'></a>
```toml
[[Terra]]
ChainID = 'Bombay-12' # Example
Enabled = true # Default
BlockRate = '6s' # Default
BlocksUntilTxTimeout = 30 # Default
ConfirmPollPeriod = '1s' # Default
FallbackGasPriceULuna = '0.015' # Default
FCDURL = 'http://terra.com' # Example
GasLimitMultiplier = '1.5' # Default
MaxMsgsPerBatch = 100 # Default
OCR2CachePollPeriod = '4s' # Default
OCR2CacheTTL = '1m' # Default
TxMsgTimeout = '10m' # Default
```


### ChainID<a id='Terra-ChainID'></a>
```toml
ChainID = 'Bombay-12' # Example
```
ChainID is the Terra chain ID. Mandatory.

### Enabled<a id='Terra-Enabled'></a>
```toml
Enabled = true # Default
```
Enabled enables this chain.

### BlockRate<a id='Terra-BlockRate'></a>
```toml
BlockRate = '6s' # Default
```
BlockRate is the average time between blocks.

### BlocksUntilTxTimeout<a id='Terra-BlocksUntilTxTimeout'></a>
```toml
BlocksUntilTxTimeout = 30 # Default
```
BlocksUntilTxTimeout is the number of blocks to wait before giving up on the tx getting confirmed.

### ConfirmPollPeriod<a id='Terra-ConfirmPollPeriod'></a>
```toml
ConfirmPollPeriod = '1s' # Default
```
ConfirmPollPeriod sets how often check for tx confirmation.

### FallbackGasPriceULuna<a id='Terra-FallbackGasPriceULuna'></a>
```toml
FallbackGasPriceULuna = '0.015' # Default
```
FallbackGasPriceULuna sets a fallback gas price to use when the estimator is not available.

### FCDURL<a id='Terra-FCDURL'></a>
```toml
FCDURL = 'http://terra.com' # Example
```
FCDURL sets the FCD URL.

### GasLimitMultiplier<a id='Terra-GasLimitMultiplier'></a>
```toml
GasLimitMultiplier = '1.5' # Default
```
GasLimitMultiplier scales the estimated gas limit.

### MaxMsgsPerBatch<a id='Terra-MaxMsgsPerBatch'></a>
```toml
MaxMsgsPerBatch = 100 # Default
```
MaxMsgsPerBatch limits the numbers of mesages per transaction batch.

### OCR2CachePollPeriod<a id='Terra-OCR2CachePollPeriod'></a>
```toml
OCR2CachePollPeriod = '4s' # Default
```
OCR2CachePollPeriod is the rate to poll for the OCR2 state cache.

### OCR2CacheTTL<a id='Terra-OCR2CacheTTL'></a>
```toml
OCR2CacheTTL = '1m' # Default
```
OCR2CacheTTL is the stale OCR2 cache deadline.

### TxMsgTimeout<a id='Terra-TxMsgTimeout'></a>
```toml
TxMsgTimeout = '10m' # Default
```
TxMsgTimeout is the maximum age for resending transaction before they expire.

## Terra.Nodes<a id='Terra-Nodes'></a>
```toml
[[Terra.Nodes]]
Name = 'primary' # Example
TendermintURL = 'http://tender.mint' # Example
```


### Name<a id='Terra-Nodes-Name'></a>
```toml
Name = 'primary' # Example
```
Name is a unique (per-chain) identifier for this node.

### TendermintURL<a id='Terra-Nodes-TendermintURL'></a>
```toml
TendermintURL = 'http://tender.mint' # Example
```
TendermintURL is the HTTP(S) tendermint endpoint for this node.

