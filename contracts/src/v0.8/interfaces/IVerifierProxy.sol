// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

interface IVerifierProxy {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the correct verifier.
   * @param signedReport The encoded data to be verified.
   * @return verifierResponse The encoded response from the verifier.
   */
  function verify(bytes memory signedReport) external returns (bytes memory verifierResponse);

  /**
   * @notice Sets a new verifier for a config digest
   * @param currentConfigDigest The current config digest
   * @param newConfigDigest The config digest to set
   * reports for a given config digest.
   */
  function setVerifier(bytes32 currentConfigDigest, bytes32 newConfigDigest) external;

  /**
   * @notice Sets the verifier address to initialized
   * @param verifierAddr The address of the verifier contract that we want to initialize
   */
  function initializeVerifier(address verifierAddr) external;

  /**
   * @notice Removes a verifier
   * @param configDigest The config digest of the verifier to remove
   */
  function unsetVerifier(bytes32 configDigest) external;

  /**
   * @notice Retrieves the verifier address that verifies reports
   * for a config digest.
   * @param configDigest The config digest to query for
   * @return verifierAddr The address of the verifier contract that verifies
   * reports for a given config digest.
   */
  function getVerifier(bytes32 configDigest) external view returns (address verifierAddr);
}
