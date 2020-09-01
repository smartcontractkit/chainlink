// package logger exports multiple type loggers:
// - the *Logger type is a wrapper over uber/zap#SugaredLogger with added conditional utilities.
// - the Default instance is exported and used in all top-level functions (similar to godoc.org/log). This is sufficient for most use-cases.
// - ProductionLogger() builds a *Logger which stores logs on disk.
// - CreateTestLogger() prints log lines formatted for test runners. Use this in your tests!
// - CreateMemoryTestLogger() records logs in memory. It's useful for making assertions on the produced logs.
package logger
