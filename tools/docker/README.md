# Using docker-compose for local development
The docker-compose configuration present in this directory allows for a user to quickly setup all of chainlink's services to perform actions like integration tests, acceptance tests, and development across multiple services.

# Requirements
- [docker-compose](https://docs.docker.com/compose/install/)

# Services
The following services are exposed by the docker-compose configuration, and are referred to by `<service name>` in consequent descriptions.

## node
The `node` service is the chainlink node, configured to run in development mode by default. When running, you can access the node's front end by navigating to `localhost:6688` and using the credentials defined in `../secrets/apicredentials`.  The first line is the username, the second is the password.

## devnet
The `devnet` service spins up an ethereum blockchain in development mode, it's used by the `node` service to do blockchain interactions.

## explorer
The `explorer` service contains the Chainlink Explorer, a service that allows users to search and discover chainlink specific transactions, jobs, and more. The served explorer site is accessible from `localhost:3000`. 

## Setup
The following commands assume that you're executing `docker-compose` commands with the current working directory being `tools/docker`.
A full description of how to run `docker-compose` can be found in its [web documentation](https://docs.docker.com/compose/).

## Build
Before being able to run our docker containers, we'll need to build their corresponding images. Make sure to re-build your images to reflect repository changes.

### Building all images
The following command will build all internal images, one by one. Any external images will instead be fetched.
```sh
docker-compose build
```

You can build images in parallel to speed up build times with the `--parallel` flag.
```sh
docker-compose build --parallel
```

### Build a single image
The following command will build a single image, along with all of its dependent images.
```sh
docker-compose build <service name>
```

## Startup
### Start all services
The following command will start up all services and their dependent child services.
```sh 
docker-compose up
```

Adding the `-d` flag will detach the spun up services from your terminal.
```sh 
docker-compose up -d
```

### Start a single service
The following command will start up a single service along with any service dependencies it has.
```sh
docker-compose up <service name>
```

Adding the `-d` flag will detach the spun up services from your terminal.
```sh
docker-compose up -d
```

## Teardown
### Shutdown running services
```sh
docker-compose down
```

### Shutdown running services along with their volumes
This will remove all volumes of the spun-up services. Useful if you want to completely wipe state before running them again (database state, blockchain state, etc).
```sh
docker-compose down -v
```