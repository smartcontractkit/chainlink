package template

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type RelayerDB struct {
	Schema string
}

// resolve resolves the template with the given RelayerDB
func resolve(out io.Writer, in io.Reader, val RelayerDB) error {
	unresolved, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	tmpl, err := template.New("schema-resolver").Parse(string(unresolved))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", unresolved, err)
	}
	err = tmpl.Execute(out, val)
	return err
}

var migrationSuffix = ".tmpl.sql"

func generateMigrations(rootDir string, tmpDir string, val RelayerDB) ([]string, error) {
	err := os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	var migrations = []string{}
	var resolverFunc = func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, migrationSuffix) {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}
		outPath := filepath.Join(tmpDir, strings.Replace(filepath.Base(path), migrationSuffix, ".sql", 1))
		out, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outPath, err)
		}
		defer out.Close()
		err = resolve(out, bytes.NewBuffer(b), val)
		if err != nil {
			return fmt.Errorf("failed to resolve template %s: %w", path, err)
		}
		migrations = append(migrations, outPath)
		return nil
	}

	err = filepath.Walk(rootDir, resolverFunc)
	if err != nil {
		return nil, err
	}
	return migrations, nil
}
