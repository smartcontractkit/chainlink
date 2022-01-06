package relay

import (
	"encoding/json"

	"github.com/smartcontractkit/chainlink/core/services/relay/types"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	evmCh "github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/sqlx"

	"github.com/pkg/errors"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"go.uber.org/multierr"
)

var (
	SupportedRelayers = map[types.Network]struct{}{
		types.EVM:    {},
		types.Solana: {},
	}
	_ types.Relayer = &evm.Relayer{}
	_ types.Relayer = &solana.Relayer{}
)

type delegate struct {
	relayers map[types.Network]types.Relayer
	ks       keystore.Master
}

func NewDelegate(db *sqlx.DB, ks keystore.Master, chainSet evmCh.ChainSet, lggr logger.Logger) *delegate {
	return &delegate{
		ks: ks,
		relayers: map[types.Network]types.Relayer{
			types.EVM:    evm.NewRelayer(db, chainSet, lggr),
			types.Solana: solana.NewRelayer(lggr),
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

func (d delegate) NewOCR2Provider(externalJobID uuid.UUID, s interface{}) (types.OCR2Provider, error) {
	// We expect trusted input
	spec := s.(*job.OffchainReporting2OracleSpec)
	choice := spec.Relay
	switch choice {
	case types.EVM:
		var config evm.RelayConfig
		err := json.Unmarshal(spec.RelayConfig.Bytes(), &config)
		if err != nil {
			return nil, err
		}

		return d.relayers[types.EVM].NewOCR2Provider(externalJobID, evm.OCR2Spec{
			ID:            spec.ID,
			IsBootstrap:   spec.IsBootstrapPeer,
			ContractID:    spec.ContractID,
			TransmitterID: spec.TransmitterID,
			ChainID:       config.ChainID.ToInt(),
		})
	case types.Solana:
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

		return d.relayers[types.Solana].NewOCR2Provider(externalJobID, solana.OCR2Spec{
			ID:                 spec.ID,
			IsBootstrap:        spec.IsBootstrapPeer,
			NodeEndpointHTTP:   config.NodeEndpointHTTP,
			NodeEndpointWS:     config.NodeEndpointWS,
			ProgramID:          programID,
			StateID:            stateID,
			StoreProgramID:     storeProgramID,
			TransmissionsID:    transmissionsID,
			TransmissionSigner: transmissionSigner,
		})
	default:
		return nil, errors.Errorf("unknown relayer network type: %s", spec.Relay)
	}
}
