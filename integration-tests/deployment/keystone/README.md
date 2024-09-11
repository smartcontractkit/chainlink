# Programmatic deployment of Keystone

The current scope of this package is the ability to deploy and configure the Capability Registry, OCR3 and Forwarder contracts.

It builds on the `Environment` abstraction introduced by `chainlink-deployments`. The original concept of Environment was to delineate dev vs testnet vs prod.
A deployment would inject the necessary configuration (eg simulated chain vs testnet chain) as appropriate via the `Environment` abstraction, and the
deployment implementation be agnostic to these details.

Here we use an extended concept `MultiDonEnvironent`, which is simply and `Environment` in which different dons (nodesets) host specific capabilities. 

The entry point to the deployment is the `Deploy` func. The arguments to the this func are environment dependent, and vary from one deployment to another.

```
type DeployRequest struct {
	RegistryChainSel uint64
	Menv             deployment.MultiDonEnvironment

	DonToCapabilities map[string][]kcr.CapabilitiesRegistryCapability                   // from external source; the key is a human-readable name. TODO consider using the 'sortedHash' of p2pkeys as the key rather than a name
	NodeIDToNop       map[string]capabilities_registry.CapabilitiesRegistryNodeOperator // TODO maybe should be derivable from OffchainClient interface but doesn't seem to be notion of NOP in job distributor
	OCR3Config        *OracleConfigSource                                               // TODO: probably should be a map of don to config; but currently we only have one wf don therefore one config
}

type DeployResponse struct {
	Changeset *deployment.ChangesetOutput
	DonInfos  map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

func Deploy(ctx context.Context, lggr logger.Logger, req DeployRequest) (*DeployResponse, error)
```


In order to make this all work we need a mapping what nodes run which capabilities, which nodes belong to what nop, and ocr configuration for consensus. These are represented by `DonToCapabilities, NodeIDToNop, OCR3Config`, respectively.

The first and last are simple artifacts, there is no apriori constraints on there values.

However, `NodeIDToNop` is glue between the external system that tracks real-world nodes and nops. For us, this system is CLO. This means that values in `NodeIDToNop` need to be derived coherently from data that is sourced from prod CLO.

## Obtaining and parsing CLO data

A real deployment requires real data from CLO.

### Requirements
- clo access to [stage](https://feeds-manager.main.stage.cldev.sh/sign-in), [prod](https://feeds-manager.main.prod.cldev.sh/sign-in) (ask in #topic-keystone-clo and tag Joey Punzel)
- [clo cli](https://github.com/smartcontractkit/feeds-manager#chainlink-orchestrator-api-client)
- stage & prod configuration for the cli


As hinted above, CLO is the system that knows about nodes and node operators. One of our goals is to configure the registry contract with the nodes and nops. So we have to faithfully plumb the values in CLO to our deployment.

For the time being, it was not possible to do this programmatically in golang.

So the next best this is to us the existing clo cli to snapshot the relevant state in a consumable format.

The state is represented in `clo/models/models*go`.  Example data is in `clo/testdata/keystone_nops.json`

First, ensure you can login to the CLO instances [stage](https://feeds-manager.main.stage.cldev.sh/sign-in), [prod](https://feeds-manager.main.prod.cldev.sh/sign-in) (ask in #topic-keystone-clo and tag Joey Punzel)

Next, you need the clo cli:
See to build and install
https://github.com/smartcontractkit/feeds-manager#chainlink-orchestrator-api-client

Now, you need config for stage and prod env, eg `~/.fmsclient/stage.yaml` and `~/.fmsclient/prod.yaml`

`~/.fmsclient/prod.yaml` :
```
EMAIL: "your.name@smartcontract.com"
PASSWORD: 'XXXredacted'
BASE_URL: "https://gql.feeds-manager.main.prod.cldev.sh"
```
`~/.fmsclient/stage.yaml`
```
EMAIL: "your.name@smartcontract.com"
PASSWORD: 'XXXredacted'
BASE_URL: "https://gql.feeds-manager.main.stage.cldev.sh"
```



Now run the cli to get *all* the node operators
```
./bin/fmscli --config ~/.fmsclient/prod.yaml get nodeOperators > /some/file.json
```

The output of this will be a JSON serialization of `[]*models.NodeOperator` and should same from as the testdata `clo/testdata/keystone_nops.json` This test data was post filtered to only contain keystone node operators.

In order to make the data in `/some/file` useful for a deployment, you need to filter to only contain `keystone` nodes. This can be heuristically with `CapabilityNodeSet` func. See the test in `clo/don_nodeset_test.go` for an example.


This filtered data, the map of don -> []capabilities is enough to fully specify the offchain data required for a deployment.

See `deploy_test:Test_Deploy/memory_chains_clo_offchain` for an explicit example

## Chain configuration

All the tests use in memory chains with simulated backends. Real deployments are in `chainlink-deployment`, which load chain-specific configuration (such as rpc endpoints) to instantiate real chains.


## Jobs
Are not handled programmatically yet. They are managed manually in CLO with help from [RDD](https://github.com/smartcontractkit/reference-data-directory#workflows) 