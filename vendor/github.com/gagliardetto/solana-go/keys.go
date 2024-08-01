// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solana

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"sort"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type PrivateKey []byte

func MustPrivateKeyFromBase58(in string) PrivateKey {
	out, err := PrivateKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PrivateKeyFromBase58(privkey string) (PrivateKey, error) {
	res, err := base58.Decode(privkey)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PrivateKeyFromSolanaKeygenFile(file string) (PrivateKey, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read keygen file: %w", err)
	}

	var values []byte
	err = json.Unmarshal(content, &values)
	if err != nil {
		return nil, fmt.Errorf("decode keygen file: %w", err)
	}

	return PrivateKey([]byte(values)), nil
}

func (k PrivateKey) String() string {
	return base58.Encode(k)
}

func NewRandomPrivateKey() (PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(crypto_rand.Reader)
	if err != nil {
		return nil, err
	}
	var publicKey PublicKey
	copy(publicKey[:], pub)
	return PrivateKey(priv), nil
}

func (k PrivateKey) Sign(payload []byte) (Signature, error) {
	p := ed25519.PrivateKey(k)
	signData, err := p.Sign(crypto_rand.Reader, payload, crypto.Hash(0))
	if err != nil {
		return Signature{}, err
	}

	var signature Signature
	copy(signature[:], signData)

	return signature, err
}

func (k PrivateKey) PublicKey() PublicKey {
	p := ed25519.PrivateKey(k)
	pub := p.Public().(ed25519.PublicKey)

	var publicKey PublicKey
	copy(publicKey[:], pub)

	return publicKey
}

// PK is a convenience alias for PublicKey
type PK = PublicKey

func (p PublicKey) Verify(message []byte, signature Signature) bool {
	pub := ed25519.PublicKey(p[:])
	return ed25519.Verify(pub, message, signature[:])
}

type PublicKey [PublicKeyLength]byte

func PublicKeyFromBytes(in []byte) (out PublicKey) {
	byteCount := len(in)
	if byteCount == 0 {
		return
	}

	max := PublicKeyLength
	if byteCount < max {
		max = byteCount
	}

	copy(out[:], in[0:max])
	return
}

// MPK is a convenience alias for MustPublicKeyFromBase58
func MPK(in string) PublicKey {
	return MustPublicKeyFromBase58(in)
}

func MustPublicKeyFromBase58(in string) PublicKey {
	out, err := PublicKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PublicKeyFromBase58(in string) (out PublicKey, err error) {
	val, err := base58.Decode(in)
	if err != nil {
		return out, fmt.Errorf("decode: %w", err)
	}

	if len(val) != PublicKeyLength {
		return out, fmt.Errorf("invalid length, expected %v, got %d", PublicKeyLength, len(val))
	}

	copy(out[:], val)
	return
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return []byte(base58.Encode(p[:])), nil
}

func (p *PublicKey) UnmarshalText(data []byte) error {
	return p.Set(string(data))
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(p[:]))
}

func (p *PublicKey) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*p, err = PublicKeyFromBase58(s)
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", s, err)
	}
	return
}

// MarshalBSON implements the bson.Marshaler interface.
func (p PublicKey) MarshalBSON() ([]byte, error) {
	return bson.Marshal(p.String())
}

// UnmarshalBSON implements the bson.Unmarshaler interface.
func (p *PublicKey) UnmarshalBSON(data []byte) (err error) {
	var s string
	if err := bson.Unmarshal(data, &s); err != nil {
		return err
	}

	*p, err = PublicKeyFromBase58(s)
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", s, err)
	}
	return nil
}

// MarshalBSONValue implements the bson.ValueMarshaler interface.
func (p PublicKey) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(p.String())
}

// UnmarshalBSONValue implements the bson.ValueUnmarshaler interface.
func (p *PublicKey) UnmarshalBSONValue(t bsontype.Type, data []byte) (err error) {
	var s string
	if err := bson.Unmarshal(data, &s); err != nil {
		return err
	}

	*p, err = PublicKeyFromBase58(s)
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", s, err)
	}
	return nil
}

func (p PublicKey) Equals(pb PublicKey) bool {
	return p == pb
}

// IsAnyOf checks if p is equals to any of the provided keys.
func (p PublicKey) IsAnyOf(keys ...PublicKey) bool {
	for _, k := range keys {
		if p.Equals(k) {
			return true
		}
	}
	return false
}

