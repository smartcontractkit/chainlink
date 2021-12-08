package delegate

import (
	"encoding/json"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/sqlx"

	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/relay"
	evmrelay "github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"go.uber.org/multierr"
)

var (
	_ relay.Relayer = &evmrelay.Relayer{}
	_ relay.Relayer = &solana.Relayer{}
)

type delegate struct {
	relayers map[relay.Network]relay.Relayer
	ks       keystore.Master
}

func NewRelayDelegate(db *sqlx.DB, ks keystore.Master, chainSet evm.ChainSet, lggr logger.Logger) *delegate {
	return &delegate{
		ks: ks,
		relayers: map[relay.Network]relay.Relayer{
			relay.EVM:    evmrelay.NewRelayer(db, chainSet, lggr),
			relay.Solana: solana.NewRelayer(lggr),
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
	// We expect trusted input
	spec := s.(*job.OffchainReporting2OracleSpec)
	choice := spec.Relay
	switch choice {
	case relay.EVM:
		var config evmrelay.RelayConfig
		err := json.Unmarshal(spec.RelayConfig.Bytes(), &config)
		if err != nil {
			return nil, err
		}

		return d.relayers[relay.EVM].NewOCR2Provider(externalJobID, evmrelay.OCR2Spec{
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
		validatorProgramID, err := solanaGo.PublicKeyFromBase58(config.ValidatorProgramID)
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
		return d.relayers[relay.Solana].NewOCR2Provider(externalJobID, solana.OCR2Spec{
			ID:                 spec.ID,
			IsBootstrap:        spec.IsBootstrapPeer,
			NodeEndpointRPC:    config.NodeEndpointRPC,
			NodeEndpointWS:     config.NodeEndpointWS,
			ProgramID:          programID,
			StateID:            stateID,
			ValidatorProgramID: validatorProgramID,
			TransmissionsID:    transmissionsID,
			TransmissionSigner: transmissionSigner,
		})
	default:
		return nil, errors.Errorf("unknown relayer network type: %s", spec.Relay)
	}
}
