// ragep2p is a simple p2p networking library that provides best-effort,
// self-healing, message-oriented, authenticated, encrypted bidirectional
// communication streams between peers.
//
// Concepts
//
// ragep2p peers are identified by their PeerID. PeerIDs are public keys.
// PeerIDs' string representation is compatible with libp2p to ease migration
// from libp2p to ragep2p.
//
// ragep2p provides Host and Stream abstractions.
//
// A Host allows users to establish Streams with other peers identified by their
// PeerID.  The host will transparently handle peer discovery, secure connection
// (re)establishment, multiplexing streams over the connection and rate
// limiting.
//
// A Stream allows users to send binary messages to another peer. Messages are
// delivered on a best-effort basis. If the underlying connection to the other
// peer drops or isn't fast enough or the other peer has not opened a
// corresponding stream or ..., messages may get dropped. Streams are
// self-healing: users don't need to close and re-open streams even if the
// underlying connection drops or the other peer becomes unavailable. We
// guarantee that messages that are delivered are delivered in FIFO order and
// without modifications.
//
// Peer discovery
//
// ragep2p will handle peer discovery (i.e. associating network addresses like
// 1.2.3.4:1337 with PeerIDs) automatically. Upon construction of a Host, a
// Discoverer is passed in, which is then used by the Host for this purpose.
//
// If multiple network addresses are discovered for a PeerID, ragep2p will try
// sequentially dialing all of them until a connection is successfully
// established.
//
// Thread Safety
//
// All public functions on Host and Stream are thread-safe.
//
// It is safe to double-Close(), though all but the first Close() may return an
// error.
//
// Allocations
//
// We allocate a buffer for each message received. In principle, this could allo
// an adversary to force a recipient to run out of memory. To defend against
// this, we put limits on the length of messages and rate limit messages,
// thereby also limiting adversarially-controlled allocations.
//
// Security
//
// Note: Users don't need to care about the details of how these security
// measures work, only what properties they provide.
//
// ragep2p's security model assumes that all Streams on the local Host behave
// honestly and cooperatively. Since many Streams are multiplexed over a single
// connection, a single "bad" Stream could completely exhaust the entire
// connection preventing other Streams from delivering messages as well. Other
// network participants, however, are not assumed to behave honestly and we
// attempt to defend against fingerprinting, impersonation, MitM, resource
// exhaustion, tarpitting, etc.
//
// ragep2p uses the Ed25519 signature algorithm and sha2 hash function.
//
// ragep2p attempts to prevent fingerprinting of ragep2p nodes. ragep2p will not
// respond on a connection until it has received a valid knock message
// constructed over the PeerID of the connection initiator and the PeerID of the
// connection receiver. (knocks carry a signature for authentication, though
// it's important to note that by their uni-directional nature a knock does not
// constitute a proper handshake and can be replayed.)
//
// ragep2p connections are authenticated and encrypted using mutual TLS 1.3,
// using the crypto/tls package from Go's standard library. TLS is used with
// ephemeral certificates using the keypair corresponding to the Host's PeerID.
//
// ragep2p tries to defend against resource exchaustion attacks. In particular,
// we enforce maximum Stream counts per peer, maximum lengths for various
// messages, apply rate limiting at the tcp connection level as well as at the
// individual Stream level, and have a constant bound on the number of buffered
// messages per Stream.
//
// ragep2p defends against tarpitting, i.e. other peers that intentionally
// read/write from the underlying connection slowly. Host.NewStream(),
// Stream.Close(), Stream.SendMessage(), and Stream.Receive() return immediately
// and do any potential resulting communication asynchronously in the
// background. Host.Close() terminates after at most a few seconds.
//
package ragep2p
