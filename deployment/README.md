## Overview
The deployment Go module serves as a product agnostic set of environment
abstractions used to deploy and configure products including both on/offchain
dependencies. The environment abstractions allow for complex and critical
deployment/configuration logic to be tested against ephemeral environments
and then exposed for use in persistent environments like testnet/mainnet.

## Table of Contents
- [Address Book](#address-book)
- [View](#view)
- [Environment](#environment)
- [Job Distributor](#job-distributor)
- [Changesets](#changsets)
- [Directory Structure](#directory-structure)
- [Integration Testing](#integration-testing)
- [FAQ](#faq)

## Address Book
An [address book](https://github.com/smartcontractkit/chainlink/blob/develop/deployment/address_book.go#L79) represents
a set of versioned, onchain addresses for a given product across all blockchain families. The primary key 
is a family agnostic [chain-selector](https://github.com/smartcontractkit/chain-selectors) chain identifier combined with a unique 
address within the chain. Anything which is globally addressable on the chain can be used for the address field, for example
EVM smart contract addresses, Aptos objectIDs/accounts, Solana programs/accounts etc. 
The address book holds the minimum amount of information to derive the onchain state of the system from 
the chains themselves as a source of truth.  

It is recommended that you define a State struct holding Go bindings for each onchain component and author a
translation layer between the address book and the state struct. Think of it like an expanded/wrapped address book  
which enables read/writing to those objects. See an example [here](https://github.com/smartcontractkit/chainlink/blob/develop/deployment/ccip/state.go#L205).
This way, given an address book, you can easily read/write state or generate a [View](##View).
Note that for contract upgrades its expected that you would have versioned field names in the State struct until the v1 is fully removed, this
way you can easily test upgrades and support multiple versions of contracts. 

## View
A [view](https://github.com/smartcontractkit/chainlink/blob/develop/deployment/changeset.go#L35) is a function which
serializes the state of the system into a JSON object. This is useful for exporting to other systems (like a UI, docs, DS&A etc).
You can generate it however you see fit, but a straightforward way is to translate
the address book into a State structure and then serialize that using Go bindings. 

## Environment
An [environment](https://github.com/smartcontractkit/chainlink/blob/develop/deployment/environment.go#L71) represents
the existing state of the system including onchain (including non-EVMs) and offchain components. Conceptually it contains
a set of pointers and interfaces to dereference those pointers from the source of truth.
The onchain "pointers" are an address book of existing addresses and the Chains field holds
clients to read/write to those addresses. The offchain "pointers" are a set of nodeIDs and 
the Offchain client (interface to the [job-distributor](##Job Distributor)) to read/write to them.

## Job Distributor
The job distributor is a product agnostic in-house service for 
managing jobs and CL nodes. It is required to use if you want to 
manage your system through chainlink deployments.

## Changsets
A [changeset](https://github.com/smartcontractkit/chainlink/blob/develop/deployment/changeset.go#L21) is a 
Go function which describes a set of changes to be applied to an environment given some configuration:
```go
type ChangeSet func(e Environment, config interface{}) (ChangesetOutput, error)
```
For example, changesets might include:
- Deploying a new contract
- Deploying 2 contracts where the second contract depends on the first's address
- Deploying a contract and creating a job spec where the job spec points to the deployed contract address
- Creating an MCMS proposal to set a billing parameter on a contract 
- Deploying a full system from scratch
  - Mainly useful for integration tests
- Modifying a contract not yet owned by MCMS via the deployer key
 
Once sufficient changesets are built and tested, the ongoing maintenance
of a product should be just invoking existing changesets with new configuration.
The configuration can be environment/product specific, for example
specific chain addresses, chain selectors, data sources etc. The outputs are 
a set of diff artifacts to be applied to the environment (MCMS proposals, job specs, addresses created).  
You can use the changeset for side effects only and return no artifacts.
An example would be making an onchain change with the deployer key instead of an MCMS proposal, 
however that should generally be uncommon. Usually we'd expect an initial deployment to produce
a set of addresses and job specs (likely pointing to those addresses) and then from that point forward
we'd expect to use MCMS proposals to make changes. 

TODO: Add various examples in deployment/example.

## Directory structure

/deployment
- package name `deployment`
- Product agnostic environment abstractions and helpers using those
  abstractions

/deployment/environment/memory
- package name `memory`
- In-memory environment for fast integration testing
- EVM only

/deployment/environment/devenv
- package name `devenv`
- Docker environment for higher fidelity testing
- Support non-EVMs (yet to be implemented)

/deployment/common
- Deploymnet/configuration/view logic for product agnostic
contracts (like MCMS, LinkToken etc) which can be shared
by products.

/deployment/product/internal
- Internal building blocks for changesets  

/deployment/product/view
- Hold readonly mappings Go bindings to json marshallable objects. 
- Used to generate a view of the system.

/deployment/product/changeset
- Think of this as the public API for deployment and configuration
of your product. 
- All the changesets should have an associated test using a memory or devenv
environment.
- package name `changeset` imported as `<package>changeset`
 
## Integration testing
Integration tests should live in the integration-tests/go.mod module and leverage
the deployment module for product deployment and configuration. The integration tests
should only depend on deployment/<product>/changeset and deployment/environment. 

## FAQ
### Should my changeset be idempotent? 
It depends on the use case and is at your discretion. In many cases
it would be beneficial to make it idempotent so that if anything goes wrong
you can re-run it without side effects. However, it's possible that the onchain contract
design doesn't allow for idempotency so in that case you'd have to be prepared
with recovery changesets if something goes wrong as re-running it would not be an option.
