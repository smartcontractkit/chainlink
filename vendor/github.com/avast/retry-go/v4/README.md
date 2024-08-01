# retry

[![Release](https://img.shields.io/github/release/avast/retry-go.svg?style=flat-square)](https://github.com/avast/retry-go/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
![GitHub Actions](https://github.com/avast/retry-go/actions/workflows/workflow.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/avast/retry-go?style=flat-square)](https://goreportcard.com/report/github.com/avast/retry-go)
[![GoDoc](https://godoc.org/github.com/avast/retry-go?status.svg&style=flat-square)](http://godoc.org/github.com/avast/retry-go)
[![codecov.io](https://codecov.io/github/avast/retry-go/coverage.svg?branch=master)](https://codecov.io/github/avast/retry-go?branch=master)
[![Sourcegraph](https://sourcegraph.com/github.com/avast/retry-go/-/badge.svg)](https://sourcegraph.com/github.com/avast/retry-go?badge)

Simple library for retry mechanism

Slightly inspired by
[Try::Tiny::Retry](https://metacpan.org/pod/Try::Tiny::Retry)

# SYNOPSIS

HTTP GET with retry:

    url := "http://example.com"
    var body []byte

    err := retry.Do(
    	func() error {
    		resp, err := http.Get(url)
    		if err != nil {
    			return err
    		}
    		defer resp.Body.Close()
    		body, err = ioutil.ReadAll(resp.Body)
    		if err != nil {
    			return err
    		}
    		return nil
    	},
    )

    if err != nil {
    	// handle error
    }

    fmt.Println(string(body))

HTTP GET with retry with data:

    url := "http://example.com"

    body, err := retry.DoWithData(
    	func() ([]byte, error) {
    		resp, err := http.Get(url)
    		if err != nil {
    			return nil, err
    		}
    		defer resp.Body.Close()
    		body, err := ioutil.ReadAll(resp.Body)
    		if err != nil {
    			return nil, err
    		}

    		return body, nil
    	},
    )

    if err != nil {
    	// handle error
    }

    fmt.Println(string(body))

[More examples](https://github.com/avast/retry-go/tree/master/examples)

# SEE ALSO

* [giantswarm/retry-go](https://github.com/giantswarm/retry-go) - slightly
complicated interface.

* [sethgrid/pester](https://github.com/sethgrid/pester) - only http retry for
http calls with retries and backoff

* [cenkalti/backoff](https://github.com/cenkalti/backoff) - Go port of the
exponential backoff algorithm from Google's HTTP Client Library for Java. Really
complicated interface.

* [rafaeljesus/retry-go](https://github.com/rafaeljesus/retry-go) - looks good,
slightly similar as this package, don't have 'simple' `Retry` method

* [matryer/try](https://github.com/matryer/try) - very popular package,
nonintuitive interface (for me)

# BREAKING CHANGES

* 4.0.0

    - infinity retry is possible by set `Attempts(0)` by PR [#49](https://github.com/avast/retry-go/pull/49)

* 3.0.0

    - `DelayTypeFunc` accepts a new parameter `err` - this breaking change affects only your custom Delay Functions. This change allow [make delay functions based on error](examples/delay_based_on_error_test.go).

* 1.0.2 -> 2.0.0

    - argument of `retry.Delay` is final delay (no multiplication by `retry.Units` anymore)
    - function `retry.Units` are removed
    - [more about this breaking change](https://github.com/avast/retry-go/issues/7)

* 0.3.0 -> 1.0.0

    - `retry.Retry` function are changed to `retry.Do` function
    - `retry.RetryCustom` (OnRetry) and `retry.RetryCustomWithOpts` functions are now implement via functions produces Options (aka `retry.OnRetry`)

## Usage

#### func  BackOffDelay

```go
func BackOffDelay(n uint, _ error, config *Config) time.Duration
```
BackOffDelay is a DelayType which increases delay between consecutive retries

#### func  Do

```go
func Do(retryableFunc RetryableFunc, opts ...Option) error
```

#### func  DoWithData

```go
func DoWithData[T any](retryableFunc RetryableFuncWithData[T], opts ...Option) (T, error)
```

#### func  FixedDelay

```go
func FixedDelay(_ uint, _ error, config *Config) time.Duration
```
FixedDelay is a DelayType which keeps delay the same through all iterations

#### func  IsRecoverable

```go
func IsRecoverable(err error) bool
```
IsRecoverable checks if error is an instance of `unrecoverableError`

#### func  RandomDelay

```go
func RandomDelay(_ uint, _ error, config *Config) time.Duration
```
RandomDelay is a DelayType which picks a random delay up to config.maxJitter

#### func  Unrecoverable

```go
func Unrecoverable(err error) error
```
Unrecoverable wraps an error in `unrecoverableError` struct

#### type Config

```go
type Config struct {
}
```


#### type DelayTypeFunc

```go
type DelayTypeFunc func(n uint, err error, config *Config) time.Duration
```

DelayTypeFunc is called to return the next delay to wait after the retriable
function fails on `err` after `n` attempts.

#### func  CombineDelay

```go
func CombineDelay(delays ...DelayTypeFunc) DelayTypeFunc
```
CombineDelay is a DelayType the combines all of the specified delays into a new
DelayTypeFunc

#### type Error

```go
type Error []error
```

Error type represents list of errors in retry

#### func (Error) As

```go
func (e Error) As(target interface{}) bool
```

#### func (Error) Error

```go
func (e Error) Error() string
```
Error method return string representation of Error It is an implementation of
error interface

#### func (Error) Is

```go
func (e Error) Is(target error) bool
```

#### func (Error) Unwrap

```go
func (e Error) Unwrap() error
```
Unwrap the last error for compatibility with `errors.Unwrap()`. When you need to
unwrap all errors, you should use `WrappedErrors()` instead.

    err := Do(
    	func() error {
    		return errors.New("original error")
    	},
    	Attempts(1),
    )

    fmt.Println(errors.Unwrap(err)) # "original error" is printed

Added in version 4.2.0.

#### func (Error) WrappedErrors

```go
func (e Error) WrappedErrors() []error
```
WrappedErrors returns the list of errors that this Error is wrapping. It is an
implementation of the `errwrap.Wrapper` interface in package
[errwrap](https://github.com/hashicorp/errwrap) so that `retry.Error` can be
used with that library.

#### type OnRetryFunc

```go
type OnRetryFunc func(attempt uint, err error)
```

Function signature of OnRetry function

#### type Option

```go
type Option func(*Config)
```

Option represents an option for retry.

#### func  Attempts

```go
func Attempts(attempts uint) Option
```
Attempts set count of retry. Setting to 0 will retry until the retried function
succeeds. default is 10

#### func  AttemptsForError

```go
func AttemptsForError(attempts uint, err error) Option
```
AttemptsForError sets count of retry in case execution results in given `err`
Retries for the given `err` are also counted against total retries. The retry
will stop if any of given retries is exhausted.

added in 4.3.0

#### func  Context

```go
func Context(ctx context.Context) Option
```
Context allow to set context of retry default are Background context

example of immediately cancellation (maybe it isn't the best example, but it
describes behavior enough; I hope)

    ctx, cancel := context.WithCancel(context.Background())
    cancel()

    retry.Do(
    	func() error {
    		...
    	},
    	retry.Context(ctx),
    )

#### func  Delay

```go
func Delay(delay time.Duration) Option
```
Delay set delay between retry default is 100ms

#### func  DelayType

```go
func DelayType(delayType DelayTypeFunc) Option
```
DelayType set type of the delay between retries default is BackOff

#### func  LastErrorOnly

```go
func LastErrorOnly(lastErrorOnly bool) Option
```
return the direct last error that came from the retried function default is
false (return wrapped errors with everything)

#### func  MaxDelay

```go
func MaxDelay(maxDelay time.Duration) Option
```
MaxDelay set maximum delay between retry does not apply by default

#### func  MaxJitter

```go
func MaxJitter(maxJitter time.Duration) Option
```
MaxJitter sets the maximum random Jitter between retries for RandomDelay

#### func  OnRetry

```go
func OnRetry(onRetry OnRetryFunc) Option
```
OnRetry function callback are called each retry

log each retry example:

    retry.Do(
    	func() error {
    		return errors.New("some error")
    	},
    	retry.OnRetry(func(n uint, err error) {
    		log.Printf("#%d: %s\n", n, err)
    	}),
    )

#### func  RetryIf

```go
func RetryIf(retryIf RetryIfFunc) Option
```
RetryIf controls whether a retry should be attempted after an error (assuming
there are any retry attempts remaining)

skip retry if special error example:

    retry.Do(
    	func() error {
    		return errors.New("special error")
    	},
    	retry.RetryIf(func(err error) bool {
    		if err.Error() == "special error" {
    			return false
    		}
    		return true
    	})
    )

By default RetryIf stops execution if the error is wrapped using
`retry.Unrecoverable`, so above example may also be shortened to:

    retry.Do(
    	func() error {
    		return retry.Unrecoverable(errors.New("special error"))
    	}
    )

#### func  UntilSucceeded

```go
func UntilSucceeded() Option
```
UntilSucceeded will retry until the retried function succeeds. Equivalent to
setting Attempts(0).

#### func  WithTimer

```go
func WithTimer(t Timer) Option
```
WithTimer provides a way to swap out timer module implementations. This
primarily is useful for mocking/testing, where you may not want to explicitly
wait for a set duration for retries.

example of augmenting time.After with a print statement

    type struct MyTimer {}

    func (t *MyTimer) After(d time.Duration) <- chan time.Time {
        fmt.Print("Timer called!")
        return time.After(d)
    }

    retry.Do(
        func() error { ... },
    	   retry.WithTimer(&MyTimer{})
    )

#### func  WrapContextErrorWithLastError

```go
func WrapContextErrorWithLastError(wrapContextErrorWithLastError bool) Option
```
WrapContextErrorWithLastError allows the context error to be returned wrapped
with the last error that the retried function returned. This is only applicable
when Attempts is set to 0 to retry indefinitly and when using a context to
cancel / timeout

default is false

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    retry.Do(
    	func() error {
    		...
    	},
    	retry.Context(ctx),
    	retry.Attempts(0),
    	retry.WrapContextErrorWithLastError(true),
    )

#### type RetryIfFunc

```go
type RetryIfFunc func(error) bool
```

Function signature of retry if function

#### type RetryableFunc

```go
type RetryableFunc func() error
```

Function signature of retryable function

#### type RetryableFuncWithData

```go
type RetryableFuncWithData[T any] func() (T, error)
```

Function signature of retryable function with data

#### type Timer

```go
type Timer interface {
	After(time.Duration) <-chan time.Time
}
```

Timer represents the timer used to track time for a retry.

## Contributing

Contributions are very much welcome.

### Makefile

Makefile provides several handy rules, like README.md `generator` , `setup` for prepare build/dev environment, `test`, `cover`, etc...

Try `make help` for more information.

### Before pull request

> maybe you need `make setup` in order to setup environment

please try:
* run tests (`make test`)
* run linter (`make lint`)
* if your IDE don't automaticaly do `go fmt`, run `go fmt` (`make fmt`)

### README

README.md are generate from template [.godocdown.tmpl](.godocdown.tmpl) and code documentation via [godocdown](https://github.com/robertkrimen/godocdown).

Never edit README.md direct, because your change will be lost.
