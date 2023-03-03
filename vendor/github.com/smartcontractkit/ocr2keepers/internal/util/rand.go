package util

import (
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"encoding/binary"
	"math/rand"
)

type keyedCryptoRandSource struct {
	stream cipher.Stream
}

func NewKeyedCryptoRandSource(key [16]byte) rand.Source {
	var iv [16]byte // zero IV is fine here
	block, err := aes.NewCipher(key[:])
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
	_, err := crand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return int64(binary.LittleEndian.Uint64(b[:]) & (1<<63 - 1))
}

func (_ cryptoRandSource) Seed(_ int64) {
	panic("cryptoRandSource.Seed: Not supported")
}
