package evm

import (
	"fmt"
	"github.com/smartcontractkit/smart-contract-spec/internal/utils"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Field struct {
	Name         string
	SolidityName string
	Type         GoType
}

type GoType struct {
	IsStruct      bool
	StructDetails Struct
	PrimitiveType string
}

func (g GoType) Print() string {
	if g.IsStruct {
		return g.StructDetails.Name
	}
	return g.PrimitiveType
}

type Struct struct {
	Name   string
	Fields []Field
}

type Function struct {
	Name                string
	SolidityName        string
	Input               *Struct
	Output              *GoType
	RequiresTransaction bool
	IsPayable           bool
	HasOutput           bool
	HasInput            bool
}

type CodeDetails struct {
	ABI                string
	Functions          []Function
	Structs            *map[string]Struct
	ContractStructName string
	ContractName       string
}

// TODO review structs param - at least should
func solidityStructToGoStruct(abiType abi.Type, structs *map[string]Struct, owner string, nameProvider func() string) (Struct, error) {
	if abiType.TupleType == nil {
		return Struct{}, nil
	}

	var structName string
	if abiType.TupleRawName != "" {
		name, err := solidityNameToGoName(abiType.TupleRawName, owner)
		if err != nil {
			return Struct{}, err
		}
		structName = name
	} else {
		structName = nameProvider()
	}

	requiredStruct, exists := (*structs)[structName]
	if exists {
		return requiredStruct, nil
	}

	fields := []Field{}
	for i := 0; i < len(abiType.TupleElems); i++ {
		fieldAbiType := abiType.TupleElems[i]
		fieldAbiName := abiType.TupleRawNames[i]
		goFieldName, err := solidityNameToGoName(fieldAbiName, owner)
		if err != nil {
			return Struct{}, err
		}

		owner := fmt.Sprintf("%s.%s", owner, structName)

		structNameProvider := func() string {
			return fmt.Sprintf("%s%s%d", structName, "Inner", i)
		}

		goType, err := abiTypeToGoType(*fieldAbiType, structs, owner, structNameProvider)
		fields = append(fields, Field{
			Name:         goFieldName,
			Type:         goType,
			SolidityName: fieldAbiName,
		})
	}

	s := Struct{
		Name:   structName,
		Fields: fields,
	}
	(*structs)[structName] = s
	return s, nil
}

// For simplicity every contract method will receive a single struct and return a single struct
// TODO CR already receives all inputs within a struct but as response it can return an struct or a primitive value. We should fix that as well.
func getParam(arguments abi.Arguments, structs *map[string]Struct, contractName string, structName string) (*Struct, error) {
	if len(arguments) == 0 {
		return nil, nil
	}
	fields := []Field{}
	owner := fmt.Sprintf("%s.%s", contractName, contractName)
	for _, argument := range arguments {
		GoName, err := solidityNameToGoName(argument.Name, owner)
		nameProvider := func() string {
			return structName
		}
		goType, err := abiTypeToGoType(argument.Type, structs, owner, nameProvider)
		if err != nil {
			return nil, err
		}
		fields = append(fields, Field{
			Name:         GoName,
			Type:         goType,
			SolidityName: argument.Name,
		})
	}
	//If there's only one parameter of type struct do not wrap it in another struct.
	if len(fields) == 1 && fields[0].Type.IsStruct {
		return &fields[0].Type.StructDetails, nil
	}
	s := Struct{
		Fields: fields,
		Name:   structName,
	}
	(*structs)[structName] = s
	return &s, nil

}

func abiTypeToGoType(abiType abi.Type, structs *map[string]Struct, owner string, nameProvider func() string) (GoType, error) {
	goStruct, err := solidityStructToGoStruct(abiType, structs, owner, nameProvider)
	if err != nil {
		return GoType{}, err
	}
	return GoType{
		IsStruct:      abiType.TupleType != nil,
		StructDetails: goStruct,
		PrimitiveType: SolidityPrimitiveTypeToGoType(abiType),
	}, nil
}

func ConvertABIToCodeDetails(contractName string, abiJSON string) (CodeDetails, error) {
	// Parse the ABI
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	// Initialize a CodeDetails struct
	codeDetails := CodeDetails{
		ABI:                abiJSON,
		Functions:          make([]Function, 0),
		ContractStructName: contractName,
		ContractName:       contractName,
	}

	structs := &map[string]Struct{}

	// Iterate over all the methods in the ABI
	for methodName, method := range parsedABI.Methods {

		//TODO fix this, logic may not be the same
		owner := fmt.Sprintf("%s.%s", contractName, methodName)
		fixedMethodName, err := solidityNameToGoName(methodName, owner)
		if err != nil {
			return CodeDetails{}, err
		}
		// Initialize a Function struct for this method
		inputParam, err := getParam(method.Inputs, structs, owner, fixedMethodName+"Input")
		if err != nil {
			return CodeDetails{}, err
		}
		outputParam, err := getParam(method.Outputs, structs, owner, fixedMethodName+"Output")
		if err != nil {
			return CodeDetails{}, err
		}
		requiresTransaction, isPayable := requiresTransactionAndPayable(method.StateMutability)
		functionDetails := Function{
			Name:                fixedMethodName,
			SolidityName:        methodName,
			Input:               inputParam,
			Output:              convertToPrimitiveIfNeeded(outputParam),
			RequiresTransaction: requiresTransaction,
			IsPayable:           isPayable,
		}

		// Append the function details to the contract
		codeDetails.Functions = append(codeDetails.Functions, functionDetails)
	}

	codeDetails.Structs = structs
	return codeDetails, nil
}

func convertToPrimitiveIfNeeded(param *Struct) *GoType {
	if param == nil {
		return nil
	}
	if len(param.Fields) == 1 {
		field := param.Fields[0]
		if !field.Type.IsStruct {
			return &field.Type
		}
	}
	return &GoType{
		IsStruct:      true,
		StructDetails: *param,
	}
}

// isReadOnly takes the stateMutability string and returns true if the method is read-only, false if it requires a transaction
func requiresTransactionAndPayable(stateMutability string) (bool, bool) {
	requiresTransaction := false
	isPayable := false
	switch stateMutability {
	case "view", "pure":
		requiresTransaction = false
	case "payable":
		isPayable = true
		requiresTransaction = true
	case "nonpayable":
		requiresTransaction = true
	}
	return requiresTransaction, isPayable
}

func solidityNameToGoName(name string, owner string) (string, error) {
	if !utils.IsValidIdentifier(name) {
		return "", fmt.Errorf("invalid identifier for parameter: %s in %s, only digits and letters and underscore are allows", name, owner)
	}
	//Case of a single primitive return type.
	if name == "" {
		name = "Value"
	}
	if paramNeedsRenaming(name) {
		return getExpectedParameterName(name), nil
	}

	return name, nil
}

func getExpectedParameterName(paramName string) string {
	return utils.CapitalizeFirstLetter(utils.RemoveUnderscore(paramName))
}

func paramNeedsRenaming(paramName string) bool {
	return utils.StartsWithLowerCase(paramName) || utils.ContainsUnderscore(paramName)
}
