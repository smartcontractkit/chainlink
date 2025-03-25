---
"chainlink": patch
---

Implement support for TimestampedStreamValue data types in LLO (RWAs) #added

Support encoding into evm_abi_unpacked or evm_streamlined report formats.

ABI must specify how to encode both types, as such:

```json
// Encodes the timestamp as uint64 and data payload as int192
{
    "abi": [[{"type":"uint64"},{"type":"int192"}]],
}
```

The first element of the array encodes the timestamp, the second encodes the data payload.

Users may suppress one or the other entirely by using the special keyword "bytes0" e.g.

```json
// Encodes only the data payload
{
    "abi": [[{"type":"bytes0"},{"type":"int192"}]],
}
```
