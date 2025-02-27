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
InsecureFastScrypt = false # Default
RootDir = '~/.chainlink' # Default
ShutdownGracePeriod = '5s' # Default
```


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
CCIP = true # Default
MultiFeedsManagers = false # Default
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

### CCIP
```toml
CCIP = true # Default
```
CCIP enables the CCIP service.

### MultiFeedsManagers
```toml
MultiFeedsManagers = false # Default
```
MultiFeedsManagers enables support for multiple feeds manager connections.

## Database
```toml
[Database]
DefaultIdleInTxSessionTimeout = '1h' # Default
DefaultLockTimeout = '15s' # Default
DefaultQueryTimeout = '10s' # Default
LogQueries = false # Default
MaxIdleConns = 10 # Default
MaxOpenConns = 100 # Default
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
MaxOpenConns = 100 # Default
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
UniConn = false # Default
Logging = false # Default
BufferSize = 100 # Default
MaxBatchSize = 50 # Default
SendInterval = '500ms' # Default
SendTimeout = '10s' # Default
UseBatchSend = true # Default
```


### UniConn
```toml
UniConn = false # Default
```
UniConn toggles which ws connection style is used.

### Logging
```toml
Logging = false # Default
```
Logging toggles verbose logging of the raw telemetry messages being sent.

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

## TelemetryIngress.Endpoints
```toml
[[TelemetryIngress.Endpoints]] # Example
Network = 'EVM' # Example
ChainID = '111551111' # Example
ServerPubKey = 'test-pub-key-111551111-evm' # Example
URL = 'localhost-111551111-evm:9000' # Example
```


### Network
```toml
Network = 'EVM' # Example
```
Network aka EVM, Solana, Starknet

### ChainID
```toml
ChainID = '111551111' # Example
```
ChainID of the network

### ServerPubKey
```toml
ServerPubKey = 'test-pub-key-111551111-evm' # Example
```
ServerPubKey is the public key of the telemetry server.

### URL
```toml
URL = 'localhost-111551111-evm:9000' # Example
```
URL is where to send telemetry.

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
Level determines only what is printed on the screen/console. This configuration does not apply to the logs that are recorded in a file (see [`Log.File`](#logfile) for more details).

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
MaxSize determines the log file's max size before file rotation. Having this not set or set to a value smaller than 1Mb will disable logging to disk. If your disk doesn't have enough disk space, the logging will pause and the application will log errors until space is available again.

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
AuthenticationMethod = 'local' # Default
AllowOrigins = 'http://localhost:3000,http://localhost:6688' # Default
BridgeCacheTTL = '0s' # Default
BridgeResponseURL = 'https://my-chainlink-node.example.com:6688' # Example
HTTPWriteTimeout = '10s' # Default
HTTPPort = 6688 # Default
SecureCookies = true # Default
SessionTimeout = '15m' # Default
SessionReaperExpiration = '240h' # Default
HTTPMaxSize = '32768b' # Default
StartTimeout = '15s' # Default
ListenIP = '0.0.0.0' # Default
```


### AuthenticationMethod
```toml
AuthenticationMethod = 'local' # Default
```
AuthenticationMethod defines which pluggable auth interface to use for user login and role assumption. Options include 'local' and 'ldap'. See docs for more details

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

### HTTPMaxSize
```toml
HTTPMaxSize = '32768b' # Default
```
HTTPMaxSize defines the maximum size for HTTP requests and responses made by the node server.

### StartTimeout
```toml
StartTimeout = '15s' # Default
```
StartTimeout defines the maximum amount of time the node will wait for a server to start.

### ListenIP
```toml
ListenIP = '0.0.0.0' # Default
```
ListenIP specifies the IP to bind the HTTP server to

## WebServer.LDAP
```toml
[WebServer.LDAP]
ServerTLS = true # Default
SessionTimeout = '15m0s' # Default
QueryTimeout = '2m0s' # Default
BaseUserAttr = 'uid' # Default
BaseDN = 'dc=custom,dc=example,dc=com' # Example
UsersDN = 'ou=users' # Default
GroupsDN = 'ou=groups' # Default
ActiveAttribute = '' # Default
ActiveAttributeAllowedValue = '' # Default
AdminUserGroupCN = 'NodeAdmins' # Default
EditUserGroupCN = 'NodeEditors' # Default
RunUserGroupCN = 'NodeRunners' # Default
ReadUserGroupCN = 'NodeReadOnly' # Default
UserApiTokenEnabled = false # Default
UserAPITokenDuration = '240h0m0s' # Default
UpstreamSyncInterval = '0s' # Default
UpstreamSyncRateLimit = '2m0s' # Default
```
Optional LDAP config if WebServer.AuthenticationMethod is set to 'ldap'
LDAP queries are all parameterized to support custom LDAP 'dn', 'cn', and attributes

### ServerTLS
```toml
ServerTLS = true # Default
```
ServerTLS defines the option to require the secure ldaps

### SessionTimeout
```toml
SessionTimeout = '15m0s' # Default
```
SessionTimeout determines the amount of idle time to elapse before session cookies expire. This signs out GUI users from their sessions.

### QueryTimeout
```toml
QueryTimeout = '2m0s' # Default
```
QueryTimeout defines how long queries should wait before timing out, defined in seconds

### BaseUserAttr
```toml
BaseUserAttr = 'uid' # Default
```
BaseUserAttr defines the base attribute used to populate LDAP queries such as "uid=$", default is example

### BaseDN
```toml
BaseDN = 'dc=custom,dc=example,dc=com' # Example
```
BaseDN defines the base LDAP 'dn' search filter to apply to every LDAP query, replace example,com with the appropriate LDAP server's structure

### UsersDN
```toml
UsersDN = 'ou=users' # Default
```
UsersDN defines the 'dn' query to use when querying for the 'users' 'ou' group

### GroupsDN
```toml
GroupsDN = 'ou=groups' # Default
```
GroupsDN defines the 'dn' query to use when querying for the 'groups' 'ou' group

### ActiveAttribute
```toml
ActiveAttribute = '' # Default
```
ActiveAttribute is an optional user field to check truthiness for if a user is valid/active. This is only required if the LDAP provider lists inactive users as members of groups

### ActiveAttributeAllowedValue
```toml
ActiveAttributeAllowedValue = '' # Default
```
ActiveAttributeAllowedValue is the value to check against for the above optional user attribute

### AdminUserGroupCN
```toml
AdminUserGroupCN = 'NodeAdmins' # Default
```
AdminUserGroupCN is the LDAP 'cn' of the LDAP group that maps the core node's 'Admin' role

### EditUserGroupCN
```toml
EditUserGroupCN = 'NodeEditors' # Default
```
EditUserGroupCN is the LDAP 'cn' of the LDAP group that maps the core node's 'Edit' role

### RunUserGroupCN
```toml
RunUserGroupCN = 'NodeRunners' # Default
```
RunUserGroupCN is the LDAP 'cn' of the LDAP group that maps the core node's 'Run' role

### ReadUserGroupCN
```toml
ReadUserGroupCN = 'NodeReadOnly' # Default
```
ReadUserGroupCN is the LDAP 'cn' of the LDAP group that maps the core node's 'Read' role

### UserApiTokenEnabled
```toml
UserApiTokenEnabled = false # Default
```
UserApiTokenEnabled enables the users to issue API tokens with the same access of their role

### UserAPITokenDuration
```toml
UserAPITokenDuration = '240h0m0s' # Default
```
UserAPITokenDuration is the duration of time an API token is active for before expiring

### UpstreamSyncInterval
```toml
UpstreamSyncInterval = '0s' # Default
```
UpstreamSyncInterval is the interval at which the background LDAP sync task will be called. A '0s' value disables the background sync being run on an interval. This check is already performed during login/logout actions, all sessions and API tokens stored in the local ldap tables are updated to match the remote server

### UpstreamSyncRateLimit
```toml
UpstreamSyncRateLimit = '2m0s' # Default
```
UpstreamSyncRateLimit defines a duration to limit the number of query/API calls to the upstream LDAP provider. It prevents the sync functionality from being called multiple times within the defined duration

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
ListenIP = '0.0.0.0' # Default
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

### ListenIP
```toml
ListenIP = '0.0.0.0' # Default
```
ListenIP specifies the IP to bind the HTTPS server to

## JobPipeline
```toml
[JobPipeline]
ExternalInitiatorsEnabled = false # Default
MaxRunDuration = '10m' # Default
MaxSuccessfulRuns = 10000 # Default
ReaperInterval = '1h' # Default
ReaperThreshold = '24h' # Default
ResultWriteQueueDepth = 100 # Default
VerboseLogging = true # Default
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

### VerboseLogging
```toml
VerboseLogging = true # Default
```
VerboseLogging enables detailed logging of pipeline execution steps.
This can be useful for debugging failed runs without relying on the UI
or database.

You may disable if this results in excessive log volume.

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
CaptureAutomationCustomTelemetry = true # Default
AllowNoBootstrappers = false # Default
DefaultTransactionQueueDepth = 1 # Default
SimulateTransactions = false # Default
TraceLogging = false # Default
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

### CaptureAutomationCustomTelemetry
```toml
CaptureAutomationCustomTelemetry = true # Default
```
CaptureAutomationCustomTelemetry toggles collecting automation specific telemetry

### AllowNoBootstrappers
```toml
AllowNoBootstrappers = false # Default
```
AllowNoBootstrappers enables single-node consensus without bootstrapper nodes (i.e. f=0, n=1)

### DefaultTransactionQueueDepth
```toml
DefaultTransactionQueueDepth = 1 # Default
```
DefaultTransactionQueueDepth controls the queue size for `DropOldestStrategy` in OCR2. Set to 0 to use `SendEvery` strategy instead.

### SimulateTransactions
```toml
SimulateTransactions = false # Default
```
SimulateTransactions enables transaction simulation for OCR2.

### TraceLogging
```toml
TraceLogging = false # Default
```
TraceLogging enables trace level logging.

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
TraceLogging = false # Default
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

### TraceLogging
```toml
TraceLogging = false # Default
```
TraceLogging enables trace level logging.

## P2P
```toml
[P2P]
IncomingMessageBufferSize = 10 # Default
OutgoingMessageBufferSize = 10 # Default
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw' # Example
TraceLogging = false # Default
```
P2P has a versioned networking stack. Currenly only `[P2P.V2]` is supported.
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

## P2P.V2
```toml
[P2P.V2]
Enabled = true # Default
AnnounceAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
DefaultBootstrappers = ['12D3KooWMHMRLQkgPbFSYHwD3NBuwtS1AmxhvKVUrcfyaGDASR4U@1.2.3.4:9999', '12D3KooWM55u5Swtpw9r8aFLQHEtw7HR4t44GdNs654ej5gRs2Dh@example.com:1234'] # Example
DeltaDial = '15s' # Default
DeltaReconcile = '1m' # Default
ListenAddresses = ['1.2.3.4:9999', '[a52d:0:a88:1274::abcd]:1337'] # Example
```


### Enabled
```toml
Enabled = true # Default
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

