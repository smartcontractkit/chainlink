package player_idx

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/commontypes"
)

var MarshalLen = int(math.Round(math.Log(float64(MaxPlayer))/math.Log(2))) / 8

func (pi *PlayerIdx) Marshal() []byte {
	return RawMarshal(pi.idx)
}

func RawMarshal(toMarshal Int) []byte {

	tm := uint64(toMarshal)
	rv := make([]byte, MarshalLen)
	for i := 0; tm != 0; i++ {
		rv[i] = uint8(tm & 0xFF)
		tm >>= 8
	}
	return rv
}

func Unmarshal(d []byte) (*PlayerIdx, []byte, error) {
	i, rem, err := RawUnmarshal(d)
	if err != nil {
		return nil, nil, err
	}
	if i == 0 {
		return nil, nil, errors.Errorf("player index must not be zero")
	}
	return &PlayerIdx{i}, rem, nil
}

func RawUnmarshal(d []byte) (Int, []byte, error) {
	if len(d) < MarshalLen {
		errMsg := "wrong length for marshalled player idx: expected %d, got %d"
		return 0, nil, errors.Errorf(errMsg, MarshalLen, len(d))
	}
	i := int64(0)
	for j := 0; j < MarshalLen; j++ {
		i <<= 8
		i += int64(d[j])
	}
	return Int(i), d[MarshalLen:], nil
}

func (pi PlayerIdx) Equal(pi2 *PlayerIdx) bool {
	return pi.idx == pi2.idx
}

func (pi PlayerIdx) NonZero() error {
	if pi.idx == 0 {
		return errors.Errorf("player index cannot be zero")
	}
	return nil
}

func (pi *PlayerIdx) Check() (*PlayerIdx, error) {
	if pi == nil {
		return nil, errors.Errorf("player index cannot be nil")
	}
	if err := pi.NonZero(); err != nil {
		return nil, err
	}
	return pi, nil
}

func (pi PlayerIdx) AtMost(n Int) bool {
	return pi.NonZero() == nil && pi.idx <= n
}

func (pi PlayerIdx) OracleID() commontypes.OracleID {
	return commontypes.OracleID(pi.idx - 1)
}

func (pi PlayerIdx) String() string {
	return fmt.Sprintf("player %d", pi.idx)
}
