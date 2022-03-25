package relay

import (
	"context"
	"encoding/json"

	solanaGo "github.com/gagliardetto/solana-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink-solana/pkg/solana"
	"github.com/smartcontractkit/chainlink-terra/pkg/terra"
	"go.uber.org/multierr"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/core/services/relay/types"
)

var (
	SupportedRelayers = map[types.Network]struct{}{
		types.EVM:    {},
		types.Solana: {},
		types.Terra:  {},
	}
	_ types.RelayerCtx = &evm.Relayer{}
	_ types.RelayerCtx = &solana.Relayer{}
	_ types.RelayerCtx = &terra.Relayer{}
)

type delegate struct {
	relayers map[types.Network]types.RelayerCtx
	ks       keystore.Master
}

// NewDelegate creates a master relay delegate which manages "relays" which are OCR2 median reporting plugins
// for various chains (evm and non-evm). nil Relayers will be disabled.
func NewDelegate(ks keystore.Master) *delegate {
	d := &delegate{
		ks:       ks,
		relayers: map[types.Network]types.RelayerCtx{},
	}
	return d
}

// AddRelayer registers the relayer r, or a disabled placeholder if nil.
// NOT THREAD SAFE
func (d delegate) AddRelayer(n types.Network, r types.RelayerCtx) {
	d.relayers[n] = r
}

// Start starts all relayers it manages.
func (d delegate) Start(ctx context.Context) error {
	var err error
	for _, r := range d.relayers {
		err = multierr.Combine(err, r.Start(ctx))
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

// OCR2ProviderArgs contains the minimal parameters to create a OCR2 Provider.
type OCR2ProviderArgs struct {
	ID              int32
	ContractID      string
	TransmitterID   null.String
	Relay           types.Network
	RelayConfig     job.JSONConfig
	IsBootstrapPeer bool
	Plugin          job.OCR2PluginType
}

// NewOCR2Provider creates a new OCR2 provider instance.
func (d delegate) NewOCR2Provider(externalJobID uuid.UUID, s interface{}) (types.OCR2ProviderCtx, error) {
	// We expect trusted input
	spec := s.(*OCR2ProviderArgs)
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
			Plugin:        spec.Plugin,
		})
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
			ChainID:            config.ChainID,
			ProgramID:          programID,
			StateID:            stateID,
			StoreProgramID:     storeProgramID,
			TransmissionsID:    transmissionsID,
			TransmissionSigner: transmissionSigner,
		})
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
		})
	default:
		return nil, errors.Errorf("unknown relayer network type: %s", spec.Relay)
	}
}
