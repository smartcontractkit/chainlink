package log

/*
	This file just exports internal symbols, for the purpose of unit-testing.
*/

type LogsOnBlock = logsOnBlock
type MockILogPool = mockILogPool

func NewLogPool() *logPool {
	return newLogPool()
}
