package evm

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/verifier_v0_5_0"
)

type ConfigSetIteratorAdapter struct {
	*verifier_v0_5_0.VerifierConfigSetIterator
}

func (a *ConfigSetIteratorAdapter) GetEvent() *verifier_v0_5_0.VerifierConfigSet {
	return a.Event
}

type ConfigUpdatedIteratorAdapter struct {
	*verifier_v0_5_0.VerifierConfigUpdatedIterator
}

func (a *ConfigUpdatedIteratorAdapter) GetEvent() *verifier_v0_5_0.VerifierConfigUpdated {
	return a.Event
}

type ConfigActivatedIteratorAdapter struct {
	*verifier_v0_5_0.VerifierConfigActivatedIterator
}

func (a *ConfigActivatedIteratorAdapter) GetEvent() *verifier_v0_5_0.VerifierConfigActivated {
	return a.Event
}

type ConfigDeactivatedIteratorAdapter struct {
	*verifier_v0_5_0.VerifierConfigDeactivatedIterator
}

func (a *ConfigDeactivatedIteratorAdapter) GetEvent() *verifier_v0_5_0.VerifierConfigDeactivated {
	return a.Event
}

// VerifierReader wraps the actual contract
type VerifierReader struct {
	*verifier_v0_5_0.VerifierFilterer
	*verifier_v0_5_0.VerifierCaller
}

// Override methods that need to return our custom interfaces
func (r *VerifierReader) FilterConfigSet(
	opts *bind.FilterOpts,
	configDigest [][32]byte,
) (LogIterator[verifier_v0_5_0.VerifierConfigSet], error) {
	iter, err := r.VerifierFilterer.FilterConfigSet(opts, configDigest)
	if err != nil {
		return nil, err
	}

	return &ConfigSetIteratorAdapter{iter}, nil
}

func (r *VerifierReader) FilterConfigUpdated(
	opts *bind.FilterOpts,
	configDigest [][32]byte,
) (LogIterator[verifier_v0_5_0.VerifierConfigUpdated], error) {
	iter, err := r.VerifierFilterer.FilterConfigUpdated(opts, configDigest)
	if err != nil {
		return nil, err
	}

	return &ConfigUpdatedIteratorAdapter{iter}, nil
}

func (r *VerifierReader) FilterConfigActivated(
	opts *bind.FilterOpts,
	configDigest [][32]byte,
) (LogIterator[verifier_v0_5_0.VerifierConfigActivated], error) {
	iter, err := r.VerifierFilterer.FilterConfigActivated(opts, configDigest)
	if err != nil {
		return nil, err
	}

	return &ConfigActivatedIteratorAdapter{iter}, nil
}

func (r *VerifierReader) FilterConfigDeactivated(
	opts *bind.FilterOpts,
	configDigest [][32]byte,
) (LogIterator[verifier_v0_5_0.VerifierConfigDeactivated], error) {
	iter, err := r.VerifierFilterer.FilterConfigDeactivated(opts, configDigest)
	if err != nil {
		return nil, err
	}

	return &ConfigDeactivatedIteratorAdapter{iter}, nil
}

// NewVerifierReader creates a new wrapper
func NewVerifierReader(contract *verifier_v0_5_0.Verifier) *VerifierReader {
	return &VerifierReader{
		VerifierFilterer: &contract.VerifierFilterer,
		VerifierCaller:   &contract.VerifierCaller,
	}
}
