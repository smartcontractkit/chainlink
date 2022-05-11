package utils_test

import (
	"testing"

	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/stretchr/testify/assert"
)

func TestVerifyPasswordComplexity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		password string
		errors   []error
	}{
		{"QWERTYuiop123!@#", []error{}},
		{"QQQQWERTYuiop123!@#", []error{utils.ErrPasswordRepeatedChars}},
		{"abcB123+!@", []error{utils.ErrPasswordMinUppercase}},
		{"ABCd123+!@", []error{utils.ErrPasswordMinLowercase}},
		{"ABCzxc1+!@", []error{utils.ErrPasswordMinNumbers}},
		{"aB2+", []error{
			utils.ErrPasswordMinLength,
			utils.ErrPasswordMinUppercase,
			utils.ErrPasswordMinLowercase,
			utils.ErrPasswordMinNumbers,
		}},
	}

	for _, test := range tests {
		test := test

		t.Run(test.password, func(t *testing.T) {
			t.Parallel()

			err := utils.VerifyPasswordComplexity(test.password)
			if len(test.errors) == 0 {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.ErrorContains(t, err, "password does not meet the requirements")
				for _, subErr := range test.errors {
					assert.ErrorContains(t, err, subErr.Error())
				}
			}
		})
	}
}
