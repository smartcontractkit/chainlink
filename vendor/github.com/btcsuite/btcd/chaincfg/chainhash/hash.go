// Copyright (c) 2013-2016 The btcsuite developers
// Copyright (c) 2015 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package chainhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// HashSize of array used to store hashes.  See Hash.
const HashSize = 32

// MaxHashStringSize is the maximum length of a Hash hash string.
const MaxHashStringSize = HashSize * 2

var (
	// TagBIP0340Challenge is the BIP-0340 tag for challenges.
	TagBIP0340Challenge = []byte("BIP0340/challenge")

	// TagBIP0340Aux is the BIP-0340 tag for aux data.
	TagBIP0340Aux = []byte("BIP0340/aux")

	// TagBIP0340Nonce is the BIP-0340 tag for nonces.
	TagBIP0340Nonce = []byte("BIP0340/nonce")

	// TagTapSighash is the tag used by BIP 341 to generate the sighash
	// flags.
	TagTapSighash = []byte("TapSighash")

	// TagTagTapLeaf is the message tag prefix used to compute the hash
	// digest of a tapscript leaf.
	TagTapLeaf = []byte("TapLeaf")

	// TagTapBranch is the message tag prefix used to compute the
	// hash digest of two tap leaves into a taproot branch node.
	TagTapBranch = []byte("TapBranch")

	// TagTapTweak is the message tag prefix used to compute the hash tweak
	// used to enable a public key to commit to the taproot branch root
	// for the witness program.
	TagTapTweak = []byte("TapTweak")

	// precomputedTags is a map containing the SHA-256 hash of the BIP-0340
	// tags.
	precomputedTags = map[string]Hash{
		string(TagBIP0340Challenge): sha256.Sum256(TagBIP0340Challenge),
		string(TagBIP0340Aux):       sha256.Sum256(TagBIP0340Aux),
		string(TagBIP0340Nonce):     sha256.Sum256(TagBIP0340Nonce),
		string(TagTapSighash):       sha256.Sum256(TagTapSighash),
		string(TagTapLeaf):          sha256.Sum256(TagTapLeaf),
		string(TagTapBranch):        sha256.Sum256(TagTapBranch),
		string(TagTapTweak):         sha256.Sum256(TagTapTweak),
	}
)

// ErrHashStrSize describes an error that indicates the caller specified a hash
// string that has too many characters.
var ErrHashStrSize = fmt.Errorf("max hash string length is %v bytes", MaxHashStringSize)

// Hash is used in several of the bitcoin messages and common structures.  It
// typically represents the double sha256 of data.
type Hash [HashSize]byte

// String returns the Hash as the hexadecimal string of the byte-reversed
// hash.
func (hash Hash) String() string {
	for i := 0; i < HashSize/2; i++ {
		hash[i], hash[HashSize-1-i] = hash[HashSize-1-i], hash[i]
	}
	return hex.EncodeToString(hash[:])
}

// CloneBytes returns a copy of the bytes which represent the hash as a byte
// slice.
//
// NOTE: It is generally cheaper to just slice the hash directly thereby reusing
// the same bytes rather than calling this method.
func (hash *Hash) CloneBytes() []byte {
	newHash := make([]byte, HashSize)
	copy(newHash, hash[:])

	return newHash
}

// SetBytes sets the bytes which represent the hash.  An error is returned if
// the number of bytes passed in is not HashSize.
func (hash *Hash) SetBytes(newHash []byte) error {
	nhlen := len(newHash)
	if nhlen != HashSize {
		return fmt.Errorf("invalid hash length of %v, want %v", nhlen,
			HashSize)
	}
	copy(hash[:], newHash)

	return nil
}

// IsEqual returns true if target is the same as hash.
func (hash *Hash) IsEqual(target *Hash) bool {
	if hash == nil && target == nil {
		return true
	}
	if hash == nil || target == nil {
		return false
	}
	return *hash == *target
}

// NewHash returns a new Hash from a byte slice.  An error is returned if
// the number of bytes passed in is not HashSize.
func NewHash(newHash []byte) (*Hash, error) {
	var sh Hash
	err := sh.SetBytes(newHash)
	if err != nil {
		return nil, err
	}
	return &sh, err
}

// TaggedHash implements the tagged hash scheme described in BIP-340. We use
// sha-256 to bind a message hash to a specific context using a tag:
// sha256(sha256(tag) || sha256(tag) || msg).
func TaggedHash(tag []byte, msgs ...[]byte) *Hash {
	// Check to see if we've already pre-computed the hash of the tag. If
	// so then this'll save us an extra sha256 hash.
	shaTag, ok := precomputedTags[string(tag)]
	if !ok {
		shaTag = sha256.Sum256(tag)
	}

	// h = sha256(sha256(tag) || sha256(tag) || msg)
	h := sha256.New()
	h.Write(shaTag[:])
	h.Write(shaTag[:])

	for _, msg := range msgs {
		h.Write(msg)
	}

	taggedHash := h.Sum(nil)

	// The function can't error out since the above hash is guaranteed to
	// be 32 bytes.
	hash, _ := NewHash(taggedHash)

	return hash
}

// NewHashFromStr creates a Hash from a hash string.  The string should be
// the hexadecimal string of a byte-reversed hash, but any missing characters
// result in zero padding at the end of the Hash.
func NewHashFromStr(hash string) (*Hash, error) {
	ret := new(Hash)
	err := Decode(ret, hash)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Decode decodes the byte-reversed hexadecimal string encoding of a Hash to a
// destination.
func Decode(dst *Hash, src string) error {
	// Return error if hash string is too long.
	if len(src) > MaxHashStringSize {
		return ErrHashStrSize
	}

	// Hex decoder expects the hash to be a multiple of two.  When not, pad
	// with a leading zero.
	var srcBytes []byte
	if len(src)%2 == 0 {
		srcBytes = []byte(src)
	} else {
		srcBytes = make([]byte, 1+len(src))
		srcBytes[0] = '0'
		copy(srcBytes[1:], src)
	}

	// Hex decode the source bytes to a temporary destination.
	var reversedHash Hash
	_, err := hex.Decode(reversedHash[HashSize-hex.DecodedLen(len(srcBytes)):], srcBytes)
	if err != nil {
		return err
	}

	// Reverse copy from the temporary hash to destination.  Because the
	// temporary was zeroed, the written result will be correctly padded.
	for i, b := range reversedHash[:HashSize/2] {
		dst[i], dst[HashSize-1-i] = reversedHash[HashSize-1-i], b
	}

	return nil
}
