// Package tdh2 implements the TDH2 protocol (Shoup and Gennaro, 2001: https://www.shoup.net/papers/thresh1.pdf).
package tdh2

import (
	"bytes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/smartcontractkit/tdh2/go/tdh2/internal/group"
	"github.com/smartcontractkit/tdh2/go/tdh2/internal/group/nist"
	"github.com/smartcontractkit/tdh2/go/tdh2/internal/group/share"
)

var (
	// defaultHash is the default hash function used. Note, its output size
	// determines the input size in TDH2.
	defaultHash = sha256.New
	// InputSize determines the size of messages and labels.
	InputSize = defaultHash().Size()
)

func parseGroup(group string) (group.Group, error) {
	switch group {
	case nist.NewP256().String():
		return nist.NewP256(), nil
	}
	return nil, fmt.Errorf("unsupported group: %q", group)
}

// PrivateShare is a node's private share. It extends group.s share.PriShare.
type PrivateShare struct {
	group group.Group
	index int
	v     group.Scalar
}

func (s PrivateShare) String() string {
	return fmt.Sprintf("grp:%s idx:%d", s.group.String(), s.index)
}

func (s PrivateShare) Index() int {
	return s.index
}

// mulPoint returns a new point with value v*p, where v is a private scalar.
// If p==nil, the returned point has value v*BasePoint.
func (s *PrivateShare) mulPoint(p group.Point) group.Point {
	return s.group.Point().Mul(s.v, p)
}

// mulScalar returns a new scalar with value v*a where v is a private scalar.
func (s *PrivateShare) mulScalar(a group.Scalar) group.Scalar {
	return s.group.Scalar().Mul(s.v, a)
}

func (p *PrivateShare) Clear() {
	p.group = nil
	p.index = 0
	p.v.Zero()
}

// privateShareRaw is used for PrivateShare (un)marshaling.
type privateShareRaw struct {
	Group string
	Index int
	V     []byte
}

func (s PrivateShare) Marshal() ([]byte, error) {
	v, err := s.v.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal V: %w", err)
	}
	return json.Marshal(&privateShareRaw{
		Group: s.group.String(),
		Index: s.index,
		V:     v,
	})
}

func (s *PrivateShare) Unmarshal(data []byte) error {
	var raw privateShareRaw
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("cannot unmarshal data: %w", err)
	}

	s.group, err = parseGroup(raw.Group)
	if err != nil {
		return fmt.Errorf("cannot parse group: %w", err)
	}

	s.index = raw.Index
	s.v = s.group.Scalar()
	if err = s.v.UnmarshalBinary(raw.V); err != nil {
		return fmt.Errorf("cannot unmarshal: %w", err)
	}
	return nil
}

// PubliKey defines a public and verification key.
type PublicKey struct {
	group  group.Group
	g_bar  group.Point
	h      group.Point
	hArray []group.Point
}

func (a *PublicKey) Equal(b *PublicKey) bool {
	if a.group.String() != b.group.String() || !a.g_bar.Equal(b.g_bar) || !a.h.Equal(b.h) {
		return false
	}
	if len(a.hArray) != len(b.hArray) {
		return false
	}
	for i := range a.hArray {
		if !a.hArray[i].Equal(b.hArray[i]) {
			return false
		}
	}
	return true
}

// MasterSecret keeps the master secret of a TDH2 instance.
type MasterSecret struct {
	group group.Group
	s     group.Scalar
}

func (m *MasterSecret) String() string {
	return fmt.Sprintf("group:%s value:hidden", m.group)
}

func (m *MasterSecret) Clear() {
	m.group = nil
	m.s.Zero()
}

// masterSecretRaw is used for MasterSecret (un)marshaling.
type masterSecretRaw struct {
	Group string
	S     []byte
}

func (m *MasterSecret) Marshal() ([]byte, error) {
	s, err := m.s.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal s: %w", err)
	}
	return json.Marshal(&masterSecretRaw{
		Group: m.group.String(),
		S:     s,
	})
}

func (m *MasterSecret) Unmarshal(data []byte) error {
	var raw masterSecretRaw
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot unmarshal data: %w", err)
	}
	var err error
	m.group, err = parseGroup(raw.Group)
	if err != nil {
		return fmt.Errorf("cannot parse group: %w", err)
	}
	m.s = m.group.Scalar()
	if err := m.s.UnmarshalBinary(raw.S); err != nil {
		return fmt.Errorf("cannot unmarshal s: %w", err)
	}
	return nil
}

