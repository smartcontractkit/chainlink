package matrix

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/githuboutput"
	"github.com/smartcontractkit/chainlink/v2/tools/ci/internal/paths"
)

// SuiteType identifies predefined test matrix suites.
type SuiteType string

const (
	SuiteCRESmoke      SuiteType = "cre-smoke"
	SuiteCRERegression SuiteType = "cre-regression"
	SuiteCREMixedEnv   SuiteType = "cre-mixed-env"
	SuiteCCIP          SuiteType = "ccip"
)

const (
	DefaultCRESmokeDir          = "system-tests/tests/smoke/cre"
	DefaultCRESmokePattern      = `^TestCRE_.*_E2E$`
	DefaultCRERegressionDir     = "system-tests/tests/regression/cre"
	DefaultCRERegressionPattern = `^(Test|Example).*`
	DefaultCCIPDir              = "integration-tests/smoke/ccip"
	DefaultCCIPPattern          = `^TestCCIP_.*_E2E$`
	DefaultCRERunnerSpec        = "cpu=16/ram=64/family=m7i+m8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs"
	DefaultCCIPRunnerSpec       = "cpu=8/ram=64/family=r6i+r7i+r8i/spot=co/image=ubuntu24-full-x64/extras=s3-cache+tmpfs"
)

var perTestRegressionConfigs = map[string]string{
	"Test_CRE_V2_Stellar_Regression": "configs/workflow-gateway-don-stellar.toml",
}

const defaultRegressionConfig = "configs/workflow-gateway-capabilities-don.toml"

var mixedEnvTests = []string{
	"TestCRE_V2_Suite_Bucket_A_E2E",
	"TestCRE_V2_Suite_Bucket_B_E2E",
	"TestCRE_V2_EVM_Read_HeavyCalls_E2E",
	"TestCRE_V2_EVM_Read_StateQueries_E2E",
	"TestCRE_V2_EVM_Read_TxArtifacts_E2E",
}

// CCIPConfig holds per-test configuration overrides for CCIP smoke tests.
type CCIPConfig struct {
	Timeout             string
	SelectedNetwork     string
	RMNRageProxyVersion string
	RMNAFN2ProxyVersion string
	JobTimeout          int
}

var defaultCCIPConfig = CCIPConfig{
	Timeout:         "15m",
	SelectedNetwork: "SIMULATED_1,SIMULATED_2",
}

var perTestCCIPConfigs = map[string]CCIPConfig{
	"TestCCIP_RMN_GlobalCurseTwoMessagesOnTwoLanes_E2E": {
		Timeout:             "15m",
		SelectedNetwork:     "SIMULATED_1,SIMULATED_2",
		RMNRageProxyVersion: "master-amd6416f5d86",
		RMNAFN2ProxyVersion: "master-amd64-10b42b2",
	},
	"TestCCIP_JobSpecs_E2E": {
		Timeout:         "15m",
		SelectedNetwork: "SIMULATED_1,SIMULATED_2",
		JobTimeout:      20,
	},
}

// TopologyConfig defines the topology name and configuration file path for a test run.
type TopologyConfig struct {
	Topology string `json:"topology"`
	Configs  string `json:"configs"`
}

var defaultTopology = TopologyConfig{
	Topology: "workflow-gateway-capabilities",
	Configs:  "configs/workflow-gateway-capabilities-don.toml",
}

