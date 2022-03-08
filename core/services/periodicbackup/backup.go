package periodicbackup

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	filePattern        = "cl_backup_%s.dump"
	minBackupFrequency = time.Minute

	excludedDataFromTables = []string{
		"pipeline_runs",
		"pipeline_task_runs",
	}
)

type backupResult struct {
	size            int64
	path            string
	maskedArguments []string
	pgDumpArguments []string
}

type (
	DatabaseBackup interface {
		services.ServiceCtx
		RunBackup(version string) error
	}

	databaseBackup struct {
		logger          logger.Logger
		databaseURL     url.URL
		mode            config.DatabaseBackupMode
		frequency       time.Duration
		outputParentDir string
		done            chan bool
		utils.StartStopOnce
	}

	Config interface {
		DatabaseBackupMode() config.DatabaseBackupMode
		DatabaseBackupFrequency() time.Duration
		DatabaseBackupURL() *url.URL
		DatabaseBackupDir() string
		DatabaseURL() url.URL
		RootDir() string
	}
)

// NewDatabaseBackup instantiates a *databaseBackup
func NewDatabaseBackup(config Config, lggr logger.Logger) (DatabaseBackup, error) {
	lggr = lggr.Named("DatabaseBackup")
	dbUrl := config.DatabaseURL()
	dbBackupUrl := config.DatabaseBackupURL()
	if dbBackupUrl != nil {
		dbUrl = *dbBackupUrl
	}

	outputParentDir := filepath.Join(config.RootDir(), "backup")
	if config.DatabaseBackupDir() != "" {
		dir, err := filepath.Abs(config.DatabaseBackupDir())
		if err != nil {
			return nil, errors.Errorf("failed to get path for DATABASE_BACKUP_DIR (%s) - please set it to a valid directory path", config.DatabaseBackupDir())
		}
		outputParentDir = dir
	}

	return &databaseBackup{
		lggr,
		dbUrl,
		config.DatabaseBackupMode(),
		config.DatabaseBackupFrequency(),
		outputParentDir,
		make(chan bool),
		utils.StartStopOnce{},
	}, nil
}

// Start starts DatabaseBackup.
func (backup *databaseBackup) Start(context.Context) error {
	return backup.StartOnce("DatabaseBackup", func() (err error) {
		ticker := time.NewTicker(backup.frequency)
		if backup.frequency == 0 {
			backup.logger.Info("Periodic database backups are disabled; DATABASE_BACKUP_FREQUENCY was set to 0")
			// Stopping the ticker means it will never fire, effectively disabling periodic backups
			ticker.Stop()
		} else if backup.frequencyIsTooSmall() {
			return errors.Errorf("Database backup frequency (%s=%v) is too small. Please set it to at least %s (or set to 0 to disable periodic backups)", "DATABASE_BACKUP_FREQUENCY", backup.frequency, minBackupFrequency)
		}

		go func() {
			for {
				select {
				case <-backup.done:
					ticker.Stop()
					return
				case <-ticker.C:
					backup.logger.Infow("Starting automatic database backup, this can take a while. To disable periodic backups, set DATABASE_BACKUP_FREQUENCY=0. To disable database backups entirely, set DATABASE_BACKUP_MODE=none.")
					//nolint:errcheck
					backup.RunBackup(static.Version)
				}
			}
		}()

		return nil
	})
}

func (backup *databaseBackup) Close() error {
	return backup.StopOnce("DatabaseBackup", func() (err error) {
		backup.done <- true
		return nil
	})
}

func (backup *databaseBackup) frequencyIsTooSmall() bool {
	return backup.frequency < minBackupFrequency
}

func (backup *databaseBackup) RunBackup(version string) error {
	backup.logger.Debugw("Starting backup", "mode", backup.mode, "url", backup.databaseURL.Redacted(), "directory", backup.outputParentDir)
	startAt := time.Now()
	result, err := backup.runBackup(version)
	duration := time.Since(startAt)
	if err != nil {
		backup.logger.Errorw("Backup failed", "duration", duration, "err", err)
		return err
	}
	backup.logger.Infow("Backup completed successfully.", "duration", duration, "fileSize", result.size, "filePath", result.path)
	return nil
}

func (backup *databaseBackup) runBackup(version string) (*backupResult, error) {
	err := os.MkdirAll(backup.outputParentDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create directories on the path: %s", backup.outputParentDir)
	}
	tmpFile, err := ioutil.TempFile(backup.outputParentDir, "cl_backup_tmp_")
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create a tmp file")
	}
	err = os.Remove(tmpFile.Name())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to remove the tmp file before running backup")
	}

	args := []string{
		backup.databaseURL.String(),
		"-f", tmpFile.Name(),
		"-F", "c", // format: custom (zipped)
	}

	if backup.mode == config.DatabaseBackupModeLite {
		for _, table := range excludedDataFromTables {
			args = append(args, fmt.Sprintf("--exclude-table-data=%s", table))
		}
	}

	maskArgs := func(args []string) []string {
		masked := make([]string, len(args))
		copy(masked, args)
		masked[0] = backup.databaseURL.Redacted()
		return masked
	}

	maskedArgs := maskArgs(args)
	backup.logger.Debugf("Running pg_dump with: %v", maskedArgs)

	cmd := exec.Command(
		"pg_dump", args...,
	)

	_, err = cmd.Output()

	if err != nil {
		partialResult := &backupResult{
			size:            0,
			path:            "",
			maskedArguments: maskedArgs,
			pgDumpArguments: args,
		}
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			return partialResult, errors.Wrapf(err, "pg_dump failed with output: %s", string(ee.Stderr))
		}
		return partialResult, errors.Wrap(err, "pg_dump failed")
	}

	if version == "" {
		version = "unknown"
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
		size:            file.Size(),
		path:            finalFilePath,
		maskedArguments: maskedArgs,
		pgDumpArguments: args,
	}, nil
}
