package proposalutils

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
)

const INDENT = "    "

var (
	_                  Analyzer = BytesAndAddressAnalyzer
	_                  Analyzer = ChainSelectorAnalyzer
	chainSelectorRegex          = regexp.MustCompile(`[cC]hain([sS]el)?.*$`)
)

// Argument is a unit of decoded proposal. Calling Describe on it returns human-readable representation of its content.
// Some implementations are recursive (arrays, structs) and require attention to formatting.
type Argument interface {
	Describe(ctx *ArgumentContext) string
}

// Analyzer is an extension point of proposal decoding.
// You can implement your own Analyzer which returns your own Argument instance.
type Analyzer func(argName string, argAbi *abi.Type, argVal interface{}, analyzers []Analyzer) Argument

// ArgumentContext is a storage for context that may need to Argument during its description.
// Refer to BytesAndAddressAnalyzer and ChainSelectorAnalyzer for usage examples
type ArgumentContext struct {
	ChainSelector uint64
	Ctx           map[string]interface{}
}

func NewArgumentContext(chainSelector uint64, keysAndValues ...interface{}) *ArgumentContext {
	ctx := make(map[string]interface{})
	var key string
	for i, kv := range keysAndValues {
		if i%2 == 0 {
			key = kv.(string)
		} else {
			ctx[key] = kv
		}
	}
	return &ArgumentContext{
		ChainSelector: chainSelector,
		Ctx:           ctx,
	}
}

func ContextGet[T any](ctx *ArgumentContext, key string) (T, error) {
	ctxElemRaw, ok := ctx.Ctx[key]
	if !ok {
		return *new(T), fmt.Errorf("context element %s not found", key)
	}
	ctxElem, ok := ctxElemRaw.(T)
	if !ok {
		return *new(T), fmt.Errorf("context element %s type mismatch (expected: %T, was: %T)", key, ctxElem, *new(T))
	}
	return ctxElem, nil
}

type NamedArgument struct {
	Name  string
	Value Argument
}

func (n NamedArgument) Describe(context *ArgumentContext) string {
	return fmt.Sprintf("%s: %s", n.Name, n.Value.Describe(context))
}

type ArrayArgument struct {
	Elements []Argument
}

func (a ArrayArgument) Describe(context *ArgumentContext) string {
	indented := false
	elementsDescribed := make([]string, len(a.Elements))
	for _, arg := range a.Elements {
		argDescribed := arg.Describe(context)
		indented = indented || strings.Contains(argDescribed, INDENT)
		elementsDescribed = append(elementsDescribed, argDescribed)
	}
	description := strings.Builder{}
	if indented {
		// Write each element in new line + indentation
		description.WriteString("[\n")
		for i, elem := range elementsDescribed {
			description.WriteString(IndentString(elem))
			if i < len(a.Elements)-1 {
				description.WriteString(",\n")
			}
		}
		description.WriteString("\n]")
	} else {
		// Write elements in one line
		description.WriteString("[")
		for i, elem := range elementsDescribed {
			description.WriteString(elem)
			if i < len(a.Elements)-1 {
				description.WriteString(",")
			}
		}
		description.WriteString("]")
	}
	return description.String()
}

type StructArgument struct {
	Fields []NamedArgument
}

func (s StructArgument) Describe(context *ArgumentContext) string {
	description := strings.Builder{}
	if len(s.Fields) >= 2 {
		// Pretty format struct with indentation
		description.WriteString("{\n")
		for _, arg := range s.Fields {
			description.WriteString(IndentString(arg.Describe(context)))
			description.WriteString("\n")
		}
		description.WriteString("}")
	} else {
		// Struct in one line
		description.WriteString("{ ")
		for i, arg := range s.Fields {
			description.WriteString(arg.Describe(context))
			if i < len(s.Fields)-1 {
				description.WriteString(", ")
			}
		}
		description.WriteString(" }")
	}
	return description.String()
}

type SimpleArgument struct {
	Value string
}

func (s SimpleArgument) Describe(_ *ArgumentContext) string {
	return s.Value
}

type ChainSelectorArgument struct {
	Value uint64
}

