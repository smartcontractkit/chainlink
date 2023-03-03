package anon

import (
	"bytes"
	"errors"

	"go.dedis.ch/kyber/v3"
)

// unlinkable ring signature
type uSig struct {
	C0 kyber.Scalar
	S  []kyber.Scalar
}

// linkable ring signature
type lSig struct {
	C0  kyber.Scalar
	S   []kyber.Scalar
	Tag kyber.Point
}

func signH1pre(suite Suite, linkScope []byte, linkTag kyber.Point,
	message []byte) kyber.XOF {
	H1pre := suite.XOF(message) // m
	if linkScope != nil {
		_, _ = H1pre.Write(linkScope) // L
		tag, _ := linkTag.MarshalBinary()
		_, _ = H1pre.Write(tag) // ~y
	}
	return H1pre
}

func signH1(suite Suite, H1pre kyber.XOF, PG, PH kyber.Point) kyber.Scalar {
	H1 := H1pre.Clone()
	PGb, _ := PG.MarshalBinary()
	_, _ = H1.Write(PGb)
	if PH != nil {
		PHb, _ := PH.MarshalBinary()
		_, _ = H1.Write(PHb)
	}
	return suite.Scalar().Pick(H1)
}

// Sign creates an optionally anonymous, optionally linkable
// signature on a given message.
//
// The caller supplies one or more public keys representing an anonymity set,
// and the private key corresponding to one of those public keys.
// The resulting signature proves to a verifier that the owner of
// one of these public keys signed the message,
// without revealing which key-holder signed the message,
// offering anonymity among the members of this explicit anonymity set.
// The other users whose keys are listed in the anonymity set need not consent
// or even be aware that they have been included in an anonymity set:
// anyone having a suitable public key may be "conscripted" into a set.
//
// If the provided anonymity set contains only one public key (the signer's),
// then this function produces a traditional non-anonymous signature,
// equivalent in both size and performance to a standard ElGamal signature.
//
// The caller may request either unlinkable or linkable anonymous signatures.
// If linkScope is nil, this function generates an unlinkable signature,
// which contains no information about which member signed the message.
// The anonymity provided by unlinkable signatures is forward-secure,
// in that a signature reveals nothing about which member generated it,
// even if all members' private keys are later released.
// For cryptographic background on unlinkable anonymity-set signatures -
// also known as ring signatures or ad-hoc group signatures -
// see Rivest, "How to Leak a Secret" at
// https://people.csail.mit.edu/rivest/pubs/RST01.pdf.
//
// If the caller passes a non-nil linkScope,
// the resulting anonymous signature will be linkable.
// This means that given two signatures produced using the same linkScope,
// a verifier will be able to tell whether
// the same or different anonymity set members produced those signatures.
// In particular, verifying a linkable signature yields a linkage tag.
// This linkage tag has a 1-to-1 correspondence with the signer's public key
// within a given linkScope, but is cryptographically unlinkable
// to either the signer's public key or to linkage tags in other scopes.
// The provided linkScope may be an arbitrary byte-string;
// the only significance these scopes have is whether they are equal or unequal.
// For details on the linkable signature algorithm this function implements,
// see Liu/Wei/Wong,
// "Linkable Spontaneous Anonymous Group Signature for Ad Hoc Groups" at
// https://eprint.iacr.org/2004/027.
//
// Linkage tags may be used to protect against sock-puppetry or Sybil attacks
// in situations where a verifier needs to know how many distinct members
// of an anonymity set are present or signed messages in a given context.
// It is cryptographically hard for one anonymity set member
// to produce signatures with different linkage tags in the same scope.
// An important and fundamental downside, however, is that
// linkable signatures do NOT offer forward-secure anonymity.
// If an anonymity set member's private key is later released,
// it is trivial to check whether or not that member produced a given signature.
// Also, anonymity set members who did NOT sign a message could
// (voluntarily or under coercion) prove that they did not sign it,
// e.g., simply by signing some other message in that linkage context
// and noting that the resulting linkage tag comes out different.
// Thus, linkable anonymous signatures are not appropriate to use
// in situations where there may be significant risk
// that members' private keys may later be compromised,
// or that members may be persuaded or coerced into revealing whether or not
// they produced a signature of interest.
//
func Sign(suite Suite, message []byte,
	anonymitySet Set, linkScope []byte, mine int, privateKey kyber.Scalar) []byte {

	// Note that Rivest's original ring construction directly supports
	// heterogeneous rings containing public keys of different types -
	// e.g., a mixture of RSA keys and DSA keys with varying parameters.
	// Our ring signature construction currently supports
	// only homogeneous rings containing compatible keys
	// drawn from the cipher suite (e.g., the same elliptic curve).
	// The upside to this constrint is greater flexibility:
	// e.g., we also easily obtain linkable ring signatures,
	// which are not readily feasible with the original ring construction.

	n := len(anonymitySet)           // anonymity set size
	L := []kyber.Point(anonymitySet) // public keys in anonymity set
	pi := mine

	// If we want a linkable ring signature, produce correct linkage tag,
	// as a pseudorandom base point multiplied by our private key.
	// Liu's scheme specifies the linkScope as a hash of the ring;
	// this is one reasonable choice of linkage scope,
	// but there are others, so we parameterize this choice.
	var linkBase, linkTag kyber.Point
	if linkScope != nil {
		linkStream := suite.XOF(linkScope)
		linkBase = suite.Point().Pick(linkStream)
		linkTag = suite.Point().Mul(privateKey, linkBase)
	}

	// First pre-hash the parameters to H1
	// that are invariant for different ring positions,
	// so that we don't have to hash them many times.
	H1pre := signH1pre(suite, linkScope, linkTag, message)

	// Pick a random commit for my ring position
	u := suite.Scalar().Pick(suite.RandomStream())
	var UB, UL kyber.Point
	UB = suite.Point().Mul(u, nil)
	if linkScope != nil {
		UL = suite.Point().Mul(u, linkBase)
	}

	// Build the challenge ring
	s := make([]kyber.Scalar, n)
	c := make([]kyber.Scalar, n)
	c[(pi+1)%n] = signH1(suite, H1pre, UB, UL)
	var P, PG, PH kyber.Point
	P = suite.Point()
	PG = suite.Point()
	if linkScope != nil {
		PH = suite.Point()
	}
	for i := (pi + 1) % n; i != pi; i = (i + 1) % n {
		s[i] = suite.Scalar().Pick(suite.RandomStream())
		PG.Add(PG.Mul(s[i], nil), P.Mul(c[i], L[i]))
		if linkScope != nil {
			PH.Add(PH.Mul(s[i], linkBase), P.Mul(c[i], linkTag))
		}
		c[(i+1)%n] = signH1(suite, H1pre, PG, PH)
		//fmt.Printf("s%d %s\n",i,s[i].String())
		//fmt.Printf("c%d %s\n",(i+1)%n,c[(i+1)%n].String())
	}
	s[pi] = suite.Scalar()
	s[pi].Mul(privateKey, c[pi]).Sub(u, s[pi]) // s_pi = u - x_pi c_pi

	// Encode and return the signature
	buf := bytes.Buffer{}
	if linkScope != nil { // linkable ring signature
		sig := lSig{c[0], s, linkTag}
		_ = suite.Write(&buf, &sig)
	} else { // unlinkable ring signature
		sig := uSig{c[0], s}
		_ = suite.Write(&buf, &sig)
	}
	return buf.Bytes()
}

