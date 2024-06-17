package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/store/migrate/internal"
)

func Test_loadManifest(t *testing.T) {
	/*
		coreMigration := func(fileName string) string {
			return filepath.Join(migrate.MIGRATIONS_DIR, fileName)
		}
		pluginMigration := func(pluginKind, pluginVariant, fileName string) string {
			return filepath.Join(migrate.PLUGIN_MIGRATIONS_DIR, pluginKind, pluginVariant, fileName)
		}
	*/
	type args struct {
		txt string
	}
	tests := []struct {
		name    string
		args    args
		want    internal.Manifest
		wantErr bool
	}{
		{
			name: "valid manifest",
			args: args{
				txt: `
core/store/migrate/plugins/app/appname/0001_xyz.sql	
core/store/migrate/plugins/relayer/evm/0001_a.sql				
core/store/migrate/migrations/0001_initial_schema.sql
`,
			},
			want: internal.Manifest{
				Entries: []internal.ManifestEntry{
					internal.ManifestEntry{
						Type:          "plugin",
						PluginKind:    "app",
						PluginVariant: "appname",
						Version:       1,
					},
					internal.ManifestEntry{
						Type:          "plugin",
						PluginKind:    "relayer",
						PluginVariant: "evm",
						Version:       1,
					},
					internal.ManifestEntry{
						Type:    "core",
						Version: 1,
					},
				},
			},
		},
		{
			name: "invalid root directory",
			args: args{
				txt: `
wrong/prefix/plugins/relayer/evm/0001_a.sql				
core/store/migrate/migrations/0001_initial_schema.sql
`,
			},
			wantErr: true,
		},

		{
			name: "invalid migration name, no version",
			args: args{
				txt: `			
core/store/migrate/migrations/initial_schema.sql
`,
			},
			wantErr: true,
		},
		{
			name: "invalid migration name, version not 4 digits",
			args: args{
				txt: `			
core/store/migrate/migrations/10_initial_schema.sql
`,
			},
			wantErr: true,
		},
		{
			name: "invalid core migration sub directory",
			args: args{
				txt: `			
core/store/migrate/migrations/WRONG/0001_initial_schema.sql
`,
			},
			wantErr: true,
		},
		{
			name: "invalid plugin migration: missing plugin kind, variant",
			args: args{
				txt: `			
core/store/migrate/plugins/0001_initial_schema.sql
`,
			},
			wantErr: true,
		},
		{
			name: "invalid plugin migration: not a relayer or app",
			args: args{
				txt: `			
core/store/migrate/plugins/wrongkind/variant/0001_initial_schema.sql
`,
			},
			wantErr: true,
		},
		{
			name: "plugin migration",
			args: args{
				txt: `core/store/migrate/plugins/relayer/evm/0001_a.sql
`,
			},
			want: internal.Manifest{
				Entries: []internal.ManifestEntry{
					internal.ManifestEntry{
						Type:          "plugin",
						PluginKind:    "relayer",
						PluginVariant: "evm",
						Version:       1,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := internal.LoadManifest(tt.args.txt)
			require.Equal(t, err != nil, tt.wantErr, "LoadManifest() error = %v, wantErr %v", err, tt.wantErr)
			require.Equal(t, len(tt.want.Entries), len(got.Entries))
			for i, entry := range got.Entries {
				assert.Equal(t, tt.want.Entries[i].Type, entry.Type, "type mismatch %d", i)
				assert.Equal(t, tt.want.Entries[i].PluginKind, entry.PluginKind, "plugin kind mismatch %d", i)
				assert.Equal(t, tt.want.Entries[i].PluginVariant, entry.PluginVariant, "plugin variant mismatch %d", i)
				assert.Equal(t, tt.want.Entries[i].Version, entry.Version, "version mismatch %d", i)
			}
		})
	}
}
