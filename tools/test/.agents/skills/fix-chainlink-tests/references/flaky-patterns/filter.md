<symptom>
The dominant flake pattern in simulated-chain tests that enable `Feature.LogPoller = true`. Error message contains `"failed to retrieve log value pointer of block N: not found"` and the stack trace points to a `FilterXxx` call that immediately follows a `backend.Commit()`. Note: Raw geth bindings do NOT have this race, only interface types backed by LogPoller.
</symptom>

<fix_a_receipt_parsing>

For one-shot events where you only need a value emitted at creation (e.g. `SubscriptionCreated`, `RequestSent`): parse the tx receipt directly instead of calling `FilterXxx`.
```go
// AFTER (deterministic):
tx, err := coordinator.CreateSubscription(auth)
require.NoError(t, err)
backend.Commit()
receipt, err := backend.Client().TransactionReceipt(ctx, tx.Hash())
require.NoError(t, err)
require.Equal(t, uint64(1), receipt.Status)
var subID *big.Int
for _, log := range receipt.Logs {
    if log.Address != coordinatorAddress {
        continue
    }
    // SubscriptionCreated(uint64 indexed subId, address owner): Topics[1] = subId
    subID = new(big.Int).SetBytes(log.Topics[1].Bytes())
    break
}
require.NotNil(t, subID, "no SubscriptionCreated log in receipt")
```

</fix_a_receipt_parsing>

<fix_b_non_fatal_filter>

For diagnostic/verification filters called inside a polling loop: a transient LogPoller error must not crash the test — it should retry.
```go
// AFTER (retries):
require.Eventually(t, func() bool {
    // LogPoller may not have indexed the latest block yet; skip and retry.
    it, err := coordinator.FilterRandomWordsForced(nil, ids, subs, addrs)
    if err == nil {
        for it.Next() {
            require.Equal(t, expected, it.Event.Field)
        }
    }
    return utils.IsEmpty(commitment[:])
}, timeout, tick)
```

</fix_b_non_fatal_filter>

<fix_c_dynamic_reference>

If `require.Eventually` commits new blocks on each iteration, compute the reference block number inside the closure so it doesn't become stale.
```go
// AFTER (dynamic):
require.Eventually(t, func() bool {
    backend.Commit()
    tip, err := backend.Client().HeaderByNumber(ctx, nil)
    if err != nil || tip == nil || tip.Number.Uint64() < 256 {
        return false
    }
    _, err = bhsContract.GetBlockhash(nil, new(big.Int).SetUint64(tip.Number.Uint64()-256))
    return err == nil
}, testutils.WaitTimeoutCustom(t, 5*time.Minute), time.Second)
```

</fix_c_dynamic_reference>