package abi

import (
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/umbracle/ethgo"
)

// batch of predefined reflect types
var (
	boolT         = reflect.TypeOf(bool(false))
	uint8T        = reflect.TypeOf(uint8(0))
	uint16T       = reflect.TypeOf(uint16(0))
	uint32T       = reflect.TypeOf(uint32(0))
	uint64T       = reflect.TypeOf(uint64(0))
	int8T         = reflect.TypeOf(int8(0))
	int16T        = reflect.TypeOf(int16(0))
	int32T        = reflect.TypeOf(int32(0))
	int64T        = reflect.TypeOf(int64(0))
	addressT      = reflect.TypeOf(ethgo.Address{})
	stringT       = reflect.TypeOf("")
	dynamicBytesT = reflect.SliceOf(reflect.TypeOf(byte(0)))
	functionT     = reflect.ArrayOf(24, reflect.TypeOf(byte(0)))
	tupleT        = reflect.TypeOf(map[string]interface{}{})
	bigIntT       = reflect.TypeOf(new(big.Int))
)

// Kind represents the kind of abi type
type Kind int

const (
	// KindBool is a boolean
	KindBool Kind = iota

	// KindUInt is an uint
	KindUInt

	// KindInt is an int
	KindInt

	// KindString is a string
	KindString

	// KindArray is an array
	KindArray

	// KindSlice is a slice
	KindSlice

	// KindAddress is an address
	KindAddress

	// KindBytes is a bytes array
	KindBytes

	// KindFixedBytes is a fixed bytes
	KindFixedBytes

	// KindFixedPoint is a fixed point
	KindFixedPoint

	// KindTuple is a tuple
	KindTuple

	// KindFunction is a function
	KindFunction
)

func (k Kind) String() string {
	names := [...]string{
		"Bool",
		"Uint",
		"Int",
		"String",
		"Array",
		"Slice",
		"Address",
		"Bytes",
		"FixedBytes",
		"FixedPoint",
		"Tuple",
		"Function",
	}

	return names[k]
}

// TupleElem is an element of a tuple
type TupleElem struct {
	Name    string
	Elem    *Type
	Indexed bool
}

// Type is an ABI type
type Type struct {
	kind  Kind
	size  int
	elem  *Type
	tuple []*TupleElem
	t     reflect.Type
}

func NewTupleType(inputs []*TupleElem) *Type {
	return &Type{
		kind:  KindTuple,
		tuple: inputs,
		t:     tupleT,
	}
}

func NewTupleTypeFromArgs(inputs []*ArgumentStr) (*Type, error) {
	elems := []*TupleElem{}
	for _, i := range inputs {
		typ, err := NewTypeFromArgument(i)
		if err != nil {
			return nil, err
		}
		elems = append(elems, &TupleElem{
			Name:    i.Name,
			Elem:    typ,
			Indexed: i.Indexed,
		})
	}
	return NewTupleType(elems), nil
}

// ParseLog parses a log using this type
func (t *Type) ParseLog(log *ethgo.Log) (map[string]interface{}, error) {
	return ParseLog(t, log)
}

// Decode decodes an object using this type
func (t *Type) Decode(input []byte) (interface{}, error) {
	return Decode(t, input)
}

// DecodeStruct decodes an object using this type to the out param
func (t *Type) DecodeStruct(input []byte, out interface{}) error {
	return DecodeStruct(t, input, out)
}

// Encode encodes an object using this type
func (t *Type) Encode(v interface{}) ([]byte, error) {
	return Encode(v, t)
}

func (t *Type) String() string {
	return t.Format(false)
}

// String returns the raw representation of the type
func (t *Type) Format(includeArgs bool) string {
	switch t.kind {
	case KindTuple:
		rawAux := []string{}
		for _, i := range t.TupleElems() {
			name := i.Elem.Format(includeArgs)
			if i.Indexed {
				name += " indexed"
			}
			if includeArgs {
				if i.Name != "" {
					name += " " + i.Name
				}
			}
			rawAux = append(rawAux, name)
		}
		return fmt.Sprintf("tuple(%s)", strings.Join(rawAux, ","))

	case KindArray:
		return fmt.Sprintf("%s[%d]", t.elem.Format(includeArgs), t.size)

	case KindSlice:
		return fmt.Sprintf("%s[]", t.elem.Format(includeArgs))

	case KindBytes:
		return "bytes"

	case KindFixedBytes:
		return fmt.Sprintf("bytes%d", t.size)

	case KindString:
		return "string"

	case KindBool:
		return "bool"

	case KindAddress:
		return "address"

	case KindFunction:
		return "function"

	case KindUInt:
		return fmt.Sprintf("uint%d", t.size)

	case KindInt:
		return fmt.Sprintf("int%d", t.size)

	default:
		panic(fmt.Errorf("BUG: abi type not found %s", t.kind.String()))
	}
}

