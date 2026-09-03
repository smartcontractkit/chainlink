package confidentialrelay

import (
	"strconv"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"

	confidentialrelaytypes "github.com/smartcontractkit/chainlink-common/pkg/capabilities/v2/actions/confidentialrelay"
)

const capabilityCallDomain = "call-capability"
const secretsGetDomain = "get-secrets"

// capExecKey is the deterministic cache key for a capability-exec request,
// built from its logical identity: the (workflow, execution, step, capability)
// tuple the relay-DON signature binds to. Avoids hashing: the fields are
// required non-empty by Validate, so a plain join is stable and debuggable.
func capExecKey(p confidentialrelaytypes.CapabilityRequestParams) string {
	return strings.Join([]string{capabilityCallDomain, p.WorkflowID, p.ExecutionID, p.ReferenceID, p.CapabilityID}, "/")
}

// secretsKey is the deterministic cache key for a secrets-get request, built
// from its logical identity: workflow, execution, callback id.
func secretsKey(p confidentialrelaytypes.SecretsRequestParams) string {
	return strings.Join([]string{secretsGetDomain, p.WorkflowID, p.ExecutionID, strconv.Itoa(int(p.CallbackID))}, "/")
}

// responseMemo caches a completed relay-DON signed result so the enclave's
// retry loop (re-fan-out after a gateway rotation or timeout) returns the
// cached result without re-executing the capability or vault fetch. Only
// completed results are cached; in-flight calls are not deduplicated, so a
// burst of concurrent same-identity requests each execute (preserving parallel
// dispatch).
//
// It uses patrickmn/go-cache, the string-keyed TTL cache already used
// elsewhere in chainlink (core/services/ocr2/plugins/promwrapper). The cap-
// and sec-prefixed keys share one cache instance. go-cache runs its own
// background cleanup goroutine, so there is no separate sweep to stop.
type responseMemo struct {
	c *cache.Cache
}

// newResponseMemo builds the memo with the configured TTL and cleanup interval
// (ConfidentialCompute.RelayResponseCache in cresettings).
func newResponseMemo(ttl, cleanupInterval time.Duration) *responseMemo {
	return &responseMemo{c: cache.New(ttl, cleanupInterval)}
}

func (m *responseMemo) getCap(key string) any {
	if v, ok := m.c.Get(key); ok {
		return v
	}
	return nil
}

func (m *responseMemo) putCap(key string, signed any) {
	m.c.SetDefault(key, signed)
}

func (m *responseMemo) getSec(key string) any {
	if v, ok := m.c.Get(key); ok {
		return v
	}
	return nil
}

func (m *responseMemo) putSec(key string, signed any) {
	m.c.SetDefault(key, signed)
}
