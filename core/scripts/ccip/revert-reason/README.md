## Setup

Before starting:

1. Create `.env` file based on the example `.env.example` next to the CLI binary.

2. If you plan to resolve revert reasons from a transaction hash, you will need:
   1. An EVM chain endpoint URL
   2. A from address

The endpoint URL can be a locally running node (archive mode), or an externally hosted one like
[alchemy](https://www.alchemy.com/).



To see all available commands, run the following:
```bash
go run main.go --help
```


## Usage

Decoding an error code string (offline):

```bash
> ./ccip-revert-reason reason --from-error "0x4e487b710000000000000000000000000000000000000000000000000000000000000032"
2022/12/05 15:18:33 Using config file .env
Decoded error: Assertion failure
If you access an array, bytesN or an array slice at an out-of-bounds or negative index (i.e. x[i] where i >= x.length or i < 0).%                    
```


Resolving from a transaction hash (`NODE_URL` and `FROM_ADDRESS` env vars need to be defined)

```bash
> ./ccip-revert-reason reason "0x4e487b710000000000000"
2022/12/05 15:18:33 Using config file .env
Decoded error: Assertion failure
If you access an array, bytesN or an array slice at an out-of-bounds or negative index (i.e. x[i] where i >= x.length or i < 0).%                    
```