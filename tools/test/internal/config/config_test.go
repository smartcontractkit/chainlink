package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostgresContainerName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		conf *App
		want string
	}{
		{
			name: "run mode",
			conf: &App{PackageSlug: "core_services"},
			want: "test_core_services",
		},
		{
			name: "diagnose worker",
			conf: &App{DiagnoseMode: true, WorkerIndex: 1, PackageSlug: "core_services"},
			want: "iteration_1_core_services",
		},
		{
			name: "diagnose parallel worker",
			conf: &App{DiagnoseMode: true, WorkerIndex: 3, PackageSlug: "core_services"},
			want: "iteration_3_core_services",
		},
		{
			name: "missing slug defaults",
			conf: &App{},
			want: "test_pkgs",
		},
		{
			name: "diagnose defaults worker index",
			conf: &App{DiagnoseMode: true, PackageSlug: "core_services"},
			want: "iteration_1_core_services",
		},
		{
			name: "nil config",
			conf: nil,
			want: "test_pkgs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, tt.conf.PostgresContainerName())
		})
	}
}
