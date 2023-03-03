package dkg

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/protobuf"
	"github.com/smartcontractkit/ocr2vrf/internal/util"
	dkg_types "github.com/smartcontractkit/ocr2vrf/types"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type PluginConfig struct {
	offchainConfig offchainConfig
	onchainConfig  onchainConfig
}

type offchainConfig struct {
	epks []kyber.Point
	spks []kyber.Point

	encryptionGroup anon.Suite
	translator      point_translation.PubKeyTranslation
}

type onchainConfig struct{ contract.KeyID }

func (o *offchainConfig) MarshalBinary() ([]byte, error) {
	if len(o.epks) > int(player_idx.MaxPlayer) || len(o.epks) > math.MaxInt32 {
		return nil, fmt.Errorf("too many players")
	}
	if olen, elen := len(o.spks), len(o.epks); olen != elen {
		errMsg := "num public keys don't match; len(epks)=%d, len(spks)=%d"
		return nil, fmt.Errorf(errMsg, olen, elen)
	}
	epks, spks := make([][]byte, 0, len(o.epks)), make([][]byte, 0, len(o.epks))
	for i, epk := range o.epks {
		epkb, err := epk.MarshalBinary()
		if err != nil {
			errMsg := "could not marshal %dth encryption key %v"
			return nil, util.WrapErrorf(err, errMsg, i, epk)
		}
		epks = append(epks, epkb)
		spkb, err := o.spks[i].MarshalBinary()
		if err != nil {
			errMsg := "could not marshal %dth signing key %v"
			return nil, util.WrapErrorf(err, errMsg, i, o.spks[i])
		}
		spks = append(spks, spkb)
	}
	if o.encryptionGroup == nil {
		errMsg := "attempt to marshal DKG offchain config with nil encryption group"
		return nil, fmt.Errorf(errMsg)
	}
	if o.translator == nil {
		errMsg := "attempt to marshal DKG offchain config with nil group translator"
		return nil, fmt.Errorf(errMsg)
	}

	return proto.Marshal(&protobuf.OffchainConfig{
		EncryptionPKs:   epks,
		SignaturePKs:    spks,
		EncryptionGroup: o.encryptionGroup.String(),
		Translator:      o.translator.Name(),
	})

}

func (o *onchainConfig) Marshal() []byte {
	return append([]byte{}, o.KeyID[:]...)
}

func unmarshalBinaryOffchainConfig(
	offchainBinaryConfig []byte,
) (*offchainConfig, error) {
	p := &protobuf.OffchainConfig{}
	if err := proto.Unmarshal(offchainBinaryConfig, p); err != nil {
		errMsg := "could not deserialize dkg offchain config binary 0x%x"
		return nil, util.WrapErrorf(err, errMsg, offchainBinaryConfig)
	}
	if p.EncryptionPKs == nil || p.SignaturePKs == nil ||
		p.EncryptionGroup == "" || p.Translator == "" {
		errMsg := "empty field while deserializing offchain config: %+v"
		return nil, fmt.Errorf(errMsg, p)
	}
	nepk, nspk := len(p.EncryptionPKs), len(p.SignaturePKs)
	if nepk != nspk {
		errMsg := "num encryption PKs (%d) and signature PKs (%d) should match"
		return nil, fmt.Errorf(errMsg, nepk, nspk)
	}
	encGgroup, ok := encryptionGroupRegistry[p.EncryptionGroup]
	if !ok {
		errMsg := "unrecognized encryption-group name, '%s'"
		return nil, fmt.Errorf(errMsg, p.EncryptionGroup)
	}
	epks, spks := make([]kyber.Point, 0, nepk), make([]kyber.Point, 0, nspk)
	for i, bepk := range p.EncryptionPKs {
		epk := encGgroup.Point()
		if err := epk.UnmarshalBinary(bepk); err != nil {
			errMsg := "could not unmarshal encryption key 0x%x in group %s"
			return nil, util.WrapErrorf(err, errMsg, bepk, p.EncryptionGroup)
		}
		epks = append(epks, epk)
		spk := SigningGroup.Point()
		bspk := p.SignaturePKs[i]
		if err := spk.UnmarshalBinary(bspk); err != nil {
			errMsg := "could not unmarshal signing key 0x%x in group %s"
			return nil, util.WrapErrorf(err, errMsg, bspk, "ed25519")
		}
		spks = append(spks, spk)
	}
	translator, ok := translatorRegistry[p.Translator]
	if !ok {
		return nil, fmt.Errorf("unrecognized translator name: %s", p.Translator)
	}
	return &offchainConfig{
		epks,
		spks,
		encGgroup,
		translator,
	}, nil
}

