package jobs

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/smartcontractkit/chainlink-common/pkg/settings"
)

// CombineCRESettingsFiles wraps settings.CombineTOMLFiles in order to make copies of org files.
// The copies allow us to bridge a transition period where different versions normalize the IDs in different ways.
func CombineCRESettingsFiles(wd string, fs fs.FS) ([]byte, error) {
	if err := os.CopyFS(wd, fs); err != nil {
		return nil, fmt.Errorf("failed to copy files: %w", err)
	}
	if err := copyOrgFiles(wd); err != nil {
		return nil, fmt.Errorf("failed to make copies of Org files: %w", err)
	}

	return settings.CombineTOMLFiles(os.DirFS(wd))
}

// copyOrgFiles makes two copies of each org file. One with an `org_` prefix and one lower-case.
func copyOrgFiles(dir string) error {
	dir = filepath.Join(dir, "org")
	err := fs.WalkDir(os.DirFS(dir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := strings.TrimSuffix(d.Name(), ".toml")
		if strings.HasPrefix(name, "org_") {
			return fmt.Errorf("invalid org id %s: must not have org_ prefix", name)
		}
		prefixed := "org_" + name
		err = copyFile(dir, name, prefixed)
		if err != nil {
			return fmt.Errorf("failed to copy file %s to %s: %w", name, prefixed, err)
		}
		lower := strings.ToLower(name)
		if name != lower { // no need for a copy
			err = copyFile(dir, name, lower)
			if err != nil {
				return fmt.Errorf("failed to copy file %s to %s: %w", name, lower, err)
			}
		}

		return nil
	})
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("error walking %s: %w", dir, err)
	}

	return nil
}

func copyFile(dir, src, dst string) error {
	s, err := os.Open(filepath.Join(dir, src+".toml"))
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", src, err)
	}
	defer s.Close()
	d, err := os.Create(filepath.Join(dir, dst+".toml"))
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", dst, err)
	}
	defer d.Close()
	_, err = io.Copy(d, s)

	return err
}