// ToPointer returns a pointer to the pubkey.
func (p PublicKey) ToPointer() *PublicKey {
	return &p
}

func (p PublicKey) Bytes() []byte {
	return []byte(p[:])
}

// Check if a `Pubkey` is on the ed25519 curve.
func (p PublicKey) IsOnCurve() bool {
	return IsOnCurve(p[:])
}

var zeroPublicKey = PublicKey{}

// IsZero returns whether the public key is zero.
// NOTE: the System Program public key is also zero.
func (p PublicKey) IsZero() bool {
	return p == zeroPublicKey
}

func (p *PublicKey) Set(s string) (err error) {
	*p, err = PublicKeyFromBase58(s)
	if err != nil {
		return fmt.Errorf("invalid public key %s: %w", s, err)
	}
	return
}

func (p PublicKey) String() string {
	return base58.Encode(p[:])
}

// Short returns a shortened pubkey string,
// only including the first n chars, ellipsis, and the last n characters.
// NOTE: this is ONLY for visual representation for humans,
// and cannot be used for anything else.
func (p PublicKey) Short(n int) string {
	return formatShortPubkey(n, p)
}

func formatShortPubkey(n int, pubkey PublicKey) string {
	str := pubkey.String()
	if n > (len(str)/2)-1 {
		n = (len(str) / 2) - 1
	}
	if n < 2 {
		n = 2
	}
	return str[:n] + "..." + str[len(str)-n:]
}

type PublicKeySlice []PublicKey

// UniqueAppend appends the provided pubkey only if it is not
// already present in the slice.
// Returns true when the provided pubkey wasn't already present.
func (slice *PublicKeySlice) UniqueAppend(pubkey PublicKey) bool {
	if !slice.Has(pubkey) {
		slice.Append(pubkey)
		return true
	}
	return false
}

func (slice *PublicKeySlice) Append(pubkeys ...PublicKey) {
	*slice = append(*slice, pubkeys...)
}

func (slice PublicKeySlice) Has(pubkey PublicKey) bool {
	for _, key := range slice {
		if key.Equals(pubkey) {
			return true
		}
	}
	return false
}

func (slice PublicKeySlice) Len() int {
	return len(slice)
}

func (slice PublicKeySlice) Less(i, j int) bool {
	return bytes.Compare(slice[i][:], slice[j][:]) < 0
}

func (slice PublicKeySlice) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Sort sorts the slice.
func (slice PublicKeySlice) Sort() {
	sort.Sort(slice)
}

// Dedupe returns a new slice with all duplicate pubkeys removed.
func (slice PublicKeySlice) Dedupe() PublicKeySlice {
	slice.Sort()
	deduped := make(PublicKeySlice, 0)
	for i := 0; i < len(slice); i++ {
		if i == 0 || !slice[i].Equals(slice[i-1]) {
			deduped = append(deduped, slice[i])
		}
	}
	return deduped
}

// Contains returns true if the slice contains the provided pubkey.
func (slice PublicKeySlice) Contains(pubkey PublicKey) bool {
	for _, key := range slice {
		if key.Equals(pubkey) {
			return true
		}
	}
	return false
}

// ContainsAll returns true if all the provided pubkeys are present in the slice.
func (slice PublicKeySlice) ContainsAll(pubkeys PublicKeySlice) bool {
	for _, pubkey := range pubkeys {
		if !slice.Contains(pubkey) {
			return false
		}
	}
	return true
}

// ContainsAny returns true if any of the provided pubkeys are present in the slice.
func (slice PublicKeySlice) ContainsAny(pubkeys PublicKeySlice) bool {
	for _, pubkey := range pubkeys {
		if slice.Contains(pubkey) {
			return true
		}
	}
	return false
}

func (slice PublicKeySlice) ToBase58() []string {
	out := make([]string, len(slice))
	for i, pubkey := range slice {
		out[i] = pubkey.String()
	}
	return out
}

func (slice PublicKeySlice) ToBytes() [][]byte {
	out := make([][]byte, len(slice))
	for i, pubkey := range slice {
		out[i] = pubkey.Bytes()
	}
	return out
}

func (slice PublicKeySlice) ToPointers() []*PublicKey {
	out := make([]*PublicKey, len(slice))
	for i, pubkey := range slice {
		out[i] = pubkey.ToPointer()
	}
	return out
}

