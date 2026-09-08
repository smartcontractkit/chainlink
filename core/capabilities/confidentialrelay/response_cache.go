package confidentialrelay

import (
	"strconv"
	"strings"

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
