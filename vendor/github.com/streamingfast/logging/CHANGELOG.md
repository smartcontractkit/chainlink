# Change log

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Next

### Changed

* The default text `encoder` use to encode log entries now emits the level when coloring is disabled.
* **Deprecated** `logging.IsTraceEnabled`, define your logger and `Tracer` directly with `var zlog, tracer = logging.PackageLogger(<shortName>, "...")` instead of separately, `tracer.Enabled()` can then be used to determine if tracing should be enabled (can be enable dynamically).
* **Deprecated** `logging.TestingOverride`, use `logging.InstantiateLoggers` directly.
* **Deprecated** `logging.Overidde`, use `logging.InstantiateLoggers` directly and use the `logging.WithDefaultSpec` to configure the various loggers.` instead.
* **Deprecated** `logging.Register`, use `var zlog, _ = logging.PackageLogger(<shortName>, "...")` instead.
* **Deprecated** `logging.RegisterOnUpdate`, use `logging.LoggerOnUpdate` instead (will probably be removed actually entirely since it's not needed anymore).
* **Deprecated** `logging.WithServiceName`, no replacement yet (will be `logging.LoggerServiceName` in a future release, if unspecified `shortName` will be used).

### Removed

* **BREAKING CHANGE** Removed `logging.Handler`, it has been moved to `dtracing` package to limit transitive depdendencies on this project, you should use `dtracing.NewAddTraceIDAwareLoggerMiddleware` instead.

## 2020-03-21

### Changed

* License changed to Apache 2.0