// publicKeyRaw is used for PublicKey (un)marshaling.
type publicKeyRaw struct {
	Group  string
	G_bar  []byte
	H      []byte
	HArray [][]byte
}

func (p PublicKey) Marshal() ([]byte, error) {
	gbar, err := p.g_bar.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshaling G_bar: %w", err)
	}

	h, err := p.h.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshaling H: %w", err)
	}

	harray := [][]byte{}
	for _, h := range p.hArray {
		d, err := h.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("cannot marshal H: %w", err)
		}
		harray = append(harray, d)
	}

	return json.Marshal(&publicKeyRaw{
		Group:  p.group.String(),
		G_bar:  gbar,
		H:      h,
		HArray: harray,
	})
}

func (p *PublicKey) Unmarshal(data []byte) error {
	var raw publicKeyRaw
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("cannot unmarshal data: %w", err)
	}

	p.group, err = parseGroup(raw.Group)
	if err != nil {
		return fmt.Errorf("cannot parse group: %w", err)
	}

	p.g_bar = p.group.Point()
	if err = p.g_bar.UnmarshalBinary(raw.G_bar); err != nil {
		return fmt.Errorf("unmarshaling G_bar: %w", err)
	}

	p.h = p.group.Point()
	if err = p.h.UnmarshalBinary(raw.H); err != nil {
		return fmt.Errorf("unmarshaling H: %w", err)
	}

	p.hArray = []group.Point{}
	for _, h := range raw.HArray {
		new := p.group.Point()
		if err = new.UnmarshalBinary(h); err != nil {
			return fmt.Errorf("cannot unmarshal point: %w", err)
		}
		p.hArray = append(p.hArray, new)
	}

	return nil
}

// GenerateKeys generates a master secret, public key, and secret key shares according to TDH2 paper.
// It takes cryptographic group to be used, master secret to be used (if nil, a new secret is generated),
// the total number of nodes n, the number of shares sufficient for decryption k, and a randomness source.
// It returns the master secret (either passed or generated), public key, and secret key shares.
func GenerateKeys(grp group.Group, ms *MasterSecret, k, n int, rand cipher.Stream) (*MasterSecret, *PublicKey, []*PrivateShare, error) {
	if k > n {
		return nil, nil, nil, fmt.Errorf("threshold is higher than total number of nodes")
	}
	if k <= 0 {
		return nil, nil, nil, fmt.Errorf("threshold has to be positive")
	}
	if ms != nil && grp.String() != ms.group.String() {
		return nil, nil, nil, fmt.Errorf("inconsistent groups")
	}

	var s group.Scalar
	if ms != nil {
		s = ms.s
	}
	poly := share.NewPriPoly(grp, k, s, rand)
	x := poly.Secret()
	if ms != nil && !x.Equal(ms.s) {
		return nil, nil, nil, fmt.Errorf("generated wrong secret")
	}

	HArray := make([]group.Point, n)
	shares := poly.Shares(n)
	privShares := []*PrivateShare{}
	// IDs are assigned consecutively from 0.
	for i, s := range shares {
		if i != s.I {
			return nil, nil, nil, fmt.Errorf("share index=%d, expect=%d", s.I, i)
		}
		HArray[i] = grp.Point().Mul(s.V, nil)
		privShares = append(privShares, &PrivateShare{grp, s.I, s.V})
	}

	return &MasterSecret{
			group: grp,
			s:     x},
		&PublicKey{
			group:  grp,
			g_bar:  grp.Point().Pick(rand),
			h:      grp.Point().Mul(x, nil),
			hArray: HArray,
		}, privShares, nil
}

// Redeal re-deals private shares such that new quorums can decrypt old ciphertexts. It takes the
// previous public key and master secret as well as the number of nodes sufficient for decrypt k,
// the total number of nodes n, and a randomness source. It returns a new public key and private shares.
// The master secret passed corresponds to the public key returned. The old public key can still be used
// for encryption but it cannot be used for share verification (the new key has to be used instead).
func Redeal(pk *PublicKey, ms *MasterSecret, k, n int, rand cipher.Stream) (*PublicKey, []*PrivateShare, error) {
	if ms == nil {
		return nil, nil, fmt.Errorf("nil secret")
	}
	_, new, shares, err := GenerateKeys(pk.group, ms, k, n, rand)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot generate keys: %w", err)
	}
	return &PublicKey{
		group:  pk.group,
		g_bar:  pk.g_bar,
		h:      pk.h,
		hArray: new.hArray,
	}, shares, nil
}

