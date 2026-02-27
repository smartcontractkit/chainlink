package mercury

import (
	"github.com/ethereum/go-ethereum/common/hexutil"

	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

func MustHexToConfigDigest(s string) (cd ocrtypes.ConfigDigest) {
	b := hexutil.MustDecode(s)
	var err error
	cd, err = ocrtypes.BytesToConfigDigest(b)
	if err != nil {
		panic(err)
	}
	return
}
