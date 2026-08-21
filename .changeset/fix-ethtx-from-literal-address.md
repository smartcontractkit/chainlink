---
"chainlink": patch
---

#bugfix Reject inputs that the JSON decoder only partially consumes in the
pipeline `JSONWithVarExprs` getter. Previously, a literal Ethereum address
like `0x...` in an `ethtx` task's `from` field would be parsed as the JSON
number `0`, returned as `int64(0)`, and fail downstream with
`AddressSliceParam: cannot convert int64`. The getter now requires the
decoder to consume the entire input, letting `NonemptyString` handle the
literal address as intended. Resolves #21768.
