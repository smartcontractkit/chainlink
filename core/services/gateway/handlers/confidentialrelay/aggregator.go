package confidentialrelay

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"slices"

	jsonrpc "github.com/smartcontractkit/chainlink-common/pkg/jsonrpc2"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var (
	errInsufficientResponsesForQuorum = errors.New("insufficient valid responses to reach quorum")
	errQuorumUnobtainable             = errors.New("quorum unobtainable")
)

type aggregator struct{}

func (a *aggregator) Aggregate(resps map[string]jsonrpc.Response[json.RawMessage], donF int, donMembersCount int, l logger.Logger) (*jsonrpc.Response[json.RawMessage], error) {
	// F+1 (QuorumFPlusOne) is sufficient because each relay node calls the
	// target DON (Vault or capability) through CRE's standard capability
	// dispatch, which includes DON-level consensus. Every honest relay node
	// receives the same consensus-aggregated response and performs deterministic
	// translation, producing byte-identical outputs. F+1 matching responses
	// therefore guarantees at least one honest node vouched for the result.
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
	}

	var qualifiedDigests []string
	for sha, n := range shaToCount {
		if n >= requiredQuorum {
			qualifiedDigests = append(qualifiedDigests, sha)
		}
	}
	if len(qualifiedDigests) > 0 {
		slices.Sort(qualifiedDigests)
		want := qualifiedDigests[0]
		for _, k := range slices.Sorted(maps.Keys(resps)) {
			r := resps[k]
			sha, err := r.Digest()
			if err != nil {
				continue
			}
			if sha == want {
				out := r
				return &out, nil
			}
		}
	}

	remainingResponses := donMembersCount - len(resps)
	if maxShaToCount+remainingResponses < requiredQuorum {
		l.Warnw("quorum unattainable for request", "requiredQuorum", requiredQuorum, "remainingResponses", remainingResponses, "maxShaToCount", maxShaToCount)
		return nil, fmt.Errorf("%w: requiredQuorum=%d, maxShaToCount=%d, remainingResponses=%d", errQuorumUnobtainable, requiredQuorum, maxShaToCount, remainingResponses)
	}

	return nil, errInsufficientResponsesForQuorum
}
