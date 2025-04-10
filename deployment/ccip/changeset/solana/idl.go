package solana

import (
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	. "github.com/gagliardetto/utilz"
)

// https://github.com/project-serum/anchor/blob/97e9e03fb041b8b888a9876a7c0676d9bb4736f3/ts/src/idl.ts
type IDL struct {
	Version string `json:"version"`
	Name    string `json:"name"`
	// Docs         []string         `json:"docs"` // @custom
	Instructions []IdlInstruction `json:"instructions"`
	State        *IdlState        `json:"state,omitempty"`
	Accounts     IdlTypeDefSlice  `json:"accounts,omitempty"`
	Types        IdlTypeDefSlice  `json:"types,omitempty"`
	Events       []IdlEvent       `json:"events,omitempty"`
	Errors       []IdlErrorCode   `json:"errors,omitempty"`
	Constants    []IdlConstant    `json:"constants,omitempty"`

	Metadata *IdlMetadata `json:"metadata,omitempty"` // NOTE: deprecated
}

// TODO: write generator
type IdlConstant struct {
	Name  string
	Type  IdlType
	Value string
}

type IdlMetadata struct {
	Address string `json:"address"`
}

type IdlTypeDefSlice []IdlTypeDef

func (named IdlTypeDefSlice) GetByName(name string) *IdlTypeDef {
	for i := range named {
		v := named[i]
		if v.Name == name {
			return &v
		}
	}
	return nil
}

// Validate validates and IDL
func (idl *IDL) Validate() error {
	// TODO
	return nil
}

type IdlEvent struct {
	Name   string          `json:"name"`
	Fields []IdlEventField `json:"fields"`
}

type IdlEventField struct {
	Name  string  `json:"name"`
	Type  IdlType `json:"type"`
	Index bool    `json:"index"`
}

type IdlInstruction struct {
	Name     string              `json:"name"`
	Docs     []string            `json:"docs"` // @custom
	Accounts IdlAccountItemSlice `json:"accounts"`
	Args     []IdlField          `json:"args"`
}

type IdlAccountItemSlice []IdlAccountItem

func (slice IdlAccountItemSlice) NumAccounts() (count int) {

	for _, item := range slice {
		if item.IdlAccount != nil {
			count++
		}

		if item.IdlAccounts != nil {
			count += item.IdlAccounts.Accounts.NumAccounts()
		}
	}

	return count
}

func (slice IdlAccountItemSlice) Walk(
	parentGroupPath string,
	previousIndex *int,
	parentGroup *IdlAccounts,
	callback func(string, int, *IdlAccounts, *IdlAccount) bool,
) {
	defaultVal := -1
	if previousIndex == nil {
		previousIndex = &defaultVal
	}
	for _, item := range slice {
		item.Walk(parentGroupPath, previousIndex, parentGroup, callback)
	}
}

type IdlState struct {
	Struct  IdlTypeDef       `json:"struct"`
	Methods []IdlStateMethod `json:"methods"`
}

type IdlStateMethod = IdlInstruction

// type IdlAccountItem = IdlAccount | IdlAccounts;
type IdlAccountItem struct {
	IdlAccount  *IdlAccount
	IdlAccounts *IdlAccounts
}

func (item IdlAccountItem) Walk(
	parentGroupPath string,
	previousIndex *int,
	parentGroup *IdlAccounts,
	callback func(string, int, *IdlAccounts, *IdlAccount) bool,
) {
	defaultVal := -1
	if previousIndex == nil {
		previousIndex = &defaultVal
	}
	if item.IdlAccount != nil {
		*previousIndex++
		doContinue := callback(parentGroupPath, *previousIndex, parentGroup, item.IdlAccount)
		if !doContinue {
			return
		}
	}

	if item.IdlAccounts != nil {
		var thisGroupName string
		if parentGroupPath == "" {
			thisGroupName = item.IdlAccounts.Name
		} else {
			thisGroupName = parentGroupPath + "/" + item.IdlAccounts.Name
		}
		item.IdlAccounts.Accounts.Walk(thisGroupName, previousIndex, item.IdlAccounts, callback)
	}
}

// TODO: verify with examples
func (env *IdlAccountItem) UnmarshalJSON(data []byte) error {

	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp == nil {
		return fmt.Errorf("envelope is nil: %v", env)
	}

	switch v := temp.(type) {
	case map[string]interface{}:
		{
			// Ln(ShakespeareBG("::IdlAccountItem"))
			// spew.Dump(v)

			if len(v) == 0 {
				return nil
			}

			// Multiple accounts:
			if _, ok := v["accounts"]; ok {
				if err := TranscodeJSON(temp, &env.IdlAccounts); err != nil {
					return err
				}
			}
			// Single account:
			// TODO: check both isMut and isSigner
			if _, ok := v["isMut"]; ok {
				if err := TranscodeJSON(temp, &env.IdlAccount); err != nil {
					return err
				}
			}

			// panic(Sf("what is this?:\n%s", spew.Sdump(temp)))
		}
	default:
		return fmt.Errorf("Unknown kind: %s", spew.Sdump(temp))
	}

	return nil
}

