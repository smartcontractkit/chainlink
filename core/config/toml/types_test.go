package toml

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		{"pgtestdb instance with short password", "postgresql://pgtdbuser:short@localhost:5432/testdb_tpl_abc_inst_def?sslmode=disable", ""},
		{"chainlink_test with short password", "postgresql://postgres:short@localhost:5432/chainlink_test?sslmode=disable", ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			url := testutils.MustParseURL(t, test.url)
			err := validateDBURL(*url)
			if test.wantErr == "" {
				assert.NoError(t, err)
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
			collectorTarget: new("https://testing.collector.dev"),
			mode:            new("tls"),
			wantErr:         false,
		},
		{
			name:            "valid http address in unencrypted mode",
			collectorTarget: new("https://localhost:4317"),
			mode:            new("unencrypted"),
			wantErr:         true,
			errMsg:          "CollectorTarget: invalid value (https://localhost:4317): must be a valid local URI",
		},
		// Tracing.Mode = 'tls'
		{
			name:            "valid localhost address",
			collectorTarget: new("localhost:4317"),
			mode:            new("tls"),
			wantErr:         false,
		},
		{
			name:            "valid docker address",
			collectorTarget: new("otel-collector:4317"),
			mode:            new("tls"),
			wantErr:         false,
		},
		{
			name:            "valid IP address",
			collectorTarget: new("192.168.1.1:4317"),
			mode:            new("tls"),
			wantErr:         false,
		},
		{
			name:            "invalid port",
			collectorTarget: new("localhost:invalid"),
			wantErr:         true,
			mode:            new("tls"),
			errMsg:          "CollectorTarget: invalid value (localhost:invalid): must be a valid URI",
		},
		{
			name:            "invalid address",
			collectorTarget: new("invalid address"),
			wantErr:         true,
			mode:            new("tls"),
			errMsg:          "CollectorTarget: invalid value (invalid address): must be a valid URI",
		},
		{
			name:            "nil CollectorTarget",
			collectorTarget: new(""),
			wantErr:         true,
			mode:            new("tls"),
			errMsg:          "CollectorTarget: invalid value (): must be a valid URI",
		},
		// Tracing.Mode = 'unencrypted'
		{
			name:            "valid localhost address",
			collectorTarget: new("localhost:4317"),
			mode:            new("unencrypted"),
			wantErr:         false,
		},
		{
			name:            "valid docker address",
			collectorTarget: new("otel-collector:4317"),
			mode:            new("unencrypted"),
			wantErr:         false,
		},
		{
			name:            "valid IP address",
			collectorTarget: new("192.168.1.1:4317"),
			mode:            new("unencrypted"),
			wantErr:         false,
		},
		{
			name:            "invalid port",
			collectorTarget: new("localhost:invalid"),
			wantErr:         true,
			mode:            new("unencrypted"),
			errMsg:          "CollectorTarget: invalid value (localhost:invalid): must be a valid local URI",
		},
		{
			name:            "invalid address",
			collectorTarget: new("invalid address"),
			wantErr:         true,
			mode:            new("unencrypted"),
			errMsg:          "CollectorTarget: invalid value (invalid address): must be a valid local URI",
		},
		{
			name:            "nil CollectorTarget",
			collectorTarget: new(""),
			wantErr:         true,
			mode:            new("unencrypted"),
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
				Enabled:         new(true),
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
			samplingRatio: new(0.0),
			wantErr:       false,
		},
		{
			name:          "valid upper bound",
			samplingRatio: new(1.0),
			wantErr:       false,
		},
		{
			name:          "valid value",
			samplingRatio: new(0.5),
			wantErr:       false,
		},
		{
			name:          "invalid negative value",
			samplingRatio: new(-0.1),
			wantErr:       true,
			errMsg:        configutils.ErrInvalid{Name: "SamplingRatio", Value: -0.1, Msg: "must be between 0 and 1"}.Error(),
		},
		{
			name:          "invalid value greater than 1",
			samplingRatio: new(1.1),
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
				Enabled:       new(true),
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
			tlsCertPath: new("/etc/ssl/certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "relative file path",
			tlsCertPath: new("certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "excessively long file path",
			tlsCertPath: new(strings.Repeat("z", 4097)),
			wantErr:     true,
			errMsg:      "TLSCertPath: invalid value (" + strings.Repeat("z", 4097) + "): must be a valid file path",
		},
		{
			name:        "empty file path",
			tlsCertPath: new(""),
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
			tlsCertPath: new("/etc/ssl/certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "relative file path",
			tlsCertPath: new("certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "excessively long file path",
			tlsCertPath: new(strings.Repeat("z", 4097)),
			wantErr:     false,
		},
		{
			name:        "empty file path",
			tlsCertPath: new(""),
			wantErr:     false,
		},
	}

	for _, tt := range tls_tests {
		t.Run(tt.name, func(t *testing.T) {
			tracing := &Tracing{
				Mode:        new("tls"),
				TLSCertPath: tt.tlsCertPath,
				Enabled:     new(true),
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
				Mode:        new("unencrypted"),
				TLSCertPath: tt.tlsCertPath,
				Enabled:     new(true),
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
			mode:        new("tls"),
			tlsCertPath: new("/path/to/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "tls mode without TLS path",
			mode:        new("tls"),
			tlsCertPath: nil,
			wantErr:     true,
			errMsg:      "TLSCertPath: missing: must be set when Tracing.Mode is tls",
		},
		{
			name:        "unencrypted mode with TLS path",
			mode:        new("unencrypted"),
			tlsCertPath: new("/path/to/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "unencrypted mode without TLS path",
			mode:        new("unencrypted"),
			tlsCertPath: nil,
			wantErr:     false,
		},
		{
			name:        "invalid mode",
			mode:        new("unknown"),
			tlsCertPath: nil,
			wantErr:     true,
			errMsg:      "Mode: invalid value (unknown): must be either 'tls' or 'unencrypted'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracing := &Tracing{
				Enabled:     new(true),
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
			tlsCertPath: new("/etc/ssl/certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "relative file path",
			tlsCertPath: new("certs/cert.pem"),
			wantErr:     false,
		},
		{
			name:        "excessively long file path",
			tlsCertPath: new(strings.Repeat("z", 4097)),
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

func TestEthKeys_TOMLSerialization(t *testing.T) {
	t.Parallel()
	t.Run("encode", func(t *testing.T) {
		ethKeysWrapper := EthKeys{
			Keys: []*EthKey{
				{JSON: new(models.Secret("key1")), Password: new(models.Secret("pass1")), ID: new(1)},
				{JSON: new(models.Secret("key2")), Password: new(models.Secret("pass2")), ID: new(99)},
			},
		}

		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		err := enc.Encode(ethKeysWrapper)
		require.NoError(t, err)

		var decoded EthKeys
		err = toml.NewDecoder(strings.NewReader(buf.String())).Decode(&decoded)
		require.NoError(t, err)
		assert.Len(t, decoded.Keys, len(ethKeysWrapper.Keys))
		for i, key := range ethKeysWrapper.Keys {
			// have to compare the GoString() of the Secret because it is redacted
			assert.Equal(t, key.JSON.GoString(), decoded.Keys[i].JSON.GoString())
			assert.Equal(t, key.Password.GoString(), decoded.Keys[i].Password.GoString())
			assert.Equal(t, *key.ID, *decoded.Keys[i].ID)
		}
	})
	t.Run("decode", func(t *testing.T) {
		var decoded2 EthKeys
		btoml := `[[Keys]]
JSON = '{k:v}'
ID = 1337
Password = 'something'`
		err := toml.Unmarshal([]byte(btoml), &decoded2)
		require.NoError(t, err)
		assert.Len(t, decoded2.Keys, 1)
		assert.Equal(t, 1337, *decoded2.Keys[0].ID)
		assert.Equal(t, models.NewSecret("something"), decoded2.Keys[0].Password)
		assert.Equal(t, models.NewSecret("{k:v}"), decoded2.Keys[0].JSON)
	})
}

func TestSolKeys_TOMLSerialization(t *testing.T) {
	t.Parallel()

	t.Run("encode", func(t *testing.T) {
		solKeys := SolKeys{
			Keys: []*SolKey{
				{JSON: new(models.Secret("solkey1")), Password: new(models.Secret("pass1")), ID: new("devnet")},
				{JSON: new(models.Secret("solkey2")), Password: new(models.Secret("pass2")), ID: new("mainnet")},
			},
		}

		var buf bytes.Buffer
		enc := toml.NewEncoder(&buf)
		err := enc.Encode(solKeys)
		require.NoError(t, err)

		var decoded SolKeys
		err = toml.NewDecoder(strings.NewReader(buf.String())).Decode(&decoded)
		require.NoError(t, err)
		assert.Len(t, solKeys.Keys, len(decoded.Keys))
		for i, key := range solKeys.Keys {
			assert.Equal(t, key.JSON.GoString(), decoded.Keys[i].JSON.GoString())
			assert.Equal(t, key.Password.GoString(), decoded.Keys[i].Password.GoString())
			assert.Equal(t, *key.ID, *decoded.Keys[i].ID)
		}
	})

	t.Run("decode", func(t *testing.T) {
		var decoded SolKeys
		btoml := `[[Keys]]
JSON = '{k:v}'
ID = "devnet"
Password = 'secret'`
		err := toml.Unmarshal([]byte(btoml), &decoded)
		require.NoError(t, err)
		assert.Len(t, decoded.Keys, 1)
		assert.Equal(t, "devnet", *decoded.Keys[0].ID)
		assert.Equal(t, models.NewSecret("secret"), decoded.Keys[0].Password)
		assert.Equal(t, models.NewSecret("{k:v}"), decoded.Keys[0].JSON)
	})
}

func TestSolKeys_SetFrom(t *testing.T) {
	t.Parallel()

	solKeysWrapper1 := &SolKeys{}
	solKeysWrapper2 := SolKeys{
		Keys: []*SolKey{
			{
				JSON:     new(models.Secret("solkey1")),
				Password: new(models.Secret("pass1")),
				ID:       new("devnet"),
			},
		},
	}

	err := solKeysWrapper1.SetFrom(&solKeysWrapper2)
	require.NoError(t, err)
	assert.Equal(t, solKeysWrapper2, *solKeysWrapper1)
}

func TestEthKeys_SetFrom(t *testing.T) {
	ethKeysWrapper1 := &EthKeys{}
	ethKeysWrapper2 := EthKeys{
		Keys: []*EthKey{
			{JSON: new(models.Secret("key1")), Password: new(models.Secret("pass1")), ID: new(1)},
		},
	}

	err := ethKeysWrapper1.SetFrom(&ethKeysWrapper2)
	require.NoError(t, err)
	assert.Equal(t, ethKeysWrapper2, *ethKeysWrapper1)
}

func TestEthKeys_SetFrom_multipleSecretsFiles(t *testing.T) {
	t.Parallel()
	// Secrets files are applied in order and must union: a key from an earlier
	// -s file has to survive a later file that only carries other chains' keys.
	base := &EthKeys{Keys: []*EthKey{
		{JSON: new(models.Secret("key1")), Password: new(models.Secret("pass1")), ID: new(1)},
	}}
	disjoint := &EthKeys{Keys: []*EthKey{
		{JSON: new(models.Secret("key56")), Password: new(models.Secret("pass56")), ID: new(56)},
	}}

	require.NoError(t, base.SetFrom(disjoint))

	ids := make([]int, len(base.Keys))
	for _, k := range base.Keys {
		ids = append(ids, *k.ID)
	}
	assert.Equal(t, []int{1, 56}, ids, "keys from earlier secrets files must not be discarded")

	// Union must not weaken the no-overrides guarantee the -s flag documents.
	dupe := &EthKeys{Keys: []*EthKey{
		{JSON: new(models.Secret("other")), Password: new(models.Secret("otherpass")), ID: new(1)},
	}}
	require.Error(t, base.SetFrom(dupe))
	assert.Len(t, base.Keys, 2)
}

func TestEthKeys_validateMerge_nilID(t *testing.T) {
	// A secrets file may omit ID, and merging must survive it: a missing field
	// is a validation error, not a crash.
	base := &EthKeys{}
	noID := &EthKeys{Keys: []*EthKey{
		{JSON: new(models.Secret("key1")), Password: new(models.Secret("pass1"))},
	}}
	require.NotPanics(t, func() {
		require.NoError(t, base.SetFrom(noID))
	})
}

func TestSolKeys_SetFrom_multipleSecretsFiles(t *testing.T) {
	t.Parallel()
	base := &SolKeys{Keys: []*SolKey{
		{JSON: new(models.Secret("key1")), Password: new(models.Secret("pass1")), ID: new("devnet")},
	}}
	disjoint := &SolKeys{Keys: []*SolKey{
		{JSON: new(models.Secret("key2")), Password: new(models.Secret("pass2")), ID: new("mainnet")},
	}}

	require.NoError(t, base.SetFrom(disjoint))

	ids := make([]string, len(base.Keys))
	for _, k := range base.Keys {
		ids = append(ids, *k.ID)
	}
	assert.Equal(t, []string{"devnet", "mainnet"}, ids, "keys from earlier secrets files must not be discarded")

	dupe := &SolKeys{Keys: []*SolKey{
		{JSON: new(models.Secret("other")), Password: new(models.Secret("otherpass")), ID: new("devnet")},
	}}
	require.Error(t, base.SetFrom(dupe))
	assert.Len(t, base.Keys, 2)

	require.NotPanics(t, func() {
		noID := &SolKeys{Keys: []*SolKey{{JSON: new(models.Secret("k"))}}}
		require.NoError(t, (&SolKeys{}).SetFrom(noID))
	})
}

func TestAptosKeys_SetFrom_multipleSecretsFiles(t *testing.T) {
	t.Parallel()
	base := &AptosKeys{Keys: []*AptosKey{
		{JSON: new(models.Secret("key1")), Password: new(models.Secret("pass1")), ID: new(uint64(1))},
	}}
	disjoint := &AptosKeys{Keys: []*AptosKey{
		{JSON: new(models.Secret("key2")), Password: new(models.Secret("pass2")), ID: new(uint64(2))},
	}}

	require.NoError(t, base.SetFrom(disjoint))

	ids := make([]uint64, len(base.Keys))
	for _, k := range base.Keys {
		ids = append(ids, *k.ID)
	}
	assert.Equal(t, []uint64{1, 2}, ids, "keys from earlier secrets files must not be discarded")

	dupe := &AptosKeys{Keys: []*AptosKey{
		{JSON: new(models.Secret("other")), Password: new(models.Secret("otherpass")), ID: new(uint64(1))},
	}}
	require.Error(t, base.SetFrom(dupe))
	assert.Len(t, base.Keys, 2)

	require.NotPanics(t, func() {
		noID := &AptosKeys{Keys: []*AptosKey{{JSON: new(models.Secret("k"))}}}
		require.NoError(t, (&AptosKeys{}).SetFrom(noID))
	})
}

func TestStellarKeys_SetFrom_multipleSecretsFiles(t *testing.T) {
	t.Parallel()
	base := &StellarKeys{Keys: []*StellarKey{
		{JSON: new(commonconfig.SecretString("key1")), Password: new(commonconfig.SecretString("pass1")), ID: new("testnet")},
	}}
	disjoint := &StellarKeys{Keys: []*StellarKey{
		{JSON: new(commonconfig.SecretString("key2")), Password: new(commonconfig.SecretString("pass2")), ID: new("pubnet")},
	}}

	require.NoError(t, base.SetFrom(disjoint))

	ids := make([]string, len(base.Keys))
	for _, k := range base.Keys {
		ids = append(ids, *k.ID)
	}
	assert.Equal(t, []string{"testnet", "pubnet"}, ids, "keys from earlier secrets files must not be discarded")

	dupe := &StellarKeys{Keys: []*StellarKey{
		{JSON: new(commonconfig.SecretString("other")), Password: new(commonconfig.SecretString("otherpass")), ID: new("testnet")},
	}}
	require.Error(t, base.SetFrom(dupe))
	assert.Len(t, base.Keys, 2)

	require.NotPanics(t, func() {
		noID := &StellarKeys{Keys: []*StellarKey{{JSON: new(commonconfig.SecretString("k"))}}}
		require.NoError(t, (&StellarKeys{}).SetFrom(noID))
	})
}

// A key with only some of JSON/Password/ID set must be rejected, not silently
// accepted. All four key types share the "all fields must be nil or non-nil"
// rule.
func TestKeys_ValidateConfig_partialFields(t *testing.T) {
	secret := new(models.Secret("s"))
	stellarSecret := new(commonconfig.SecretString("s"))

	for _, tt := range []struct {
		name string
		cfg  interface{ ValidateConfig() error }
	}{
		{"EthKey missing ID", &EthKey{JSON: secret, Password: secret}},
		{"EthKey missing Password", &EthKey{JSON: secret, ID: new(1)}},
		{"SolKey missing ID", &SolKey{JSON: secret, Password: secret}},
		{"SolKey missing Password", &SolKey{JSON: secret, ID: new("devnet")}},
		{"AptosKey missing ID", &AptosKey{JSON: secret, Password: secret}},
		{"StellarKey missing ID", &StellarKey{JSON: stellarSecret, Password: stellarSecret}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			require.Error(t, tt.cfg.ValidateConfig())
		})
	}

	// All-nil remains valid: an absent key section is not an error.
	require.NoError(t, (&EthKey{}).ValidateConfig())
	require.NoError(t, (&SolKey{}).ValidateConfig())
	require.NoError(t, (&AptosKey{}).ValidateConfig())
	require.NoError(t, (&StellarKey{}).ValidateConfig())
}

func TestBridgeStatusReporter_ValidateConfig(t *testing.T) {
	testCases := []struct {
		name        string
		config      *BridgeStatusReporter
		expectError bool
		errorMsg    string
	}{
		{
			name: "disabled with nil fields",
			config: &BridgeStatusReporter{
				Enabled:              new(false),
				StatusPath:           nil,
				PollingInterval:      nil,
				IgnoreInvalidBridges: nil,
				IgnoreJoblessBridges: nil,
			},
			expectError: false,
		},
		{
			name: "disabled with empty fields",
			config: &BridgeStatusReporter{
				Enabled:              new(false),
				StatusPath:           new(""),
				PollingInterval:      durationPtr(0),
				IgnoreInvalidBridges: new(false),
				IgnoreJoblessBridges: new(true),
			},
			expectError: false,
		},
		{
			name: "disabled with valid fields",
			config: &BridgeStatusReporter{
				Enabled:              new(false),
				StatusPath:           new("/status"),
				PollingInterval:      durationPtr(5 * time.Minute),
				IgnoreInvalidBridges: new(true),
				IgnoreJoblessBridges: new(false),
			},
			expectError: false,
		},
		{
			name: "nil enabled (defaults to disabled)",
			config: &BridgeStatusReporter{
				Enabled:              nil,
				StatusPath:           new("/status"),
				PollingInterval:      durationPtr(5 * time.Minute),
				IgnoreInvalidBridges: new(true),
				IgnoreJoblessBridges: new(false),
			},
			expectError: false,
		},
		// Enabled valid cases with auto-defaulting
		{
			name: "enabled with valid config",
			config: &BridgeStatusReporter{
				Enabled:              new(true),
				StatusPath:           new("/status"),
				PollingInterval:      durationPtr(5 * time.Minute),
				IgnoreInvalidBridges: new(true),
				IgnoreJoblessBridges: new(false),
			},
			expectError: false,
		},
		{
			name: "enabled with nil fields - should fail validation",
			config: &BridgeStatusReporter{
				Enabled:              new(true),
				StatusPath:           nil,
				PollingInterval:      nil,
				IgnoreInvalidBridges: nil,
				IgnoreJoblessBridges: nil,
			},
			expectError: true,
			errorMsg:    "must be set",
		},
		{
			name: "enabled with empty status path - should auto-default",
			config: &BridgeStatusReporter{
				Enabled:              new(true),
				StatusPath:           new(""),
				PollingInterval:      durationPtr(5 * time.Minute),
				IgnoreInvalidBridges: new(true),
				IgnoreJoblessBridges: new(false),
			},
			expectError: false,
		},
		{
			name: "enabled with zero polling interval - should fail validation",
			config: &BridgeStatusReporter{
				Enabled:              new(true),
				StatusPath:           new("/status"),
				PollingInterval:      durationPtr(0),
				IgnoreInvalidBridges: new(true),
				IgnoreJoblessBridges: new(false),
			},
			expectError: true,
			errorMsg:    "must be greater than or equal to: 1m",
		},
		{
			name: "enabled with polling interval less than 1 minute - should fail validation",
			config: &BridgeStatusReporter{
				Enabled:              new(true),
				StatusPath:           new("/status"),
				PollingInterval:      durationPtr(30 * time.Second),
				IgnoreInvalidBridges: new(true),
				IgnoreJoblessBridges: new(false),
			},
			expectError: true,
			errorMsg:    "must be greater than or equal to: 1m",
		},
		{
			name: "enabled with polling interval exactly 1 minute",
			config: &BridgeStatusReporter{
				Enabled:              new(true),
				StatusPath:           new("/status"),
				PollingInterval:      durationPtr(1 * time.Minute),
				IgnoreInvalidBridges: new(true),
				IgnoreJoblessBridges: new(false),
			},
			expectError: false,
		},
		{
			name: "enabled with all fields missing - should fail validation",
			config: &BridgeStatusReporter{
				Enabled:              new(true),
				StatusPath:           new(""),
				PollingInterval:      durationPtr(0),
				IgnoreInvalidBridges: nil,
				IgnoreJoblessBridges: nil,
			},
			expectError: true,
			errorMsg:    "must be greater than or equal to: 1m",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.ValidateConfig()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)

				// Verify defaults are set when enabled
				if tc.config.Enabled != nil && *tc.config.Enabled {
					assert.NotNil(t, tc.config.StatusPath)
					assert.NotEmpty(t, *tc.config.StatusPath)
					assert.NotNil(t, tc.config.PollingInterval)
					assert.GreaterOrEqual(t, tc.config.PollingInterval.Duration(), time.Minute)
					assert.NotNil(t, tc.config.IgnoreInvalidBridges)
					assert.NotNil(t, tc.config.IgnoreJoblessBridges)
				}
			}
		})
	}
}

func durationPtr(d time.Duration) *commonconfig.Duration {
	cd := *commonconfig.MustNewDuration(d)
	return &cd
}

func TestMetering_ValidateConfig(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		config      *Metering
		expectError bool
		errorMsg    string
	}{
		{
			name:        "disabled with all nil fields",
			config:      &Metering{},
			expectError: false,
		},
		{
			name: "records enabled with non-empty NodeID",
			config: &Metering{
				MeterRecordsEnabled: new(true),
				NodeID:              new("clp-cre-wf-zone-a-1"),
			},
			expectError: false,
		},
		{
			name: "records enabled with nil NodeID",
			config: &Metering{
				MeterRecordsEnabled: new(true),
				NodeID:              nil,
			},
			expectError: true,
			errorMsg:    "NodeID",
		},
		{
			name: "records enabled with empty NodeID",
			config: &Metering{
				MeterRecordsEnabled: new(true),
				NodeID:              new(""),
			},
			expectError: true,
			errorMsg:    "NodeID",
		},
		{
			name: "snapshots enabled without records enabled",
			config: &Metering{
				MeterSnapshotsEnabled: new(true),
				NodeID:                new("clp-cre-wf-zone-a-1"),
			},
			expectError: true,
			errorMsg:    "requires MeterRecordsEnabled to be true",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.config.ValidateConfig()
			if tc.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
