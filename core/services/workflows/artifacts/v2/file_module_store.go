package v2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	defaultCacheSubdir    = "workflow-module-cache"
	binaryFileName        = "binary.wasm"
	engineVersionFileName = "engine_version.txt"
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
	if err := checkCacheDirWritable(cacheDir); err != nil {
		return nil, err
	}
	return &FileModuleStore{cacheDir: cacheDir}, nil
}

func checkCacheDirWritable(dir string) error {
	f, err := os.CreateTemp(dir, ".writecheck-*")
	if err != nil {
		return fmt.Errorf("module cache directory is not writable: %w", err)
	}
	name := f.Name()
	defer os.Remove(name)
	if _, err := f.Write([]byte{0}); err != nil {
		_ = f.Close()
		return fmt.Errorf("module cache directory is not writable: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("module cache directory is not writable: %w", err)
	}
	return nil
}

// CacheDir returns the resolved on-disk root for the module cache.
func (s *FileModuleStore) CacheDir() string { return s.cacheDir }

func (s *FileModuleStore) workflowDir(workflowID string) string {
	return filepath.Join(s.cacheDir, workflowID)
}

func (s *FileModuleStore) StoreModule(workflowID string, module []byte, engineVersion string) error {
	dir := s.workflowDir(workflowID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create workflow cache directory: %w", err)
	}

	if err := atomicWrite(filepath.Join(dir, binaryFileName), module); err != nil {
		return fmt.Errorf("failed to write module binary: %w", err)
	}
	if err := atomicWrite(filepath.Join(dir, engineVersionFileName), []byte(engineVersion)); err != nil {
		return fmt.Errorf("failed to write engine version: %w", err)
	}
	return nil
}

// GetModule returns the on-disk path of the cached binary together with the engine version
// recorded at write time. ok is false only when no binary is cached for workflowID.
// A missing engine version file is treated as empty string (legacy cache entries).
func (s *FileModuleStore) GetModule(workflowID string) (string, string, bool, error) {
	dir := s.workflowDir(workflowID)
	binPath := filepath.Join(dir, binaryFileName)
	if _, err := os.Stat(binPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", "", false, nil
		}
		return "", "", false, fmt.Errorf("failed to stat module binary: %w", err)
	}
	verBytes, err := os.ReadFile(filepath.Join(dir, engineVersionFileName))
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", "", false, fmt.Errorf("failed to read engine version: %w", err)
	}
	return binPath, string(verBytes), true, nil
}

func (s *FileModuleStore) DeleteModule(workflowID string) error {
	dir := s.workflowDir(workflowID)
	if err := os.RemoveAll(dir); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to delete module cache: %w", err)
	}
	return nil
}

func atomicWrite(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}
