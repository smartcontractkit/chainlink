# Chainlink Go Implementation

## Developer setup

After creating the standard Golang directory structure and configuring your $GOPATH, clone this repo and perform the following:

```bash
$ brew install dep
$ dep ensure # Installs dependencies into ./vendor
$ go build -o chainlink
$ ./chainlink
```

### Solidity Development setup

```bash
$ cd solidity
$ yarn install
$ truffle test
```

### direnv

We use [direnv](https://github.com/direnv/direnv/) to set up PATH and aliases 
for a friendlier developer experience. Here is an example `.envrc` that we use:

```bash
PATH_add tmp
PATH_add solidity/node_modules/.bin
```
