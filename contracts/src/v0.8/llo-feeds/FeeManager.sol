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
  mapping(address => mapping(bytes32 => mapping(address => uint256))) public subscriberDiscounts;

  //the total discount that can be applied to a fee, 1e18 = 100% discount
  uint256 private constant PERCENTAGE_SCALAR = 1e18;

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
  uint256 public nativePremium;

  //the error thrown if the discount or premium is invalid
  error InvalidPremium();

  //the error thrown if the token is invalid
  error InvalidToken();

  //the error thrown if the discount is invalid
  error InvalidDiscount();

  //the error thrown if the address is invalid
  error InvalidAddress();

  // Events emitted upon state change
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint256 discount);
  event NativePremiumSet(uint256 newPremium);

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
  function updateSubscriberDiscount(
    address subscriber,
    bytes32 feedId,
    address token,
    uint256 discount
  ) external onlyOwner {
    //make sure the discount is not greater than the total discount that can be applied
    if (discount > PERCENTAGE_SCALAR) revert InvalidDiscount();
    //make sure the token is either LINK or native
    if (token != LINK_ADDRESS && token != NATIVE_ADDRESS) revert InvalidAddress();

    subscriberDiscounts[subscriber][feedId][token] = discount;

    //emit the event
    emit SubscriberDiscountUpdated(subscriber, feedId, token, discount);
  }

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
      fee.amount = (report.readUint256(NATIVE_FEE_INDEX) * (PERCENTAGE_SCALAR + nativePremium)) / PERCENTAGE_SCALAR;
    }

    //decode the feedId from the report to calculate the discount being applied
    bytes32 feedId = bytes32(report);

    //set the fee amount to the discounted fee, rounding down
    fee.amount =
      fee.amount -
      ((fee.amount * subscriberDiscounts[subscriber][feedId][quoteFeeAddress]) / PERCENTAGE_SCALAR);

    //return the fee
    return fee;
  }

  // @inheritdoc IFeeManager
  function setNativePremium(uint256 premium) external onlyOwner {
    if (premium > PERCENTAGE_SCALAR) revert InvalidPremium();

    nativePremium = premium;

    //emit the event
    emit NativePremiumSet(premium);
  }
}
