package read

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
)

type ErrRead struct {
	Err    error
	Batch  bool
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

func newErrorFromCall(err error, call Call, block string, batch bool) ErrRead {
	return ErrRead{
		Err:   err,
		Batch: batch,
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

func (e ErrRead) Error() string {
	var builder strings.Builder

	builder.WriteString("[rpc error]")
	builder.WriteString(fmt.Sprintf(" batch: %T;", e.Batch))
	builder.WriteString(fmt.Sprintf(" err: %s;", e.Err.Error()))

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

func (e ErrRead) Unwrap() error {
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
