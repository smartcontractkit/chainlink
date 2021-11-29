package delegate

import (
	"errors"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/chainlink/core/services/relay/ethereum"
	"github.com/smartcontractkit/chainlink/core/services/relay/solana"
	"go.uber.org/multierr"
)

type delegate struct {
	relayers relay.Relayers
}

func NewRelayer(config relay.Config) *delegate {
	return &delegate{
		relayers: relay.Relayers{
			relay.Ethereum: ethereum.NewRelayer(config),
			relay.Solana:   solana.NewRelayer(config),
		},
	}
}

// A delegate relayer on start will start all relayers it manages.
func (d delegate) Start() error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Start())
	}
	return err
}

// A delegate relayer on close will close all relayers it manages.
func (d delegate) Close() error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Close())
	}
	return err
}

// A delegate relayer is healthy if all relayers it manages are ready.
func (d delegate) Ready() error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Ready())
	}
	return err
}

// A delegate relayer is healthy if all relayers it manages are healthy.
func (d delegate) Healthy() error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Healthy())
	}
	return err
}

func (d delegate) NewOCR2Provider(externalJobID uuid.UUID, s interface{}) (relay.OCR2Provider, error) {
	// TODO [relay]: make a new specific OffchainReporting2RelayOracleSpec
	spec, ok := s.(job.OffchainReporting2OracleSpec)
	if !ok {
		return nil, errors.New("unsuccessful cast to 'job.OffchainReporting2OracleSpec'")
	}

	// TODO [relay]: make a relay network choice depending on job spec
	network := "solana"
	choice := relay.Network(network)

	switch choice {
	case relay.Ethereum:
		return d.relayers[choice].NewOCR2Provider(externalJobID, ethereum.OCR2Spec{
			ID:                      spec.ID,
			ChainID:                 spec.EVMChainID,
			ContractAddress:         spec.ContractAddress,
			EncryptedOCRKeyBundleID: spec.EncryptedOCRKeyBundleID,
			TransmitterAddress:      spec.TransmitterAddress,
		})
	case relay.Solana:
		return d.relayers[choice].NewOCR2Provider(externalJobID, solana.OCR2Spec{
			ID:                      spec.ID,
			ChainID:                 spec.EVMChainID,
			NodeEndpointRPC:         "", // TODO [relay]: add validator url from job spec
			NodeEndpointWS:          "", // TODO [relay]: add validator url from job spec
			ContractAddress:         "", // TODO [relay]: add contract address from job spec
			EncryptedOCRKeyBundleID: spec.EncryptedOCRKeyBundleID,
			TransmitterAddress:      "", // TODO [relay]: add transmitter address from job spec
		})
	default:
		return nil, fmt.Errorf("unknown relayer network type: %s", network)
	}
}
