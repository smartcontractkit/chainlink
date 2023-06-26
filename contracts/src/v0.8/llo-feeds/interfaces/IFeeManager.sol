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
    function setSubscriberDiscount(address subscriber, bytes32 feedId, address token, uint16 discount) external;

    /**
     * @notice Removes a subscriber from the fee manager
     * @param subscriber address of the subscriber
     * @param feedId feed id to apply the discount to
     * @param token token to apply the discount to
     */
    function removeSubscriberDiscount(address subscriber, bytes32 feedId, address token) external;

    /**
     * @notice Gets the fee from a report. If the sender is a subscriber, they will receive a discount.
	 * @param sender sender address trying to verify
	 * @param signedReport signed report to verify
	 * @param feeMetadata any metadata required to fetch the fee
	 * @return feeData fee data containing token and amount
     */
    function getFee(address sender, bytes calldata signedReport, bytes calldata feeMetadata) external returns (Common.Asset memory feeData);

    /**
     * @notice Updates the subscriber address for a discount
     * @param newSubscriberAddress new subscriber address
     * @param feedId feed id the discount is applied to
     */
    function updateSubscriberDiscountAddress(address newSubscriberAddress, bytes32 feedId) external;
}