// Removed returns the elements that are present in `a` but not in `b`.
func (a PublicKeySlice) Removed(b PublicKeySlice) PublicKeySlice {
	var diff PublicKeySlice
	for _, pubkey := range a {
		if !b.Contains(pubkey) {
			diff = append(diff, pubkey)
		}
	}
	return diff.Dedupe()
}

// Added returns the elements that are present in `b` but not in `a`.
func (a PublicKeySlice) Added(b PublicKeySlice) PublicKeySlice {
	return b.Removed(a)
}

// Intersect returns the intersection of two PublicKeySlices, i.e. the elements
// that are in both PublicKeySlices.
// The returned PublicKeySlice is sorted and deduped.
func (prev PublicKeySlice) Intersect(next PublicKeySlice) PublicKeySlice {
	var intersect PublicKeySlice
	for _, pubkey := range prev {
		if next.Contains(pubkey) {
			intersect = append(intersect, pubkey)
		}
	}
	return intersect.Dedupe()
}

// Equals returns true if the two PublicKeySlices are equal (same order of same keys).
func (slice PublicKeySlice) Equals(other PublicKeySlice) bool {
	if len(slice) != len(other) {
		return false
	}
	for i, pubkey := range slice {
		if !pubkey.Equals(other[i]) {
			return false
		}
	}
	return true
}

// Same returns true if the two slices contain the same public keys,
// but not necessarily in the same order.
func (slice PublicKeySlice) Same(other PublicKeySlice) bool {
	if len(slice) != len(other) {
		return false
	}
	for _, pubkey := range slice {
		if !other.Contains(pubkey) {
			return false
		}
	}
	return true
}

// Split splits the slice into chunks of the specified size.
func (slice PublicKeySlice) Split(chunkSize int) []PublicKeySlice {
	divided := make([]PublicKeySlice, 0)
	if len(slice) == 0 || chunkSize < 1 {
		return divided
	}
	if len(slice) == 1 {
		return append(divided, slice)
	}

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		divided = append(divided, slice[i:end])
	}

	return divided
}

// Last returns the last element of the slice.
// Returns nil if the slice is empty.
func (slice PublicKeySlice) Last() *PublicKey {
	if len(slice) == 0 {
		return nil
	}
	return slice[len(slice)-1].ToPointer()
}

// First returns the first element of the slice.
// Returns nil if the slice is empty.
func (slice PublicKeySlice) First() *PublicKey {
	if len(slice) == 0 {
		return nil
	}
	return slice[0].ToPointer()
}

// GetAddedRemoved compares to the `next` pubkey slice, and returns
// two slices:
// - `added` is the slice of pubkeys that are present in `next` but NOT present in `previous`.
// - `removed` is the slice of pubkeys that are present in `previous` but are NOT present in `next`.
func (prev PublicKeySlice) GetAddedRemoved(next PublicKeySlice) (added PublicKeySlice, removed PublicKeySlice) {
	return next.Removed(prev), prev.Removed(next)
}

// GetAddedRemovedPubkeys accepts two slices of pubkeys (`previous` and `next`), and returns
// two slices:
// - `added` is the slice of pubkeys that are present in `next` but NOT present in `previous`.
// - `removed` is the slice of pubkeys that are present in `previous` but are NOT present in `next`.
func GetAddedRemovedPubkeys(previous PublicKeySlice, next PublicKeySlice) (added PublicKeySlice, removed PublicKeySlice) {
	added = make(PublicKeySlice, 0)
	removed = make(PublicKeySlice, 0)

	for _, prev := range previous {
		if !next.Has(prev) {
			removed = append(removed, prev)
		}
	}

	for _, nx := range next {
		if !previous.Has(nx) {
			added = append(added, nx)
		}
	}

	return
}

var nativeProgramIDs = PublicKeySlice{
	BPFLoaderProgramID,
	BPFLoaderDeprecatedProgramID,
	FeatureProgramID,
	ConfigProgramID,
	StakeProgramID,
	VoteProgramID,
	Secp256k1ProgramID,
	SystemProgramID,
	SysVarClockPubkey,
	SysVarEpochSchedulePubkey,
	SysVarFeesPubkey,
	SysVarInstructionsPubkey,
	SysVarRecentBlockHashesPubkey,
	SysVarRentPubkey,
	SysVarRewardsPubkey,
	SysVarSlotHashesPubkey,
	SysVarSlotHistoryPubkey,
	SysVarStakeHistoryPubkey,
}

