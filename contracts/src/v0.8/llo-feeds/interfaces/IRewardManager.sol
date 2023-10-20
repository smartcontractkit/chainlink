// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.0/contracts/interfaces/IERC165.sol";
import {Common} from "../../libraries/Common.sol";

interface IRewardManager is IERC165 {
  /**
   * @notice Record the fee received for a particular pool
   * @param payments array of structs containing pool id and amount
   * @param payee the user the funds should be retrieved from
   */
  function onFeePaid(FeePayment[] calldata payments, address payee) external;

  /**
   * @notice Claims the rewards in a specific pool
   * @param poolIds array of poolIds to claim rewards for
   */
  function claimRewards(bytes32[] calldata poolIds) external;

  /**
   * @notice Set the RewardRecipients and weights for a specific pool. This should only be called once per pool Id. Else updateRewardRecipients should be used.
   * @param poolId poolId to set RewardRecipients and weights for
   * @param rewardRecipientAndWeights array of each RewardRecipient and associated weight
   */
  function setRewardRecipients(bytes32 poolId, Common.AddressAndWeight[] calldata rewardRecipientAndWeights) external;

  /**
   * @notice Updates a subset the reward recipients for a specific poolId. The collective weight of the recipients should add up to the recipients existing weights. Any recipients with a weight of 0 will be removed.
   * @param poolId the poolId to update
   * @param newRewardRecipients array of new reward recipients
   */
  function updateRewardRecipients(bytes32 poolId, Common.AddressAndWeight[] calldata newRewardRecipients) external;

  /**
   * @notice Pays all the recipients for each of the pool ids
   * @param poolId the pool id to pay recipients for
   * @param recipients array of recipients to pay within the pool
   */
  function payRecipients(bytes32 poolId, address[] calldata recipients) external;

  /**
   * @notice Sets the fee manager. This needs to be done post construction to prevent a circular dependency.
   * @param newFeeManager address of the new verifier proxy
   */
  function setFeeManager(address newFeeManager) external;

  /**
   * @notice Gets a list of pool ids which have reward for a specific recipient.
   * @param recipient address of the recipient to get pool ids for
   * @param startIndex the index to start from
   * @param endIndex the index to stop at
   */
  function getAvailableRewardPoolIds(
    address recipient,
    uint256 startIndex,
    uint256 endIndex
  ) external view returns (bytes32[] memory);

  /**
   * @notice The structure to hold a fee payment notice
   * @param poolId the poolId receiving the payment
   * @param amount the amount being paid
   */
  struct FeePayment {
    bytes32 poolId;
    uint192 amount;
  }
}
