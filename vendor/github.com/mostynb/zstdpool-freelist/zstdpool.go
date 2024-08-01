// Copyright 2020 Mostyn Bramley-Moore.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package zstdpool provides pools for github.com/klauspost/compress/zstd's
// Encoder and Decoder types, which do not leak goroutines as naive usage
// of sync.Pool might.
package zstdpool

import (
	"errors"
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
)

type decoder struct {
	d    *zstd.Decoder
	next *decoder
}

// DecoderPool implements a non-leaky pool of zstd.Decoders, since sync.Pool
// can leak goroutines when used with zstd.Decoder.
type DecoderPool struct {
	opts []zstd.DOption

	mu        sync.Mutex
	head      *decoder
	available int
}

type encoder struct {
	e    *zstd.Encoder
	next *encoder
}

// EncoderPool implements a non-leaky pool of zstd.Encoders, since sync.Pool
// can leak goroutines when used with zstd.Encoder.
type EncoderPool struct {
	opts []zstd.EOption

	mu        sync.Mutex
	head      *encoder
	available int
}

// NewDecoderPool returns a DecoderPool that pools *zstd.Decoders created
// with the specified zstd.DOptions.
func NewDecoderPool(opts ...zstd.DOption) DecoderPool {
	return DecoderPool{opts: opts}
}

// Get returns a new (or reset) *zstd.Decoder for reading from r from the
// pool. The *zstd.Decoder should not be Close()'ed before being returned to
// the pool with Put().
//
// Note that the decoder.IOReadCloser() should not be Close()'ed before
// the *zstd.Decoder is returned to the pool. Consider using GetReadCloser
// instead of Get if you want to use the decoder.IOReadCloser().
func (p *DecoderPool) Get(r io.Reader) (*zstd.Decoder, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.head == nil {
		return zstd.NewReader(r, p.opts...)
	}

	head := p.head
	p.head = head.next
	p.available--

	head.d.Reset(r)

	return head.d, nil
}

// GetReadCloser returns an io.ReadCloser that decompresses zstandard data
// from r and returns the underlying *zstd.Decoder to this DecoderPool on
// Close().
func (p *DecoderPool) GetReadCloser(r io.Reader) (*DecoderReadCloser, error) {
	dec, err := p.Get(r)
	if err != nil {
		return nil, err
	}

	return &DecoderReadCloser{p: p, d: dec, rc: dec.IOReadCloser()}, nil
}

// Put adds an unused *zstd.Decoder to the pool. You should only add decoders
// to the pool that were returned by the pool's Get function, or were created
// with same zstd.DOptions as the pool.
func (p *DecoderPool) Put(d *zstd.Decoder) {
	d.Reset(nil) // Free up reference to the underlying io.Reader.

	p.mu.Lock()
	defer p.mu.Unlock()

	p.available++

	if p.head == nil {
		p.head = &decoder{d: d}
		return
	}

	dec := &decoder{d: d, next: p.head}
	p.head = dec
	return
}

// TargetSize functions take the current size of a pool, and return the
// target size (which must not be larger than the current size), and should
// not be negative.
type TargetSize func(currentSize int) (targetSize int)

var errTargetTooLarge = errors.New("TargetSize functions should not return a value larger than the current size")

var errNegativeTarget = errors.New("TargetSize functions should not return negative values")

var errMiscount = errors.New("internal error: miscount while resizing pool")

// Resize takes a TargetSize function ts, which it asks for a target size to
// reduce the pool size to. It returns the original size of the pool, the new
// size of the pool, and an error if something went wrong.
func (p *DecoderPool) Resize(ts TargetSize) (old, new int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	old = p.available

	targetSize := ts(old)
	if targetSize == old {
		return old, old, nil
	}

	if targetSize > p.available {
		return old, old, errTargetTooLarge
	}

	if targetSize < 0 {
		return old, old, errNegativeTarget
	}

	for p.available > targetSize {
		if p.head == nil {
			return old, p.available, errMiscount
		}

		p.head.d.Reset(nil)
		p.head = p.head.next
		p.available--
	}

	return old, targetSize, nil
}

// DecoderReadCloser implements io.Readcloser by wrapping a *zstd.Decoder
// and returning it to a DecoderPool when Close() is called.
type DecoderReadCloser struct {
	p  *DecoderPool
	d  *zstd.Decoder
	rc io.ReadCloser
}

// Close does not close the underlying *zstd.Decoder, but instead returns
// it to the DecoderPool.
func (c *DecoderReadCloser) Close() {
	c.p.Put(c.d)
}

// Read wraps the *zstd.Decoder's Read function.
func (c *DecoderReadCloser) Read(b []byte) (int, error) {
	return c.d.Read(b)
}

// WriteTo wraps the *zstd.Decoder's WriteTo function.
func (c *DecoderReadCloser) WriteTo(w io.Writer) (int64, error) {
	return c.d.WriteTo(w)
}

// NewEncoderPool returns an EncoderPool that pools *zstd.Encoders created
// with the specified zstd.EOptions.
func NewEncoderPool(opts ...zstd.EOption) EncoderPool {
	return EncoderPool{opts: opts}
}

// Get returns a new (or reset) *zstd.Encoder for writing to w from the
// pool. The *zstd.Encoder should be Close()'ed before being returned to
// the pool with Put().
func (p *EncoderPool) Get(w io.Writer) (*zstd.Encoder, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.head == nil {
		return zstd.NewWriter(w, p.opts...)
	}

	head := p.head
	p.head = head.next
	p.available--

	head.e.Reset(w)

	return head.e, nil
}

// Put adds an unused *zstd.Encoder to the pool. You should only add encoders
// to the pool that were returned by the pool's Get function, or were created
// with same zstd.EOptions as the pool.
func (p *EncoderPool) Put(e *zstd.Encoder) {
	e.Reset(nil) // Free up reference to the underlying io.Writer.

	p.mu.Lock()
	defer p.mu.Unlock()

	p.available++

	if p.head == nil {
		p.head = &encoder{e: e}
		return
	}

	enc := &encoder{e: e, next: p.head}
	p.head = enc
	return
}

// Resize takes a TargetSize function ts, which it asks for a target size to
// reduce the pool size to. It returns the original size of the pool, the new
// size of the pool, and an error if something went wrong.
func (p *EncoderPool) Resize(ts TargetSize) (old, new int, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	old = p.available

	targetSize := ts(old)
	if targetSize == old {
		return old, old, nil
	}

	if targetSize > p.available {
		return old, old, errTargetTooLarge
	}

	if targetSize < 0 {
		return old, old, errNegativeTarget
	}

	for p.available > targetSize {
		if p.head == nil {
			return old, p.available, errMiscount
		}

		p.head.e.Reset(nil)
		p.head = p.head.next
		p.available--
	}

	return old, targetSize, nil
}
