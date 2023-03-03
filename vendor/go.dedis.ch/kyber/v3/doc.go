/*
Package kyber provides a toolbox of advanced cryptographic primitives,
for applications that need more than straightforward signing and encryption.
This top level package defines the interfaces to cryptographic primitives
designed to be independent of specific cryptographic algorithms,
to facilitate upgrading applications to new cryptographic algorithms
or switching to alternative algorithms for experimentation purposes.

Abstract Groups

This toolkits public-key crypto API includes a kyber.Group interface
supporting a broad class of group-based public-key primitives
including DSA-style integer residue groups and elliptic curve groups. Users of
this API can write higher-level crypto algorithms such as zero-knowledge
proofs without knowing or caring exactly what kind of group, let alone which
precise security parameters or elliptic curves, are being used. The kyber.Group
interface supports the standard algebraic operations on group elements and
scalars that nontrivial public-key algorithms tend to rely on. The interface
uses additive group terminology typical for elliptic curves, such that point
addition is homomorphically equivalent to adding their (potentially secret)
scalar multipliers. But the API and its operations apply equally well to
DSA-style integer groups.

As a trivial example, generating a public/private keypair is as simple as:

    suite := suites.MustFind("Ed25519")		// Use the edwards25519-curve
	a := suite.Scalar().Pick(suite.RandomStream()) // Alice's private key
	A := suite.Point().Mul(a, nil)          // Alice's public key

The first statement picks a private key (Scalar) from a the suites's source of
cryptographic random or pseudo-random bits, while the second performs elliptic
curve scalar multiplication of the curve's standard base point (indicated by the
'nil' argument to Mul) by the scalar private key 'a'. Similarly, computing a
Diffie-Hellman shared secret using Alice's private key 'a' and Bob's public key
'B' can be done via:

	S := suite.Point().Mul(a, B)		// Shared Diffie-Hellman secret

Note that we use 'Mul' rather than 'Exp' here because the library uses
the additive-group terminology common for elliptic curve crypto,
rather than the multiplicative-group terminology of traditional
integer groups - but the two are semantically equivalent and the
interface itself works for both elliptic curve and integer groups.

Higher-level Building Blocks

Various sub-packages provide several specific
implementations of these cryptographic interfaces.
In particular, the 'group/mod' sub-package provides implementations
of modular integer groups underlying conventional DSA-style algorithms.
The `group/nist` package provides NIST-standardized elliptic curves built on
the Go crypto library.
The 'group/edwards25519' sub-package provides the kyber.Group interface
using the popular Ed25519 curve.

Other sub-packages build more interesting high-level cryptographic tools
atop these primitive interfaces, including:

- share: Polynomial commitment and verifiable Shamir secret splitting
for implementing verifiable 't-of-n' threshold cryptographic schemes.
This can be used to encrypt a message so that any 2 out of 3 receivers
must work together to decrypt it, for example.

- proof: An implementation of the general Camenisch/Stadler framework
for discrete logarithm knowledge proofs.
This system supports both interactive and non-interactive proofs
of a wide variety of statements such as,
"I know the secret x associated with public key X
or I know the secret y associated with public key Y",
without revealing anything about either secret
or even which branch of the "or" clause is true.

- sign: The sign directory contains different signature schemes.

- sign/anon provides anonymous and pseudonymous public-key encryption and signing,
where the sender of a signed message or the receiver of an encrypted message
is defined as an explicit anonymity set containing several public keys
rather than just one. For example, a member of an organization's board of trustees
might prove to be a member of the board without revealing which member she is.

- sign/cosi provides collective signature algorithm, where a bunch of signers create a
unique, compact and efficiently verifiable signature using the Schnorr signature as a basis.

- sign/eddsa provides a kyber-native implementation of the EdDSA signature scheme.

- sign/schnorr provides a basic vanilla Schnorr signature scheme implementation.

- shuffle: Verifiable cryptographic shuffles of ElGamal ciphertexts,
which can be used to implement (for example) voting or auction schemes
that keep the sources of individual votes or bids private
without anyone having to trust more than one of the shuffler(s) to shuffle
votes/bids honestly.

Target Use-cases

As should be obvious, this library is intended to be used by
developers who are at least moderately knowledgeable about
cryptography. If you want a crypto library that makes it easy to
implement "basic crypto" functionality correctly - i.e., plain
public-key encryption and signing - then
[NaCl secretbox](https://godoc.org/golang.org/x/crypto/nacl/secretbox)
may be a better choice. This toolkit's purpose is to make it possible
- and preferably easy - to do slightly more interesting things that
most current crypto libraries don't support effectively. The one
existing crypto library that this toolkit is probably most comparable
to is the Charm rapid prototyping library for Python
(https://charm-crypto.com/category/charm).

This library incorporates and/or builds on existing code from a variety of
sources, as documented in the relevant sub-packages.

Reporting Security Problems

This library is offered as-is, and without a guarantee. It will need an
independent security review before it should be considered ready for use in
security-critical applications. If you integrate Kyber into your application it
is YOUR RESPONSIBILITY to arrange for that audit.

If you notice a possible security problem, please report it
to dedis-security@epfl.ch.

*/
package kyber
