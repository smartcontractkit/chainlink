package vrf

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	vrf_types "github.com/smartcontractkit/ocr2vrf/types"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"

	"github.com/smartcontractkit/ocr2vrf/gethwrappers/testdkgstub"
	"github.com/smartcontractkit/ocr2vrf/internal/common/ocr"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
	"github.com/smartcontractkit/ocr2vrf/internal/crypto/point_translation"
	"github.com/smartcontractkit/ocr2vrf/internal/dkg/contract"
	"github.com/smartcontractkit/ocr2vrf/internal/util"
)

type testDKG interface {
	contract.DKG
	TransmitKeyToVRFContract() error
}

type stubDKG struct {
	pk           kyber.Point
	binaryPk     []byte
	sk           kyber.Scalar
	contract     *testdkgstub.TestDKGStub
	auth         *bind.TransactOpts
	client       util.CommittingClient
	address      common.Address
	keyID        contract.KeyID
	signers      []common.Address
	transmitters []common.Address
}

var _ contract.DKG = (*stubDKG)(nil)
var _ testDKG = (*stubDKG)(nil)

func newStubDKG(
	auth *bind.TransactOpts,
	client util.CommittingClient,
	keyID contract.KeyID,
	pkGroup kyber.Group,
	sk kyber.Scalar,
	pk kyber.Point,
	signers []common.Address,
	transmitters []common.Address,

) contract.OnchainContract {
	key, err := pk.MarshalBinary()
	if err != nil {
		panic(util.WrapError(err, "could not marshal key"))
	}
	address, tx, stubContract, err := testdkgstub.DeployTestDKGStub(
		auth,
		client,
		key,
		keyID,
	)
	if err != nil {
		panic(util.WrapError(err, "could not deploy stub DKG"))
	}
	client.Commit()
	if err := util.CheckStatus(context.TODO(), tx, client); err != nil {
		panic(util.WrapError(err, "deployment of stub DKG failed"))
	}
	return contract.OnchainContract{
		&stubDKG{pk,
			key,
			sk, stubContract,
			auth,
			client,
			address,
			keyID,
			signers,
			transmitters,
		},
		pkGroup,
	}
}

func (d *stubDKG) GetKey(
	_ context.Context,
	_ contract.KeyID,
	_ [32]byte,
) (contract.OnchainKeyData, error) {
	return contract.OnchainKeyData{PublicKey: d.binaryPk}, nil
}

func (d *stubDKG) AddClient(
	_ context.Context,
	_ [32]byte,
	clientAddress common.Address,
) error {
	tx, err := d.contract.AddClient(d.auth, d.keyID, clientAddress)
	if err != nil {
		return util.WrapErrorf(
			err,
			"could not add client 0x%x to stub DKG",
			clientAddress,
		)
	}
	d.client.Commit()
	if err := util.CheckStatus(context.TODO(), tx, d.client); err != nil {
		return util.WrapErrorf(
			err,
			"failed to add client 0x%x to stub DKG",
			clientAddress,
		)
	}
	return nil
}

func (d *stubDKG) Address() common.Address {
	return d.address
}

func (d *stubDKG) CurrentCommittee(
	ctx context.Context,
) (ocr.OCRCommittee, error) {
	return vrf_types.OCRCommittee{d.signers, d.transmitters}, nil
}

func (d *stubDKG) TransmitKeyToVRFContract() error {
	tx, err := d.contract.KeyGenerated(
		d.auth,
		testdkgstub.KeyDataStructKeyData{d.binaryPk, nil},
	)
	if err != nil {
		return util.WrapErrorf(
			err,
			"failed in call to KeyGenerated",
		)
	}
	d.client.Commit()
	if err := util.CheckStatus(context.TODO(), tx, d.client); err != nil {
		return util.WrapErrorf(
			err,
			"failed to generate key",
		)
	}
	return nil
}

func (d *stubDKG) InitiateDKG(
	_ context.Context,
	_ ocr.OCRCommittee,
	_ player_idx.Int,
	_ contract.KeyID,
	_ contract.EncryptionPublicKeys,
	_ contract.SigningPublicKeys,
	_ anon.Suite,
	_ point_translation.PubKeyTranslation,
) error {
	panic("implement me")
}
