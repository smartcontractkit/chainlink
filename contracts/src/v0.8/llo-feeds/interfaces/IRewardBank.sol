// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../shared/vendor/IERC165.sol";

interface IRewardBank is IERC165 {

    /**
     * @notice Add a contract address to the authorized contracts list
     * @param contractAddress address of the contract to add
     */
    function addAuthorizedContract(address contractAddress) external;

    /**
      * @notice Remove a contract address from the authorized contracts list
      * @param contractAddress address of the contract to remove
      */
    function removeAuthorizedContract(address contractAddress) external;

    /**
    * @notice Record the fee received for a particular pool
    * @param poolId poolId of the report being verified
    * @param fee struct with the asset address and amount forwarded to the FeeManager
    */
    function onFeePaid(bytes32 poolId, Asset calldata fee) external;

    /**
     * @notice Claims the rewards for a specific asset in the pot for a specific pool id
     * @param poolIds array of poolIds to claim rewards for
     */
    function claimRewards(bytes32[] calldata poolIds) external;

     /**
      * @notice Updates a subset the reward recipients for a specific poolId. The collective weight of the recipients should add up to the recipients existing weights. Any recipients with a weight of 0 will be removed.
      * @param poolIds an array containing the single poolId to update
      * @param newRewardRecipients array of new reward recipients
      */
    function updateRewardRecipients(bytes32[] calldata poolIds, RewardRecipientAndWeight[] calldata newRewardRecipients) external;

    /**
    * @notice Set the RewardRecipients and weights for a specific pool. This should only be called once per pool Id. Else updateRewardRecipients should be used.
    * @param poolId poolId to set RewardRecipients and weights for
    * @param rewardRecipientAndWeights array of each RewardRecipient and associated weight
    */
    function setRewardRecipients(bytes32 poolId, RewardRecipientAndWeight[] calldata rewardRecipientAndWeights) external;

    /**
     * @notice Gets a list of pool ids which have rewards for a specific recipient.
     * @param Recipient address of the recipient to get pool ids for
     */
    function getAvailableRewardPoolIds(address Recipient) external view returns (bytes32[] memory);

    // @notice The asset struct to hold the address of the fee and the amount
    struct Asset {
        address assetAddress;
        uint256 amount;
    }

    // @notice Struct to hold the address of a reward recipient and their weight to determine what percentage of the pot they are entitled too
    struct RewardRecipientAndWeight {
        address rewardRecipientAddress;
        uint16 weight;
    }

}