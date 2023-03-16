[//]: # (Documentation generated from docs/secrets.toml - DO NOT EDIT.)

This document describes the TOML format for secrets.

Each secret has an alternative corresponding environment variable.

See also [CONFIG.md](CONFIG.md)

## Example

```toml
[Database]
URL = 'postgresql://user:pass@localhost:5432/dbname?sslmode=disable' # Required

[Password]
Keystore = 'keystore_pass' # Required
```

## Database
```toml
[Database]
URL = "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" # Example
BackupURL = "postgresql://user:pass@read-replica.example.com:5432/dbname?sslmode=disable" # Example
AllowSimplePasswords = false # Default
```


### URL
```toml
URL = "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" # Example
```
URL is the PostgreSQL URI to connect to your database. Chainlink nodes require Postgres versions >= 11. See
[Running a Chainlink Node](https://docs.chain.link/docs/running-a-chainlink-node/#set-the-remote-database_url-config) for an example.

Environment variable: `CL_DATABASE_URL`

### BackupURL
```toml
BackupURL = "postgresql://user:pass@read-replica.example.com:5432/dbname?sslmode=disable" # Example
```
BackupURL is where the automatic database backup will pull from, rather than the main DATABASE_URL. It is recommended
to set this value to a read replica if you have one to avoid excessive load on the main database.

Environment variable: `CL_DATABASE_BACKUP_URL`

### AllowSimplePasswords
```toml
AllowSimplePasswords = false # Default
```
AllowSimplePasswords skips the password complexity check normally enforced on URL & BackupURL.

Environment variable: `CL_DATABASE_ALLOW_SIMPLE_PASSWORDS`

## Explorer
```toml
[Explorer]
AccessKey = "access_key" # Example
Secret = "secret" # Example
```


### AccessKey
```toml
AccessKey = "access_key" # Example
```
AccessKey is the access key for authenticating with the Explorer.

Environment variable: `CL_EXPLORER_ACCESS_KEY`

### Secret
```toml
Secret = "secret" # Example
```
Secret is the secret for authenticating with the Explorer.

Environment variable: `CL_EXPLORER_SECRET`

## Password
```toml
[Password]
Keystore = "keystore_pass" # Example
VRF = "VRF_pass" # Example
```


### Keystore
```toml
Keystore = "keystore_pass" # Example
```
Keystore is the password for the node's account.

Environment variable: `CL_PASSWORD_KEYSTORE`

### VRF
```toml
VRF = "VRF_pass" # Example
```
VRF is the password for the vrf keys.

Environment variable: `CL_PASSWORD_VRF`

## Pyroscope
```toml
[Pyroscope]
AuthToken = "pyroscope-token" # Example
```


### AuthToken
```toml
AuthToken = "pyroscope-token" # Example
```
AuthToken is the API key for the Pyroscope server.

Environment variable: `CL_PYROSCOPE_AUTH_TOKEN`

## Mercury
```toml
[Mercury]
```
Mercury credentials are needed if running OCR2 jobs in mercury mode. 0 or
more Mercury credentials may be specified. URLs must be unique.

## Mercury.Credentials
```toml
[[Mercury.Credentials]]
URL = "http://example.com/reports" # Example
Username = "exampleusername" # Example
Password = "examplepassword" # Example
```


### URL
```toml
URL = "http://example.com/reports" # Example
```
URL is the URL of the mercury endpoint

### Username
```toml
Username = "exampleusername" # Example
```
Username is used for basic auth with the mercury endpoint

### Password
```toml
Password = "examplepassword" # Example
```
Password is used for basic auth with the mercury endpoint

