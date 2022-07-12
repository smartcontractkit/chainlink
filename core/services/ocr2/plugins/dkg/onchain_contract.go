package dkg

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	dkgwrapper "github.com/smartcontractkit/ocr2vrf/gethwrappers/dkg"
	ocr2vrfTypes "github.com/smartcontractkit/ocr2vrf/types"
	"go.dedis.ch/kyber/v3/sign/anon"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
)

type onchainContract struct {
	wrapper      *dkgwrapper.DKG
	dkgAddress   common.Address
	vrfCommittte ocr2vrfTypes.OCRCommittee
}

var _ dkg.DKG = &onchainContract{}

func newOnchainDKGClient(dkgAddress string, ethClient evmclient.Client) (*onchainContract, error) {
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

func (o *onchainContract) AddClient(
	ctx context.Context,
	keyID [32]byte,
	clientAddress common.Address,
) error {
	// TODO: implement!!
	panic("unimplemented!")
}

func (o *onchainContract) Address() common.Address {
	return o.dkgAddress
}

func (o *onchainContract) CurrentCommittee(
	ctx context.Context,
) (ocr2vrfTypes.OCRCommittee, error) {
	return o.vrfCommittte, nil
}

func (o *onchainContract) InitiateDKG(
	ctx context.Context,
	committee ocr2vrfTypes.OCRCommittee,
	f uint8,
	keyID dkg.KeyID,
	epks dkg.EncryptionPublicKeys,
	spks dkg.SigningPublicKeys,
	encGroup anon.Suite,
	translator dkg.PubKeyTranslation,
) error {
	panic("unimplemented!")
}
