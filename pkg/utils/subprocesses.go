package utils

import "sync"

// Subprocesses is an abstraction over the following pattern:
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
