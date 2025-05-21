package read

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
)

type readType string

const (
	batchReadType  readType = "BatchGetLatestValue"
	singleReadType readType = "GetLatestValue"
	eventReadType  readType = "QueryKey"
)

type Error struct {
	Err    error
	Type   readType
	Detail *readDetail
	Result *string
}

type readDetail struct {
	Address        string
	Contract       string
	Method         string
	Params, RetVal any
	Block          string
}

func newErrorFromCall(err error, call Call, block string, tp readType) Error {
	return Error{
		Err:  err,
		Type: tp,
		Detail: &readDetail{
			Address:  call.ContractAddress.Hex(),
			Contract: call.ContractName,
			Method:   call.ReadName,
			Params:   call.Params,
			RetVal:   call.ReturnVal,
			Block:    block,
		},
	}
}

func (e Error) Error() string {
	var builder strings.Builder

	builder.WriteString("[read error]")
	builder.WriteString(fmt.Sprintf(" err: %s;", e.Err.Error()))
	builder.WriteString(fmt.Sprintf(" type: %s;", e.Type))

	if e.Detail != nil {
		builder.WriteString(fmt.Sprintf(" block: %s;", e.Detail.Block))
		builder.WriteString(fmt.Sprintf(" address: %s;", e.Detail.Address))
		builder.WriteString(fmt.Sprintf(" contract-name: %s;", e.Detail.Contract))
		builder.WriteString(fmt.Sprintf(" read-name: %s;", e.Detail.Method))
		builder.WriteString(fmt.Sprintf(" params: %+v;", e.Detail.Params))
		builder.WriteString(fmt.Sprintf(" expected return type: %s;", reflect.TypeOf(e.Detail.RetVal)))

		if e.Result != nil {
			builder.WriteString(fmt.Sprintf("encoded result: %s;", *e.Result))
		}
	}

	return builder.String()
}

func (e Error) Unwrap() error {
	return e.Err
}

type MultiCallError struct {
	Err    error
	Type   readType
	Detail *callsReadDetail
	Result *string
}

type callsReadDetail struct {
	Calls []Call
	Block string
}

func newErrorFromCalls(err error, calls []Call, block string, tp readType) MultiCallError {
	return MultiCallError{
		Err:  err,
		Type: tp,
		Detail: &callsReadDetail{
			Calls: calls,
			Block: block,
		},
	}
}

func (e MultiCallError) Error() string {
	var builder strings.Builder

	builder.WriteString("[read error]")
	builder.WriteString(fmt.Sprintf(" err: %s;", e.Err.Error()))
	builder.WriteString(fmt.Sprintf(" type: %s;", e.Type))

	if e.Detail != nil {
		builder.WriteString(fmt.Sprintf(" block: %s;", e.Detail.Block))
		for _, call := range e.Detail.Calls {
			builder.WriteString(fmt.Sprintf(" address: %s;", call.ContractAddress.Hex()))
			builder.WriteString(fmt.Sprintf(" contract-name: %s;", call.ContractName))
			builder.WriteString(fmt.Sprintf(" read-name: %s;", call.ReadName))
			builder.WriteString(fmt.Sprintf(" params: %+v;", call.Params))
			builder.WriteString(fmt.Sprintf(" expected return type: %s;", reflect.TypeOf(call.ReturnVal)))
		}

		if e.Result != nil {
			builder.WriteString(fmt.Sprintf("encoded result: %s;", *e.Result))
		}
	}

	return builder.String()
}

func (e MultiCallError) Unwrap() error {
	return e.Err
}

type ConfigError struct {
	Msg string
}

func newMissingReadIdentifierErr(readIdentifier string) ConfigError {
	return ConfigError{
		Msg: fmt.Sprintf("[no configured reader] read-identifier: '%s'", readIdentifier),
	}
}

func newMissingContractErr(readIdentifier, contract string) ConfigError {
	return ConfigError{
		Msg: fmt.Sprintf("[no configured reader] read-identifier: %s; contract: %s;", readIdentifier, contract),
	}
}

func newMissingReadNameErr(readIdentifier, contract, readName string) ConfigError {
	return ConfigError{
		Msg: fmt.Sprintf("[no configured reader] read-identifier: %s; contract: %s; read-name: %s;", readIdentifier, contract, readName),
	}
}

func newUnboundAddressErr(address, contract, readName string) ConfigError {
	return ConfigError{
		Msg: fmt.Sprintf("[address not bound] address: %s; contract: %s; read-name: %s;", address, contract, readName),
	}
}

func (e ConfigError) Error() string {
	return e.Msg
}

type FilterError struct {
	Err    error
	Action string
	Filter logpoller.Filter
}

func (e FilterError) Error() string {
	return fmt.Sprintf("[logpoller filter error] action: %s; err: %s; filter: %+v;", e.Action, e.Err.Error(), e.Filter)
}

func (e FilterError) Unwrap() error {
	return e.Err
}

type NoContractExistsError struct {
	Err     error
	Address common.Address
}

func (e NoContractExistsError) Error() string {
	return fmt.Sprintf("%s: contract does not exist at address: %s", e.Err.Error(), e.Address)
}

func (e NoContractExistsError) Unwrap() error {
	return e.Err
}
