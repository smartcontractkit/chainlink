package starknet

import (
	"fmt"
	"math/big"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/loop/internal/pb"
)

type Signature struct {
	sig *pb.StarknetSignature
}

func (s *Signature) Bytes() ([]byte, error) {
	return proto.Marshal(s.sig)
}

func (s *Signature) Ints() (x *big.Int, y *big.Int, err error) {
	if s.sig == nil {
		return nil, nil, fmt.Errorf("signature uninitialized")
	}

	return s.sig.X.Int(), s.sig.Y.Int(), nil
}

// b is expected to encode x,y components in accordance with [signature.Bytes]
func SignatureFromBytes(b []byte) (*Signature, error) {
	starkPb := &pb.StarknetSignature{}
	err := proto.Unmarshal(b, starkPb)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling bytes to signature: %w", err)
	}

	return &Signature{
		sig: starkPb,
	}, nil
}

// x,y must be non-negative numbers
func SignatureFromBigInts(x *big.Int, y *big.Int) (*Signature, error) {
	if x.Cmp(big.NewInt(0)) < 0 || y.Cmp(big.NewInt(0)) < 0 {
		return nil, fmt.Errorf("Cannot create signature from negative values (x,y), (%v, %v)", x, y)
	}

	starkPb := &pb.StarknetSignature{
		X: pb.NewBigIntFromInt(x),
		Y: pb.NewBigIntFromInt(y),
	}
	return &Signature{
		sig: starkPb,
	}, nil
}
