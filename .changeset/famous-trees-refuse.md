---
"chainlink": patch
---

Support new report type 'evm_streamlined'. #added

This new report type is designed to be as small and optimized as possible to minimize report size and calldata.

Reports are encoded as such:

(no FeedID specified in opts)
```
<32 bits> channel ID
<64 bits> unsigned report timestamp nanoseconds
<bytes>   report data as packed ABI encoding
```

(FeedID specified in opts)
```
<256 bits> feed ID
<64 bits> unsigned report timestamp nanoseconds
<bytes>   report data as packed ABI encoding
```

Report contexts are encoded as such:

```
// Equivalent to abi.encodePacked(digest, len(report), report, len(sigs), sigs...)
// bytes32 config digest
// packed uint16 len report
// packed bytes report
// packed uint8 len sigs
// packed bytes sigs
```

See report_codec_evm_streamlined_test.go for examples.

