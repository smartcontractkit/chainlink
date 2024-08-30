package config

import (
	"encoding/json"
)

type HTTPServerConfig struct {
	Host                 string
	Port                 uint16
	TLSEnabled           bool
	TLSCertPath          string
	TLSKeyPath           string
	Path                 string
	ContentTypeHeader    string
	ReadTimeoutMillis    uint32
	WriteTimeoutMillis   uint32
	RequestTimeoutMillis uint32
	MaxRequestBytes      int64
}

type WebSocketServerConfig struct {
	HTTPServerConfig
	HandshakeTimeoutMillis uint32
}

type WebSocketClientConfig struct {
	HandshakeTimeoutMillis uint32
}

type GatewayConfig struct {
	UserServerConfig        HTTPServerConfig
	NodeServerConfig        WebSocketServerConfig
	ConnectionManagerConfig ConnectionManagerConfig
	Dons                    []DONConfig
}

type ConnectionManagerConfig struct {
	AuthGatewayId             string
	AuthTimestampToleranceSec uint32
	AuthChallengeLen          uint32
	HeartbeatIntervalSec      uint32
}

type DONConfig struct {
	DonId         string
	HandlerName   string
	HandlerConfig json.RawMessage
	Members       []NodeConfig
	F             int
}

type NodeConfig struct {
	Name    string
	Address string
}
