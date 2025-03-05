package evmtesting

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/go-viper/mapstructure/v2"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/interfacetests"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/bindings"
	evmcodec "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/codec"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

const contractName = "ChainReaderTester"

// Wraps EVMChainComponentsInterfaceTester to rely on the EVM bindings generated for CR/CW instead of going directly to CR/CW. This way we can reuse all existing tests. Transformation between expected
// contract names and read keys will be done here as well as invocation delegation to generated code.
func WrapContractReaderTesterWithBindings(t *testing.T, wrapped *EVMChainComponentsInterfaceTester[*testing.T]) interfacetests.ChainComponentsInterfaceTester[*testing.T] {
	// Tests not yet supported by EVM bindings.
	wrapped.DisableTests([]string{
		interfacetests.ContractReaderGetLatestValueAsValuesDotValue, interfacetests.ContractReaderGetLatestValueNoArgumentsAndPrimitiveReturnAsValuesDotValue, interfacetests.ContractReaderGetLatestValueNoArgumentsAndSliceReturnAsValueDotValue,
		interfacetests.ContractReaderGetLatestValueGetsLatestForEvent, interfacetests.ContractReaderGetLatestValueBasedOnConfidenceLevelForEvent,
		interfacetests.ContractReaderGetLatestValueReturnsNotFoundWhenNotTriggeredForEvent, interfacetests.ContractReaderGetLatestValueWithFilteringForEvent, interfacetests.ContractReaderBatchGetLatestValue, interfacetests.ContractReaderBatchGetLatestValueMultipleContractNamesSameFunction,
		interfacetests.ContractReaderBatchGetLatestValueDifferentParamsResultsRetainOrder, interfacetests.ContractReaderBatchGetLatestValueDifferentParamsResultsRetainOrderMultipleContracts, interfacetests.ContractReaderBatchGetLatestValueNoArgumentsPrimitiveReturn,
		interfacetests.ContractReaderBatchGetLatestValueSetsErrorsProperly, interfacetests.ContractReaderBatchGetLatestValueNoArgumentsWithSliceReturn, interfacetests.ContractReaderBatchGetLatestValueWithModifiersOwnMapstructureOverride,
		interfacetests.ContractReaderQueryKeyNotFound, interfacetests.ContractReaderQueryKeyReturnsData, interfacetests.ContractReaderQueryKeyReturnsDataAsValuesDotValue, interfacetests.ContractReaderQueryKeyReturnsDataAsValuesDotValue,
		interfacetests.ContractReaderQueryKeyCanFilterWithValueComparator, interfacetests.ContractReaderQueryKeyCanLimitResultsWithCursor,
		interfacetests.ContractReaderQueryKeysNotFound, interfacetests.ContractReaderQueryKeysReturnsData, interfacetests.ContractReaderQueryKeysReturnsDataTwoEventTypes, interfacetests.ContractReaderQueryKeysReturnsDataAsValuesDotValue,
		interfacetests.ContractReaderQueryKeysCanFilterWithValueComparator, interfacetests.ContractReaderQueryKeysCanLimitResultsWithCursor,
		ContractReaderQueryKeyFilterOnDataWordsWithValueComparator, ContractReaderQueryKeyOnDataWordsWithValueComparatorOnNestedField,
		ContractReaderQueryKeyFilterOnDataWordsWithValueComparatorOnDynamicField, ContractReaderQueryKeyFilteringOnDataWordsUsingValueComparatorsOnFieldsWithManualIndex,
		// TODO BCFR-1073 - Fix flaky tests
		interfacetests.ContractReaderGetLatestValueBasedOnConfidenceLevel,
	})
	wrapped.SetChainReaderConfigSupplier(func(t *testing.T) types.ChainReaderConfig {
		return getChainReaderConfig(wrapped)
	})
	wrapped.SetChainWriterConfigSupplier(func(t *testing.T) types.ChainWriterConfig {
		return getChainWriterConfig(t, wrapped)
	})
	return newBindingClientTester(wrapped)
}

