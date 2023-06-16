// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {IFeeManager} from "./interfaces/IFeeManager.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";
import {ByteUtil} from "../libraries/internal/ByteUtil.sol";

/*
 * @title FeeManager
 * @author Austin Born
 * @author Michael Fletcher
 * @notice This contract is used for the handling of fees required for users verifying reports.
 */
contract FeeManager is IFeeManager, ConfirmedOwner, TypeAndVersionInterface {
    //using for bytes manipulation
    using ByteUtil for bytes;

    //list of subscribers and their discounts subscriberDiscounts[subscriber][feedId][token]
    mapping(address => mapping(bytes32 => mapping(address => uint256))) private subscriberDiscounts;

    //the total discount that can be applied to a fee, 10000 = 100% discount
    uint16 private constant TOTAL_DISCOUNT = 10000;

    //the length of the fee data in bytes
    address private immutable LINK_ADDRESS;

    //the index of the link fee data in the report
    uint16 private constant LINK_FEE_INDEX = 32 + 4 + 24 +  24 + 24 + 8 + 32 + 8;

    //the index of the link fee data in the report
    uint16 private constant NATIVE_FEE_INDEX = 32 + 4 + 24 +  24 + 24 + 8 + 32 + 8 + 32;

    //the index of the  fee data in the quote
    uint16 private constant QUOTE_METADATA_FEE_ADDRESS_INDEX = 0;

    //the error thrown if the discount exceeds the maximum allowed
    error InvalidDiscount();

    //the struct to store the fee data
    struct FeeData {
        // The token address of the fee
        address token;
        // The amount of the fee
        uint256 amount;
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
        // The fee for the report in link
        uint256 linkFee;
        // The fee for the report in native
        uint256 nativeFee;
    }

    /**
     * @notice Construct the FeeManager contract
     * @param linkAddress The address of the LINK token
     */
    constructor(address linkAddress) ConfirmedOwner(msg.sender) {
        //set the link address
        LINK_ADDRESS = linkAddress;
    }

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
    function setSubscriberDiscount(address subscriber, bytes32 feedId, address token, uint16 discount) external onlyOwner {
        //make sure the discount is not greater than the total discount that can be applied
        if(discount > TOTAL_DISCOUNT) revert InvalidDiscount();

        subscriberDiscounts[subscriber][feedId][token] = discount;
    }

    // @inheritdoc IFeeManager
    function removeSubscriberDiscount(address subscriber, bytes32 feedId, address token) external onlyOwner {
        delete subscriberDiscounts[subscriber][feedId][token];
    }

    // @inheritdoc IFeeManager
    function getFee(address sender, bytes calldata signedReport, bytes calldata quoteMetadata) external view returns (bytes memory feeData) {
        //The quote
        FeeData memory fee;

        //any report without a fee will default to 0
        if (signedReport.length < LINK_FEE_INDEX) {
            return abi.encode(fee);
        }

        //decode the quoteMetadata to get the desired fee
        address quoteFeeAddress = quoteMetadata.readAddress(QUOTE_METADATA_FEE_ADDRESS_INDEX);

        //calculate either the LINK fee or native fee if it's within the report
        if(quoteFeeAddress == LINK_ADDRESS) {
            fee.token = LINK_ADDRESS;
            fee.amount = signedReport.readUint256(LINK_FEE_INDEX);
        } else {
            fee.token = address(0);
            fee.amount = signedReport.readUint256(NATIVE_FEE_INDEX);
        }

        //decode the feedId from the report to calculate the discount being applied
        bytes32 feedId = bytes32(signedReport);

        //set the fee amount to the discounted fee, rounding down
        fee.amount = fee.amount - (fee.amount * subscriberDiscounts[sender][feedId][quoteFeeAddress] / TOTAL_DISCOUNT);

        //return the encoded fee data
        return abi.encode(fee);
    }
}