# Chainlink cluster
Example CL nodes cluster for system level tests

Enter the shell
```
nix develop
```

# Develop

## New cluster
We are using [devspace](https://www.devspace.sh/docs/getting-started/installation?x0=3)

Configure the cluster, see `deployments.app.helm.values` and [values.yaml](./values.yaml) comments

`Caveat`: currently it's working only with `dockerhub`
```
nix develop
cd charts/chainlink-cluster

# set your unique namespace if it's a new cluster
devspace use namespace cl-cluster

export DEVSPACE_IMAGE="..." - any dockerhub registry name, ex. "registry/app", you should be logged in

devspace deploy
```

Connect to your environment
```
devspace dev
cd chainlink
make chainlink
make chainlink-local-start
```
Fix something in the code locally, it'd automatically sync, rebuild it inside container and run again
```
make chainlink
make chainlink-local-start
```
If you need to update the whole cluster run `deploy` again with a new set of images
```
devspace deploy
```
Destroy the cluster
```
devspace purge
```

## Debug existing cluster
If you need to debug CL node that is already deployed change `dev.app.container` and `dev.app.labelSelector` in [devspace.yaml](devspace.yaml) and run:
```
devspace dev
```

## Automatic file sync
When you run `devspace dev` your files described in `dev.app.sync` of [devspace.yaml](devspace.yaml) will be uploaded to the switched container

After that all the changes will be synced automatically

# Helm
## Install
```
helm install cl-cluster .
```

## Helm Test
```
helm test cl-cluster
```

## Uninstall
```
helm uninstall cl-cluster
```