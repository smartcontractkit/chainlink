package v2

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink/cfgtest"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
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
