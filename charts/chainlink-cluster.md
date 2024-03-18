# chainlink-cluster

## Chainlink cluster

Example CL nodes cluster for system level tests

Install `kubefwd` (no nixpkg for it yet, planned)

```
brew install txn2/tap/kubefwd
```

If you want to build images you need [docker](https://docs.docker.com/engine/install/) service running

Enter the shell (from the root project dir)

```
nix develop
```

## Develop

### New cluster

We are using [devspace](https://www.devspace.sh/docs/getting-started/installation?x0=3)

Configure the cluster, see `deployments.app.helm.values` and [values.yaml](chainlink-cluster/values.yaml) comments for more details

Configure your `cluster` setup (one time setup, internal usage only)

```
export DEVSPACE_IMAGE="..."
cd charts/chainlink-cluster
./setup.sh ${my-personal-namespace-name-crib}
```

Build and deploy current commit

```
devspace deploy
```

Default `ttl` is `72h`, use `ttl` command to update if you need more time

Valid values are `1h`, `2m`, `3s`, etc. Go time format is invalid `1h2m3s`

```
devspace run ttl ${namespace} 120h
```

If you don't need a build use

```
devspace deploy --skip-build
```

To deploy particular commit (must be in registry) use

```
devspace deploy --skip-build ${short_sha_of_image}
```

Forward ports to check UI or run tests

```
devspace run connect ${my-personal-namespace-name-crib}
```

Connect to your environment, by replacing container with label `node-1` with your local repository files

```
devspace dev -p node
make chainlink
make chainlink-local-start
```

Fix something in the code locally, it'd automatically sync, rebuild it inside container and run again

```
make chainlink
make chainlink-local-start
```

Reset the pod to original image

```
devspace reset pods
```

Destroy the cluster

```
devspace purge
```

### Running load tests

Check this [doc](../integration-tests/load/ocr.md)

If you used `devspace dev ...` always use `devspace reset pods` to switch the pods back

### Debug existing cluster

If you need to debug CL node that is already deployed change `dev.app.container` and `dev.app.labelSelector` in [devspace.yaml](chainlink-cluster/devspace.yaml) if they are not default and run:

```
devspace dev -p node
```

### Automatic file sync

When you run `devspace dev` your files described in `dev.app.sync` of [devspace.yaml](chainlink-cluster/devspace.yaml) will be uploaded to the switched container

After that all the changes will be synced automatically

Check `.profiles` to understand what is uploaded in profiles `runner` and `node`

## Helm

If you would like to use `helm` directly, please uncomment data in `values.yaml`

### Install from local files

```
helm install -f values.yaml cl-cluster .
```

Forward all apps (in another terminal)

```
sudo kubefwd svc -n cl-cluster
```

Then you can connect and run your tests

### Install from release

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

### Create a new release

Bump version in `Chart.yml` add your changes and add `helm_release` label to any PR to trigger a release

### Helm Test

```
helm test cl-cluster
```

### Uninstall

```
helm uninstall cl-cluster
```

## Grafana dashboard

We are using [Grabana](https://github.com/K-Phoen/grabana) lib to create dashboards programmatically

You can select `PANELS_INCLUDED`, options are `core`, `wasp`, comma separated

You can also select dashboard platform in `INFRA_PLATFORM` either `kubernetes` or `docker`

```
export LOKI_TENANT_ID=promtail
export LOKI_URL=...
export GRAFANA_URL=...
export GRAFANA_TOKEN=...
export PROMETHEUS_DATA_SOURCE_NAME=Thanos
export LOKI_DATA_SOURCE_NAME=Loki
export PANELS_INCLUDED=core,wasp
export INFRA_PLATFORM=kubernetes|docker
export GRAFANA_FOLDER=DashboardCoreDebug
export DASHBOARD_NAME=ChainlinkClusterDebug

go run dashboard/cmd/dashboard_deploy.go
```

Open Grafana folder `DashboardCoreDebug` and find dashboard `ChainlinkClusterDebug`
