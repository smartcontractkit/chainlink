//go:build integration

package web_test

func ptr[T any](t T) *T { return &t }
