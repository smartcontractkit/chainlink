package updater

import (
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

// Mock implementations
type mockModuleOperator struct {
	version         module.Version
	updateTime      time.Time
	org             string
	repo            string
	sha             string
	err             error
	modulesToUpdate []string // Add this new field
}

func (m *mockModuleOperator) GetLatestVersion(modulePath string) (module.Version, error) {
	if m.err != nil {
		return module.Version{}, m.err
	}
	return m.version, nil
}

func (m *mockModuleOperator) GetModuleInfo(modulePath string) (module.Version, time.Time, error) {
	if m.err != nil {
		return module.Version{}, time.Time{}, m.err
	}
	return m.version, m.updateTime, nil
}

func (m *mockModuleOperator) ParseModulePathParts(modulePath string) (org, repo string, err error) {
	if m.err != nil {
		return "", "", m.err
	}
	if m.org == "" && m.repo == "" {
		return "smartcontractkit", "chainlink", nil // Default values for tests
	}
	return m.org, m.repo, nil
}

func (m *mockModuleOperator) GetGitInfo(remote, branch string) (string, time.Time, error) {
	if m.err != nil {
		return "", time.Time{}, m.err
	}
	if m.sha == "" {
		m.sha = "abcdef123456" // default test SHA
	}
	return m.sha, m.updateTime, nil
}

// Add the missing method
func (m *mockModuleOperator) UpdateRequiredVersions(modFile *modfile.File, newVersion string) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.modulesToUpdate != nil {
		return m.modulesToUpdate, nil
	}
	// Default behavior: return the modules that have local replace directives
	var modules []string
	for _, rep := range modFile.Replace {
		if rep.New.Version == "" { // Local replace has empty version
			modules = append(modules, rep.Old.Path)
		}
	}
	return modules, nil
}

type mockSystemOperator struct {
	files map[string][]byte
	err   error
}

func newMockSystemOperator() *mockSystemOperator {
	return &mockSystemOperator{
		files: make(map[string][]byte),
	}
}

func (m *mockSystemOperator) ReadFile(path string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	content, ok := m.files[path]
	if !ok {
		return nil, fmt.Errorf("file not found: %s", path)
	}
	return content, nil
}

func (m *mockSystemOperator) WriteFile(path string, data []byte, perm os.FileMode) error {
	if m.err != nil {
		return m.err
	}
	m.files[path] = data
	return nil
}

