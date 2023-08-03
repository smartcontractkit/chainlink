package templates

import (
	"github.com/google/uuid"
)

type NodeSecretsTemplate struct {
	PgDbName string
	PgHost   string
	PgPort   string
}

func (c NodeSecretsTemplate) String() (string, error) {
	tpl := `
[Database]
URL = 'postgresql://postgres:test@{{ .PgHost }}:{{ .PgPort }}/{{ .PgDbName }}?sslmode=disable' # Required
AllowSimplePasswords = true

[Password]
Keystore = '................' # Required

[Mercury.Credentials.cred1]
# URL = 'http://host.docker.internal:3000/reports'
URL = 'localhost:1338'
Username = 'node'
Password = 'nodepass'
`
	return MarshallTemplate(c, uuid.NewString(), tpl)
}
