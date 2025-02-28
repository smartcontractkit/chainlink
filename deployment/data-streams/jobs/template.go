package jobs

import (
	"embed"
	"math/big"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

func newTemplate() (*template.Template, error) {
	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"times": func(i int64) string {
			t := big.NewInt(10)
			t.Exp(t, big.NewInt(i), nil)
			return t.String()
		},
	}
	return template.New("ds").Funcs(funcMap).ParseFS(templatesFS, "templates/*.tmpl")
}

func renderTemplate(fname string, data any) (string, error) {
	tmpl, err := newTemplate()
	if err != nil {
		return "", err
	}
	b := new(strings.Builder)
	err = tmpl.ExecuteTemplate(b, fname, data)
	return b.String(), err
}
