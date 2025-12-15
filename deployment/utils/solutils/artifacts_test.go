package solutils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestDownloadProgramArtifacts(t *testing.T) {
	tests := []struct {
		name        string
		setupServer func() *httptest.Server
		wantFiles   []string
		wantErr     string
	}{
		{
			name: "successful download and extraction",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Create a test tar.gz with some files
					w.Header().Set("Content-Type", "application/gzip")

					gzWriter := gzip.NewWriter(w)
					defer gzWriter.Close()

					tarWriter := tar.NewWriter(gzWriter)
					defer tarWriter.Close()

					// Add test files to the tar
					testFiles := map[string]string{
						"program1.so": "fake program 1 content",
						"program2.so": "fake program 2 content",
						"config.json": `{"test": "config"}`,
					}

					for filename, content := range testFiles {
						header := &tar.Header{
							Name:     filename,
							Size:     int64(len(content)),
							Typeflag: tar.TypeReg,
						}

						err := tarWriter.WriteHeader(header)
						if err != nil {
							t.Errorf("Failed to write tar header: %v", err)
							return
						}

						_, err = tarWriter.Write([]byte(content))
						if err != nil {
							t.Errorf("Failed to write tar content: %v", err)
							return
						}
					}
				}))
			},
			wantFiles: []string{"program1.so", "program2.so", "config.json"},
		},
		{
			name: "server returns 404",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			wantErr: "download failed with status 404",
		},
		{
			name: "server returns 500",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			wantErr: "download failed with status 500",
		},
		{
			name: "invalid gzip content",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/gzip")
					_, err := w.Write([]byte("invalid gzip content"))
					if err != nil {
						t.Errorf("Failed to write gzip content: %v", err)
						return
					}
				}))
			},
			wantErr: "gzip",
		},
		{
			name: "empty tar archive",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/gzip")

					gzWriter := gzip.NewWriter(w)
					defer gzWriter.Close()

					tarWriter := tar.NewWriter(gzWriter)
					defer tarWriter.Close()
					// Don't add any files - empty archive
				}))
			},
			wantFiles: []string{}, // No files expected
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			server := tt.setupServer()
			defer server.Close()

			// Create temporary directory for extraction
			tempDir := t.TempDir()

			// Execute
			err := downloadProgramArtifacts(
				t.Context(), server.URL, tempDir, logger.Test(t),
			)

			// Assert
			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)

				// Check that expected files were created
				for _, expectedFile := range tt.wantFiles {
					filePath := filepath.Join(tempDir, expectedFile)
					assert.FileExists(t, filePath, "Expected file %s to exist", expectedFile)

					// Verify file is not empty (except for empty archive test)
					if len(tt.wantFiles) > 0 {
						info, err := os.Stat(filePath)
						require.NoError(t, err)
						assert.Positive(t, info.Size(), "File %s should not be empty", expectedFile)
					}
				}

				// Check that no unexpected files were created
				entries, err := os.ReadDir(tempDir)
				require.NoError(t, err)
				assert.Len(t, entries, len(tt.wantFiles), "Unexpected number of files extracted")
			}
		})
	}
}

func TestDownloadProgramArtifacts_ContextCancellation(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow server
		select {
		case <-r.Context().Done():
			return
		case <-make(chan struct{}):
			// This will never be reached due to context cancellation
		}
	}))
	defer server.Close()

	// Create a context that gets cancelled immediately
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // Cancel immediately

	err := downloadProgramArtifacts(ctx, server.URL, t.TempDir(), logger.Test(t))
	require.Error(t, err)
	require.ErrorContains(t, err, "context canceled")
}

func TestDownloadProgramArtifacts_InvalidURL(t *testing.T) {
	tempDir := t.TempDir()

	err := downloadProgramArtifacts(t.Context(), "http://invalid-url", tempDir, logger.Test(t))
	require.ErrorContains(t, err, "dial tcp: lookup invalid-url: no such host")
}

func TestDownloadProgramArtifacts_NonExistentTargetDir(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/gzip")

		gzWriter := gzip.NewWriter(w)
		defer gzWriter.Close()

		tarWriter := tar.NewWriter(gzWriter)
		defer tarWriter.Close()

		// Add a test file
		content := "test content"
		header := &tar.Header{
			Name:     "test.so",
			Size:     int64(len(content)),
			Typeflag: tar.TypeReg,
		}

		err := tarWriter.WriteHeader(header)
		if err != nil {
			t.Errorf("Failed to write tar header: %v", err)
			return
		}
		_, err = tarWriter.Write([]byte(content))
		if err != nil {
			t.Errorf("Failed to write tar content: %v", err)
			return
		}
	}))
	defer server.Close()

	// Use a non-existent directory path
	nonExistentDir := "/tmp/non_existent_parent_dir_12345/target"

	err := downloadProgramArtifacts(t.Context(), server.URL, nonExistentDir, logger.Test(t))
	require.NoError(t, err) // Should succeed because MkdirAll creates parent directories

	// Verify the file was created
	assert.FileExists(t, filepath.Join(nonExistentDir, "test.so"))

	// Cleanup
	os.RemoveAll("/tmp/non_existent_parent_dir_12345")
}
