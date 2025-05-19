package updater

import (
	"errors"
	"testing"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				RepoRemote:  DefaultRepoRemote,
				BranchTrunk: DefaultBranchTrunk,
				OrgName:     DefaultOrgName,
				RepoName:    DefaultRepoName,
			},
			wantErr: false,
		},
		{
			name: "version flag bypasses validation",
			config: &Config{
				ShowVersion: true,
			},
			wantErr: false,
		},
		{
			name: "missing repo remote",
			config: &Config{
				BranchTrunk: DefaultBranchTrunk,
				OrgName:     DefaultOrgName,
				RepoName:    DefaultRepoName,
			},
			wantErr: true,
		},
		{
			name: "missing branch trunk",
			config: &Config{
				RepoRemote: DefaultRepoRemote,
				OrgName:    DefaultOrgName,
				RepoName:   DefaultRepoName,
			},
			wantErr: true,
		},
		{
			name: "missing org name",
			config: &Config{
				RepoRemote:  DefaultRepoRemote,
				BranchTrunk: DefaultBranchTrunk,
				RepoName:    DefaultRepoName,
			},
			wantErr: true,
		},
		{
			name: "missing repo name",
			config: &Config{
				RepoRemote:  DefaultRepoRemote,
				BranchTrunk: DefaultBranchTrunk,
				OrgName:     DefaultOrgName,
			},
			wantErr: true,
		},
		{
			name: "invalid remote characters",
			config: &Config{
				OrgName:     "test",
				RepoName:    "test",
				RepoRemote:  "origin!@#",
				BranchTrunk: "main",
			},
			wantErr: true,
		},
		{
			name: "invalid branch characters",
			config: &Config{
				OrgName:     "test",
				RepoName:    "test",
				RepoRemote:  "origin",
				BranchTrunk: "main!@#",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConfig_ValidateErrorType(t *testing.T) {
	cfg := &Config{
		RepoRemote:  "invalid*remote",
		BranchTrunk: "develop",
		OrgName:     "test",
		RepoName:    "test",
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("expected error due to invalid repo remote, got nil")
		return
	}

	if !errors.Is(err, ErrInvalidConfig) {
		t.Errorf("expected error to be ErrInvalidConfig, got: %v", err)
	}
}

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantCfg *Config
		wantErr bool
	}{
		{
			name: "default flags",
			args: []string{},
			wantCfg: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "develop",
			},
			wantErr: false,
		},
		{
			name: "show version",
			args: []string{"-version"},
			wantCfg: &Config{
				ShowVersion: true,
				RepoRemote:  "origin",
				BranchTrunk: "develop",
			},
			wantErr: false,
		},
		{
			name: "custom remote and branch",
			args: []string{"-repo-remote", "upstream", "-branch-trunk", "main"},
			wantCfg: &Config{
				RepoRemote:  "upstream",
				BranchTrunk: "main",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFlags(tt.args, "test-version")
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			if got.RepoRemote != tt.wantCfg.RepoRemote {
				t.Errorf("ParseFlags() RepoRemote = %v, want %v", got.RepoRemote, tt.wantCfg.RepoRemote)
			}
			if got.BranchTrunk != tt.wantCfg.BranchTrunk {
				t.Errorf("ParseFlags() BranchTrunk = %v, want %v", got.BranchTrunk, tt.wantCfg.BranchTrunk)
			}
			if got.ShowVersion != tt.wantCfg.ShowVersion {
				t.Errorf("ParseFlags() ShowVersion = %v, want %v", got.ShowVersion, tt.wantCfg.ShowVersion)
			}
		})
	}
}
