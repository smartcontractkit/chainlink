package ccipcapability

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type RelayGetter interface {
	GetIDToRelayerMap() (map[types.RelayID]loop.Relayer, error)
}

type Delegate struct {
	lggr            logger.Logger
	registrarConfig plugins.RegistrarConfig
	pipelineRunner  pipeline.Runner
	relayGetter     RelayGetter
	capRegistry     CapabilityRegistry
	keystore        keystore.Master
	ds              sqlutil.DataSource
	peerWrapper     *ocrcommon.SingletonPeerWrapper

	isNewlyCreatedJob bool
}

func NewDelegate(
	lggr logger.Logger,
	registrarConfig plugins.RegistrarConfig,
	pipelineRunner pipeline.Runner,
	relayGetter RelayGetter,
	registrySyncer CapabilityRegistry,
	keystore keystore.Master,
	ds sqlutil.DataSource,
	peerWrapper *ocrcommon.SingletonPeerWrapper,
) *Delegate {
	return &Delegate{
		lggr:            lggr,
		registrarConfig: registrarConfig,
		pipelineRunner:  pipelineRunner,
		relayGetter:     relayGetter,
		capRegistry:     registrySyncer,
		ds:              ds,
		keystore:        keystore,
		peerWrapper:     peerWrapper,
	}
}

func (d *Delegate) JobType() job.Type {
	return job.CCIP
}

func (d *Delegate) BeforeJobCreated(job.Job) {
	// This is only called first time the job is created
	d.isNewlyCreatedJob = true
}

func (d *Delegate) ServicesForSpec(ctx context.Context, spec job.Job) (services []job.ServiceCtx, err error) {
	// TODO: validate spec.

	// In general there should only be one P2P key but the node may have multiple.
	// The job spec should specify the correct P2P key to use.
	peerID, err := p2pkey.MakePeerID(spec.CCIPSpec.P2PKeyID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to make peer ID from provided spec p2p id: %s", spec.CCIPSpec.P2PKeyID)
	}

	p2pID, err := d.keystore.P2P().Get(peerID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all p2p keys")
	}

	ocrKeys := make(map[chaintype.ChainType]ocr2key.KeyBundle)
	for chainType, bundleAny := range spec.CCIPSpec.OCRKeyBundleIDs {
		ct := chaintype.ChainType(chainType)
		if !chaintype.IsSupportedChainType(ct) {
			return nil, errors.Errorf("unsupported chain type: %s", chainType)
		}

		bundleID, ok := bundleAny.(string)
		if !ok {
			return nil, errors.New("OCRKeyBundleIDs must be a map of chain types to OCR key bundle IDs")
		}

		bundle, err := d.keystore.OCR2().Get(bundleID)
		if err != nil {
			return nil, errors.Wrapf(err, "OCR key bundle with ID %s not found", bundleID)
		}

		ocrKeys[ct] = bundle
	}

	transmitterKeys := make(map[types.RelayID]string)
	for relayIDStr, transmitterIDAny := range spec.CCIPSpec.TransmitterIDs {
		var relayID types.RelayID
		if err := relayID.UnmarshalString(relayIDStr); err != nil {
			return nil, errors.Wrapf(err, "invalid relay ID specified in transmitter ids mapping: %s", relayIDStr)
		}

		transmitterID, ok := transmitterIDAny.(string)
		if !ok {
			return nil, errors.New("transmitter id is not a string")
		}

		switch relayID.Network {
		case types.NetworkEVM:
			ethKey, err := d.keystore.Eth().Get(ctx, transmitterID)
			if err != nil {
				return nil, errors.Wrapf(err, "eth transmitter key with ID %s not found", transmitterID)
			}

			transmitterKeys[relayID] = ethKey.String()
		case types.NetworkCosmos:
			cosmosKey, err := d.keystore.Cosmos().Get(transmitterID)
			if err != nil {
				return nil, errors.Wrapf(err, "cosmos transmitter key with ID %s not found", transmitterID)
			}

			transmitterKeys[relayID] = cosmosKey.String()
		case types.NetworkSolana:
			solKey, err := d.keystore.Solana().Get(transmitterID)
			if err != nil {
				return nil, errors.Wrapf(err, "solana transmitter key with ID %s not found", transmitterID)
			}

			transmitterKeys[relayID] = solKey.String()
		case types.NetworkStarkNet:
			starkKey, err := d.keystore.StarkNet().Get(transmitterID)
			if err != nil {
				return nil, errors.Wrapf(err, "starknet transmitter key with ID %s not found", transmitterID)
			}

			transmitterKeys[relayID] = starkKey.String()
		default:
			return nil, errors.Errorf("unsupported network: %s", relayID.Network)
		}
	}

	relayers, err := d.relayGetter.GetIDToRelayerMap()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all relayers")
	}

	// TODO: can we use the same DB for all plugin instances?
	// if all queries are scoped by config digest, should be OK.
	ocrDB := ocr2.NewDB(d.ds, spec.ID, 0, d.lggr)

	// TODO: pass in home chain reader
	hcr := &homeChainReader{}

	return []job.ServiceCtx{
		hcr,
		&launcher{
			ocrKeyBundles:          ocrKeys,
			transmitters:           transmitterKeys,
			relayers:               relayers,
			capRegistry:            d.capRegistry,
			p2pID:                  p2pID,
			peerWrapper:            d.peerWrapper,
			jobID:                  spec.ID,
			externalJobID:          spec.ExternalJobID,
			capabilityVersion:      spec.CCIPSpec.CapabilityVersion,
			capabilityLabelledName: spec.CCIPSpec.CapabilityLabelledName,
			lggr:                   d.lggr,
			homeChainReader:        hcr,
			isNewlyCreatedJob:      d.isNewlyCreatedJob,
			relayConfigs:           spec.CCIPSpec.RelayConfigs,
			pluginConfig:           spec.CCIPSpec.PluginConfig,
			db:                     ocrDB,
		}}, nil
}

func (d *Delegate) AfterJobCreated(spec job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(ctx context.Context, spec job.Job) error {
	// TODO: shut down needed services?
	return nil
}