func newBindingClientTester(wrapped *EVMChainComponentsInterfaceTester[*testing.T]) bindingClientTester {
	bindingsMapping := newBindingsMapping()
	return bindingClientTester{
		ChainComponentsInterfaceTester: wrapped,
		bindingsMapping:                &bindingsMapping,
	}
}

func newBindingsMapping() bindingsMapping {
	contractReaderProxy := bindingContractReaderProxy{}
	chainWriterProxy := bindingChainWriterProxy{}
	methodNameMappingByContract := make(map[string]map[string]string)
	methodNameMappingByContract[interfacetests.AnyContractName] = map[string]string{
		interfacetests.MethodTakingLatestParamsReturningTestStruct: "GetElementAtIndex",
		interfacetests.MethodReturningSeenStruct:                   "ReturnSeen",
		interfacetests.MethodReturningAlterableUint64:              "GetAlterablePrimitiveValue",
		interfacetests.MethodReturningUint64:                       "GetPrimitiveValue",
		interfacetests.MethodReturningUint64Slice:                  "getSliceValue",
		interfacetests.MethodSettingStruct:                         "AddTestStruct",
		interfacetests.MethodSettingUint64:                         "SetAlterablePrimitiveValue",
		interfacetests.MethodTriggeringEvent:                       "TriggerEvent",
		interfacetests.MethodTriggeringEventWithDynamicTopic:       "TriggerEventWithDynamicTopic",
	}
	methodNameMappingByContract[interfacetests.AnySecondContractName] = map[string]string{
		interfacetests.MethodReturningUint64: "GetDifferentPrimitiveValue",
	}

	bm := bindingsMapping{
		contractNameMapping: map[string]string{
			interfacetests.AnyContractName:       contractName,
			interfacetests.AnySecondContractName: contractName,
		},
		methodNameMappingByContract: methodNameMappingByContract,
		contractReaderProxy:         &contractReaderProxy,
		chainWriterProxy:            &chainWriterProxy,
		chainReaderTesters:          map[string]*bindings.ChainReaderTester{},
	}
	contractReaderProxy.bm = &bm
	chainWriterProxy.bm = &bm
	bm.createDelegates()
	return bm
}

func getChainReaderConfig(wrapped *EVMChainComponentsInterfaceTester[*testing.T]) types.ChainReaderConfig {
	testStruct := interfacetests.CreateTestStruct[*testing.T](0, wrapped)
	chainReaderConfig := bindings.NewChainReaderConfig()
	chainReaderConfig.Contracts["ChainReaderTester"].Configs["ReturnSeen"] = &types.ChainReaderDefinition{
		CacheEnabled:      false,
		ChainSpecificName: "returnSeen",
		ReadType:          0,
		InputModifications: codec.ModifiersConfig{
			&codec.HardCodeModifierConfig{
				OnChainValues: map[string]any{
					"BigField":              testStruct.BigField.String(),
					"AccountStruct.Account": hexutil.Encode(testStruct.AccountStruct.Account),
				},
			},
		},
		OutputModifications: codec.ModifiersConfig{
			&codec.HardCodeModifierConfig{OffChainValues: map[string]any{"ExtraField": interfacetests.AnyExtraValue}},
		},
	}
	return chainReaderConfig
}

func getChainWriterConfig(t *testing.T, wrapped *EVMChainComponentsInterfaceTester[*testing.T]) types.ChainWriterConfig {
	return bindings.NewChainWriterConfig(*assets.NewWei(big.NewInt(1000000000000000000)), 2_000_000, wrapped.Helper.Accounts(t)[1].From)
}

func (b bindingClientTester) Name() string {
	return "generated bindings"
}

type bindingClientTester struct {
	interfacetests.ChainComponentsInterfaceTester[*testing.T]
	bindingsMapping *bindingsMapping
}

func (b bindingClientTester) GetContractReader(t *testing.T) commontypes.ContractReader {
	contractReader := b.ChainComponentsInterfaceTester.GetContractReader(t)
	if b.bindingsMapping.contractReaderProxy.ContractReader == nil {
		b.bindingsMapping.contractReaderProxy.ContractReader = contractReader
		b.addDefaultBindings(t)
		for _, tester := range b.bindingsMapping.chainReaderTesters {
			tester.ContractReader = contractReader
		}
	}
	return b.bindingsMapping.contractReaderProxy
}

