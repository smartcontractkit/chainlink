# LogTrigger registration: EOA / invalid-address misclassified as system error

Context for a bug where registering a LogTrigger against a non-contract address
(EOA) is retried indefinitely as a **system error**, even though the root cause is
**user input** (bad workflow trigger config).

Observed log (staging Sepolia):

```
trigger registration failed with system error; will retry
err="[14]Unavailable: failed to register log-tracking: 'rpc error: code = Unknown desc = [3]InvalidArgument: one or more addresses are not contracts: index 0: 0x5F14A6f6Cb2b6BBaC603277f4Cf498D9C6464883 (EOA, not a contract)' for triggerID: ..."
```

Logger: `TriggerPublisher` at `trigger_publisher.go:309`.

---

## Expected behavior

| Origin | `TriggerPublisher` behavior |
|--------|----------------------------|
| **User error** (bad address, bad topics, etc.) | Store `registrationErr`, log "will **not** retry", short-circuit future quorum attempts |
| **System error** (RPC down, internal failure) | Do **not** store registration, log "will retry", call `underlying.RegisterTrigger()` again on every 2F+1 quorum |

EOA-not-a-contract is a workflow configuration mistake → should be **user error**, no retry.

---

## Call chain

### 1. Correct classification at source

**File:** `chainlink-evm/pkg/logpoller/log_poller.go` — `RegisterFilter()`

For each filter address, `CodeAt` is called. Empty bytecode → EOA:

```go
if len(code) == 0 {
    invalidAddresses = append(invalidAddresses,
        fmt.Sprintf("index %d: %s (EOA, not a contract)", i, addr.Hex()))
}
if len(invalidAddresses) > 0 {
    return caperrors.NewPublicUserError(addressesError, caperrors.InvalidArgument)
}
```

This is the right origin (`OriginUser`, code `InvalidArgument`).

### 2. EVM relay — pass-through

**File:** `chainlink-evm/pkg/relay/evm_service.go` — `RegisterLogTracking()`

```go
return e.chain.LogPoller().RegisterFilter(ctx, lpfilter)
```

No reclassification.

### 3. gRPC boundary — type information lost

Capability plugin and chain node communicate over gRPC (loop relayer).

**Server:** `chainlink-common/pkg/loop/internal/relayer/evm.go` — `evmServer.RegisterLogTracking()`
returns the `caperrors.Error` as a plain Go `error`.

**Client:** `EVMClient.RegisterLogTracking()` → `net.WrapRPCErr(err)`.

`caperrors.Error` does **not** implement `GRPCStatus()`. gRPC therefore maps it to:

```
rpc error: code = Unknown desc = [3]InvalidArgument: one or more addresses are not contracts: ...
```

Notes:

- `[3]InvalidArgument` is `capabilityError.Error()` string format (`[code]Name: msg`), not gRPC `codes.InvalidArgument`.
- `WrapRPCErr` preserves the gRPC status wrapper; it does **not** restore a `caperrors.Error` type on the client.

### 4. Re-wrapped as system error — EVM LogTrigger

**File:** `capabilities/chain_capabilities/evm/trigger/trigger.go` — `RegisterLogTrigger()`

```go
if err = lts.EVMService.RegisterLogTracking(ctx, filterQuery); err != nil {
    registerError := fmt.Errorf("failed to register log-tracking: '%w' for triggerID: %s, ...", err, triggerID, ...)
    var lpError caperrors.Error
    if errors.As(err, &lpError) {
        if lpError.Origin() == caperrors.OriginUser {
            return nil, caperrors.NewPublicUserError(registerError, lpError.Code())
        }
    }
    // errors.As fails on gRPC-wrapped error → always hits this branch in production
    return nil, caperrors.NewPublicSystemError(registerError, caperrors.Unavailable)
}
```

`errors.As(err, &lpError)` fails because the client error is a gRPC `wrappedError`, not
`caperrors.Error`. The user error from log poller is therefore upgraded to
`[14]Unavailable` / `OriginSystem`.

Unit test `fail to register log-tracking user error` in `trigger_test.go` mocks
`EVMService.RegisterLogTracking` returning `caperrors.Error` **directly** (no gRPC) —
that path works; production does not.

### 5. Retry decision — TriggerPublisher

**File:** `trigger_publisher.go` — `Receive()` / `MethodRegisterTrigger`

After 2F+1 quorum, `underlying.RegisterTrigger()` is called. On error:

```go
if errors.As(registerErr, &capErr) && capErr.Origin() == caperrors.OriginUser {
    p.registrations[key] = &pubRegState{registrationErr: registerErr} // will NOT retry
} else {
    // line 309 — system error path, NO entry in p.registrations
    p.lggr.Errorw("trigger registration failed with system error; will retry", ...)
}
```

Because step 4 produced `OriginSystem`:

- Log says "will retry"
- `(workflowID, triggerID)` is **not** stored in `p.registrations`
- Every subsequent 2F+1 quorum calls `underlying.RegisterTrigger()` again (slow path)
- Contributes to LogTrigger receiver channel load alongside registration refresh traffic

---

## Why the nested `[3]InvalidArgument` does not save us

`caperrors.DeserializeErrorFromString()` expects the **serialization** format:

```
Public:User:InvalidArgument:<message>
```

The string embedded in the gRPC error uses the **Error()** format:

```
[3]InvalidArgument: one or more addresses are not contracts: ...
```

So naive string deserialization on the gRPC message will not recover origin either.

---

## Impact

1. **Infinite retry** on permanently invalid trigger configs (EOA address, etc.)
2. **Slow-path registration** on every quorum (not the cheap "already exists" refresh path)
3. **LogTrigger channel pressure** — extra work competing with `TriggerEventAck` processing
4. **Misleading ops signal** — `platform_trigger_publisher_register_trigger_total{outcome="error"}` with system-error log pattern

---

## Fix directions (not implemented here)

Pick one or combine:

1. **`trigger.go`:** Unwrap gRPC errors and recover `caperrors.Error` (e.g. parse `[N]CodeName:` prefix from status message, or add a dedicated unwrap helper in `chainlink-common/pkg/loop/internal/net`).
2. **gRPC layer:** Implement `GRPCStatus()` on `caperrors.Error` (or an interceptor) so user errors cross the hop with correct gRPC code / metadata.
3. **`trigger_publisher.go` (defensive):** Treat known user-error substrings (`not a contract`, `EOA`, etc.) as user errors — last resort, fragile.

Preferred: fix at (1) or (2) so all EVM service user errors propagate correctly, not only EOA.

---

## Key file references

| Layer | File |
|-------|------|
| EOA detection | `chainlink-evm/pkg/logpoller/log_poller.go` (~311–326) |
| Relay | `chainlink-evm/pkg/relay/evm_service.go` (~166–176) |
| gRPC client | `chainlink-common/pkg/loop/internal/relayer/evm.go` (~260–262) |
| gRPC error wrap | `chainlink-common/pkg/loop/internal/net/errors.go` — `WrapRPCErr` |
| Re-classification | `capabilities/chain_capabilities/evm/trigger/trigger.go` (~301–313) |
| Retry policy | `trigger_publisher.go` (~300–311) |
| cap error format | `chainlink-common/pkg/capabilities/errors/error.go` — `Error()`, `Origin()` |