type IdlAccount struct {
	Docs     []string `json:"docs"` // @custom
	Name     string   `json:"name"`
	IsMut    bool     `json:"isMut"`
	IsSigner bool     `json:"isSigner"`
	Optional bool     `json:"optional"` // @custom
}

// A nested/recursive version of IdlAccount.
type IdlAccounts struct {
	Name     string              `json:"name"`
	Docs     []string            `json:"docs"` // @custom
	Accounts IdlAccountItemSlice `json:"accounts"`
}

type IdlField struct {
	Name string   `json:"name"`
	Docs []string `json:"docs"` // @custom
	Type IdlType  `json:"type"`
}

type IdlTypeAsString string

const (
	IdlTypeBool      IdlTypeAsString = "bool"
	IdlTypeU8        IdlTypeAsString = "u8"
	IdlTypeI8        IdlTypeAsString = "i8"
	IdlTypeU16       IdlTypeAsString = "u16"
	IdlTypeI16       IdlTypeAsString = "i16"
	IdlTypeU32       IdlTypeAsString = "u32"
	IdlTypeI32       IdlTypeAsString = "i32"
	IdlTypeU64       IdlTypeAsString = "u64"
	IdlTypeI64       IdlTypeAsString = "i64"
	IdlTypeU128      IdlTypeAsString = "u128"
	IdlTypeI128      IdlTypeAsString = "i128"
	IdlTypeBytes     IdlTypeAsString = "bytes"
	IdlTypeString    IdlTypeAsString = "string"
	IdlTypePublicKey IdlTypeAsString = "publicKey"

	// Custom additions:
	IdlTypeUnixTimestamp IdlTypeAsString = "unixTimestamp"
	IdlTypeHash          IdlTypeAsString = "hash"
	IdlTypeDuration      IdlTypeAsString = "duration"

	// | IdlTypeVec
	// | IdlTypeOption
	// | IdlTypeDefined;
)

type IdlTypeVec struct {
	Vec IdlType `json:"vec"`
}

type IdlTypeOption struct {
	Option IdlType `json:"option"`
}

// User defined type.
type IdlTypeDefined struct {
	Defined string `json:"defined"`
}

// Wrapper type:
type IdlTypeArray struct {
	Thing IdlType
	Num   int
}

func (env *IdlType) UnmarshalJSON(data []byte) error {

	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp == nil {
		return fmt.Errorf("envelope is nil: %v", env)
	}

	switch v := temp.(type) {
	case string:
		{
			env.asString = IdlTypeAsString(v)
		}
	case map[string]interface{}:
		{
			// Ln(PurpleBG("::IdlType"))
			// spew.Dump(v)

			if len(v) == 0 {
				return nil
			}

			if _, ok := v["vec"]; ok {
				var target IdlTypeVec
				if err := TranscodeJSON(temp, &target); err != nil {
					return err
				}
				env.asIdlTypeVec = &target
			}
			if _, ok := v["option"]; ok {
				var target IdlTypeOption
				if err := TranscodeJSON(temp, &target); err != nil {
					return err
				}
				env.asIdlTypeOption = &target
			}
			if _, ok := v["defined"]; ok {
				var target IdlTypeDefined
				if err := TranscodeJSON(temp, &target); err != nil {
					return err
				}
				env.asIdlTypeDefined = &target
			}
			if got, ok := v["array"]; ok {

				if _, ok := got.([]interface{}); !ok {
					panic(Sf("array is not in expected format:\n%s", spew.Sdump(got)))
				}
				arrVal := got.([]interface{})
				if len(arrVal) != 2 {
					panic(Sf("array is not of expected length:\n%s", spew.Sdump(got)))
				}
				var target IdlTypeArray
				if err := TranscodeJSON(arrVal[0], &target.Thing); err != nil {
					return err
				}

				target.Num = int(arrVal[1].(float64))

				env.asIdlTypeArray = &target
			}
			// panic(Sf("what is this?:\n%s", spew.Sdump(temp)))
		}
	default:
		return fmt.Errorf("Unknown kind: %s", spew.Sdump(temp))
	}

	return nil
}