func (b bindingClientTester) addDefaultBindings(t *testing.T) {
	defaultBindings := b.ChainComponentsInterfaceTester.GetBindings(t)
	for _, binding := range defaultBindings {
		chainReaderTester := b.bindingsMapping.chainReaderTesters[binding.Address]
		if chainReaderTester == nil {
			chainReaderTester = &bindings.ChainReaderTester{
				BoundContract: binding,
				ChainWriter:   b.bindingsMapping.chainWriterProxy.ContractWriter,
			}
			b.bindingsMapping.chainReaderTesters[binding.Address] = chainReaderTester
		} else {
			chainReaderTester.ChainWriter = b.bindingsMapping.chainWriterProxy.ContractWriter
		}
	}
}

func (b bindingClientTester) GetChainWriter(t *testing.T) commontypes.ContractWriter {
	chainWriter := b.ChainComponentsInterfaceTester.GetContractWriter(t)
	if b.bindingsMapping.chainWriterProxy.ContractWriter == nil {
		b.addDefaultBindings(t)
		for _, tester := range b.bindingsMapping.chainReaderTesters {
			tester.ChainWriter = chainWriter
		}
		b.bindingsMapping.chainWriterProxy.ContractWriter = chainWriter
	}
	return b.bindingsMapping.chainWriterProxy
}

type bindingsMapping struct {
	contractNameMapping         map[string]string
	methodNameMappingByContract map[string]map[string]string
	delegates                   map[string]*Delegate
	chainReaderTesters          map[string]*bindings.ChainReaderTester
	contractReaderProxy         *bindingContractReaderProxy
	chainWriterProxy            *bindingChainWriterProxy
}

type bindingContractReaderProxy struct {
	commontypes.ContractReader
	bm *bindingsMapping
}

type bindingChainWriterProxy struct {
	commontypes.ContractWriter
	bm *bindingsMapping
}

func (b bindingContractReaderProxy) Bind(ctx context.Context, boundContracts []commontypes.BoundContract) error {
	updatedBindings := b.bm.translateContractNames(boundContracts)
	for _, updatedBinding := range updatedBindings {
		b.bm.chainReaderTesters[updatedBinding.Address] = &bindings.ChainReaderTester{
			BoundContract:  updatedBinding,
			ContractReader: b.ContractReader,
			ChainWriter:    b.bm.chainWriterProxy.ContractWriter,
		}
	}
	return b.ContractReader.Bind(ctx, updatedBindings)
}

func (b bindingsMapping) translateContractNames(boundContracts []commontypes.BoundContract) []commontypes.BoundContract {
	updatedBindings := make([]commontypes.BoundContract, 0, len(boundContracts))
	for _, boundContract := range boundContracts {
		updatedBindings = append(updatedBindings, commontypes.BoundContract{
			Address: boundContract.Address,
			Name:    b.translateContractName(boundContract.Name),
		})
	}
	return updatedBindings
}

func (b bindingContractReaderProxy) Close() error {
	return b.ContractReader.Close()
}

func (b bindingContractReaderProxy) GetLatestValue(ctx context.Context, readKey string, confidenceLevel primitives.ConfidenceLevel, params, returnVal any) error {
	delegate, err := b.bm.getBindingDelegate(readKey)
	if err != nil {
		return err
	}
	output, err := delegate.apply(ctx, readKey, params, confidenceLevel)
	if err != nil {
		return err
	}
	if output == nil {
		return nil
	}
	err = convertStruct(output, returnVal)
	if err != nil {
		return err
	}
	return nil
}

