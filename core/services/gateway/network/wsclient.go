package network

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

type WebSocketClient interface {
	Connect(ctx context.Context, url *url.URL) (*websocket.Conn, error)
}

type WebSocketClientConfig struct {
	HandshakeTimeoutMillis uint32
}

type webSocketClient struct {
	initiator ConnectionInitiator
	dialer    *websocket.Dialer
	lggr      logger.Logger
}

func NewWebSocketClient(config WebSocketClientConfig, initiator ConnectionInitiator, lggr logger.Logger) WebSocketClient {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Duration(config.HandshakeTimeoutMillis) * time.Millisecond,
	}
	client := &webSocketClient{
		initiator: initiator,
		dialer:    dialer,
		lggr:      logger.Named(lggr, "WebSocketClient"),
	}
	return client
}

func (c *webSocketClient) Connect(ctx context.Context, url *url.URL) (*websocket.Conn, error) {
	authHeader, err := c.initiator.NewAuthHeader(url)
	if err != nil {
		return nil, err
	}
	authHeaderStr := base64.StdEncoding.EncodeToString(authHeader)

	hdr := make(http.Header)
	hdr.Add(WsServerHandshakeAuthHeaderName, authHeaderStr)

	conn, resp, err := c.dialer.DialContext(ctx, url.String(), hdr)

	if err != nil {
		c.lggr.Errorf("WebSocketClient: couldn't connect to %s: %v", url.String(), err)
		c.tryCloseConn(conn)
		return nil, err
	}

	challengeStr := resp.Header.Get(WsServerHandshakeChallengeHeaderName)
	if challengeStr == "" {
		c.lggr.Error("WebSocketClient: empty challenge")
		c.tryCloseConn(conn)
		return nil, err
	}
	challenge, err := base64.StdEncoding.DecodeString(challengeStr)
	if err != nil {
		c.lggr.Errorf("WebSocketClient: couldn't decode challenge: %s: %v", challengeStr, err)
		c.tryCloseConn(conn)
		return nil, err
	}

	response, err := c.initiator.ChallengeResponse(url, challenge)
	if err != nil {
		c.lggr.Errorw("WebSocketClient: couldn't generate challenge response", "err", err)
		c.tryCloseConn(conn)
		return nil, err
	}

	if err = conn.WriteMessage(websocket.BinaryMessage, response); err != nil {
		c.lggr.Errorw("WebSocketClient: couldn't send challenge response", "err", err)
		c.tryCloseConn(conn)
		return nil, err
	}
	return conn, nil
}

func (c *webSocketClient) tryCloseConn(conn *websocket.Conn) {
	if conn != nil {
		err := conn.Close()
		if err != nil {
			c.lggr.Errorf("WebSocketClient: error closing connection %v", err)
		}
	}
}
