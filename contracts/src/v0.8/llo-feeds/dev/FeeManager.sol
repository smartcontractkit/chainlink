// SPDX-License-Identifier: MIT
pragma solidity 0.8.16;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {IFeeManager} from "./interfaces/IFeeManager.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.0/contracts/interfaces/IERC165.sol";
import {Common} from "../../libraries/Common.sol";
import {IRewardManager} from "./interfaces/IRewardManager.sol";
import {IWERC20} from "../../shared/interfaces/IWERC20.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.0/contracts/interfaces/IERC20.sol";
import {Math} from "../../vendor/openzeppelin-solidity/v4.8.0/contracts/utils/math/Math.sol";

/**
 * @title FeeManager
 * @author Michael Fletcher
 * @author Austin Born
 * @notice This contract is used for the handling of fees required for users verifying reports.
 */
contract FeeManager is IFeeManager, ConfirmedOwner, TypeAndVersionInterface {
  /// @notice list of subscribers and their discounts subscriberDiscounts[subscriber][feedId][token]
  mapping(address => mapping(bytes32 => mapping(address => uint256))) public s_subscriberDiscounts;

  /// @notice the total discount that can be applied to a fee, 1e18 = 100% discount
  uint256 private constant PERCENTAGE_SCALAR = 1e18;

  /// @notice the LINK token address
  address private immutable i_linkAddress;

  /// @notice the native token address
  address private immutable i_nativeAddress;

  /// @notice the proxy address
  address private immutable i_proxyAddress;

  /// @notice the reward manager address
  IRewardManager private immutable i_rewardManager;

  // @notice the mask to apply to get the report version
  bytes32 private constant REPORT_VERSION_MASK = 0xffff000000000000000000000000000000000000000000000000000000000000;

  // @notice the different report versions
  bytes32 private constant REPORT_V1 = 0x0001000000000000000000000000000000000000000000000000000000000000;
  bytes32 private constant REPORT_V2 = 0x0002000000000000000000000000000000000000000000000000000000000000;
  bytes32 private constant REPORT_V3 = 0x0003000000000000000000000000000000000000000000000000000000000000;

  /// @notice the surcharge fee to be paid if paying in native
  uint256 public s_nativeSurcharge;

  /// @notice the error thrown if the discount or surcharge is invalid
  error InvalidSurcharge();

  /// @notice the error thrown if the token is invalid
  error InvalidToken();

  /// @notice the error thrown if the discount is invalid
  error InvalidDiscount();

  /// @notice the error thrown if the address is invalid
  error InvalidAddress();

  /// @notice thrown if msg.value is supplied with a bad quote
  error InvalidDeposit();

  /// @notice thrown if a report has expired
  error ExpiredReport();

  /// @notice thrown if a report has no quote
  error InvalidQuote();

  /// @notice thrown if a report has an invalid version
  error InvalidReportVersion();

  // @notice Thrown when the caller is not authorized
  error Unauthorized();

  /// @notice Emitted whenever a subscriber's discount is updated
  /// @param subscriber address of the subscriber to update discounts for
  /// @param feedId Feed ID for the discount
  /// @param token Token address for the discount
  /// @param discount Discount to apply, in relation to the PERCENTAGE_SCALAR
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint256 discount);

  /// @notice Emitted when updating the native surcharge
  /// @param newSurcharge Surcharge amount to apply relative to PERCENTAGE_SCALAR
  event NativeSurchargeUpdated(uint256 newSurcharge);

  /// @notice Emits when this contract does not have enough LINK to send to the reward manager when paying in native
  /// @param configDigest Config digest of the report
  /// @param linkQuantity Amount of LINK required to pay the reward
  /// @param nativeQuantity Amount of native required to pay the reward
  event InsufficientLink(bytes32 indexed configDigest, uint256 linkQuantity, uint256 nativeQuantity);

  /// @notice Emitted when funds are withdrawn
  /// @param assetAddress Address of the asset withdrawn
  /// @param quantity Amount of the asset withdrawn
  event Withdraw(address adminAddress, address assetAddress, uint256 quantity);

  /**
   * @notice Construct the FeeManager contract
   * @param _linkAddress The address of the LINK token
   * @param _nativeAddress The address of the wrapped ERC-20 version of the native token (represents fee in native or wrapped)
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

    i_linkAddress = _linkAddress;
    i_nativeAddress = _nativeAddress;
    i_proxyAddress = _proxyAddress;
    i_rewardManager = IRewardManager(_rewardManagerAddress);

    IERC20(i_linkAddress).approve(address(i_rewardManager), type(uint256).max);
  }

  modifier onlyOwnerOrProxy() {
    if (msg.sender != owner() && msg.sender != i_proxyAddress) revert Unauthorized();
    _;
  }

  /// @inheritdoc TypeAndVersionInterface
  function typeAndVersion() external pure override returns (string memory) {
    return "FeeManager 0.0.1";
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return interfaceId == this.processFee.selector;
  }

  /// @inheritdoc IFeeManager
  function processFee(bytes calldata payload, address subscriber) external payable onlyOwnerOrProxy {
    if (subscriber == address(this)) revert InvalidAddress();

    //decode the report from the payload
    (, bytes memory report) = abi.decode(payload, (bytes32[3], bytes));

    //get the feedId from the report
    bytes32 feedId = bytes32(report);

    //v2 doesn't need a quote payload, so skip the decoding if the report is a v1 report
    Quote memory quote;
    if (getReportVersion(feedId) != REPORT_V1) {
      //all reports greater than v1 should have a quote payload
      (, , , , , bytes memory quoteBytes) = abi.decode(
        payload,
        (bytes32[3], bytes, bytes32[], bytes32[], bytes32, bytes)
      );

      //decode the quote from the bytes
      (quote) = abi.decode(quoteBytes, (Quote));
    }

    //decode the fee, it will always be native or LINK
    (Common.Asset memory fee, Common.Asset memory reward) = getFeeAndReward(msg.sender, report, quote);

    //keep track of change in case of any over payment
    uint256 change;

    //wrap the amount required to pay the fee
    if (msg.value != 0) {
      //quote must be in native with enough to cover the fee
      if (fee.assetAddress != i_nativeAddress) revert InvalidDeposit();
      if (fee.amount > msg.value) revert InvalidDeposit();

      //wrap the amount required to pay the fee & approve
      IWERC20(i_nativeAddress).deposit{value: fee.amount}();

      unchecked {
        //msg.value is always >= to fee.amount
        change = msg.value - fee.amount;
      }
    }

    //get the config digest which is the first 32 bytes of the payload
    bytes32 configDigest = bytes32(payload);

    //some users might not be billed
    if (fee.amount != 0) {
      //if the fee is in LINK, transfer directly from the subscriber to the reward manager
      if (fee.assetAddress == i_linkAddress) {
        //distributes the fee
        i_rewardManager.onFeePaid(configDigest, subscriber, reward.amount);
      } else {
        //if the fee is in native wrapped, transfer to this contract in exchange for the equivalent amount of LINK excluding the surcharge
        if (msg.value == 0) {
          IERC20(fee.assetAddress).transferFrom(subscriber, address(this), fee.amount);
        }

        //check that the contract has enough LINK before paying the fee
        if (reward.amount > IERC20(i_linkAddress).balanceOf(address(this))) {
          // If not enough LINK on this contract to forward for rewards, fire this event and
          // call onFeePaid out-of-band to pay out rewards
          emit InsufficientLink(configDigest, reward.amount, fee.amount);
        } else {
          //bill the payee and distribute the fee using the config digest as the key
          i_rewardManager.onFeePaid(configDigest, address(this), reward.amount);
        }
      }
    }

    // a refund may be needed if the payee has paid in excess of the fee
    if (change != 0) {
      payable(subscriber).transfer(change);
    }
  }

  /// @inheritdoc IFeeManager
  function getFeeAndReward(
    address subscriber,
    bytes memory report,
    Quote memory quote
  ) public view returns (Common.Asset memory, Common.Asset memory) {
    Common.Asset memory fee;
    Common.Asset memory reward;

    //get the feedId from the report
    bytes32 feedId = bytes32(report);

    //the report needs to be a support version
    bytes32 reportVersion = getReportVersion(feedId);

    //version 1 of the reports don't require quotes, so the fee will be 0
    if (reportVersion == REPORT_V1) {
      fee.assetAddress = i_nativeAddress;
      reward.assetAddress = i_linkAddress;
      return (fee, reward);
    }

    //verify the quote payload is a supported token
    if (quote.quoteAddress != i_nativeAddress && quote.quoteAddress != i_linkAddress) {
      revert InvalidQuote();
    }

    //decode the report depending on the version
    uint256 linkQuantity;
    uint256 nativeQuantity;
    uint256 expiresAt;
    if (reportVersion == REPORT_V2) {
      (, , , , expiresAt, linkQuantity, nativeQuantity) = abi.decode(
        report,
        (bytes32, uint32, int192, uint32, uint32, uint192, uint192)
      );
    } else if (reportVersion == REPORT_V3) {
      (, , , , , , expiresAt, linkQuantity, nativeQuantity) = abi.decode(
        report,
        (bytes32, uint32, int192, int192, int192, uint32, uint32, uint192, uint192)
      );
    } else {
      revert InvalidReportVersion();
    }

    //read the timestamp bytes from the report data and verify it has not expired
    if (expiresAt < block.timestamp) {
      revert ExpiredReport();
    }

    //the reward is always set in LINK
    reward.assetAddress = i_linkAddress;
    reward.amount = linkQuantity;

    //calculate either the LINK fee or native fee if it's within the report
    if (quote.quoteAddress == i_linkAddress) {
      fee.assetAddress = reward.assetAddress;
      fee.amount = reward.amount;
    } else {
      fee.assetAddress = i_nativeAddress;
      fee.amount = Math.ceilDiv(nativeQuantity * (PERCENTAGE_SCALAR + s_nativeSurcharge), PERCENTAGE_SCALAR);
    }

    //get the discount being applied
    uint256 discount = s_subscriberDiscounts[subscriber][feedId][quote.quoteAddress];

    //apply the discount to the fee, rounding up
    fee.amount = fee.amount - ((fee.amount * discount) / PERCENTAGE_SCALAR);

    //apply the discount to the reward, rounding down
    reward.amount = reward.amount - Math.ceilDiv(reward.amount * discount, PERCENTAGE_SCALAR);

    //return the fee
    return (fee, reward);
  }

  /// @inheritdoc IFeeManager
  function setFeeRecipients(
    bytes32 configDigest,
    Common.AddressAndWeight[] calldata rewardRecipientAndWeights
  ) external onlyOwnerOrProxy {
    i_rewardManager.setRewardRecipients(configDigest, rewardRecipientAndWeights);
  }

  /// @inheritdoc IFeeManager
  function setNativeSurcharge(uint256 surcharge) external onlyOwner {
    if (surcharge > PERCENTAGE_SCALAR) revert InvalidSurcharge();

    s_nativeSurcharge = surcharge;

    emit NativeSurchargeUpdated(surcharge);
  }

  /// @inheritdoc IFeeManager
  function updateSubscriberDiscount(
    address subscriber,
    bytes32 feedId,
    address token,
    uint256 discount
  ) external onlyOwner {
    //make sure the discount is not greater than the total discount that can be applied
    if (discount > PERCENTAGE_SCALAR) revert InvalidDiscount();
    //make sure the token is either LINK or native
    if (token != i_linkAddress && token != i_nativeAddress) revert InvalidAddress();

    s_subscriberDiscounts[subscriber][feedId][token] = discount;

    emit SubscriberDiscountUpdated(subscriber, feedId, token, discount);
  }

  /// @inheritdoc IFeeManager
  function withdraw(address assetAddress, uint256 quantity) external onlyOwner {
    //address 0 is used to withdraw native in the context of withdrawing
    if (assetAddress == address(0)) {
      payable(owner()).transfer(quantity);
      return;
    }

    //withdraw the requested asset
    IERC20(assetAddress).transfer(owner(), quantity);

    //emit event when funds are withdrawn
    emit Withdraw(msg.sender, assetAddress, quantity);
  }

  function linkAvailableForPayment() external view returns (uint256) {
    //return the amount of LINK this contact has available to pay rewards
    return IERC20(i_linkAddress).balanceOf(address(this));
  }

  /**
   * @notice Gets the current version of the report that is encoded as the last two bytes of the feed
   * @param feedId feed id to get the report version for
   */
  function getReportVersion(bytes32 feedId) internal pure returns (bytes32) {
    return REPORT_VERSION_MASK & feedId;
  }
}