var perTestTopologies = map[string][]TopologyConfig{
	"TestCRE_V2_Suite_Bucket_B_E2E": {
		{Topology: "workflow-gateway-capabilities", Configs: "configs/workflow-gateway-capabilities-don.toml"},
		{Topology: "workflow-gateway-capabilities-vault-jwt_auth-enabled", Configs: "configs/workflow-gateway-capabilities-don-vault-jwt_auth-enabled.toml"},
		{Topology: "workflow-gateway-capabilities-vault-optimizations-enabled", Configs: "configs/workflow-gateway-capabilities-don-vault-optimizations-enabled.toml"},
		{Topology: "workflow-gateway-capabilities-vault-stall-purge", Configs: "configs/workflow-gateway-capabilities-don-vault-stall-purge.toml"},
	},
	"TestCRE_V2_Aptos_Suite_E2E": {
		{Topology: "workflow-gateway-aptos", Configs: "configs/workflow-gateway-don-aptos.toml"},
	},
	"TestCRE_V2_Stellar_Suite_E2E": {
		{Topology: "workflow-gateway-stellar", Configs: "configs/workflow-gateway-don-stellar.toml"},
	},
	"TestCRE_V2_Solana_Write_E2E": {
		{Topology: "workflow", Configs: "configs/workflow-don-solana.toml"},
	},
	"TestCRE_V2_Solana_LogTrigger_E2E": {
		{Topology: "workflow", Configs: "configs/workflow-don-solana.toml"},
	},
	"TestCRE_V2_Solana_Read_Accounts_E2E": {
		{Topology: "workflow", Configs: "configs/workflow-don-solana.toml"},
	},
	"TestCRE_V2_Solana_Read_Block_E2E": {
		{Topology: "workflow", Configs: "configs/workflow-don-solana.toml"},
	},
	"TestCRE_V2_Solana_Read_Tx_E2E": {
		{Topology: "workflow", Configs: "configs/workflow-don-solana.toml"},
	},
	"TestCRE_V2_Sharding_E2E": {
		{Topology: "workflow-gateway-sharded", Configs: "configs/workflow-gateway-sharded-don.toml"},
	},
	"TestCRE_V2_ShardManualAssignment_E2E": {
		{Topology: "workflow-gateway-sharded-manual", Configs: "configs/workflow-gateway-sharded-manual.toml"},
	},
	"TestCRE_V2_ShardRingOCROverrides_E2E": {
		{Topology: "workflow-gateway-sharded-ringocr-overrides", Configs: "configs/workflow-gateway-sharded-ringocr-overrides.toml"},
	},
	"TestCRE_V2_Module_Cache_E2E": {
		{Topology: "workflow-gateway-cache-test", Configs: "configs/workflow-gateway-don-cache-test.toml"},
	},
	"TestCRE_V2_HTTP_Action_Multi_Gateway_E2E": {
		{Topology: "workflow-gateway-capabilities-multi-gateway", Configs: "configs/workflow-gateway-capabilities-multi-gateway-don.toml"},
	},
}

// GetTopologiesForTest returns the topologies and config files for a given test name,
// falling back to the default topology when the test has no explicit entry.
func GetTopologiesForTest(testName string) []TopologyConfig {
	if cfgs, ok := perTestTopologies[testName]; ok {
		return cfgs
	}
	return []TopologyConfig{defaultTopology}
}

// Entry represents one test job item in the matrix JSON array.
type Entry struct {
	TestName            string `json:"test_name"`
	TestID              string `json:"test_id"`
	RunsOn              string `json:"runs_on"`
	Topology            string `json:"topology,omitempty"`
	Configs             string `json:"configs,omitempty"`
	Timeout             string `json:"timeout,omitempty"`
	SelectedNetwork     string `json:"selected_network,omitempty"`
	RMNRageProxyVersion string `json:"rmn_rageproxy_version,omitempty"`
	RMNAFN2ProxyVersion string `json:"rmn_afn2proxy_version,omitempty"`
	JobTimeout          int    `json:"job_timeout,omitempty"`
}

// SuiteOptions specifies parameters for building a test suite matrix.
type SuiteOptions struct {
	Dir        string
	Pattern    string
	RunID      string
	RunAttempt string
	RunnerSpec string
}