func unmarshalBinaryOnchainConfig(onchainBinaryConfig []byte) (rv onchainConfig, err error) {
	if len(onchainBinaryConfig) != len(contract.KeyID{}) {
		return rv, fmt.Errorf("onchainConfig binary is wrong length")
	}
	copy(rv.KeyID[:], onchainBinaryConfig)
	return rv, nil
}

func unmarshalPluginConfig(offchainBinaryConfig, onchainBinaryConfig []byte) (*PluginConfig, error) {
	offchainConfig, err := unmarshalBinaryOffchainConfig(offchainBinaryConfig)
	if err != nil {
		errMsg := "while unmarshaling offchaincomponent of config"
		return nil, util.WrapError(err, errMsg)
	}
	hydratedOnchainConfig, err := unmarshalBinaryOnchainConfig(onchainBinaryConfig)
	if err != nil {
		errMsg := "while unmarshaling onchaincomponent of config"
		return nil, util.WrapError(err, errMsg)
	}
	return &PluginConfig{*offchainConfig, hydratedOnchainConfig}, nil
}

type NewDKGArgs struct {
	t                          player_idx.Int
	selfIdx                    *player_idx.PlayerIdx
	cfgDgst                    types.ConfigDigest
	esk                        kyber.Scalar
	epks                       []kyber.Point
	ssk                        kyber.Scalar
	spks                       []kyber.Point
	keyID                      contract.KeyID
	keyConsumer                KeyConsumer
	encryptionGroup            anon.Suite
	translationGroup           kyber.Group
	translator                 point_translation.PubKeyTranslation
	contract                   contract.OnchainContract
	logger                     commontypes.Logger
	randomness                 io.Reader
	db                         dkg_types.DKGSharePersistence
	xxxTestingOnlySigningGroup anon.Suite
}

func (p *PluginConfig) NewDKGArgs(
	d types.ConfigDigest,
	l *localArgs,
	oID commontypes.OracleID,
	n, t player_idx.Int,
) (*NewDKGArgs, error) {
	oc := p.offchainConfig
	translationGroup, err := oc.translator.TargetGroup(oc.encryptionGroup)
	if err != nil {
		errMsg := "could not determine translation target group"
		return nil, util.WrapError(err, errMsg)
	}
	players, err := player_idx.PlayerIdxs(n)
	if err != nil {
		return nil, util.WrapError(err, "could not determine local player index")
	}
	selfIdx := players[oID]
	return &NewDKGArgs{
		t,
		selfIdx,
		d,
		l.esk,
		oc.epks,
		l.ssk,
		oc.spks,
		p.onchainConfig.KeyID,
		l.keyConsumer,
		oc.encryptionGroup,
		translationGroup,
		oc.translator,
		l.contract,
		l.logger,
		l.randomness,
		l.shareDB,
		nil,
	}, nil
}

type localArgs struct {
	esk         kyber.Scalar
	ssk         kyber.Scalar
	keyID       contract.KeyID
	contract    contract.OnchainContract
	logger      commontypes.Logger
	keyConsumer KeyConsumer
	randomness  io.Reader
	shareDB     dkg_types.DKGSharePersistence
}

func (o *offchainConfig) String() string {
	epks := make([]string, len(o.epks))
	spks := make([]string, len(o.spks))
	for i := range epks {
		epk, err := o.epks[i].MarshalBinary()
		if err != nil {
			epks[i] = "unmarshallable: " + err.Error()
		} else {
			epks[i] = hexutil.Encode(epk)
		}
		spk, err := o.spks[i].MarshalBinary()
		if err != nil {
			epks[i] = "unmarshallable: " + err.Error()
		} else {
			spks[i] = hexutil.Encode(spk)
		}
	}
	return fmt.Sprintf(`PluginConfig{
  epks: %s,
  spks: %s,
  encryptionGroup: %s,
  translator: %s,
}`,
		strings.Join(epks, ", "),
		strings.Join(spks, ", "),
		o.encryptionGroup,
		o.translator,
	)
}
