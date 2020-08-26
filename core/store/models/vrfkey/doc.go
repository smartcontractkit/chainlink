// Package vrfkey tracks the secret keys associated with VRF proofs. It
// is a separate package from ../store to increase encapsulation of the keys,
// reduce the risk of them leaking, and reduce confusion between VRF keys and
// ethereum keys.
//
// The three types, PrivateKey (private_key.go), PublicKey (public_key.go) and
// EncryptedVRFKey (serialzation.go) are all aspects of the one keypair.
//
// The details of the secret key in a keypair should remain private to this
// package. If you need to access the secret key, you should add a method to
// PrivateKey which does the crypto requiring it, without leaking the secret.
// See MakeMarshaledProof in private_key.go, for an example.
//
// PrivateKey#PublicKey represents the associated public key, and, in the
// context of a VRF, represents a public commitment to a particular "verifiable
// random function."
//
// EncryptedVRFKey is used to store a public/private keypair in a database,
// in encrypted form.
//
// Usage
//
// Call vrfkey.CreateKey() to generate a PrivateKey with cryptographically
// secure randomness.
//
// Call PrivateKey#Encrypt(passphrase) to create a representation of the key
// which is encrypted for storage.
//
// Call MakeMarshaledProof(privateKey, seed) to generate the VRF proof for the given
// seed and private key. The proof is marshaled in the format expected by the
// on-chain verification mechanism in VRF.sol. If you want to know the VRF
// output independently of the on-chain verification mechanism, you can get it
// from vrf.UnmarshalSolidityProof(p).Output.
package vrfkey
