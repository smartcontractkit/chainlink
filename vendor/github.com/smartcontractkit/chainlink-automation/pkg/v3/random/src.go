package random

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"

	"golang.org/x/crypto/sha3"
)

var (
	newCipherFn = aes.NewCipher
	randReadFn  = crand.Read
)

// Generates a randomness source derived from the prefix and seq # so
// that it's the same across the network for the same input.
func GetRandomKeySource(prefix []byte, seq uint64) [16]byte {
	// similar key building as libocr transmit selector
	hash := sha3.NewLegacyKeccak256()
	hash.Write(prefix[:])
	temp := make([]byte, 8)
	binary.LittleEndian.PutUint64(temp, seq)
	hash.Write(temp)

	var keyRandSource [16]byte
	copy(keyRandSource[:], hash.Sum(nil))
	return keyRandSource
}

type keyedCryptoRandSource struct {
	stream cipher.Stream
}

func NewKeyedCryptoRandSource(key [16]byte) rand.Source {
	var iv [16]byte // zero IV is fine here
	block, err := newCipherFn(key[:])
	if err != nil {
		// assertion
		panic(err)
	}
	return &keyedCryptoRandSource{cipher.NewCTR(block, iv[:])}
}

const int63Mask = 1<<63 - 1

func (crs *keyedCryptoRandSource) Int63() int64 {
	var buf [8]byte
	crs.stream.XORKeyStream(buf[:], buf[:])
	return int64(binary.LittleEndian.Uint64(buf[:]) & int63Mask)
}

func (crs *keyedCryptoRandSource) Seed(seed int64) {
	panic("keyedCryptoRandSource.Seed: Not supported")
}

type cryptoRandSource struct{}

func NewCryptoRandSource() rand.Source {
	return cryptoRandSource{}
}

func (_ cryptoRandSource) Int63() int64 {
	var b [8]byte
	_, err := randReadFn(b[:])
	if err != nil {
		panic(err)
	}
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

func (_ cryptoRandSource) Seed(_ int64) {
	panic("cryptoRandSource.Seed: Not supported")
}
