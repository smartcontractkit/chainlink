package hash

import (
	hashAlg "crypto/sha256"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const Size = hashAlg.Size

type Hash [Size]byte

func GetHash(s []byte) Hash {
	return hashAlg.Sum256(s)
}

var Zero Hash

type Hashes []Hash

func MakeHashes() Hashes { return Hashes{} }

func (hs *Hashes) Add(h Hash) {
	*hs = append(*hs, h)
}

func (hs Hashes) String() string {
	shs := make([]string, len(hs))
	for i, h := range hs {
		shs[i] = h.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(shs, ", "))
}

func (h Hash) String() string {
	return hexutil.Encode(h[:])
}