func (c ChainSelectorArgument) Describe(_ *ArgumentContext) string {
	chainID, err := chain_selectors.GetChainIDFromSelector(c.Value)
	if err != nil {
		return fmt.Sprintf("%d (unknown)", c.Value)
	}
	family, err := chain_selectors.GetSelectorFamily(c.Value)
	if err != nil {
		return fmt.Sprintf("%d (unknown)", c.Value)
	}
	chainInfo, err := chain_selectors.GetChainDetailsByChainIDAndFamily(chainID, family)
	if err != nil {
		return fmt.Sprintf("%d (unknown)", c.Value)
	}
	return fmt.Sprintf("%d (%s)", c.Value, chainInfo.ChainName)
}

type BytesArgument struct {
	Value []byte
}

func (a BytesArgument) Describe(_ *ArgumentContext) string {
	return hexutil.Encode(a.Value)
}

type AddressArgument struct {
	Value common.Address
}

func (a AddressArgument) Describe(ctx *ArgumentContext) string {
	description := a.Value.Hex() + " (address, unknown)"
	addressBook, err := ContextGet[deployment.AddressBook](ctx, "AddressBook")
	if err != nil {
		return description
	}
	addresses, err := addressBook.Addresses()
	if err != nil {
		return description
	}
	typeAndVersion, ok := addresses[ctx.ChainSelector][a.Value.Hex()]
	if ok {
		return fmt.Sprintf("%s (address, %d, %s)", a.Value.Hex(), ctx.ChainSelector, typeAndVersion.String())
	}
	// it can be under another chain
	for chainSel, addresses := range addresses {
		typeAndVersion, ok := addresses[a.Value.Hex()]
		if ok {
			return fmt.Sprintf("%s (address, %d, %s)", a.Value.Hex(), chainSel, typeAndVersion.String())
		}
	}
	return description
}

type DecodedCall struct {
	Address string
	Method  string
	Inputs  []NamedArgument
	Outputs []NamedArgument
}

func (d *DecodedCall) Describe(context *ArgumentContext) string {
	description := strings.Builder{}
	description.WriteString(fmt.Sprintf("Proposal chain: %s\n", ChainSelectorArgument{Value: context.ChainSelector}.Describe(context)))
	description.WriteString(fmt.Sprintf("Address: %s (%s)\n", AddressArgument{Value: common.HexToAddress(d.Address)}.Describe(context)))
	description.WriteString(fmt.Sprintf("Method: %s\n", d.Method))
	description.WriteString(fmt.Sprintf("Inputs:\n%s\n", IndentString(d.describeArguments(d.Inputs, context))))
	description.WriteString(fmt.Sprintf("Outputs:\n%s\n", IndentString(d.describeArguments(d.Outputs, context))))
	return description.String()
}

func (d *DecodedCall) describeArguments(arguments []NamedArgument, context *ArgumentContext) string {
	description := strings.Builder{}
	for _, argument := range arguments {
		description.WriteString(argument.Describe(context))
		description.WriteRune('\n')
	}
	return description.String()
}

type TxCallDecoder struct {
	Analyzers []Analyzer
}

func NewTxCallDecoder(extraAnalyzers []Analyzer) *TxCallDecoder {
	analyzers := make([]Analyzer, 0, len(extraAnalyzers)+2)
	analyzers = append(analyzers, extraAnalyzers...)
	analyzers = append(analyzers, BytesAndAddressAnalyzer)
	analyzers = append(analyzers, ChainSelectorAnalyzer)
	return &TxCallDecoder{Analyzers: analyzers}
}

func (p *TxCallDecoder) Analyze(abi *abi.ABI, data []byte) (DecodedCall, error) {
	methodID, methodData := data[:4], data[4:]
	method, err := abi.MethodById(methodID)
	if err != nil {
		return DecodedCall{}, err
	}
	outs := make(map[string]interface{})
	err = method.Outputs.UnpackIntoMap(outs, methodData)
	if err != nil {
		return DecodedCall{}, err
	}
	args := make(map[string]interface{})
	err = method.Inputs.UnpackIntoMap(args, methodData)
	if err != nil {
		return DecodedCall{}, err
	}
	return p.analyzeMethodCall(method, args, outs)
}

