package knockingtls

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/sec"
	p2ptls "github.com/libp2p/go-libp2p-tls"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/loghelper"
	"github.com/smartcontractkit/chainlink/libocr/offchainreporting/types"
	"golang.org/x/crypto/ed25519"
)

const ID = "cl_knockingtls/1.0.0"

type KnockingTLSTransport struct {
	tls        *p2ptls.Transport 	allowlist  []peer.ID         	privateKey *p2pcrypto.Ed25519PrivateKey
	myId       peer.ID
	logger     types.Logger
}

var (
	errInvalidConnection = errors.New("invalid connection")
	errInvalidSignature  = errors.New("invalid signature")
)

func (c *KnockingTLSTransport) SecureInbound(ctx context.Context, insecure net.Conn) (sec.SecureConn, error) {
	signatureReceived := make([]byte, ed25519.SignatureSize)

	logger := loghelper.MakeLoggerWithContext(c.logger, types.LogFields{
		"remoteAddr":   insecure.RemoteAddr(),
		"localAddr":    insecure.LocalAddr(),
		"allowlistLen": len(c.allowlist),
	})

	n, err := insecure.Read(signatureReceived)
	if err != nil {
		insecure.Close()
		return nil, fmt.Errorf("can't read sig: %w", err)
	}

	if n < ed25519.SignatureSize {
												err = insecure.Close()
		return nil, fmt.Errorf("can't read sig: %w", err)
	}

	var wg sync.WaitGroup
	admissionChan := make(chan peer.ID)
		for _, peerId := range c.allowlist {
		wg.Add(1)
		go func(id peer.ID) {
			defer wg.Done()
			pk, err := id.ExtractPublicKey()
			if err != nil {
				return
			}
			verified, err := pk.Verify([]byte(c.myId.Pretty()), signatureReceived)
			if err != nil {
				return
			}
			if verified {
				admissionChan <- id
			}
		}(peerId)
	}

	go func() {
				wg.Wait()
		close(admissionChan)
	}()

		admittedId, ok := <-admissionChan

		if !ok {
		insecure.Close()
		return nil, errInvalidSignature
	} else {
		sconn, err := c.tls.SecureInbound(ctx, insecure)

		if err != nil {
			logger.Error("tls connection errored", types.LogFields{
				"err":           err,
				"sigVerifiedAs": admittedId,
			})
		}
		return sconn, err
	}
}

func (c *KnockingTLSTransport) SecureOutbound(ctx context.Context, insecure net.Conn, p peer.ID) (sec.SecureConn, error) {
	sig, err := c.privateKey.Sign([]byte(p.Pretty()))
	if err != nil {
		insecure.Close()
		return nil, err
	}

	if len(sig) != ed25519.SignatureSize {
		insecure.Close()
		return nil, errors.New("sign returned invalid sig")
	}

	n, err := insecure.Write(sig)
	if err != nil {
		insecure.Close()
		return nil, err
	}
	if n != len(sig) {
		insecure.Close()
		return nil, errors.New("can't send all tag")
	}

	return c.tls.SecureOutbound(ctx, insecure, p)
}

func (c *KnockingTLSTransport) UpdateAllowlist(allowlist []peer.ID) {
	c.logger.Debug("allowlist updated", types.LogFields{
		"old": c.allowlist,
		"new": allowlist,
	})
	c.allowlist = allowlist
}

func (c *KnockingTLSTransport) GetAllowlist() []peer.ID {
	return c.allowlist
}

func NewKnockingTLS(logger types.Logger, myPrivKey p2pcrypto.PrivKey, allowlist ...peer.ID) (*KnockingTLSTransport, error) {
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
		tls:        tls,
		allowlist:  allowlist,
		privateKey: ed25515Key,
		myId:       id,
		logger: loghelper.MakeLoggerWithContext(logger, types.LogFields{
			"id": "KnockingTLS",
		}),
	}, nil
}

var _ sec.SecureTransport = &KnockingTLSTransport{}
