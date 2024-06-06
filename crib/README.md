# Crib DevSpace Setup

CRIB is a DevSpace-based method to launch chainlink cluster for system level tests.

## Initial Setup

### Prerequisites
- **Kubefwd**: Ensure you have `kubefwd` installed for port forwarding. Although no nixpkg exists yet, you can install it via Homebrew:
    ```bash
    brew install txn2/tap/kubefwd
    ```
- **Docker**: If you plan to build images, ensure the [Docker](https://docs.docker.com/engine/install/) service is running.
- **VPN**: Ensure you are connected to the VPN to access the necessary resources.

### Dev Environment Setup
```bash
# Enter nix shell, which contains necessary tools
nix develop 
 
# Copy the environment example file and fill in the necessary values
cd crib/
cp .env.example .env
```

## CRIB Initialization

The CRIB initialization script sets up your development environment for deploying a Chainlink cluster. This script automates several tasks, such as configuring AWS profiles, updating Kubernetes configurations, and logging into Docker and Helm registries.

### Script Overview
The script performs the following actions:
1. **Environment Setup**: Sources the `.env` file containing necessary environment variables.
2. **AWS Profile Configuration**: Sets up or updates the AWS profile for ECR registry access.
3. **AWS Authentication**: Ensures AWS credentials are valid and logs into AWS SSO if needed.
4. **Kubernetes Configuration**: Updates the kubeconfig for EKS access and sets the Kubernetes context.
5. **Docker and Helm Registry Login**: Logs into AWS ECR to allow pulling and pushing of Docker images and Helm charts.

### Usage Instructions

Execute the script with the desired namespace as an argument:
```bash
./cribbit.sh crib-yournamespace
```
**Note**: The namespace must begin with `crib-` unless overridden.

### Troubleshooting

- **Missing Environment Variables**: If any environment variables are missing, the script will terminate early. Make sure all required variables are defined in your `.env` file.
- **AWS Credentials Not Detected**: If AWS credentials cannot be verified, the script will attempt to log in through SSO. Follow the prompts to complete this process.
- **Docker Daemon Not Running**: Ensure that the Docker service is running before executing the script. This is required for Docker and Helm registry logins.

## Cluster Configuration and Deployment

### Configuring New Clusters
We use [Devspace](https://www.devspace.sh/docs/getting-started/installation?x0=3) for cluster deployment. Review the settings in `deployments.app.helm.values` and [values.yaml](../charts/chainlink-cluster/values.yaml) for detailed configuration options.


### Usage Examples
```bash
# Deploy the current repository state
devspace deploy 

# Deploy a specific image tag already available in ECR
devspace deploy --override-image-tag "<image-tag>" 

# Deploy a public ECR image tag
DEVSPACE_IMAGE=public.ecr.aws/chainlink/chainlink devspace deploy --override-image-tag 2.9.0

# Set the time-to-live (TTL) for the namespace, once this time has passed the namespace along with all its associated resources will be deleted
# Valid values are `1h`, `2m`, `3s`, etc. Go time format is invalid `1h2m3s`
devspace run ttl ${namespace} 120hA # Default 72h

# Forward ports to check UI or run tests
devspace run connect ${my-personal-namespace-name-crib}

# List ingress hostnames
devspace run ingress-hosts

# Destroy the cluster
devspace purge
```

## Load Testing
Deploy the dashboard and run load tests as described in the [testing documentation](../integration-tests/load/ocr/README.md):

**NOTE:** If you used `devspace dev ...` always use `devspace reset pods` to switch the pods back

# Helm

If you would like to use `helm` directly, uncomment data in `values.yaml`

## Install from local files
Deploy the cluster with the following command:
```bash
helm install -f values.yaml cl-cluster .
```

In another terminal, forward all apps:
```bash
sudo kubefwd svc -n cl-cluster
```

Then you can connect and run your tests

# Grafana dashboard

We use the [Grabana](https://github.com/K-Phoen/grabana) library to create dashboards programmatically.

## Dashboard Platform Configuration
Dashboard platform selection can be done with `INFRA_PLATFORM` environment variable.

The available options are: 
  - `kubernetes` 
  - `docker`

## Panel Selection 
Non-default panel selection can be done with `PANELS_INCLUDED` environment variable.

A comma-separated list of panel names can be provided to include only those panels in the dashboard. 

## Dashboard Deployment
```
export LOKI_TENANT_ID=promtail
export LOKI_URL=...
export GRAFANA_URL=...
export GRAFANA_TOKEN=...
export PROMETHEUS_DATA_SOURCE_NAME=Thanos
export LOKI_DATA_SOURCE_NAME=Loki
export INFRA_PLATFORM=kubernetes
export GRAFANA_FOLDER=DashboardCoreDebug
export DASHBOARD_NAME=CL-Cluster

devspace run dashboard_deploy
```

Open Grafana folder `DashboardCoreDebug` and find dashboard `ChainlinkClusterDebug`.

## Load Testing

To deploy your dashboard and run load [tests](../../integration-tests/load/), see [the integration test README](../../integration-tests/README.md).

```
devspace run dashboard_deploy
devspace run workload
devspace run dashboard_test
```

## Local Testing

Go to [dashboard-lib](../dashboard-lib) and link the modules locally.

```
cd dashboard
pnpm link --global
cd crib/dashboard/tests
pnpm link --global dashboard-tests
```

Then run the tests with commands mentioned above
