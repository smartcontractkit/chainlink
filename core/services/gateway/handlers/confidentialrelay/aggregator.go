package confidentialrelay

import (
	"encoding/json"
	"errors"
	"strconv"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	errInsufficientResponsesForQuorum = errors.New("insufficient valid responses to reach quorum")
	errQuorumUnobtainable             = errors.New("quorum unobtainable")
)

type aggregator struct{}

func (a *aggregator) Aggregate(resps map[string]jsonrpc.Response[json.RawMessage], donF int, donMembersCount int, l logger.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	// F+1 is sufficient: each honest node independently validates the enclave's
	// Nitro attestation, so F+1 matching responses guarantees at least one
	// honest node vouched for the result.
	requiredQuorum := donF + 1

	if len(resps) < requiredQuorum {
		return nil, errInsufficientResponsesForQuorum
	}

	shaToCount := map[string]int{}
	maxShaToCount := 0
	for _, r := range resps {
		sha, err := r.Digest()
		if err != nil {
			l.Errorw("failed to compute digest of response during quorum validation, skipping...", "error", err)
			continue
		}
		shaToCount[sha]++
		if shaToCount[sha] > maxShaToCount {
			maxShaToCount = shaToCount[sha]
		}
		if shaToCount[sha] >= requiredQuorum {
			return &r, nil
		}
	}

	remainingResponses := donMembersCount - len(resps)
	if maxShaToCount+remainingResponses < requiredQuorum {
		l.Warnw("quorum unattainable for request", "requiredQuorum", requiredQuorum, "remainingResponses", remainingResponses, "maxShaToCount", maxShaToCount)
		return nil, errors.New(errQuorumUnobtainable.Error() + ". RequiredQuorum=" + strconv.Itoa(requiredQuorum) + ". maxShaToCount=" + strconv.Itoa(maxShaToCount) + " remainingResponses=" + strconv.Itoa(remainingResponses))
	}

	return nil, errInsufficientResponsesForQuorum
}
