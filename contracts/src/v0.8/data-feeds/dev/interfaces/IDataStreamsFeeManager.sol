// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

/// @notice IDataStreamsFeeManager
/// FeeManager of the Data Streams service which handles fees
/// to the Data Streams verifier. getFeeAndReward allows the
/// user to pass in their address, a report, and the intended
/// ERC-20 token address to get the final quoted fee amount in
/// return.

interface IDataStreamsFeeManager {
  struct Asset {
    address assetAddress;
    uint256 amount;
  }

  struct Quote {
    address quoteAddress;
  }

  /**
   * @notice Return address of the RewardManager contract
   * @return rewardManager address
   */
  function i_rewardManager() external view returns (address rewardManager);

  /**
   * @notice Calculate the applied fee and the reward from a report. If the sender is a subscriber, they will receive a discount.
   * @param subscriber address trying to verify
   * @param report report to calculate the fee for
   * @param quoteAddress address of the quote payment token
   * @return (fee, reward, totalDiscount) fee and the reward data with the discount applied
   */
  function getFeeAndReward(
    address subscriber,
    bytes memory report,
    address quoteAddress
  ) external returns (Asset memory, Asset memory, uint256);

  struct FeeAndReward {
    bytes32 configDigest;
    Asset fee;
    Asset reward;
    uint256 appliedDiscount;
  }
}