// Elem returns the elem value for slice and arrays
func (t *Type) Elem() *Type {
	return t.elem
}

// Size returns the size of the type
func (t *Type) Size() int {
	return t.size
}

// TupleElems returns the elems of the tuple
func (t *Type) TupleElems() []*TupleElem {
	return t.tuple
}

// GoType returns the go type
func (t *Type) GoType() reflect.Type {
	return t.t
}

// Kind returns the kind of the type
func (t *Type) Kind() Kind {
	return t.kind
}

func (t *Type) isVariableInput() bool {
	return t.kind == KindSlice || t.kind == KindBytes || t.kind == KindString
}

func (t *Type) isDynamicType() bool {
	if t.kind == KindTuple {
		for _, elem := range t.tuple {
			if elem.Elem.isDynamicType() {
				return true
			}
		}
		return false
	}
	return t.kind == KindString || t.kind == KindBytes || t.kind == KindSlice || (t.kind == KindArray && t.elem.isDynamicType())
}

func parseType(arg *ArgumentStr) (string, error) {
	if !strings.HasPrefix(arg.Type, "tuple") {
		return arg.Type, nil
	}

	if len(arg.Components) == 0 {
		return "tuple()", nil
	}

	// parse the arg components from the tuple
	str := []string{}
	for _, i := range arg.Components {
		aux, err := parseType(i)
		if err != nil {
			return "", err
		}
		if i.Indexed {
			str = append(str, aux+" indexed "+i.Name)
		} else {
			str = append(str, aux+" "+i.Name)
		}
	}
	return fmt.Sprintf("tuple(%s)%s", strings.Join(str, ","), strings.TrimPrefix(arg.Type, "tuple")), nil
}

// NewTypeFromArgument parses an abi type from an argument
func NewTypeFromArgument(arg *ArgumentStr) (*Type, error) {
	str, err := parseType(arg)
	if err != nil {
		return nil, err
	}
	return NewType(str)
}

// NewType parses a type in string format
func NewType(s string) (*Type, error) {
	l := newLexer(s)
	l.nextToken()

	return readType(l)
}

// MustNewType parses a type in string format or panics if its invalid
func MustNewType(s string) *Type {
	t, err := NewType(s)
	if err != nil {
		panic(err)
	}
	return t
}

func getTypeSize(t *Type) int {
	if t.kind == KindArray && !t.elem.isDynamicType() {
		if t.elem.kind == KindArray || t.elem.kind == KindTuple {
			return t.size * getTypeSize(t.elem)
		}
		return t.size * 32
	} else if t.kind == KindTuple && !t.isDynamicType() {
		total := 0
		for _, elem := range t.tuple {
			total += getTypeSize(elem.Elem)
		}
		return total
	}
	return 32
}

var typeRegexp = regexp.MustCompile("^([[:alpha:]]+)([[:digit:]]*)$")

func expectedToken(t tokenType) error {
	return fmt.Errorf("expected token %s", t.String())
}

func notExpectedToken(t tokenType) error {
	return fmt.Errorf("token '%s' not expected", t.String())
}

func readType(l *lexer) (*Type, error) {
	var tt *Type

	tok := l.nextToken()

	isTuple := false
	if tok.typ == tupleToken {
		if l.nextToken().typ != lparenToken {
			return nil, expectedToken(lparenToken)
		}
		isTuple = true
	} else if tok.typ == lparenToken {
		isTuple = true
	}
	if isTuple {
		var next token
		elems := []*TupleElem{}
		for {

			name := ""
			indexed := false

			elem, err := readType(l)
			if err != nil {
				if l.current.typ == rparenToken && len(elems) == 0 {
					// empty tuple 'tuple()'
					break
				}
				return nil, fmt.Errorf("failed to decode type: %v", err)
			}

			switch l.peek.typ {
			case strToken:
				l.nextToken()
				name = l.current.literal

			case indexedToken:
				l.nextToken()
				indexed = true
				if l.peek.typ == strToken {
					l.nextToken()
					name = l.current.literal
				}
			}

			elems = append(elems, &TupleElem{
				Name:    name,
				Elem:    elem,
				Indexed: indexed,
			})

			next = l.nextToken()
			if next.typ == commaToken {
				continue
			} else if next.typ == rparenToken {
				break
			} else {
				return nil, notExpectedToken(next.typ)
			}
		}
		tt = &Type{kind: KindTuple, tuple: elems, t: tupleT}

	} else if tok.typ != strToken {
		return nil, expectedToken(strToken)

	} else {
		// Check normal types
		elem, err := decodeSimpleType(tok.literal)
		if err != nil {
			return nil, err
		}
		tt = elem
	}

	// check for arrays at the end of the type
	for {
		if l.peek.typ != lbracketToken {
			break
		}

		l.nextToken()
		n := l.nextToken()

		var tAux *Type
		if n.typ == rbracketToken {
			tAux = &Type{kind: KindSlice, elem: tt, t: reflect.SliceOf(tt.t)}

		} else if n.typ == numberToken {
			size, err := strconv.ParseUint(n.literal, 10, 32)
			if err != nil {
				return nil, fmt.Errorf("failed to read array size '%s': %v", n.literal, err)
			}

			tAux = &Type{kind: KindArray, elem: tt, size: int(size), t: reflect.ArrayOf(int(size), tt.t)}
			if l.nextToken().typ != rbracketToken {
				return nil, expectedToken(rbracketToken)
			}
		} else {
			return nil, notExpectedToken(n.typ)
		}

		tt = tAux
	}
	return tt, nil
}

