// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {IFeeManager} from "./interfaces/IFeeManager.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";
import {ByteUtil} from "../libraries/internal/ByteUtil.sol";
import {Common} from "../libraries/internal/Common.sol";

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

  //the link token address
  address private immutable LINK_ADDRESS;

  //the native token address
  address private immutable NATIVE_ADDRESS;

  //the index of the link fee data in the report
  uint16 private constant LINK_FEE_INDEX = 32 + 4 + 24 + 24 + 24 + 8 + 32 + 8;

  //the index of the link fee data in the report
  uint16 private constant NATIVE_FEE_INDEX = 32 + 4 + 24 + 24 + 24 + 8 + 32 + 8 + 32;

  //the index of the fee data in the quote
  uint16 private constant QUOTE_METADATA_FEE_ADDRESS_INDEX = 0;

  //the premium fee to be paid if paying in native
  uint16 private nativePremium;

  //the error thrown if the discount or premium is invalid
  error InvalidPremium();

  //the error thrown if the token is invalid
  error InvalidToken();

  //the error thrown if the discount is invalid
  error InvalidDiscount();

  //the error thrown if the address is invalid
  error InvalidAddress();

  /**
   * @notice Construct the FeeManager contract
   * @param linkAddress The address of the LINK token
   */
  constructor(address linkAddress, address nativeAddress) ConfirmedOwner(msg.sender) {
    //set the link address
    LINK_ADDRESS = linkAddress;
    //set the native address
    NATIVE_ADDRESS = nativeAddress;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "FeeManager 0.0.1";
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return interfaceId == this.getFee.selector;
  }

  // @inheritdoc IFeeManager
  function setSubscriberDiscount(
    address subscriber,
    bytes32 feedId,
    address token,
    uint16 discount
  ) external onlyOwner {
    //make sure the discount is not greater than the total discount that can be applied
    if (discount > TOTAL_DISCOUNT) revert InvalidDiscount();
    //make sure the token is either LINK or native
    if (token != LINK_ADDRESS && token != NATIVE_ADDRESS) revert InvalidAddress();

    subscriberDiscounts[subscriber][feedId][token] = discount;
  }

  // @inheritdoc IFeeManager
  function removeSubscriberDiscount(address subscriber, bytes32 feedId, address token) external onlyOwner {
    delete subscriberDiscounts[subscriber][feedId][token];
  }

  // Error message when an offset is out of bounds
  error InvalidOffset(uint256 expected, uint256 actual);

  // @inheritdoc IFeeManager
  function getFee(
    address subscriber,
    bytes calldata report,
    bytes calldata quoteMetadata
  ) external view returns (Common.Asset memory asset) {
    //The quote
    Common.Asset memory fee;

    //without a quote the fee will default to 0
    if (quoteMetadata.length == 0) {
      return fee;
    }

    //any report without a fee will default to 0
    if (report.length <= LINK_FEE_INDEX) {
      return fee;
    }

    //decode the quoteMetadata to get the desired fee
    address quoteFeeAddress = quoteMetadata.readAddress(QUOTE_METADATA_FEE_ADDRESS_INDEX);

    //calculate either the LINK fee or native fee if it's within the report
    if (quoteFeeAddress == LINK_ADDRESS) {
      fee.assetAddress = LINK_ADDRESS;
      fee.amount = report.readUint256(LINK_FEE_INDEX);
    } else {
      fee.assetAddress = NATIVE_ADDRESS;
      fee.amount = (report.readUint256(NATIVE_FEE_INDEX) * (TOTAL_DISCOUNT + nativePremium)) / TOTAL_DISCOUNT;
    }

    //decode the feedId from the report to calculate the discount being applied
    bytes32 feedId = bytes32(report);

    //set the fee amount to the discounted fee, rounding down
    fee.amount =
      fee.amount -
      ((fee.amount * subscriberDiscounts[subscriber][feedId][quoteFeeAddress]) / TOTAL_DISCOUNT);

    //return the fee
    return fee;
  }

  // @inheritdoc IFeeManager
  function setNativePremium(uint16 premium) external onlyOwner {
    if (premium > TOTAL_DISCOUNT) revert InvalidPremium();

    nativePremium = premium;
  }
}
