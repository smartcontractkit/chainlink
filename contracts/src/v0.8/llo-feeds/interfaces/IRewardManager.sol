// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../shared/vendor/IERC165.sol";
import {Common} from "../../libraries/internal/Common.sol";

interface IRewardManager is IERC165 {

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
    * @param payee the user the funds should be deposited from
    * @param fee struct with the asset address and amount forwarded to the FeeManager
    */
    function onFeePaid(bytes32 poolId, address payee, Common.Asset calldata fee) external;

    /**
     * @notice Claims the rewards in a specific pool
     * @param poolIds array of poolIds to claim reward-manager for
     */
    function claimRewards(bytes32[] calldata poolIds) external;

     /**
      * @notice Updates a subset the reward recipients for a specific poolId. The collective weight of the recipients should add up to the recipients existing weights. Any recipients with a weight of 0 will be removed.
      * @param poolIds an array containing the single poolId to update
      * @param newRewardRecipients array of new reward recipients
      */
    function updateRewardRecipients(bytes32[] calldata poolIds, Common.AddressAndWeight[] calldata newRewardRecipients) external;

    /**
    * @notice Set the RewardRecipients and weights for a specific pool. This should only be called once per pool Id. Else updateRewardRecipients should be used.
    * @param poolId poolId to set RewardRecipients and weights for
    * @param rewardRecipientAndWeights array of each RewardRecipient and associated weight
    */
    function setRewardRecipients(bytes32 poolId, Common.AddressAndWeight[] calldata rewardRecipientAndWeights) external;

    /**
     * @notice Gets a list of pool ids which have reward-manager for a specific recipient.
     * @param recipient address of the recipient to get pool ids for
     */
    function getAvailableRewardPoolIds(address recipient) external view returns (bytes32[] memory);

    /**
     * @notice Updates the billing address for a specific recipient
     * @param newBillingAddress address of the recipient to update billing address for
     */
    function updateBillingAddress(address newBillingAddress) external;

    /**
     * @notice Pays all the recipients for each of the pool ids
     * @param poolIds an array containing the single poolId to pay
     * @param recipients array of recipients to pay within the pool
     */
    function payRecipients(bytes32[] calldata poolIds, address[] calldata recipients) external;

    /**
     * @notice Sets the verifier proxy. This needs to be done post construction to prevent a circular dependency.
     * @param newVerifierProxy address of the new verifier proxy
     */
    function setVerifierProxy(address newVerifierProxy) external;
}