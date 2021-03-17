package periodicbackup

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/static"
)

var (
	filePattern        = "cl_backup_%s.tar.gz"
	minBackupFrequency = time.Minute
)

type backupResult struct {
	size int64
	path string
}

type (
	DatabaseBackup interface {
		Start() error
		Close() error
		RunBackupGracefully(version string)
	}

	databaseBackup struct {
		logger          *logger.Logger
		databaseURL     url.URL
		frequency       time.Duration
		outputParentDir string
		done            chan bool
	}

	Config interface {
		DatabaseBackupFrequency() time.Duration
		DatabaseBackupURL() *url.URL
		DatabaseURL() url.URL
		RootDir() string
	}
)

func NewDatabaseBackup(config Config, logger *logger.Logger) DatabaseBackup {

	dbUrl := config.DatabaseURL()
	dbBackupUrl := config.DatabaseBackupURL()
	if dbBackupUrl != nil {
		dbUrl = *dbBackupUrl
	}
	return &databaseBackup{
		logger,
		dbUrl,
		config.DatabaseBackupFrequency(),
		config.RootDir(),
		make(chan bool),
	}
}

func (backup databaseBackup) Start() error {

	if backup.frequencyIsTooSmall() {
		return errors.Errorf("Database backup frequency (%s=%v) is too small. Please set it to at least %s", "DATABASE_BACKUP_FREQUENCY", backup.frequency, minBackupFrequency)
	}

	ticker := time.NewTicker(backup.frequency)

	go func() {
		for {
			select {
			case <-backup.done:
				ticker.Stop()
				return
			case <-ticker.C:
				backup.RunBackupGracefully(static.Version)
			}
		}
	}()

	return nil
}

func (backup databaseBackup) Close() error {
	backup.done <- true
	return nil
}

func (backup *databaseBackup) frequencyIsTooSmall() bool {
	return backup.frequency < minBackupFrequency
}

func (backup *databaseBackup) RunBackupGracefully(version string) {
	backup.logger.Info("DatabaseBackup: Starting database backup...")
	startAt := time.Now()
	result, err := backup.runBackup(version)
	duration := time.Since(startAt)
	if err != nil {
		backup.logger.Errorw("DatabaseBackup: Failed", "duration", duration, "error", err)
	} else {
		backup.logger.Infow("DatabaseBackup: Database backup finished successfully.", "duration", duration, "fileSize", result.size, "filePath", result.path)
	}
}

func (backup *databaseBackup) runBackup(version string) (*backupResult, error) {

	tmpFile, err := ioutil.TempFile(backup.outputParentDir, "db_backup")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a tmp file")
	}
	err = os.Remove(tmpFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to remove the tmp file before running backup")
	}

	cmd := exec.Command(
		"pg_dump", backup.databaseURL.String(),
		"-f", tmpFile.Name(),
		"-F", "t", // format: tar
	)

	_, err = cmd.Output()

	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			return nil, errors.Wrap(err, fmt.Sprintf("pg_dump failed with output: %s", string(ee.Stderr)))
		}
		return nil, errors.Wrap(err, "pg_dump failed")
	}

	finalFilePath := filepath.Join(backup.outputParentDir, fmt.Sprintf(filePattern, version))
	_ = os.Remove(finalFilePath)
	err = os.Rename(tmpFile.Name(), finalFilePath)
	if err != nil {
		_ = os.Remove(tmpFile.Name())
		return nil, errors.Wrap(err, "Failed to rename the temp file to the final backup file")
	}

	file, err := os.Stat(finalFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to access the final backup file")
	}

	return &backupResult{
		size: file.Size(),
		path: finalFilePath,
	}, nil
}