func (b bindingChainWriterProxy) SubmitTransaction(ctx context.Context, contract, method string, args any, transactionID string, toAddress string, meta *commontypes.TxMeta, value *big.Int) error {
	chainReaderTesters := b.bm.chainReaderTesters[toAddress]
	switch contract {
	case interfacetests.AnyContractName, interfacetests.AnySecondContractName:
		switch method {
		case interfacetests.MethodSettingStruct:
			bindingsInput := bindings.AddTestStructInput{}
			_ = convertStruct(args, &bindingsInput)
			return chainReaderTesters.AddTestStruct(ctx, bindingsInput, transactionID, toAddress, meta)
		case interfacetests.MethodSettingUint64:
			bindingsInput := bindings.SetAlterablePrimitiveValueInput{}
			_ = convertStruct(args, &bindingsInput)
			return chainReaderTesters.SetAlterablePrimitiveValue(ctx, bindingsInput, transactionID, toAddress, meta)
		case interfacetests.MethodTriggeringEvent:
			bindingsInput := bindings.TriggerEventInput{}
			_ = convertStruct(args, &bindingsInput)
			return chainReaderTesters.TriggerEvent(ctx, bindingsInput, transactionID, toAddress, meta)
		case interfacetests.MethodTriggeringEventWithDynamicTopic:
			bindingsInput := bindings.TriggerEventWithDynamicTopicInput{}
			_ = convertStruct(args, &bindingsInput)
			return chainReaderTesters.TriggerEventWithDynamicTopic(ctx, bindingsInput, transactionID, toAddress, meta)
		default:
			return errors.New("No logic implemented for method: " + method)
		}
	default:
		return errors.New("contract with id not supported " + contract)
	}
}

func (b *bindingChainWriterProxy) GetTransactionStatus(ctx context.Context, transactionID string) (commontypes.TransactionStatus, error) {
	return b.ContractWriter.GetTransactionStatus(ctx, transactionID)
}

func removeAddressFromReadIdentifier(s string) string {
	index := strings.Index(s, "-")
	if index == -1 {
		return s
	}
	return s[index+1:]
}

func (b *bindingsMapping) createDelegates() {
	delegates := make(map[string]*Delegate)
	boundContract := commontypes.BoundContract{Address: "", Name: contractName}
	methodTakingLatestParamsKey := removeAddressFromReadIdentifier(boundContract.ReadIdentifier(b.methodNameMappingByContract[interfacetests.AnyContractName][interfacetests.MethodTakingLatestParamsReturningTestStruct]))
	delegates[methodTakingLatestParamsKey] = b.createDelegateForMethodTakingLatestParams()
	methodReturningAlterableUint64Key := removeAddressFromReadIdentifier(boundContract.ReadIdentifier(b.methodNameMappingByContract[interfacetests.AnyContractName][interfacetests.MethodReturningAlterableUint64]))
	delegates[methodReturningAlterableUint64Key] = b.createDelegateForMethodReturningAlterableUint64()
	methodReturningSeenStructKey := removeAddressFromReadIdentifier(boundContract.ReadIdentifier(b.methodNameMappingByContract[interfacetests.AnyContractName][interfacetests.MethodReturningSeenStruct]))
	delegates[methodReturningSeenStructKey] = b.createDelegateForMethodReturningSeenStruct()
	methodReturningUint64Key := removeAddressFromReadIdentifier(boundContract.ReadIdentifier(b.methodNameMappingByContract[interfacetests.AnyContractName][interfacetests.MethodReturningUint64]))
	delegates[methodReturningUint64Key] = b.createDelegateForMethodReturningUint64()
	methodReturningUint64SliceKey := removeAddressFromReadIdentifier(boundContract.ReadIdentifier(b.methodNameMappingByContract[interfacetests.AnyContractName][interfacetests.MethodReturningUint64Slice]))
	delegates[methodReturningUint64SliceKey] = b.createDelegateForMethodReturningUint64Slice()
	methodReturningDifferentUint64Key := removeAddressFromReadIdentifier(boundContract.ReadIdentifier(b.methodNameMappingByContract[interfacetests.AnySecondContractName][interfacetests.MethodReturningUint64]))
	delegates[methodReturningDifferentUint64Key] = b.createDelegateForSecondContractMethodReturningUint64()
	b.delegates = delegates
}