## Capabilities.RateLimit
```toml
[Capabilities.RateLimit]
GlobalRPS = 200 # Default
GlobalBurst = 200 # Default
PerSenderRPS = 100 # Default
PerSenderBurst = 100 # Default
```


### GlobalRPS
```toml
GlobalRPS = 200 # Default
```
GlobalRPS is the global rate limit for the dispatcher.

### GlobalBurst
```toml
GlobalBurst = 200 # Default
```
GlobalBurst is the global burst limit for the dispatcher.

### PerSenderRPS
```toml
PerSenderRPS = 100 # Default
```
PerSenderRPS is the per-sender rate limit for the dispatcher.

### PerSenderBurst
```toml
PerSenderBurst = 100 # Default
```
PerSenderBurst is the per-sender burst limit for the dispatcher.

## Capabilities.WorkflowRegistry
```toml
[Capabilities.WorkflowRegistry]
Address = '0x0' # Example
NetworkID = 'evm' # Default
ChainID = '1' # Default
MaxBinarySize = '20.00mb' # Default
MaxEncryptedSecretsSize = '26.40kb' # Default
MaxConfigSize = '50.00kb' # Default
```


### Address
```toml
Address = '0x0' # Example
```
Address is the address for the workflow registry contract.

### NetworkID
```toml
NetworkID = 'evm' # Default
```
NetworkID identifies the target network where the remote registry is located.

### ChainID
```toml
ChainID = '1' # Default
```
ChainID identifies the target chain id where the remote registry is located.

### MaxBinarySize
```toml
MaxBinarySize = '20.00mb' # Default
```
MaxBinarySize is the maximum size of a binary that can be fetched from the registry.

### MaxEncryptedSecretsSize
```toml
MaxEncryptedSecretsSize = '26.40kb' # Default
```
MaxEncryptedSecretsSize is the maximum size of encrypted secrets that can be fetched from the given secrets url.

### MaxConfigSize
```toml
MaxConfigSize = '50.00kb' # Default
```
MaxConfigSize is the maximum size of a config that can be fetched from the given config url.

## Capabilities.ExternalRegistry
```toml
[Capabilities.ExternalRegistry]
Address = '0x0' # Example
NetworkID = 'evm' # Default
ChainID = '1' # Default
```


### Address
```toml
Address = '0x0' # Example
```
Address is the address for the capabilities registry contract.

### NetworkID
```toml
NetworkID = 'evm' # Default
```
NetworkID identifies the target network where the remote registry is located.

### ChainID
```toml
ChainID = '1' # Default
```
ChainID identifies the target chain id where the remote registry is located.

## Capabilities.Dispatcher
```toml
[Capabilities.Dispatcher]
SupportedVersion = 1 # Default
ReceiverBufferSize = 10000 # Default
```


### SupportedVersion
```toml
SupportedVersion = 1 # Default
```
SupportedVersion is the version of the version of message schema.

### ReceiverBufferSize
```toml
ReceiverBufferSize = 10000 # Default
```
ReceiverBufferSize is the size of the buffer for incoming messages.

## Capabilities.Dispatcher.RateLimit
```toml
[Capabilities.Dispatcher.RateLimit]
GlobalRPS = 800 # Default
GlobalBurst = 1000 # Default
PerSenderRPS = 10 # Default
PerSenderBurst = 50 # Default
```


### GlobalRPS
```toml
GlobalRPS = 800 # Default
```
GlobalRPS is the global rate limit for the dispatcher.

### GlobalBurst
```toml
GlobalBurst = 1000 # Default
```
GlobalBurst is the global burst limit for the dispatcher.

### PerSenderRPS
```toml
PerSenderRPS = 10 # Default
```
PerSenderRPS is the per-sender rate limit for the dispatcher.

### PerSenderBurst
```toml
PerSenderBurst = 50 # Default
```
PerSenderBurst is the per-sender burst limit for the dispatcher.

## Capabilities.Peering
```toml
[Capabilities.Peering]
IncomingMessageBufferSize = 10 # Default
OutgoingMessageBufferSize = 10 # Default
PeerID = '12D3KooWMoejJznyDuEk5aX6GvbjaG12UzeornPCBNzMRqdwrFJw' # Example
TraceLogging = false # Default
```


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

