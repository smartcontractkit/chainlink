package dkg

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3/sign/anon"

	"github.com/smartcontractkit/ocr2vrf/dkg"
	dkgwrapper "github.com/smartcontractkit/ocr2vrf/gethwrappers/dkg"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
)

type onchainContract struct {
	wrapper    *dkgwrapper.DKG
	dkgAddress common.Address
}

var _ dkg.DKG = &onchainContract{}

func NewOnchainDKGClient(dkgAddress string, ethClient evmclient.Client) (dkg.DKG, error) {
	dkgAddr := common.HexToAddress(dkgAddress)
	wrapper, err := dkgwrapper.NewDKG(dkgAddr, ethClient)
	if err != nil {
		return nil, errors.Wrap(err, "new dkg wrapper")
	}
	return &onchainContract{
		wrapper:    wrapper,
		dkgAddress: dkgAddr,
	}, nil
}

func (o *onchainContract) GetKey(
	ctx context.Context,
	keyID dkg.KeyID,
	configDigest [32]byte,
) (dkg.OnchainKeyData, error) {
	keyData, err := o.wrapper.GetKey(&bind.CallOpts{
		Context: ctx,
	}, keyID, configDigest)
	if err != nil {
		return dkg.OnchainKeyData{}, errors.Wrap(err, "wrapper GetKey")
	}
	return dkg.OnchainKeyData{
		PublicKey: keyData.PublicKey,
		Hashes:    keyData.Hashes,
	}, nil
}

func (o *onchainContract) Address() common.Address {
	return o.dkgAddress
}

func (o *onchainContract) CurrentCommittee(ctx context.Context) (ocr2vrftypes.OCRCommittee, error) {
	// NOTE: this is only ever used in tests in the ocr2vrf repo.
	// Since this isn't really used for production DKG running,
	// there's no point in implementing it.
	panic("unimplemented")
}

func (o *onchainContract) InitiateDKG(
	ctx context.Context,
	committee ocr2vrftypes.OCRCommittee,
	f ocr2vrftypes.PlayerIdxInt,
	keyID dkg.KeyID,
	epks dkg.EncryptionPublicKeys,
	spks dkg.SigningPublicKeys,
	encGroup anon.Suite,
	translator ocr2vrftypes.PubKeyTranslation,
) error {
	// NOTE: this is only ever used in tests, the idea here is to call setConfig
	// on the DKG contract to get the OCR process going. Since this isn't really
	// used for production DKG running, there's no point in implementing it.
	panic("unimplemented")
}

func (o *onchainContract) AddClient(
	ctx context.Context,
	keyID [32]byte,
	clientAddress common.Address,
) error {
	// NOTE: this is only ever used in tests in the ocr2vrf repo.
	// Since this isn't really used for production DKG running,
	// there's no point in implementing it.
	panic("unimplemented")
}