func (p *TxCallDecoder) analyzeMethodCall(method *abi.Method, args map[string]interface{}, outs map[string]interface{}) (DecodedCall, error) {
	inputs := make([]NamedArgument, len(method.Inputs))
	for i, input := range method.Inputs {
		arg, ok := args[input.Name]
		if !ok {
			return DecodedCall{}, fmt.Errorf("missing argument '%s'", input.Name)
		}
		inputs[i] = NamedArgument{
			Name:  input.Name,
			Value: p.analyzeArg(input.Name, &input.Type, arg),
		}
	}
	outputs := make([]NamedArgument, len(method.Outputs))
	for i, output := range method.Outputs {
		out, ok := outs[output.Name]
		if !ok {
			return DecodedCall{}, fmt.Errorf("missing output '%s'", output.Name)
		}
		outputs[i] = NamedArgument{
			Name:  output.Name,
			Value: p.analyzeArg(output.Name, &output.Type, out),
		}
	}
	return DecodedCall{
		Method:  method.String(),
		Inputs:  inputs,
		Outputs: outputs,
	}, nil
}

func (p *TxCallDecoder) analyzeArg(argName string, argAbi *abi.Type, argVal interface{}) Argument {
	if len(p.Analyzers) > 0 {
		for _, analyzer := range p.Analyzers {
			arg := analyzer(argName, argAbi, argVal, p.Analyzers)
			if arg != nil {
				return arg
			}
		}
	}
	// Struct analyzer
	if argAbi.T == abi.TupleTy {
		return p.analyzeStruct(argAbi, argVal)
	}
	// Array analyzer
	if argAbi.T == abi.SliceTy || argAbi.T == abi.ArrayTy {
		return p.analyzeArray(argName, argAbi, argVal)
	}
	// Fallback
	return SimpleArgument{Value: fmt.Sprintf("%v", argVal)}
}

func (p *TxCallDecoder) analyzeStruct(argAbi *abi.Type, argVal interface{}) StructArgument {
	argTyp := argAbi.GetType()
	fields := make([]NamedArgument, argTyp.NumField())
	for i := 0; i < argTyp.NumField(); i++ {
		if !argTyp.Field(i).IsExported() {
			continue
		}
		argFieldName := argTyp.Field(i).Name
		argFieldAbi := argAbi.TupleElems[i]
		argFieldTyp := reflect.ValueOf(argVal).FieldByName(argFieldName)
		argument := p.analyzeArg(argFieldName, argFieldAbi, argFieldTyp.Interface())
		fields[i] = NamedArgument{
			Name:  argFieldName,
			Value: argument,
		}
	}
	return StructArgument{
		Fields: fields,
	}
}

func (p *TxCallDecoder) analyzeArray(argName string, argAbi *abi.Type, argVal interface{}) ArrayArgument {
	argTyp := reflect.ValueOf(argVal)
	elements := make([]Argument, argTyp.Len())
	for i := 0; i < argTyp.Len(); i++ {
		argElemTyp := argTyp.Index(i)
		argument := p.analyzeArg(argName, argAbi.Elem, argElemTyp.Interface())
		elements[i] = argument
	}
	return ArrayArgument{
		Elements: elements,
	}
}

func BytesAndAddressAnalyzer(_ string, argAbi *abi.Type, argVal interface{}, _ []Analyzer) Argument {
	if argAbi.T == abi.FixedBytesTy || argAbi.T == abi.BytesTy || argAbi.T == abi.AddressTy {
		argArrTyp := reflect.ValueOf(argVal)
		argArr := make([]byte, argArrTyp.Len())
		for i := 0; i < argArrTyp.Len(); i++ {
			argArr[i] = byte(argArrTyp.Index(i).Uint())
		}
		if argAbi.T == abi.AddressTy {
			return AddressArgument{Value: common.BytesToAddress(argArr)}
		}
		return BytesArgument{Value: argArr}
	}
	return nil
}

func ChainSelectorAnalyzer(argName string, argAbi *abi.Type, argVal interface{}, _ []Analyzer) Argument {
	if argAbi.GetType().Kind() == reflect.Uint64 && chainSelectorRegex.MatchString(argName) {
		return ChainSelectorArgument{Value: argVal.(uint64)}
	}
	return nil
}

func IndentString(s string) string {
	return IndentStringWith(s, INDENT)
}

func IndentStringWith(s string, indent string) string {
	result := &strings.Builder{}
	components := strings.Split(s, "\n")
	for i, component := range components {
		result.WriteString(indent)
		result.WriteString(component)
		if i < len(components)-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}
