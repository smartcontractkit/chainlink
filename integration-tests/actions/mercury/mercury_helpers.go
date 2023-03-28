package mercury

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/integration-tests/testsetups/mercury"
)

func GenFeedIds(count int) [][32]byte {
	var feedIds [][32]byte
	for i := 0; i < count; i++ {
		feedIds = append(feedIds, mercury.StringToByte32(fmt.Sprintf("feed-%d", i)))
	}
	return feedIds
}
