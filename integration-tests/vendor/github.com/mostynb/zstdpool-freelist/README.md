# zstdpool

Zstdpool provides flexible pool implementations for the Encoder and Decoder
types in github.com/klauspost/compress/zstd which do not leak goroutines.

## Why not use sync.Pool?

`zstd.Encoder` leaks goroutines if it is garbage collected without `Close()`
being called first. So we can't safely put an unclosed encoder in a sync.Pool.
But encoders cannot be reused after being closed, so we can't put a closed
encoder in a sync.Pool either.

`zstd.Decoder`s can be reused after being closed, so you can close them before
placing them in a sync.Pool, but doing so frees resources that we would want
to keep until the decoder is no longer used.

These problems might be possible to work around with finalizers, but it is
difficult to confirm is working as expected, and could silently break
with internal changes in github.com/klauspost/compress/zstd.

## Status

This code is not yet well tested, and the API may change at any time.

## Usage

### Encoding

```Go
var encPool = zstdpool.NewEncoderPool()

func compressStream(in io.Reader, out io.Writer) error {
	enc, err := encPool.Get(out)
	if err != nil {
		return err
	}
	defer encPool.Put(enc)

	_, err = enc.ReadFrom(in)
	return err
}
```

### Decoding

```Go
var decPool = zstdpool.NewDecoderPool()

func decompressStream(in io.Reader, out io.Writer) error {
	dec, err := decPool.Get(in)
	if err != nil {
		return err
	}
	defer decPool.Put(dec)

	_, err = dec.WriteTo(out)
	return err
}
```

## Contributions

Contributions are always welcome.

## License

This code is released under the [Apache 2.0 license](LICENSE).
