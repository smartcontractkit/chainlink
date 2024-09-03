// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {Common} from "../../libraries/Common.sol";
import {IERC165} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";

interface IDestinationVerifier is IERC165 {
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
}
