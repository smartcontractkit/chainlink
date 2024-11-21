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

func (m *mockSystemOperator) WriteFile(path string, data []byte, perm uint32) error {
	if m.err != nil {
		return m.err
	}
	m.files[path] = data
	return nil
}

func TestUpdater_Run(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		gitOp    *mockGitOperator
		sysOp    *mockSystemOperator
		wantErr  bool
		wantFile string
	}{
		{
			name: "successful update",
			config: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "main",
			},
			gitOp: &mockGitOperator{
				sha:        "abc123def456",
				commitDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				org:        "smartcontractkit",
				repo:       "chainlink",
			},
			sysOp: func() *mockSystemOperator {
				m := newMockSystemOperator()
				m.files["go.mod"] = []byte(`module github.com/smartcontractkit/chainlink/v2
require test.com/mod v1.0.0
`)
				return m
			}(),
			wantErr: false,
		},
		{
			name: "handles module with local replace",
			config: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "main",
			},
			gitOp: &mockGitOperator{
				sha:        "abc123def456",
				commitDate: time.Now(),
				org:        "smartcontractkit",
				repo:       "chainlink",
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

	u := New(&mockGitOperator{}, sysOp, cfg)
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
