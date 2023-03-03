package contract

import (
	"context"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/ocr2vrf/internal/common/ocr"
	"github.com/smartcontractkit/ocr2vrf/internal/util"

	"go.dedis.ch/kyber/v3"
)

type KeyID [32]byte

type OnchainContract struct {
	DKG      DKG
	KeyGroup kyber.Group
}

type DKG interface {
	GetKey(
		ctx context.Context,
		keyID KeyID,
		configDigest [32]byte,
	) (OnchainKeyData, error)

	AddClient(
		ctx context.Context,
		keyID [32]byte,
		clientAddress common.Address,
	) error

	Address() common.Address

	CurrentCommittee(ctx context.Context) (ocr.OCRCommittee, error)
}

type OnchainKeyData struct {
	PublicKey []byte
	Hashes    [][32]byte
}

func (o OnchainContract) KeyData(
	ctx context.Context, keyID KeyID, cfgDgst types.ConfigDigest,
) (KeyData, error) {
	kd, err := o.DKG.GetKey(ctx, keyID, cfgDgst)
	if err != nil {
		return KeyData{}, util.WrapError(err, "could not retrieve key from contract")
	}
	if len(kd.PublicKey) == 0 || len(kd.Hashes) == 0 {
		return KeyData{}, nil
	}
	return MakeKeyDataFromOnchainKeyData(kd, o.KeyGroup)
}

func (o OnchainContract) Address() common.Address {
	return o.DKG.Address()
}
