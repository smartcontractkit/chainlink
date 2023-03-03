package gateway

import (
	"net/http"
)

type options struct {
	client       *http.Client
	chainID      string
	errorHandler func(e error) error
	baseUrl      string
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

func (fso *funcOption) apply(do *options) {
	fso.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// Option configures how we set up the connection.
type Option interface {
	apply(*options)
}

func WithHttpClient(client http.Client) Option {
	return newFuncOption(func(o *options) {
		o.client = &client
	})
}

func WithChain(chainID string) Option {
	return newFuncOption(func(o *options) {
		o.chainID = chainID
	})
}

func WithBaseURL(baseURL string) Option {
	return newFuncOption(func(o *options) {
		o.baseUrl = baseURL
	})
}

// WithErrorHandler returns an Option to set the error handler to be used.
func WithErrorHandler(f func(e error) error) Option {
	return newFuncOption(func(o *options) {
		o.errorHandler = f
	})
}
