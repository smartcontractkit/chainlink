/*
The relay package contains the components necessary to implement the "autorelay"
feature.

Warning: the internal interfaces are unstable.

System Components:
- A discovery service to discover public relays.
- An AutoNAT client used to determine if the node is behind a NAT/firewall.
- One or more autonat services, instances of `AutoNATServices`. These are used
  by the autonat client.
- One or more relays, instances of `RelayHost`.
- The AutoRelay service. This is the service that actually:

AutoNATService: https://github.com/libp2p/go-libp2p-autonat-svc
AutoNAT: https://github.com/libp2p/go-libp2p-autonat

How it works:
- `AutoNATService` instances are instantiated in the bootstrappers (or other
  well known publicly reachable hosts)
- `AutoRelay`s are constructed with `libp2p.New(libp2p.Routing(makeDHT))`
  They passively discover autonat service instances and test dialability of
  their listen address set through them.  When the presence of NAT is detected,
  they discover relays through the DHT, connect to some of them and begin
  advertising relay addresses.  The new set of addresses is propagated to
  connected peers through the `identify/push` protocol.
*/
package relay
