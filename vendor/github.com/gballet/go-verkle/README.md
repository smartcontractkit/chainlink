[![Go Version](https://img.shields.io/badge/go-v1.19-green.svg)](https://golang.org/dl/)
[![Lint and Test](https://github.com/gballet/go-verkle/actions/workflows/go.yml/badge.svg)](https://github.com/gballet/go-verkle/actions/workflows/go.yml)
[![DeepSource](https://deepsource.io/gh/gballet/go-verkle.svg/?label=active+issues&show_trend=true&token=OjuF5Q2HbKzpWY8LgWuffNZp)](https://deepsource.io/gh/gballet/go-verkle/?ref=repository-badge)
[![goreports](https://goreportcard.com/badge/github.com/gballet/go-verkle)](https://goreportcard.com/report/github.com/gballet/go-verkle)
[![API Reference](https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667)](https://pkg.go.dev/github.com/gballet/go-verkle)
[![Block replay](https://github.com/gballet/go-verkle/actions/workflows/block_replay.yml/badge.svg)](https://github.com/gballet/go-verkle/actions/workflows/block_replay.yml)


# go-verkle

> A Go implementation of Verkle Tree datastructure defined in the [spec](https://github.com/crate-crypto/verkle-trie-ref/tree/master/verkle). 


## Test & Benchmarks

To run the tests and benchmarks, run the following commands:
```bash
$ go test ./...
```

To run the benchmarks:
```bash
go test ./... -bench=. -run=none -benchmem
```

## Security

If you find any security vulnerability, please don't open a GH issue and contact repo owners directly.


## LICENSE

[License](LICENSE).
