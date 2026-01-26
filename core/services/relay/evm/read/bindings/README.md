This folder contains CCIP gobindings copied (unused code has been trimmed) from v1.6.0 version 
(https://github.com/smartcontractkit/chainlink-ccip/blob/main/chains/evm/gobindings/generated/v1_6_0/offramp/offramp.go and
https://github.com/smartcontractkit/chainlink-ccip/blob/main/chains/evm/gobindings/generated/v1_6_0/onramp/onramp.go)
They have been inlined here to avoid dependency on chainlink-ccip.

The root cause of importing these gobindings is to have a more efficient event decoding (https://smartcontract-it.atlassian.net/browse/CCIP-5348). 
Once CCIP is migrated to v1.7.0 this optimization (and associated gobindings here) can be removed.