func decodeSimpleType(str string) (*Type, error) {
	match := typeRegexp.FindStringSubmatch(str)
	if len(match) == 0 {
		return nil, fmt.Errorf("type format is incorrect. Expected 'type''bytes' but found '%s'", str)
	}
	match = match[1:]

	var err error
	t := match[0]

	bytes := 0
	ok := false

	if bytesStr := match[1]; bytesStr != "" {
		bytes, err = strconv.Atoi(bytesStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bytes '%s': %v", bytesStr, err)
		}
		ok = true
	}

	// int and uint without bytes default to 256, 'bytes' may
	// have or not, the rest dont have bytes
	if t == "int" || t == "uint" {
		if !ok {
			bytes = 256
		}
	} else if t != "bytes" && ok {
		return nil, fmt.Errorf("type %s does not expect bytes", t)
	}

	switch t {
	case "uint":
		var k reflect.Type
		switch bytes {
		case 8:
			k = uint8T
		case 16:
			k = uint16T
		case 32:
			k = uint32T
		case 64:
			k = uint64T
		default:
			if bytes%8 != 0 {
				panic(fmt.Errorf("number of bytes has to be M mod 8"))
			}
			k = bigIntT
		}
		return &Type{kind: KindUInt, size: int(bytes), t: k}, nil

	case "int":
		var k reflect.Type
		switch bytes {
		case 8:
			k = int8T
		case 16:
			k = int16T
		case 32:
			k = int32T
		case 64:
			k = int64T
		default:
			if bytes%8 != 0 {
				panic(fmt.Errorf("number of bytes has to be M mod 8"))
			}
			k = bigIntT
		}
		return &Type{kind: KindInt, size: int(bytes), t: k}, nil

	case "byte":
		bytes = 1
		fallthrough

	case "bytes":
		if bytes == 0 {
			return &Type{kind: KindBytes, t: dynamicBytesT}, nil
		}
		return &Type{kind: KindFixedBytes, size: int(bytes), t: reflect.ArrayOf(int(bytes), reflect.TypeOf(byte(0)))}, nil

	case "string":
		return &Type{kind: KindString, t: stringT}, nil

	case "bool":
		return &Type{kind: KindBool, t: boolT}, nil

	case "address":
		return &Type{kind: KindAddress, t: addressT, size: 20}, nil

	case "function":
		return &Type{kind: KindFunction, size: 24, t: functionT}, nil

	default:
		return nil, fmt.Errorf("unknown type '%s'", t)
	}
}

type tokenType int

const (
	eofToken tokenType = iota
	strToken
	numberToken
	tupleToken
	lparenToken
	rparenToken
	lbracketToken
	rbracketToken
	commaToken
	indexedToken
	invalidToken
)

func (t tokenType) String() string {
	names := [...]string{
		"eof",
		"string",
		"number",
		"tuple",
		"(",
		")",
		"[",
		"]",
		",",
		"indexed",
		"<invalid>",
	}
	return names[t]
}

type token struct {
	typ     tokenType
	literal string
}

type lexer struct {
	input        string
	current      token
	peek         token
	position     int
	readPosition int
	ch           byte
}

func newLexer(input string) *lexer {
	l := &lexer{input: input}
	l.readChar()
	return l
}

func (l *lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *lexer) nextToken() token {
	l.current = l.peek
	l.peek = l.nextTokenImpl()
	return l.current
}

func (l *lexer) nextTokenImpl() token {
	var tok token

	// skip whitespace
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}

	switch l.ch {
	case ',':
		tok.typ = commaToken
	case '(':
		tok.typ = lparenToken
	case ')':
		tok.typ = rparenToken
	case '[':
		tok.typ = lbracketToken
	case ']':
		tok.typ = rbracketToken
	case 0:
		tok.typ = eofToken
	default:
		if isLetter(l.ch) {
			tok.literal = l.readIdentifier()
			if tok.literal == "tuple" {
				tok.typ = tupleToken
			} else if tok.literal == "indexed" {
				tok.typ = indexedToken
			} else {
				tok.typ = strToken
			}

			return tok
		} else if isDigit(l.ch) {
			return token{numberToken, l.readNumber()}
		} else {
			tok.typ = invalidToken
		}
	}

	l.readChar()
	return tok
}

func (l *lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
