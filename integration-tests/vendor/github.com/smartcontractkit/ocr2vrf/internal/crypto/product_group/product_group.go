package product_group

import (
	"bytes"
	"crypto/cipher"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/sign/anon"
)

type ProductGroup struct {
	factorGroups []kyber.Group

	generators []kyber.Point

	r cipher.Stream
}

var _ anon.Suite = (*ProductGroup)(nil)

func NewProductGroup(factorGroups []kyber.Group, generators []kyber.Point, r cipher.Stream) (*ProductGroup, error) {
	if len(factorGroups) < 1 {
		return nil, errors.Errorf("need at least one slot in product group")
	}
	if len(factorGroups) != len(generators) {
		return nil, errors.Errorf("must be exactly one generator for each factor group")
	}
	scalarType := reflect.TypeOf(factorGroups[0].Scalar())
	for groupIdx, g := range factorGroups {
		if reflect.TypeOf(g.Scalar()) != scalarType {
			return nil, errors.Errorf("%dth group uses a different scalar type %T (expected %T)",
				groupIdx, reflect.TypeOf(g.Scalar()), scalarType)
		}
		if reflect.TypeOf(g.Point()) != reflect.TypeOf(generators[groupIdx]) {
			return nil, errors.Errorf("%dth generator has wrong type (expected %T, got %T)",
				groupIdx, g.Point(), generators[groupIdx])
		}
		if generators[groupIdx].Equal(g.Point().Null()) {
			return nil, errors.Errorf("%dth generator is zero point", groupIdx)
		}

	}
	return &ProductGroup{factorGroups, generators, r}, nil
}

func (pg *ProductGroup) String() string {
	var substrings, gsubstrings []string
	for factorIdx, g := range pg.factorGroups {
		substrings = append(substrings, pg.generators[factorIdx].String())
		gsubstrings = append(gsubstrings, g.String())
	}

	return fmt.Sprintf("ProductGroup〈(%s)〉 ⊂ (%s)", strings.Join(substrings, ","),
		strings.Join(gsubstrings, "×"))
}

func (pg *ProductGroup) ScalarLen() int {
	return pg.factorGroups[0].ScalarLen()
}

func (pg *ProductGroup) Scalar() kyber.Scalar {
	return pg.factorGroups[0].Scalar()
}

func (pg *ProductGroup) PointLen() (rv int) {
	for _, g := range pg.factorGroups {
		rv += g.PointLen()
	}
	return
}

func (pg *ProductGroup) Point() kyber.Point {
	rv := &PointTuple{nil, pg}
	for _, g := range pg.factorGroups {
		rv.elements = append(rv.elements, g.Point())
	}
	return rv
}

type PointTuple struct {
	elements []kyber.Point
	group    *ProductGroup
}

var _ kyber.Point = (*PointTuple)(nil)

func (pg *ProductGroup) NewPoint(elements []kyber.Point) (*PointTuple, error) {
	if len(elements) != len(pg.generators) {
		return nil, errors.Errorf("needed %d-tuple, got %d-tuple", len(pg.generators), len(elements))
	}
	for pidx, p := range elements {
		exemplar := pg.generators[pidx]
		if reflect.TypeOf(p) != reflect.TypeOf(exemplar) {
			return nil, errors.Errorf("needed point of type %T, got point of type %T", exemplar, p)
		}
	}
	return &PointTuple{elements, pg}, nil
}

func (pt *PointTuple) MarshalBinary() (data []byte, err error) {
	data, err = pt.WireMarshal()
	if err != nil {
		return nil, err
	}
	data = append(data, pt.group.String()...)
	return data, nil
}

func (pt *PointTuple) WireMarshal() (data []byte, err error) {
	for i, p := range pt.elements {
		d, err := p.MarshalBinary()
		if err != nil {
			return nil, errors.Wrapf(err, "from marshalling %dth element of %s", i, pt)
		}
		data = append(data, d...)
	}
	return data, nil
}

