package wsrpc

import (
	"context"
	"crypto/ed25519"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/internal/message"
)

//go:generate mockery --name Logger --output ./mocks/ --case=underscore
type Logger interface {
	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
}

var _ Conn = &websocket.Conn{}

//go:generate mockery --name Conn --output ./mocks/ --case=underscore
type Conn interface {
	SetWriteDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, p []byte) (err error)
	Close() error
}

type UniClientConn struct {
	connMu sync.Mutex
	conn   Conn

	tlsConfig *tls.Config
	target    string
	connector func(ctx context.Context, target string, tlsConfig *tls.Config) (Conn, error)
	lggr      Logger
}

// DialUniWithContext will blocks until connection is established or context expires.
func DialUniWithContext(ctx context.Context, lggr Logger, target string, privKey ed25519.PrivateKey, serverPubKey ed25519.PublicKey) (*UniClientConn, error) {
	pubs := credentials.PublicKeys{serverPubKey}
	tlsConfig, err := credentials.NewClientTLSConfig(privKey, &pubs)
	if err != nil {
		return nil, err
	}
	conn, err := retryConnectWithBackoff(ctx, lggr, func(ctx2 context.Context) (Conn, error) {
		return connect(ctx2, target, tlsConfig)
	})
	if err != nil {
		return nil, err
	}
	return &UniClientConn{conn: conn, tlsConfig: tlsConfig, target: target, lggr: lggr, connector: connect}, nil
}

func connect(ctx context.Context, target string, tlsConfig *tls.Config) (Conn, error) {
	d := websocket.Dialer{
		TLSClientConfig:  tlsConfig,
		HandshakeTimeout: 45 * time.Second,
	}
	url := fmt.Sprintf("wss://%s", target)
	conn, _, err := d.DialContext(ctx, url, http.Header{})
	if err != nil {
		return nil, fmt.Errorf("error while dialing %w", err)
	}
	return conn, nil
}

func max(d1 time.Duration, d2 time.Duration) time.Duration {
	if d1 > d2 {
		return d1
	}
	return d2
}

// reconnect will retry forever to connect unless cancelled
// assumes caller holds conn lock.
func retryConnectWithBackoff(ctx context.Context, lggr Logger, connect func(ctx2 context.Context) (Conn, error)) (Conn, error) {
	reconnectWait := time.Second
	for {
		freshConn, err := connect(ctx)
		if err != nil {
			lggr.Warnf("error connecting %v, waiting then retrying", err)
			// If ctx is cancelled, return.
			// Otherwise, wait to reconnect and try again.
			select {
			case <-ctx.Done():
				lggr.Warnf("ctx error %v reconnecting", ctx.Err())
				return nil, ctx.Err()
			case <-time.After(reconnectWait):
			}
			reconnectWait = max(reconnectWait*2, 1*time.Minute)
			continue
		}
		return freshConn, nil
	}
}

// Invoke will try forever to send the message to the websocket unless context is cancelled.
// It reconnects on write/read errors to retry.
func (uc *UniClientConn) Invoke(ctx context.Context, method string, args interface{}, reply interface{}) error {
	uc.connMu.Lock()
	defer uc.connMu.Unlock()

	callID := uuid.NewString()
	req, err := message.NewRequest(callID, method, args)
	if err != nil {
		return err
	}
	reqBytes, err := MarshalProtoMessage(req)
	if err != nil {
		return err
	}
	var resBytes []byte
	deadline, isDeadline := ctx.Deadline()
	for {
		if isDeadline {
			_ = uc.conn.SetWriteDeadline(deadline)
		}
		err = uc.conn.WriteMessage(websocket.BinaryMessage, reqBytes)
		if err != nil {
			if ctx.Err() != nil {
				uc.lggr.Warnf("ctx error %v writing message", ctx.Err())
				return ctx.Err()
			}
			uc.lggr.Warnf("received error %v writing message, reconnecting", err)
			freshConn, err2 := retryConnectWithBackoff(ctx, uc.lggr, func(ctx2 context.Context) (Conn, error) {
				return uc.connector(ctx2, uc.target, uc.tlsConfig)
			})
			if err2 != nil {
				return err2
			}
			uc.conn = freshConn
			continue
		}
		if isDeadline {
			_ = uc.conn.SetReadDeadline(deadline)
		}
		_, resBytes, err = uc.conn.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				uc.lggr.Warnf("ctx error %v reading message", ctx.Err())
				return ctx.Err()
			}
			uc.lggr.Warnf("received error %v reading message, reconnecting", err)
			freshConn, err2 := retryConnectWithBackoff(ctx, uc.lggr, func(ctx2 context.Context) (Conn, error) {
				return uc.connector(ctx2, uc.target, uc.tlsConfig)
			})
			if err2 != nil {
				return err2
			}
			uc.conn = freshConn
			continue
		}
		break
	}
	msg := message.Message{}
	if err := UnmarshalProtoMessage(resBytes, &msg); err != nil {
		return err
	}
	switch r := msg.Exchange.(type) {
	case *message.Message_Response:
		if r == nil || r.Response.Payload == nil {
			return errors.New("response payload is nil")
		}
		return UnmarshalProtoMessage(r.Response.Payload, reply)
	default:
		return errors.New("unexpected message type")
	}
}

func (uc *UniClientConn) Close() error {
	uc.connMu.Lock()
	defer uc.connMu.Unlock()
	return uc.conn.Close()
}