func TestUpdater_Run(t *testing.T) {
	testTime := time.Date(2024, 11, 22, 18, 21, 10, 0, time.UTC)
	testSHA := "ac7a7395feed"

	tests := []struct {
		name     string
		config   *Config
		modOp    *mockModuleOperator
		sysOp    *mockSystemOperator
		wantErr  bool
		wantFile string
	}{
		{
			name: "successful update",
			config: &Config{
				ModulesToUpdate: []string{"github.com/smartcontractkit/chainlink/v2"},
			},
			modOp: &mockModuleOperator{
				version: module.Version{
					Path:    "github.com/smartcontractkit/chainlink/v2",
					Version: "v2.0.0",
				},
				updateTime: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/v2 v2.0.0
`)
				return m
			}(),
			wantErr: false,
		},
		{
			name: "handles module with local replace",
			config: &Config{
				ModulesToUpdate: []string{"github.com/smartcontractkit/chainlink/v2"},
			},
			modOp: &mockModuleOperator{
				version: module.Version{
					Path:    "github.com/smartcontractkit/chainlink/v2",
					Version: "v2.0.0",
				},
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/v2 v2.0.0
replace github.com/smartcontractkit/chainlink/v2 => ../
`)
				return m
			}(),
			wantErr: false,
		},
		{
			name: "v1 module update",
			config: &Config{
				ModulesToUpdate: []string{"github.com/example/mod"},
			},
			modOp: &mockModuleOperator{
				version: module.Version{
					Path:    "github.com/example/mod",
					Version: "v1.2.3",
				},
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/example/mod v1.0.0
`)
				return m
			}(),
			wantErr: false,
		},
		{
			name: "updates v2 module with timestamp",
			config: &Config{
				ModulesToUpdate: []string{"github.com/smartcontractkit/chainlink/v2"},
				RepoRemote:      "origin",
				BranchTrunk:     "develop",
			},
			modOp: &mockModuleOperator{
				sha:        "ac7a7395feed",                                   // Set exact SHA
				updateTime: time.Date(2024, 11, 22, 18, 21, 10, 0, time.UTC), // Set exact time
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/v2 v2.0.0-20241119120536-03115e80382d
`)
				return m
			}(),
			wantErr: false,
			wantFile: `module test

require github.com/smartcontractkit/chainlink/v2 v2.0.0-20241122182110-ac7a7395feed
`,
		},
		{
			name: "updates v0 module with timestamp",
			config: &Config{
				ModulesToUpdate: []string{"github.com/smartcontractkit/chainlink/deployment"},
				RepoRemote:      "origin",
				BranchTrunk:     "develop",
			},
			modOp: &mockModuleOperator{
				sha:        testSHA,
				updateTime: testTime,
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/deployment v0.0.0-20241119120536-03115e80382d
`)
				return m
			}(),
			wantErr: false,
			wantFile: `module test

require github.com/smartcontractkit/chainlink/deployment v0.0.0-20241122182110-ac7a7395feed
`,
		},
		{
			name: "handles multiple modules with different versions",
			config: &Config{
				ModulesToUpdate: []string{
					"github.com/smartcontractkit/chainlink/v2",
					"github.com/smartcontractkit/chainlink/deployment",
				},
				RepoRemote:  "origin",
				BranchTrunk: "develop",
			},
			modOp: &mockModuleOperator{
				sha:        testSHA,
				updateTime: testTime,
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require (
    github.com/smartcontractkit/chainlink/v2 v2.0.0-20241119120536-03115e80382d
    github.com/smartcontractkit/chainlink/deployment v0.0.0-20241119120536-03115e80382d
)
`)
				return m
			}(),
			wantErr: false,
			wantFile: `module test

require (
	github.com/smartcontractkit/chainlink/v2 v2.0.0-20241122182110-ac7a7395feed
	github.com/smartcontractkit/chainlink/deployment v0.0.0-20241122182110-ac7a7395feed
)
`,
		},
		{
			name: "updates v3 module with timestamp",
			config: &Config{
				ModulesToUpdate: []string{"github.com/smartcontractkit/chainlink/v3"},
				RepoRemote:      "origin",
				BranchTrunk:     "develop",
			},
			modOp: &mockModuleOperator{
				sha:        testSHA,
				updateTime: testTime,
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/v3 v3.0.0-20241119120536-03115e80382d
`)
				return m
			}(),
			wantErr: false,
			wantFile: `module test

require github.com/smartcontractkit/chainlink/v3 v3.0.0-20241122182110-ac7a7395feed
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.modOp, tt.sysOp, tt.config)
			err := u.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Updater.Run() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantFile != "" {
				got := string(tt.sysOp.files["go.mod"])
				if got != tt.wantFile {
					t.Errorf("File content mismatch\nGot:\n%s\nWant:\n%s", got, tt.wantFile)
				}
			}
		})
	}
}

func TestUpdater_FindLocalReplaceModules(t *testing.T) {
	sysOp := newMockSystemOperator()
	sysOp.files["go.mod"] = []byte(`
module test
require (
    github.com/smartcontractkit/chainlink/v2 v2.0.0
    github.com/other/repo v1.0.0
)
replace (
    github.com/smartcontractkit/chainlink/v2 => ../
    github.com/other/repo => ../other
)`)

	cfg := &Config{
		OrgName:  "smartcontractkit",
		RepoName: "chainlink",
	}

	u := New(&mockModuleOperator{}, sysOp, cfg)
	modules, err := u.findLocalReplaceModules()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if len(modules) != 1 {
		t.Errorf("expected 1 module, got %d", len(modules))
		return
	}
	if modules[0] != "github.com/smartcontractkit/chainlink/v2" {
		t.Errorf("expected chainlink module, got %s", modules[0])
	}
}