func (b *bindingsMapping) createDelegateForMethodTakingLatestParams() *Delegate {
	delegate := Delegate{inputType: reflect.TypeOf(bindings.GetElementAtIndexInput{})}
	delegate.delegateFunc = func(ctx context.Context, readyKey string, input *any, level primitives.ConfidenceLevel) (any, error) {
		methodInvocation := func(ctx context.Context, readKey string, input *bindings.GetElementAtIndexInput, level primitives.ConfidenceLevel) (any, error) {
			chainReaderTester := b.GetChainReaderTester(readKey)
			return chainReaderTester.GetElementAtIndex(ctx, *input, level)
		}
		return invokeSpecificMethod(ctx, b.translateReadKey(readyKey), (*input).(*bindings.GetElementAtIndexInput), level, methodInvocation)
	}
	return &delegate
}

func (b *bindingsMapping) createDelegateForMethodReturningAlterableUint64() *Delegate {
	delegate := Delegate{}
	delegate.delegateFunc = func(ctx context.Context, readyKey string, input *any, level primitives.ConfidenceLevel) (any, error) {
		methodInvocation := func(ctx context.Context, readKey string, input any, level primitives.ConfidenceLevel) (any, error) {
			chainReaderTester := b.GetChainReaderTester(readKey)
			return chainReaderTester.GetAlterablePrimitiveValue(ctx, level)
		}
		return invokeSpecificMethod(ctx, b.translateReadKey(readyKey), nil, level, methodInvocation)
	}
	return &delegate
}

func (b *bindingsMapping) createDelegateForMethodReturningSeenStruct() *Delegate {
	delegate := Delegate{inputType: reflect.TypeOf(bindings.ReturnSeenInput{})}
	delegate.delegateFunc = func(ctx context.Context, readyKey string, input *any, level primitives.ConfidenceLevel) (any, error) {
		methodInvocation := func(ctx context.Context, readKey string, input *bindings.ReturnSeenInput, level primitives.ConfidenceLevel) (any, error) {
			chainReaderTester := b.GetChainReaderTester(readKey)
			return chainReaderTester.ReturnSeen(ctx, *input, level)
		}
		return invokeSpecificMethod(ctx, b.translateReadKey(readyKey), (*input).(*bindings.ReturnSeenInput), level, methodInvocation)
	}
	return &delegate
}

func (b *bindingsMapping) createDelegateForMethodReturningUint64() *Delegate {
	delegate := Delegate{}
	delegate.delegateFunc = func(ctx context.Context, readyKey string, input *any, level primitives.ConfidenceLevel) (any, error) {
		methodInvocation := func(ctx context.Context, readKey string, input any, level primitives.ConfidenceLevel) (any, error) {
			chainReaderTester := b.GetChainReaderTester(readKey)
			return chainReaderTester.GetPrimitiveValue(ctx, level)
		}
		return invokeSpecificMethod(ctx, b.translateReadKey(readyKey), nil, level, methodInvocation)
	}
	return &delegate
}

func (b *bindingsMapping) createDelegateForMethodReturningUint64Slice() *Delegate {
	delegate := Delegate{}
	delegate.delegateFunc = func(ctx context.Context, readyKey string, input *any, level primitives.ConfidenceLevel) (any, error) {
		methodInvocation := func(ctx context.Context, readKey string, input any, level primitives.ConfidenceLevel) (any, error) {
			chainReaderTester := b.GetChainReaderTester(readKey)
			return chainReaderTester.GetSliceValue(ctx, level)
		}
		return invokeSpecificMethod(ctx, b.translateReadKey(readyKey), nil, level, methodInvocation)
	}
	return &delegate
}

func (b *bindingsMapping) createDelegateForSecondContractMethodReturningUint64() *Delegate {
	delegate := Delegate{}
	delegate.delegateFunc = func(ctx context.Context, readyKey string, input *any, level primitives.ConfidenceLevel) (any, error) {
		methodInvocation := func(ctx context.Context, readKey string, input any, level primitives.ConfidenceLevel) (any, error) {
			chainReaderTester := b.GetChainReaderTester(readKey)
			return chainReaderTester.GetDifferentPrimitiveValue(ctx, level)
		}
		return invokeSpecificMethod(ctx, b.translateReadKey(readyKey), nil, level, methodInvocation)
	}
	return &delegate
}

