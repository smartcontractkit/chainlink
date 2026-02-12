package topologyviz

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/system-tests/lib/cre"
	envconfig "github.com/smartcontractkit/chainlink/system-tests/lib/cre/environment/config"
)

const (
	ClassSingleDON = "single-don"
	ClassMultiDON  = "multi-don"
	ClassSharded   = "sharded"
)

type CapabilityPlacement struct {
	RawFlag    string  `json:"raw_flag"`
	BaseFlag   string  `json:"base_flag"`
	ChainID    *uint64 `json:"chain_id,omitempty"`
	RemoteFrom bool    `json:"remote_from_workflow"`
}

type DONSummary struct {
	Name                      string                `json:"name"`
	DONTypes                  []string              `json:"don_types"`
	NodeCount                 int                   `json:"node_count"`
	NodeRoles                 []string              `json:"node_roles"`
	ExposesRemoteCapabilities bool                  `json:"exposes_remote_capabilities"`
	ShardIndex                uint                  `json:"shard_index,omitempty"`
	SupportedEVMChains        []uint64              `json:"supported_evm_chains,omitempty"`
	SupportedSolChains        []string              `json:"supported_sol_chains,omitempty"`
	Capabilities              []CapabilityPlacement `json:"capabilities,omitempty"`
}

type TopologySummary struct {
	ConfigRef string       `json:"config_ref"`
	Topology  string       `json:"topology"`
	InfraType string       `json:"infra_type"`
	DONs      []DONSummary `json:"dons"`
}

type Artifacts struct {
	ASCIIPath    string
	MarkdownPath string
}

type capabilityMatrixRow struct {
	Capability string
	ByDON      map[string]string
}

func BuildSummary(cfg *envconfig.Config, configRef string) (*TopologySummary, error) {
	if cfg == nil {
		return nil, errors.New("config must not be nil")
	}

	dons := make([]DONSummary, 0, len(cfg.NodeSets))
	for _, ns := range cfg.NodeSets {
		if ns == nil {
			continue
		}
		dons = append(dons, summarizeDON(ns))
	}
	sort.Slice(dons, func(i, j int) bool { return strings.ToLower(dons[i].Name) < strings.ToLower(dons[j].Name) })

	infraType := "unknown"
	if cfg.Infra != nil {
		infraType = cfg.Infra.Type
	}
	topologyClass := classifyTopology(dons)

	return &TopologySummary{
		ConfigRef: configRef,
		Topology:  topologyClass,
		InfraType: infraType,
		DONs:      dons,
	}, nil
}

func WriteArtifacts(summary *TopologySummary, outputDir string) (*Artifacts, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, errors.Wrap(err, "failed to create topology output directory")
	}

	asciiPath := filepath.Join(outputDir, "topology-summary.txt")
	mdPath := filepath.Join(outputDir, "topology.md")
	jsonPath := filepath.Join(outputDir, "topology.json")

	if err := os.WriteFile(asciiPath, []byte(RenderASCII(summary)), 0o644); err != nil { //nolint:gosec // we want the documentation to be readable by everyone
		return nil, errors.Wrap(err, "failed to write topology ascii summary")
	}
	if err := os.WriteFile(mdPath, []byte(RenderMarkdown(summary)), 0o644); err != nil { //nolint:gosec // we want the documentation to be readable by everyone
		return nil, errors.Wrap(err, "failed to write topology markdown")
	}

	// Remove legacy JSON artifact if present from previous runs.
	if _, err := os.Stat(jsonPath); err == nil {
		_ = os.Remove(jsonPath)
	}

	return &Artifacts{ASCIIPath: asciiPath, MarkdownPath: mdPath}, nil
}

