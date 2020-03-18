# Running cypress tests interactively

Our cypress tests require a lot of backing services to be run, so usage of docker-compose is expected.

# Setup

1. Start the acceptance environment

```sh
# from the git repo root
cd tools/docker
./compose acceptance
```

2. In another window pane, build and start the cypress job server.

```sh
# from the git repo root
cd tools/cypress-job-server
docker build -t cjs -f Dockerfile ../..
# The flags passed make it so our services from the acceptance environment
# are able to reach this server via the cypress-job-server name.
docker run --name cypress-job-server --network docker_default cjs
```

If you get an error like:

```
docker: Error response from daemon: Conflict. The container name "/cypress-job-server" is already in use by container "8efd0b5d704eaf400a6b51e87ec2aaeff16e867292258f401a3e5d19fe1add10". You have to remove (or rename) that container to be able to reuse that name.
```

Then you can run

```sh
docker rm <containerHash>
```

and then retry step 2.

# Running

Now we can interactively run the cypress test runner.

```sh
# from the git repo root
cd integration
yarn cypress open
```
