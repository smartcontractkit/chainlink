---
"chainlink": patch
---

#bugfix Bridge names are now treated as case-insensitive everywhere they are
read. The pipeline runtime was casting the spec-provided name straight to
`bridges.BridgeName` and looking it up as-is, so a job referencing a bridge
created with mixed case (e.g. `OpenWeatherMap`) would fail at run time with
`could not find bridge with name`. The GraphQL `createBridge` mutation had
the same issue on the write path. Both now route the name through
`bridges.ParseBridgeName`, matching the JSON-API behavior and the documented
contract that bridge names are case-insensitive. Resolves #9785.
