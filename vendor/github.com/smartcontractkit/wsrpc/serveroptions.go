package wsrpc

import (
	"crypto/ed25519"

	"github.com/smartcontractkit/wsrpc/credentials"
)

// A ServerOption sets options such as credentials, codec and keepalive parameters, etc.
type ServerOption interface {
	apply(*serverOptions)
}

type serverOptions struct {
	// Buffer
	writeBufferSize int
	readBufferSize  int

	// Transport Credentials
	creds credentials.TransportCredentials

	// The address that the healthcheck will run on
	healthcheckAddr string
}

// funcServerOption wraps a function that modifies serverOptions into an
// implementation of the ServerOption interface.
type funcServerOption struct {
	f func(*serverOptions)
}

func newFuncServerOption(f func(*serverOptions)) *funcServerOption {
	return &funcServerOption{
		f: f,
	}
}

func (fdo *funcServerOption) apply(do *serverOptions) {
	fdo.f(do)
}

// Creds returns a ServerOption that sets credentials for server connections.
func Creds(privKey ed25519.PrivateKey, pubKeys []ed25519.PublicKey) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		pubs := credentials.PublicKeys(pubKeys)

		config, err := credentials.NewServerTLSConfig(privKey, &pubs)
		if err != nil {
			return
		}

		o.creds = credentials.NewTLS(config, &pubs)
	})
}

// WriteBufferSize specifies the I/O write buffer size in bytes. If a buffer
// size is zero, then a useful default size is used.
func WriteBufferSize(s int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.writeBufferSize = s
	})
}

// WriteBufferSize specifies the I/O read buffer size in bytes. If a buffer
// size is zero, then a useful default size is used.
func ReadBufferSize(s int) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.readBufferSize = s
	})
}

var defaultServerOptions = serverOptions{
	writeBufferSize: 4096,
	readBufferSize:  4096,
}

// WithHealthcheck specifies whether to run a healthcheck endpoint. If a url
// is not provided, a healthcheck endpoint is not started.
func WithHealthcheck(addr string) ServerOption {
	return newFuncServerOption(func(o *serverOptions) {
		o.healthcheckAddr = addr
	})
}