## Capabilities.Peering.V2
```toml
[Capabilities.Peering.V2]
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

## Capabilities.GatewayConnector
```toml
[Capabilities.GatewayConnector]
ChainIDForNodeKey = '11155111' # Example
NodeAddress = '0x68902d681c28119f9b2531473a417088bf008e59' # Example
DonID = 'example_don' # Example
WSHandshakeTimeoutMillis = 1000 # Example
AuthMinChallengeLen = 10 # Example
AuthTimestampToleranceSec = 10 # Example
```


### ChainIDForNodeKey
```toml
ChainIDForNodeKey = '11155111' # Example
```
ChainIDForNodeKey is the ChainID of the network associated with a private key to be used for authentication with Gateway nodes

### NodeAddress
```toml
NodeAddress = '0x68902d681c28119f9b2531473a417088bf008e59' # Example
```
NodeAddress is the address of the desired private key to be used for authentication with Gateway nodes

### DonID
```toml
DonID = 'example_don' # Example
```
DonID is the Id of the Don

### WSHandshakeTimeoutMillis
```toml
WSHandshakeTimeoutMillis = 1000 # Example
```
WSHandshakeTimeoutMillis is Websocket handshake timeout

### AuthMinChallengeLen
```toml
AuthMinChallengeLen = 10 # Example
```
AuthMinChallengeLen is the minimum number of bytes in authentication challenge payload

### AuthTimestampToleranceSec
```toml
AuthTimestampToleranceSec = 10 # Example
```
AuthTimestampToleranceSec is Authentication timestamp tolerance

## Capabilities.GatewayConnector.Gateways
```toml
[[Capabilities.GatewayConnector.Gateways]]
ID = 'example_gateway' # Example
URL = 'wss://localhost:8081/node' # Example
```


### ID
```toml
ID = 'example_gateway' # Example
```
ID of the Gateway

### URL
```toml
URL = 'wss://localhost:8081/node' # Example
```
URL of the Gateway

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

## Tracing
```toml
[Tracing]
Enabled = false # Default
CollectorTarget = 'localhost:4317' # Example
NodeID = 'NodeID' # Example
SamplingRatio = 1.0 # Example
Mode = 'tls' # Default
TLSCertPath = '/path/to/cert.pem' # Example
```


### Enabled
```toml
Enabled = false # Default
```
Enabled turns trace collection on or off. On requires an OTEL Tracing Collector.

### CollectorTarget
```toml
CollectorTarget = 'localhost:4317' # Example
```
CollectorTarget is the logical address of the OTEL Tracing Collector.

### NodeID
```toml
NodeID = 'NodeID' # Example
```
NodeID is an unique name for this node relative to any other node traces are collected for.

### SamplingRatio
```toml
SamplingRatio = 1.0 # Example
```
SamplingRatio is the ratio of traces to sample for this node.

### Mode
```toml
Mode = 'tls' # Default
```
Mode is a string value. `tls` or `unencrypted` are the only values allowed. If set to `unencrypted`, `TLSCertPath` can be unset, meaning traces will be sent over plaintext to the collector.

### TLSCertPath
```toml
TLSCertPath = '/path/to/cert.pem' # Example
```
TLSCertPath is the file path to the TLS certificate used for secure communication with an OTEL Tracing Collector.

## Tracing.Attributes
```toml
[Tracing.Attributes]
env = 'test' # Example
```
Tracing.Attributes are user specified key-value pairs to associate in the context of the traces

### env
```toml
env = 'test' # Example
```
env is an example user specified key-value pair

## Mercury
```toml
[Mercury]
VerboseLogging = false # Default
```


### VerboseLogging
```toml
VerboseLogging = false # Default
```
VerboseLogging enables detailed logging of mercury/LLO operations. These logs
can be expensive since they may serialize large structs, so they are disabled
by default.

## Mercury.Cache
```toml
[Mercury.Cache]
LatestReportTTL = "1s" # Default
MaxStaleAge = "1h" # Default
LatestReportDeadline = "5s" # Default
```
Mercury.Cache controls settings for the price retrieval cache querying a mercury server

### LatestReportTTL
```toml
LatestReportTTL = "1s" # Default
```
LatestReportTTL controls how "stale" we will allow a price to be e.g. if
set to 1s, a new price will always be fetched if the last result was
from 1 second ago or older.

Another way of looking at it is such: the cache will _never_ return a
price that was queried from now-LatestReportTTL or before.

Setting to zero disables caching entirely.

### MaxStaleAge
```toml
MaxStaleAge = "1h" # Default
```
MaxStaleAge is that maximum amount of time that a value can be stale
before it is deleted from the cache (a form of garbage collection).

This should generally be set to something much larger than
LatestReportTTL. Setting to zero disables garbage collection.

### LatestReportDeadline
```toml
LatestReportDeadline = "5s" # Default
```
LatestReportDeadline controls how long to wait for a response from the
mercury server before retrying. Setting this to zero will wait indefinitely.

## Mercury.TLS
```toml
[Mercury.TLS]
CertFile = "/path/to/client/certs.pem" # Example
```
Mercury.TLS controls client settings for when the node talks to traditional web servers or load balancers.

### CertFile
```toml
CertFile = "/path/to/client/certs.pem" # Example
```
CertFile is the path to a PEM file of trusted root certificate authority certificates

## Mercury.Transmitter
```toml
[Mercury.Transmitter]
Protocol = "wsrpc" # Default
TransmitQueueMaxSize = 100_000 # Default
TransmitTimeout = "5s" # Default
TransmitConcurrency = 100 # Default
ReaperFrequency = "1h" # Default
ReaperMaxAge = "48h" # Default
```
Mercury.Transmitter controls settings for the mercury transmitter

### Protocol
```toml
Protocol = "wsrpc" # Default
```
Protocol is the protocol to use for the transmitter.

Options are either:
- "wsrpc" for the legacy websocket protocol
- "grpc" for the gRPC protocol

### TransmitQueueMaxSize
```toml
TransmitQueueMaxSize = 100_000 # Default
```
TransmitQueueMaxSize controls the size of the transmit queue. This is scoped
per OCR instance. If the queue is full, the transmitter will start dropping
the oldest messages in order to make space.

This is useful if mercury server goes offline and the nop needs to buffer
transmissions.

### TransmitTimeout
```toml
TransmitTimeout = "5s" # Default
```
TransmitTimeout controls how long the transmitter will wait for a response
when sending a message to the mercury server, before aborting and considering
the transmission to be failed.

### TransmitConcurrency
```toml
TransmitConcurrency = 100 # Default
```
TransmitConcurrency is the max number of concurrent transmits to each server.

Only has effect with LLO jobs.

### ReaperFrequency
```toml
ReaperFrequency = "1h" # Default
```
ReaperFrequency controls how often the stale transmission reaper will run.
Setting to 0 disables the reaper.

### ReaperMaxAge
```toml
ReaperMaxAge = "48h" # Default
```
ReaperMaxAge controls how old a transmission can be before it is considered
stale. Setting to 0 disables the reaper.

## Telemetry
```toml
[Telemetry]
Enabled = false # Default
Endpoint = 'example.com/collector' # Example
CACertFile = 'cert-file' # Example
InsecureConnection = false # Default
TraceSampleRatio = 0.01 # Default
EmitterBatchProcessor = true # Default
EmitterExportTimeout = '1s' # Default
```
Telemetry holds OTEL settings.
This data includes open telemetry metrics, traces, & logs.
It does not currently include prometheus metrics or standard out logs, but may in the future.

### Enabled
```toml
Enabled = false # Default
```
Enabled turns telemetry collection on or off.

### Endpoint
```toml
Endpoint = 'example.com/collector' # Example
```
Endpoint of the OTEL Collector.

### CACertFile
```toml
CACertFile = 'cert-file' # Example
```
CACertFile is the file path of the TLS certificate used for secure communication with the OTEL Collector.
Required unless InescureConnection is true.

### InsecureConnection
```toml
InsecureConnection = false # Default
```
InsecureConnection bypasses the TLS CACertFile requirement and uses an insecure connection instead.
Only available in dev mode.

### TraceSampleRatio
```toml
TraceSampleRatio = 0.01 # Default
```
TraceSampleRatio is the rate at which to sample traces. Must be between 0 and 1.

### EmitterBatchProcessor
```toml
EmitterBatchProcessor = true # Default
```
EmitterBatchProcessor enables batching for telemetry events

### EmitterExportTimeout
```toml
EmitterExportTimeout = '1s' # Default
```
EmitterExportTimeout sets timeout for exporting telemetry events

## Telemetry.ResourceAttributes
```toml
[Telemetry.ResourceAttributes]
foo = "bar" # Example
```
ResourceAttributes are global metadata to include with all telemetry.

### foo
```toml
foo = "bar" # Example
```
foo is an example resource attribute

## EVM
EVM defaults depend on ChainID:

<details><summary>Ethereum Mainnet (1)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x514910771AF9Ca656af840dff83E8264EcF986CA'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
OperatorFactoryAddress = '0x3E64Cd889482443324F91bFA9c84fE72A511f48A'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '9m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 10500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ethereum Ropsten (3)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x20fE562d797A42Dcb3399062AE9546cd06f63280'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ethereum Rinkeby (4)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x01BE23585060835E02B77ef475b0Cc51aA1e0709'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ethereum Goerli (5)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x326C977E6efc84E512bB9C30f76E30c160eD06FB'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Optimism Mainnet (10)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '13m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>RSK Mainnet (30)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x14AdaE34beF7ca957Ce2dDe5ADD97ea050123827'
LogBackfillBatchSize = 1000
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 mwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>RSK Testnet (31)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x8bBbd80981FE76d44854D8DF305e8985c19f0e78'
LogBackfillBatchSize = 1000
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 mwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ethereum Kovan (42)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0xa36085F69e2889c224210F603D836748e7dC0088'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
OperatorFactoryAddress = '0x8007e24251b1D2Fc518Eb843A701d9cD21fe0aA3'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>BSC Mainnet (56)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x404460C6A5EdE2D891e8297795264fDe62ADBB75'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 2
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '45s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '2s'
DatabaseTimeout = '2s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '500ms'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>OKX Testnet (65)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>OKX Mainnet (66)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Astar Shibuya (81)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 100
FinalityTagEnabled = true
LinkContractAddress = '0xe74037112db8807B3B4B3895F5790e5bc1866a29'
LogBackfillBatchSize = 1000
LogPollInterval = '6s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>BSC Testnet (97)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x84b9B910527Ad5C03A9Ca831909E21e236EA7b06'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 2
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '40s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = false
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '2s'
DatabaseTimeout = '2s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '500ms'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Gnosis Mainnet (100)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'gnosis'
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0xE2e73A1c69ecF83F464EFCE6A5be353a37cA09b2'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '2m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Shibarium Mainnet (109)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x71052BAe71C25C78E37fD12E5ff1101A71d9018F'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '10 gwei'
PriceMax = '10 micro'
PriceMin = '2 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = true
FeeCapDefault = '10 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 2
PollInterval = '3s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Heco Mainnet (128)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x404460C6A5EdE2D891e8297795264fDe62ADBB75'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 2
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '2s'
DatabaseTimeout = '2s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '500ms'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Hashkey Testnet (133)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x8418c4d7e8e17ab90232DC72150730E6c4b84F57'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '1 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 2
PollInterval = '8s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Polygon Mainnet (137)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 500
FinalityTagEnabled = true
LinkContractAddress = '0xb0897686c545045aFc77CF20eC7A532E3120E0F1'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '6m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 5000
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '20 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Sonic Mainnet (146)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x71052BAe71C25C78E37fD12E5ff1101A71d9018F'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 500
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 10
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '2s'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Shibarium Testnet (157)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x44637eEfD71A090990f89faEC7022fc74B2969aD'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '10 gwei'
PriceMax = '10 micro'
PriceMin = '2 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = true
FeeCapDefault = '10 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 2
PollInterval = '3s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Hashkey Mainnet (177)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x71052BAe71C25C78E37fD12E5ff1101A71d9018F'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '1 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 2
PollInterval = '8s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>XLayer Sepolia (195)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'xlayer'
FinalityDepth = 500
FinalityTagEnabled = false
LinkContractAddress = '0x724593f6FCb0De4E6902d4C55D7C74DaA2AF0E55'
LogBackfillBatchSize = 1000
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '12m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 15
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 mwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '20 mwei'
BumpPercent = 40
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 12
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>XLayer Mainnet (196)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'xlayer'
FinalityDepth = 500
FinalityTagEnabled = false
LinkContractAddress = '0x8aF9711B44695a5A081F25AB9903DDB73aCf8FA9'
LogBackfillBatchSize = 1000
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '6m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 15
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '100 mwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '100 mwei'
BumpPercent = 40
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 12
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Bsquared Mainnet (223)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 2000
FinalityTagEnabled = true
LinkContractAddress = '0x709229D9587886a1eDFeE6b5cE636E1D70d1cE39'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h10m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Fantom Mainnet (250)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x6F43FF82CCA38001B6699a8AC47A2d0E66939407'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 2
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 3800000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Fraxtal Mainnet (252)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Kroma Mainnet (255)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'kroma'
FinalityDepth = 400
FinalityTagEnabled = true
LinkContractAddress = '0xC1F6f7622ad37C3f46cDF6F8AA0344ADE80BF450'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x4200000000000000000000000000000000000005'

[HeadTracker]
HistoryDepth = 400
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>ZKsync Goerli (280)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zksync'
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0xD29F4Cc763A064b6C563B8816f09351b3Fbb61A0'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '18.446744073709551615 ether'
PriceMin = '0'
LimitDefault = 100000000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'zksync'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Hedera Mainnet (295)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'hedera'
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x7Ce6bb2Cc2D3Fd45a974Da6a0F29236cb9513a98'
LogBackfillBatchSize = 1000
LogPollInterval = '10s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '2m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = true
BumpMin = '10 gwei'
BumpPercent = 20
BumpThreshold = 0
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Hedera Testnet (296)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'hedera'
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x90a386d59b9A6a4795a011e8f032Fc21ED6FEFb6'
LogBackfillBatchSize = 1000
LogPollInterval = '10s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '2m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = true
BumpMin = '10 gwei'
BumpPercent = 20
BumpThreshold = 0
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>ZKsync Sepolia (300)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zksync'
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x23A1aFD896c8c8876AF46aDc38521f4432658d1e'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '18.446744073709551615 ether'
PriceMin = '0'
LimitDefault = 100000000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'zksync'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 11000000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>ZKsync Mainnet (324)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zksync'
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x52869bae3E091e36b0915941577F2D47d8d8B534'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '18.446744073709551615 ether'
PriceMin = '0'
LimitDefault = 100000000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'zksync'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 11000000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Optimism Goerli (420)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = false
LinkContractAddress = '0xdc2CC710e42857672E7907CF474a69B63B93089f'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 60
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Worldchain Mainnet (480)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 2500
FinalityTagEnabled = true
LinkContractAddress = '0x915b648e994d5f31059B38223b9fbe98ae185473'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h30m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Metis Rinkeby (588)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'metis'
FinalityDepth = 10
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Astar Mainnet (592)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'astar'
FinalityDepth = 100
FinalityTagEnabled = true
LinkContractAddress = '0x31EFB841d5e0b4082F7E1267dab8De1b853f2A9d'
LogBackfillBatchSize = 1000
LogPollInterval = '6s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '100 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Mode Sepolia (919)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x925a4bfE64AE2bFAC8a02b35F78e60C29743755d'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '120 gwei'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 60
EIP1559DynamicFees = true
FeeCapDefault = '120 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 2
PollInterval = '3s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Klaytn Testnet (1001)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '750 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Metis Mainnet (1088)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0xd2FE54D1E5F568eB710ba9d898Bf4bD02C7c0353'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Polygon Zkevm Mainnet (1101)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zkevm'
FinalityDepth = 500
FinalityTagEnabled = false
LinkContractAddress = '0xdB7A504CF869484dd6aC5FaF925c8386CBF7573D'
LogBackfillBatchSize = 1000
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '6m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 15
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 40
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>WeMix Mainnet (1111)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'wemix'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x80f1FcdC96B55e459BF52b998aBBE2c364935d69'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '40s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '100 gwei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>WeMix Testnet (1112)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'wemix'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x3580c7A817cCD41f7e02143BFa411D4EeAE78093'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '40s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '100 gwei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = false
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Core Testnet (1114)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 7
FinalityTagEnabled = false
LinkContractAddress = '0x6C475841d1D7871940E93579E5DBaE01634e17aA'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '35 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Core Mainnet (1116)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 27
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '30 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Bsquared Testnet (1123)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 2000
FinalityTagEnabled = true
LinkContractAddress = '0x436a1907D9e6a65E6db73015F08f9C66F6B63E45'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h10m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Unichain Testnet (1301)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 2000
FinalityTagEnabled = true
LinkContractAddress = '0xda40816f278Cd049c137F6612822D181065EBfB4'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '45m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '2s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Simulated (1337)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '100'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '0s'
ResendAfterThreshold = '0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FixedPrice'
PriceDefault = '1 gwei'
PriceMax = '1 gwei'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 0
EIP1559DynamicFees = false
FeeCapDefault = '1 gwei'
TipCapDefault = '1 mwei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 10
MaxBufferSize = 100
SamplingInterval = '0s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = false
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Soneium Mainnet (1868)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x32D8F819C8080ae44375F8d383Ffd39FC642f3Ec'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '2h0m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '1 mwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 60
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Soneium Sepolia (1946)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x7ea13478Ea3961A0e8b538cb05a9DF0477c79Cd2'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '2h0m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '1 mwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 60
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ronin Mainnet (2020)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x3902228D6A3d2Dc44731fD9d45FeE6a61c722D0b'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '1 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ronin Saigon (2021)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x5bB50A6888ee6a67E22afFDFD9513be7740F1c15'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '1 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Kroma Sepolia (2358)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'kroma'
FinalityDepth = 400
FinalityTagEnabled = true
LinkContractAddress = '0xa75cCA5b404ec6F4BB6EC4853D177FE7057085c8'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x4200000000000000000000000000000000000005'

[HeadTracker]
HistoryDepth = 400
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Polygon Zkevm Cardona (2442)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zkevm'
FinalityDepth = 500
FinalityTagEnabled = false
LinkContractAddress = '0x5576815a38A3706f37bf815b261cCc7cCA77e975'
LogBackfillBatchSize = 1000
LogPollInterval = '30s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '12m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 40
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Fraxtal Testnet (2522)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Botanix Testnet (3636)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x7311DED199CC28D80E58e81e8589aa160199FCD2'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = false
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Fantom Testnet (4002)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0xfaFedb041c0DD4fA2Dc0d87a6B0979Ee6FA7af5F'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 2
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 3800000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Merlin Mainnet (4200)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zkevm'
FinalityDepth = 1000
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '4s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Worldchain Testnet (4801)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 2500
FinalityTagEnabled = true
LinkContractAddress = '0xC82Ea35634BcE95C394B6BC00626f827bB0F4801'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h30m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Mantle Mainnet (5000)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 1200
FinalityTagEnabled = true
LinkContractAddress = '0xfe36cF0B43aAe49fBc5cFC5c0AF22a623114E043'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '40m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '120 gwei'
PriceMin = '1 gwei'
LimitDefault = 80000000000
LimitMax = 100000000000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 60
EIP1559DynamicFees = true
FeeCapDefault = '120 gwei'
TipCapDefault = '0'
TipCapMin = '0'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 200
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
EIP1559FeeCapBufferBlocks = 0
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 1250
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Mantle Sepolia (5003)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 1200
FinalityTagEnabled = true
LinkContractAddress = '0x22bdEdEa0beBdD7CfFC95bA53826E55afFE9DE04'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h0m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '120 gwei'
PriceMin = '1 gwei'
LimitDefault = 80000000000
LimitMax = 100000000000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 60
EIP1559DynamicFees = true
FeeCapDefault = '120 gwei'
TipCapDefault = '0'
TipCapMin = '0'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 200
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
EIP1559FeeCapBufferBlocks = 0
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 1250
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Klaytn Mainnet (8217)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '750 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Base Mainnet (8453)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x88Fb150BDc53A65fe94Dea0c9BA0a6dAf8C6e196'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '15m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Monad Testnet (10143)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 120
FinalityTagEnabled = false
LogBackfillBatchSize = 100
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
LogBroadcasterEnabled = false
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '2s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '4s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Gnosis Chiado (10200)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'gnosis'
FinalityDepth = 100
FinalityTagEnabled = false
LinkContractAddress = '0xDCA67FD8324990792C0bfaE95903B8A64097754F'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '2m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '500 gwei'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>L3X Mainnet (12324)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x79f531a3D07214304F259DC28c7191513223bcf3'
LogBackfillBatchSize = 1000
LogPollInterval = '10s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>L3X Sepolia (12325)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0xa71848C99155DA0b245981E5ebD1C94C4be51c43'
LogBackfillBatchSize = 1000
LogPollInterval = '10s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ethereum Holesky (17000)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x685cE6742351ae9b618F383883D6d1e0c5A31B4B'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '1 micro'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 2
PollInterval = '3s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Mode Mainnet (34443)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x183E3691EfF3524B2315D3703D94F922CbE51F54'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '120 gwei'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 60
EIP1559DynamicFees = true
FeeCapDefault = '120 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 2
PollInterval = '3s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Lens Sepolia (37111)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zksync'
FinalityDepth = 40
FinalityTagEnabled = false
LinkContractAddress = '0x7f1b9eE544f9ff9bB521Ab79c205d79C55250a36'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '10m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '7m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '70 gwei'
PriceMax = '2 micro'
PriceMin = '70 gwei'
LimitDefault = 100000000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 40
BumpThreshold = 1
EIP1559DynamicFees = false
FeeCapDefault = '2 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'zksync'

[HeadTracker]
HistoryDepth = 250
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Arbitrum Mainnet (42161)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0xf97f4df75117a78c1A5a0DBb814Af92458539FB4'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 14500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Celo Mainnet (42220)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'celo'
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0xd07294e6E917e07dfDcee882dd1e2565085C2ae0'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '1m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '5 gwei'
PriceMax = '500 gwei'
PriceMin = '5 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '2 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 12
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Avalanche Fuji (43113)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 2
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '1m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = false
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Avalanche Mainnet (43114)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x5947BB275c521040051D82396192181b413227A3'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 2
FinalizedBlockOffset = 2
NoNewFinalizedHeadsThreshold = '1m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Celo Testnet (44787)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'celo'
FinalityDepth = 2750
FinalityTagEnabled = true
LinkContractAddress = '0x32E08557B14FaD8908025619797221281D439071'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '1m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '45m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '1 micro'
PriceMin = '5 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 200
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Zircuit Sepolia (48899)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zircuit'
FinalityDepth = 1000
FinalityTagEnabled = true
LinkContractAddress = '0xDEE94506570cA186BC1e3516fCf4fd719C312cCD'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '40m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = true
Threshold = 90
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 60
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Zircuit Mainnet (48900)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zircuit'
FinalityDepth = 1000
FinalityTagEnabled = true
LinkContractAddress = '0x5D6d033B4FbD2190D99D930719fAbAcB64d2439a'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '40m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = true
Threshold = 90
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Sonic Testnet (57054)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x61876F0429726D7777B46f663e1C9ab75d08Fc56'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 500
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 10
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '2s'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ink Mainnet (57073)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 3000
FinalityTagEnabled = true
LinkContractAddress = '0x71052BAe71C25C78E37fD12E5ff1101A71d9018F'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h0m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '2s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Linea Goerli (59140)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 15
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 40
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Linea Sepolia (59141)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 900
FinalityTagEnabled = false
LinkContractAddress = '0xF64E6E064a71B45514691D397ad4204972cD6508'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
Threshold = 50
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 1000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Linea Mainnet (59144)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 300
FinalityTagEnabled = false
LinkContractAddress = '0xa18152629128738a5c081eb226335FEd4B9C95e9'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
Threshold = 50
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '400 mwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 40
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 350
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Metis Sepolia (59902)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x9870D6a0e05F867EAAe696e106741843F7fD116D'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'SuggestedPrice'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>BOB Mainnet (60808)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 3150
FinalityTagEnabled = true
LinkContractAddress = '0x5aB885CDa7216b163fb6F813DEC1E1532516c833'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h50m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Treasure Mainnet (61166)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zksync'
FinalityDepth = 50
FinalityTagEnabled = true
LogBackfillBatchSize = 1000
LogPollInterval = '10s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '5m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '25 mwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'zksync'

[HeadTracker]
HistoryDepth = 250
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Polygon Mumbai (80001)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 500
FinalityTagEnabled = false
LinkContractAddress = '0x326C977E6efc84E512bB9C30f76E30c160eD06FB'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 5000
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '20 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Polygon Amoy (80002)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 500
FinalityTagEnabled = false
LinkContractAddress = '0x0Fd9e8d3aF1aaee056EB9e802c3A762a667b1904'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 5
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 10
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '12m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 5000
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '20 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Berachain Testnet (80084)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x52CEEed7d3f8c6618e4aaD6c6e555320d0D83271'
LogBackfillBatchSize = 1000
LogPollInterval = '6s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '5m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Berachain Mainnet (80094)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 10
FinalityTagEnabled = false
LinkContractAddress = '0x52CEEed7d3f8c6618e4aaD6c6e555320d0D83271'
LogBackfillBatchSize = 1000
LogPollInterval = '6s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '5m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Blast Mainnet (81457)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x93202eC683288a9EA75BB829c6baCFb2BfeA9013'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '120 gwei'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 60
EIP1559DynamicFees = true
FeeCapDefault = '120 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 4
PollInterval = '4s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Base Goerli (84531)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 60
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Base Sepolia (84532)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0xE4aB69C077896252FAFBD49EFD26B5D171A32410'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '12m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 60
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Bitlayer Testnet (200810)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 21
FinalityTagEnabled = false
LinkContractAddress = '0x2A5bACb2440BC17D53B7b9Be73512dDf92265e48'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '1 gwei'
PriceMax = '1 gwei'
PriceMin = '40 mwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '1 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Bitlayer Mainnet (200901)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 21
FinalityTagEnabled = false
LinkContractAddress = '0x56B275c0Ec034a229a1deD8DB17089544bc276D9'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '1 gwei'
PriceMax = '1 gwei'
PriceMin = '40 mwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '1 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Arbitrum Rinkeby (421611)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x615fBe6372676474d9e6933d310469c9b68e9726'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Arbitrum Goerli (421613)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0xd14838A68E8AFBAdE5efb411d5871ea0011AFd28'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 14500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Arbitrum Sepolia (421614)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0xb1D4538B4571d411F07960EF2838Ce337FE1E80E'
LogBackfillBatchSize = 1000
LogPollInterval = '1s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 5
EIP1559DynamicFees = false
FeeCapDefault = '1 micro'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 0
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 14500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Scroll Sepolia (534351)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'scroll'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x231d45b53C905c3d6201318156BDC725c9c3B9B1'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = true
DetectionApiUrl = 'https://sepolia-venus.scroll.io'

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '1 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x5300000000000000000000000000000000000002'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Scroll Mainnet (534352)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'scroll'
FinalityDepth = 10
FinalityTagEnabled = true
LinkContractAddress = '0x548C6944cba02B9D1C0570102c89de64D258d3Ac'
LogBackfillBatchSize = 1000
LogPollInterval = '5s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = true
DetectionApiUrl = 'https://venus.scroll.io'

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '1 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 24
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x5300000000000000000000000000000000000002'

[HeadTracker]
HistoryDepth = 50
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Merlin Testnet (686868)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zkevm'
FinalityDepth = 1000
FinalityTagEnabled = false
LogBackfillBatchSize = 1000
LogPollInterval = '4s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 100
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '3m0s'

[Transactions.AutoPurge]
Enabled = true
MinAttempts = 3

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 2000
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ink Testnet (763373)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 3000
FinalityTagEnabled = true
LinkContractAddress = '0x3423C922911956b1Ccbc2b5d4f38216a6f4299b4'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h0m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '2s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>BOB Testnet (808813)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 3150
FinalityTagEnabled = true
LinkContractAddress = '0xcd2AfB2933391E35e8682cbaaF75d9CA7339b183'
LogBackfillBatchSize = 1000
LogPollInterval = '3s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '1h50m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'FeeHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 100
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '4s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Treasure Testnet (978658)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'zksync'
FinalityDepth = 50
FinalityTagEnabled = true
LogBackfillBatchSize = 1000
LogPollInterval = '10s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '10m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '115792089237316195423570985008687907853269984665.640564039457584007913129639935 tether'
PriceMin = '100 mwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'zksync'

[HeadTracker]
HistoryDepth = 250
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Ethereum Sepolia (11155111)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x779877A7B0D9E8603169DdbD7836e478b4624789'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.1 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 4
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 50

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = false
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 10500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Optimism Sepolia (11155420)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0xE4aB69C077896252FAFBD49EFD26B5D171A32410'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '40s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '15m0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '30s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = true
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 60
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 10
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 1
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 6500000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Corn Mainnet (21000000)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x7311DED199CC28D80E58e81e8589aa160199FCD2'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'Arbitrum'
PriceDefault = '100 mwei'
PriceMax = '1 micro'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 1000000000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Corn Testnet (21000001)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'arbitrum'
FinalityDepth = 50
FinalityTagEnabled = true
LinkContractAddress = '0x996EfAb6011896Be832969D91E9bc1b3983cfdA1'
LogBackfillBatchSize = 1000
LogPollInterval = '15s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'Arbitrum'
PriceDefault = '100 mwei'
PriceMax = '1 micro'
PriceMin = '0'
LimitDefault = 500000
LimitMax = 1000000000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'arbitrum'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Blast Sepolia (168587773)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
ChainType = 'optimismBedrock'
FinalityDepth = 200
FinalityTagEnabled = true
LinkContractAddress = '0x02c359ebf98fc8BF793F970F9B8302bb373BdF32'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 3
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '3m0s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

[BalanceMonitor]
Enabled = true

[GasEstimator]
Mode = 'BlockHistory'
PriceDefault = '20 gwei'
PriceMax = '120 gwei'
PriceMin = '1 gwei'
LimitDefault = 500000
LimitMax = 500000
LimitMultiplier = '1'
LimitTransfer = 21000
EstimateLimit = false
BumpMin = '100 wei'
BumpPercent = 20
BumpThreshold = 60
EIP1559DynamicFees = true
FeeCapDefault = '120 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[GasEstimator.DAOracle]
OracleType = 'opstack'
OracleAddress = '0x420000000000000000000000000000000000000F'

[HeadTracker]
HistoryDepth = 300
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 4
PollInterval = '4s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Harmony Mainnet (1666600000)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x218532a12a389a4a92fC0C5Fb22901D1c19198aA'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
```

