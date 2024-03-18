# Running Smoke Tests Locally

pre-reqs
* install node.js (18 preferred)
* k3d

## TLDR;
run this once
```shell
echo '127.0.0.1 k3d-myregistry.localhost' | sudo tee -a /etc/hosts
```
run to create k3d resources
```shell
k3d registry create myregistry.localhost --port 5001
k3d cluster create --registry-use k3d-myregistry.localhost:5001
```
run to create new build
```shell
cd ~/go/src/github.com/chainlink
env DOCKER_DEFAULT_PLATFORM=linux/amd64 docker buildx build --platform linux/amd64 -f ./core/chainlink.Dockerfile --build-arg ENVIRONMENT=release --build-arg COMMIT_SHA=$(git rev-parse HEAD) -t smartcontract/chainlink:develop-$(git rev-parse HEAD) .
export CHAINLINK_VERSION=develop-$(git rev-parse HEAD)
docker tag docker.io/smartcontract/chainlink:$CHAINLINK_VERSION k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink:$CHAINLINK_VERSION
docker push k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink:$CHAINLINK_VERSION
export CHAINLINK_IMAGE=k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink
```
run the tests
 ```shell
cd ~/go/src/github.com/chainlink
make test_smoke_simulated args="--focus-file=auto_ocr_test.go"
 ``` 

## Already have the initial stuff set up and just want to rebuild and run
build+run
```shell
cd ~/go/src/github.com/chainlink
env DOCKER_DEFAULT_PLATFORM=linux/amd64 docker buildx build --platform linux/amd64 -f ./core/chainlink.Dockerfile --build-arg ENVIRONMENT=release --build-arg COMMIT_SHA=$(git rev-parse HEAD) -t smartcontract/chainlink:develop-$(git rev-parse HEAD) .
export CHAINLINK_VERSION=develop-$(git rev-parse HEAD)
export TEST_LOG_LEVEL="debug"
docker tag docker.io/smartcontract/chainlink:$CHAINLINK_VERSION k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink:$CHAINLINK_VERSION
docker push k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink:$CHAINLINK_VERSION
export CHAINLINK_IMAGE=k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink
make test_smoke_simulated args="--focus-file=auto_ocr_test.go"
 ``` 

## Step by Step

1. Build a Docker image of the chainlink repo:

   ```shell
   env DOCKER_DEFAULT_PLATFORM=linux/amd64 docker buildx build --platform linux/amd64 -f ./core/chainlink.Dockerfile --build-arg ENVIRONMENT=release --build-arg COMMIT_SHA=$(git rev-parse HEAD) -t smartcontract/chainlink:develop-$(git rev-parse HEAD) .
   ```
   last line of the output will have something like
   `=> => naming to docker.io/smartcontract/chainlink:develop-a4caf33ce0ed6b841294c5ef06563c1cd4de6dfc`
   use the tag at the end in the next command
   ```shell
   export CHAINLINK_VERSION=develop-a4caf33ce0ed6b841294c5ef06563c1cd4de6dfc
   ```
2. Set up a Kubernetes cluster locally and a Docker registry which Kubernetes can pull from:

   ```shell
   k3d registry create myregistry.localhost --port 5001
   k3d cluster create --registry-use k3d-myregistry.localhost:5001
   ```

3. Add these lines to the `/etc/hosts` file

    ```shell
    # Added for k3d registry
    127.0.0.1 k3d-myregistry.localhost
    ```

4. **Tag** the Docker image with the appropriate name:

   ```shell
   docker tag docker.io/smartcontract/chainlink:$CHAINLINK_VERSION k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink:$CHAINLINK_VERSION
   ```

5. **Push** the Docker image we created in step 1 to the new registry:

   ```shell 
   docker push k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink:$CHAINLINK_VERSION
   ```

6. Before actually running the tests, we need to set environment variables as follows:

   ```shell
   export CHAINLINK_IMAGE=k3d-myregistry.localhost:5001/docker.io/smartcontract/chainlink
   ```

7. (**Optional**) In case you want to use the changes you make locally in both repos (i.e. if you want to test the changes in the chainlink repo in accordance to the changes you made in the chainlink-testing-framework), you need to run the following in the chainlink/integration-tests folder:
   ```shell
   go mod edit -replace github.com/smartcontractkit/chainlink-testing-framework=~/go/src/github.com/chainlink-testing-framework
   ```

8. Finally, run make test_smoke inside the chainlink repo root folder to run all the integration tests. In case you want to only run some specific tests, you can do so with something like
    ```shell
    make test_smoke_simulated args="--focus-file=keeper_test.go"
    ``` 

this setup doc was modified from [notion](https://www.notion.so/chainlink/Setting-up-Integration-Tests-Framework-Locally-dc0e3db7718b45ad9249e97d7ef74c51)
