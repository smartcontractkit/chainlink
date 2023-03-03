/*
Package knockingTLS is a wrapper around the TLS transport to provide server identity privacy from unauthorized clients.

In TLS, the server reveals its identity (public key) to connecting clients before authenticate the client, which could
cause privacy issues. This package mandates the client to authenticate to the server first, by sending a
small "knock" as the first message from the client before a handshake takes place.

In the current implementation, a knock is simply a 64 byte Ed25519 signature, which adds minimal overhead.
This wrapper works with all version of TLS.
*/
package knockingtls
