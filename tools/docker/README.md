# Using docker-compose for local development

The docker-compose configuration present in this directory allows for a user to quickly setup all of chainlink's services to perform actions like integration tests, acceptance tests, and development across multiple services.

# Requirements

- [docker-compose](https://docs.docker.com/compose/install/)

# Using the compose script

Inside the `chainlink/tools/docker` directory, there is a helper script that is included which should cover all cases of integration / acceptance / development needs acroos multiple services. To see a list of available commands, perform the following:

```sh
cd tools/docker
./compose help
```

## Examples

### Acceptance testing

Acceptance can be accomplished by using the `acceptance` command.

```sh
./compose acceptance
```

- The chainlink node can be reached at `http://localhost:6688`

Credentials for logging into the operator-ui can be found [here](../../tools/secrets/apicredentials)

###

### Doing local development on the core node

Doing quick, iterative changes on the core codebase can still be achieved within the compose setup with the `cld` or `cldo` commands.
The `cld` command will bring up the services that a chainlink node needs to connect to (parity/geth, postgres), and then attach the users terminal to a docker container containing the host's chainlink repository bind-mounted inside the container at `/usr/local/src/chainlink`. What this means is that any changes made within the host's repository will be synchronized to the container, and vice versa for changes made within the container at `/usr/local/src/chainlink`.

This enables a user to make quick changes on either the container or the host, run `cldev` within the attached container, check the new behaviour of the re-built node, and repeat this process until the desired results are achieved.

```sh
./compose cld
#
# Now you are inside the container
cldev # cldev without the "core" postfix simply calls the core node cli
#
# NAME:
#    main - CLI for Chainlink
#
# USAGE:
#    main [global options] command [command options] [arguments...]
#
# VERSION:
#    unset@unset
#
# COMMANDS:
#    admin              Commands for remotely taking admin related actions
#    bridges            Commands for Bridges communicating with External Adapters
#    config             Commands for the node's configuration
#    jobs               Commands for managing Jobs
#    node, local        Commands for admin actions that must be run locally
#    runs               Commands for managing Runs
#    txs                Commands for handling Ethereum transactions
#    agreements, agree  Commands for handling service agreements
#    attempts, txas     Commands for managing Ethereum Transaction Attempts
#    createextrakey     Create a key in the node's keystore alongside the existing key; to create an original key, just run the node
#    initiators         Commands for managing External Initiators
#    help, h            Shows a list of commands or help for one command
#
# GLOBAL OPTIONS:
#    --json, -j     json output as opposed to table
#    --help, -h     show help
#    --version, -v  print the version
cldev core # import our testing key and api credentials, then start the node
#
# ** Importing default key 0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f
# 2019-12-11T20:31:18Z [INFO]  Locking postgres for exclusive access with 500ms timeout orm/orm.go:74        #
# 2019-12-11T20:31:18Z [WARN]  pq: relation "migrations" does not exist           migrations/migrate.go:149
# ** Running node
# 2019-12-11T20:31:20Z [INFO]  Starting Chainlink Node 0.7.0 at commit 7324e9c476ed6b5c0a08d5a38779d4a6bfbb3880 cmd/local_client.go:27
# ...
# ...
```

### Cleaning up

To remove any containers, volumes, and networks related to our docker-compose setup, we can run the `clean` command:

```sh
./compose clean
```

### Running your own commands based off of docker-compose

The following commands allow you do just about anything:

```sh
./compose <subcommand>
./compose dev <subcommand>
```

For example, to see what our compose configuration looks like:

```sh
./compose config # base config
```

Or, to run just an ethereum node:

```sh
./compose up devnet # start a parity devnet node
```

```sh
GETH_MODE=true ./compose up devnet # start a geth devnet node
```

# Environment Variables

For more information regarding environment variables, the docker [documentation](https://docs.docker.com/compose/environment-variables/) explains it in great detail.
All of the environment variables listed under the `environment` key in each service contains a default entry under the `.env` file of this directory. Additional environment variables can be added by using the `chainlink-variables.env` file. Both files are further expanded upon below.

## Overriding existing variables

The existing variables listed under the `environment` key in each service can be overridden by setting a shell environment variable of the same key. For example, referring to `ETH_CHAIN_ID` variable under the `node` service, the default value of `ETH_CHAIN_ID` in `.env` is `34055`. If we wanted to change this to `1337`, we could set a shell variable to override this value.

```sh
export ETH_CHAIN_ID=1337
./compose acceptance # ETH_CHAIN_ID now has the value of 1337, instead of the default value of 34055
```

## Adding new environment variables

What if we want to add new environment variables that are not listed under the `environment` key of a service? `docker-compose` provides us with a way to pass our own variables that are not defined under the `environment` key by using an [env_file](https://docs.docker.com/compose/compose-file/#env_file). We can see from our `docker-compose.yaml` file that there is an env file under the name of `chainlink-variables.env`. In this file, you can specify any extra environment variables that you'd like to pass to the associated container.

For example, lets say we want to pass the variable `ALLOW_ORIGINS` defined in `store/orm/schema.go`, so that we can serve our api from a different port without getting CORS errors. We can't pass this in as a shell variable, as the variable is not defined under the `environment` key under the `node` service. What we can do though, is specify `ALLOW_ORIGINS` in `chainlink-variables.env`, which will get passed to the container.

```sh
# assuming that we're in the tools/docker directory

# Add our custom environment variable
echo "ALLOW_ORIGINS=http://localhost:1337" > chainlink-variables.env

# now the node will allow requests from the origin of http://localhost:1337 rather than the default value of http://localhost:3000,http://localhost:6688
./compose acceptance
```

# Following logs

The `logs` command will allow you to follow the logs of any running service. For example:

```bash
./compose up node # starts the node service and all it's dependencies, including devnet, the DB...
./compose logs devnet # shows the blockchain logs
# ^C to exit
./compose logs # shows the combined logs of all running services
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

## Logging into via the frontend results in HTTP Status Forbidden (403)

This is most likely due to the (Allow Origins access policy](https://docs.chain.link/docs/configuration-variables#section-allow-origins). Make sure you are using 'http://localhost' (not 127.0.0.1), or try disabling ALLOW_ORIGINS.

```
# Disable ALLOW_ORIGINS for testing
echo "ALLOW_ORIGINS=*" >> chainlink-variables.env
```

# Using the dockerized development environment

The dockerized development environment provides an alternative development and testing environment to the docker-compose setup as described above. The goals for this environment are to:

- create a development environment that is easily configured by interview candidates, potential contributors, etc.
- contain all dependencies in a single docker image
- contain sensible, pre-configured defaults

The entire chainlink repo is bind-mounted so any changes will take effect immediately - this makes the env good for TDD. Node modules are also bind-mounted, so you shouldn't have to install many deps after launching the container. Go deps are not bind-mounted, so you will have to install those after starting the container. You should only need to do this once, as long as you re-use the container.

The docker env contains direnv, so whatever changes you make locally to your (bind-mounted) `.envrc` will be reflected in the docker container. The container is built with a default ENV that should require minimal changes for basic testing and development.

### Building the dev environment

```bash
# build the image and tag it as chainlink-develop
docker build ./tools/docker/ -t chainlink-develop:latest -f ./tools/docker/develop.Dockerfile
# create the image
docker container create -v /home/ryan/chainlink/chainlink:/root/chainlink --name chainlink-dev chainlink-develop:latest
# if you want to access the db, chain, node, or from the host... expose the relevant ports
# This could also be used in case you want to run some services in the container, and others directly
# on the host
docker container create -v /home/ryan/chainlink/chainlink:/root/chainlink --name chainlink-dev -p 5432:5432 -p 6688:6688 -p 6689:6689 -p 3000:3000 -p 3001:3001 -p 8545:8545 -p 8546:8546 chainlink-develop:latest
# start the container (this will run in the background until you stop it)
docker start chainlink-dev
```

### Connecting to the dev environment

```bash
# connect to the container by opening bash prompts - you can open as many as you'd like
docker exec -ti chainlink-dev bash
```

### Run services / tests inside container

\$ --> inside container bash prompt

This is nothing new, just a demonstration that you should be able to run all the commands/tests/services you normally do for development/testing, but now inside of the docker container. As mentioned above, if you want access to these services on the host machine, you will have to expose their ports.

```bash
# install deps and chainlink
$ make install

# run go tests
$ make testdb
$ go test ./...

# run evm tests
$ cd contracts
$ pnpm test

# start geth
$ geth --dev --datadir ./tools/gethnet/datadir --mine --ipcdisable --dev.period 2 --unlock 0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f --password ./tools/clroot/password.txt --config ./tools/gethnet/config.toml

# run chainlink node (will require changing env vars from defaults)
$ chainlink local node -a ./tools/secrets/apicredentials -p ./tools/secrets/password.txt
```

### Included Tooling:

This image contains the following additional tools:

- geth, openethereum, ganache
- delve, gofuzz
- slither, echidna
- web3.py
