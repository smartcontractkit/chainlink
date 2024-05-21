### OCR Load tests

## Setup
These tests can connect to any cluster create with [chainlink-cluster](../../../charts/chainlink-cluster/README.md)

Create your cluster, if you already have one just use `kubefwd`
```
kubectl create ns cl-cluster
devspace use namespace cl-cluster
devspace deploy
sudo kubefwd svc -n cl-cluster
```

Change environment connection configuration [here](../../../charts/chainlink-cluster/connect.toml)

If you haven't changed anything in [devspace.yaml](../../crib/devspace.yaml) then default connection configuration will work

## Usage

```
export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestOCRLoad
go test -v -run TestOCRVolume
```

Check test configuration [here](config.toml)