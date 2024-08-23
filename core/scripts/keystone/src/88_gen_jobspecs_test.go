package src

import (
	"fmt"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func (d *donHostSpec) ToString() string {
	var result string
	result += "Bootstrap:\n"
	result += "Host: " + d.bootstrap.host + "\n"
	result += d.bootstrap.spec.ToString()
	result += "\n\nOracles:\n"
	for i, oracle := range d.oracles {
		if i != 0 {
			result += "--------------------------------\n"
		}
		result += fmt.Sprintf("Oracle %d:\n", i)
		result += "Host: " + oracle.host + "\n"
		result += oracle.spec.ToString()
		result += "\n\n"
	}
	return result
}

func TestGenSpecs(t *testing.T) {
	pubkeysPath := "./testdata/PublicKeys.json"
	nodeListPath := "./testdata/NodeList.txt"
	chainID := int64(11155111)
	p2pPort := int64(6690)
	contractAddress := "0xB29934624cAe3765E33115A9530a13f5aEC7fa8A"

	specs := genSpecs(pubkeysPath, nodeListPath, "../templates", chainID, p2pPort, contractAddress)
	snaps.MatchSnapshot(t, specs.ToString())
}
