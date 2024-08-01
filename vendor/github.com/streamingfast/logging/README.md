# StreamingFast Logging library
[![reference](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://pkg.go.dev/github.com/streamingfast/logging)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

This is the logging library used as part of **[StreamingFast](https://github.com/streamingfast/streamingfast)**.


## Usage

In all library packages (by convention, use the import path):

```go
var zlog *zap.Logger

func init() {
	logging.Register("github.com/path/to/my/package", &zlog)
}
```

In `main` packages:

```go
var zlog *zap.Logger

func setupLogger() {
	logging.Register("main", &zlog)

	logging.Set(logging.MustCreateLogger())
	// Optionally set a different logger here and there,
	// using a regexp matching the registered names:
	//logging.Set(zap.NewNop(), "eosdb")
}
```

In tests (to avoid a race between the `init()` statements)

```go
func init() {
	if os.Getenv("DEBUG") != "" {
		logging.Override(logging.MustCreateLoggerWithLevel("test", zap.NewAtomicLevelAt(zap.DebugLevel)), ""))
	}
}
```

You can switch log levels dynamically, by poking the port 1065 like this:

On listening servers (port 1065, hint: logs!)


* `curl http://localhost:1065/ -XPUT -d '{"level": "debug"}'`


## Contributing

**Issues and PR in this repo related strictly to the streamingfast logging library.**

Report any protocol-specific issues in their
[respective repositories](https://github.com/streamingfast/streamingfast#protocols)

**Please first refer to the general
[StreamingFast contribution guide](https://github.com/streamingfast/streamingfast/blob/master/CONTRIBUTING.md)**,
if you wish to contribute to this code base.

## License

[Apache 2.0](LICENSE)
