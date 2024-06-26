package cmd

import (
	"os"
	"strings"
	"testing"
)

func TestReplaceEnvPlaceholders(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		input      string
		want       string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "All variables set",
			input:   "Path: ${{ env.PATH }}, Home: ${{ env.HOME }}",
			want:    "Path: /usr/bin, Home: /home/user",
			wantErr: false,
		},
		{
			name:       "One variable unset",
			input:      "Path: ${{ env.PATH }}, Unset: ${{ env.UNSET_VAR }}",
			want:       "Path: /usr/bin, Unset: ${{ env.UNSET_VAR }}",
			wantErr:    true,
			wantErrMsg: "environment variable 'UNSET_VAR' not set or is empty",
		},
		{
			name:       "Multiple variables unset",
			input:      "Unset1: ${{ env.UNSET_VAR1 }}, Unset2: ${{ env.UNSET_VAR2 }}",
			want:       "Unset1: ${{ env.UNSET_VAR1 }}, Unset2: ${{ env.UNSET_VAR2 }}",
			wantErr:    true,
			wantErrMsg: "environment variable 'UNSET_VAR1' not set or is empty, environment variable 'UNSET_VAR2' not set or is empty",
		},
	}

	// Set environment variables for the test
	os.Setenv("PATH", "/usr/bin")
	os.Setenv("HOME", "/home/user")
	defer os.Unsetenv("PATH")
	defer os.Unsetenv("HOME")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := replaceEnvPlaceholders(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("replaceEnvPlaceholders() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !strings.Contains(err.Error(), tt.wantErrMsg) {
				t.Errorf("replaceEnvPlaceholders() error = %v, wantErrMsg %v", err, tt.wantErrMsg)
			}
			if got != tt.want {
				t.Errorf("replaceEnvPlaceholders() got = %v, want %v", got, tt.want)
			}
		})
	}

	// Unset the variables to clean up after tests
	os.Unsetenv("UNSET_VAR1")
	os.Unsetenv("UNSET_VAR2")
}
