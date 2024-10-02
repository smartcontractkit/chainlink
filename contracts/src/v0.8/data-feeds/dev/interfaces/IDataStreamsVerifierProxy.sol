// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice IDataStreamsVerifierProxy
/// VerifierProxy of the Data Streams service. This
/// is the user gateway that exposes methods to
/// verify a set of off-chain reports that have been
/// signed by the Chainlink DON to prove authenticity
/// of the data origin. A fee in either the chain's
/// native asset or LINK will be assessed based on the
/// address encoded in parameterPayload.

interface IDataStreamsVerifierProxy {
  /**
   * @notice Return address of the FeeManager contract
   * @return feeManager address
   */
  function s_feeManager() external view returns (address feeManager);

  /**
   * @notice Bulk verifies that the data encoded has been signed
   * correctly by routing to the correct verifier, and bills the user if applicable.
   * @param payloads The encoded payloads to be verified, including the signed
   * report.
   * @param parameterPayload fee metadata for billing
   * @return verifiedReports The encoded reports from the verifier.
   */
  function verifyBulk(
    bytes[] calldata payloads,
    bytes calldata parameterPayload
  ) external payable returns (bytes[] memory verifiedReports);
}
