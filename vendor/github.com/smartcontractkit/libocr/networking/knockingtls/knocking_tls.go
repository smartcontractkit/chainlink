package knockingtls

import (
	"context"
	"crypto"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"
	p2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/internal/loghelper"
	"golang.org/x/crypto/ed25519"
)

const ID = "cl_knockingtls/1.0.0"
const domainSeparator = "knockknock" + ID
const readTimeout = 1 * time.Minute
const version = byte(0x01)

// knock = version (1 byte) || pk (ed25519.PublicKeySize) || sig (ed25519.SignatureSize)
const knockSize = 1 + ed25519.PublicKeySize + ed25519.SignatureSize

type KnockingTLSTransport struct {
	tls               *p2ptls.Transport // underlying TLS transport
	allowlistMutex    sync.RWMutex      // allowlist may be accessed concurrently by libp2p (via SecureInbound) and OCR (via {Get|Update}Allowlist)
	allowlist         []peer.ID         // peer ids that are permitted
	privateKey        *p2pcrypto.Ed25519PrivateKey
	myId              peer.ID
	logger            loghelper.LoggerWithContext
	readTimeout       time.Duration
	bandwidthLimiters *Limiters
}

var _ sec.SecureTransport = (*KnockingTLSTransport)(nil)

var errInvalidSignature = errors.New("invalid signature in knock")

func buildKnockMessage(p peer.ID) ([]byte, error) {
	// defensive programming
	if len(p.Pretty()) > 128 {
		return nil, errors.New("too big id. looks suspicious")
	}
	h := crypto.SHA256.New()
	h.Write([]byte(domainSeparator))
	h.Write([]byte(p.Pretty()))

	return h.Sum(nil), nil
}

func (c *KnockingTLSTransport) SecureInbound(ctx context.Context, insecure net.Conn) (sec.SecureConn, error) {
	// always close conn unless the flag is set to false
	shouldClose := true
	defer func() {
		if shouldClose {
			insecure.Close()
		}
	}()

	knock := make([]byte, knockSize)

	logger := c.logger.MakeChild(commontypes.LogFields{
		"remoteAddr": insecure.RemoteAddr(),
		"localAddr":  insecure.LocalAddr(),
	})

	// set the read timeout so we don't block forever
	err := insecure.SetReadDeadline(time.Now().Add(c.readTimeout))
	if err != nil {
		return nil, err
	}
	n, err := insecure.Read(knock)
	if err != nil {
		return nil, fmt.Errorf("can't read knock: %w", err)
	}

	if n < knockSize {
		// We abort if the first read doesn't return knockSize bytes (which should be really rare).
		// Without this, an attacker can observe that we keep waiting until she sends the knockSize th byte, which
		// allows the attacker to fingerprint the server.
		return nil, fmt.Errorf("didn't get a full knock: got %d bytes", n)
	}

	if knock[0] != version {
		return nil, errors.New("invalid version")
	}
	// starting the 2nd byte is the actual knock
	knock = knock[1:]

	pk, err := p2pcrypto.UnmarshalEd25519PublicKey(knock[:ed25519.PublicKeySize])
	if err != nil {
		return nil, err
	}

	peerId, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return nil, err
	}

	inAllowList := false
	// wrap use of the mutex in a func so we can use defer to unlock
	func() {
		c.allowlistMutex.RLock()
		defer c.allowlistMutex.RUnlock()

		logger.Trace("verifying a knock", commontypes.LogFields{
			"allowlist": c.allowlist,
			"fromId":    peerId.Pretty(),
			"knock":     hex.EncodeToString(knock),
		})

		for i := range c.allowlist {
			if peerId == c.allowlist[i] {
				inAllowList = true
				break
			}
		}
	}()

	if !inAllowList {
		return nil, fmt.Errorf("remote peer %s not in the list", peerId.Pretty())
	}

	knockMsg, err := buildKnockMessage(c.myId)
	if err != nil {
		return nil, err
	}

	verified, err := pk.Verify(knockMsg, knock[ed25519.PublicKeySize:])
	if err != nil {
		return nil, err
	}

	if !verified {
		return nil, errInvalidSignature
	}

	// reset the timeout
	err = insecure.SetReadDeadline(time.Time{})
	if err != nil {
		return nil, err
	}

	// Wrap insecure connection with a bandwidth rate limiter.
	bandwidthLimiter, found := c.bandwidthLimiters.Find(peerId)
	if !found {
		c.logger.Error("Failed to find a rate limiter for inbound connection", commontypes.LogFields{
			"forPeerID":         peerId.Pretty(),
			"availableLimiters": c.bandwidthLimiters.Get(),
		})
		return nil, fmt.Errorf("Failed to find a rate limiter for peerID=%s in SecureInbound", peerId.Pretty())
	}

	insecureButLimited := NewRateLimitedConn(insecure, bandwidthLimiter, c.logger.MakeChild(commontypes.LogFields{
		"remotePeerID": peerId.Pretty(),
	}))

	secure, err := c.tls.SecureInbound(ctx, insecureButLimited)
	if err != nil {
		return nil, err
	}

	// enable rate limiting for the inbound connection. We only do this after
	// the TLS handshake completes to prevent a spoofing attacker from exhausting
	// an honest node's rate limit.
	insecureButLimited.EnableRateLimiting()

	// set flag to false so defer does not reset connection
	shouldClose = false
	return secure, nil
}

