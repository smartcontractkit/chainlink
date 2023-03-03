package client

import (
	"errors"
	"fmt"
)

type ConnectionMode int

const (
	LocalConnection ConnectionMode = iota
	RemoteConnection
)

// Protocol represents a URL scheme to use when fetching connection details
type Protocol int

const (
	// WS : Web Socket Protocol
	WS Protocol = iota
	// WSSUFFIX : Web Socket Protocol
	WSSUFFIX
	// WSS : Web Socket Secure Protocol
	WSS
	// HTTP : Hypertext Transfer Protocol
	HTTP
	// HTTPS : Hypertext Transfer Protocol Secure
	HTTPS
)

// URLConverter converts ports to URLs
type URLConverter struct {
	ci  ConnectionInfo
	err error
}

// NewURLConverter creates new URLConverter instance
func NewURLConverter(fp ConnectionInfo, err error) *URLConverter {
	return &URLConverter{fp, err}
}

// As converts host/port to an URL
func (m *URLConverter) As(conn ConnectionMode, proto Protocol) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	var host string
	var port uint16
	if conn == RemoteConnection {
		host = m.ci.Host
		port = m.ci.Ports.Remote
	} else {
		host = "localhost"
		port = m.ci.Ports.Local
	}
	switch proto {
	case HTTP:
		return fmt.Sprintf("http://%s:%d", host, port), nil
	case HTTPS:
		return fmt.Sprintf("https://%s:%d", host, port), nil
	case WS:
		return fmt.Sprintf("ws://%s:%d", host, port), nil
	case WSSUFFIX:
		return fmt.Sprintf("ws://%s:%d/ws", host, port), nil
	case WSS:
		return fmt.Sprintf("wss://%s:%d", host, port), nil
	default:
		return "", errors.New("unknown protocol conversion type")
	}
}
