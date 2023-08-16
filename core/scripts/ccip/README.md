# CCIP Setting up a DON
To configure a DON and all contracts from nothing we need a few things

- [ ]  1 private key
- [ ]  k8s cluster (or other way of running a set of nodes)
- [ ]  2 blockchains
- [ ]  1pass access
- [ ]  vpn access
- [ ]  (optional) access to CCIP Shared Vault in 1password
- [ ]  (optional) Grafana access

The scripts are in beta state and will change in the future!

To see k8s config please check example for beta [k8s config](https://github.com/smartcontractkit/infra-k8s/tree/main/projects/chainlink/files/chainlink-clusters/clc-ocr-multichain-ccip-beta/config.yaml) in the `infra-k8s` repo. NOTE: replace `beta` with `prod` when using the prod cluster.

## Funding keys

Chainlink has an **[internal faucet](https://internal-apps-ui.main.prod.cldev.sh/faucets)** to fund both the native gas token and Link tokens. Most relevant testnets are supported. Fund the wallet associated with the private key on both networks. Make sure to reserve ~0.5 eth (or chain specific equivalent) per CL node since you’ll also be funding the nodes from this wallet. You will also need a few link tokens to fill the pools, you can get those from the faucet as well.

## Node credentials

To communicate with the nodes in the k8s cluster you will need to have their urls & passwords. These can be found in 1pass by searching the node name e.g `ccip beta` or if you have access to `CCIP Shared Vault` there is ready to go json for every environment. Infra is the owner of the k8s instances and will populate 1pass with these credentials when they spin up new nodes. To test access simply go to the url you see in 1pass and test it. This should be first and second to last time you open the node UI (we’ll need to open the bootstrap node UI once more)

## Setting up the script

The script needs the node credentials discussed in the last step. To supply them we transform them into a JSON file and put that file in `core/scripts/ccip/json/credentials/` in the following format. You can get current node configuration in 1pass for ex. `CCIP Beta testnet credentials`

```json
{
  "Env": "production",
  "Bootstrap": {
    "URL": "https://",
    "Email": "admin@chain.link",
    "Password": "pass",
    "RemoteIP": "https://"
  },
  "Nodes": [
    {
      "URL": "https://",
      "Email": "admin@chain.link",
      "Password": "pass",
      "RemoteIP": "https://"
    },
    ...
  ]
}
```

The name of this file should correspond to an `Environment` in `core/scripts/ccip/dione/dione.go`. You can freely create new environments when needed. An example of a filename would be `staging-beta.json`.

Next we’re going to set the environment variables needed to run the various scripts in `core/scripts/ccip/ccip_runner_test.go`.

- `CL_DATABASE_URL` this is because all tests require a database url. Atm. any random value will do here for ex. `_test` The script does not use a db.
- `OWNER_KEY` this is the private key of the owner account. This account will own all contracts. This should be the key that was funded in the first step.

`TestCCIP` requires two more variables set

- `SEED_KEY` a second private key used to derive 10 sender keys from. Only used in batching tests.
- `COMMAND` the command to run. This test is an entry point to various commands that are explained below the function itself.

At the top of [ccip_runner_test.go](ccip_runner_test.go) you’ll find three variables. These determine the two blockchains and the environment (the json you just configured) for all scripts to use.

## Deployments with Rhea

We’ll be deploying two new lanes, one source → dest and one dest → source. Rhea is bidirectional so you just run it once to set both direction lanes.

We keep separately blockchain config in`EVMChainConfig` object and lane specific one in `EvmDeploymentConfig`

Create a new config and set the `ChainId`, `EthUrl` (through the secrets file) in `chainIdToRPC` map. In addition you will need to add rpc wss. Example:

[secrets.go](/secrets/secrets.go)

```go
var chainIdToRPC = map[uint64]string{
	5:        "wss://...",
	43113:    "wss://...",
	420:      "wss://...",
	11155111: "wss://...",
}
```

## Deploying to the new chain

We need to go to `core/scripts/ccip/rhea/deployments/` and create two new `EvmChainConfig` objects in one of environments files (beta/prod). When naming them please use the Environment name as well as the blockchain. This chain config will contain information about both the blockchain and the deployed contracts.

Set `ChainId`, `GasSettings`, `SupportedTokens` and `DeploySettings` `FeeTokens` , `WrappedNative` - reuse them if there is configuration for specific chain already. Example for Sepolia on Beta:

```go
var Beta_Sepolia = rhea.EVMChainConfig{
	ChainId: 11155111,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: false,
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x779877A7B0D9E8603169DdbD7836e478b4624789"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.LINK.Price(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x097D90c9d3E0B50Ca60e1ae45F6A81010f9FB534"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         TokenPrices[rhea.WETH],
			TokenPoolType: rhea.LockRelease,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress(""),
	Afn:           gethcommon.HexToAddress(""),
	PriceRegistry: gethcommon.HexToAddress(""),
	Confirmations: 4,
	DeploySettings: rhea.ChainDeploySettings{
		DeployAFN:           true,
		DeployTokenPools:    true,
		DeployRouter:        true,
		DeployPriceRegistry: true,
		DeployedAtBlock:     0,
	},
}
```

Then we need to add another two `EvmDeploymentConfig` objects for the lane setup.

In the DeploySettings we’ll set everything to true except for the `DeployPingPongDapp`. We can deploy dapps at a later point, at this time we focus on the CCIP contracts. This means we’ll deploy new versions of all the core CCIP contracts when we deploy. When doing partial deployments (e.g. there is already an AFN we want to re-use) we could simply set the AFN property to false and it would re-use the address that we set in the `EvmChainConfig`. The `DeployedAt` property can be set to 0 for now. We’ll set this to the block number when we started the deployment at a later point.

It should look like this:

```go
// Lanes
var Beta_SepoliaToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		CommitStore:  gethcommon.HexToAddress(""),
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployCommitStore:  true,
			DeployRamp:         true,
			DeployPingPongDapp: true,
			DeployedAtBlock:    0,
		},
	},
}
```

Add new chain to `Chains` and lanes to specific `ChainMapping`. Example for Avax Fuji <-> Sepolia:
```go
var BetaChains = map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.AvaxFuji:       {ChainConfig: Beta_AvaxFuji},
	rhea.Sepolia:        {ChainConfig: Beta_Sepolia},
}

var BetaChainMapping = map[rhea.Chain]map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.Sepolia: {
		rhea.AvaxFuji:       Beta_SepoliaToAvaxFuji,
	},
	rhea.AvaxFuji: {
		rhea.Sepolia:        Beta_AvaxFujiToSepolia,
	},
}
```

We now set the newly created chain configs as `sourceChain` and `destChain` in [ccip_runner_test.go](ccip_runner_test.go) . Make sure to set the ENV as well.

### Deploying

Now we’re ready to deploy the contracts! The deployment is split into `Chains` part and `Lane` part.

Run `TestRheaDeployChains` for `Chains` part to deploy needed chains.
Run output should be written to console & `./json/deployments/env/chain/....`
We need to populate the `EvmChainConfig` with the new chain data before running `Lane` part.
Run `TestRheaDeployLane` for `Lane` part to deploy a single lane. 
The console should keep you up to date on what it’s doing.
The scripts will deploy and configure everything for you. After it’s done it will print out the contract locations.
Now populate `rhea.EvmDeploymentConfig` for specific lane you just set up.

We can now also put in the `DeployedAt` value. This is only used in the jobspecs to make sure we replay all txs until contract creation. Any block number around the deployment will work.

Whole deployment data is saved into the file for specific environment `/json/deployments/env` with all contracts deployed per lane and chain.

## Blue-green deployment

Blue-green (aka zero downtime deployment) is split into three phases:
* **Deploy Upgrade Lane** - deploys CommitStore, OnRamp and OffRamp contracts and attach them to the UpgradeRouter, but without any changes to the current deployment. After this step, "Upgrade Lane" should be fully operational and work in parallel to the current deployment.
* **Promote Upgrade Lane** - registers contracts from Upgrade Lane in [Prod/Non-Upgrade] Routers on both ends. This is the point in which we switch routing from current deployment to the upgrade/new deployment.
* **(optional) Clean up** - removes old jobspecs (after pending requests are processed)

### How to deploy

#### Deploy Upgrade Lane

Start by making sure that you have proper source/destination chains, environment and version in place
```go
var (
	sourceChain        = rhea.ArbitrumGoerli
	destChain          = rhea.Sepolia
	ENV                = dione.StagingBeta
	currentVersion     = "0.3.0+core2.1.0"        // Current Version of jobspecs. You can find it by logging to the NOP and listing jobs. If there is no suffix, then currentVersion is empty
	upgradeLaneVersion = "0.4.0-beta.0+core2.1.0" // Version that will be used to differentiate current deployment from upgrade deployment.  
)
```
At this point `currentVersion` should point to the currently deployed lane's version. `upgradeLaneVersion` should point to next version what will be created during deployment.
Please follow [this guide](https://github.com/smartcontractkit/releng/blob/9f7154b48d0a1f74e6f1ddf65fd0b76bc4c0a5aa/docs/vault/notes/guides.release-processes.ccip.ccip-core-node.md) to generate a new version. Once created, just copy it to the `upgradeLaneVersion` variable.

**Notice** At this moment, ARM doesn't support blue-green deployments / multiple contracts instances per single lane. Therefore, we need to deploy MockARM and use it in `UpgradeLane`.
You can do it by setting following options on every chain we are deploying.
```go
var Beta_Sepolia = rhea.EVMChainConfig{
        // ...
	ARMConfig: nil,
	DeploySettings: rhea.ChainDeploySettings{
		DeployARM:           true,
                // ...
	},
}
```
Now you can deploy mocked ARM contracts by calling `TestRheaDeployChains`. If current configuration points to the real ARM, back it up - you are going to need it later, because at the end of the day, promoted lane has to use real ARM.
Please update `$env.go` file with updated arm addresses like this:
```go
var Beta_Sepolia = rhea.EVMChainConfig{
	ARM:           gethcommon.HexToAddress("0xAddressFromJsonFile"),
}
```

We need to let scripts know that we want to deploy Upgrade Lane by adding `EVMLaneConfig` under `UpgradeLaneConfig` field for each end. 
For instance, Sepolia <> Arbitrum deployment requires adding `UpgradeLaneConfig` to both `$ENV_SepoliaToArbitrumGoerli` and `$ENV_ArbitrumGoerliToSepolia`
```go
var Beta_SepoliaToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		// No changes ...
	},
	UpgradeLaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x0"),
		OffRamp:      gethcommon.HexToAddress("0x0"),
		CommitStore:  gethcommon.HexToAddress("0x0"),
		PingPongDapp: gethcommon.HexToAddress("0x0"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployCommitStore:  true,
			DeployRamp:         true,
			DeployPingPongDapp: true,
			DeployedAtBlock:    0,
		},
	},
}
```

We are ready to deploy Upgrade Lane by calling `TestRheaDeployUpgradeLane` from [ccip_runner_test.go](ccip_runner_test.go). This might take a couple of minutes to finish.
Run output should be written to console & `./json/deployments/env/upgrade-lane/....`. We need to populate the `UpgradeLaneConfig` with the contracts created during deployment. 
Make sure you set deployment options to false to avoid accidental redeployment.
```go
var Beta_SepoliaToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		// No changes ...
	},
	UpgradeLaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xABC"),
		OffRamp:      gethcommon.HexToAddress("0xDEF"),
		CommitStore:  gethcommon.HexToAddress("0xGHI"),
		PingPongDapp: gethcommon.HexToAddress("0xJKL"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployCommitStore:  false,
			DeployRamp:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    123456,
		},
	},
}
```

So far so good. At this point we should have a fully operational Upgrade Lane working on a real environment! Let's verify if everything is up and running. 
I recommend doing the following steps:
* Go to the NOP and verify if there are jobspecs for both current deployment and upgrade deployment. Jobs' names contain version.
* Check PingPongDaaps in block explorer and verify if traffic flows between them - you might need around ~20 minutes to observe it on both ends
* Verify OffRamps in block explorer, you should be able to see DON executing there
* Logs and dashboard similar to the classic deployment

Great! You are ready to move on and promote your deployment!

### Promote Upgrade Lane

Second phase of deployment includes attaching new contracts (OnRamp and OffRamp) to the Router and switching real traffic there. 
We will use `TestRheaPromoteUpgradeLaneDeployment` from [ccip_runner_test.go](ccip_runner_test.go).
Run output should be written to console &  `./json/deployments/env/lane/....`. We need to replace `LaneConfig` with `UpgradeLaneConfig`, it's better to copy those values from `json` files generated as they are source of the truth in this case.
Also, please remove `UpgradeLaneConfig` to avoid confusion, stating that there is nothing waiting to promoted on the Upgrade Lane.
```go
var Beta_SepoliaToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
            OnRamp:       gethcommon.HexToAddress("0xABC"),
            OffRamp:      gethcommon.HexToAddress("0xDEF"),
            CommitStore:  gethcommon.HexToAddress("0xGHI"),
            PingPongDapp: gethcommon.HexToAddress("0xPreviousAddressProbably"),
            DeploySettings: rhea.LaneDeploySettings{
                DeployCommitStore:  false,
                DeployRamp:         false,
                DeployPingPongDapp: false,
                DeployedAtBlock:    123456,
            },       
	}
}
```

Verify if everything is okay. If so, then you can apply the same steps to all lanes.
Once all lanes are upgraded and using the same versions of contracts you can upgrade versions and commit them.
```go
var (
    currentVersion     = upgradeLaneVersion
    upgradeLaneVersion = "" // Will be filled during next deployment
)
```


### Switch back to the Real ARM

This has to be coordinated with the research team. Start with bringing back values of the real ARM contracts to the chain config

```go
var Beta_Sepolia = rhea.EVMChainConfig{
	ARM:           gethcommon.HexToAddress("0xPointingToTheRealARM"),
	ARMConfig:     "previousValue"
}
```

Run `TestUpdateLaneARMAddress` and let know research team about newly deployed contracts' addresses - they have to update those in the offchain part of the ARM. Once it's done you might need ~30 minutes for ARM start blessing roots. 
In case of ARM cursing lanes, you might need to use `TestCCIP.uncurseSourceARM` to remove curse on a particular lane.

## Checking the deployment

To verify the deployment we have the Metis cli tool we ran before. Run it again but now on the freshly deployed contracts. It should output a long list of tables that should **only** contains `true`. Two that can be false at this point are `CommitStore OCR2 configured`  and `OffRamp OCR2 configured` - we set them in the next step `Setting the Don` below. Any `false` other than that indicates something went wrong in the configuration. A full clean deployment is fairly well tested so I do not suspect anything to be wrong at this point. If anything is wrong it should be fairly descriptive about what is wrong and the fixes are almost always easy to manually do.

## Setting the DON

At this point we only have to make the nodes aware of the new contracts by submitting job specs and setting the OCR2 config on the CommitStore and OffRamp.

The nodes infra setup should have the proper `EVM_CHAINS` set and therefore have sending keys for the relevant chains. If this is not the case you need to get this sorted first! This text assumes the nodes & keys are properly set and there are no orphaned keys (keys for chains without running nodes).

### Dione

Dione is the tool that will handle all the node communication for us. The script allows for many things, you should configure what it should do yourself with the don functions. First we should get some information from the nodes, namely the sending keys, the OCR2 keys and the PeerIDs. 
This can be done with specified `TestDionePopulateNodeKeys` script.

The `WriteToFile` call will write the params to disk in the  `/json/nodes/` folder its also needed to be called when cluster is changed.

```json
{
  "Bootstrap": {
    "EthKeys": null,
    "PeerID": "peerID",
    "OCRKeys": {
      "data": null
    }
  },
  "Nodes": [
    {
      "EthKeys": {
        "5": "0x"
      },
      "PeerID": "peerID",
      "OCRKeys": {
        "data": [
      ...
        ]
      }
    }
    ]
}
```

You might see references to `starknet`, `terra` or `solana`, this is normal.

### Funding nodes.

To run CCIP we need funds on cll nodes eoa’s. To fund the chain addresses use `TestFundNode`

You need to do it for both ways. Just swap `sourceChain` and `destChain`

This will fund the EthKeys we just retrieved from the nodes. You can change the amount, this will give each node 4 gas fee tokens. The script does not wait for any tx confirmations and does no handle any errors. That does make it **very** fast. Don’t worry if it looks like it did nothing, just visit your address on the block explorers of the configured chains and check if everything went through. Give it a minute, depending on the chain.

`TestPrintNodeBalances` can be run to verify node balances.

### Bootstrap job

At this point we hit something that is not automated (yet..). You will have to manually submit the bootstrap job spec. Luckily, we can generate the jobspec so we simply copy-paste the spec in the node UI.

To print all the jobspecs to the console run `TestCCIP` with the env `COMMAND="printSpecs"`

It will print all the job specs but don’t worry about the relay/execution specs, we’ve fully automated those! Open the Bootstrap UI in a browser and copy-paste in the boostrap spec. When you look at the spec you’ll probably see some error about the contract not being configured. This is fine, we’ll do that soon.

### Job specs

We can come back to the Dione and set up nodes jobs.

Just run `TestDione`Code below:

```go
// TestDione can be run as a test with the following config
// OWNER_KEY  private key used to deploy all contracts and is used as default in all single user tests.
func TestDione(t *testing.T) {
	checkOwnerKeyAndSetupChain(t)

	don := dione.NewDON(ENV, logger.TestLogger(t))
	don.ClearAllJobs(helpers.ChainName(int64(SOURCE.ChainConfig.ChainId)), helpers.ChainName(int64(DESTINATION.ChainConfig.ChainId)))
	don.AddTwoWaySpecs(SOURCE, DESTINATION)

	// Sometimes jobs don't get added correctly. This script looks for missing jobs
	// and attempts to add them.
	don.AddMissingSpecs(DESTINATION, SOURCE)
	don.AddMissingSpecs(SOURCE, DESTINATION)
}
```

This will add **all** the needed jobs, so 4 per node: 2 relay and 2 execution for both ways. It will also clear any old job specs so you don’t have to delete any jobs manually from current lane setup. Sometimes nodes may not able to handle all of requests to create jobs or you may have some connection issues - then `AddMissingSpecs` will check each of nodes to add any missing jobs.

### Setting OCR2 Config

Now we’ve deployed the job specs we have to configure the contracts OCR2 config. To do this run the `TestCCIP`script with the env `COMMAND="setConfig"`. This needs to be done both ways, so swap `sourceChain` and `destChain` and re-run the script.

Congratulations! You’ve now deployed CCIP 🎉🥳

## **Validating with Metis**

For verifying state of the lane we use Metis cli tool. To check the configuration state of the `SOURCE` and `DESTINATION` blockchains set them in path [metis/main.go](metis/main.go) you should now be able to run it with `s` option. We will be creating new SOURCE and DESTINATION configurations in the next step. Remember about environment variables we had set above.

## Validating deployment with Ping Pong

You might want to validate and celebrate the deployment with a never-ending ping pong session. To do so run `TestCCIP` with the env `COMMAND="deployPingPong"` and manually set the contract addresses in the `EvmChainConfig`. Then run `TestCCIP` with `startPingPong`. If everything went well you now have an endless ping pong (or at least until the funding runs out).

## Grafana

You might want to check some logs, this can be done in Grafana. You can take advantage of CCIP dashboard for existing lanes.

[Grafana CCIP Dashboard](https://chainlinklabs.grafana.net/d/R5zMiJv4z/ccip-dashboard)

You can also manually catch up things in Explorer for new blockchains.

```shell
{cluster="staging-us-west-2-main", pod=~"clc-ocr-multichain-ccip-beta-nodes-.*"} |~ "❗"
```

Looking for ❗will give you all the attempts that nodes do to put a tx on chain.