### Overview
The deployment package in the integration-tests Go module serves
as a product agnostic set of environment abstractions used
to deploy and configure products including both on/offchain
dependencies. The environment abstractions allow for
complex and critical deployment/configuration logic to be tested
against ephemeral environments and then exposed for use in persistent
environments like testnet/mainnet.

### Directory structure

/deployment
- package name `deployment`
- Product agnostic environment abstractions and helpers using those
  abstractions

/deployment/memory
- package name `memory`
- In-memory environment for fast integration testing
- EVM only

/deployment/devenv
- Coming soon
- package name `devenv`
- Docker environment for higher fidelity testing
- Support non-EVMs (yet to be implemented)

/deployment/ccip
- package name `ccipdeployment`
- Files and tests per product deployment/configuration workflows
- Tests can use deployment/memory for fast integration testing
- TODO: System state representation is defined here, need to define
  an interface to comply with for all products.

/deployment/ccip/changeset
- package name `changeset` imported as `ccipchangesets`
- These function like scripts describing state transitions
  you wish to apply to _persistent_ environments like testnet/mainnet
- Ordered list of Go functions following the format
```Go
0001_descriptive_name.go
func Apply0001(env deployment.Environment, c ccipdeployment.Config) (deployment.ChangesetOutput, error)
{
    // Deploy contracts, generate MCMS proposals, generate
    // job specs according to contracts etc.
    return deployment.ChangesetOutput{}, nil
}
0001_descriptive_name_test.go
func TestApply0001(t *testing.T)
{
    // Set up memory env
    // Apply0001 function
    // Take the artifacts from ChangeSet output
    // Apply them to the memory env
    // Send traffic, run assertions etc.
}
```
- Changesets are exposed and applied via a different repo. 