// Encrypt a message with a label (see p15 of the paper).
func Encrypt(pk *PublicKey, msg []byte, label []byte, rand cipher.Stream) (*Ciphertext, error) {
	r := pk.group.Scalar().Pick(rand)
	s := pk.group.Scalar().Pick(rand)

	h, err := hash1(pk.group.String(), pk.group.Point().Mul(r, pk.h))
	if err != nil {
		return nil, fmt.Errorf("cannot hash: %w", err)
	}
	c, err := xor(h, msg)
	if err != nil {
		return nil, fmt.Errorf("cannot xor: %w", err)
	}

	u := pk.group.Point().Mul(r, nil)
	w := pk.group.Point().Mul(s, nil)
	u_bar := pk.group.Point().Mul(r, pk.g_bar)
	w_bar := pk.group.Point().Mul(s, pk.g_bar)
	e, err := hash2(c, label, u, w, u_bar, w_bar, pk.group)
	if err != nil {
		return nil, fmt.Errorf("cannot generate e: %w", err)
	}
	f := pk.group.Scalar().Add(s, pk.group.Scalar().Mul(r, e.Clone()))

	return &Ciphertext{
		group: pk.group,
		c:     c,
		label: label,
		u:     u,
		u_bar: u_bar,
		e:     e,
		f:     f,
	}, nil
}

// VerifyShare verifies the correctness of the decryption share obtained from node i.
// The caller has to ensure that the ciphertext is validated.
func VerifyShare(pk *PublicKey, ctxt *Ciphertext, share *DecryptionShare) error {
	if pk.group.String() != ctxt.group.String() {
		return fmt.Errorf("incorrect ciphertext group: %q", ctxt.group)
	}

	if pk.group.String() != share.group.String() {
		return fmt.Errorf("incorrect share group: %q", share.group)
	}

	if err := checkEi(pk, ctxt, share); err != nil {
		return fmt.Errorf("failed format validity check: %w", err)
	}

	return nil
}

// checkEi checks the validity of param E_i to ensure that it is a DH triple (formula 3 on p13).
func checkEi(pk *PublicKey, ctxt *Ciphertext, share *DecryptionShare) error {
	g := pk.group
	ui_hat := g.Point().Sub(g.Point().Mul(share.f_i, ctxt.u), g.Point().Mul(share.e_i, share.u_i))
	if share.index < 0 || share.index >= len(pk.hArray) {
		return fmt.Errorf("invalid share index")
	}
	hi_hat := g.Point().Sub(g.Point().Mul(share.f_i, nil), g.Point().Mul(share.e_i, pk.hArray[share.index]))
	ei, err := hash4(share.u_i, ui_hat, hi_hat, pk.group)
	if err != nil {
		return fmt.Errorf("cannot generate e_i: %w", err)
	}
	if !share.e_i.Equal(ei) {
		return fmt.Errorf("error during the verification of E_i")
	}
	return nil
}

// Ciphertext defines a ciphertext as output from the Encryption algorithm.
type Ciphertext struct {
	group group.Group
	c     []byte
	label []byte
	u     group.Point
	u_bar group.Point
	e     group.Scalar
	f     group.Scalar
}

// Verify checks if the ciphertext matches the public key
// (i.e., it checks the validity of param e -- see formula 4 on p15).
func (c *Ciphertext) Verify(pk *PublicKey) error {
	if c.group.String() != pk.group.String() {
		return fmt.Errorf("group mismatch")
	}
	w := pk.group.Point().Sub(pk.group.Point().Mul(c.f, nil), pk.group.Point().Mul(c.e, c.u))
	w_bar := pk.group.Point().Sub(pk.group.Point().Mul(c.f, pk.g_bar), pk.group.Point().Mul(c.e, c.u_bar))
	e, err := hash2(c.c, c.label, c.u, w, c.u_bar, w_bar, pk.group)
	if err != nil {
		return fmt.Errorf("cannot compute e: %w", err)
	}
	if !c.e.Equal(e) {
		return fmt.Errorf("wrong e")
	}
	return nil
}

func (a *Ciphertext) Equal(b *Ciphertext) bool {
	if a.group.String() != b.group.String() ||
		!bytes.Equal(a.c, b.c) ||
		!bytes.Equal(a.label, b.label) ||
		!a.u.Equal(b.u) ||
		!a.u_bar.Equal(b.u_bar) ||
		!a.e.Equal(b.e) ||
		!a.f.Equal(b.f) {
		return false
	}
	return true

}

