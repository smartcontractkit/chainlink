# libocr

libocr consists of a Go library and a set of Solidity smart contracts that implement the *Chainlink Offchain Reporting Protocol*, a [Byzantine fault tolerant](https://en.wikipedia.org/wiki/Byzantine_fault) protocol that allows a set of oracles to generate *offchain* an aggregate report of the oracles' observations of some underlying data source. This report is then transmitted to an onchain contract in a single transaction.

You may also be interested in [libocr's integration into the actual Chainlink node](https://github.com/smartcontractkit/chainlink/tree/develop/core/services/offchainreporting).


## Protocol Description

Protocol execution mostly happens offchain over a peer to peer network between Chainlink nodes. The nodes regularly elect a new leader node who drives the rest of the protocol. The protocol is designed to choose each leader fairly and quickly rotate away from leaders that aren’t making progress towards timely onchain reports.

The leader regularly requests followers to provide freshly signed observations and aggregates them into a report. It then sends the aggregate report back to the followers and asks them to attest to the report's validity by signing it. If a quorum of followers approves the report, the leader assembles a final report with the quorum's signatures and broadcasts it to all followers.

The nodes then attempt to transmit the final report to the smart contract according to a randomized schedule. Finally, the smart contract verifies that a quorum of nodes signed the report and exposes the median value to consumers.


## Organization
```
.
├── contract: Ethereum smart contracts
├── gethwrappers: go-ethereum bindings for the OCR1 contracts, generated with abigen
├── gethwrappers2: go-ethereum bindings for the OCR2 contracts, generated with abigen
├── networking: p2p networking layer
├── offchainreporting: offchain reporting protocol version 1
├── offchainreporting2: offchain reporting protocol version 2
├── permutation: helper package for generating permutations
└── subprocesses: helper package for managing go routines
```