func (c *KnockingTLSTransport) SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (sec.SecureConn, error) {
	// always close conn unless the flag is set to false
	shouldClose := true
	defer func() {
		if shouldClose {
			insecure.Close()
		}
	}()

	pk, err := c.privateKey.GetPublic().Raw()
	if err != nil || len(pk) != ed25519.PublicKeySize {
		return nil, errors.New("can't get PK")
	}

	knockMsg, err := buildKnockMessage(p)
	if err != nil {
		return nil, err
	}
	sig, err := c.privateKey.Sign(knockMsg)
	if err != nil || len(sig) != ed25519.SignatureSize {
		return nil, errors.New("can't sign")
	}

	// knock = version || pk || sig
	knock := []byte{version}
	knock = append(knock, pk...)
	knock = append(knock, sig...)

	n, err := insecure.Write(knock)
	if err != nil {
		return nil, err
	}
	if n != knockSize {
		return nil, errors.New("can't send all tag")
	}

	// Wrap insecure connection with a bandwidth rate limiter.
	bandwidthLimiter, found := c.bandwidthLimiters.Find(p)
	if !found {
		c.logger.Error("Failed to find a rate limiter for outbound connection", commontypes.LogFields{
			"forPeerID":         p.Pretty(),
			"availableLimiters": c.bandwidthLimiters.Get(),
		})
		return nil, fmt.Errorf("Failed to find a rate limiter for peerID=%s in SecureOutbound", p.Pretty())
	}
	insecureButLimited := NewRateLimitedConn(insecure, bandwidthLimiter, c.logger.MakeChild(commontypes.LogFields{
		"remotePeerID": p.Pretty(),
	}))

	secure, err := c.tls.SecureOutbound(ctx, insecureButLimited, p)
	if err != nil {
		return nil, err
	}
	// enable rate limiting for the inbound connection. We only do this after
	// the TLS handshake completes to prevent a spoofing attacker from exhausting
	// an honest node's rate limit.
	insecureButLimited.EnableRateLimiting()

	// set the flag to false so defer doesn't close the conn
	shouldClose = false
	return secure, nil
}

func (c *KnockingTLSTransport) UpdateAllowlist(allowlist []peer.ID) {
	c.allowlistMutex.Lock()
	defer c.allowlistMutex.Unlock()

	c.logger.Debug("allowlist updated", commontypes.LogFields{
		"old": c.allowlist,
		"new": allowlist,
	})
	c.allowlist = allowlist
}

// NewKnockingTLS creates a TLS transport. Allowlist is a list of peer IDs that this transport should accept handshake from.
func NewKnockingTLS(logger commontypes.Logger, myPrivKey p2pcrypto.PrivKey, bandwidthLimiters *Limiters, allowlist ...peer.ID) (*KnockingTLSTransport, error) {
	ed25515Key, ok := myPrivKey.(*p2pcrypto.Ed25519PrivateKey)
	if !ok {
		return nil, errors.New("only support ed25519 key")
	}
	if allowlist == nil {
		allowlist = []peer.ID{}
	}

	tls, err := p2ptls.New(myPrivKey)
	if err != nil {
		return nil, err
	}

	id, err := peer.IDFromPrivateKey(myPrivKey)
	if err != nil {
		return nil, err
	}

	return &KnockingTLSTransport{
		tls:            tls,
		allowlistMutex: sync.RWMutex{},
		allowlist:      allowlist,
		privateKey:     ed25515Key,
		myId:           id,
		logger: loghelper.MakeRootLoggerWithContext(logger).MakeChild(commontypes.LogFields{
			"id": "KnockingTLS",
		}),
		readTimeout:       readTimeout,
		bandwidthLimiters: bandwidthLimiters,
	}, nil
}