// SetupOptions specifies parameters for generating all integration test matrices in one go.
type SetupOptions struct {
	RunID            string
	RunAttempt       string
	CRESmoke         bool
	CRESmokeDir      string
	CRERegression    bool
	CRERegressionDir string
	CREMixedEnv      bool
	CCIP             bool
	CCIPDir          string
}

// BuildSuiteMatrix generates entries for a specific test suite.
func BuildSuiteMatrix(suite SuiteType, opts SuiteOptions) ([]Entry, error) {
	runID := opts.RunID
	if runID == "" {
		runID = "0"
	}
	attempt := opts.RunAttempt
	if attempt == "" {
		attempt = "1"
	}

	switch suite {
	case SuiteCRESmoke:
		dir := opts.Dir
		if dir == "" {
			dir = DefaultCRESmokeDir
		}
		pattern := opts.Pattern
		if pattern == "" {
			pattern = DefaultCRESmokePattern
		}
		runner := opts.RunnerSpec
		if runner == "" {
			runner = DefaultCRERunnerSpec
		}

		testNames, err := ScanDir(dir, pattern)
		if err != nil {
			return nil, err
		}
		return BuildMatrix(testNames, runID, attempt, runner), nil

	case SuiteCRERegression:
		dir := opts.Dir
		if dir == "" {
			dir = DefaultCRERegressionDir
		}
		pattern := opts.Pattern
		if pattern == "" {
			pattern = DefaultCRERegressionPattern
		}
		runner := opts.RunnerSpec
		if runner == "" {
			runner = DefaultCRERunnerSpec
		}

		testNames, err := ScanDir(dir, pattern)
		if err != nil {
			return nil, err
		}

		entries := make([]Entry, 0, len(testNames))
		for idx, name := range testNames {
			runsOn := fmt.Sprintf("runs-on=%s-%d-%s/%s", runID, idx, attempt, runner)
			cfg := defaultRegressionConfig
			if c, ok := perTestRegressionConfigs[name]; ok {
				cfg = c
			}
			entries = append(entries, Entry{
				TestName: name,
				TestID:   strconv.Itoa(idx),
				RunsOn:   runsOn,
				Configs:  cfg,
			})
		}
		return entries, nil

	case SuiteCREMixedEnv:
		runner := opts.RunnerSpec
		if runner == "" {
			runner = DefaultCRERunnerSpec
		}
		entries := make([]Entry, 0, len(mixedEnvTests))
		for idx, name := range mixedEnvTests {
			runsOn := fmt.Sprintf("runs-on=%s-%d-%s/%s", runID, idx, attempt, runner)
			entries = append(entries, Entry{
				TestName: name,
				TestID:   strconv.Itoa(idx),
				RunsOn:   runsOn,
			})
		}
		return entries, nil

	case SuiteCCIP:
		dir := opts.Dir
		if dir == "" {
			dir = DefaultCCIPDir
		}
		pattern := opts.Pattern
		if pattern == "" {
			pattern = DefaultCCIPPattern
		}
		runner := opts.RunnerSpec
		if runner == "" {
			runner = DefaultCCIPRunnerSpec
		}

		testNames, err := ScanDir(dir, pattern)
		if err != nil {
			return nil, err
		}

		entries := make([]Entry, 0, len(testNames))
		for idx, name := range testNames {
			runsOn := fmt.Sprintf("runs-on=%s-%d-%s/%s", runID, idx, attempt, runner)
			cfg := defaultCCIPConfig
			if c, ok := perTestCCIPConfigs[name]; ok {
				cfg = c
			}
			entries = append(entries, Entry{
				TestName:            name,
				TestID:              strconv.Itoa(idx),
				RunsOn:              runsOn,
				Timeout:             cfg.Timeout,
				SelectedNetwork:     cfg.SelectedNetwork,
				RMNRageProxyVersion: cfg.RMNRageProxyVersion,
				RMNAFN2ProxyVersion: cfg.RMNAFN2ProxyVersion,
				JobTimeout:          cfg.JobTimeout,
			})
		}
		return entries, nil

	default:
		return nil, fmt.Errorf("unknown suite type: %q", suite)
	}
}

