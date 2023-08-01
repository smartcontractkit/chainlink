package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var NodeSecretsTemplate = `
[Database]
URL = 'postgresql://postgres:test@{{ .PgHost }}:{{ .PgPort }}/cl-node?sslmode=disable' # Required
AllowSimplePasswords = true

[Password]
Keystore = '................' # Required

[Mercury.Credentials.cred1]
# URL = 'http://host.docker.internal:3000/reports'
URL = 'localhost:1338'
Username = 'node'
Password = 'nodepass'
`

func ExecuteNodeSecretsTemplate(pgHost string, pgPort string) (string, error) {
	data := struct {
		PgHost string
		PgPort string
	}{
		PgHost: pgHost,
		PgPort: pgPort,
	}

	t, err := template.New("node-secrets").Parse(NodeSecretsTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)

	return buf.String(), err
}
