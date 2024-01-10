package toml

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	configutils "github.com/smartcontractkit/chainlink/v2/core/utils/config"
)

func TestMercurySecrets_valid(t *testing.T) {
	ms := MercurySecrets{
		Credentials: map[string]MercuryCredentials{
			"cred1": {
				URL:      models.MustSecretURL("https://facebook.com"),
				Username: models.NewSecret("new user1"),
				Password: models.NewSecret("new password1"),
			},
			"cred2": {
				URL:      models.MustSecretURL("HTTPS://GOOGLE.COM"),
				Username: models.NewSecret("new user1"),
				Password: models.NewSecret("new password2"),
			},
			"cred3": {
				LegacyURL: models.MustSecretURL("https://abc.com"),
				URL:       models.MustSecretURL("HTTPS://GOOGLE1.COM"),
				Username:  models.NewSecret("new user1"),
				Password:  models.NewSecret("new password2"),
			},
		},
	}

	err := ms.ValidateConfig()
	assert.NoError(t, err)
}

func TestMercurySecrets_duplicateURLs(t *testing.T) {
	ms := MercurySecrets{
		Credentials: map[string]MercuryCredentials{
			"cred1": {
				URL:      models.MustSecretURL("HTTPS://GOOGLE.COM"),
				Username: models.NewSecret("new user1"),
				Password: models.NewSecret("new password1"),
			},
			"cred2": {
				URL:      models.MustSecretURL("HTTPS://GOOGLE.COM"),
				Username: models.NewSecret("new user2"),
				Password: models.NewSecret("new password2"),
			},
		},
	}

	err := ms.ValidateConfig()
	assert.Error(t, err)
	assert.Equal(t, "URL: invalid value (https://GOOGLE.COM): duplicate - must be unique", err.Error())
}

func TestMercurySecrets_emptyURL(t *testing.T) {
	ms := MercurySecrets{
		Credentials: map[string]MercuryCredentials{
			"cred1": {
				URL:      nil,
				Username: models.NewSecret("new user1"),
				Password: models.NewSecret("new password1"),
			},
		},
	}

	err := ms.ValidateConfig()
	assert.Error(t, err)
	assert.Equal(t, "URL: missing: must be provided and non-empty", err.Error())
}

func Test_validateDBURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{"no user or password", "postgresql://foo.example.com:5432/chainlink?application_name=Test+Application", "DB URL must be authenticated; plaintext URLs are not allowed"},
		{"with user and no password", "postgresql://myuser@foo.example.com:5432/chainlink?application_name=Test+Application", "DB URL must be authenticated; password is required"},
		{"with user and password of insufficient length", "postgresql://myuser:shortpw@foo.example.com:5432/chainlink?application_name=Test+Application", fmt.Sprintf("%s	%s\n", utils.ErrMsgHeader, "password is less than 16 characters long")},
		{"with no user and password of sufficient length", "postgresql://:thisisareallylongpassword@foo.example.com:5432/chainlink?application_name=Test+Application", ""},
		{"with user and password of sufficient length", "postgresql://myuser:thisisareallylongpassword@foo.example.com:5432/chainlink?application_name=Test+Application", ""},
		{"with user and password of insufficient length as params", "postgresql://foo.example.com:5432/chainlink?application_name=Test+Application&password=shortpw&user=myuser", fmt.Sprintf("%s	%s\n", utils.ErrMsgHeader, "password is less than 16 characters long")},
		{"with no user and password of sufficient length as params", "postgresql://foo.example.com:5432/chainlink?application_name=Test+Application&password=thisisareallylongpassword", ""},
		{"with user and password of sufficient length as params", "postgresql://foo.example.com:5432/chainlink?application_name=Test+Application&password=thisisareallylongpassword&user=myuser", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := testutils.MustParseURL(t, test.url)
			err := validateDBURL(*url)
			if test.wantErr == "" {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, test.wantErr)
			}
		})
	}
}

