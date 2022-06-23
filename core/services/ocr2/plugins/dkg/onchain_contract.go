package dkg

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/ocr2vrf/dkg"
	dkgwrapper "github.com/smartcontractkit/ocr2vrf/gethwrappers/dkg"
	dkgcontract "github.com/smartcontractkit/ocr2vrf/pkg/dkg/contract"

	evmclient "github.com/smartcontractkit/chainlink/core/chains/evm/client"
)

type onchainContract struct {
	wrapper    *dkgwrapper.DKG
	dkgAddress common.Address
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
	keyID dkgcontract.KeyID,
	configDigest [32]byte,
) (dkgcontract.OnchainKeyData, error) {
	keyData, err := o.wrapper.GetKey(&bind.CallOpts{
		Context: ctx,
	}, keyID, configDigest)
	if err != nil {
		return dkgcontract.OnchainKeyData{}, errors.Wrap(err, "wrapper GetKey")
	}
	return dkgcontract.OnchainKeyData{
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
