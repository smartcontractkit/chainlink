/*
The simulated backend cannot access old blocks and will return an error if
anything other than `latest`, `nil`, or the latest block are passed to
`CallContract`.

The simulated client avoids the old block error from the simulated backend by
passing `nil` to `CallContract` when calling `CallContext` or `BatchCallContext`
and will not return an error when an old block is used.
*/
package client
