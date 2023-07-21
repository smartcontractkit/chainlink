// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../ConfirmedOwner.sol";
import {IFeeManager} from "./interfaces/IFeeManager.sol";
import {TypeAndVersionInterface} from "../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../shared/vendor/IERC165.sol";
import {Common} from "../libraries/internal/Common.sol";
import {IRewardManager} from "./interfaces/IRewardManager.sol";
import {IWERC20} from "../shared/vendor/IWERC20.sol";
import {IERC20} from "../shared/vendor/IERC20.sol";
import {Math} from "../shared/vendor/Math.sol";

/*
 * @title FeeManager
 * @author Austin Born
 * @author Michael Fletcher
 * @notice This contract is used for the handling of fees required for users verifying reports.
 */
contract FeeManager is IFeeManager, ConfirmedOwner, TypeAndVersionInterface {
  //list of subscribers and their discounts subscriberDiscounts[subscriber][feedId][token]
  mapping(address => mapping(bytes32 => mapping(address => uint256))) public subscriberDiscounts;

  //the total discount that can be applied to a fee, 1e18 = 100% discount
  uint256 private constant PERCENTAGE_SCALAR = 1e18;

  //the link token address
  address private immutable linkAddress;

  //the native token address
  address private immutable nativeAddress;

  //the proxy address
  address private immutable proxyAddress;

  //the reward manager address
  IRewardManager private immutable rewardManager;

  //the report packed length, each field is packed into 32 bytes
  uint16 private constant DEFAULT_REPORT_LENGTH = 32 + 32 + 32 + 32 + 32 + 32 + 32 + 32;

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

  //thrown if msg.value is supplied with a bad quote
  error InvalidDeposit();

  //thrown if a report has expired
  error ExpiredReport();

  // Events emitted upon state change
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint256 discount);
  event NativePremiumSet(uint256 newPremium);
  event InsufficientLink(bytes32 indexed configDigest, uint256 linkQuantity, uint256 nativeQuantity);

  struct Quote {
    address quoteAddress;
  }

  /**
   * @notice Construct the FeeManager contract
   * @param _linkAddress The address of the LINK token
   * @param _nativeAddress The address of the NATIVE token
   * @param _proxyAddress The address of the proxy contract
   * @param _rewardManagerAddress The address of the reward manager contract
   */
  constructor(
    address _linkAddress,
    address _nativeAddress,
    address _proxyAddress,
    address _rewardManagerAddress
  ) ConfirmedOwner(msg.sender) {
    if (
      _linkAddress == address(0) ||
      _nativeAddress == address(0) ||
      _proxyAddress == address(0) ||
      _rewardManagerAddress == address(0)
    ) revert InvalidAddress();

    linkAddress = _linkAddress;
    nativeAddress = _nativeAddress;
    proxyAddress = _proxyAddress;
    rewardManager = IRewardManager(_rewardManagerAddress);
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "FeeManager 0.0.1";
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return interfaceId == this.processFee.selector;
  }

  modifier onlyOwnerOrProxy() {
    require(msg.sender == owner() || msg.sender == proxyAddress, "Only owner or proxy");
    _;
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
    if (token != linkAddress && token != nativeAddress) revert InvalidAddress();

    subscriberDiscounts[subscriber][feedId][token] = discount;

    //emit the event
    emit SubscriberDiscountUpdated(subscriber, feedId, token, discount);
  }

  /**
   * @notice Calculate the applied fee and the reward from a report. If the sender is a subscriber, they will receive a discount.
   * @param subscriber address trying to verify
   * @param report report to calculate the fee for
   * @param quote any metadata required to fetch the fee
   * @return the fee and the reward data
   */
  function getFeeAndReward(
    address subscriber,
    bytes memory report,
    Quote memory quote
  ) public view returns (Common.Asset memory, Common.Asset memory) {
    //The fee and reward
    Common.Asset memory feeQuantity;
    Common.Asset memory rewardQuantity;

    //any report without a fee does not need to be processed
    if (report.length <= DEFAULT_REPORT_LENGTH) {
      return (feeQuantity, rewardQuantity);
    }

    //decode the fee
    (, , , , , , , , uint256 linkQuantity, uint256 nativeQuantity, uint256 expiresAt) = abi.decode(
      report,
      (bytes32, uint32, int192, int192, int192, uint64, bytes32, uint64, uint256, uint256, uint32)
    );

    // Read the timestamp bytes from the report data and verify it has not expired
    if (expiresAt < block.timestamp) {
      revert ExpiredReport();
    }

    //without a quote the fee will default to billing in link if a quote is not provided
    address quoteFeeAddress;
    if (quote.quoteAddress == address(0)) {
      quoteFeeAddress = linkAddress;
    } else {
      //decode the quoteMetadata to get the desired asset to pay the quote in
      quoteFeeAddress = quote.quoteAddress;
    }

    //the fee paid is always in LINK
    rewardQuantity.assetAddress = linkAddress;
    rewardQuantity.amount = linkQuantity;

    //calculate either the LINK fee or native fee if it's within the report
    if (quoteFeeAddress == linkAddress) {
      feeQuantity.assetAddress = rewardQuantity.assetAddress;
      feeQuantity.amount = rewardQuantity.amount;
    } else {
      feeQuantity.assetAddress = nativeAddress;
      feeQuantity.amount = (nativeQuantity * (PERCENTAGE_SCALAR + nativePremium)) / PERCENTAGE_SCALAR;
    }

    //decode the feedId from the report to calculate the discount being applied
    bytes32 feedId = bytes32(report);
    uint256 discount = subscriberDiscounts[subscriber][feedId][quoteFeeAddress];

    //apply the discount to the fee, rounding up
    feeQuantity.amount = feeQuantity.amount - ((feeQuantity.amount * discount) / PERCENTAGE_SCALAR);

    //apply the discount to the reward, rounding down
    rewardQuantity.amount = rewardQuantity.amount - Math.ceilDiv(rewardQuantity.amount * discount, PERCENTAGE_SCALAR);

    //return the fee
    return (feeQuantity, rewardQuantity);
  }

  // @inheritdoc IFeeManager
  function processFee(bytes calldata payload, address subscriber) external payable onlyOwnerOrProxy {
    //decode the payload
    (, bytes memory report, , , , bytes memory quoteBytes) = abi.decode(
      payload,
      (bytes32[3], bytes, bytes32[], bytes32[], bytes32, bytes)
    );

    //reports without quotes are valid, decode the quote if there are quote bytes
    Quote memory quote;
    if (quoteBytes.length > 0) {
      quote = abi.decode(quoteBytes, (Quote));
    }

    //decode the fee, it will always be native or link
    (Common.Asset memory fee, Common.Asset memory reward) = getFeeAndReward(msg.sender, report, quote);

    //keep track of change in case of any over payment
    uint256 change;

    //wrap the amount required to pay the fee
    if (msg.value > 0) {
      //quote must be in native with enough to cover the fee
      if (fee.assetAddress != nativeAddress) revert InvalidDeposit();
      if (fee.amount > msg.value) revert InvalidDeposit();

      //wrap the amount required to pay the bill & approve
      IWERC20(nativeAddress).deposit{value: fee.amount}();

      unchecked {
        //msg.value is always >= to fee.amount
        change = msg.value - fee.amount;
      }
    }

    //get the config digest which is the first 32 bytes of the payload
    bytes32 configDigest = bytes32(payload);

    //some users might not be billed
    if (fee.amount > 0) {
      //if the fee is in link, we're transferring directly from the subscriber, else the contract is covering the link
      if (fee.assetAddress == linkAddress) {
        //bill the payee and distribute the fee
        rewardManager.onFeePaid(configDigest, subscriber, reward);
      } else {
        //if the fee is in native wrapped, we're transferring to this contract in exchange for the equivalent amount of link (minus the native premium)
        if (msg.value == 0) {
          IERC20(fee.assetAddress).transferFrom(msg.sender, address(this), fee.amount);
        }

        //check we have enough link before paying the fee
        if (reward.amount > IERC20(linkAddress).balanceOf(address(this))) {
          //approve the transfer of link required to verify the report to the reward manager
          IERC20(linkAddress).approve(address(rewardManager), reward.amount);

          //bill the payee and distribute the fee using the config digest as the key
          rewardManager.onFeePaid(configDigest, address(this), reward);
        } else {
          //contract does not have enough link
          emit InsufficientLink(configDigest, reward.amount, fee.amount);
        }
      }
    }

    //we may need to refund if the payee paid in excess of the fee
    if (change > 0) {
      payable(subscriber).transfer(change);
    }
  }

  // @inheritdoc IFeeManager
  function setFeeRecipients(
    bytes32 configDigest,
    Common.AddressAndWeight[] calldata rewardRecipientAndWeights
  ) external onlyOwnerOrProxy {
    rewardManager.setRewardRecipients(configDigest, rewardRecipientAndWeights);
  }

  // @inheritdoc IFeeManager
  function setNativePremium(uint256 premium) external onlyOwner {
    if (premium > PERCENTAGE_SCALAR) revert InvalidPremium();

    nativePremium = premium;

    //emit the event
    emit NativePremiumSet(premium);
  }

  // @inheritdoc IFeeManager
  function withdraw(address assetAddress, uint256 quantity) external onlyOwner {
    //address 0 is used to withdraw native
    if (assetAddress == address(0)) {
      payable(owner()).transfer(quantity);
      return;
    }

    //withdraw the requested asset
    IERC20(assetAddress).transfer(owner(), quantity);
  }
}
