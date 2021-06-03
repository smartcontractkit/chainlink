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
	"github.com/smartcontractkit/chainlink/core/service"
	"github.com/smartcontractkit/chainlink/core/static"
	"github.com/smartcontractkit/chainlink/core/store/orm"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	filePattern        = "cl_backup_%s.dump"
	minBackupFrequency = time.Minute

	excludedDataFromTables = []string{
		"job_runs",
		"task_runs",
		"eth_task_run_txes",
		"run_requests",
		"run_results",
		"sync_events",
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
		service.Service
		RunBackupGracefully(version string)
	}

	databaseBackup struct {
		logger          *logger.Logger
		databaseURL     url.URL
		mode            orm.DatabaseBackupMode
		frequency       time.Duration
		outputParentDir string
		done            chan bool
		utils.StartStopOnce
	}

	Config interface {
		DatabaseBackupMode() orm.DatabaseBackupMode
		DatabaseBackupFrequency() time.Duration
		DatabaseBackupURL() *url.URL
		DatabaseBackupDir() string
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

	outputParentDir := filepath.Join(config.RootDir(), "backup")
	if config.DatabaseBackupDir() != "" {
		dir, err := filepath.Abs(config.DatabaseBackupDir())
		if err != nil {
			logger.Errorf("Invalid path for DATABASE_BACKUP_DIR (%s) - please set it to a valid directory path", config.DatabaseBackupDir())
		}
		outputParentDir = dir
	}

	return &databaseBackup{
		logger,
		dbUrl,
		config.DatabaseBackupMode(),
		config.DatabaseBackupFrequency(),
		outputParentDir,
		make(chan bool),
		utils.StartStopOnce{},
	}
}

func (backup *databaseBackup) Start() error {
	return backup.StartOnce("DatabaseBackup", func() (err error) {
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

func (backup *databaseBackup) RunBackupGracefully(version string) {
	backup.logger.Debugw("DatabaseBackup: Starting database backup...", "mode", backup.mode, "url", backup.databaseURL.String(), "directory", backup.outputParentDir)
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

	err := os.MkdirAll(backup.outputParentDir, os.ModePerm)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("DatabaseBackup: Failed to create directories on the path: %s", backup.outputParentDir))
	}
	tmpFile, err := ioutil.TempFile(backup.outputParentDir, "cl_backup_tmp_")
	if err != nil {
		return nil, errors.Wrap(err, "DatabaseBackup: Failed to create a tmp file")
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

	if backup.mode == orm.DatabaseBackupModeLite {
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
	backup.logger.Debugf("DatabaseBackup: Running pg_dump with: %v", maskedArgs)

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
		if ee, ok := err.(*exec.ExitError); ok {
			return partialResult, errors.Wrap(err, fmt.Sprintf("pg_dump failed with output: %s", string(ee.Stderr)))
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
