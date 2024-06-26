package evm

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
)

// 4 digit version prefix to match the goose versioning
//
//go:embed [0-9][0-9][0-9][0-9]_*.tmpl.sql
var embeddedTmplFS embed.FS

var MigrationRootDir = "."

type Cfg struct {
	Schema  string
	ChainID *big.Big
}

func (c Cfg) Validate() error {
	if c.Schema == "" {
		return fmt.Errorf("schema is required")
	}
	if c.ChainID == nil {
		return fmt.Errorf("chain id is required")
	}
	return nil
}

var migrationSuffix = ".tmpl.sql"

func resolve(out io.Writer, in string, val Cfg) error {
	if err := val.Validate(); err != nil {
		return err
	}
	id := fmt.Sprintf("init_%s_%s", val.Schema, val.ChainID)
	tmpl, err := template.New(id).Option("missingkey=error").Parse(in)
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", in, err)
	}
	err = tmpl.Execute(out, val)
	if err != nil {
		return fmt.Errorf("failed to execute template %s: %w", in, err)
	}
	return nil
}

func generateMigrations(fsys fs.FS, rootDir string, tmpDir string, val Cfg) ([]string, error) {
	err := os.MkdirAll(tmpDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	b := make([]byte, 1024*1024)
	var migrations = []string{}
	var resolverFunc fs.WalkDirFunc
	resolverFunc = func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, migrationSuffix) {
			return nil
		}
		f, err := fsys.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file %s: %w", path, err)
		}
		defer f.Close()
		n, err := f.Read(b)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}
		content := b[:n]
		outPath := filepath.Join(tmpDir, strings.Replace(filepath.Base(path), migrationSuffix, ".sql", 1))
		out, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", outPath, err)
		}
		defer out.Close()
		err = resolve(out, string(content), val)
		if err != nil {
			return fmt.Errorf("failed to resolve template %s: %w", path, err)
		}
		migrations = append(migrations, outPath)
		return nil
	}

	err = fs.WalkDir(fsys, rootDir, resolverFunc)
	if err != nil {
		return nil, err
	}
	return migrations, nil
}
