package txtar

import (
	"io/fs"
	"os"
	"path/filepath"
)

type RecurseOpt bool

const (
	Recurse   RecurseOpt = true
	NoRecurse RecurseOpt = false
)

type DirVisitor struct {
	rootDir string
	cb      func(path string) error
	recurse RecurseOpt
}

func (d *DirVisitor) Walk() error {
	return filepath.WalkDir(d.rootDir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !de.IsDir() {
			return nil
		}

		isRootDir, err := d.isRootDir(de)
		if err != nil {
			return err
		}

		// If we're not recursing, skip all other directories except the root.
		if !bool(d.recurse) && !isRootDir {
			return nil
		}

		matches, err := fs.Glob(os.DirFS(path), "*txtar")
		if err != nil {
			return err
		}

		if len(matches) > 0 {
			return d.cb(path)
		}

		return nil
	})
}

func (d *DirVisitor) isRootDir(de fs.DirEntry) (bool, error) {
	fi, err := os.Stat(d.rootDir)
	if err != nil {
		return false, err
	}

	fi2, err := de.Info()
	if err != nil {
		return false, err
	}
	return os.SameFile(fi, fi2), nil
}

func NewDirVisitor(rootDir string, recurse RecurseOpt, cb func(path string) error) *DirVisitor {
	return &DirVisitor{
		rootDir: rootDir,
		cb:      cb,
		recurse: recurse,
	}
}
