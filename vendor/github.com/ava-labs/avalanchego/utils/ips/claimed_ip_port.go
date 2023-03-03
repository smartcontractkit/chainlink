// Copyright (C) 2019-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package ips

import (
	"crypto/x509"
)

// Can't import these from wrappers package due to circular import.
const (
	intLen  = 4
	longLen = 8
	ipLen   = 18
	// Certificate length, signature length, IP, timestamp
	baseIPCertDescLen = 2*intLen + ipLen + longLen
)

// A self contained proof that a peer is claiming ownership of an IPPort at a
// given time.
type ClaimedIPPort struct {
	// The peer's certificate.
	Cert *x509.Certificate
	// The peer's claimed IP and port.
	IPPort IPPort
	// The time the peer claimed to own this IP and port.
	Timestamp uint64
	// [Cert]'s signature over the IPPort and timestamp.
	// This is used in the networking library to ensure that this IPPort was
	// actually claimed by the peer in question, and not by a malicious peer
	// trying to get us to dial bogus IPPorts.
	Signature []byte
}

// Returns the length of the byte representation of this ClaimedIPPort.
func (i *ClaimedIPPort) BytesLen() int {
	// See wrappers.PackPeerTrackInfo.
	return baseIPCertDescLen + len(i.Cert.Raw) + len(i.Signature)
}
