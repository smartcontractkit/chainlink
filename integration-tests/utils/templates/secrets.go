package templates

import "github.com/google/uuid"

// NodeSecretsTemplate are used as text templates because of secret redacted fields of chainlink.Secrets
// secret fields can't be marshalled as a plain text
type NodeSecretsTemplate struct {
	PgDbName   string
	PgHost     string
	PgPort     string
	PgPassword string
}

func (c NodeSecretsTemplate) String() (string, error) {
	tpl := `
[Database]
URL = 'postgresql://postgres:{{ .PgPassword }}@{{ .PgHost }}:{{ .PgPort }}/{{ .PgDbName }}?sslmode=disable' # Required

[Password]
Keystore = '................' # Required

[Mercury.Credentials.cred1]
# URL = 'http://host.docker.internal:3000/reports'
URL = 'localhost:1338'
Username = 'node'
Password = 'nodepass'
`
	return MarshalTemplate(c, uuid.NewString(), tpl)
}
