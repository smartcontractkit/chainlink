package relay

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/core/services/relay/types"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	"github.com/smartcontractkit/chainlink/core/services/keystore"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"go.uber.org/multierr"
)

var (
	SupportedRelayers = map[types.Network]struct{}{
		types.EVM:    {},
		types.Solana: {},
		types.Terra:  {},
	}
	_ types.Relayer = &evm.Relayer{}
	_ types.Relayer = &solana.Relayer{}
	_ types.Relayer = &terra.Relayer{}
)

type delegate struct {
	relayers map[types.Network]types.Relayer
	ks       keystore.Master
}

// NewDelegate creates a master relay delegate which manages "relays" which are OCR2 median reporting plugins
// for various chains (evm and non-evm). nil Relayers will be disabled.
func NewDelegate(ks keystore.Master) *delegate {
	d := &delegate{
		ks:       ks,
		relayers: map[types.Network]types.Relayer{},
	}
	return d
}

// AddRelayer registers the relayer r, or a disabled placeholder if nil.
// NOT THREAD SAFE
func (d delegate) AddRelayer(n types.Network, r types.Relayer) {
	d.relayers[n] = r
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

func (d delegate) NewOCR2Provider(externalJobID uuid.UUID, s interface{}, contractReady chan struct{}) (types.OCR2Provider, error) {
	// We expect trusted input
	spec := s.(*job.OffchainReporting2OracleSpec)
	choice := spec.Relay
	switch choice {
	case types.EVM:
		r, exists := d.relayers[types.EVM]
		if !exists {
			return nil, errors.New("no EVM relay found; is EVM enabled?")
		}

		var config evm.RelayConfig
		err := json.Unmarshal(spec.RelayConfig.Bytes(), &config)
		if err != nil {
			return nil, err
		}

		return r.NewOCR2Provider(externalJobID, evm.OCR2Spec{
			ID:            spec.ID,
			IsBootstrap:   spec.IsBootstrapPeer,
			ContractID:    spec.ContractID,
			TransmitterID: spec.TransmitterID,
			ChainID:       config.ChainID.ToInt(),
		}, contractReady)
	case types.Solana:
		r, exists := d.relayers[types.Solana]
		if !exists {
			return nil, errors.New("no Solana relay found; is Solana enabled?")
		}

		var config solana.RelayConfig
		err := json.Unmarshal(spec.RelayConfig.Bytes(), &config)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'spec.RelayConfig' unmarshal")
		}

		// use state account as contract ID (unique for each feed, program account is not)
		stateID, err := solanaGo.PublicKeyFromBase58(spec.ContractID)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'solana.PublicKeyFromBase58' for 'spec.ContractID")
		}

		programID, err := solanaGo.PublicKeyFromBase58(config.OCR2ProgramID)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.OCR2ProgramID")
		}

		storeProgramID, err := solanaGo.PublicKeyFromBase58(config.StoreProgramID)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.StateID")
		}

		transmissionsID, err := solanaGo.PublicKeyFromBase58(config.TransmissionsID)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'solana.PublicKeyFromBase58' for 'spec.RelayConfig.TransmissionsID")
		}

		var transmissionSigner solana.TransmissionSigner
		if !spec.IsBootstrapPeer {
			if !spec.TransmitterID.Valid {
				return nil, errors.New("transmitterID is required for non-bootstrap jobs")
			}
			transmissionSigner, err = d.ks.Solana().Get(spec.TransmitterID.String)
			if err != nil {
				return nil, err
			}
		}

		return r.NewOCR2Provider(externalJobID, solana.OCR2Spec{
			ID:                 spec.ID,
			IsBootstrap:        spec.IsBootstrapPeer,
			NodeEndpointHTTP:   config.NodeEndpointHTTP,
			ProgramID:          programID,
			StateID:            stateID,
			StoreProgramID:     storeProgramID,
			TransmissionsID:    transmissionsID,
			TransmissionSigner: transmissionSigner,
			UsePreflight:       config.UsePreflight,
			Commitment:         config.Commitment,
			PollingInterval:    config.PollingInterval,
			PollingCtxTimeout:  config.PollingCtxTimeout,
			StaleTimeout:       config.StaleTimeout,
		}, contractReady)
	case types.Terra:
		r, exists := d.relayers[types.Terra]
		if !exists {
			return nil, errors.New("no Terra relay found; is Terra enabled?")
		}

		var config terra.RelayConfig
		err := json.Unmarshal(spec.RelayConfig.Bytes(), &config)
		if err != nil {
			return nil, errors.Wrap(err, "error on 'spec.RelayConfig' unmarshal")
		}

		if !spec.IsBootstrapPeer {
			if !spec.TransmitterID.Valid {
				return nil, errors.New("transmitterID is required for non-bootstrap jobs")
			}
		}

		return r.NewOCR2Provider(externalJobID, terra.OCR2Spec{
			RelayConfig:   config,
			ID:            spec.ID,
			IsBootstrap:   spec.IsBootstrapPeer,
			ContractID:    spec.ContractID,
			TransmitterID: spec.TransmitterID.String,
		}, contractReady)
	default:
		return nil, errors.Errorf("unknown relayer network type: %s", spec.Relay)
	}
}
