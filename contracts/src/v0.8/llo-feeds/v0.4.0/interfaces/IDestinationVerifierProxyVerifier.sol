// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

interface IDestinationVerifierProxyVerifier is IERC165 {
  /**
   * @notice Verifies that the data encoded has been signed correctly using the signatures included within the payload.
   * @param signedReport The encoded data to be verified.
   * @param parameterPayload The encoded parameters to be used in the verification and billing process.
   * @param sender The address that requested to verify the contract.Used for logging and applying the fee.
   * @dev Verification is typically only done through the proxy contract so we can't just use msg.sender.
   * @return verifierResponse The encoded verified response.
   */
  function verify(
    bytes calldata signedReport,
    bytes calldata parameterPayload,
    address sender
  ) external payable returns (bytes memory verifierResponse);

  /**
   * @notice Bulk verifies that the data encoded has been signed correctly using the signatures included within the payload.
   * @param signedReports The encoded data to be verified.
   * @param parameterPayload The encoded parameters to be used in the verification and billing process.
   * @param sender The address that requested to verify the contract. Used for logging and applying the fee.
   * @dev Verification is typically only done through the proxy contract so we can't just use msg.sender.
   * @return verifiedReports The encoded verified responses.
   */
  function verifyBulk(
    bytes[] calldata signedReports,
    bytes calldata parameterPayload,
    address sender
  ) external payable returns (bytes[] memory verifiedReports);

  /*
   * @notice Returns the reward manager
   * @return IDestinationRewardManager
   */
  // solhint-disable-next-line func-name-mixedcase
  function s_feeManager() external view returns (address);

  /**
   * @notice Returns the access controller
   * @return IDestinationFeeManager
   */
  // solhint-disable-next-line func-name-mixedcase
  function s_accessController() external view returns (address);
}
