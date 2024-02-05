// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Common} from "../libraries/Common.sol";
import {AccessControllerInterface} from "../../shared/interfaces/AccessControllerInterface.sol";
import {IVerifierFeeManager} from "./IVerifierFeeManager.sol";

interface IVerifierProxy {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the correct verifier, and bills the user if applicable.
   * @param payload The encoded data to be verified, including the signed
   * report.
   * @param parameterPayload fee metadata for billing
   * @return verifierResponse The encoded report from the verifier.
   */
  function verify(
    bytes calldata payload,
    bytes calldata parameterPayload
  ) external payable returns (bytes memory verifierResponse);

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

  /**
   * @notice Sets the verifier address initially, allowing `setVerifier` to be set by this Verifier in the future
   * @param verifierAddress The address of the verifier contract to initialize
   */
  function initializeVerifier(address verifierAddress) external;

  /**
   * @notice Sets a new verifier for a config digest
   * @param currentConfigDigest The current config digest
   * @param newConfigDigest The config digest to set
   * @param addressesAndWeights The addresses and weights of reward recipients
   * reports for a given config digest.
   */
  function setVerifier(
    bytes32 currentConfigDigest,
    bytes32 newConfigDigest,
    Common.AddressAndWeight[] memory addressesAndWeights
  ) external;

  /**
   * @notice Removes a verifier for a given config digest
   * @param configDigest The config digest of the verifier to remove
   */
  function unsetVerifier(bytes32 configDigest) external;

  /**
   * @notice Retrieves the verifier address that verifies reports
   * for a config digest.
   * @param configDigest The config digest to query for
   * @return verifierAddress The address of the verifier contract that verifies
   * reports for a given config digest.
   */
  function getVerifier(bytes32 configDigest) external view returns (address verifierAddress);

  /**
   * @notice Called by the admin to set an access controller contract
   * @param accessController The new access controller to set
   */
  function setAccessController(AccessControllerInterface accessController) external;

  /**
   * @notice Updates the fee manager
   * @param feeManager The new fee manager
   */
  function setFeeManager(IVerifierFeeManager feeManager) external;
}
