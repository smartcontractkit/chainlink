package change_base_group

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
	"reflect"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/ocr2vrf/internal/crypto/ciphertext/schnorr"

	"go.dedis.ch/kyber/v3"
)

func NewChangeBaseGroup(group schnorr.Suite, base kyber.Point) (*changeBaseGroup, error) {
	if !belongs(group, base) {
		return nil, errors.Errorf("point does not belong to group")
	}
	if base == nil || base.Equal(group.Point().Null()) {
		return nil, errors.Errorf("base point cannot be zero")
	}
	return &changeBaseGroup{group, base}, nil
}

type changeBaseGroup struct {
	schnorr.Suite
	base kyber.Point
}

var _ schnorr.Suite = (*changeBaseGroup)(nil)

func (g *changeBaseGroup) equal(og *changeBaseGroup) bool {
	return g.Suite == og.Suite && g.base.Equal(og.base)
}

func (g *changeBaseGroup) String() string {
	return fmt.Sprintf("&changeBaseGroup{Suite: %s; base: %s}", g.Suite, g.base)
}

type changeBasePoint struct {
	point kyber.Point
	group *changeBaseGroup
}

var _ kyber.Point = (*changeBasePoint)(nil)

func (g *changeBaseGroup) Point() kyber.Point {
	return &changeBasePoint{g.Suite.Point().Null(), g}
}

func (g *changeBaseGroup) Lift(p kyber.Point) (*changeBasePoint, error) {
	if !belongs(g.Suite, p) {
		return nil, errors.Errorf(
			"attempt to lift point when it doesn't belong to the underlying group (type %T but need %T)",
			p, g.Suite.Point(),
		)
	}
	return &changeBasePoint{p.Clone(), g}, nil
}

func (p *changeBasePoint) description() string {
	return fmt.Sprintf("〈%s〉⊂ (%s)", p.group.base, p.group.Suite)
}

func (p *changeBasePoint) MarshalBinary() ([]byte, error) {
	rv, err := p.point.MarshalBinary()
	if err != nil {
		return nil, errors.Wrapf(err, "could not marshal underlying point")
	}
	rv = append(rv, []byte(p.description())...)
	return rv, nil
}

func (p *changeBasePoint) UnmarshalBinary(data []byte) error {
	if len(data) < len(p.description()) {
		return errors.Errorf("marshal data cannot contain description")
	}
	dstart := len(data) - len(p.description())
	if !bytes.Equal(data[dstart:], []byte(p.description())) {
		return errors.Errorf("binary data does not contain group description")
	}
	if err := p.point.UnmarshalBinary(data[:dstart]); err != nil {
		return errors.Wrapf(err, "could not unmarshal underlying point")
	}
	return nil
}

func (p *changeBasePoint) String() string {
	return fmt.Sprintf("%s ∈ %s", p.point, p.description())
}

func (p *changeBasePoint) MarshalSize() int {
	return p.point.MarshalSize() + len(p.description())
}

func (p *changeBasePoint) MarshalTo(w io.Writer) (int, error) {
	n, err := p.point.MarshalTo(w)
	if err != nil {
		return n, errors.Wrapf(err, "could not marshal underlying point")
	}
	dn, err := w.Write([]byte(p.description()))
	return n + dn, errors.Wrapf(err, "could not write description for marshalling")
}

func (p *changeBasePoint) Equal(s2 kyber.Point) bool {
	s2cb, ok := s2.(*changeBasePoint)
	return ok && p.point.Equal(s2cb.point) && p.group.base.Equal(s2cb.group.base)
}

func (p *changeBasePoint) Mul(s kyber.Scalar, p2 kyber.Point) kyber.Point {
	if p2 == nil {
		p2 = p.group.base
	} else {
		if !p.group.equal(p2.(*changeBasePoint).group) {
			panic("attempt to multiply into different group")
		}
		p2 = p2.(*changeBasePoint).point
	}
	_ = p.point.Mul(s, p2)
	return p
}

func (p *changeBasePoint) Add(a kyber.Point, b kyber.Point) kyber.Point {
	if !p.group.equal(a.(*changeBasePoint).group) {
		panic("attempt to add into different group")
	}
	if !p.group.equal(b.(*changeBasePoint).group) {
		panic("attempt to add into different group")
	}
	_ = p.point.Add(a.(*changeBasePoint).point, b.(*changeBasePoint).point)
	return p
}

func (p *changeBasePoint) Null() kyber.Point {
	_ = p.point.Null()
	return p
}

func (p *changeBasePoint) Pick(rand cipher.Stream) kyber.Point {
	_ = p.point.Pick(rand)
	return p
}

func belongs(g kyber.Group, p kyber.Point) bool {
	return reflect.TypeOf(p) == reflect.TypeOf(g.Point())
}

func (p *changeBasePoint) Set(p2 kyber.Point) kyber.Point                 { panic("not implemented") }
func (p *changeBasePoint) EmbedLen() int                                  { panic("not implemented") }
func (p *changeBasePoint) Embed(data []byte, r cipher.Stream) kyber.Point { panic("not implemented") }
func (p *changeBasePoint) Data() ([]byte, error)                          { panic("not implemented") }
func (p *changeBasePoint) Base() kyber.Point                              { panic("not implemented") }
func (p *changeBasePoint) Clone() kyber.Point                             { panic("not implemented") }
func (p *changeBasePoint) UnmarshalFrom(r io.Reader) (int, error)         { panic("not implemented") }
func (p *changeBasePoint) Sub(a kyber.Point, b kyber.Point) kyber.Point   { panic("not implemented") }
func (p *changeBasePoint) Neg(a kyber.Point) kyber.Point                  { panic("not implemented") }