</p></details>

<details><summary>Harmony Testnet (1666700000)</summary><p>

```toml
AutoCreateKey = true
BlockBackfillDepth = 10
BlockBackfillSkip = false
FinalityDepth = 50
FinalityTagEnabled = false
LinkContractAddress = '0x8b12Ac23BFe11cAb03a634C1F117D64a7f2cFD3e'
LogBackfillBatchSize = 1000
LogPollInterval = '2s'
LogKeepBlocksDepth = 100000
LogPrunePageSize = 0
BackupLogPollerBlockDelay = 100
MinIncomingConfirmations = 1
MinContractPayment = '0.00001 link'
NonceAutoSync = true
NoNewHeadsThreshold = '30s'
LogBroadcasterEnabled = true
RPCDefaultBatchSize = 250
RPCBlockQueryDelay = 1
FinalizedBlockOffset = 0
NoNewFinalizedHeadsThreshold = '0s'

[Transactions]
Enabled = true
ForwardersEnabled = false
MaxInFlight = 16
MaxQueued = 250
ReaperInterval = '1h0m0s'
ReaperThreshold = '168h0m0s'
ResendAfterThreshold = '1m0s'

[Transactions.AutoPurge]
Enabled = false

[Transactions.TransactionManagerV2]
Enabled = false

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
EstimateLimit = false
BumpMin = '5 gwei'
BumpPercent = 20
BumpThreshold = 3
EIP1559DynamicFees = false
FeeCapDefault = '100 gwei'
TipCapDefault = '1 wei'
TipCapMin = '1 wei'

[GasEstimator.BlockHistory]
BatchSize = 25
BlockHistorySize = 8
CheckInclusionBlocks = 12
CheckInclusionPercentile = 90
TransactionPercentile = 60

[GasEstimator.FeeHistory]
CacheTimeout = '10s'

[HeadTracker]
HistoryDepth = 100
MaxBufferSize = 3
SamplingInterval = '1s'
MaxAllowedFinalityDepth = 10000
FinalityTagBypass = true
PersistenceEnabled = true

[NodePool]
PollFailureThreshold = 5
PollInterval = '10s'
SelectionMode = 'HighestHead'
SyncThreshold = 5
LeaseDuration = '0s'
NodeIsSyncingEnabled = false
FinalizedBlockPollInterval = '5s'
EnforceRepeatableRead = true
DeathDeclarationDelay = '1m0s'
NewHeadsPollInterval = '0s'

[OCR]
ContractConfirmations = 4
ContractTransmitterTransmitTimeout = '10s'
DatabaseTimeout = '10s'
DeltaCOverride = '168h0m0s'
DeltaCJitterOverride = '1h0m0s'
ObservationGracePeriod = '1s'

[OCR2]
[OCR2.Automation]
GasLimit = 5400000

[Workflow]
GasLimitDefault = 400000
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
ChainType = 'arbitrum' # Example
```
ChainType is automatically detected from chain ID. Set this to force a certain chain type regardless of chain ID.
Available types: `arbitrum`, `celo`, `gnosis`, `hedera`, `kroma`, `metis`, `optimismBedrock`, `scroll`, `wemix`, `xlayer`, `zksync`

