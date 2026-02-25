package v2

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultCacheSubdir = "workflow-module-cache"
	binaryFileName     = "binary.wasm"
	binaryIDFileName   = "binary_id"
)

type FileModuleStore struct {
	cacheDir string
}

func NewFileModuleStore(cacheDir string) (*FileModuleStore, error) {
	if cacheDir == "" {
		cacheDir = filepath.Join(os.TempDir(), defaultCacheSubdir)
	}
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create module cache directory: %w", err)
	}
	return &FileModuleStore{cacheDir: cacheDir}, nil
}

func (s *FileModuleStore) workflowDir(workflowID string) string {
	return filepath.Join(s.cacheDir, workflowID)
}

func (s *FileModuleStore) StoreModule(workflowID string, binaryID string, module []byte) error {
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

	idPath := filepath.Join(dir, binaryIDFileName)
	tmpID := idPath + ".tmp"
	if err := os.WriteFile(tmpID, []byte(binaryID), 0o600); err != nil {
		return fmt.Errorf("failed to write binary ID: %w", err)
	}
	if err := os.Rename(tmpID, idPath); err != nil {
		os.Remove(tmpID)
		return fmt.Errorf("failed to finalize binary ID: %w", err)
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

func (s *FileModuleStore) GetBinaryID(workflowID string) (string, bool, error) {
	p := filepath.Join(s.workflowDir(workflowID), binaryIDFileName)
	data, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to read binary ID: %w", err)
	}
	return string(data), true, nil
}

func (s *FileModuleStore) DeleteModule(workflowID string) error {
	dir := s.workflowDir(workflowID)
	if err := os.RemoveAll(dir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete module cache: %w", err)
	}
	return nil
}
