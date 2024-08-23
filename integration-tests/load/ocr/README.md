### OCR Load tests

## Setup CRIB
CRIB CORE user documentation is available in [CORE CRIB - Deploy & Access Instructions](https://smartcontract-it.atlassian.net/wiki/spaces/TT/pages/597197209/CORE+CRIB+-+Deploy+Access+Instructions)
```shell
devspace deploy --debug --profile local-dev-simulated-core-ocr1 --skip-build
```

## Usage

Create `overrides.toml` in this directory
```toml
[CRIB]
namespace = "$your_crib_namespace_here"
# only Geth is supported right now
network_name = "geth"
nodes = 5

[Logging.Loki]
tenant_id="promtail"
endpoint="..."
basic_auth_secret="..."
```
Run the tests

Set `K8S_STAGING_INGRESS_SUFFIX` when run locally (`export K8S_STAGING_INGRESS_SUFFIX=$(op read op://CRIB/secrets/K8S_STAGING_INGRESS_SUFFIX)`)

```
go test -v -run TestOCRLoad
go test -v -run TestOCRVolume
```