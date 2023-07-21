// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../shared/vendor/IERC165.sol";
import {Common} from "../../libraries/internal/Common.sol";

interface IFeeManager is IERC165 {
  /**
   * @notice Adds a subscriber to the fee manager
   * @param subscriber address of the subscriber
   * @param discount discount to be applied to the fee
   * @param feedId feed id to apply the discount to
   * @param token token to apply the discount to
   */
  function updateSubscriberDiscount(address subscriber, bytes32 feedId, address token, uint256 discount) external;

  /**
   * @notice Processes the fee for a report, billing the subscriber and paying the reward manager
   * @param payload report and quote to process the fee for
   * @param subscriber address of the fee will be applied
   */
  function processFee(bytes calldata payload, address subscriber) external payable;

  /**
   * @notice Sets the native premium
   * @param premium premium to be paid if paying in native
   */
  function setNativePremium(uint256 premium) external;

  /**
   * @notice Sets the fee recipients within the reward manager
   * @param configDigest digest of the configuration
   * @param rewardRecipientAndWeights the address and weights of all the recipients to receive rewards
   */
  function setFeeRecipients(
    bytes32 configDigest,
    Common.AddressAndWeight[] calldata rewardRecipientAndWeights
  ) external;

  /**
   * @notice Withdraws any native rewards to the owner address
   * @param quantity quantity of native tokens to withdraw, address(0) is native
   * @param quantity quantity to withdraw
   */
  function withdraw(address assetAddress, uint256 quantity) external;
}
