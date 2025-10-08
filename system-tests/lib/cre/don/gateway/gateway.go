package gateway

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/pelletier/go-toml/v2"

	"github.com/smartcontractkit/chainlink-deployments-framework/offchain/jd"
	jobv1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/job"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/don/jobs"
	"github.com/smartcontractkit/chainlink/v2/core/services/gateway/config"
)

func CreateJobs(ctx context.Context, jd *jd.JobDistributor, donTopology *cre.DonTopology, gatewayConfigs map[cre.NodeUUID]config.GatewayConfig) error {
	jobSpecs := make(cre.DonJobs, 0)

	header := `
type = "gateway"
schemaVersion = 1
externalJobID = "%s"
name = "cre-gateway"
forwardingAllowed = false
`

	for nodeUUID, gc := range gatewayConfigs {
		jobSpec := fmt.Sprintf(header, uuid.NewString())

		type wrapper struct {
			GC config.GatewayConfig `json:"gatewayConfig"  toml:"gatewayConfig"`
		}

		gatewayNode, found := donTopology.Dons.NodeWithUUID(nodeUUID)
		if !found {
			return fmt.Errorf("could not find gateway node with UUID %s in DON topology", nodeUUID)
		}

		tomlStr, mErr := toml.Marshal(wrapper{GC: gc})
		if mErr != nil {
			return fmt.Errorf("failed to marshal gateway config to toml: %w", mErr)
		}

		// hack for json.RawMessage that otherwise outputs a byte array instead of JSON string in toml, which cannot be parsed by gateway
		replaced, rErr := expandConfigByteArray(string(tomlStr), []string{"gatewayConfig", "Dons", "Handlers", "Config"})
		if rErr != nil {
			return fmt.Errorf("failed to expand config byte arrays: %w", rErr)
		}

		jobSpec += "\n" + replaced
		jobSpecs = append(jobSpecs, &jobv1.ProposeJobRequest{
			NodeId: gatewayNode.JobDistributorDetails.NodeID,
			Spec:   jobSpec,
		})
	}

	return jobs.Create(ctx, jd, donTopology, jobSpecs)
}

// ExpandConfigByteArray finds lines like `Config = [10, 109, ...]` and replaces them
// with TOML tables under the given path, using the bytes as TOML text.
// Example path: []string{"gatewayConfig","Dons","Handlers","Config"}
func expandConfigByteArray(tomlDoc string, path []string) (string, error) {
	re := regexp.MustCompile(`(?m)^(\s*)Config\s*=\s*\[([0-9,\s]+)\]\s*$`)
	return re.ReplaceAllStringFunc(tomlDoc, func(line string) string {
		m := re.FindStringSubmatch(line)
		if m == nil {
			return line
		}
		indent := m[1]
		nums := m[2]

		// parse the byte array back to text
		bs, err := parseByteArray(nums)
		if err != nil {
			// if parsing fails, keep original line to avoid breaking output
			return line
		}
		snippet := string(bs)

		// rewrite snippet under full path
		expanded := embedUnderPath(snippet, path)

		// keep the same indentation before each emitted line
		var out strings.Builder
		for _, l := range strings.Split(expanded, "\n") {
			if len(strings.TrimSpace(l)) == 0 {
				out.WriteString("\n")
				continue
			}
			out.WriteString(indent)
			out.WriteString(l)
			out.WriteString("\n")
		}
		return strings.TrimRight(out.String(), "\n")
	}), nil
}

func parseByteArray(s string) ([]byte, error) {
	fields := strings.Split(s, ",")
	buf := bytes.NewBuffer(nil)
	for _, f := range fields {
		t := strings.TrimSpace(f)
		if t == "" {
			continue
		}
		n, err := strconv.Atoi(t)
		if err != nil || n < 0 || n > 255 {
			return nil, fmt.Errorf("invalid byte: %q", t)
		}
		buf.WriteByte(byte(n))
	}
	return buf.Bytes(), nil
}

// embedUnderPath prefixes any table headers in snippet with path,
// and adds a `[path...]` header for the root keys.
func embedUnderPath(snippet string, path []string) string {
	base := "[" + strings.Join(path, ".") + "]"
	var out strings.Builder
	out.WriteString(base)
	out.WriteString("\n")

	for _, raw := range strings.Split(strings.ReplaceAll(snippet, "\r\n", "\n"), "\n") {
		line := strings.TrimRight(raw, " \t")
		trim := strings.TrimSpace(line)
		if trim == "" || strings.HasPrefix(trim, "#") {
			out.WriteString(line)
			out.WriteString("\n")
			continue
		}
		// table header inside snippet? e.g. [NodeRateLimiter]
		if strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			inner := strings.TrimSuffix(strings.TrimPrefix(trim, "["), "]")
			out.WriteString("[" + strings.Join(path, ".") + "." + inner + "]\n")
			continue
		}
		// regular key/value line
		out.WriteString(line)
		out.WriteString("\n")
	}
	return strings.TrimRight(out.String(), "\n")
}
