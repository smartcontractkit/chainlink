package updater

import (
	"fmt"
	"testing"
	"time"
)

// Mock implementations
type mockGitOperator struct {
	sha        string
	commitDate time.Time
	org        string
	repo       string
	err        error
}

func (m *mockGitOperator) GetSHA(remote, branch string) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.sha, nil
}

func (m *mockGitOperator) GetCommitDate(sha string) (time.Time, error) {
	if m.err != nil {
		return time.Time{}, m.err
	}
	return m.commitDate, nil
}

func (m *mockGitOperator) GetRepoInfo(remote string) (org, repo string, err error) {
	if m.err != nil {
		return "", "", m.err
	}
	if m.org == "" && m.repo == "" {
		return "smartcontractkit", "chainlink", nil // Default values for tests
	}
	return m.org, m.repo, nil
}

type mockSystemOperator struct {
	files    map[string][]byte
	walkFn   func(root string, fn func(path string, isDir bool) error) error // Update signature
	commands []string
	err      error
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

func (m *mockSystemOperator) WriteFile(path string, data []byte, perm uint32) error {
	if m.err != nil {
		return m.err
	}
	m.files[path] = data
	return nil
}

func (m *mockSystemOperator) Walk(root string, fn func(path string, isDir bool) error) error {
	if m.walkFn != nil {
		return m.walkFn(root, fn) // Pass through the callback
	}
	return nil
}

func (m *mockSystemOperator) Chdir(dir string) error { return m.err }
func (m *mockSystemOperator) Getwd() (string, error) { return "/mock/dir", m.err }
func (m *mockSystemOperator) RunCommand(name string, args ...string) error {
	if m.err != nil {
		return m.err
	}
	m.commands = append(m.commands, name+" "+fmt.Sprint(args))
	return nil
}

func TestUpdater_Run(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		gitOp    *mockGitOperator
		sysOp    *mockSystemOperator
		wantErr  bool
		wantFile string // Expected file content after update
	}{
		{
			name: "successful update",
			config: &Config{
				ModulesToUpdate: []string{"github.com/example/module"},
				RepoRemote:     "origin",
				BranchTrunk:    "main",
				RootPath:       ".",
			},
			gitOp: &mockGitOperator{
				sha:        "abc123def456", // 12 chars
				commitDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module test
require github.com/example/module v0.0.0-20230101000000-123456789abc`)
				m.walkFn = func(root string, fn func(path string, isDir bool) error) error {
					return fn("go.mod", false)
				}
				return m
			}(),
			wantErr: false,
		},
		{
			name: "handles multiple go.mod files",
			config: &Config{
				ModulesToUpdate: []string{"test.com/mod"},
				RepoRemote:     "origin",
				BranchTrunk:    "main",
			},
			gitOp: &mockGitOperator{
				sha:        "def456789012", // 12 chars
				commitDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["dir1/go.mod"] = []byte(`require test.com/mod v1.0.0`)
				m.files["dir2/go.mod"] = []byte(`require test.com/mod v1.0.0`)
				paths := []string{"dir1/go.mod", "dir2/go.mod"}
				var i int
				m.walkFn = func(root string, fn func(path string, isDir bool) error) error {
					for ; i < len(paths); i++ {
						if err := fn(paths[i], false); err != nil {
							return err
						}
					}
					return nil
				}
				return m
			}(),
			wantErr: false,
		},
		{
			name: "handles module with version suffix",
			config: &Config{
				ModulesToUpdate: []string{"test.com/mod/v2"},
				RepoRemote:     "origin",
				BranchTrunk:    "main",
			},
			gitOp: &mockGitOperator{
				sha:        "abc123def456", // 12 chars
				commitDate: time.Now(),
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`require test.com/mod/v2 v2.0.0`)
				m.walkFn = func(root string, fn func(path string, isDir bool) error) error {
					return fn("go.mod", false)
				}
				return m
			}(),
			wantErr: false,
		},
		{
			name: "finds and updates org modules with local replaces",
			config: &Config{
				UpdateOrgModules: true,
				RepoRemote:     "origin",
				BranchTrunk:    "main",
				OrgName:        "smartcontractkit",
				RepoName:       "chainlink",
			},
			gitOp: &mockGitOperator{
				sha:        "abc123def456", // 12 chars
				commitDate: time.Now(),
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				// Root go.mod for repo detection
				m.files["go.mod"] = []byte(`module github.com/smartcontractkit/chainlink/v2`)
				// Module with local replace
				m.files["core/go.mod"] = []byte(`
module test
require (
    github.com/smartcontractkit/chainlink/v2 v2.0.0
)
replace github.com/smartcontractkit/chainlink/v2 => ../
`)
				return m
			}(),
			wantErr: false,
		},
		{
			name: "skips non-org modules with local replaces",
			config: &Config{
				UpdateOrgModules: true,
				RepoRemote:     "origin",
				BranchTrunk:    "main",
				OrgName:        "smartcontractkit",
				RepoName:       "chainlink",
			},
			gitOp: &mockGitOperator{
				sha:        "abc123",
				commitDate: time.Now(),
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`
module test
require (
    github.com/other/repo v1.0.0
)
replace github.com/other/repo => ../other
`)
				return m
			}(),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := New(tt.gitOp, tt.sysOp, tt.config)
			err := u.Run()
			if (err != nil) != tt.wantErr {
				t.Errorf("Updater.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdater_UpdateMultipleModules(t *testing.T) {
	sysOp := newMockSystemOperator()
	sysOp.files["go.mod"] = []byte(`
module test
require (
	mod1.com v1.0.0
	mod2.com v1.0.0
)`)
	sysOp.walkFn = func(root string, fn func(path string, isDir bool) error) error {
		return fn("go.mod", false)
	}

	gitOp := &mockGitOperator{
		sha:        "abc123def456", // 12 chars
		commitDate: time.Now(),
	}

	cfg := &Config{
		ModulesToUpdate: []string{"mod1.com", "mod2.com"},
		RepoRemote:     "origin",
		BranchTrunk:    "main",
	}

	u := New(gitOp, sysOp, cfg)
	if err := u.Run(); err != nil {
		t.Errorf("unexpected error: %v", err)
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

    // Setup Walk to properly handle the callback
    sysOp.walkFn = func(root string, fn func(path string, isDir bool) error) error {
        return fn("go.mod", false) // Call the callback with our test file
    }

    cfg := &Config{
        UpdateOrgModules: true,
        OrgName:        "smartcontractkit",
        RepoName:       "chainlink",
    }

    u := New(&mockGitOperator{}, sysOp, cfg)
    modules, err := u.findLocalReplaceModules(".")
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