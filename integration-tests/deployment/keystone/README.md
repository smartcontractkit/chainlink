# Programmatic deployment of Keystone

The current scope of this package is the ability to deploy and configure the Capability Registry, OCR3 and Forwarder contracts.

It builds on the `Environment` abstraction introduced by `chainlink-deployments`. The concept of Environment is used delineate dev vs testnet vs prod.
A deployment injects the necessary configuration (eg simulated chain vs testnet chain) as appropriate via the `Environment` abstraction, and the
deployment implementation be agnostic to these details.
 

The entry point to the deployment is the `Deploy` func. The arguments to the this func are environment dependent, and vary from one deployment to another.

```
type DeployRequest struct {
	RegistryChainSel uint64
	Env              *deployment.Environment

	Dons       []DonCapabilities   // externally sourced based on the environment
	OCR3Config *OracleConfigSource // TODO: probably should be a map of don to config; but currently we only have one wf don therefore one config
}

type DeployResponse struct {
	Changeset *deployment.ChangesetOutput
	DonInfos  map[string]capabilities_registry.CapabilitiesRegistryDONInfo
}

func Deploy(ctx context.Context, lggr logger.Logger, req DeployRequest) (*DeployResponse, error) 
```


In order to make this all work we need a mapping what nodes run which capabilities, which nodes belong to what don, and ocr configuration for consensus. The first two are represented by `Dons, OCR3Config`, respectively.

The mapping for nodes->capability is an external artifact that is declare in configuration for the given environment. The mapping between nodes and Dons is also configuration, however it is constrained by 
real world data about the nodes themselves, such as the p2pkeys and so forth.

For keystone, this constraint boils down to an integration point with CLO, which is the current system of record all Node/NOP metadata (as well the Jobs themselves).

Therefore, in order to the system to work, we need to source data from CLO.

# CLO integration

The integration with CLO is contained `clo` package. It defines a minimal, keystone-specific, translation of the CLO data model to the new Job Distributor model. This is needed because the `Environment` abstraction relies on the JD data model and API (via `OffchainClient`).

However, at the time of writing, it was not feasible to programmatically access the CLO API within our deployment (KS-454).

For the time being, there are manual steps to obtain the metadata from CLO as described below.

## Obtaining and parsing CLO data

A real deployment requires real data from CLO.

### Requirements
- clo access to [stage](https://feeds-manager.main.stage.cldev.sh/sign-in), [prod](https://feeds-manager.main.prod.cldev.sh/sign-in) (ask in #topic-keystone-clo and tag Joey Punzel)
- [clo cli](https://github.com/smartcontractkit/feeds-manager#chainlink-orchestrator-api-client)
- stage & prod configuration for the cli


As discussed above, CLO is the system that knows about nodes and node operators. One of our goals is to configure the registry contract with the nodes and nops. So we have to faithfully plumb the values in CLO to our deployment.

For the time being, it is not possible to do this programmatically in golang.T he next best this is to us the existing clo cli to snapshot the relevant state in a consumable format.

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