package multistream

import (
	"bytes"
	"errors"
	"io"
)

// ErrNotSupported is the error returned when the muxer does not support
// the protocol specified for the handshake.
var ErrNotSupported = errors.New("protocol not supported")

// ErrNoProtocols is the error returned when the no protocols have been
// specified.
var ErrNoProtocols = errors.New("no protocols specified")

// SelectProtoOrFail performs the initial multistream handshake
// to inform the muxer of the protocol that will be used to communicate
// on this ReadWriteCloser. It returns an error if, for example,
// the muxer does not know how to handle this protocol.
func SelectProtoOrFail(proto string, rwc io.ReadWriteCloser) error {
	errCh := make(chan error, 1)
	go func() {
		var buf bytes.Buffer
		delimWrite(&buf, []byte(ProtocolID))
		delimWrite(&buf, []byte(proto))
		_, err := io.Copy(rwc, &buf)
		errCh <- err
	}()
	// We have to read *both* errors.
	err1 := readMultistreamHeader(rwc)
	err2 := readProto(proto, rwc)
	if werr := <-errCh; werr != nil {
		return werr
	}
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// SelectOneOf will perform handshakes with the protocols on the given slice
// until it finds one which is supported by the muxer.
func SelectOneOf(protos []string, rwc io.ReadWriteCloser) (string, error) {
	if len(protos) == 0 {
		return "", ErrNoProtocols
	}

	// Use SelectProtoOrFail to pipeline the /multistream/1.0.0 handshake
	// with an attempt to negotiate the first protocol. If that fails, we
	// can continue negotiating the rest of the protocols normally.
	//
	// This saves us a round trip.
	switch err := SelectProtoOrFail(protos[0], rwc); err {
	case nil:
		return protos[0], nil
	case ErrNotSupported: // try others
	default:
		return "", err
	}
	for _, p := range protos[1:] {
		err := trySelect(p, rwc)
		switch err {
		case nil:
			return p, nil
		case ErrNotSupported:
		default:
			return "", err
		}
	}
	return "", ErrNotSupported
}

func handshake(rw io.ReadWriter) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- delimWriteBuffered(rw, []byte(ProtocolID))
	}()

	if err := readMultistreamHeader(rw); err != nil {
		return err
	}
	return <-errCh
}

func readMultistreamHeader(r io.Reader) error {
	tok, err := ReadNextToken(r)
	if err != nil {
		return err
	}

	if tok != ProtocolID {
		return errors.New("received mismatch in protocol id")
	}
	return nil
}

func trySelect(proto string, rwc io.ReadWriteCloser) error {
	err := delimWriteBuffered(rwc, []byte(proto))
	if err != nil {
		return err
	}
	return readProto(proto, rwc)
}

func readProto(proto string, r io.Reader) error {
	tok, err := ReadNextToken(r)
	if err != nil {
		return err
	}

	switch tok {
	case proto:
		return nil
	case "na":
		return ErrNotSupported
	default:
		return errors.New("unrecognized response: " + tok)
	}
}
