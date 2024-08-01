package net

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrConnAccept struct {
	ID   uint32
	Name string
	Err  error
}

func (e ErrConnAccept) Error() string {
	return fmt.Sprintf("failed to accept %s server connection %d: %s", e.Name, e.ID, e.Err)
}

func (e ErrConnAccept) Unwrap() error {
	return e.Err
}

type ErrConnDial struct {
	ID   uint32
	Name string
	Err  error
}

func (e ErrConnDial) Error() string {
	return fmt.Sprintf("failed to dial %s client connection %d: %s", e.Name, e.ID, e.Err)
}

func (e ErrConnDial) Unwrap() error {
	return e.Err
}

// isErrTerminal returns true if the grpc [status] [codes.Code] indicates that the plugin connection has terminated and
// must be refreshed.
func isErrTerminal(err error) bool {
	switch status.Code(err) {
	case codes.Unavailable, codes.Canceled:
		return true
	case codes.OK, codes.Unknown, codes.InvalidArgument, codes.DeadlineExceeded, codes.NotFound, codes.AlreadyExists,
		codes.PermissionDenied, codes.ResourceExhausted, codes.FailedPrecondition, codes.Aborted, codes.OutOfRange,
		codes.Unimplemented, codes.Internal, codes.DataLoss, codes.Unauthenticated:
		return false
	}
	return false
}

func WrapRPCErr(err error) error {
	if err == nil {
		return nil
	}
	return &wrappedError{err: err, status: status.Convert(err)}
}

type wrappedError struct {
	err    error
	status *status.Status
}

func (w *wrappedError) Error() string {
	return w.err.Error()
}

func (w *wrappedError) Is(target error) bool {
	s := status.Convert(target)
	return w.status.Code() == s.Code() && strings.Contains(w.status.Message(), s.Message())
}

func (w *wrappedError) GRPCStatus() *status.Status {
	return w.status
}
