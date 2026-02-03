# CHIP Test Sink

A lightweight gRPC server that implements the `ChipIngress` interface so system tests can assert on CloudEvents without touching production ingress.

## Typical Flow

1. **Create channels** for the event types you care about (e.g., user logs or workflow base messages).
2. **Build a publish handler**. If you only need user logs and base messages, you can reuse `t_helpers.GetPublishFn`, which demuxes those two event types into channels. For any other event types, provide your own `PublishFn` and fan out as needed. When you demux N types but only assert on M (< N), make sure to drain the unused channels (e.g., with `IgnoreUserLogs`) so the sink never blocks once buffers fill.
3. **Start the sink** using `t_helpers.StartChipTestSink`, which blocks until the listener is ready.
4. **Wire assertions** with helpers such as `WaitForUserLog`, `WatchWorkflowLogs`, `WaitForBaseMessage`, or `WatchBaseMessages`.
5. **Shut down** the sink in `t.Cleanup` to release the listener.

The sink can optionally forward events to a real Chip Ingress endpoint by setting `Config.UpstreamEndpoint`.

## Helper Summary

- `StartChipTestSink(t, publishFn)` – boots the server on the default Chip Ingress port and waits for readiness.
- `GetPublishFn(logger, userLogsCh, baseMessageCh)` – returns a `chiptestsink.PublishFn` that demultiplexes CloudEvents into typed channels.
- `WaitForUserLog(ctx, logger, userLogsCh, needle)` – blocks until a user log containing `needle` arrives or the context ends.
- `WaitForBaseMessage(ctx, logger, baseMessageCh, needle)` – same as above for base messages.
- `WatchWorkflowLogs(t, logger, userLogsCh, baseMessageCh, poisonNeedle, successNeedle, timeout)` – combines timeout + cancel-cause logic so tests fail fast on poison logs and pass when the expected user log appears.
- `WatchBaseMessages(t, logger, baseMessageCh, needle, timeout)` – asserts that a base message shows up before the timeout.
- `FailOnBaseMessage(ctx, cancelCause, t, logger, baseMessageCh, needle)` – cancels the shared context immediately once the poison message is observed.
- `IgnoreUserLogs(ctx, userLogsCh)` – drains user log traffic when a test does not care about it, preventing channel back-pressure.
- `WaitForServerStart(t, started, errCh)` – utility used by the sink starter to ensure the listener is bound before tests proceed.

## Example Snippet

```go
userLogsCh := make(chan *workflowevents.UserLogs, 1000)
baseMessageCh := make(chan *commonevents.BaseMessage, 1000)
publishFn := t_helpers.GetPublishFn(testLogger, userLogsCh, baseMessageCh)
server := t_helpers.StartChipTestSink(t, publishFn)
t.Cleanup(func() { server.Shutdown(t.Context()) })

t_helpers.WatchWorkflowLogs(t, testLogger, userLogsCh, baseMessageCh, "Workflow Engine initialization failed", "expected user log", 2*time.Minute)
```

This setup mirrors how current CRE regression and smoke tests consume the test sink to validate workflow behavior.
