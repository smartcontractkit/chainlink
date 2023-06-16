// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../shared/vendor/IERC165.sol";

interface IRewardBank is IERC165 {

    /**
     * @notice Updates the verifier proxy address
     * @param _verifierProxyAddr address of the new verifier proxy
     */
    function setVerifierProxy(address _verifierProxyAddr) external;

    /**
    * @notice Record the fee received for a `verify` request
    * @param configDigest config digest of the report being verified
    * @param fee struct with the asset address and amount forwarded to the FeeManager
    */
    function onFeePaid(bytes32 configDigest, Asset calldata fee) external;

    /**
    * @notice Distributes all the rewards in the specified pot to the RewardRecipients
    * @param configDigest config digest of the report being verified
    * @param assetAddress address of the assetAddress to distribute
    */
    function payRewardRecipients(bytes32 configDigest, address assetAddress) external;

    /**
    * @notice Set the RewardRecipients and weights for a specific feed config
    * @param configDigest config digest to set RewardRecipients and weights for
    * @param rewardRecipientAndWeights array of each RewardRecipient and associated weight
    */
    function setRewardRecipients(bytes32 configDigest, RewardRecipientAndWeight[] calldata rewardRecipientAndWeights) external;

    // @notice The asset struct to hold the address of the asset and the amount
    struct Asset {
        address assetAddress;
        uint256 amount;
    }

    // @notice Struct to hold the address of a RewardRecipient and its weight to determine what percentage of the pot it receives
    struct RewardRecipientAndWeight {
        address rewardRecipientAddress;
        uint8 weight;
    }
}