package multiexp

import "errors"

var ErrTooManyGoRoutines = errors.New("cannot configure more than 1024 go routines")
