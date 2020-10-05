package permutation

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"math/rand"
)

func Permutation(n int, key [16]byte) []int {
	var result []int
	for i := 0; i < n; i++ {
		result = append(result, i)
	}

	r := rand.New(newCryptoRandSource(key))
	r.Shuffle(n, func(i int, j int) {
		result[i], result[j] = result[j], result[i]
	})
	return result
}

type cryptoRandSource struct {
	stream cipher.Stream
}

func newCryptoRandSource(key [16]byte) *cryptoRandSource {
	var iv [16]byte 	block, err := aes.NewCipher(key[:])
	if err != nil {
				panic(err)
	}
	return &cryptoRandSource{cipher.NewCTR(block, iv[:])}
}

const int63Mask = 1<<63 - 1

func (crs *cryptoRandSource) Int63() int64 {
	var buf [8]byte
	crs.stream.XORKeyStream(buf[:], buf[:])
	return int64(binary.LittleEndian.Uint64(buf[:]) & int63Mask)
}

func (crs *cryptoRandSource) Seed(seed int64) {
	panic("cryptoRandSource.Seed: Not supported")
}
