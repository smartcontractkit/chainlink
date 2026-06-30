package txtar

import (
	"io/fs"
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
	root := filepath.Clean(d.rootDir)
	return filepath.WalkDir(root, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !de.IsDir() {
			return nil
		}

		if !bool(d.recurse) && filepath.Clean(path) != root {
			return fs.SkipDir
		}

		matches, err := filepath.Glob(filepath.Join(path, "*txtar"))
		if err != nil {
			return err
		}

		if len(matches) > 0 {
			if err := d.cb(path); err != nil {
				return err
			}
		}

		if !bool(d.recurse) {
			return fs.SkipDir
		}

		return nil
	})
}

func NewDirVisitor(rootDir string, recurse RecurseOpt, cb func(path string) error) *DirVisitor {
	return &DirVisitor{
		rootDir: rootDir,
		cb:      cb,
		recurse: recurse,
	}
}
