package templates

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var NodeSecretsTemplate = `
[Database]
URL = 'postgresql://postgres:{{ .PgPassword }}@{{ .PgHost }}:{{ .PgPort }}/{{ .PgDbName }}?sslmode=disable' # Required

[Password]
Keystore = 'mysecretpassword' # Required

[Mercury.Credentials.cred1]
# URL = 'http://host.docker.internal:3000/reports'
URL = 'localhost:1338'
Username = 'node'
Password = 'nodepass'
`

func ExecuteNodeSecretsTemplate(pgPassowrd, pgDbName, pgHost, pgPort string) (string, error) {
	data := struct {
		PgDbName   string
		PgHost     string
		PgPort     string
		PgPassword string
	}{
		PgDbName:   pgDbName,
		PgHost:     pgHost,
		PgPort:     pgPort,
		PgPassword: pgPassowrd,
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
