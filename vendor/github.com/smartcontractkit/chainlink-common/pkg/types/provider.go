package types

import (
	"slices"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
)

// The bootstrap jobs only watch config.
type ConfigProvider interface {
	Service
	OffchainConfigDigester() ocrtypes.OffchainConfigDigester
	ContractConfigTracker() ocrtypes.ContractConfigTracker
}

// Plugin is an alias for PluginProvider, for compatibility.
// Deprecated
type Plugin = PluginProvider

// PluginProvider provides common components for any OCR2 plugin.
// It watches config and is able to transmit.
type PluginProvider interface {
	ConfigProvider
	ContractTransmitter() ocrtypes.ContractTransmitter
	ChainReader() ChainReader
	Codec() Codec
}

type OCR3ContractTransmitter interface {
	OCR3ContractTransmitter() ocr3types.ContractTransmitter[[]byte]
}

// General error types for providers to return--can be used to wrap more specific errors.
// These should work with or without LOOP enabled, to help the client decide how to handle
// an error. The structure of any wrapped errors would normally be automatically flattened
// to a single string, making it difficult for the client to respond to different categories
// of errors in different ways. This lessons the need for doing our own custom parsing of
// error strings.

type InvalidArgumentError string

func (e InvalidArgumentError) Error() string {
	return string(e)
}

func (e InvalidArgumentError) GRPCStatus() *status.Status {
	return status.New(codes.InvalidArgument, e.Error())
}

func (e InvalidArgumentError) Is(target error) bool {
	if e == target {
		return true
	}

	return grpcErrorHasTypeAndMessage(target, string(e), codes.InvalidArgument)
}

type UnimplementedError string

func (e UnimplementedError) Error() string {
	return string(e)
}

func (e UnimplementedError) GRPCStatus() *status.Status {
	return status.New(codes.Unimplemented, e.Error())
}

func (e UnimplementedError) Is(target error) bool {
	if e == target {
		return true
	}

	return grpcErrorHasTypeAndMessage(target, string(e), codes.Unimplemented)
}

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

func (e InternalError) GRPCStatus() *status.Status {
	return status.New(codes.Internal, e.Error())
}

func (e InternalError) Is(target error) bool {
	if e == target {
		return true
	}

	return grpcErrorHasTypeAndMessage(target, string(e), codes.Internal)
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return string(e)
}

func (e NotFoundError) GRPCStatus() *status.Status {
	return status.New(codes.NotFound, e.Error())
}

func (e NotFoundError) Is(target error) bool {
	if e == target {
		return true
	}

	return grpcErrorHasTypeAndMessage(target, string(e), codes.NotFound)
}

func grpcErrorHasTypeAndMessage(target error, msg string, code codes.Code) bool {
	s, ok := status.FromError(target)
	if !ok || s.Code() != code {
		return false
	}

	errs := strings.Split(s.Message(), ":")
	return slices.ContainsFunc(errs, func(err string) bool {
		return strings.Trim(err, " ") == msg
	})
}
