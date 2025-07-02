# PoR: an implementation of the OCR3 plugin that enables the safe decision (and distribution) of mintable amounts per chain for PoR.

Background:
There are multiple tracked chains, e.g., `chains = [chain_A, chain_B, ...]`

The EA, at each oracle, communicates (reads/interacts) with these chains, in an attempt to obtain the latest `blockNumber` for each chain, and calculate mintable information for a particular query (a query is a mapping of `chain -> blockNumber` pairs).

The goal of this plugin is to generate honest reports based on new information from multiple EAs.
First it obtains information from the EAs on the latest status on the chains being tracked, namely the latest block number they are aware of. Based on this, the plugin settles on a query: a map of (chain -> safe block number) pairs (a block number is safe if is is guaranteed that chain has at least that many blocks).
Second, it queries the oracles for information (mintable amount + a block number) to report chain by chain.
These two steps are done in a pipelined fashion in every round of the plugin.s

General plugin workflow:

Every round of consensus has:
1. The leader
2. The followers (other oracles).

Every round of consensus:

1. The leader runs the `Query` method.
2. Followers (and leader) run the `Observation(...)` method.
3. The `ValidateObservation(query, observation)` method is run for each observation.
4. The `ObservationQuorum(observations)` method is (continuously) run on increasing sets of observations received by the leader (until it returns true)
5. The `Outcome(observations)` method is run on every oracle on a set of `observations` s.t. `ObservationQuorum(observations) = true`
6. The `Reports(...)` method is run on every oracle
7. The `ShouldAcceptAttestedReport` method is run for each attested report on every oracle when the report attestation is gathered.
8. The `ShouldTransmitAcceptedReport` method is run for each attested report, which is not filtered out by `ShouldAcceptAttestedReport` on every oracle, right before the oracle sends the report for transmission.        