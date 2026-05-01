---
"chainlink": patch
---

#internal

Skip flaky `TestRMN_TwoMessagesOneSourceChainCursed` in `integration-tests/smoke/ccip/ccip_rmn_test.go`. The test is tracked at #18031 and matches the pattern already used for other quarantined RMN tests in the same file.
