// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @title Abstract base contract for verification of OCR3 report
/// @dev Defines the main interface functions to be implemented by signature specific contract variants.
/// @dev Defines the common method used to compute a Keccak-256 digest over a report and its context information.
abstract contract MultiFormatAttestationVerifierBase {
    /// @notice Verifies and stores the list of provided public keys for later use within `verifyAttestation(...)`.
    ///         Reverts if an invalid number of keys was provided or any key is found invalid.
    /// @param n The number of keys expected to be set by this call. Must match the actual number of keys present in
    ///        the `keys` parameter. The maximum number of keys supported is 32 (based on the width of the
    ///        attribution bitmask).
    /// @param keys A concatenation of `n` public keys. The exact format of the key depends on the signature scheme used
    ///        for verification (for example, for ECDSA, addresses of 20 bytes each would be used).
    function _setVerificationKeys(uint8 n, bytes calldata keys) internal virtual;

    /// @notice Verifies the attestation for the given report.
    ///         Reverts if the attestation could not be verified successfully.
    /// @param configDigest configuration digest of this configuration
    /// @param n The total number of oracles. Used to verify the attribution bitmask as part of the attestation data.
    /// @param f Maximum number of faulty/dishonest oracles the protocol can tolerate while still working correctly.
    ///        An attestation generated from exactly `f + 1` signatures from different oracles is expected.
    /// @param report report data to be attested
    /// @param attestation the attestation data for the reports
    function _verifyAttestation(
        bytes32 configDigest,
        uint64 seqNr,
        uint8 n,
        uint8 f,
        bytes memory report,
        bytes calldata attestation
    ) internal virtual;

    /// @notice Compute the cryptographic hash over the given report in combination with the configuration digest and
    ///         sequence number.
    /// @param configDigest configuration digest of this configuration
    /// @param seqNr the sequence number the given report was attested in
    /// @param report the report data
    function _hashReport(bytes32 configDigest, uint64 seqNr, bytes memory report) internal pure returns (bytes32) {
        return keccak256(abi.encode(configDigest, seqNr, keccak256(report)));
    }
}