package matrix

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// TopologyConfig defines the topology name and configuration file path for a test run.
type TopologyConfig struct {
	Topology string `json:"topology"`
	Configs  string `json:"configs"`
}

var defaultTopology = TopologyConfig{
	Topology: "workflow-gateway-capabilities",
	Configs:  "configs/workflow-gateway-capabilities-don.toml",
}

var perTestTopologies = map[string]TopologyConfig{
	"TestCRE_V2_Aptos_Suite_E2E": {
		Topology: "workflow-gateway-aptos",
		Configs:  "configs/workflow-gateway-don-aptos.toml",
	},
	"TestCRE_V2_Stellar_Suite_E2E": {
		Topology: "workflow-gateway-stellar",
		Configs:  "configs/workflow-gateway-don-stellar.toml",
	},
	"TestCRE_V2_Solana_Write_E2E": {
		Topology: "workflow",
		Configs:  "configs/workflow-don-solana.toml",
	},
	"TestCRE_V2_Solana_LogTrigger_E2E": {
		Topology: "workflow",
		Configs:  "configs/workflow-don-solana.toml",
	},
	"TestCRE_V2_Solana_Read_Accounts_E2E": {
		Topology: "workflow",
		Configs:  "configs/workflow-don-solana.toml",
	},
	"TestCRE_V2_Solana_Read_Block_E2E": {
		Topology: "workflow",
		Configs:  "configs/workflow-don-solana.toml",
	},
	"TestCRE_V2_Solana_Read_Tx_E2E": {
		Topology: "workflow",
		Configs:  "configs/workflow-don-solana.toml",
	},
	"TestCRE_V2_Sharding_E2E": {
		Topology: "workflow-gateway-sharded",
		Configs:  "configs/workflow-gateway-sharded-don.toml",
	},
	"TestCRE_V2_ShardManualAssignment_E2E": {
		Topology: "workflow-gateway-sharded-manual",
		Configs:  "configs/workflow-gateway-sharded-manual.toml",
	},
	"TestCRE_V2_ShardRingOCROverrides_E2E": {
		Topology: "workflow-gateway-sharded-ringocr-overrides",
		Configs:  "configs/workflow-gateway-sharded-ringocr-overrides.toml",
	},
	"TestCRE_V2_Module_Cache_E2E": {
		Topology: "workflow-gateway-cache-test",
		Configs:  "configs/workflow-gateway-don-cache-test.toml",
	},
	"TestCRE_V2_HTTP_Action_Multi_Gateway_E2E": {
		Topology: "workflow-gateway-capabilities-multi-gateway",
		Configs:  "configs/workflow-gateway-capabilities-multi-gateway-don.toml",
	},
}

// GetTopologyForTest returns the topology and config file for a given test name.
func GetTopologyForTest(testName string) TopologyConfig {
	if cfg, ok := perTestTopologies[testName]; ok {
		return cfg
	}
	return defaultTopology
}

// Entry represents one test job item in the matrix JSON array.
type Entry struct {
	TestName string `json:"test_name"`
	TestID   string `json:"test_id"`
	RunsOn   string `json:"runs_on"`
	Topology string `json:"topology"`
	Configs  string `json:"configs"`
}

// ScanDir scans all Go test files in dir and discovers test/example function names matching rePattern.
func ScanDir(dir string, rePattern string) ([]string, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", rePattern, err)
	}

	resolvedDir := dir
	if _, statErr := os.Stat(resolvedDir); statErr != nil && !filepath.IsAbs(dir) {
		if alt := filepath.Join("../..", dir); func() bool { _, e := os.Stat(alt); return e == nil }() {
			resolvedDir = alt
		}
	}

	fset := token.NewFileSet()
	//nolint:staticcheck // parser.ParseDir is used for fast AST test discovery without heavy type-checking
	pkgs, err := parser.ParseDir(fset, resolvedDir, func(fi fs.FileInfo) bool {
		return strings.HasSuffix(fi.Name(), "_test.go")
	}, parser.SkipObjectResolution)
	if err != nil {
		return nil, fmt.Errorf("failed to parse directory %q: %w", resolvedDir, err)
	}

	seen := make(map[string]struct{})
	var testNames []string

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || fn.Recv != nil {
					continue
				}

				name := fn.Name.Name
				if re.MatchString(name) {
					if _, exists := seen[name]; !exists {
						seen[name] = struct{}{}
						testNames = append(testNames, name)
					}
				}
			}
		}
	}

	if len(testNames) == 0 {
		return nil, fmt.Errorf("no matching test functions found in %q matching pattern %q", dir, rePattern)
	}

	sort.Strings(testNames)
	return testNames, nil
}

// BuildMatrix converts testNames into a slice of Entry structs with unique runs-on specifiers.
func BuildMatrix(testNames []string, runID, runAttempt, runnerSpec string) []Entry {
	entries := make([]Entry, 0, len(testNames))
	for idx, name := range testNames {
		runsOn := fmt.Sprintf("runs-on=%s-%d-%s/%s", runID, idx, runAttempt, runnerSpec)
		topCfg := GetTopologyForTest(name)
		entries = append(entries, Entry{
			TestName: name,
			TestID:   name,
			RunsOn:   runsOn,
			Topology: topCfg.Topology,
			Configs:  topCfg.Configs,
		})
	}
	return entries
}

// WriteOutput formats the matrix entries as JSON and writes to w and optionally $GITHUB_OUTPUT.
func WriteOutput(w io.Writer, entries []Entry, writeGithubOutput bool) error {
	matrixJSON, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to marshal matrix JSON: %w", err)
	}

	if _, err := fmt.Fprintln(w, string(matrixJSON)); err != nil {
		return err
	}

	if writeGithubOutput {
		githubOutputFile := os.Getenv("GITHUB_OUTPUT")
		if githubOutputFile != "" {
			f, err := os.OpenFile(filepath.Clean(githubOutputFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
			if err != nil {
				return fmt.Errorf("failed to open GITHUB_OUTPUT file %s: %w", githubOutputFile, err)
			}
			defer f.Close()

			if _, err := fmt.Fprintf(f, "matrix=%s\n", string(matrixJSON)); err != nil {
				return fmt.Errorf("failed to write to GITHUB_OUTPUT: %w", err)
			}
		}
	}

	return nil
}