func RenderASCII(summary *TopologySummary) string {
	var b strings.Builder
	b.WriteString("\nDON TOPOLOGY OVERVIEW\n")
	b.WriteString(fmt.Sprintf("Config: %s\n", summary.ConfigRef))
	b.WriteString(fmt.Sprintf("Class: %s\n", summary.Topology))
	b.WriteString(fmt.Sprintf("Infra: %s\n", summary.InfraType))
	b.WriteString("\n")

	b.WriteString(RenderASCIIDONTable(summary))
	b.WriteString("\n")
	b.WriteString(RenderASCIICapabilityMatrix(summary))

	return b.String()
}

func RenderASCIIStartSummary(summary *TopologySummary) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Topology: %s (%s)\n", summary.ConfigRef, summary.Topology))
	b.WriteString(RenderASCIICapabilityMatrix(summary))
	return b.String()
}

func RenderASCIICapabilityMatrix(summary *TopologySummary) string {
	rows := buildCapabilityMatrix(summary.DONs)
	if len(rows) == 0 {
		return "Capability Matrix: no capabilities declared\n"
	}

	headers := make([]string, 0, len(summary.DONs)+1)
	headers = append(headers, "Capability")
	for _, don := range summary.DONs {
		headers = append(headers, don.Name+" DON")
	}

	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		if len(row.Capability) > widths[0] {
			widths[0] = len(row.Capability)
		}
		for idx, don := range summary.DONs {
			if l := len(row.ByDON[don.Name]); l > widths[idx+1] {
				widths[idx+1] = l
			}
		}
	}

	var b strings.Builder
	b.WriteString("Capability Matrix\n")
	b.WriteString(buildBorder(widths))
	b.WriteString(buildRow(headers, widths))
	b.WriteString(buildBorder(widths))
	for _, row := range rows {
		values := []string{row.Capability}
		for _, don := range summary.DONs {
			values = append(values, row.ByDON[don.Name])
		}
		b.WriteString(buildRow(values, widths))
	}
	b.WriteString(buildBorder(widths))
	return b.String()
}

func RenderASCIIDONTable(summary *TopologySummary) string {
	headers := []string{"DON", "Types", "Nodes", "EVM Chains", "Attributes"}
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}

	rows := make([][]string, 0, len(summary.DONs))
	for _, don := range summary.DONs {
		attrs := make([]string, 0, 2)
		if don.ExposesRemoteCapabilities {
			attrs = append(attrs, "remote-capabilities")
		}
		attributesText := "-"
		if len(attrs) > 0 {
			attributesText = strings.Join(attrs, ",")
		}

		evmText := "-"
		if len(don.SupportedEVMChains) > 0 {
			evmText = joinUint64(don.SupportedEVMChains)
		}

		row := []string{
			don.Name,
			strings.Join(don.DONTypes, ","),
			strconv.Itoa(don.NodeCount),
			evmText,
			attributesText,
		}
		for i, val := range row {
			if len(val) > widths[i] {
				widths[i] = len(val)
			}
		}
		rows = append(rows, row)
	}

	var b strings.Builder
	b.WriteString("DONs\n")
	b.WriteString(buildBorder(widths))
	b.WriteString(buildRow(headers, widths))
	b.WriteString(buildBorder(widths))
	for _, row := range rows {
		b.WriteString(buildRow(row, widths))
	}
	b.WriteString(buildBorder(widths))
	return b.String()
}

