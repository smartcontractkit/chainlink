// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {IERC165} from "../../shared/vendor/IERC165.sol";

interface IRewardBank is IERC165 {

    /**
    * @notice Record the fee received for a `verify` request
    * @param configDigest config digest of the report being verified
    * @param fee struct with the asset address and amount forwarded to the FeeManager
    */
    function onFeePaid(bytes32 configDigest, Asset calldata fee) external;

    /**
    * @notice Distributes all the rewards in the specified pot to the Payees
    * @param configDigest config digest of the report being verified
    * @param assetAddress address of the assetAddress to distribute
    */
    function payPayees(bytes32 configDigest, address assetAddress) external;

    /**
    * @notice Set the Payees and weights for a specific feed config
    * @param configDigest config digest to set Payees and weights for
    * @param PayeeAndWeights array of each Payee and associated weight
    */
    function setPayees(bytes32 configDigest, PayeeAndWeight[] calldata PayeeAndWeights) external;

    // @notice The asset struct to hold the address of the asset and the amount
    struct Asset {
        address assetAddress;
        uint256 amount;
    }

    // @notice Struct to hold the address of a Payee and its weight to determine what percentage of the pot it receives
    struct PayeeAndWeight {
        address PayeeAddress;
        uint8 weight;
    }
}