// Wrapper type:
type IdlType struct {
	asString         IdlTypeAsString
	asIdlTypeVec     *IdlTypeVec
	asIdlTypeOption  *IdlTypeOption
	asIdlTypeDefined *IdlTypeDefined
	asIdlTypeArray   *IdlTypeArray
}

func (env *IdlType) IsString() bool {
	return env.asString != ""
}
func (env *IdlType) IsIdlTypeVec() bool {
	return env.asIdlTypeVec != nil
}
func (env *IdlType) IsIdlTypeOption() bool {
	return env.asIdlTypeOption != nil
}
func (env *IdlType) IsIdlTypeDefined() bool {
	return env.asIdlTypeDefined != nil
}
func (env *IdlType) IsArray() bool {
	return env.asIdlTypeArray != nil
}

// Getters:
func (env *IdlType) GetString() IdlTypeAsString {
	return env.asString
}
func (env *IdlType) GetIdlTypeVec() *IdlTypeVec {
	return env.asIdlTypeVec
}
func (env *IdlType) GetIdlTypeOption() *IdlTypeOption {
	return env.asIdlTypeOption
}
func (env *IdlType) GetIdlTypeDefined() *IdlTypeDefined {
	return env.asIdlTypeDefined
}
func (env *IdlType) GetArray() *IdlTypeArray {
	return env.asIdlTypeArray
}

type IdlTypeDef struct {
	Name string       `json:"name"`
	Type IdlTypeDefTy `json:"type"`
}

type IdlTypeDefTyKind string

const (
	IdlTypeDefTyKindStruct IdlTypeDefTyKind = "struct"
	IdlTypeDefTyKindEnum   IdlTypeDefTyKind = "enum"
)

// TODO:
type IdlTypeDefTyStruct struct {
	Kind IdlTypeDefTyKind `json:"kind"` // == "struct"

	Fields *IdlTypeDefStruct `json:"fields,omitempty"`
}

// TODO:
type IdlTypeDefTyEnum struct {
	Kind IdlTypeDefTyKind `json:"kind"` // == "enum"

	Variants IdlEnumVariantSlice `json:"variants,omitempty"`
}

type IdlTypeDefTy struct {
	Kind IdlTypeDefTyKind `json:"kind"`

	Fields   *IdlTypeDefStruct   `json:"fields,omitempty"`
	Variants IdlEnumVariantSlice `json:"variants,omitempty"`
}

type IdlEnumVariantSlice []IdlEnumVariant

func (slice IdlEnumVariantSlice) IsAllUint8() bool {
	for _, elem := range slice {
		if !elem.IsUint8() {
			return false
		}
	}
	return true
}

func (slice IdlEnumVariantSlice) IsSimpleEnum() bool {
	return slice.IsAllUint8()
}

type IdlTypeDefStruct = []IdlField

type IdlEnumVariant struct {
	Name   string         `json:"name"`
	Docs   []string       `json:"docs"` // @custom
	Fields *IdlEnumFields `json:"fields,omitempty"`
}

func (variant *IdlEnumVariant) IsUint8() bool {
	// it's a simple uint8 if there is no fields data
	return variant.Fields == nil
}

// TODO
// type IdlEnumFields = IdlEnumFieldsNamed | IdlEnumFieldsTuple;
type IdlEnumFields struct {
	IdlEnumFieldsNamed *IdlEnumFieldsNamed
	IdlEnumFieldsTuple *IdlEnumFieldsTuple
}

type IdlEnumFieldsNamed []IdlField

type IdlEnumFieldsTuple []IdlType

// TODO: verify with examples
func (env *IdlEnumFields) UnmarshalJSON(data []byte) error {

	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp == nil {
		return fmt.Errorf("envelope is nil: %v", env)
	}

	switch v := temp.(type) {
	case []interface{}:
		{
			// Ln(LimeBG("::IdlEnumFields"))
			// spew.Dump(v)

			if len(v) == 0 {
				return nil
			}

			firstItem := v[0]

			if _, ok := firstItem.(map[string]interface{})["name"]; ok {
				// TODO:
				// If has `name` field, then it's most likely a IdlEnumFieldsNamed.
				if err := TranscodeJSON(temp, &env.IdlEnumFieldsNamed); err != nil {
					return err
				}
			} else {
				if err := TranscodeJSON(temp, &env.IdlEnumFieldsTuple); err != nil {
					return err
				}
			}

			// panic(Sf("what is this?:\n%s", spew.Sdump(temp)))
		}
	default:
		return fmt.Errorf("Unknown kind: %s", spew.Sdump(temp))
	}

	return nil
}

type IdlErrorCode struct {
	Code int    `json:"code"`
	Name string `json:"name"`
	Msg  string `json:"msg,omitempty"`
}