// Transforms a readKey from ChainReader using the generic testing config to the actual config being used with go bindings which is the auto-generated from the solidity contract.
func (b bindingsMapping) translateReadKey(key string) string {
	var updatedKey = key
	parts := strings.Split(key, "-")
	contractName := parts[1]
	methodName := parts[2]
	for testConfigName, bindingsName := range b.contractNameMapping {
		if contractName == testConfigName {
			updatedKey = strings.Replace(updatedKey, testConfigName, bindingsName, 1)
		}
	}
	for testConfigName, bindingsName := range b.methodNameMappingByContract[contractName] {
		if methodName == testConfigName {
			updatedKey = strings.Replace(updatedKey, testConfigName, bindingsName, 1)
		}
	}
	return updatedKey
}

// Transforms a readKey from ChainReader using the generic testing config to the actual config being used with go bindings which is the auto-generated from the solidity contract.
func (b bindingsMapping) translateContractName(contractName string) string {
	for testContractName, bindingsName := range b.contractNameMapping {
		if contractName == testContractName {
			return bindingsName
		}
	}
	return contractName
}

func invokeSpecificMethod[T any](ctx context.Context, readKey string, input T, level primitives.ConfidenceLevel, methodInvocation func(ctx context.Context, readKey string, input T, level primitives.ConfidenceLevel) (any, error)) (any, error) {
	return methodInvocation(ctx, readKey, input, level)
}

func (b bindingsMapping) getBindingDelegate(readKey string) (*Delegate, error) {
	translatedKey := removeAddressFromReadIdentifier(b.translateReadKey(readKey))
	delegate := b.delegates[translatedKey]

	if delegate == nil {
		return nil, fmt.Errorf("delegate not found for readerKey %s", translatedKey)
	}
	return delegate, nil
}

func (b bindingsMapping) GetChainReaderTester(key string) *bindings.ChainReaderTester {
	address := key[0:strings.Index(key, "-")]
	return b.chainReaderTesters[address]
}

type Delegate struct {
	inputType    reflect.Type
	delegateFunc func(context.Context, string, *any, primitives.ConfidenceLevel) (any, error)
}

func (d Delegate) getInput(input any) (*any, error) {
	if input == nil {
		return nil, nil
	}
	adaptedInput := reflect.New(d.inputType).Interface()
	err := convertStruct(input, adaptedInput)
	if err != nil {
		return nil, err
	}
	return &adaptedInput, nil
}

