package v2

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultCacheSubdir = "workflow-module-cache"
	binaryFileName     = "binary.wasm"
)

type FileModuleStore struct {
	cacheDir string
}

// NewFileModuleStore opens the on-disk module cache rooted at cacheDir when non-empty,
// or at os.TempDir()/workflow-module-cache when cacheDir is empty.
// If cleanOnStartup is true, the resolved directory is removed first so the process starts
// with an empty cache (workflow registry sync repopulates it).
func NewFileModuleStore(cacheDir string, cleanOnStartup bool) (*FileModuleStore, error) {
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), defaultCacheSubdir)
	}
	if cleanOnStartup {
		if err := os.RemoveAll(cacheDir); err != nil {
			return nil, fmt.Errorf("failed to clear module cache directory: %w", err)
		}
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create module cache directory: %w", err)
	}
	return &FileModuleStore{cacheDir: cacheDir}, nil
}

func (s *FileModuleStore) workflowDir(workflowID string) string {
	return filepath.Join(s.cacheDir, workflowID)
}

func (s *FileModuleStore) StoreModule(workflowID string, module []byte) error {
	dir := s.workflowDir(workflowID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create workflow cache directory: %w", err)
	}

	binaryPath := filepath.Join(dir, binaryFileName)
	tmpBinary := binaryPath + ".tmp"
	if err := os.WriteFile(tmpBinary, module, 0o600); err != nil {
		return fmt.Errorf("failed to write module binary: %w", err)
	}
	if err := os.Rename(tmpBinary, binaryPath); err != nil {
		os.Remove(tmpBinary)
		return fmt.Errorf("failed to finalize module binary: %w", err)
	}

	return nil
}

func (s *FileModuleStore) GetModulePath(workflowID string) (string, bool, error) {
	p := filepath.Join(s.workflowDir(workflowID), binaryFileName)
	if _, err := os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to stat module binary: %w", err)
	}
	return p, true, nil
}

func (s *FileModuleStore) DeleteModule(workflowID string) error {
	dir := s.workflowDir(workflowID)
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete module cache: %w", err)
	}
	return nil
}
