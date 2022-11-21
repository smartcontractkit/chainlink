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

## Table of contents

- [Database](#Database)
- [Explorer](#Explorer)
- [Password](#Password)
- [Pyroscope](#Pyroscope)
- [Mercury](#Mercury)
	- [Credentials](#Mercury-Credentials)

## Database<a id='Database'></a>
```toml
[Database]
URL = "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" # Example
BackupURL = "postgresql://user:pass@read-replica.example.com:5432/dbname?sslmode=disable" # Example
AllowSimplePasswords = false # Default
```


### URL<a id='Database-URL'></a>
```toml
URL = "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" # Example
```
URL is the PostgreSQL URI to connect to your database. Chainlink nodes require Postgres versions >= 11. See
[Running a Chainlink Node](https://docs.chain.link/docs/running-a-chainlink-node/#set-the-remote-database_url-config) for an example.

Environment variable: `CL_DATABASE_URL`

### BackupURL<a id='Database-BackupURL'></a>
```toml
BackupURL = "postgresql://user:pass@read-replica.example.com:5432/dbname?sslmode=disable" # Example
```
BackupURL is where the automatic database backup will pull from, rather than the main DATABASE_URL. It is recommended
to set this value to a read replica if you have one to avoid excessive load on the main database.

Environment variable: `CL_DATABASE_BACKUP_URL`

### AllowSimplePasswords<a id='Database-AllowSimplePasswords'></a>
```toml
AllowSimplePasswords = false # Default
```
AllowSimplePasswords skips the password complexity check normally enforced on URL & BackupURL.

Environment variable: `CL_DATABASE_ALLOW_SIMPLE_PASSWORDS`

## Explorer<a id='Explorer'></a>
```toml
[Explorer]
AccessKey = "access_key" # Example
Secret = "secret" # Example
```


### AccessKey<a id='Explorer-AccessKey'></a>
```toml
AccessKey = "access_key" # Example
```
AccessKey is the access key for authenticating with the Explorer.

Environment variable: `CL_EXPLORER_ACCESS_KEY`

### Secret<a id='Explorer-Secret'></a>
```toml
Secret = "secret" # Example
```
Secret is the secret for authenticating with the Explorer.

Environment variable: `CL_EXPLORER_SECRET`

## Password<a id='Password'></a>
```toml
[Password]
Keystore = "keystore_pass" # Example
VRF = "VRF_pass" # Example
```


### Keystore<a id='Password-Keystore'></a>
```toml
Keystore = "keystore_pass" # Example
```
Keystore is the password for the node's account.

Environment variable: `CL_PASSWORD_KEYSTORE`

### VRF<a id='Password-VRF'></a>
```toml
VRF = "VRF_pass" # Example
```
VRF is the password for the vrf keys.

Environment variable: `CL_PASSWORD_VRF`

## Pyroscope<a id='Pyroscope'></a>
```toml
[Pyroscope]
AuthToken = "pyroscope-token" # Example
```


### AuthToken<a id='Pyroscope-AuthToken'></a>
```toml
AuthToken = "pyroscope-token" # Example
```
AuthToken is the API key for the Pyroscope server.

Environment variable: `CL_PYROSCOPE_AUTH_TOKEN`

## Mercury<a id='Mercury'></a>
```toml
[Mercury]
```
Mercury credentials are needed if running OCR2 jobs in mercury mode. 0 or
more Mercury credentials may be specified. URLs must be unique.

## Mercury.Credentials<a id='Mercury-Credentials'></a>
```toml
[[Mercury.Credentials]]
URL = "http://example.com/reports" # Example
Username = "exampleusername" # Example
Password = "examplepassword" # Example
```


### URL<a id='Mercury-Credentials-URL'></a>
```toml
URL = "http://example.com/reports" # Example
```
URL is the URL of the mercury endpoint

### Username<a id='Mercury-Credentials-Username'></a>
```toml
Username = "exampleusername" # Example
```
Username is used for basic auth with the mercury endpoint

### Password<a id='Mercury-Credentials-Password'></a>
```toml
Password = "examplepassword" # Example
```
Password is used for basic auth with the mercury endpoint

