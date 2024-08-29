# Using docker-compose for local development

The docker-compose configuration present in this directory allows to quickly run a local development environment.

# Requirements

- [docker-compose](https://docs.docker.com/compose/install/)

# Using the compose script

You should use the script `compose` located in the `chainlink/tools/docker` directory.
To see a list of available commands, perform the following:

```sh
cd tools/docker
./compose help
```

## Env vars

### .env file
.env file is used to set the environment variables for the docker-compose commands

### Compose script env vars
The following env vars are used for the compose script :
- `WITH_OBSERVABILITY=true` to enable grafana, prometheus and alertmanager
- `GETH_MODE=true` to use geth instead of parity
- `CHAIN_ID=<number>` to specify the chainID (default is 34055 for parity and 1337 for geth)
- `HTTPURL=<url>` to specify the RPC node HTTP url (default is set if you use geth or parity)
- `WSURL=<url>` to specify the RPC node WS url (default is set if you use geth or parity)

If you specify both `HTTPURL` and `WSURL`, it won't run the devnet RPC node.

for example :
```sh
CHAIN_ID=11155111 WSURL=wss://eth.sepolia HTTPURL=https://eth.sepolia ./compose dev
```

```sh
WITH_OBSERVABILITY=true ./compose up
```

## Dev

Will run one node with a postgres database and by default a devnet RPC node that can be either geth or parity.

```sh
./compose dev
```

The chainlink node can be reached at `http://localhost:6688`

Credentials for logging into the operator-ui can be found [here](../../tools/secrets/apicredentials)

## Up

Runs all services including two nodes with two postgres databases and by default a devnet RPC node that can be either geth or parity.

```sh
./compose up
```

## Cleaning up

To remove any containers, volumes, and networks related to our docker-compose setup, we can run the `clean` command:

```sh
./compose clean
```

## Logs

You can use logs command to see the logs of a specific service or all services.

```sh
./compose logs node # shows the logs of the node service
```

```sh
./compose logs # shows the combined logs of all running services
```

## Connecting to the dev environment

```sh
# connect to the container by opening bash prompts
./compose connect
```

# Troubleshooting

## My storage space is full! How do I clean up docker's disk usage?

```
docker system prune
```

## The build process takes up a lot of resources / brings my computer to a halt

The default configuration tries to build everything in parallel. You can avoid this by clearing the Docker Compose build options.

```
# Use docker compose's default build configuration
export DOCKER_COMPOSE_BUILD_OPTS=""
```

## Logging from a container is hidden

Sometimes docker-compose does not show logging from some docker containers. This can be solved by using the docker command directly.

```
# List docker instances
docker ps
# CONTAINER ID        IMAGE                     COMMAND                  CREATED             STATUS              PORTS                                                                                                 NAMES
# 41410c9d79d8        smartcontract/chainlink   "chainlink node star…"   2 minutes ago       Up 2 minutes        0.0.0.0:6688->6688/tcp                                                                                chainlink-node
# f7e657e101d8        smartcontract/devnet      "/bin/parity --confi…"   47 hours ago        Up 2 minutes        5001/tcp, 8080/tcp, 8082-8083/tcp, 8180/tcp, 8546/tcp, 30303/tcp, 0.0.0.0:8545->8545/tcp, 30303/udp   parity

# Follow logs using name of container
docker logs -f chainlink-node
```