func (d Delegate) apply(ctx context.Context, readKey string, input any, confidenceLevel primitives.ConfidenceLevel) (any, error) {
	adaptedInput, err := d.getInput(input)
	if err != nil {
		return nil, err
	}
	output, err := d.delegateFunc(ctx, readKey, adaptedInput, confidenceLevel)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// Utility function to converted original types from and to bindings expected types.
func convertStruct(src any, dst any) error {
	if reflect.TypeOf(src).Kind() == reflect.Ptr && reflect.TypeOf(dst).Kind() == reflect.Ptr && reflect.TypeOf(src).Elem() == reflect.TypeOf(interfacetests.LatestParams{}) && reflect.TypeOf(dst).Elem() == reflect.TypeOf(bindings.GetElementAtIndexInput{}) {
		value := src.(*interfacetests.LatestParams).I
		dst.(*bindings.GetElementAtIndexInput).I = big.NewInt(int64(value))
		return nil
	}
	decoder, err := createDecoder(dst)
	if err != nil {
		return err
	}
	err = decoder.Decode(src)
	if err != nil {
		return err
	}
	switch {
	case reflect.TypeOf(dst).Elem() == reflect.TypeOf(interfacetests.TestStructWithExtraField{}):
		destTestStruct := dst.(*interfacetests.TestStructWithExtraField)
		if destTestStruct != nil {
			auxTestStruct := &interfacetests.TestStruct{}
			decoder, _ := createDecoder(auxTestStruct)
			_ = decoder.Decode(src)
			destTestStruct.TestStruct = *auxTestStruct
			sourceTestStruct := src.(bindings.TestStruct)
			destTestStruct.BigField = sourceTestStruct.BigField
			destTestStruct.NestedStaticStruct.Inner.I = int(sourceTestStruct.NestedStaticStruct.Inner.IntVal)
			destTestStruct.NestedStaticStruct.FixedBytes = sourceTestStruct.NestedStaticStruct.FixedBytes
			destTestStruct.NestedDynamicStruct.Inner.I = int(sourceTestStruct.NestedDynamicStruct.Inner.IntVal)
			destTestStruct.NestedDynamicStruct.FixedBytes = sourceTestStruct.NestedDynamicStruct.FixedBytes
			destTestStruct.ExtraField = interfacetests.AnyExtraValue
		}
	case reflect.TypeOf(dst).Elem() == reflect.TypeOf(interfacetests.TestStruct{}):
		destTestStruct := dst.(*interfacetests.TestStruct)
		if destTestStruct != nil {
			sourceTestStruct := src.(bindings.TestStruct)
			destTestStruct.BigField = sourceTestStruct.BigField
			destTestStruct.NestedStaticStruct.Inner.I = int(sourceTestStruct.NestedStaticStruct.Inner.IntVal)
			destTestStruct.NestedStaticStruct.FixedBytes = sourceTestStruct.NestedStaticStruct.FixedBytes
			destTestStruct.NestedDynamicStruct.Inner.I = int(sourceTestStruct.NestedDynamicStruct.Inner.IntVal)
			destTestStruct.NestedDynamicStruct.FixedBytes = sourceTestStruct.NestedDynamicStruct.FixedBytes
		}
	case reflect.TypeOf(src) == reflect.TypeOf(interfacetests.TestStruct{}) && reflect.TypeOf(dst) == reflect.TypeOf(&bindings.AddTestStructInput{}):
		destTestStruct := dst.(*bindings.AddTestStructInput)
		if destTestStruct != nil {
			sourceTestStruct := src.(interfacetests.TestStruct)
			destTestStruct.BigField = sourceTestStruct.BigField
			destTestStruct.NestedStaticStruct.Inner.IntVal = int64(sourceTestStruct.NestedStaticStruct.Inner.I)
			destTestStruct.NestedStaticStruct.FixedBytes = sourceTestStruct.NestedStaticStruct.FixedBytes
			destTestStruct.NestedDynamicStruct.Inner.IntVal = int64(sourceTestStruct.NestedDynamicStruct.Inner.I)
			destTestStruct.NestedDynamicStruct.FixedBytes = sourceTestStruct.NestedDynamicStruct.FixedBytes
		}
	case reflect.TypeOf(src) == reflect.TypeOf(interfacetests.TestStruct{}) && reflect.TypeOf(dst) == reflect.TypeOf(&bindings.ReturnSeenInput{}):
		destTestStruct := dst.(*bindings.ReturnSeenInput)
		if destTestStruct != nil {
			sourceTestStruct := src.(interfacetests.TestStruct)
			destTestStruct.BigField = sourceTestStruct.BigField
			destTestStruct.NestedStaticStruct.Inner.IntVal = int64(sourceTestStruct.NestedStaticStruct.Inner.I)
			destTestStruct.NestedStaticStruct.FixedBytes = sourceTestStruct.NestedStaticStruct.FixedBytes
			destTestStruct.NestedDynamicStruct.Inner.IntVal = int64(sourceTestStruct.NestedDynamicStruct.Inner.I)
			destTestStruct.NestedDynamicStruct.FixedBytes = sourceTestStruct.NestedDynamicStruct.FixedBytes
		}
	}
	return err
}

func createDecoder(dst any) (*mapstructure.Decoder, error) {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringToByteArrayHook,
		),
		Result: dst,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	return decoder, err
}

func stringToByteArrayHook(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() == reflect.String && to == reflect.TypeOf([]byte{}) {
		return evmcodec.EVMAddressModifier{}.DecodeAddress(data.(string))
	}
	if from == reflect.TypeOf([]byte{}) && to.Kind() == reflect.String {
		return evmcodec.EVMAddressModifier{}.EncodeAddress(data.([]byte))
	}
	return data, nil
}