// Decrypt decrypts a ciphertext using a secret key share x_i according to TDH2 paper.
// The caller has to ensure that the ciphertext is validated.
func (ctxt *Ciphertext) Decrypt(group group.Group, x_i *PrivateShare, rand cipher.Stream) (*DecryptionShare, error) {
	if group.String() != ctxt.group.String() {
		return nil, fmt.Errorf("incorrect ciphertext group: %q", ctxt.group)
	}
	if group.String() != x_i.group.String() {
		return nil, fmt.Errorf("incorrect share group: %q", x_i.group)
	}

	s_i := group.Scalar().Pick(rand)
	u_i := x_i.mulPoint(ctxt.u)
	u_hat := group.Point().Mul(s_i, ctxt.u)
	h_hat := group.Point().Mul(s_i, nil)
	e_i, err := hash4(u_i, u_hat, h_hat, group)
	if err != nil {
		return nil, fmt.Errorf("cannot generate e_i: %w", err)
	}
	f_i := group.Scalar().Add(s_i, x_i.mulScalar(e_i.Clone()))
	return &DecryptionShare{
		group: group,
		index: x_i.index,
		u_i:   u_i,
		e_i:   e_i,
		f_i:   f_i,
	}, nil
}

// CombineShares combines a set of decryption shares and returns the decrypted message.
// The caller has to ensure that the ciphertext is validated, the decryption shares are valid,
// all the shares are distinct and the number of them is at least k.
func (c *Ciphertext) CombineShares(group group.Group, shares []*DecryptionShare, k, n int) ([]byte, error) {
	if group.String() != c.group.String() {
		return nil, fmt.Errorf("incorrect ciphertext group: %q", c.group)
	}

	if len(shares) < k {
		return nil, fmt.Errorf("too few shares")
	}

	pubShares := []*share.PubShare{}
	for _, s := range shares {
		if group.String() != s.group.String() {
			return nil, fmt.Errorf("incorrect share group: %q", s.group)
		}
		pubShares = append(pubShares, &share.PubShare{
			I: s.index,
			V: s.u_i,
		})
	}

	arg, err := share.RecoverCommit(group, pubShares, k, n)
	if err != nil {
		return nil, fmt.Errorf("cannot recover secret: %w", err)
	}

	h, err := hash1(c.group.String(), arg)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal %q: %w", arg, err)
	}

	return xor(h, c.c)
}

// ciphertextRaw is used for Ciphertext (un)marshaling.
type ciphertextRaw struct {
	Group string
	C     []byte
	Label []byte
	U     []byte
	U_bar []byte
	E     []byte
	F     []byte
}

func (c Ciphertext) Marshal() ([]byte, error) {
	u, err := c.u.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal U: %w", err)
	}
	ubar, err := c.u_bar.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal U_bar: %w", err)
	}
	f, err := c.f.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal F: %w", err)
	}
	e, err := c.e.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal E: %w", err)
	}
	return json.Marshal(&ciphertextRaw{
		Group: c.group.String(),
		C:     c.c,
		Label: c.label,
		U:     u,
		U_bar: ubar,
		E:     e,
		F:     f,
	})
}

func (c *Ciphertext) Unmarshal(data []byte) error {
	var raw ciphertextRaw
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("cannot unmarshal data: %w", err)
	}
	c.c = raw.C
	c.label = raw.Label
	c.group, err = parseGroup(raw.Group)
	if err != nil {
		return fmt.Errorf("cannot parse group: %w", err)
	}
	c.e = c.group.Scalar()
	if err = c.e.UnmarshalBinary(raw.E); err != nil {
		return fmt.Errorf("cannot unmarshal E: %w", err)
	}
	c.u = c.group.Point()
	if err = c.u.UnmarshalBinary(raw.U); err != nil {
		return fmt.Errorf("cannot unmarshal U: %w", err)
	}
	c.u_bar = c.group.Point()
	if err = c.u_bar.UnmarshalBinary(raw.U_bar); err != nil {
		return fmt.Errorf("cannot unmarshal U_bar: %w", err)
	}
	c.f = c.group.Scalar()
	if err = c.f.UnmarshalBinary(raw.F); err != nil {
		return fmt.Errorf("cannot unmarshal F: %w", err)
	}
	return nil
}