// GenerateSetupMatrices generates all requested matrices for integration-tests gating/setup.
func GenerateSetupMatrices(opts SetupOptions) (map[string][]Entry, error) {
	result := make(map[string][]Entry)

	if opts.CRESmoke {
		entries, err := BuildSuiteMatrix(SuiteCRESmoke, SuiteOptions{
			Dir:        opts.CRESmokeDir,
			RunID:      opts.RunID,
			RunAttempt: opts.RunAttempt,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build cre-smoke matrix: %w", err)
		}
		result["cre-matrix"] = entries
	}

	if opts.CRERegression {
		entries, err := BuildSuiteMatrix(SuiteCRERegression, SuiteOptions{
			Dir:        opts.CRERegressionDir,
			RunID:      opts.RunID,
			RunAttempt: opts.RunAttempt,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build cre-regression matrix: %w", err)
		}
		result["cre-regression-matrix"] = entries
	}

	if opts.CREMixedEnv {
		entries, err := BuildSuiteMatrix(SuiteCREMixedEnv, SuiteOptions{
			RunID:      opts.RunID,
			RunAttempt: opts.RunAttempt,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build cre-mixed-env matrix: %w", err)
		}
		result["cre-mixed-env-matrix"] = entries
	}

	if opts.CCIP {
		entries, err := BuildSuiteMatrix(SuiteCCIP, SuiteOptions{
			Dir:        opts.CCIPDir,
			RunID:      opts.RunID,
			RunAttempt: opts.RunAttempt,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to build ccip matrix: %w", err)
		}
		result["ccip-matrix"] = entries
	}

	return result, nil
}

// ScanDir scans all Go test files in dir and discovers test/example function names matching rePattern.
func ScanDir(dir string, rePattern string) ([]string, error) {
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern %q: %w", rePattern, err)
	}

	resolvedDir := paths.ResolveFromRepoRoot(dir)

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
				if name == "TestMain" {
					continue
				}
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
// Tests with multiple registered topologies produce one entry per topology.
func BuildMatrix(testNames []string, runID, runAttempt, runnerSpec string) []Entry {
	var entries []Entry
	for _, name := range testNames {
		for _, topCfg := range GetTopologiesForTest(name) {
			idx := len(entries)
			entries = append(entries, Entry{
				TestName: name,
				TestID:   strconv.Itoa(idx),
				RunsOn:   fmt.Sprintf("runs-on=%s-%d-%s/%s", runID, idx, runAttempt, runnerSpec),
				Topology: topCfg.Topology,
				Configs:  topCfg.Configs,
			})
		}
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
		if err := githuboutput.AppendVar("matrix", string(matrixJSON)); err != nil {
			return fmt.Errorf("failed to write matrix to GITHUB_OUTPUT: %w", err)
		}
	}

	return nil
}

// WriteMultiOutput formats map of matrices as JSON and writes to w and optionally $GITHUB_OUTPUT.
func WriteMultiOutput(w io.Writer, matrices map[string][]Entry, writeGithubOutput bool) error {
	data, err := json.Marshal(matrices)
	if err != nil {
		return fmt.Errorf("failed to marshal matrices JSON: %w", err)
	}

	if _, err := fmt.Fprintln(w, string(data)); err != nil {
		return err
	}

	if writeGithubOutput {
		vars := make(map[string]string, len(matrices))
		for key, entries := range matrices {
			val, err := json.Marshal(entries)
			if err != nil {
				return fmt.Errorf("failed to marshal %s JSON: %w", key, err)
			}
			vars[key] = string(val)
		}
		if err := githuboutput.AppendVars(vars); err != nil {
			return fmt.Errorf("failed to write matrices to GITHUB_OUTPUT: %w", err)
		}
	}

	return nil
}
