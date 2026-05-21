package workflow

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// WriteArtifactTarball writes an uncompressed tar archive containing the given regular files
// at the archive root using each file's base name (flat layout). Empty paths are skipped.
func WriteArtifactTarball(dstPath string, filePaths []string) error {
	out, cErr := os.Create(dstPath)
	if cErr != nil {
		return errors.Wrap(cErr, "failed to create tarball")
	}
	defer out.Close()

	tw := tar.NewWriter(out)
	defer tw.Close()

	for _, p := range filePaths {
		if p == "" {
			continue
		}
		fi, err := os.Stat(p)
		if err != nil {
			return errors.Wrapf(err, "stat %q", p)
		}
		if !fi.Mode().IsRegular() {
			continue
		}

		hdr, err := tar.FileInfoHeader(fi, filepath.Base(p))
		if err != nil {
			return errors.Wrapf(err, "tar header for %q", p)
		}
		hdr.Name = filepath.Base(p)
		if wErr := tw.WriteHeader(hdr); wErr != nil {
			return errors.Wrapf(wErr, "write tar header for %q", p)
		}

		f, err := os.Open(p)
		if err != nil {
			return errors.Wrapf(err, "open %q", p)
		}
		if _, cErr := io.Copy(tw, f); cErr != nil {
			f.Close()
			return errors.Wrapf(cErr, "copy %q into tarball", p)
		}
		if cErr := f.Close(); cErr != nil {
			return errors.Wrapf(cErr, "close %q", p)
		}
	}

	return nil
}