### FinalityDepth
```toml
FinalityDepth = 50 # Default
```
FinalityDepth is the number of blocks after which an ethereum transaction is considered "final". Note that the default is automatically set based on chain ID, so it should not be necessary to change this under normal operation.
BlocksConsideredFinal determines how deeply we look back to ensure that transactions are confirmed onto the longest chain
There is not a large performance penalty to setting this relatively high (on the order of hundreds)
It is practically limited by the number of heads we store in the database and should be less than this with a comfortable margin.
If a transaction is mined in a block more than this many blocks ago, and is reorged out, we will NOT retransmit this transaction and undefined behaviour can occur including gaps in the nonce sequence that require manual intervention to fix.
Therefore, this number represents a number of blocks we consider large enough that no re-org this deep will ever feasibly happen.

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

### FinalityTagEnabled
```toml
FinalityTagEnabled = false # Default
```
FinalityTagEnabled means that the chain supports the finalized block tag when querying for a block. If FinalityTagEnabled is set to true for a chain, then FinalityDepth field is ignored.
Finality for a block is solely defined by the finality related tags provided by the chain's RPC API. This is a placeholder and hasn't been implemented yet.

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
LogBackfillBatchSize = 1000 # Default
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

### LogPrunePageSize
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
LogPrunePageSize = 0 # Default
```
LogPrunePageSize defines size of the page for pruning logs. Controls how many logs/blocks (at most) are deleted in a single prune tick. Default value 0 means no paging, delete everything at once.

### BackupLogPollerBlockDelay
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
BackupLogPollerBlockDelay = 100 # Default
```
BackupLogPollerBlockDelay works in conjunction with Feature.LogPoller. Controls the block delay of Backup LogPoller, affecting how far behind the latest finalized block it starts and how often it runs.
BackupLogPollerDelay=0 will disable Backup LogPoller (_not recommended for production environment_).

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
RPCDefaultBatchSize = 250 # Default
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