// DecryptionShare defines a decryption share
type DecryptionShare struct {
	group group.Group
	index int
	u_i   group.Point
	e_i   group.Scalar
	f_i   group.Scalar
}

// TODO(pszal): test + fix tests which currently ignore share equality
func (a *DecryptionShare) Equal(b *DecryptionShare) bool {
	if a.group.String() != b.group.String() ||
		a.index != b.index ||
		!a.u_i.Equal(b.u_i) ||
		!a.e_i.Equal(b.e_i) ||
		!a.f_i.Equal(b.f_i) {
		return false
	}
	return true
}

func (d DecryptionShare) Index() int {
	return d.index
}

// decryptionShareRaw is used for DecryptionShare (un)marshaling.
type decryptionShareRaw struct {
	Group string
	Index int
	U_i   []byte
	E_i   []byte
	F_i   []byte
}

func (d DecryptionShare) Marshal() ([]byte, error) {
	u, err := d.u_i.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal U_i: %w", err)
	}
	f, err := d.f_i.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal F_i: %w", err)
	}
	e, err := d.e_i.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("cannot marshal E_i: %w", err)
	}
	return json.Marshal(&decryptionShareRaw{
		Group: d.group.String(),
		Index: d.index,
		U_i:   u,
		E_i:   e,
		F_i:   f,
	})
}

func (d *DecryptionShare) Unmarshal(data []byte) error {
	var raw decryptionShareRaw
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return fmt.Errorf("cannot unmarshal data: %w", err)
	}
	d.index = raw.Index
	d.group, err = parseGroup(raw.Group)
	if err != nil {
		return fmt.Errorf("cannot parse group: %w", err)
	}
	d.e_i = d.group.Scalar()
	if err = d.e_i.UnmarshalBinary(raw.E_i); err != nil {
		return fmt.Errorf("cannot unmarshal E_i: %w", err)
	}
	d.u_i = d.group.Point()
	if err = d.u_i.UnmarshalBinary(raw.U_i); err != nil {
		return fmt.Errorf("cannot unmarshal U_i: %w", err)
	}
	d.f_i = d.group.Scalar()
	if err = d.f_i.UnmarshalBinary(raw.F_i); err != nil {
		return fmt.Errorf("cannot unmarshal F_i: %w", err)
	}
	return nil
}

// hash is a generic hash function
func hash(msg []byte) []byte {
	h := defaultHash()
	h.Write(msg)
	return h.Sum(nil)
}

// hash1 is an implementation of the H_1 hash function (see p15 of the paper).
func hash1(group string, g group.Point) ([]byte, error) {
	point, err := concatenate(group, g)
	if err != nil {
		return nil, fmt.Errorf("cannot concatenate points: %w", err)
	}
	return hash(append([]byte("tdh2hash1"), point...)), nil
}

// hash2 is an implementation of the H_2 hash function (see p15 of the paper).
func hash2(msg, label []byte, g1, g2, g3, g4 group.Point, group group.Group) (group.Scalar, error) {
	if len(msg) != len(label) || len(msg) != InputSize {
		return nil, fmt.Errorf("message and label must be %dB long", InputSize)
	}

	points, err := concatenate(group.String(), g1, g2, g3, g4)
	if err != nil {
		return nil, fmt.Errorf("cannot concatenate points: %w", err)
	}
	input := []byte("tdh2hash2")
	for _, arg := range [][]byte{msg, label, points} {
		input = append(input, arg...)
	}

	return group.Scalar().SetBytes(hash(input)), nil
}

// hash4 is an implementation of the H_4 hash function (see p15 of the paper).
func hash4(g1, g2, g3 group.Point, group group.Group) (group.Scalar, error) {
	points, err := concatenate(group.String(), g1, g2, g3)
	if err != nil {
		return nil, fmt.Errorf("cannot concatenate points: %w", err)
	}
	h := hash(append([]byte("tdh2hash4"), points...))

	return group.Scalar().SetBytes(h), nil
}

// concatenate marshals and concatenates points (elements of a group). It is
// used in hash functions.
func concatenate(group string, points ...group.Point) ([]byte, error) {
	final := group
	for _, point := range points {
		p, err := point.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("cannot marshal point=%v err=%v", point, err)
		}
		final += "," + hex.EncodeToString(p)
	}
	return []byte(final), nil
}

// xor computes and returns XOR between two slices.
func xor(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}
	buf := make([]byte, len(a))
	for i := range a {
		buf[i] = a[i] ^ b[i]
	}
	return buf, nil
}
