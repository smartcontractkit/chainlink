// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

interface IRewardBank {

    /**
    * @notice Record the fee received for a `verify` request
    * @param configDigest config digest of the report being verified
    * @param fee struct with the asset address and amount forwarded to the FeeManager
    */
    function onFeePaid(bytes32 configDigest, Asset calldata fee) external;

    /**
    * @notice Distributes all the rewards in the specified pot to the NOPs
    * @param configDigest config digest of the report being verified
    * @param assetAddress address of the assetAddress to distribute
    */
    function payNOPs(bytes32 configDigest, address assetAddress) external;

    /**
    * @notice Set the NOPs and weights for a specific feed config
    * @param configDigest config digest to set NOPs and weights for
    * @param NOPAndWeights array of each NOP and associated weight
    */
    function setNOPs(bytes32 configDigest, NOPAndWeight[] calldata NOPAndWeights) external;

    // @notice The asset struct to hold the address of the asset and the amount
    struct Asset {
        address assetAddress;
        uint256 amount;
    }

    // @notice Struct to hold the address of a NOP and its weight to determine what percentage of the pot it receives
    struct NOPAndWeight {
        address NOPAddress;
        uint8 weight;
    }
}