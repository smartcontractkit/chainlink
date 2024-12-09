package updater

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// mockGitExecutor simulates git commands for testing as an alternative to systemGitExecutor
type mockGitExecutor struct {
	sha  string
	time time.Time
}

func (m *mockGitExecutor) Command(ctx context.Context, args ...string) ([]byte, error) {
	switch args[0] {
	case "ls-remote":
		// Validate inputs
		if len(args) != 3 {
			return nil, fmt.Errorf("expected 3 args for ls-remote, got %d", len(args))
		}
		remote := args[1]
		if remote == "invalid*remote" {
			return nil, fmt.Errorf("%w: git remote '%s' contains invalid characters", ErrInvalidConfig, remote)
		}
		// Return a full 40-character SHA
		fullSHA := fmt.Sprintf("%s%s", m.sha, strings.Repeat("0", 40-len(m.sha)))
		return []byte(fullSHA + "\trefs/heads/" + args[2] + "\n"), nil
	case "show":
		if len(args) != 4 {
			return nil, fmt.Errorf("unexpected show args: %v", args)
		}
		// Use the full SHA for validation
		fullSHA := fmt.Sprintf("%s%s", m.sha, strings.Repeat("0", 40-len(m.sha)))
		if args[3] != fullSHA {
			return nil, fmt.Errorf("unexpected SHA: got %s, want %s", args[3], fullSHA)
		}
		return []byte(m.time.Format(gitTimeFormat) + "\n"), nil
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
	// Use a full 40-character SHA
	testSHA := "ac7a7395feed" + strings.Repeat("0", 28)

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
				RepoRemote:  "origin",
				BranchTrunk: "develop",
				OrgName:     "smartcontractkit",
				RepoName:    "chainlink",
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
			wantFile: fmt.Sprintf(`module test

require github.com/smartcontractkit/chainlink/v2 v2.0.0-20241122182110-%s

replace github.com/smartcontractkit/chainlink/v2 => ../
`, testSHA[:12]),
		},
		{
			name: "handles module with local replace",
			config: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "develop",
				OrgName:     "smartcontractkit",
				RepoName:    "chainlink",
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
				RepoRemote:  "origin",
				BranchTrunk: "develop",
				OrgName:     "example",
				RepoName:    "mod",
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
				RepoRemote:  "origin",
				BranchTrunk: "develop",
				OrgName:     "smartcontractkit",
				RepoName:    "chainlink",
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/v2 v2.0.0-20241119120536-03115e80382d
`)
				return m
			}(),
			wantErr: false,
			wantFile: fmt.Sprintf(`module test

require github.com/smartcontractkit/chainlink/v2 v2.0.0-20241122182110-%s

replace github.com/smartcontractkit/chainlink/v2 => ../
`, testSHA[:12]),
		},
		{
			name: "updates v0 module with timestamp",
			config: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "develop",
				OrgName:     "smartcontractkit",
				RepoName:    "chainlink",
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/deployment v0.0.0-20241119120536-03115e80382d
`)
				return m
			}(),
			wantErr: false,
			wantFile: fmt.Sprintf(`module test

require github.com/smartcontractkit/chainlink/deployment v0.0.0-20241122182110-%s

replace github.com/smartcontractkit/chainlink/deployment => ../
`, testSHA[:12]),
		},
		{
			name: "handles multiple modules with different versions",
			config: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "develop",
				OrgName:     "smartcontractkit",
				RepoName:    "chainlink",
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
			wantFile: fmt.Sprintf(`module test

require (
	github.com/smartcontractkit/chainlink/v2 v2.0.0-20241122182110-%s
	github.com/smartcontractkit/chainlink/deployment v0.0.0-20241122182110-%s
)

replace github.com/smartcontractkit/chainlink/v2 => ../

replace github.com/smartcontractkit/chainlink/deployment => ../
`, testSHA[:12], testSHA[:12]),
		},
		{
			name: "updates v3 module with timestamp",
			config: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "develop",
				OrgName:     "smartcontractkit",
				RepoName:    "chainlink",
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/smartcontractkit/chainlink/v3 v3.0.0-20241119120536-03115e80382d
`)
				return m
			}(),
			wantErr: false,
			wantFile: fmt.Sprintf(`module test

require github.com/smartcontractkit/chainlink/v3 v3.0.0-20241122182110-%s

replace github.com/smartcontractkit/chainlink/v3 => ../
`, testSHA[:12]),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.config, tt.sysOp)

			// Add local replace directive for modules that should be updated
			if !strings.Contains(tt.name, "v1 module update") {
				modContent := string(tt.sysOp.files["go.mod"])
				for _, module := range []string{"v2", "v3", "deployment"} {
					modulePath := fmt.Sprintf("github.com/%s/%s/%s",
						tt.config.OrgName, tt.config.RepoName, module)
					if strings.Contains(modContent, modulePath) &&
						!strings.Contains(modContent, "replace "+modulePath) {
						modContent += fmt.Sprintf("\nreplace %s => ../\n", modulePath)
					}
				}
				tt.sysOp.files["go.mod"] = []byte(modContent)
			}

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
