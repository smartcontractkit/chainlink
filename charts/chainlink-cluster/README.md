# Chainlink cluster
Example CL nodes cluster for system level tests

Install `kubefwd` (no nixpkg for it yet, planned)
```
brew install txn2/tap/kubefwd
```

Enter the shell (from the root project dir)
```
nix develop
```

# Develop

## New cluster
We are using [devspace](https://www.devspace.sh/docs/getting-started/installation?x0=3)

Configure the cluster, see `deployments.app.helm.values` and [values.yaml](./values.yaml) comments

Set your registry for the image, example for `ECR`:
```
aws ecr get-login-password --region us-west-2 | docker login --username AWS --password-stdin ${aws_account}.dkr.ecr.us-west-2.amazonaws.com
export DEVSPACE_IMAGE="${aws_account}.dkr.ecr.us-west-2.amazonaws.com/chainlink-devspace"
```
Enter the shell and deploy
```
# set your unique namespace if it's a new cluster
devspace use namespace cl-cluster
devspace deploy
```
If you don't need a build use
```
devspace deploy --skip-build
```

Connect to your environment
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
If you need to update the whole cluster run `deploy` again with a new set of images
```
devspace reset pods
devspace deploy
```
Destroy the cluster
```
devspace purge
```

If you need to run some system level tests inside k8s use `runner` profile:
```
devspace dev -p runner
```

If you used `devspace dev ...` always use `devspace reset pods` to switch the pods back

## Debug existing cluster
If you need to debug CL node that is already deployed change `dev.app.container` and `dev.app.labelSelector` in [devspace.yaml](devspace.yaml) if they are not default and run:
```
devspace dev -p node
or
devspace dev -p runner
```

## Automatic file sync
When you run `devspace dev` your files described in `dev.app.sync` of [devspace.yaml](devspace.yaml) will be uploaded to the switched container

After that all the changes will be synced automatically

Check `.profiles` to understand what is uploaded in profiles `runner` and `node`

# Helm
If you would like to use `helm` directly, please uncomment data in `values-raw-helm.yaml`
## Install from local files
```
helm install -f values-raw-helm.yaml cl-cluster .
```
Forward all apps (in another terminal)
```
sudo kubefwd svc
```
Then you can connect and run your tests

## Install from release
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
helm install -f values-raw-helm.yaml cl-cluster chainlink-cluster/chainlink-cluster --version v0.1.2
```

## Create a new release
Bump version in `Chart.yml` add your changes and add `helm_release` label to any PR to trigger a release

## Helm Test
```
helm test cl-cluster
```

## Uninstall
```
helm uninstall cl-cluster
```