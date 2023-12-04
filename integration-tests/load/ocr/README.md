# OCR Load tests

## Setup

These tests can connect to any cluster create with [chainlink-cluster](../../../charts/chainlink-cluster/README.md)

Create your cluster

```sh
kubectl create ns my-cluster
devspace use namespace my-cluster
devspace deploy
sudo kubefwd svc -n my-cluster
```

Change environment connection configuration [here](connection.toml)

If you haven't changed anything in [devspace.yaml](../../../charts/chainlink-cluster/devspace.yaml) then default connection configuration will work

## Usage

```sh
export LOKI_TOKEN=...
export LOKI_URL=...

go test -v -run TestOCRLoad
go test -v -run TestOCRVolume
```

Check test configuration [here](config.toml)
