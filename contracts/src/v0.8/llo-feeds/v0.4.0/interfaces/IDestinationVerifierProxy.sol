// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

interface IDestinationVerifierProxy is IERC165 {
  /**
   * @notice Verifies that the data encoded has been signed
   * correctly by routing to the verifier, and bills the user if applicable.
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
   * @notice Sets the active verifier for this proxy
   * @param verifierAddress The address of the verifier contract
   */
  function setVerifier(address verifierAddress) external;

  /**
   * @notice Used to honor the source verifierProxy feeManager interface
   * @return IVerifierFeeManager
   */
  // solhint-disable-next-line func-name-mixedcase
  function s_feeManager() external view returns (address);

  /**
   * @notice Used to honor the source verifierProxy feeManager interface
   * @return AccessControllerInterface
   */
  // solhint-disable-next-line func-name-mixedcase
  function s_accessController() external view returns (address);
}
