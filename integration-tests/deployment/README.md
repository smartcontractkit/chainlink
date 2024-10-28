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
- They should be go functions where the first argument is an
  environment and the second argument is a config struct which can be unique to the 
  changeset. The return value should be a `deployment.ChangesetOutput` and an error.
```Go
do_something.go
func DoSomethingChangeSet(env deployment.Environment, c ccipdeployment.Config) (deployment.ChangesetOutput, error)
{
    // Deploy contracts, generate MCMS proposals, generate
    // job specs according to contracts etc.
    return deployment.ChangesetOutput{}, nil
}
do_something_test.go
func TestDoSomething(t *testing.T)
{
    // Set up memory env
    // DoSomethingChangeSet function
    // Take the artifacts from ChangeSet output
    // Apply them to the memory env
    // Send traffic, run assertions etc.
}
```
- Changesets are exposed and applied via a different repo. 

/deployment/llo
- package name `llodeployment`
- Similar to /deploymet/ccip, these are product-specific deployment/configuration workflows
- Tests can use deployment/memory for fast integration testing

/deployment/llo/changeset
- package name `changeset` imported as `llochangesets`
- Similar to deployment/ccip/changesets
- These function like scripts describing state transitions
  you wish to apply to _persistent_ environments like testnet/mainnet
