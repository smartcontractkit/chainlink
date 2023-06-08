// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {IFeeManager} from "./interfaces/IFeeManager.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";

/*
 * @title FeeManager
 * @author Austin Born
 * @author Michael Fletcher
 * @notice This contract is used for the handling of fees required for users verifying reports.
 */
contract FeeManager is IFeeManager, ConfirmedOwner, TypeAndVersionInterface {

    //list of subscribers and their discounts
    mapping(address => uint16) private subscriberDiscounts;

    //the total discount that can be applied to a fee, 10000 = 100% discount
    uint16 private constant TOTAL_DISCOUNT = 10000;

    //the error thrown if the discount exceeds the maximum allowed
    error InvalidDiscount();

    struct FeeData {
        // The token address of the fee
        address token;
        // The amount of the fee
        uint256 amount;
        // The hash of the report
        bytes32 reportHash;
    }

    struct Report {
        // The feed ID the report has data for
        bytes32 feedId;
        // The time the median value was observed on
        uint32 observationsTimestamp;
        // The median value agreed in an OCR round
        int192 median;
        // The best bid value agreed in an OCR round
        int192 bid;
        // The best ask value agreed in an OCR round
        int192 ask;
        // The upper bound of the block range the median value was observed within
        uint64 blocknumberUpperBound;
        // The blockhash for the upper bound of block range (ensures correct blockchain)
        bytes32 upperBlockhash;
        // The lower bound of the block range the median value was observed within
        uint64 blocknumberLowerBound;
        // The fee for the report
        FeeData fee;
    }

    constructor() ConfirmedOwner(msg.sender) {}

    /// @inheritdoc TypeAndVersionInterface
    function typeAndVersion() external pure override returns (string memory) {
        return "FeeManager 0.0.1";
    }

    /// @inheritdoc IERC165
    function supportsInterface(bytes4 interfaceId)
    external
    pure
    override
    returns (bool)
    {
        return interfaceId == this.getFee.selector;
    }

    // @inheritdoc IFeeManager
    function setSubscriberDiscount(address subscriber, uint16 discount) external onlyOwner {
        //make sure the discount is not greater than the total discount that can be applied
        if(discount > TOTAL_DISCOUNT) revert InvalidDiscount();

        subscriberDiscounts[msg.sender] = discount;
    }

    // @inheritdoc IFeeManager
    function removeSubscriberDiscount(address subscriber) external onlyOwner {
        delete subscriberDiscounts[subscriber];
    }

    // @inheritdoc IFeeManager
    function getFee(address sender, bytes calldata signedReport, bytes calldata feeMetadata) external returns (bytes memory feeData) {
        //decode the signedReport
        Report memory report = abi.decode(signedReport, (Report));

        //set the fee amount to the discounted fee, rounding down
        report.fee.amount = report.fee.amount - (report.fee.amount * subscriberDiscounts[sender] / TOTAL_DISCOUNT);

        //return the encoded fee data
        return abi.encode(report.fee);
    }

}