func TestDatabaseSecrets_ValidateConfig(t *testing.T) {
	validUrl := commonconfig.URL(url.URL{Scheme: "https", Host: "localhost"})
	validSecretURL := *models.NewSecretURL(&validUrl)

	invalidEmptyUrl := commonconfig.URL(url.URL{})
	invalidEmptySecretURL := *models.NewSecretURL(&invalidEmptyUrl)

	invalidBackupURL := commonconfig.URL(url.URL{Scheme: "http", Host: "localhost"})
	invalidBackupSecretURL := *models.NewSecretURL(&invalidBackupURL)

	tests := []struct {
		name                string
		input               *DatabaseSecrets
		buildMode           string
		expectedErrContains []string
	}{
		{
			name: "Nil URL",
			input: &DatabaseSecrets{
				URL: nil,
			},
			expectedErrContains: []string{"URL: empty: must be provided and non-empty"},
		},
		{
			name: "Empty URL",
			input: &DatabaseSecrets{
				URL: &invalidEmptySecretURL,
			},
			expectedErrContains: []string{"URL: empty: must be provided and non-empty"},
		},
		{
			name: "Insecure Password in Production",
			input: &DatabaseSecrets{
				URL:                  &validSecretURL,
				AllowSimplePasswords: &[]bool{true}[0],
			},
			buildMode:           build.Prod,
			expectedErrContains: []string{"insecure configs are not allowed on secure builds"},
		},
		{
			name: "Invalid Backup URL with Simple Passwords Not Allowed",
			input: &DatabaseSecrets{
				URL:                  &validSecretURL,
				BackupURL:            &invalidBackupSecretURL,
				AllowSimplePasswords: &[]bool{false}[0],
			},
			expectedErrContains: []string{"missing or insufficiently complex password"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildMode := build.Mode()
			if tt.buildMode != "" {
				buildMode = tt.buildMode
			}
			err := tt.input.validateConfig(buildMode)
			if err == nil && len(tt.expectedErrContains) > 0 {
				t.Errorf("expected errors but got none")
				return
			}

			if err != nil {
				errStr := err.Error()
				for _, expectedErrSubStr := range tt.expectedErrContains {
					if !strings.Contains(errStr, expectedErrSubStr) {
						t.Errorf("expected error to contain substring %q but got %v", expectedErrSubStr, errStr)
					}
				}
			}
		})
	}
}
func TestTracing_ValidateCollectorTarget(t *testing.T) {
	tests := []struct {
		name            string
		collectorTarget *string
		mode            *string
		wantErr         bool
		errMsg          string
	}{
		{
			name:            "valid http address in tls mode",
			collectorTarget: ptr("https://testing.collector.dev"),
			mode:            ptr("tls"),
			wantErr:         false,
		},
		{
			name:            "valid http address in unencrypted mode",
			collectorTarget: ptr("https://localhost:4317"),
			mode:            ptr("unencrypted"),
			wantErr:         true,
			errMsg:          "CollectorTarget: invalid value (https://localhost:4317): must be a valid local URI",
		},
		// Tracing.Mode = 'tls'
		{
			name:            "valid localhost address",
			collectorTarget: ptr("localhost:4317"),
			mode:            ptr("tls"),
			wantErr:         false,
		},
		{
			name:            "valid docker address",
			collectorTarget: ptr("otel-collector:4317"),
			mode:            ptr("tls"),
			wantErr:         false,
		},
		{
			name:            "valid IP address",
			collectorTarget: ptr("192.168.1.1:4317"),
			mode:            ptr("tls"),
			wantErr:         false,
		},
		{
			name:            "invalid port",
			collectorTarget: ptr("localhost:invalid"),
			wantErr:         true,
			mode:            ptr("tls"),
			errMsg:          "CollectorTarget: invalid value (localhost:invalid): must be a valid URI",
		},
		{
			name:            "invalid address",
			collectorTarget: ptr("invalid address"),
			wantErr:         true,
			mode:            ptr("tls"),
			errMsg:          "CollectorTarget: invalid value (invalid address): must be a valid URI",
		},
		{
			name:            "nil CollectorTarget",
			collectorTarget: ptr(""),
			wantErr:         true,
			mode:            ptr("tls"),
			errMsg:          "CollectorTarget: invalid value (): must be a valid URI",
		},
		// Tracing.Mode = 'unencrypted'
		{
			name:            "valid localhost address",
			collectorTarget: ptr("localhost:4317"),
			mode:            ptr("unencrypted"),
			wantErr:         false,
		},
		{
			name:            "valid docker address",
			collectorTarget: ptr("otel-collector:4317"),
			mode:            ptr("unencrypted"),
			wantErr:         false,
		},
		{
			name:            "valid IP address",
			collectorTarget: ptr("192.168.1.1:4317"),
			mode:            ptr("unencrypted"),
			wantErr:         false,
		},
		{
			name:            "invalid port",
			collectorTarget: ptr("localhost:invalid"),
			wantErr:         true,
			mode:            ptr("unencrypted"),
			errMsg:          "CollectorTarget: invalid value (localhost:invalid): must be a valid local URI",
		},
		{
			name:            "invalid address",
			collectorTarget: ptr("invalid address"),
			wantErr:         true,
			mode:            ptr("unencrypted"),
			errMsg:          "CollectorTarget: invalid value (invalid address): must be a valid local URI",
		},
		{
			name:            "nil CollectorTarget",
			collectorTarget: ptr(""),
			wantErr:         true,
			mode:            ptr("unencrypted"),
			errMsg:          "CollectorTarget: invalid value (): must be a valid local URI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tlsCertPath string
			if *tt.mode == "tls" {
				tlsCertPath = "/path/to/cert.pem"
			}
			tracing := &Tracing{
				Enabled:         ptr(true),
				TLSCertPath:     &tlsCertPath,
				Mode:            tt.mode,
				CollectorTarget: tt.collectorTarget,
			}

			err := tracing.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTracing_ValidateSamplingRatio(t *testing.T) {
	tests := []struct {
		name          string
		samplingRatio *float64
		wantErr       bool
		errMsg        string
	}{
		{
			name:          "valid lower bound",
			samplingRatio: ptr(0.0),
			wantErr:       false,
		},
		{
			name:          "valid upper bound",
			samplingRatio: ptr(1.0),
			wantErr:       false,
		},
		{
			name:          "valid value",
			samplingRatio: ptr(0.5),
			wantErr:       false,
		},
		{
			name:          "invalid negative value",
			samplingRatio: ptr(-0.1),
			wantErr:       true,
			errMsg:        configutils.ErrInvalid{Name: "SamplingRatio", Value: -0.1, Msg: "must be between 0 and 1"}.Error(),
		},
		{
			name:          "invalid value greater than 1",
			samplingRatio: ptr(1.1),
			wantErr:       true,
			errMsg:        configutils.ErrInvalid{Name: "SamplingRatio", Value: 1.1, Msg: "must be between 0 and 1"}.Error(),
		},
		{
			name:          "nil SamplingRatio",
			samplingRatio: nil,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracing := Tracing{
				SamplingRatio: tt.samplingRatio,
				Enabled:       ptr(true),
			}

			err := tracing.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTracing_ValidateTLSCertPath(t *testing.T) {
	// tests for Tracing.Mode = 'tls'
	tls_tests := []struct {
		name        string
		tlsCertPath *string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid file path",
			tlsCertPath: ptr("/etc/ssl/certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "relative file path",
			tlsCertPath: ptr("certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "excessively long file path",
			tlsCertPath: ptr(strings.Repeat("z", 4097)),
			wantErr:     true,
			errMsg:      "TLSCertPath: invalid value (" + strings.Repeat("z", 4097) + "): must be a valid file path",
		},
		{
			name:        "empty file path",
			tlsCertPath: ptr(""),
			wantErr:     true,
			errMsg:      "TLSCertPath: invalid value (): must be a valid file path",
		},
	}

	// tests for Tracing.Mode = 'unencrypted'
	unencrypted_tests := []struct {
		name        string
		tlsCertPath *string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid file path",
			tlsCertPath: ptr("/etc/ssl/certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "relative file path",
			tlsCertPath: ptr("certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "excessively long file path",
			tlsCertPath: ptr(strings.Repeat("z", 4097)),
			wantErr:     false,
		},
		{
			name:        "empty file path",
			tlsCertPath: ptr(""),
			wantErr:     false,
		},
	}

	for _, tt := range tls_tests {
		t.Run(tt.name, func(t *testing.T) {
			tracing := &Tracing{
				Mode:        ptr("tls"),
				TLSCertPath: tt.tlsCertPath,
				Enabled:     ptr(true),
			}

			err := tracing.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}

	for _, tt := range unencrypted_tests {
		t.Run(tt.name, func(t *testing.T) {
			tracing := &Tracing{
				Mode:        ptr("unencrypted"),
				TLSCertPath: tt.tlsCertPath,
				Enabled:     ptr(true),
			}

			err := tracing.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTracing_ValidateMode(t *testing.T) {
	tests := []struct {
		name        string
		mode        *string
		tlsCertPath *string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "tls mode with valid TLS path",
			mode:        ptr("tls"),
			tlsCertPath: ptr("/path/to/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "tls mode without TLS path",
			mode:        ptr("tls"),
			tlsCertPath: nil,
			wantErr:     true,
			errMsg:      "TLSCertPath: missing: must be set when Tracing.Mode is tls",
		},
		{
			name:        "unencrypted mode with TLS path",
			mode:        ptr("unencrypted"),
			tlsCertPath: ptr("/path/to/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "unencrypted mode without TLS path",
			mode:        ptr("unencrypted"),
			tlsCertPath: nil,
			wantErr:     false,
		},
		{
			name:        "invalid mode",
			mode:        ptr("unknown"),
			tlsCertPath: nil,
			wantErr:     true,
			errMsg:      "Mode: invalid value (unknown): must be either 'tls' or 'unencrypted'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracing := &Tracing{
				Enabled:     ptr(true),
				Mode:        tt.mode,
				TLSCertPath: tt.tlsCertPath,
			}

			err := tracing.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMercuryTLS_ValidateTLSCertPath(t *testing.T) {
	tests := []struct {
		name        string
		tlsCertPath *string
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid file path",
			tlsCertPath: ptr("/etc/ssl/certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "relative file path",
			tlsCertPath: ptr("certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "excessively long file path",
			tlsCertPath: ptr(strings.Repeat("z", 4097)),
			wantErr:     true,
			errMsg:      "CertFile: invalid value (" + strings.Repeat("z", 4097) + "): must be a valid file path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mercury := &Mercury{
				TLS: MercuryTLS{
					CertFile: tt.tlsCertPath,
				},
			}

			err := mercury.ValidateConfig()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ptr is a utility function for converting a value to a pointer to the value.
func ptr[T any](t T) *T { return &t }
