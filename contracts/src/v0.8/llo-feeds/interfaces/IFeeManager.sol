// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {IERC165} from "../../shared/vendor/IERC165.sol";

interface IFeeManager is IERC165 {

    /**
     * @notice Adds a subscriber to the fee manager
     * @param subscriber address of the subscriber
     * @param discount discount to be applied to the fee
     */
    function setSubscriberDiscount(address subscriber, uint16 discount) external;

    /**
     * @notice Removes a subscriber from the fee manager
     * @param subscriber address of the subscriber
     */
    function removeSubscriberDiscount(address subscriber) external;

    /**
     * @notice Gets the fee from a report. If the sender is a subscriber, they will receive a discount.
	 * @param sender sender address trying to verify
	 * @param signedReport signed report to verify
	 * @param feeMetadata any metadata required to fetch the fee
	 * @return feeData fee data containing token and amount
     */
    function getFee(address sender, bytes calldata signedReport, bytes calldata feeMetadata) external returns (bytes memory feeData);

    // @notice The asset struct to hold the address of the asset and the amount
    struct Asset {
        address assetAddress;
        uint256 amount;
    }
}