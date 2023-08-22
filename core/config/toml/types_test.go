package toml

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/build"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
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

func TestValidateConfig(t *testing.T) {
	validUrl := models.URL(url.URL{Scheme: "https", Host: "localhost"})
	validSecretURL := *models.NewSecretURL(&validUrl)

	invalidEmptyUrl := models.URL(url.URL{})
	invalidEmptySecretURL := *models.NewSecretURL(&invalidEmptyUrl)

	invalidBackupURL := models.URL(url.URL{Scheme: "http", Host: "localhost"})
	invalidBackupSecretURL := *models.NewSecretURL(&invalidBackupURL)

	tests := []struct {
		name                string
		input               *DatabaseSecrets
		skip                bool
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
			skip:                !build.IsProd(),
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
			// needed while -tags test is supported
			if tt.skip {
				t.SkipNow()
			}
			err := tt.input.ValidateConfig()
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
