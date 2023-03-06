/*
Package dhtrouter implements a router that uses a DHT to store and propagate peers' address information.


Functionality

Very simply, this package allows nodes to announce their own addresses and look up other nodes' addresses.

An unexported method DHTRouter.publishHostAddr(ctx context.Context) publishes the node's own address. Specifically,
it builds and puts an announcement containing (peerid, pk, addresses, version, sig) to DHT,
where peerid = hash(pk) or peerid = pk, depending on ID format; and sig is a signature over addresses and version.

Peers validate announcements before forwarding them. A node adds (pk, addresses, version) to its DHT data store
and forwards it if 1) sig verifies against pk, 2) version is higher than the highest version it knows.
In the current implementation, the version field is simply a db backed monotonic counter.

As for exported methods, package dhtrouter provides an implementation of the interface PeerDiscoveryRouter,
which extends rhost.Routing.


Security enhancements

Compared to the built-in peer discovery protocol in libp2p/Kademlia,
package dhtrouter achieves stronger security in a few important ways.



The main problems with the built-in peer discovery protocol in libp2p/Kademlia is that it's too open, as it's designed
to be permissionless. We however are in a permissioned setting where we must ensure that unauthorized nodes cannot join
the system and only learn minimal information about the system.

In particular, libp2p/Kademlia will automatically add all peers seen over the network to their local store, and share
peer information with any nodes can connect. We cannot afford to leak such info to unauthorized nodes.

To this end, we made the following enhancements: 1) we added a wrapper around hosts that only permit DHT streams
from/to permitted IDs; 2) we leverage the filtering options in DHT so an honest node would only query permitted IDs and
only keep permitted IDs in the routing table; 3) we added a wrapper around the TLS transport to conceal node identity.
(See knockingTLS package for details.)


Security features

We claim the following security features:

* Integrity: A malicious node cannot spoof an honest one, as announcements are authenticated.

* Liveness: honest nodes can publish their address eventually, under the assumptions that an honest node can connect to
at least f+1 peers. Further, we set bucket size to be large enough so that all nodes are in the same bucket and data is
replicated by all peers.

* Peer privacy: for any peer A, only peers permitted by A can learn peer A's peers.
Further, we support namespace based permission. If peer A, B, C are in namespace NS1, and B, C, D are in namespace NS2,
peer A learns B's peers include {A, C}, but cannot learn that D is a peer of B.
*/

package dhtrouter