// https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L372
func isNativeProgramID(key PublicKey) bool {
	return nativeProgramIDs.Has(key)
}

const (
	/// Number of bytes in a pubkey.
	PublicKeyLength = 32
	// Maximum length of derived pubkey seed.
	MaxSeedLength = 32
	// Maximum number of seeds.
	MaxSeeds = 16
	/// Number of bytes in a signature.
	SignatureLength = 64

	// // Maximum string length of a base58 encoded pubkey.
	// MaxBase58Length = 44
)

// Ported from https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L159
func CreateWithSeed(base PublicKey, seed string, owner PublicKey) (PublicKey, error) {
	if len(seed) > MaxSeedLength {
		return PublicKey{}, errors.New("Max seed length exceeded")
	}

	// let owner = owner.as_ref();
	// if owner.len() >= PDA_MARKER.len() {
	//     let slice = &owner[owner.len() - PDA_MARKER.len()..];
	//     if slice == PDA_MARKER {
	//         return Err(PubkeyError::IllegalOwner);
	//     }
	// }

	b := make([]byte, 0, 64+len(seed))
	b = append(b, base[:]...)
	b = append(b, seed[:]...)
	b = append(b, owner[:]...)
	hash := sha256.Sum256(b)
	return PublicKeyFromBytes(hash[:]), nil
}

const PDA_MARKER = "ProgramDerivedAddress"

var ErrMaxSeedLengthExceeded = errors.New("Max seed length exceeded")

// Create a program address.
// Ported from https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L204
func CreateProgramAddress(seeds [][]byte, programID PublicKey) (PublicKey, error) {
	if len(seeds) > MaxSeeds {
		return PublicKey{}, ErrMaxSeedLengthExceeded
	}

	for _, seed := range seeds {
		if len(seed) > MaxSeedLength {
			return PublicKey{}, ErrMaxSeedLengthExceeded
		}
	}

	buf := []byte{}
	for _, seed := range seeds {
		buf = append(buf, seed...)
	}

	buf = append(buf, programID[:]...)
	buf = append(buf, []byte(PDA_MARKER)...)
	hash := sha256.Sum256(buf)

	if IsOnCurve(hash[:]) {
		return PublicKey{}, errors.New("invalid seeds; address must fall off the curve")
	}

	return PublicKeyFromBytes(hash[:]), nil
}

// Check if the provided `b` is on the ed25519 curve.
func IsOnCurve(b []byte) bool {
	_, err := new(edwards25519.Point).SetBytes(b)
	isOnCurve := err == nil
	return isOnCurve
}

// Find a valid program address and its corresponding bump seed.
func FindProgramAddress(seed [][]byte, programID PublicKey) (PublicKey, uint8, error) {
	var address PublicKey
	var err error
	bumpSeed := uint8(math.MaxUint8)
	for bumpSeed != 0 {
		address, err = CreateProgramAddress(append(seed, []byte{byte(bumpSeed)}), programID)
		if err == nil {
			return address, bumpSeed, nil
		}
		bumpSeed--
	}
	return PublicKey{}, bumpSeed, errors.New("unable to find a valid program address")
}

func FindAssociatedTokenAddress(
	wallet PublicKey,
	mint PublicKey,
) (PublicKey, uint8, error) {
	return findAssociatedTokenAddressAndBumpSeed(
		wallet,
		mint,
		SPLAssociatedTokenAccountProgramID,
	)
}

func findAssociatedTokenAddressAndBumpSeed(
	walletAddress PublicKey,
	splTokenMintAddress PublicKey,
	programID PublicKey,
) (PublicKey, uint8, error) {
	return FindProgramAddress([][]byte{
		walletAddress[:],
		TokenProgramID[:],
		splTokenMintAddress[:],
	},
		programID,
	)
}

// FindTokenMetadataAddress returns the token metadata program-derived address given a SPL token mint address.
func FindTokenMetadataAddress(mint PublicKey) (PublicKey, uint8, error) {
	seed := [][]byte{
		[]byte("metadata"),
		TokenMetadataProgramID[:],
		mint[:],
	}
	return FindProgramAddress(seed, TokenMetadataProgramID)
}
