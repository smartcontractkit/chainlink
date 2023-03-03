// Package bloomfilter is face-meltingly fast, thread-safe,
// marshalable, unionable, probability- and
// optimal-size-calculating Bloom filter in go
//
// https://github.com/steakknife/bloomfilter
//
// Copyright © 2014, 2015, 2018 Barry Allard
//
// MIT license
//
package v2

import (
	"compress/gzip"
	_ "encoding/gob" // make sure gob is available
	"encoding/json"
	"errors"
	"io"
	"os"
)

// ReadFrom r and overwrite f with new Bloom filter data
func (f *Filter) ReadFrom(r io.Reader) (n int64, err error) {
	f2, n, err := ReadFrom(r)
	if err != nil {
		return -1, err
	}
	f.lock.Lock()
	defer f.lock.Unlock()
	f.m = f2.m
	f.n = f2.n
	f.bits = f2.bits
	f.keys = f2.keys
	return n, nil
}

// ReadFrom Reader r into a lossless-compressed Bloom filter f
func ReadFrom(r io.Reader) (f *Filter, n int64, err error) {
	f = new(Filter)
	rawR, err := gzip.NewReader(r)
	if err != nil {
		return nil, -1, err
	}
	defer rawR.Close()
	n, err = f.UnmarshalFromReader(rawR)
	if err != nil {
		return nil, -1, err
	}
	return f, n, nil
}

// ReadFile from filename into a lossless-compressed Bloom Filter f
// Suggested file extension: .bf.gz
func ReadFile(filename string) (f *Filter, n int64, err error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, -1, err
	}
	defer r.Close()

	return ReadFrom(r)
}

// WriteTo a Writer w from lossless-compressed Bloom Filter f
func (f *Filter) WriteTo(w io.Writer) (n int64, err error) {
	f.lock.RLock()
	defer f.lock.RUnlock()

	rawW := gzip.NewWriter(w)
	defer rawW.Close()

	intN, _, err := f.MarshallToWriter(rawW)
	//intN, _, err := f.MarshallToWriter(w)
	n = int64(intN)
	return n, err
}

// WriteFile filename from a a lossless-compressed Bloom Filter f
// Suggested file extension: .bf.gz
func (f *Filter) WriteFile(filename string) (n int64, err error) {
	w, err := os.Create(filename)
	if err != nil {
		return -1, err
	}
	defer w.Close()

	return f.WriteTo(w)
}

type jsonType struct {
	Version string   `json:"version"`
	Bits    []uint64 `json:"bits"`
	Keys    []uint64 `json:"keys"`
	M       uint64   `json:"m"`
	N       uint64   `json:"n"`
}

func (f *Filter) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonType{
		string(version),
		f.bits,
		f.keys,
		f.m,
		f.n,
	})
}

func (f *Filter) UnmarshalJSON(data []byte) error {
	var j jsonType
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if j.Version != string(version) {
		return errors.New("incompatible version")
	}
	f.bits = j.Bits
	f.keys = j.Keys
	f.n = j.N
	f.m = j.M
	return nil
}

// GobDecode conforms to interface gob.GobDecoder
func (f *Filter) GobDecode(data []byte) error {
	return f.UnmarshalBinary(data)
}

// GobEncode conforms to interface gob.GobEncoder
func (f *Filter) GobEncode() ([]byte, error) {
	return f.MarshalBinary()
}
