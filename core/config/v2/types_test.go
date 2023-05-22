package v2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestCoreDefaults_notNil(t *testing.T) {
	cfgtest.AssertFieldsNotNil(t, &defaults)
}

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