func (pt *PointTuple) UnmarshalBinary(data []byte) error {
	pt.elements = make([]kyber.Point, len(pt.group.factorGroups))
	cursor := 0
	for i, g := range pt.group.factorGroups {
		pLen := g.PointLen()
		pt.elements[i] = g.Point()
		if err := pt.elements[i].UnmarshalBinary(data[cursor : cursor+pLen]); err != nil {
			trunctatedPoint := &PointTuple{
				pt.elements[:i],
				&ProductGroup{pt.group.factorGroups[:i], pt.group.generators[:i], nil},
			}
			return errors.Wrapf(err, "while unmarshalling %dth element (unmarshaled so far: %s)", i, trunctatedPoint)
		}
		cursor += pLen
	}
	if !bytes.Equal(data[cursor:], []byte(pt.group.String())) {
		return errors.Errorf("domain separator %s missing, got %s", pt.group.String(), string(data[cursor:]))
	}
	return nil
}

func (pt *PointTuple) String() string {
	var substrings []string
	for _, e := range pt.elements {
		substrings = append(substrings, e.String())
	}
	return fmt.Sprintf("(%s)", strings.Join(substrings, ","))
}

func (pt *PointTuple) MarshalSize() (rv int)                          { panic("not implemented") }
func (pt *PointTuple) MarshalTo(io.Writer) (int, error)               { panic("not implemented") }
func (pt *PointTuple) UnmarshalFrom(io.Reader) (int, error)           { panic("not implemented") }
func (pt *PointTuple) Set(p kyber.Point) kyber.Point                  { panic("not implemented") }
func (pt *PointTuple) EmbedLen() int                                  { panic("unimplemented") }
func (pt *PointTuple) Embed(data []byte, r cipher.Stream) kyber.Point { panic("unimplemented") }
func (pt *PointTuple) Data() (rv []byte, err error)                   { panic("unimplemented") }

func (pt *PointTuple) Equal(s2 kyber.Point) (rv bool) {
	p2, ok := s2.(*PointTuple)
	if !ok || len(pt.elements) != len(p2.elements) {
		return false
	}
	rv = true
	for elementIdx, e := range pt.elements {
		rv = rv && e.Equal(p2.elements[elementIdx])
	}
	return rv
}

func (pt *PointTuple) Null() kyber.Point {
	for _, e := range pt.elements {
		_ = e.Null()
	}
	return pt
}

func (pt *PointTuple) Base() kyber.Point {
	for generatorIdx, g := range pt.group.generators {
		pt.elements[generatorIdx] = g.Clone()
	}
	return pt
}

func (pt *PointTuple) Pick(rand cipher.Stream) kyber.Point {
	return pt.Mul(pt.group.Scalar().Pick(rand), nil)
}

func (pt *PointTuple) Clone() kyber.Point {
	rv := pt.group.Point().(*PointTuple)
	for elementIdx, e := range pt.elements {
		rv.elements[elementIdx] = e.Clone()
	}
	return rv
}

func mustBeMatchingPointTuple(a, b kyber.Point) *PointTuple {
	b2, ok := b.(*PointTuple)
	if !ok || len(a.(*PointTuple).elements) != len(b2.elements) {
		panic(fmt.Sprintf(
			"attempt to set product group element %s to incompatible value %s, type %T",
			a, b, b))
	}
	return b2
}

func (pt *PointTuple) Add(a, b kyber.Point) kyber.Point {
	a2 := mustBeMatchingPointTuple(pt, a)
	b2 := mustBeMatchingPointTuple(pt, b)
	for elementIdx, ae := range a2.elements {
		pt.elements[elementIdx].Add(ae, b2.elements[elementIdx])
	}
	return pt
}

func (pt *PointTuple) Sub(a, b kyber.Point) kyber.Point {
	a2 := mustBeMatchingPointTuple(pt, a)
	b2 := mustBeMatchingPointTuple(pt, b)
	for elementIdx, ae := range a2.elements {
		pt.elements[elementIdx].Sub(ae, b2.elements[elementIdx])
	}
	return pt
}

func (pt *PointTuple) Neg(a kyber.Point) kyber.Point {
	a2 := mustBeMatchingPointTuple(pt, a)
	for elementIdx, e := range a2.elements {
		pt.elements[elementIdx].Neg(e)
	}
	return pt
}

func (pt *PointTuple) Mul(s kyber.Scalar, p kyber.Point) kyber.Point {
	if pt == nil {
		panic("attempt to modify nil pointer")
	}
	if p == nil {
		p = pt.Clone().Base()
	}
	p2 := mustBeMatchingPointTuple(pt, p)
	for elementIdx, e := range p2.elements {
		pt.elements[elementIdx].Mul(s, e)
	}
	return pt
}
