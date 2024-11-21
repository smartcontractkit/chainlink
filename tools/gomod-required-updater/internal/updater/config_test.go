package updater

import "testing"

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name: "valid config",
			cfg: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "main",
			},
			wantErr: false,
		},
		{
			name: "version flag bypasses validation",
			cfg: &Config{
				ShowVersion: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
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
