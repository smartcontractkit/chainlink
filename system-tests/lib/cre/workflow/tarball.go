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
	out, err := os.Create(dstPath)
	if err != nil {
		return errors.Wrap(err, "failed to create tarball")
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
		if err := tw.WriteHeader(hdr); err != nil {
			return errors.Wrapf(err, "write tar header for %q", p)
		}

		f, err := os.Open(p)
		if err != nil {
			return errors.Wrapf(err, "open %q", p)
		}
		if _, err := io.Copy(tw, f); err != nil {
			f.Close()
			return errors.Wrapf(err, "copy %q into tarball", p)
		}
		if err := f.Close(); err != nil {
			return errors.Wrapf(err, "close %q", p)
		}
	}

	return nil
}
