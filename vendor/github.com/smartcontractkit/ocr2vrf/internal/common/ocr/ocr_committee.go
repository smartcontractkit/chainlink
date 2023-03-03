package ocr

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type OCRCommittee struct{ Signers, Transmitters []common.Address }

func (o OCRCommittee) String() string {
	sportions := make([]string, len(o.Signers))
	for i, s := range o.Signers {
		sportions[i] = s.Hex()
	}
	tportions := make([]string, len(o.Transmitters))
	for i, t := range o.Transmitters {
		tportions[i] = t.Hex()
	}
	return fmt.Sprintf("OCRCommittee{Signers: %s, Transmitters: %s}",
		strings.Join(sportions, ", "), strings.Join(tportions, ", "))
}
