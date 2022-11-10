package config

import "fmt"

// lightweight error types copied from core

type ErrInvalid struct {
	Name  string
	Value any
	Msg   string
}

func (e ErrInvalid) Error() string {
	return fmt.Sprintf("%s: invalid value %v: %s", e.Name, e.Value, e.Msg)
}

type ErrMissing struct {
	Name string
	Msg  string
}

func (e ErrMissing) Error() string {
	return fmt.Sprintf("%s: missing: %s", e.Name, e.Msg)
}

type ErrEmpty struct {
	Name string
	Msg  string
}

func (e ErrEmpty) Error() string {
	return fmt.Sprintf("%s: empty: %s", e.Name, e.Msg)
}
