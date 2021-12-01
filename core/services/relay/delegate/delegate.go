package delegate

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	solanaGo "github.com/gagliardetto/solana-go"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	"github.com/smartcontractkit/chainlink/core/services/relay/ethereum"
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
	spec, ok := s.(*job.OffchainReporting2OracleSpec)
	if !ok {
		return nil, errors.New("unsuccessful cast to 'job.OffchainReporting2OracleSpec'")
	}

	choice := spec.Relay
	switch choice {
	case relay.Ethereum:
		var config ethereum.RelayConfig
		err := json.Unmarshal(spec.RelayConfig.Bytes(), &config)
		if err != nil {
			return nil, err
		}

		return d.relayers[choice].NewOCR2Provider(externalJobID, ethereum.OCR2Spec{
			ID:             spec.ID,
			IsBootstrap:    spec.IsBootstrapPeer,
			ContractID:     spec.ContractID,
			OCRKeyBundleID: spec.OCRKeyBundleID,
			TransmitterID:  spec.TransmitterID,
			ChainID:        config.ChainID.ToInt(),
		})
	case relay.Solana:
		var config solana.RelayConfig
		err := json.Unmarshal(spec.RelayConfig.Bytes(), &config)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'spec.RelayConfig' unmarshal")
		}

		programID, err := solanaGo.PublicKeyFromBase58(spec.ContractID.ValueOrZero())
		if err != nil {
			return nil, errors.Wrap(err, "error on 'solana.PublicKeyFromBase58' for 'spec.ContractID")
		}

		stateID, err := solanaGo.PublicKeyFromBase58(config.StateID)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.StateID")
		}

		transmissionsID, err := solanaGo.PublicKeyFromBase58(config.TransmissionsID)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.TransmissionsID")
		}

		return d.relayers[choice].NewOCR2Provider(externalJobID, solana.OCR2Spec{
			ID:              spec.ID,
			IsBootstrap:     spec.IsBootstrapPeer,
			NodeEndpointRPC: config.NodeEndpointRPC,
			NodeEndpointWS:  config.NodeEndpointWS,
			ProgramID:       programID,
			StateID:         stateID,
			TransmissionsID: transmissionsID,
			// Transmitter: TODO: get solana.PrivateKey from keystore using spec specified transmitterID
			KeyBundleID: spec.OCRKeyBundleID,
		})
	default:
		return nil, fmt.Errorf("unknown relayer network type: %s", spec.Relay)
	}
}