func RenderMarkdown(summary *TopologySummary) string {
	var b strings.Builder
	b.WriteString("# DON Topology\n\n")
	b.WriteString(fmt.Sprintf("- Config: `%s`\n", summary.ConfigRef))
	b.WriteString(fmt.Sprintf("- Class: `%s`\n", summary.Topology))
	b.WriteString(fmt.Sprintf("- Infra: `%s`\n", summary.InfraType))
	b.WriteString("\n")

	b.WriteString("## Capability Matrix\n\n")
	b.WriteString("This matrix is the source of truth for capability placement by DON.\n\n")
	b.WriteString("| Capability |")
	for _, don := range summary.DONs {
		b.WriteString(fmt.Sprintf(" `%s` |", don.Name))
	}
	b.WriteString("\n|---|")
	for range summary.DONs {
		b.WriteString("---|")
	}
	b.WriteString("\n")

	rows := buildCapabilityMatrix(summary.DONs)
	for _, row := range rows {
		b.WriteString(fmt.Sprintf("| `%s` |", row.Capability))
		for _, don := range summary.DONs {
			b.WriteString(fmt.Sprintf(" `%s` |", row.ByDON[don.Name]))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")

	b.WriteString("## DONs\n\n")
	for _, don := range summary.DONs {
		b.WriteString(fmt.Sprintf("### `%s`\n\n", don.Name))
		b.WriteString(fmt.Sprintf("- Types: `%s`\n", strings.Join(don.DONTypes, "`, `")))
		b.WriteString(fmt.Sprintf("- Nodes: `%d`\n", don.NodeCount))
		if len(don.NodeRoles) > 0 {
			b.WriteString(fmt.Sprintf("- Roles: `%s`\n", strings.Join(don.NodeRoles, "`, `")))
		}
		if len(don.SupportedEVMChains) > 0 {
			b.WriteString(fmt.Sprintf("- EVM chains: `%s`\n", joinUint64(don.SupportedEVMChains)))
		}
		if len(don.SupportedSolChains) > 0 {
			b.WriteString(fmt.Sprintf("- Solana chains: `%s`\n", strings.Join(don.SupportedSolChains, "`, `")))
		}
		b.WriteString(fmt.Sprintf("- Exposes remote capabilities: `%t`\n", don.ExposesRemoteCapabilities))
		b.WriteString("\n")
	}

	return b.String()
}

func summarizeDON(ns *cre.NodeSet) DONSummary {
	rolesSet := make(map[string]struct{})
	for _, spec := range ns.NodeSpecs {
		if spec == nil {
			continue
		}
		for _, role := range spec.Roles {
			rolesSet[role] = struct{}{}
		}
	}

	donTypes := append([]string{}, ns.DONTypes...)
	sort.Strings(donTypes)
	evmChains := append([]uint64{}, ns.EVMChains()...)
	sort.Slice(evmChains, func(i, j int) bool { return evmChains[i] < evmChains[j] })
	solChains := append([]string{}, ns.SupportedSolChains...)
	sort.Strings(solChains)

	nodeCount := len(ns.NodeSpecs)
	if nodeCount == 0 {
		nodeCount = ns.Nodes
	}

	return DONSummary{
		Name:                      ns.Name,
		DONTypes:                  donTypes,
		NodeCount:                 nodeCount,
		NodeRoles:                 sortedKeys(rolesSet),
		ExposesRemoteCapabilities: ns.ExposesRemoteCapabilities,
		ShardIndex:                ns.ShardIndex,
		SupportedEVMChains:        evmChains,
		SupportedSolChains:        solChains,
		Capabilities:              summarizeCapabilities(ns),
	}
}

func summarizeCapabilities(ns *cre.NodeSet) []CapabilityPlacement {
	out := make([]CapabilityPlacement, 0, len(ns.Capabilities))
	for _, capFlag := range ns.Capabilities {
		base, chain := splitCapabilityFlag(capFlag)
		out = append(out, CapabilityPlacement{
			RawFlag:    capFlag,
			BaseFlag:   base,
			ChainID:    chain,
			RemoteFrom: ns.ExposesRemoteCapabilities && hasDONType(ns.DONTypes, cre.CapabilitiesDON),
		})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RawFlag < out[j].RawFlag })
	return out
}

func buildCapabilityMatrix(dons []DONSummary) []capabilityMatrixRow {
	type cell struct {
		remote bool
		chains []string
	}
	matrix := make(map[string]map[string]*cell)
	for _, don := range dons {
		for _, capInfo := range don.Capabilities {
			if _, exists := matrix[capInfo.BaseFlag]; !exists {
				matrix[capInfo.BaseFlag] = make(map[string]*cell)
			}
			if _, exists := matrix[capInfo.BaseFlag][don.Name]; !exists {
				matrix[capInfo.BaseFlag][don.Name] = &cell{}
			}
			c := matrix[capInfo.BaseFlag][don.Name]
			c.remote = c.remote || capInfo.RemoteFrom
			if capInfo.ChainID != nil {
				c.chains = append(c.chains, strconv.FormatUint(*capInfo.ChainID, 10))
			}
		}
	}

	caps := make([]string, 0, len(matrix))
	for cap := range matrix {
		caps = append(caps, cap)
	}
	sort.Strings(caps)

	rows := make([]capabilityMatrixRow, 0, len(caps))
	for _, cap := range caps {
		byDON := make(map[string]string, len(dons))
		for _, don := range dons {
			c, ok := matrix[cap][don.Name]
			if !ok {
				byDON[don.Name] = "-"
				continue
			}
			mode := "local"
			if c.remote {
				mode = "remote-exposed"
			}
			if len(c.chains) > 0 {
				sort.Strings(c.chains)
				mode = mode + " (" + strings.Join(uniqueStrings(c.chains), ",") + ")"
			}
			byDON[don.Name] = mode
		}
		rows = append(rows, capabilityMatrixRow{Capability: cap, ByDON: byDON})
	}

	return rows
}

func buildBorder(widths []int) string {
	var b strings.Builder
	b.WriteString("+")
	for _, w := range widths {
		b.WriteString(strings.Repeat("-", w+2))
		b.WriteString("+")
	}
	b.WriteString("\n")
	return b.String()
}

func buildRow(values []string, widths []int) string {
	var b strings.Builder
	b.WriteString("|")
	for i, v := range values {
		b.WriteString(" ")
		b.WriteString(padRight(v, widths[i]))
		b.WriteString(" |")
	}
	b.WriteString("\n")
	return b.String()
}

func padRight(v string, width int) string {
	if len(v) >= width {
		return v
	}
	return v + strings.Repeat(" ", width-len(v))
}

func classifyTopology(dons []DONSummary) string {
	hasShards := false
	capDONs := 0
	workflowDONs := 0
	for _, d := range dons {
		if d.ShardIndex > 0 || stringSliceContains(d.DONTypes, cre.ShardDON) {
			hasShards = true
		}
		if stringSliceContains(d.DONTypes, cre.CapabilitiesDON) {
			capDONs++
		}
		if stringSliceContains(d.DONTypes, cre.WorkflowDON) {
			workflowDONs++
		}
	}
	if hasShards {
		return ClassSharded
	}
	if capDONs > 0 || workflowDONs > 1 || len(dons) > 2 {
		return ClassMultiDON
	}
	return ClassSingleDON
}

func splitCapabilityFlag(flag string) (string, *uint64) {
	lastDash := strings.LastIndex(flag, "-")
	if lastDash < 0 || lastDash == len(flag)-1 {
		return flag, nil
	}
	maybeID := flag[lastDash+1:]
	chainID, err := strconv.ParseUint(maybeID, 10, 64)
	if err != nil {
		return flag, nil
	}
	return flag[:lastDash], &chainID
}

func sortedKeys(set map[string]struct{}) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func joinUint64(v []uint64) string {
	out := make([]string, 0, len(v))
	for _, n := range v {
		out = append(out, strconv.FormatUint(n, 10))
	}
	return strings.Join(out, ",")
}

func uniqueStrings(v []string) []string {
	if len(v) == 0 {
		return nil
	}
	out := []string{v[0]}
	for i := 1; i < len(v); i++ {
		if v[i] != v[i-1] {
			out = append(out, v[i])
		}
	}
	return out
}

func stringSliceContains(values []string, wanted string) bool {
	for _, v := range values {
		if strings.EqualFold(v, wanted) {
			return true
		}
	}
	return false
}

func hasDONType(donTypes []string, target string) bool {
	for _, d := range donTypes {
		if strings.EqualFold(d, target) {
			return true
		}
	}
	return false
}
