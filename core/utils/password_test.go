package utils_test

import (
	"os"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

func TestVerifyPasswordComplexity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		password       string
		mustNotcontain string
		errors         []error
	}{
		{"thispasswordislongenough", "", []error{}},
		{"exactlyrightlen1", "", []error{}},
		{"notlongenough", "", []error{errors.New("password is less than 16 characters long")}},
		{"whitespace in password is ok", "", []error{}},
		{"\t leading whitespace not ok", "", []error{utils.ErrWhitespace}},
		{"trailing whitespace not ok\n", "", []error{utils.ErrWhitespace}},
		{"contains bad string", "bad", []error{errors.New("password may not contain: \"bad\"")}},
		{"contains bAd string 2", "bad", []error{errors.New("password may not contain: \"bad\"")}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.password, func(t *testing.T) {
			t.Parallel()

			var disallowedStrings []string
			if test.mustNotcontain != "" {
				disallowedStrings = []string{test.mustNotcontain}
			}
			err := utils.VerifyPasswordComplexity(test.password, disallowedStrings...)
			if len(test.errors) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, utils.ErrMsgHeader)
				for _, subErr := range test.errors {
					assert.ErrorContains(t, err, subErr.Error())
				}
			}
		})
	}
}

func TestPasswordFromFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		password string
		err      error
	}{
		{"", utils.ErrEmptyPasswordInFile},
		{" has whitespace  ", utils.ErrPasswordWhitespace},
		{"reasonable_password", nil},
	}

	for _, test := range tests {
		test := test
		t.Run(test.password, func(t *testing.T) {
			t.Parallel()

			pwdFile, err := os.CreateTemp("", "")
			assert.NoError(t, err)
			defer os.Remove(pwdFile.Name())
			_, err = pwdFile.WriteString(test.password)
			assert.NoError(t, err)

			pwd, err := utils.PasswordFromFile(pwdFile.Name())
			if test.err != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, test.err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, pwd, test.password)
			}
		})
	}
}
