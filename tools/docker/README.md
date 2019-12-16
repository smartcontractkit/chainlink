# Using docker-compose for local development
The docker-compose configuration present in this directory allows for a user to quickly setup all of chainlink's services to perform actions like integration tests, acceptance tests, and development across multiple services.

# Requirements
- [docker-compose](https://docs.docker.com/compose/install/)

# Using the compose script
Inside the `chainlink/tools/docker` directory, there is a helper script that is included which should cover all cases of integration / acceptance / development needs acroos multiple services. To see a list of available commands, perform the following:

```sh
$ cd tools/docker
$ ./compose help
```

## Examples
### Acceptance testing
Acceptance can be accomplished by using the `acceptance` command.
```sh
./compose acceptance
```
- The explorer can be reached at `http://localhost:3001`
- The chainlink node can be reached at `http://localhost:6688`

Credentials for logging into the operator-ui can be found [here](../secrets/apicredentals)

### 
### Doing local development on the core node
Doing quick, iterative changes on the core codebase can still be achieved within the compose setup with the `cld` or `cldo` commands.
The `cld` command will bring up the services that a chainlink node needs to connect to (explorer, parity/geth, postgres), and then attach the users terminal to a docker container containing the host's chainlink repostiory bind-mounted inside the container at `/usr/local/src/chainlink`. What this means is that any changes made within the host's repository will be synchronized to the container, and vice versa for changes made within the container at `/usr/local/src/chainlink`.

This enables a user to make quick changes on either the container or the host, run `cldev` within the attached container, check the new behaviour of the re-built node, and repeat this process until the desired results are achieved.

```sh
$ ./compose cld
#
# $$ denotes that we're now in the chainlink container
$$ cldev # cldev without the "core" postfix simply calls the core node cli
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
$$ cldev core # import our testing key and api credentials, then start the node
# 
# ** Importing default key 0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f
# 2019-12-11T20:31:18Z [INFO]  Locking postgres for exclusive access with 500ms timeout orm/orm.go:74        #   
# 2019-12-11T20:31:18Z [WARN]  pq: relation "migrations" does not exist           migrations/migrate.go:149 
# ** Running node
# 2019-12-11T20:31:20Z [INFO]  Starting Chainlink Node 0.7.0 at commit 7324e9c476ed6b5c0a08d5a38779d4a6bfbb3880 cmd/local_client.go:27  
# 2019-12-11T20:31:20Z [INFO]  SGX enclave *NOT* loaded                           cmd/enclave.go:11       
# 2019-12-11T20:31:20Z [INFO]  This version of chainlink was not built with support for SGX tasks cmd/enclave.go:12       
# ...
# ...
```

`cldo` allows the user to perform the same actions above, but also applied to the operator-ui codebase and the core codebase. The operator-ui will be hosted in hot-reload/development mode at `http://localhost//3000`. To see the build progress of operator-ui, we can open another terminal to watch its output while we can mess around with the core node in the original terminal.

In the first terminal:
```sh
$ ./compose cldo
#
# $$ denotes that we're now in the chainlink container
$$ cldev # cldev without the "core" postfix simply calls the core node cli
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
$$ cldev core # import our testing key and api credentials, then start the node
# 
# ** Importing default key 0x9ca9d2d5e04012c9ed24c0e513c9bfaa4a2dd77f
# 2019-12-11T20:31:18Z [INFO]  Locking postgres for exclusive access with 500ms timeout orm/orm.go:74        #   
# 2019-12-11T20:31:18Z [WARN]  pq: relation "migrations" does not exist           migrations/migrate.go:149 
# ** Running node
# 2019-12-11T20:31:20Z [INFO]  Starting Chainlink Node 0.7.0 at commit 7324e9c476ed6b5c0a08d5a38779d4a6bfbb3880 cmd/local_client.go:27  
# 2019-12-11T20:31:20Z [INFO]  SGX enclave *NOT* loaded                           cmd/enclave.go:11       
# 2019-12-11T20:31:20Z [INFO]  This version of chainlink was not built with support for SGX tasks cmd/enclave.go:12       
# ...
# ...
```

In a new terminal:
```sh
docker logs operator-ui -f
```
You'll now have two terminals, one with the core node, one with operator-ui, with both being able to react to code changes without rebuilding their respective images.

### Running integration test suites
The integration test suite will run against a parity node by default. You can run the integration test suite with the following command:
```sh
$ ./compose test
```

If you want to run the test suite against a geth node, you can set the `GETH_MODE` environment variable.
```sh
$ GETH_MODE=true ./compose test
```

If we want to quickly test new changes we make to `integration/` or `tools/ci/ethereum_test` without re-building our images, we can use the `test:dev` command which will reflect those changes from the host file system without rebuilding those containers.
```sh
$ ./compose test:dev
```

### Cleaning up
To remove any containers, volumes, and networks related to our docker-compose setup, we can run the `clean` command:
```sh
./compose clean
```

### Running your own commands based off of docker-compose
The following commands allow you do just about anything:
```sh
$ ./compose <subcommand>
$ ./compose integ <subcommand>
$ ./compose dev:integ <subcommand>
$ ./compose dev <subcommand>
```
For example, to see what our compose configuration looks like:
```sh
$ ./compose config # base config
$ ./compose dev:integ # development integration test config
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
