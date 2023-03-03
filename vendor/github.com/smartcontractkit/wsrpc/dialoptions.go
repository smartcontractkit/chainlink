package wsrpc

import (
	"crypto/ed25519"
	"log"
	"time"

	"github.com/smartcontractkit/wsrpc/credentials"
	"github.com/smartcontractkit/wsrpc/internal/backoff"
	"github.com/smartcontractkit/wsrpc/internal/transport"
)

// dialOptions configure a Dial call. dialOptions are set by the DialOption
// values passed to Dial.
type dialOptions struct {
	copts transport.ConnectOptions
	bs    backoff.Strategy
	block bool
}

// DialOption configures how we set up the connection.
type DialOption interface {
	apply(*dialOptions)
}

// funcDialOption wraps a function that modifies dialOptions into an
// implementation of the DialOption interface.
type funcDialOption struct {
	f func(*dialOptions)
}

func (fdo *funcDialOption) apply(do *dialOptions) {
	fdo.f(do)
}

func newFuncDialOption(f func(*dialOptions)) *funcDialOption {
	return &funcDialOption{
		f: f,
	}
}

// WithTransportCredentials returns a DialOption which configures a connection
// level security credentials (e.g., TLS/SSL).
func WithTransportCreds(privKey ed25519.PrivateKey, serverPubKey ed25519.PublicKey) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		pubs := credentials.PublicKeys{serverPubKey}

		// Generate the TLS config for the client
		config, err := credentials.NewClientTLSConfig(privKey, &pubs)
		if err != nil {
			log.Println(err)

			return
		}

		o.copts.TransportCredentials = credentials.NewTLS(config, &pubs)
	})
}

// WithBlock returns a DialOption which makes caller of Dial blocks until the
// underlying connection is up. Without this, Dial returns immediately and
// connecting the server happens in background.
func WithBlock() DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.block = true
	})
}

// WithWriteTimeout returns a DialOption which sets the write timeout for a
// message to be sent on the wire.
func WithWriteTimeout(d time.Duration) DialOption {
	return newFuncDialOption(func(o *dialOptions) {
		o.copts.WriteTimeout = d
	})
}

func defaultDialOptions() dialOptions {
	return dialOptions{
		copts: transport.ConnectOptions{},
	}
}