### FinalizedBlockOffset
```toml
FinalizedBlockOffset = 0 # Default
```
FinalizedBlockOffset defines the number of blocks by which the latest finalized block will be shifted/delayed.
For example, suppose RPC returns block 100 as the latest finalized. In that case, the CL Node will treat block `100 - FinalizedBlockOffset` as the latest finalized block and `latest - FinalityDepth - FinalizedBlockOffset` in case of `FinalityTagEnabled = false.`
With `EnforceRepeatableRead = true,` RPC is considered healthy only if its most recent finalized block is larger or equal to the highest finalized block observed by the CL Node minus `FinalizedBlockOffset.`
Higher values of `FinalizedBlockOffset` with `EnforceRepeatableRead = true` reduce the number of false `FinalizedBlockOutOfSync` declarations on healthy RPCs that are slightly lagging behind due to network delays.
This may increase the number of healthy RPCs and reduce the probability that the CL Node will not have any healthy alternatives to the active RPC.
CAUTION: Setting this to values higher than 0 may delay transaction creation in products (e.g., CCIP, Automation) that base their decision on finalized on-chain events.
PoS chains with `FinalityTagEnabled=true` and batched (epochs) blocks finalization (e.g., Ethereum Mainnet) must be treated with special care as a minor increase in the `FinalizedBlockOffset` may lead to significant delays.
For example, let's say that `FinalizedBlockOffset = 1` and blocks are finalized in batches of 32.
The latest finalized block on chain is 64, so block 63 is the latest finalized for CL Node.
Block 64 will be treated as finalized by CL Node only when chain's latest finalized block is 65. As chain finalizes blocks in batches of 32,
CL Node has to wait for a whole new batch to be finalized to treat block 64 as finalized.

### LogBroadcasterEnabled
```toml
LogBroadcasterEnabled = true # Default
```
LogBroadcasterEnabled is a feature flag for LogBroadcaster, by default it's true.

### NoNewFinalizedHeadsThreshold
```toml
NoNewFinalizedHeadsThreshold = '0' # Default
```
NoNewFinalizedHeadsThreshold controls how long to wait for new finalized block before `NodePool` marks rpc endpoints as
out-of-sync. Only applicable if `FinalityTagEnabled=true`

Set to zero to disable.

## EVM.Transactions
```toml
[EVM.Transactions]
Enabled = true # Default
ForwardersEnabled = false # Default
MaxInFlight = 16 # Default
MaxQueued = 250 # Default
ReaperInterval = '1h' # Default
ReaperThreshold = '168h' # Default
ResendAfterThreshold = '1m' # Default
```


### Enabled
```toml
Enabled = true # Default
```
Enabled is a feature flag for the Transaction Manager. This flag also enables or disables the gas estimator since it is dependent on the TXM to start it.

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

## EVM.Transactions.AutoPurge
```toml
[EVM.Transactions.AutoPurge]
Enabled = false # Default
DetectionApiUrl = 'https://example.api.io' # Example
Threshold = 5 # Example
MinAttempts = 3 # Example
```


### Enabled
```toml
Enabled = false # Default
```
Enabled enables or disables automatically purging transactions that have been idenitified as terminally stuck (will never be included on-chain). This feature is only expected to be used by ZK chains.

### DetectionApiUrl
```toml
DetectionApiUrl = 'https://example.api.io' # Example
```
DetectionApiUrl configures the base url of a custom endpoint used to identify terminally stuck transactions.

### Threshold
```toml
Threshold = 5 # Example
```
Threshold configures the number of blocks a transaction has to remain unconfirmed before it is evaluated for being terminally stuck. This threshold is only applied if there is no custom API to identify stuck transactions provided by the chain.

### MinAttempts
```toml
MinAttempts = 3 # Example
```
MinAttempts configures the minimum number of broadcasted attempts a transaction has to have before it is evaluated further for being terminally stuck. This threshold is only applied if there is no custom API to identify stuck transactions provided by the chain. Ensure the gas estimator configs take more bump attempts before reaching the configured max gas price.

## EVM.Transactions.TransactionManagerV2
```toml
[EVM.Transactions.TransactionManagerV2]
Enabled = false # Default
BlockTime = '10s' # Example
CustomURL = 'https://example.api.io' # Example
DualBroadcast = false # Example
```


### Enabled
```toml
Enabled = false # Default
```
Enabled enables TransactionManagerV2.

### BlockTime
```toml
BlockTime = '10s' # Example
```
BlockTime controls the frequency of the backfill loop of TransactionManagerV2.

### CustomURL
```toml
CustomURL = 'https://example.api.io' # Example
```
CustomURL configures the base url of a custom endpoint used by the ChainDualBroadcast chain type.

### DualBroadcast
```toml
DualBroadcast = false # Example
```
DualBroadcast enables DualBroadcast functionality.

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
EstimateLimit = false # Default
BumpMin = '5 gwei' # Default
BumpPercent = 20 # Default
BumpThreshold = 3 # Default
BumpTxDepth = 16 # Example
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
- `L2Suggested` mode is deprecated and replaced with `SuggestedPrice`.
- `SuggestedPrice` is a mode which uses the gas price suggested by the rpc endpoint via `eth_gasPrice`.
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

This factor is always applied, so includes L2 transactions which uses a default gas limit of 1 and is also applied to `LimitDefault`.

### LimitTransfer
```toml
LimitTransfer = 21_000 # Default
```
LimitTransfer is the gas limit used for an ordinary ETH transfer.

### EstimateLimit
```toml
EstimateLimit = false # Default
```
EstimateLimit enables estimating gas limits for transactions. This feature respects the gas limit provided during transaction creation as an upper bound.

### BumpMin
```toml
BumpMin = '5 gwei' # Default
```
BumpMin is the minimum fixed amount of wei by which gas is bumped on each transaction attempt.

### BumpPercent
```toml
BumpPercent = 20 # Default
```
BumpPercent is the percentage by which to bump gas on a transaction that has exceeded `BumpThreshold`. The larger of `BumpPercent` and `BumpMin` is taken for gas bumps.

The `SuggestedPriceEstimator` adds the larger of `BumpPercent` and `BumpMin` on top of the price provided by the RPC when bumping a transaction's gas.

### BumpThreshold
```toml
BumpThreshold = 3 # Default
```
BumpThreshold is the number of blocks to wait for a transaction stuck in the mempool before automatically bumping the gas price. Set to 0 to disable gas bumping completely.

### BumpTxDepth
```toml
BumpTxDepth = 16 # Example
```
BumpTxDepth is the number of transactions to gas bump starting from oldest. Set to 0 for no limit (i.e. bump all). Can not be greater than EVM.Transactions.MaxInFlight. If not set, defaults to EVM.Transactions.MaxInFlight.

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

- Increase tipcap by `max(tipcap * (1 + BumpPercent), tipcap + BumpMin)`
- Increase feecap by `max(feecap * (1 + BumpPercent), feecap + BumpMin)`

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

(Only applies to EIP-1559 transactions)

## EVM.GasEstimator.DAOracle
```toml
[EVM.GasEstimator.DAOracle]
OracleType = 'opstack' # Example
OracleAddress = '0x420000000000000000000000000000000000000F' # Example
CustomGasPriceCalldata = '' # Default
```


### OracleType
```toml
OracleType = 'opstack' # Example
```
OracleType refers to the oracle family this config belongs to. Currently the available oracle types are: 'opstack', 'arbitrum', 'zksync', and 'custom_calldata'.

### OracleAddress
```toml
OracleAddress = '0x420000000000000000000000000000000000000F' # Example
```
OracleAddress is the address of the oracle contract.

### CustomGasPriceCalldata
```toml
CustomGasPriceCalldata = '' # Default
```
CustomGasPriceCalldata is optional and can be set to call a custom gas price function at the given OracleAddress.

