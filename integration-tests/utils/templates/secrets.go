package templates

import (
	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/templates"
)

// NodeSecretsTemplate are used as text templates because of secret redacted fields of chainlink.Secrets
// secret fields can't be marshalled as a plain text
type NodeSecretsTemplate struct {
	PgDbName      string
	PgHost        string
	PgPort        string
	PgPassword    string
	CustomSecrets string
}

func (c NodeSecretsTemplate) String() (string, error) {
	tpl := `
[Database]
URL = 'postgresql://postgres:{{ .PgPassword }}@{{ .PgHost }}:{{ .PgPort }}/{{ .PgDbName }}?sslmode=disable' # Required

[Password]
Keystore = '................' # Required

{{ if .CustomSecrets }}
	{{ .CustomSecrets }}
{{ else }}
[Mercury.Credentials.cred1]
URL = 'localhost:1338'
Username = 'node'
Password = 'nodepass'
{{ end }}
`
	return templates.MarshalTemplate(c, uuid.NewString(), tpl)
}
