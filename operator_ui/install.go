package main

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	const (
		owner                  = "smartcontractkit"
		repo                   = "operator-ui"
		fullRepo               = owner + "/" + repo
		tagPath                = "operator_ui/TAG"
		unpackDir              = "core/web/assets"
		downloadTimeoutSeconds = 30
	)
	// Grab first argument as root directory
	if len(os.Args) < 2 {
		log.Fatalln("Usage: install.go <root>")
	}
	rootDir := os.Args[1]

	tag := mustReadTagFile(path.Join(rootDir, tagPath))
	strippedTag := stripVersionFromTag(tag)
	assetName := fmt.Sprintf("%s-%s-%s.tgz", owner, repo, strippedTag)
	downloadUrl := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", fullRepo, tag, assetName)

	// Assuming that we're in "root/operator_ui/"
	unpackPath := filepath.Join(rootDir, unpackDir)
	err := rmrf(unpackPath)
	if err != nil {
		log.Fatalln(err)
	}

	subPath := "package/artifacts/"
	mustDownloadSubAsset(downloadUrl, downloadTimeoutSeconds, unpackPath, subPath)
}

func mustReadTagFile(file string) string {
	tagBytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalln(err)
	}
	tag := string(tagBytes)
	return strings.TrimSpace(tag)
}

func stripVersionFromTag(tag string) string {
	return strings.TrimPrefix(tag, "v")
}

func rmrf(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	err = os.Mkdir(path, 0755)
	return err
}

// Download a sub asset from a .tgz file and extract it to a destination path
func mustDownloadSubAsset(downloadUrl string, downloadTimeoutSeconds int, unpackPath string, subPath string) {
	fmt.Println("Downloading", downloadUrl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(downloadTimeoutSeconds)*time.Second)
	defer cancel()
	/* #nosec G107 */
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalln(fmt.Errorf("failed to fetch asset: %s", resp.Status))
	}

	err = decompressTgzSubpath(resp.Body, unpackPath, subPath)
	if err != nil {
		log.Fatalln(err)
	}
}

// Decompress a .tgz file to a destination path, only extracting files that are in the subpath
//
// Subpath files are extracted to the root of the destination path, rather than preserving the subpath
func decompressTgzSubpath(file io.Reader, destPath string, subPath string) error {
	// Create a gzip reader
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzr.Close()

	// Create a tar reader
	tr := tar.NewReader(gzr)

	// Iterate through the files in the tar archive
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil // End of tar archive
		case err != nil:
			return fmt.Errorf("failed to read tar file: %w", err)
		case header == nil:
			continue
		}
		// skip files that arent in the subpath
		if !strings.HasPrefix(header.Name, subPath) {
			continue
		}

		// Strip the subpath from the header name
		header.Name = strings.TrimPrefix(header.Name, subPath)

		// Target location where the dir/file should be created
		target := fmt.Sprintf("%s/%s", destPath, header.Name)

		// Check the file type
		switch header.Typeflag {
		case tar.TypeDir: // Directory
			if err := os.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			fmt.Println("Creating directory", target)
		case tar.TypeReg: // Regular file
			if err := writeFile(target, header, tr); err != nil {
				return err
			}
		}
	}
}

func writeFile(target string, header *tar.Header, tr *tar.Reader) error {
	/* #nosec G110 */
	f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, tr); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	fmt.Println("Creating file", target)
	return nil
}
