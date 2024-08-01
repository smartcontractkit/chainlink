package encodings

type Builder interface {
	Bool() TypeCodec
	Int8() TypeCodec
	Int16() TypeCodec
	Int32() TypeCodec
	Int64() TypeCodec
	Uint8() TypeCodec
	Uint16() TypeCodec
	Uint32() TypeCodec
	Uint64() TypeCodec
	String(maxLen uint) (TypeCodec, error)
	Float32() TypeCodec
	Float64() TypeCodec
	OracleID() TypeCodec
	Int(bytes uint) (TypeCodec, error)
	Uint(bytes uint) (TypeCodec, error)
	BigInt(bytes uint, signed bool) (TypeCodec, error)
}
