package anon

import (
	"crypto/subtle"
	"errors"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/key"
)

func header(suite Suite, X kyber.Point, x kyber.Scalar,
	Xb, xb []byte, anonymitySet Set) []byte {

	//fmt.Printf("Xb %s\nxb %s\n",
	//		hex.EncodeToString(Xb),hex.EncodeToString(xb))

	// Encrypt the master scalar key with each public key in the set
	S := suite.Point()
	hdr := Xb
	for i := range anonymitySet {
		Y := anonymitySet[i]
		S.Mul(x, Y) // compute DH shared secret
		seed, _ := S.MarshalBinary()
		xof := suite.XOF(seed)
		xc := make([]byte, len(xb))
		xof.XORKeyStream(xc, xb)
		hdr = append(hdr, xc...)
	}
	return hdr
}

// Create and encrypt a fresh key decryptable only by the given receivers.
// Returns the secret key and the ciphertext.
func encryptKey(suite Suite, anonymitySet Set) (k, c []byte) {
	// Choose a keypair and encode its representation
	kp := new(key.Pair)
	var Xb []byte
	kp.Gen(suite)
	Xb, _ = kp.Public.MarshalBinary()
	xb, _ := kp.Private.MarshalBinary()
	// Generate the ciphertext header
	return xb, header(suite, kp.Public, kp.Private, Xb, xb, anonymitySet)
}

// Decrypt and verify a key encrypted via encryptKey.
// On success, returns the key and the length of the decrypted header.
func decryptKey(suite Suite, ciphertext []byte, anonymitySet Set, mine int, privateKey kyber.Scalar) ([]byte, int, error) {
	// Decode the (supposed) ephemeral public key from the front
	X := suite.Point()
	var Xb []byte
	enclen := X.MarshalSize()
	if len(ciphertext) < enclen {
		return nil, 0, errors.New("ciphertext too short")
	}
	if err := X.UnmarshalBinary(ciphertext[:enclen]); err != nil {
		return nil, 0, err
	}
	Xb = ciphertext[:enclen]
	Xblen := len(Xb)

	// Decode the (supposed) master secret with our private key
	nkeys := len(anonymitySet)
	if mine < 0 || mine >= nkeys {
		panic("private-key index out of range")
	}
	seclen := suite.ScalarLen()
	if len(ciphertext) < Xblen+seclen*nkeys {
		return nil, 0, errors.New("ciphertext too short")
	}
	S := suite.Point().Mul(privateKey, X)
	seed, _ := S.MarshalBinary()
	xof := suite.XOF(seed)
	xb := make([]byte, seclen)
	secofs := Xblen + seclen*mine
	xof.XORKeyStream(xb, ciphertext[secofs:secofs+seclen])
	x := suite.Scalar()
	if err := x.UnmarshalBinary(xb); err != nil {
		return nil, 0, err
	}

	// Make sure it reproduces the correct ephemeral public key
	Xv := suite.Point().Mul(x, nil)
	if !X.Equal(Xv) {
		return nil, 0, errors.New("invalid ciphertext")
	}

	// Regenerate and check the rest of the header,
	// to ensure that that any of the anonymitySet members could decrypt it
	hdr := header(suite, X, x, Xb, xb, anonymitySet)
	hdrlen := len(hdr)
	if hdrlen != Xblen+seclen*nkeys {
		panic("wrong header size")
	}
	if subtle.ConstantTimeCompare(hdr, ciphertext[:hdrlen]) == 0 {
		return nil, 0, errors.New("invalid ciphertext")
	}

	return xb, hdrlen, nil
}

// constantTimeAllEq returns 1 iff all bytes in slice x have the value y.
// The time taken is a function of the length of the slices
// and is independent of the contents.
func constantTimeAllEq(x []byte, y byte) int {
	var z byte
	for _, b := range x {
		z |= b ^ y
	}
	return subtle.ConstantTimeByteEq(z, 0)
}

// macSize is how long the hashes are that we extract from the XOF.
// This constant of 16 is taken from the previous implementation's behavior.
const macSize = 16

// Encrypt a message for reading by any member of an explit anonymity set.
// The caller supplies one or more keys representing the anonymity set.
// If the provided set contains only one public key,
// this reduces to conventional single-receiver public-key encryption.
func Encrypt(suite Suite, message []byte,
	anonymitySet Set) []byte {

	xb, hdr := encryptKey(suite, anonymitySet)
	xof := suite.XOF(xb)

	// We now know the ciphertext layout
	hdrhi := 0 + len(hdr)
	msghi := hdrhi + len(message)
	machi := msghi + macSize
	ciphertext := make([]byte, machi)
	copy(ciphertext, hdr)

	// Now encrypt and MAC the message based on the master secret
	ctx := ciphertext[hdrhi:msghi]
	mac := ciphertext[msghi:machi]

	xof.XORKeyStream(ctx, message)
	xof = suite.XOF(ctx)
	xof.Read(mac)

	return ciphertext
}

// Decrypt a message encrypted for a particular anonymity set.
// Returns the cleartext message on success, or an error on failure.
//
// The caller provides the anonymity set for which the message is intended,
// and the private key corresponding to one of the public keys in the set.
// Decrypt verifies that the message is encrypted correctly for this set -
// in particular, that it could be decrypted by ALL of the listed members -
// before returning successfully with the decrypted message.
//
// This verification ensures that a malicious sender
// cannot de-anonymize a receiver by constructing a ciphertext incorrectly
// so as to be decryptable by only some members of the set.
// As a side-effect, this verification also ensures plaintext-awareness:
// that is, it is infeasible for a sender to construct any ciphertext
// that will be accepted by the receiver without knowing the plaintext.
//
func Decrypt(suite Suite, ciphertext []byte, anonymitySet Set, mine int, privateKey kyber.Scalar) ([]byte, error) {
	// Decrypt and check the encrypted key-header.
	xb, hdrlen, err := decryptKey(suite, ciphertext, anonymitySet,
		mine, privateKey)
	if err != nil {
		return nil, err
	}

	// Determine the message layout
	xof := suite.XOF(xb)
	if len(ciphertext) < hdrlen+macSize {
		return nil, errors.New("ciphertext too short")
	}
	hdrhi := hdrlen
	msghi := len(ciphertext) - macSize

	// Decrypt the message and check the MAC
	ctx := ciphertext[hdrhi:msghi]
	mac := ciphertext[msghi:]
	msg := make([]byte, len(ctx))
	xof.XORKeyStream(msg, ctx)
	xof = suite.XOF(ctx)
	xof.XORKeyStream(mac, mac)
	if constantTimeAllEq(mac, 0) == 0 {
		return nil, errors.New("invalid ciphertext: failed MAC check")
	}
	return msg, nil
}
