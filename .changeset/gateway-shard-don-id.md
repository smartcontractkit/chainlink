---
"chainlink": patch
---

Gateway: DON IDs for shards N>0 of a sharded DON now use the format `<DonName>_shard-<N>` (previously `<DonName>_<N>`), matching the naming used for shard DONs elsewhere. Shard 0 still uses the bare DON name. Nodes in shards N>0 must set `GatewayConnector.DonID` accordingly. #changed