## EVM.GasEstimator.LimitJobType
```toml
[EVM.GasEstimator.LimitJobType]
OCR = 100_000 # Example
OCR2 = 100_000 # Example
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

### OCR2
```toml
OCR2 = 100_000 # Example
```
OCR2 overrides LimitDefault for OCR2 jobs.

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
BatchSize = 25 # Default
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
BatchSize = 25 # Default
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

(Only applies to EIP-1559 transactions)

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

## EVM.GasEstimator.FeeHistory
```toml
[EVM.GasEstimator.FeeHistory]
CacheTimeout = '10s' # Default
```


### CacheTimeout
```toml
CacheTimeout = '10s' # Default
```
CacheTimeout is the time to wait in order to refresh the cached values stored in the FeeHistory estimator. A small jitter is applied so the timeout won't be exactly the same each time.

You want this value to be close to the block time. For slower chains, like Ethereum, you can set it to 12s, the same as the block time. For faster chains you can skip a block or two
and set it to two times the block time i.e. on Optimism you can set it to 4s. Ideally, you don't want to go lower than 1s since the RTT times of the RPC requests will be comparable to
the timeout. The estimator is already adding a buffer to account for a potential increase in prices within one or two blocks. On the other hand, slower frequency will fail to refresh
the prices and end up in stale values.

## EVM.HeadTracker
```toml
[EVM.HeadTracker]
HistoryDepth = 100 # Default
MaxBufferSize = 3 # Default
SamplingInterval = '1s' # Default
FinalityTagBypass = true # Default
MaxAllowedFinalityDepth = 10000 # Default
PersistenceEnabled = true # Default
```
The head tracker continually listens for new heads from the chain.

In addition to these settings, it log warnings if `EVM.NoNewHeadsThreshold` is exceeded without any new blocks being emitted.

### HistoryDepth
```toml
HistoryDepth = 100 # Default
```
HistoryDepth tracks the top N blocks on top of the latest finalized block to keep in the `heads` database table.
Note that this can easily result in MORE than `N + finality depth`  records since in the case of re-orgs we keep multiple heads for a particular block height.
Higher values help reduce number of RPC requests performed by TXM's Finalizer and improve TXM's Confirmer reorg protection on restarts.
At the same time, setting the value too high could lead to higher CPU consumption. The following formula could be used to calculate the optimal value: `expected_downtime_on_restart/block_time`.

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

### FinalityTagBypass
```toml
FinalityTagBypass = true # Default
```
FinalityTagBypass disables FinalityTag support in HeadTracker and makes it track blocks up to FinalityDepth from the most recent head.
It should only be used on chains with an extremely large actual finality depth (the number of blocks between the most recent head and the latest finalized block).
Has no effect if `FinalityTagsEnabled` = false

### MaxAllowedFinalityDepth
```toml
MaxAllowedFinalityDepth = 10000 # Default
```
MaxAllowedFinalityDepth - defines maximum number of blocks between the most recent head and the latest finalized block.
If actual finality depth exceeds this number, HeadTracker aborts backfill and returns an error.
Has no effect if `FinalityTagsEnabled` = false

### PersistenceEnabled
```toml
PersistenceEnabled = true # Default
```
PersistenceEnabled defines whether HeadTracker needs to store heads in the database.
Persistence is helpful on chains with large finality depth, where fetching blocks from the latest to the latest finalized takes a lot of time.
On chains with fast finality, the persistence layer does not improve the chain's load time and only consumes database resources (mainly IO).
NOTE: persistence should not be disabled for products that use LogBroadcaster, as it might lead to missed on-chain events.

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
LeaseDuration = '0s' # Default
NodeIsSyncingEnabled = false # Default
FinalizedBlockPollInterval = '5s' # Default
EnforceRepeatableRead = true # Default
DeathDeclarationDelay = '1m' # Default
NewHeadsPollInterval = '0s' # Default
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
- PriorityLevel: use the node with the smallest order number
- TotalDifficulty: use the node with the greatest total difficulty

### SyncThreshold
```toml
SyncThreshold = 5 # Default
```
SyncThreshold controls how far a node may lag behind the best node before being marked out-of-sync.
Depending on `SelectionMode`, this represents a difference in the number of blocks (`HighestHead`, `RoundRobin`, `PriorityLevel`), or total difficulty (`TotalDifficulty`).

Set to 0 to disable this check.

### LeaseDuration
```toml
LeaseDuration = '0s' # Default
```
LeaseDuration is the minimum duration that the selected "best" node (as defined by SelectionMode) will be used,
before switching to a better one if available. It also controls how often the lease check is done.
Setting this to a low value (under 1m) might cause RPC to switch too aggressively.
Recommended value is over 5m

Set to '0s' to disable

### NodeIsSyncingEnabled
```toml
NodeIsSyncingEnabled = false # Default
```
NodeIsSyncingEnabled is a flag that enables `syncing` health check on each reconnection to an RPC.
Node transitions and remains in `Syncing` state while RPC signals this state (In case of Ethereum `eth_syncing` returns anything other than false).
All of the requests to node in state `Syncing` are rejected.

Set true to enable this check

### FinalizedBlockPollInterval
```toml
FinalizedBlockPollInterval = '5s' # Default
```
FinalizedBlockPollInterval controls how often to poll RPC for new finalized blocks.
The finalized block is only used to report to the `pool_rpc_node_highest_finalized_block` metric. We plan to use it
in RPCs health assessment in the future.
If `FinalityTagEnabled = false`, poll is not performed and `pool_rpc_node_highest_finalized_block` is
reported based on latest block and finality depth.

Set to 0 to disable.

### EnforceRepeatableRead
```toml
EnforceRepeatableRead = true # Default
```
EnforceRepeatableRead defines if Core should only use RPCs whose most recently finalized block is greater or equal to
`highest finalized block - FinalizedBlockOffset`. In other words, exclude RPCs lagging on latest finalized
block.

Set false to disable

### DeathDeclarationDelay
```toml
DeathDeclarationDelay = '1m' # Default
```
DeathDeclarationDelay defines the minimum duration an RPC must be in unhealthy state before producing an error log message.
Larger values might be helpful to reduce the noisiness of health checks like `EnforceRepeatableRead = true', which might be falsely
trigger declaration of `FinalizedBlockOutOfSync` due to insignificant network delays in broadcasting of the finalized state among RPCs.
Should be greater than `FinalizedBlockPollInterval`.
Unhealthy RPC will not be picked to handle a request even if this option is set to a nonzero value.

### NewHeadsPollInterval
```toml
NewHeadsPollInterval = '0s' # Default
```
NewHeadsPollInterval define an interval for polling new block periodically using http client rather than subscribe to ws feed

Set to 0 to disable.

## EVM.NodePool.Errors
:warning: **_ADVANCED_**: _Do not change these settings unless you know what you are doing._
```toml
[EVM.NodePool.Errors]
NonceTooLow = '(: |^)nonce too low' # Example
NonceTooHigh = '(: |^)nonce too high' # Example
ReplacementTransactionUnderpriced = '(: |^)replacement transaction underpriced' # Example
LimitReached = '(: |^)limit reached' # Example
TransactionAlreadyInMempool = '(: |^)transaction already in mempool' # Example
TerminallyUnderpriced = '(: |^)terminally underpriced' # Example
InsufficientEth = '(: |^)insufficeint eth' # Example
TxFeeExceedsCap = '(: |^)tx fee exceeds cap' # Example
L2FeeTooLow = '(: |^)l2 fee too low' # Example
L2FeeTooHigh = '(: |^)l2 fee too high' # Example
L2Full = '(: |^)l2 full' # Example
TransactionAlreadyMined = '(: |^)transaction already mined' # Example
Fatal = '(: |^)fatal' # Example
ServiceUnavailable = '(: |^)service unavailable' # Example
TooManyResults = '(: |^)too many results' # Example
```
Errors enable the node to provide custom regex patterns to match against error messages from RPCs.

### NonceTooLow
```toml
NonceTooLow = '(: |^)nonce too low' # Example
```
NonceTooLow is a regex pattern to match against nonce too low errors.

### NonceTooHigh
```toml
NonceTooHigh = '(: |^)nonce too high' # Example
```
NonceTooHigh is a regex pattern to match against nonce too high errors.

### ReplacementTransactionUnderpriced
```toml
ReplacementTransactionUnderpriced = '(: |^)replacement transaction underpriced' # Example
```
ReplacementTransactionUnderpriced is a regex pattern to match against replacement transaction underpriced errors.

### LimitReached
```toml
LimitReached = '(: |^)limit reached' # Example
```
LimitReached is a regex pattern to match against limit reached errors.

### TransactionAlreadyInMempool
```toml
TransactionAlreadyInMempool = '(: |^)transaction already in mempool' # Example
```
TransactionAlreadyInMempool is a regex pattern to match against transaction already in mempool errors.

### TerminallyUnderpriced
```toml
TerminallyUnderpriced = '(: |^)terminally underpriced' # Example
```
TerminallyUnderpriced is a regex pattern to match against terminally underpriced errors.

### InsufficientEth
```toml
InsufficientEth = '(: |^)insufficeint eth' # Example
```
InsufficientEth is a regex pattern to match against insufficient eth errors.

### TxFeeExceedsCap
```toml
TxFeeExceedsCap = '(: |^)tx fee exceeds cap' # Example
```
TxFeeExceedsCap is a regex pattern to match against tx fee exceeds cap errors.

### L2FeeTooLow
```toml
L2FeeTooLow = '(: |^)l2 fee too low' # Example
```
L2FeeTooLow is a regex pattern to match against l2 fee too low errors.

### L2FeeTooHigh
```toml
L2FeeTooHigh = '(: |^)l2 fee too high' # Example
```
L2FeeTooHigh is a regex pattern to match against l2 fee too high errors.

### L2Full
```toml
L2Full = '(: |^)l2 full' # Example
```
L2Full is a regex pattern to match against l2 full errors.

### TransactionAlreadyMined
```toml
TransactionAlreadyMined = '(: |^)transaction already mined' # Example
```
TransactionAlreadyMined is a regex pattern to match against transaction already mined errors.

### Fatal
```toml
Fatal = '(: |^)fatal' # Example
```
Fatal is a regex pattern to match against fatal errors.

### ServiceUnavailable
```toml
ServiceUnavailable = '(: |^)service unavailable' # Example
```
ServiceUnavailable is a regex pattern to match against service unavailable errors.

### TooManyResults
```toml
TooManyResults = '(: |^)too many results' # Example
```
TooManyResults is a regex pattern to match an eth_getLogs error indicating the result set is too large to return

