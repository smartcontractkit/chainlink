# Chainlink cluster

Example CL nodes cluster for system level tests

# Develop
## Install from release
Note: The setup below doesn't work at the moment.

Add the repository

```
helm repo add chainlink-cluster https://raw.githubusercontent.com/smartcontractkit/chainlink/helm-release/
helm repo update
```

Set default namespace

```
kubectl create ns cl-cluster
kubectl config set-context --current --namespace cl-cluster
```

Install

```
helm install -f values.yaml cl-cluster .
```

## Create a new release
Note: The setup below doesn't work at the moment.

Bump version in `Chart.yml` add your changes and add `helm_release` label to any PR to trigger a release

## Helm Test

```
helm test cl-cluster
```

## Uninstall

```
helm uninstall cl-cluster
```