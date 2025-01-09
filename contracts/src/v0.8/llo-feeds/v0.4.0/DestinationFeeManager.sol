// SPDX-License-Identifier: MIT
pragma solidity 0.8.19;

import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {ITypeAndVersion} from "../../shared/interfaces/ITypeAndVersion.sol";
import {IERC165} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC165.sol";
import {Common} from "../libraries/Common.sol";
import {IWERC20} from "../../shared/interfaces/IWERC20.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";
import {Math} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/Math.sol";
import {SafeERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";
import {IDestinationRewardManager} from "./interfaces/IDestinationRewardManager.sol";
import {IDestinationFeeManager} from "./interfaces/IDestinationFeeManager.sol";
import {IDestinationVerifierFeeManager} from "./interfaces/IDestinationVerifierFeeManager.sol";

/**
 * @title FeeManager
 * @author Michael Fletcher
 * @author Austin Born
 * @notice This contract is used for the handling of fees required for users verifying reports.
 */
contract DestinationFeeManager is
  IDestinationFeeManager,
  IDestinationVerifierFeeManager,
  ConfirmedOwner,
  ITypeAndVersion
{
  using SafeERC20 for IERC20;

  /// @notice list of subscribers and their discounts subscriberDiscounts[subscriber][feedId][token]
  mapping(address => mapping(bytes32 => mapping(address => uint256))) public s_subscriberDiscounts;

  /// @notice map of global discounts
  mapping(address => mapping(address => uint256)) public s_globalDiscounts;

  /// @notice keep track of any subsidised link that is owed to the reward manager.
  mapping(bytes32 => uint256) public s_linkDeficit;

  /// @notice the total discount that can be applied to a fee, 1e18 = 100% discount
  uint64 private constant PERCENTAGE_SCALAR = 1e18;

  /// @notice the LINK token address
  address public immutable i_linkAddress;

  /// @notice the native token address
  address public immutable i_nativeAddress;

  /// @notice the verifier address
  mapping(address => address) public s_verifierAddressList;

  /// @notice the reward manager address
  IDestinationRewardManager public i_rewardManager;

  // @notice the mask to apply to get the report version
  bytes32 private constant REPORT_VERSION_MASK = 0xffff000000000000000000000000000000000000000000000000000000000000;

  // @notice the different report versions
  bytes32 private constant REPORT_V1 = 0x0001000000000000000000000000000000000000000000000000000000000000;

  /// @notice the surcharge fee to be paid if paying in native
  uint256 public s_nativeSurcharge;

  /// @notice the error thrown if the discount or surcharge is invalid
  error InvalidSurcharge();

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

  // @notice thrown when the caller is not authorized
  error Unauthorized();

  // @notice thrown when trying to clear a zero deficit
  error ZeroDeficit();

  /// @notice thrown when trying to pay an address that cannot except funds
  error InvalidReceivingAddress();

  /// @notice thrown when trying to bulk verify reports where theres not a matching number of poolIds
  error PoolIdMismatch();

  /// @notice Emitted whenever a subscriber's discount is updated
  /// @param subscriber address of the subscriber to update discounts for
  /// @param feedId Feed ID for the discount
  /// @param token Token address for the discount
  /// @param discount Discount to apply, in relation to the PERCENTAGE_SCALAR
  event SubscriberDiscountUpdated(address indexed subscriber, bytes32 indexed feedId, address token, uint64 discount);

  /// @notice Emitted when updating the native surcharge
  /// @param newSurcharge Surcharge amount to apply relative to PERCENTAGE_SCALAR
  event NativeSurchargeUpdated(uint64 newSurcharge);

  /// @notice Emits when this contract does not have enough LINK to send to the reward manager when paying in native
  /// @param rewards Config digest and link fees which could not be subsidised
  event InsufficientLink(IDestinationRewardManager.FeePayment[] rewards);

  /// @notice Emitted when funds are withdrawn
  /// @param adminAddress Address of the admin
  /// @param recipient Address of the recipient
  /// @param assetAddress Address of the asset withdrawn
  /// @param quantity Amount of the asset withdrawn
  event Withdraw(address adminAddress, address recipient, address assetAddress, uint192 quantity);

  /// @notice Emits when a deficit has been cleared for a particular config digest
  /// @param configDigest Config digest of the deficit cleared
  /// @param linkQuantity Amount of LINK required to pay the deficit
  event LinkDeficitCleared(bytes32 indexed configDigest, uint256 linkQuantity);

  /// @notice Emits when a fee has been processed
  /// @param configDigest Config digest of the fee processed
  /// @param subscriber Address of the subscriber who paid the fee
  /// @param fee Fee paid
  /// @param reward Reward paid
  /// @param appliedDiscount Discount applied to the fee
  event DiscountApplied(
    bytes32 indexed configDigest,
    address indexed subscriber,
    Common.Asset fee,
    Common.Asset reward,
    uint256 appliedDiscount
  );

  /**
   * @notice Construct the FeeManager contract
   * @param _linkAddress The address of the LINK token
   * @param _nativeAddress The address of the wrapped ERC-20 version of the native token (represents fee in native or wrapped)
   * @param _verifierAddress The address of the verifier contract
   * @param _rewardManagerAddress The address of the reward manager contract
   */
  constructor(
    address _linkAddress,
    address _nativeAddress,
    address _verifierAddress,
    address _rewardManagerAddress
  ) ConfirmedOwner(msg.sender) {
    if (
      _linkAddress == address(0) ||
      _nativeAddress == address(0) ||
      _verifierAddress == address(0) ||
      _rewardManagerAddress == address(0)
    ) revert InvalidAddress();

    i_linkAddress = _linkAddress;
    i_nativeAddress = _nativeAddress;
    s_verifierAddressList[_verifierAddress] = _verifierAddress;
    i_rewardManager = IDestinationRewardManager(_rewardManagerAddress);

    IERC20(i_linkAddress).approve(address(i_rewardManager), type(uint256).max);
  }

  modifier onlyVerifier() {
    if (msg.sender != s_verifierAddressList[msg.sender]) revert Unauthorized();
    _;
  }

  /// @inheritdoc ITypeAndVersion
  function typeAndVersion() external pure override returns (string memory) {
    return "DestinationFeeManager 0.4.0";
  }

  /// @inheritdoc IERC165
  function supportsInterface(bytes4 interfaceId) external pure override returns (bool) {
    return
      interfaceId == type(IDestinationFeeManager).interfaceId ||
      interfaceId == type(IDestinationVerifierFeeManager).interfaceId;
  }

  /// @inheritdoc IDestinationVerifierFeeManager
  function processFee(
    bytes32 recipient,
    bytes calldata payload,
    bytes calldata parameterPayload,
    address subscriber
  ) external payable override onlyVerifier {
    (Common.Asset memory fee, Common.Asset memory reward, uint256 appliedDiscount) = _calculateFee(
      payload,
      parameterPayload,
      subscriber
    );

    if (fee.amount == 0) {
      _tryReturnChange(subscriber, msg.value);
      return;
    }

    IDestinationFeeManager.FeeAndReward[] memory feeAndReward = new IDestinationFeeManager.FeeAndReward[](1);
    feeAndReward[0] = IDestinationFeeManager.FeeAndReward(recipient, fee, reward, appliedDiscount);

    if (fee.assetAddress == i_linkAddress) {
      _handleFeesAndRewards(subscriber, feeAndReward, 1, 0);
    } else {
      _handleFeesAndRewards(subscriber, feeAndReward, 0, 1);
    }
  }

  /// @inheritdoc IDestinationVerifierFeeManager
  function processFeeBulk(
    bytes32[] memory poolIds,
    bytes[] calldata payloads,
    bytes calldata parameterPayload,
    address subscriber
  ) external payable override onlyVerifier {
    //poolIDs are mapped to payloads, so they should be the same length
    if (poolIds.length != payloads.length) revert PoolIdMismatch();

    IDestinationFeeManager.FeeAndReward[] memory feesAndRewards = new IDestinationFeeManager.FeeAndReward[](
      payloads.length
    );

    //keep track of the number of fees to prevent over initialising the FeePayment array within _convertToLinkAndNativeFees
    uint256 numberOfLinkFees;
    uint256 numberOfNativeFees;

    uint256 feesAndRewardsIndex;
    for (uint256 i; i < payloads.length; ++i) {
      if (poolIds[i] == bytes32(0)) revert InvalidAddress();

      (Common.Asset memory fee, Common.Asset memory reward, uint256 appliedDiscount) = _calculateFee(
        payloads[i],
        parameterPayload,
        subscriber
      );

      if (fee.amount != 0) {
        feesAndRewards[feesAndRewardsIndex++] = IDestinationFeeManager.FeeAndReward(
          poolIds[i],
          fee,
          reward,
          appliedDiscount
        );

        unchecked {
          //keep track of some tallys to make downstream calculations more efficient
          if (fee.assetAddress == i_linkAddress) {
            ++numberOfLinkFees;
          } else {
            ++numberOfNativeFees;
          }
        }
      }
    }

    if (numberOfLinkFees != 0 || numberOfNativeFees != 0) {
      _handleFeesAndRewards(subscriber, feesAndRewards, numberOfLinkFees, numberOfNativeFees);
    } else {
      _tryReturnChange(subscriber, msg.value);
    }
  }

  /// @inheritdoc IDestinationFeeManager
  function getFeeAndReward(
    address subscriber,
    bytes memory report,
    address quoteAddress
  ) public view returns (Common.Asset memory, Common.Asset memory, uint256) {
    Common.Asset memory fee;
    Common.Asset memory reward;

    //get the feedId from the report
    bytes32 feedId = bytes32(report);

    //the report needs to be a support version
    bytes32 reportVersion = _getReportVersion(feedId);

    //version 1 of the reports don't require quotes, so the fee will be 0
    if (reportVersion == REPORT_V1) {
      fee.assetAddress = i_nativeAddress;
      reward.assetAddress = i_linkAddress;
      return (fee, reward, 0);
    }

    //verify the quote payload is a supported token
    if (quoteAddress != i_nativeAddress && quoteAddress != i_linkAddress) {
      revert InvalidQuote();
    }

    //decode the report depending on the version
    uint256 linkQuantity;
    uint256 nativeQuantity;
    uint256 expiresAt;
    (, , , nativeQuantity, linkQuantity, expiresAt) = abi.decode(
      report,
      (bytes32, uint32, uint32, uint192, uint192, uint32)
    );

    //read the timestamp bytes from the report data and verify it has not expired
    if (expiresAt < block.timestamp) {
      revert ExpiredReport();
    }

    //check if feed discount has been applied
    uint256 discount = s_subscriberDiscounts[subscriber][feedId][quoteAddress];

    if (discount == 0) {
      //check if a global discount has been applied
      discount = s_globalDiscounts[subscriber][quoteAddress];
    }

    //the reward is always set in LINK
    reward.assetAddress = i_linkAddress;
    reward.amount = Math.ceilDiv(linkQuantity * (PERCENTAGE_SCALAR - discount), PERCENTAGE_SCALAR);

    //calculate either the LINK fee or native fee if it's within the report
    if (quoteAddress == i_linkAddress) {
      fee.assetAddress = i_linkAddress;
      fee.amount = reward.amount;
    } else {
      uint256 surchargedFee = Math.ceilDiv(nativeQuantity * (PERCENTAGE_SCALAR + s_nativeSurcharge), PERCENTAGE_SCALAR);

      fee.assetAddress = i_nativeAddress;
      fee.amount = Math.ceilDiv(surchargedFee * (PERCENTAGE_SCALAR - discount), PERCENTAGE_SCALAR);
    }

    //return the fee
    return (fee, reward, discount);
  }

  /// @inheritdoc IDestinationVerifierFeeManager
  function setFeeRecipients(
    bytes32 configDigest,
    Common.AddressAndWeight[] calldata rewardRecipientAndWeights
  ) external onlyVerifier {
    i_rewardManager.setRewardRecipients(configDigest, rewardRecipientAndWeights);
  }

  /// @inheritdoc IDestinationFeeManager
  function setNativeSurcharge(uint64 surcharge) external onlyOwner {
    if (surcharge > PERCENTAGE_SCALAR) revert InvalidSurcharge();

    s_nativeSurcharge = surcharge;

    emit NativeSurchargeUpdated(surcharge);
  }

  /// @inheritdoc IDestinationFeeManager
  function updateSubscriberDiscount(
    address subscriber,
    bytes32 feedId,
    address token,
    uint64 discount
  ) external onlyOwner {
    //make sure the discount is not greater than the total discount that can be applied
    if (discount > PERCENTAGE_SCALAR) revert InvalidDiscount();
    //make sure the token is either LINK or native
    if (token != i_linkAddress && token != i_nativeAddress) revert InvalidAddress();

    s_subscriberDiscounts[subscriber][feedId][token] = discount;

    emit SubscriberDiscountUpdated(subscriber, feedId, token, discount);
  }

  function updateSubscriberGlobalDiscount(address subscriber, address token, uint64 discount) external onlyOwner {
    //make sure the discount is not greater than the total discount that can be applied
    if (discount > PERCENTAGE_SCALAR) revert InvalidDiscount();
    //make sure the token is either LINK or native
    if (token != i_linkAddress && token != i_nativeAddress) revert InvalidAddress();

    s_globalDiscounts[subscriber][token] = discount;

    emit SubscriberDiscountUpdated(subscriber, bytes32(0), token, discount);
  }

  /// @inheritdoc IDestinationFeeManager
  function withdraw(address assetAddress, address recipient, uint192 quantity) external onlyOwner {
    //address 0 is used to withdraw native in the context of withdrawing
    if (assetAddress == address(0)) {
      (bool success, ) = payable(recipient).call{value: quantity}("");

      if (!success) revert InvalidReceivingAddress();
      return;
    }

    //withdraw the requested asset
    IERC20(assetAddress).safeTransfer(recipient, quantity);

    //emit event when funds are withdrawn
    emit Withdraw(msg.sender, recipient, assetAddress, uint192(quantity));
  }

  /// @inheritdoc IDestinationFeeManager
  function linkAvailableForPayment() external view returns (uint256) {
    //return the amount of LINK this contact has available to pay rewards
    return IERC20(i_linkAddress).balanceOf(address(this));
  }

  /**
   * @notice Gets the current version of the report that is encoded as the last two bytes of the feed
   * @param feedId feed id to get the report version for
   */
  function _getReportVersion(bytes32 feedId) internal pure returns (bytes32) {
    return REPORT_VERSION_MASK & feedId;
  }

  function _calculateFee(
    bytes calldata payload,
    bytes calldata parameterPayload,
    address subscriber
  ) internal view returns (Common.Asset memory, Common.Asset memory, uint256) {
    if (subscriber == address(this)) revert InvalidAddress();

    //decode the report from the payload
    (, bytes memory report) = abi.decode(payload, (bytes32[3], bytes));

    //get the feedId from the report
    bytes32 feedId = bytes32(report);

    //v1 doesn't need a quote payload, so skip the decoding
    address quote;
    if (_getReportVersion(feedId) != REPORT_V1) {
      //decode the quote from the bytes
      (quote) = abi.decode(parameterPayload, (address));
    }

    //decode the fee, it will always be native or LINK
    return getFeeAndReward(subscriber, report, quote);
  }

  function _handleFeesAndRewards(
    address subscriber,
    IDestinationFeeManager.FeeAndReward[] memory feesAndRewards,
    uint256 numberOfLinkFees,
    uint256 numberOfNativeFees
  ) internal {
    IDestinationRewardManager.FeePayment[] memory linkRewards = new IDestinationRewardManager.FeePayment[](
      numberOfLinkFees
    );
    IDestinationRewardManager.FeePayment[] memory nativeFeeLinkRewards = new IDestinationRewardManager.FeePayment[](
      numberOfNativeFees
    );

    uint256 totalNativeFee;
    uint256 totalNativeFeeLinkValue;

    uint256 linkRewardsIndex;
    uint256 nativeFeeLinkRewardsIndex;

    uint256 totalNumberOfFees = numberOfLinkFees + numberOfNativeFees;
    for (uint256 i; i < totalNumberOfFees; ++i) {
      if (feesAndRewards[i].fee.assetAddress == i_linkAddress) {
        linkRewards[linkRewardsIndex++] = IDestinationRewardManager.FeePayment(
          feesAndRewards[i].configDigest,
          uint192(feesAndRewards[i].reward.amount)
        );
      } else {
        nativeFeeLinkRewards[nativeFeeLinkRewardsIndex++] = IDestinationRewardManager.FeePayment(
          feesAndRewards[i].configDigest,
          uint192(feesAndRewards[i].reward.amount)
        );
        totalNativeFee += feesAndRewards[i].fee.amount;
        totalNativeFeeLinkValue += feesAndRewards[i].reward.amount;
      }

      if (feesAndRewards[i].appliedDiscount != 0) {
        emit DiscountApplied(
          feesAndRewards[i].configDigest,
          subscriber,
          feesAndRewards[i].fee,
          feesAndRewards[i].reward,
          feesAndRewards[i].appliedDiscount
        );
      }
    }

    //keep track of change in case of any over payment
    uint256 change;

    if (msg.value != 0) {
      //there must be enough to cover the fee
      if (totalNativeFee > msg.value) revert InvalidDeposit();

      //wrap the amount required to pay the fee & approve as the subscriber paid in wrapped native
      IWERC20(i_nativeAddress).deposit{value: totalNativeFee}();

      unchecked {
        //msg.value is always >= to fee.amount
        change = msg.value - totalNativeFee;
      }
    } else {
      if (totalNativeFee != 0) {
        //subscriber has paid in wrapped native, so transfer the native to this contract
        IERC20(i_nativeAddress).safeTransferFrom(subscriber, address(this), totalNativeFee);
      }
    }

    if (linkRewards.length != 0) {
      i_rewardManager.onFeePaid(linkRewards, subscriber);
    }

    if (nativeFeeLinkRewards.length != 0) {
      //distribute subsidised fees paid in Native
      if (totalNativeFeeLinkValue > IERC20(i_linkAddress).balanceOf(address(this))) {
        // If not enough LINK on this contract to forward for rewards, tally the deficit to be paid by out-of-band LINK
        for (uint256 i; i < nativeFeeLinkRewards.length; ++i) {
          unchecked {
            //we have previously tallied the fees, any overflows would have already reverted
            s_linkDeficit[nativeFeeLinkRewards[i].poolId] += nativeFeeLinkRewards[i].amount;
          }
        }

        emit InsufficientLink(nativeFeeLinkRewards);
      } else {
        //distribute the fees
        i_rewardManager.onFeePaid(nativeFeeLinkRewards, address(this));
      }
    }

    // a refund may be needed if the payee has paid in excess of the fee
    _tryReturnChange(subscriber, change);
  }

  function _tryReturnChange(address subscriber, uint256 quantity) internal {
    if (quantity != 0) {
      payable(subscriber).transfer(quantity);
    }
  }

  /// @inheritdoc IDestinationFeeManager
  function payLinkDeficit(bytes32 configDigest) external onlyOwner {
    uint256 deficit = s_linkDeficit[configDigest];

    if (deficit == 0) revert ZeroDeficit();

    delete s_linkDeficit[configDigest];

    IDestinationRewardManager.FeePayment[] memory deficitFeePayment = new IDestinationRewardManager.FeePayment[](1);

    deficitFeePayment[0] = IDestinationRewardManager.FeePayment(configDigest, uint192(deficit));

    i_rewardManager.onFeePaid(deficitFeePayment, address(this));

    emit LinkDeficitCleared(configDigest, deficit);
  }

  /// @inheritdoc IDestinationFeeManager
  function addVerifier(address verifierAddress) external onlyOwner {
    if (verifierAddress == address(0)) revert InvalidAddress();
    //check doesn't already exist
    if (s_verifierAddressList[verifierAddress] != address(0)) revert InvalidAddress();
    s_verifierAddressList[verifierAddress] = verifierAddress;
  }

  /// @inheritdoc IDestinationFeeManager
  function removeVerifier(address verifierAddress) external onlyOwner {
    if (verifierAddress == address(0)) revert InvalidAddress();
    //check doesn't already exist
    if (s_verifierAddressList[verifierAddress] == address(0)) revert InvalidAddress();
    delete s_verifierAddressList[verifierAddress];
  }

  /// @inheritdoc IDestinationFeeManager
  function setRewardManager(address rewardManagerAddress) external onlyOwner {
    if (rewardManagerAddress == address(0)) revert InvalidAddress();

    if (!IERC165(rewardManagerAddress).supportsInterface(type(IDestinationRewardManager).interfaceId)) {
      revert InvalidAddress();
    }

    IERC20(i_linkAddress).approve(address(i_rewardManager), 0);
    i_rewardManager = IDestinationRewardManager(rewardManagerAddress);
    IERC20(i_linkAddress).approve(address(rewardManagerAddress), type(uint256).max);
  }
}
