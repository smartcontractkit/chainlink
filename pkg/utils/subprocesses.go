package utils

import "sync"

// Subprocesses is an abstraction over the following pattern of sync.WaitGroup:
//
//   var wg sync.Subprocesses
//   wg.Add(1)
//   go func() {
//     defer wg.Done()
//     // ...
//   }()
//
// Which becomes:
//
//  var subs utils.Subprocesses
//  subs.Go(func() {
//     // ...
//  })
//
// Note that it's important to not call Subprocesses.Wait() when there are
// no `Go()`ed functions in progress. This will panic.
// There are two cases when this can happen:
// 1. all the `Go()`ed functions started before the call to `Wait()` have already
// returned, maybe because a system-wide error or an already cancelled context.
// 2. Wait() gets called before any function is executed with `Go()`.
//
// Reusing a Subprocesses instance is discouraged.
// See mode details here https://pkg.go.dev/sync#WaitGroup.Add)
type Subprocesses struct {
	wg sync.WaitGroup
}

// Wait blocks until all function calls from the Go method have returned.
func (s *Subprocesses) Wait() {
	s.wg.Wait()
}

// Go calls the given function in a new goroutine.
func (s *Subprocesses) Go(f func()) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		f()
	}()
}
