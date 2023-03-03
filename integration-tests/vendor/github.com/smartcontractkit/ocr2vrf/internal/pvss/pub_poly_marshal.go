package pvss

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"
	kshare "go.dedis.ch/kyber/v3/share"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/player_idx"
)

type pubPoly struct{ *kshare.PubPoly }

func (p *pubPoly) marshal() ([]byte, error) {
	_, commits := p.Info()
	if len(commits) > int(player_idx.MaxPlayer) {
		return nil, errors.Errorf("too many coefficients to marshal")
	}
	rv := make([][]byte, 1+len(commits))
	cursor := 0

	rv[cursor] = player_idx.RawMarshal(player_idx.Int(len(commits)))
	cursor++

	pex, err := commits[0].MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "could not determine marshalled-point length")
	}
	pointLen := len(pex)

	for _, c := range commits {
		pb, err := c.MarshalBinary()
		if err != nil {
			return nil, errors.Wrap(err, "could not marshal point in coefficient commitments")
		}
		if len(pb) != pointLen {
			return nil, errors.Errorf("length of marshalled points varies")
		}
		rv[cursor] = pb
		cursor++
	}
	if cursor != len(rv) {
		panic(errors.Errorf("marshalled coefficient commitments fields out of registration"))
	}
	return bytes.Join(rv, nil), nil
}

func unmarshalPubPoly(g kyber.Group, data []byte) (commitments *pubPoly, rem []byte, err error) {

	numPoints, data, err := player_idx.RawUnmarshal(data)
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not parse number of coefficient commitments")
	}

	bytesRequired := int(numPoints) * g.PointLen()
	if bytesRequired > len(data) {
		return nil, nil, errors.Errorf(
			"data is too short to encode %d points (need %d bytes, got %d)",
			numPoints, bytesRequired, len(data))
	}
	commits := make([]kyber.Point, numPoints)
	for i := 0; i < int(numPoints); i++ {
		commits[i] = g.Point()
		if err := commits[i].UnmarshalBinary(data[:g.PointLen()]); err != nil {
			return nil, nil, errors.Wrap(err, "could not unmarshal coefficient commitment")
		}
		data = data[g.PointLen():]
	}
	return &pubPoly{kshare.NewPubPoly(g, g.Point().Base(), commits)}, data, nil
}

func (p *pubPoly) String() string {
	b, commits := p.PubPoly.Info()
	return fmt.Sprintf("&pubPoly{base: %s, commits: %s}", b, commits)
}
