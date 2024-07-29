// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Common} from "../../libraries/Common.sol";

interface IDestinationVerifier {
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

  /**
   * @notice sets off-chain reporting protocol configuration incl. participating oracles
   * @param signers addresses with which oracles sign the reports
   * @param f number of faulty oracles the system can tolerate
   * @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
   */
  function setConfig(
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights
  ) external;

  /**
   * @notice sets off-chain reporting protocol configuration incl. participating oracles
   * @param signers addresses with which oracles sign the reports
   * @param f number of faulty oracles the system can tolerate
   * @param recipientAddressesAndWeights the addresses and weights of all the recipients to receive rewards
   * @param activationTime the time at which the config was activated
   */
  function setConfigWithActivationTime(
    address[] memory signers,
    uint8 f,
    Common.AddressAndWeight[] memory recipientAddressesAndWeights,
    uint32 activationTime
  ) external;

  /**
   * @notice Sets the fee manager address
   * @param feeManager The address of the fee manager
   */
  function setFeeManager(address feeManager) external;

  /**
   * @notice Sets the access controller address
   * @param accessController The address of the access controller
   */
  function setAccessController(address accessController) external;

  /**
   * @notice Updates the config active status
   * @param donConfigId The ID of the config to update
   * @param isActive The new config active status
   */
  function setConfigActive(uint256 donConfigId, bool isActive) external;

  /**
   * @notice Removes the latest config
   */
  function removeLatestConfig() external;

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
