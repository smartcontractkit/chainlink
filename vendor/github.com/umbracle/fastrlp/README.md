
# FastRlp

FastRlp is a high performant encoding/decoding library for the RLP Ethereum format. This library is based on [fastjson](https://github.com/valyala/fastjson).

## Usage

FastRlp does not uses reflect to avoid bottlenecks. It provides a single value primitive that can be encoded or decoded into any specific type.

Encode: 

```
a := &fastrlp.Arena{}

// Encode a uint
v := a.NewUint(300)
buf := v.MarshalTo(nil)

// Encode an array
v = a.NewArray()
v.Set(a.NewUint(300))
buf = v.MarshalTo(nil)
```

You can find more examples [here](https://github.com/umbracle/fastrlp/blob/master/arena_test.go#L53).

Decode:

```
p := &fastrlp.Parser{}
v, err := p.Parse([]byte{0x01})
if err != nil {
    panic(err)
}

num, err := v.GetUint64()
if err != nil {
    panic(err)
}
fmt.Println(num)
```

## Benchmark

```
$ go-rlp-test go test -v ./. -run=XX -bench=.     
goos: linux
goarch: amd64
pkg: github.com/ferranbt/go-rlp-test
BenchmarkDecode100HeadersGeth-8      	   10000	    196183 ns/op	   32638 B/op	    1002 allocs/op
BenchmarkEncode100HeadersGeth-8      	   10000	    179328 ns/op	   88471 B/op	    1003 allocs/op
BenchmarkDecode100HeadersFastRlp-8   	   30000	     57179 ns/op	      16 B/op	       0 allocs/op
BenchmarkEncode100HeadersFastRlp-8   	   30000	     43967 ns/op	      23 B/op	       0 allocs/op
PASS
ok  	github.com/ferranbt/go-rlp-test	7.890s
```
