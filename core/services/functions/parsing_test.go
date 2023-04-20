package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseDomains(t *testing.T) {
	t.Parallel()

	t.Run("empty source code", func(t *testing.T) {
		domains := parseDomains("")
		assert.Empty(t, domains)
	})

	t.Run("parse valid and ignore invalid domains", func(t *testing.T) {
		sourceCode := `
			const x = 1;
			const u1 = "http://foo.bar.io"; // => foo.bar.io
			const u2 = "https://github.com/user/xyz?p=1"; // => github.com
			const u3 = "https://myapi.net/user/${placeholder}"; // => myapi.net
			const b0 = "http://foo.bar.io/duplicate"; // => duplicate
			const b1 = "google.com"; // no protocol
		`
		domains := parseDomains(sourceCode)
		assert.EqualValues(t, []string{
			"foo.bar.io",
			"github.com",
			"myapi.net",
		}, domains)
	})
}