// Verify checks a signature generated by Sign.
//
// The caller provides the message, anonymity set, and linkage scope
// with which the signature was purportedly produced.
// If the signature is a valid linkable signature (linkScope != nil),
// this function returns a linkage tag that uniquely corresponds
// to the signer within the given linkScope.
// If the signature is a valid unlinkable signature (linkScope == nil),
// Verify returns an empty but non-nil byte-slice instead of a linkage tag on success.
// Returns a nil linkage tag and an error if the signature is invalid.
func Verify(suite Suite, message []byte, anonymitySet Set,
	linkScope []byte, signatureBuffer []byte) ([]byte, error) {

	n := len(anonymitySet)           // anonymity set size
	L := []kyber.Point(anonymitySet) // public keys in ring

	// Decode the signature
	buf := bytes.NewBuffer(signatureBuffer)
	var linkBase, linkTag kyber.Point
	sig := lSig{}
	sig.S = make([]kyber.Scalar, n)
	if linkScope != nil { // linkable ring signature
		if err := suite.Read(buf, &sig); err != nil {
			return nil, err
		}
		linkStream := suite.XOF(linkScope)
		linkBase = suite.Point().Pick(linkStream)
		linkTag = sig.Tag
	} else { // unlinkable ring signature
		if err := suite.Read(buf, &sig.C0); err != nil {
			return nil, err
		}
		if err := suite.Read(buf, &sig.S); err != nil {
			return nil, err
		}
	}

	// Pre-hash the ring-position-invariant parameters to H1.
	H1pre := signH1pre(suite, linkScope, linkTag, message)

	// Verify the signature
	var P, PG, PH kyber.Point
	P = suite.Point()
	PG = suite.Point()
	if linkScope != nil {
		PH = suite.Point()
	}
	s := sig.S
	ci := sig.C0
	for i := 0; i < n; i++ {
		PG.Add(PG.Mul(s[i], nil), P.Mul(ci, L[i]))
		if linkScope != nil {
			PH.Add(PH.Mul(s[i], linkBase), P.Mul(ci, linkTag))
		}
		ci = signH1(suite, H1pre, PG, PH)
	}
	if !ci.Equal(sig.C0) {
		return nil, errors.New("invalid signature")
	}

	// Return the re-encoded linkage tag, for uniqueness checking
	if linkScope != nil {
		tag, _ := linkTag.MarshalBinary()
		return tag, nil
	}
	return []byte{}, nil
}
