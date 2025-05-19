package validate_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestNewCCIPSpecToml(t *testing.T) {
	tests := []struct {
		name     string
		specArgs validate.SpecArgs
		want     string
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate.NewCCIPSpecToml(tt.specArgs)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestValidatedCCIPSpec(t *testing.T) {
	type args struct {
		tomlString string
	}
	tests := []struct {
		name    string
		args    args
		wantJb  job.Job
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJb, err := validate.ValidatedCCIPSpec(tt.args.tomlString)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantJb, gotJb)
			}
		})
	}
}
