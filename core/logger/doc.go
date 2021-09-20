// Package logger exports multiple type loggers:
// - the Logger type is a wrapper over uber/zap#SugaredLogger with added conditional utilities.
// - the Default instance is exported and used in all top-level functions (similar to godoc.org/log).
//   You can use this directly but the recommended way is to create a new Logger for your module using Named(), eg:  logger.Default.Named("<my-package-name>")
// - ProductionLogger() builds a Logger which stores logs on disk.
// - CreateTestLogger() prints log lines formatted for test runners. Use this in your tests!
// - CreateMemoryTestLogger() records logs in memory. It's useful for making assertions on the produced logs.
package logger
