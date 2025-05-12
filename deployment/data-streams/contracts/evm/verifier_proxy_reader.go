package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	verifier_proxy "github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_proxy_v0_5_0"
)

type VerifierInitializedIteratorAdapter struct {
	*verifier_proxy.VerifierProxyVerifierInitializedIterator
}

func (a *VerifierInitializedIteratorAdapter) GetEvent() *verifier_proxy.VerifierProxyVerifierInitialized {
	return a.Event
}

type VerifierSetIteratorAdapter struct {
	*verifier_proxy.VerifierProxyVerifierSetIterator
}

func (a *VerifierSetIteratorAdapter) GetEvent() *verifier_proxy.VerifierProxyVerifierSet {
	return a.Event
}

type VerifierUnsetIteratorAdapter struct {
	*verifier_proxy.VerifierProxyVerifierUnsetIterator
}

func (a *VerifierUnsetIteratorAdapter) GetEvent() *verifier_proxy.VerifierProxyVerifierUnset {
	return a.Event
}

type VerifierProxyReader struct {
	*verifier_proxy.VerifierProxyFilterer
	*verifier_proxy.VerifierProxyCaller
}

func (r *VerifierProxyReader) FilterVerifierInitialized(
	opts *bind.FilterOpts,
) (LogIterator[verifier_proxy.VerifierProxyVerifierInitialized], error) {
	iter, err := r.VerifierProxyFilterer.FilterVerifierInitialized(opts)
	if err != nil {
		return nil, err
	}

	return &VerifierInitializedIteratorAdapter{iter}, nil
}

func (r *VerifierProxyReader) FilterVerifierSet(
	opts *bind.FilterOpts,
) (LogIterator[verifier_proxy.VerifierProxyVerifierSet], error) {
	iter, err := r.VerifierProxyFilterer.FilterVerifierSet(opts)
	if err != nil {
		return nil, err
	}

	return &VerifierSetIteratorAdapter{iter}, nil
}

func (r *VerifierProxyReader) FilterVerifierUnset(
	opts *bind.FilterOpts,
) (LogIterator[verifier_proxy.VerifierProxyVerifierUnset], error) {
	iter, err := r.VerifierProxyFilterer.FilterVerifierUnset(opts)
	if err != nil {
		return nil, err
	}

	return &VerifierUnsetIteratorAdapter{iter}, nil
}

// NewVerifierProxyReader creates a new wrapper
func NewVerifierProxyReader(contract *verifier_proxy.VerifierProxy) *VerifierProxyReader {
	return &VerifierProxyReader{
		VerifierProxyFilterer: &contract.VerifierProxyFilterer,
		VerifierProxyCaller:   &contract.VerifierProxyCaller,
	}
}