## EVM.OCR
```toml
[EVM.OCR]
ContractConfirmations = 4 # Default
ContractTransmitterTransmitTimeout = '10s' # Default
DatabaseTimeout = '10s' # Default
DeltaCOverride = "168h" # Default
DeltaCJitterOverride = "1h" # Default
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

### DeltaCOverride
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DeltaCOverride = "168h" # Default
```
DeltaCOverride (and `DeltaCJitterOverride`) determine the config override DeltaC.
DeltaC is the maximum age of the latest report in the contract. If the maximum age is exceeded, a new report will be
created by the report generation protocol.

### DeltaCJitterOverride
:warning: **_ADVANCED_**: _Do not change this setting unless you know what you are doing._
```toml
DeltaCJitterOverride = "1h" # Default
```
DeltaCJitterOverride is the range for jitter to add to `DeltaCOverride`.

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
Order = 100 # Default
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
WSURL is the WS(S) endpoint for this node. Required for primary nodes when `LogBroadcasterEnabled` is `true`

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

### Order
```toml
Order = 100 # Default
```
Order of the node in the pool, will takes effect if `SelectionMode` is `PriorityLevel` or will be used as a tie-breaker for `HighestHead` and `TotalDifficulty`

## EVM.OCR2.Automation
```toml
[EVM.OCR2.Automation]
GasLimit = 5400000 # Default
```


### GasLimit
```toml
GasLimit = 5400000 # Default
```
GasLimit controls the gas limit for transmit transactions from ocr2automation job.

## EVM.Workflow
```toml
[EVM.Workflow]
FromAddress = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
ForwarderAddress = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
GasLimitDefault = 400_000 # Default
```


### FromAddress
```toml
FromAddress = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
```
FromAddress is Address of the transmitter key to use for workflow writes.

### ForwarderAddress
```toml
ForwarderAddress = '0x2a3e23c6f242F5345320814aC8a1b4E58707D292' # Example
```
ForwarderAddress is the keystone forwarder contract address on chain.

### GasLimitDefault
```toml
GasLimitDefault = 400_000 # Default
```
GasLimitDefault is the default gas limit for workflow transactions.

## Cosmos
```toml
[[Cosmos]]
ChainID = 'Malaga-420' # Example
Enabled = true # Default
Bech32Prefix = 'wasm' # Default
BlockRate = '6s' # Default
BlocksUntilTxTimeout = 30 # Default
ConfirmPollPeriod = '1s' # Default
FallbackGasPrice = '0.015' # Default
GasToken = 'ucosm' # Default
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

### Bech32Prefix
```toml
Bech32Prefix = 'wasm' # Default
```
Bech32Prefix is the human-readable prefix for addresses on this Cosmos chain. See https://docs.cosmos.network/v0.47/spec/addresses/bech32.

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

### FallbackGasPrice
```toml
FallbackGasPrice = '0.015' # Default
```
FallbackGasPrice sets a fallback gas price to use when the estimator is not available.

### GasToken
```toml
GasToken = 'ucosm' # Default
```
GasToken is the token denomination which is being used to pay gas fees on this chain.

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
TxExpirationRebroadcast = false # Default
TxRetentionTimeout = '0s' # Default
SkipPreflight = true # Default
Commitment = 'confirmed' # Default
MaxRetries = 0 # Default
FeeEstimatorMode = 'fixed' # Default
ComputeUnitPriceMax = 1000 # Default
ComputeUnitPriceMin = 0 # Default
ComputeUnitPriceDefault = 0 # Default
FeeBumpPeriod = '3s' # Default
BlockHistoryPollPeriod = '5s' # Default
BlockHistorySize = 1 # Default
ComputeUnitLimitDefault = 200_000 # Default
EstimateComputeUnitLimit = false # Default
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

### TxExpirationRebroadcast
```toml
TxExpirationRebroadcast = false # Default
```
TxExpirationRebroadcast enables or disables transaction rebroadcast if expired. Expiration check is performed every `ConfirmPollPeriod`
A transaction is considered expired if the blockhash it was sent with is 150 blocks older than the latest blockhash.

### TxRetentionTimeout
```toml
TxRetentionTimeout = '0s' # Default
```
TxRetentionTimeout is the duration to retain transactions in storage after being marked as finalized or errored. Set to 0 to immediately drop transactions.

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

### BlockHistoryPollPeriod
```toml
BlockHistoryPollPeriod = '5s' # Default
```
BlockHistoryPollPeriod is the rate to poll for blocks in the block history fee estimator

### BlockHistorySize
```toml
BlockHistorySize = 1 # Default
```
BlockHistorySize is the number of blocks to take into consideration when using FeeEstimatorMode = 'blockhistory' to determine compute unit price.
If set to 1, the compute unit price will be determined by the median of the last block's compute unit prices.
If set N > 1, the compute unit price will be determined by the average of the medians of the last N blocks' compute unit prices.
DISCLAIMER: 1:1 ratio between n and RPC calls. It executes once every 'BlockHistoryPollPeriod' value.

### ComputeUnitLimitDefault
```toml
ComputeUnitLimitDefault = 200_000 # Default
```
ComputeUnitLimitDefault is the compute units limit applied to transactions unless overriden during the txm enqueue

### EstimateComputeUnitLimit
```toml
EstimateComputeUnitLimit = false # Default
```
EstimateComputeUnitLimit enables or disables compute unit limit estimations per transaction. If estimations return 0 used compute, the ComputeUnitLimitDefault value is used, if set.

## Solana.MultiNode
```toml
[Solana.MultiNode]
Enabled = false # Default
PollFailureThreshold = 5 # Default
PollInterval = '10s' # Default
SelectionMode = 'PriorityLevel' # Default
SyncThreshold = 5 # Default
NodeIsSyncingEnabled = false # Default
LeaseDuration = '1m0s' # Default
NewHeadsPollInterval = '10s' # Default
FinalizedBlockPollInterval = '10s' # Default
EnforceRepeatableRead = true # Default
DeathDeclarationDelay = '10s' # Default
NodeNoNewHeadsThreshold = '10s' # Default
NoNewFinalizedHeadsThreshold = '10s' # Default
FinalityDepth = 0 # Default
FinalityTagEnabled = true # Default
FinalizedBlockOffset = 0 # Default
```


### Enabled
```toml
Enabled = false # Default
```
Enabled enables the multinode feature.

### PollFailureThreshold
```toml
PollFailureThreshold = 5 # Default
```
PollFailureThreshold is the number of consecutive poll failures before a node is considered unhealthy.

### PollInterval
```toml
PollInterval = '10s' # Default
```
PollInterval is the rate to poll for node health.

### SelectionMode
```toml
SelectionMode = 'PriorityLevel' # Default
```
SelectionMode is the method used to select the next best node to use.

### SyncThreshold
```toml
SyncThreshold = 5 # Default
```
SyncThreshold is the number of blocks behind the best node that a node can be before it is considered out of sync.

### NodeIsSyncingEnabled
```toml
NodeIsSyncingEnabled = false # Default
```
NodeIsSyncingEnabled enables the feature to avoid sending transactions to nodes that are syncing. Not relavant for Solana.

### LeaseDuration
```toml
LeaseDuration = '1m0s' # Default
```
LeaseDuration is the max duration a node can be leased for.

### NewHeadsPollInterval
```toml
NewHeadsPollInterval = '10s' # Default
```
NewHeadsPollInterval is the rate to poll for new heads.

### FinalizedBlockPollInterval
```toml
FinalizedBlockPollInterval = '10s' # Default
```
FinalizedBlockPollInterval is the rate to poll for the finalized block.

### EnforceRepeatableRead
```toml
EnforceRepeatableRead = true # Default
```
EnforceRepeatableRead enforces the repeatable read guarantee for multinode.

### DeathDeclarationDelay
```toml
DeathDeclarationDelay = '10s' # Default
```
DeathDeclarationDelay is the duration to wait before declaring a node dead.

### NodeNoNewHeadsThreshold
```toml
NodeNoNewHeadsThreshold = '10s' # Default
```
NodeNoNewHeadsThreshold is the duration to wait before declaring a node unhealthy due to no new heads.

### NoNewFinalizedHeadsThreshold
```toml
NoNewFinalizedHeadsThreshold = '10s' # Default
```
NoNewFinalizedHeadsThreshold is the duration to wait before declaring a node unhealthy due to no new finalized heads.

### FinalityDepth
```toml
FinalityDepth = 0 # Default
```
FinalityDepth is not used when finality tags are enabled.

### FinalityTagEnabled
```toml
FinalityTagEnabled = true # Default
```
FinalityTagEnabled enables the use of finality tags.

### FinalizedBlockOffset
```toml
FinalizedBlockOffset = 0 # Default
```
FinalizedBlockOffset is the offset from the finalized block to use for finality tags.

## Solana.Nodes
```toml
[[Solana.Nodes]]
Name = 'primary' # Example
URL = 'http://solana.web' # Example
SendOnly = false # Default
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

### SendOnly
```toml
SendOnly = false # Default
```
SendOnly is a multinode config that only sends transactions to a node and does not read state

## Starknet
```toml
[[Starknet]]
ChainID = 'foobar' # Example
FeederURL = 'http://feeder.url' # Example
Enabled = true # Default
OCR2CachePollPeriod = '5s' # Default
OCR2CacheTTL = '1m' # Default
RequestTimeout = '10s' # Default
TxTimeout = '10s' # Default
ConfirmationPoll = '5s' # Default
```


### ChainID
```toml
ChainID = 'foobar' # Example
```
ChainID is the Starknet chain ID.

### FeederURL
```toml
FeederURL = 'http://feeder.url' # Example
```
FeederURL is required to get tx metadata (that the RPC can't)

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
TxTimeout = '10s' # Default
```
TxTimeout is the timeout for sending txes to an RPC endpoint.

### ConfirmationPoll
```toml
ConfirmationPoll = '5s' # Default
```
ConfirmationPoll is how often to confirmer checks for tx inclusion on chain.

## Starknet.Nodes
```toml
[[Starknet.Nodes]]
Name = 'primary' # Example
URL = 'http://stark.node' # Example
APIKey = 'key' # Example
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

### APIKey
```toml
APIKey = 'key' # Example
```
APIKey Header is optional and only required for Nethermind RPCs

