package updater

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

// mockGitExecutor simulates git commands for testing
type mockGitExecutor struct {
	sha  string
	time time.Time
}

func (m *mockGitExecutor) Command(ctx context.Context, args ...string) ([]byte, error) {
	switch args[0] {
	case "ls-remote":
		return []byte(m.sha + "\trefs/heads/develop\n"), nil
	case "show":
		return []byte(m.time.Format(gitTimeFormat)), nil
	default:
		return nil, fmt.Errorf("unexpected git command: %v", args)
	}
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
		sysOp    *mockSystemOperator
		wantErr  bool
		wantFile string
	}{
		{
			name: "successful update",
			config: &Config{
				ModulesToUpdate: []string{"github.com/smartcontractkit/chainlink/v2"},
				RepoRemote:      "origin",
				BranchTrunk:     "develop",
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
				RepoRemote:      "origin",
				BranchTrunk:     "develop",
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
				RepoRemote:      "origin",
				BranchTrunk:     "develop",
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
			u := New(tt.config, tt.sysOp)
			// Override the git executor with our mock
			u.git = &mockGitExecutor{
				sha:  testSHA,
				time: testTime,
			}
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

	u := New(cfg, sysOp)
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
