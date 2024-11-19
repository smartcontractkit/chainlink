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
				ModulesToUpdate: []string{"test.com/mod"},
				RepoRemote:     "origin",
				BranchTrunk:    "main",
			},
			wantErr: false,
		},
		{
			name: "missing modules",
			cfg: &Config{
				RepoRemote:  "origin",
				BranchTrunk: "main",
			},
			wantErr: true,
		},
		{
			name: "version flag bypasses validation",
			cfg: &Config{
				ShowVersion: true,
			},
			wantErr: false,
		},
		{
			name: "update-org-modules bypasses module validation",
			cfg: &Config{
				UpdateOrgModules: true,
				RepoRemote:      "origin",
				BranchTrunk:     "main",
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
			name: "multiple modules via CLI",
			args: []string{"-module", "mod1.com", "-module", "mod2.com"},
			wantCfg: &Config{
				ModulesToUpdate: []string{"mod1.com", "mod2.com"},
				RepoRemote:     "origin",
				BranchTrunk:    "develop",
			},
			wantErr: false,
		},
		{
			name: "show version",
			args: []string{"-version"},
			wantCfg: &Config{
				ShowVersion: true,
				RepoRemote: "origin",    // Default value
				BranchTrunk: "develop",  // Default value
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
			if len(got.ModulesToUpdate) != len(tt.wantCfg.ModulesToUpdate) {
				t.Errorf("ParseFlags() ModulesToUpdate length = %v, want %v", len(got.ModulesToUpdate), len(tt.wantCfg.ModulesToUpdate))
			}
			for i, module := range got.ModulesToUpdate {
				if module != tt.wantCfg.ModulesToUpdate[i] {
					t.Errorf("ParseFlags() ModulesToUpdate[%d] = %v, want %v", i, module, tt.wantCfg.ModulesToUpdate[i])
				}
			}
		})
	}
}
