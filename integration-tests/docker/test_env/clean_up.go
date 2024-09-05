package test_env

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/logstream"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testreporters"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/testsummary"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/osutil"
	"github.com/smartcontractkit/chainlink-testing-framework/lib/utils/runid"
)

func AttachDefaultCleanUp(l zerolog.Logger, t *testing.T, clCluster *ClCluster, showHTMLCoverageReport bool, runId *string) {
	t.Cleanup(func() {
		l.Info().Msg("Cleaning up test environment")

		runIdErr := runid.RemoveLocalRunId(runId)
		if runIdErr != nil {
			l.Warn().Msgf("Failed to remove .run.id file due to: %s (not a big deal, you can still remove it manually)", runIdErr.Error())
		}

		err := SaveCodeCoverageData(l, t, clCluster, showHTMLCoverageReport)
		if err != nil {
			l.Error().Err(err).Msg("Error handling node coverage reports")
		}
	})
}

func AttachLogStreamCleanUp(l zerolog.Logger, t *testing.T, ls *logstream.LogStream, clCluster *ClCluster, chainlinkNodeLogScannerSettings *ChainlinkNodeLogScannerSettings, collectLogs bool) {
	t.Cleanup(func() {
		l.Info().Msg("Shutting down LogStream")
		logPath, err := osutil.GetAbsoluteFolderPath("logs")
		if err == nil {
			l.Info().Str("Absolute path", logPath).Msg("LogStream logs folder location")
		}

		// flush logs when test failed or when we are explicitly told to collect logs
		flushLogStream := t.Failed() || collectLogs

		// run even if test has failed, as we might be able to catch additional problems without running the test again
		if chainlinkNodeLogScannerSettings != nil {
			logProcessor := logstream.NewLogProcessor[int](ls)

			processFn := func(log logstream.LogContent, count *int) error {
				countSoFar := count
				newCount, err := testreporters.ScanLogLine(l, string(log.Content), chainlinkNodeLogScannerSettings.FailingLogLevel, uint(*countSoFar), chainlinkNodeLogScannerSettings.Threshold, chainlinkNodeLogScannerSettings.AllowedMessages)
				if err != nil {
					return err
				}
				*count = int(newCount)
				return nil
			}

			// we cannot do parallel processing here, because ProcessContainerLogs() locks a mutex that controls whether
			// new logs can be added to the log stream, so parallel processing would get stuck on waiting for it to be unlocked
			if clCluster != nil {
			LogScanningLoop:
				for i := 0; i < len(clCluster.Nodes); i++ {
					// if something went wrong during environment setup we might not have all nodes, and we don't want an NPE
					if len(clCluster.Nodes)-1 < i || clCluster.Nodes[i] == nil {
						continue
					}
					// ignore count return, because we are only interested in the error
					_, err := logProcessor.ProcessContainerLogs(clCluster.Nodes[i].ContainerName, processFn)
					if err != nil && !strings.Contains(err.Error(), testreporters.MultipleLogsAtLogLevelErr) && !strings.Contains(err.Error(), testreporters.OneLogAtLogLevelErr) {
						l.Error().Err(err).Msg("Error processing CL node logs")
						continue
					} else if err != nil && (strings.Contains(err.Error(), testreporters.MultipleLogsAtLogLevelErr) || strings.Contains(err.Error(), testreporters.OneLogAtLogLevelErr)) {
						flushLogStream = true
						t.Errorf("Found a concerning log in Chainklink Node logs: %v", err)
						break LogScanningLoop
					}
				}
				l.Info().Msg("Finished scanning Chainlink Node logs for concerning errors")
			}
		}

		if flushLogStream {
			l.Info().Msg("Flushing LogStream logs")
			// we can't do much if this fails, so we just log the error in LogStream
			if err := ls.FlushAndShutdown(); err != nil {
				l.Error().Err(err).Msg("Error flushing and shutting down LogStream")
			}
			ls.PrintLogTargetsLocations()
			ls.SaveLogLocationInTestSummary()
		}
		l.Info().Msg("Finished shutting down LogStream")
	})
}

func AttachDbDumpingCleanup(l zerolog.Logger, t *testing.T, clCluster *ClCluster, collectLogs bool) {
	t.Cleanup(func() {
		if t.Failed() || collectLogs {
			l.Info().Msg("Dump state of all Postgres DBs used by Chainlink Nodes")

			dbDumpFolder := "db_dumps"
			dbDumpPath := fmt.Sprintf("%s/%s-%s", dbDumpFolder, t.Name(), time.Now().Format("2006-01-02T15-04-05"))
			if err := os.MkdirAll(dbDumpPath, os.ModePerm); err != nil {
				l.Error().Err(err).Msg("Error creating folder for Postgres DB dump")
				return
			}

			absDbDumpPath, err := osutil.GetAbsoluteFolderPath(dbDumpFolder)
			if err == nil {
				l.Info().Str("Absolute path", absDbDumpPath).Msg("PostgresDB dump folder location")
			}

			if clCluster != nil {
				for i := 0; i < len(clCluster.Nodes); i++ {
					// if something went wrong during environment setup we might not have all nodes, and we don't want an NPE
					if len(clCluster.Nodes)-1 < i || clCluster.Nodes[i] == nil || clCluster.Nodes[i].PostgresDb == nil {
						continue
					}

					filePath := filepath.Join(dbDumpPath, fmt.Sprintf("postgres_db_dump_%s.sql", clCluster.Nodes[i].ContainerName))
					localDbDumpFile, err := os.Create(filePath)
					if err != nil {
						l.Error().Err(err).Msg("Error creating localDbDumpFile for Postgres DB dump")
						_ = localDbDumpFile.Close()
						continue
					}

					if err := clCluster.Nodes[i].PostgresDb.ExecPgDumpFromContainer(localDbDumpFile); err != nil {
						l.Error().Err(err).Msg("Error dumping Postgres DB")
					}
					_ = localDbDumpFile.Close()
				}
				l.Info().Msg("Finished dumping state of all Postgres DBs used by Chainlink Nodes")
			}
		}
	})
}

func AttachSethCleanup(t *testing.T, sethConfig *seth.Config) {
	t.Cleanup(func() {
		if sethConfig != nil && ((t.Failed() && slices.Contains(sethConfig.TraceOutputs, seth.TraceOutput_DOT) && sethConfig.TracingLevel != seth.TracingLevel_None) || (!t.Failed() && slices.Contains(sethConfig.TraceOutputs, seth.TraceOutput_DOT) && sethConfig.TracingLevel == seth.TracingLevel_All)) {
			_ = testsummary.AddEntry(t.Name(), "dot_graphs", "true")
		}
	})
}
