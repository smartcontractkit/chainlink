package config

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/test-go/testify/assert"
)

func Test_validateDBURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{"no user or password", "postgresql://foo.example.com:5432/chainlink?application_name=Test+Application", "DB URL must be authenticated; plaintext URLs are not allowed"},
		{"with user and no password", "postgresql://myuser@foo.example.com:5432/chainlink?application_name=Test+Application", "DB URL must be authenticated; password is required"},
		{"with user and password of insufficient length", "postgresql://myuser:shortpw@foo.example.com:5432/chainlink?application_name=Test+Application", "password does not meet the requirements: must be longer than 16 characters"},
		{"with no user and password of sufficient length", "postgresql://:thisisareallylongpassword@foo.example.com:5432/chainlink?application_name=Test+Application", ""},
		{"with user and password of sufficient length", "postgresql://myuser:thisisareallylongpassword@foo.example.com:5432/chainlink?application_name=Test+Application", ""},
		{"with user and password of insufficient length as params", "postgresql://foo.example.com:5432/chainlink?application_name=Test+Application&password=shortpw&user=myuser", "password does not meet the requirements: must be longer than 16 characters"},
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
