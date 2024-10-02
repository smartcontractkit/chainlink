// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IDataFeedsFeeManager} from "./interfaces/IDataFeedsFeeManager.sol";
import {IDataFeedsRouter} from "./interfaces/IDataFeedsRouter.sol";
import {ConfirmedOwner} from "../../shared/access/ConfirmedOwner.sol";
import {EnumerableSet} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/structs/EnumerableSet.sol";
import {Math} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/Math.sol";
import {IERC20} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/interfaces/IERC20.sol";
import {TypeAndVersionInterface} from "../../interfaces/TypeAndVersionInterface.sol";

contract DataFeedsFeeManager is ConfirmedOwner, IDataFeedsFeeManager, TypeAndVersionInterface {
  using EnumerableSet for EnumerableSet.AddressSet;

  string public constant override typeAndVersion = "DataFeedsFeeManager 1.0.0";

  uint256 public constant PERCENTAGE_SCALAR = 1e18;

  address private immutable i_router;

  struct FeeTokenConfig {
    bytes32 priceFeedId; // Feed ID of the token-USD price pair for conversions
    bytes32 feeTokenDiscountConfigId; // Config ID for the discount applied to fee token
  }

  struct SenderConfig {
    EnumerableSet.AddressSet balanceAddresses; // Balance addresses the sender may pull fees from
    bytes32 senderDiscountConfigId; // Config ID for the discount applied to sender
  }

  struct BillingData {
    address feeToken; // Fee token to pull from balanceAddress
    address balanceAddress; // Address to pull balances from
  }

  EnumerableSet.AddressSet private s_feeTokenSet;
  EnumerableSet.AddressSet private s_financeAdminSet;

  mapping(address tokenAddress => FeeTokenConfig feeTokenConfig) private s_feeTokenConfigs;
  mapping(address sender => SenderConfig senderConfig) private s_senderConfigs;
  mapping(bytes32 feedId => bytes32 configId) private s_feedServiceFeeConfigs;
  mapping(bytes32 configId => mapping(Service service => uint256 fee) serviceUsdFees) private s_serviceFeeConfigs;
  mapping(bytes32 configId => uint256 feeTokenDiscount) private s_feeTokenDiscountConfigs;
  mapping(bytes32 configId => uint256 senderDiscount) private s_senderDiscountConfigs;

  event FeedServiceFeesSet(bytes32 feedId, bytes32 configId);
  event FeeProcessed(address sender, address balanceAddress, address feeToken, uint256 totalFee);
  event FeeTokenAdded(address tokenAddress, bytes32 priceFeedId, bytes32 tokenDiscountConfigId);
  event FeeTokenRemoved(address tokenAddress);
  event FinanceAdminAdded(address financeAdmin);
  event FinanceAdminRemoved(address financeAdmin);
  event ServiceFeeConfigSet(
    bytes32 configId,
    uint256 getBenchmarkUsdFee,
    uint256 getReportUsdFee,
    uint256 requestUpkeepUsdFee
  );
  event SetFeeTokenDiscount(address tokenAddress, bytes32 configId);
  event SetFeeTokenDiscountConfig(bytes32 configId, uint256 discount);
  event SetSenderDiscount(address sender, bytes32 configId);
  event SetSenderDiscountConfig(bytes32 configId, uint256 discount);
  event SpendingAddressEnabled(address balanceAddress, address spendingAddress);
  event SpendingAddressDisabled(address balanceAddress, address spendingAddress);

  error InsufficientNative();
  error NoBenchmarkReturned();
  error UnauthorizedBalanceAddress(address sender, address balanceAddress);
  error UnauthorizedFeeManagement();
  error UnauthorizedFeeProcessing();
  error UnequalArrayLengths();

  modifier authorizeFeeManagement() {
    if (msg.sender != i_router && !s_financeAdminSet.contains(msg.sender)) revert UnauthorizedFeeManagement();

    _;
  }

  modifier authorizeFeeProcessing() {
    if (msg.sender != i_router) revert UnauthorizedFeeProcessing();

    _;
  }

  constructor(address router) ConfirmedOwner(msg.sender) {
    i_router = router;

    s_financeAdminSet.add(msg.sender);
  }

  /// @inheritdoc IDataFeedsFeeManager
  function processFee(
    address sender,
    Service service,
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) external payable authorizeFeeProcessing {
    uint256 totalFee = getFee(sender, service, feedIds, billingData);

    BillingData memory feeData = abi.decode(billingData, (BillingData));

    uint256 change = 0;

    address balanceAddress = feeData.balanceAddress == address(0) ? sender : feeData.balanceAddress;

    if (feeData.feeToken == address(0)) {
      if (msg.value < totalFee) revert InsufficientNative();
      change = msg.value - totalFee;
    } else {
      if (balanceAddress != sender && !s_senderConfigs[sender].balanceAddresses.contains(balanceAddress)) {
        revert UnauthorizedBalanceAddress(sender, balanceAddress);
      }
      IERC20(feeData.feeToken).transferFrom(balanceAddress, address(this), totalFee);

      change = msg.value;
    }

    payable(sender).transfer(change);

    emit FeeProcessed(sender, balanceAddress, feeData.feeToken, totalFee);
  }

  /// @inheritdoc IDataFeedsFeeManager
  function getFee(
    address sender,
    Service service,
    bytes32[] calldata feedIds,
    bytes calldata billingData
  ) public view returns (uint256 fee) {
    uint256 totalUsdFee = 0;
    for (uint256 i; i < feedIds.length; i++) {
      totalUsdFee += s_serviceFeeConfigs[s_feedServiceFeeConfigs[feedIds[i]]][service];
    }

    BillingData memory feeData = abi.decode(billingData, (BillingData));

    bytes32[] memory priceFeedIds = new bytes32[](1);
    priceFeedIds[0] = s_feeTokenConfigs[feeData.feeToken].priceFeedId;

    (int256[] memory latestPrice, ) = IDataFeedsRouter(i_router).getBenchmarksNonbillable(priceFeedIds);
    if (latestPrice.length != 1) revert NoBenchmarkReturned();

    uint256 totalFee = _calculateTotalFee(totalUsdFee, sender, feeData.feeToken, uint256(latestPrice[0]));

    return totalFee;
  }

  function _calculateTotalFee(
    uint256 totalUsdFee,
    address sender,
    address tokenAddress,
    uint256 latestPrice
  ) private view returns (uint256 fee) {
    uint256 senderMultiplier = PERCENTAGE_SCALAR -
      s_senderDiscountConfigs[s_senderConfigs[sender].senderDiscountConfigId];
    uint256 feeTokenMultiplier = PERCENTAGE_SCALAR -
      s_feeTokenDiscountConfigs[s_feeTokenConfigs[tokenAddress].feeTokenDiscountConfigId];

    uint256 totalFee = Math.ceilDiv(
      totalUsdFee * senderMultiplier * feeTokenMultiplier,
      latestPrice * PERCENTAGE_SCALAR
    );

    return totalFee;
  }

  /// @inheritdoc IDataFeedsFeeManager
  function enableSpendingAddresses(address[] calldata spendingAddresses) external {
    for (uint256 i; i < spendingAddresses.length; i++) {
      s_senderConfigs[spendingAddresses[i]].balanceAddresses.add(msg.sender);
      emit SpendingAddressEnabled(msg.sender, spendingAddresses[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function disableSpendingAddresses(address[] calldata spendingAddresses) external {
    for (uint256 i; i < spendingAddresses.length; i++) {
      s_senderConfigs[spendingAddresses[i]].balanceAddresses.remove(msg.sender);
      emit SpendingAddressDisabled(msg.sender, spendingAddresses[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function addFinanceAdmins(address[] calldata financeAdmins) external onlyOwner {
    for (uint256 i; i < financeAdmins.length; i++) {
      s_financeAdminSet.add(financeAdmins[i]);

      emit FinanceAdminAdded(financeAdmins[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function removeFinanceAdmins(address[] calldata financeAdmins) external onlyOwner {
    for (uint256 i; i < financeAdmins.length; i++) {
      s_financeAdminSet.remove(financeAdmins[i]);

      emit FinanceAdminRemoved(financeAdmins[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function addFeeTokens(
    address[] calldata tokenAddresses,
    bytes32[] calldata priceFeedIds,
    bytes32 feeTokenDiscountConfigId
  ) external authorizeFeeManagement {
    if (tokenAddresses.length != priceFeedIds.length) revert UnequalArrayLengths();

    for (uint256 i; i < tokenAddresses.length; i++) {
      s_feeTokenSet.add(tokenAddresses[i]);

      s_feeTokenConfigs[tokenAddresses[i]] = FeeTokenConfig({
        priceFeedId: priceFeedIds[i],
        feeTokenDiscountConfigId: feeTokenDiscountConfigId
      });

      emit FeeTokenAdded(tokenAddresses[i], priceFeedIds[i], feeTokenDiscountConfigId);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function removeFeeTokens(address[] calldata tokenAddresses) external authorizeFeeManagement {
    for (uint256 i; i < tokenAddresses.length; i++) {
      s_feeTokenSet.remove(tokenAddresses[i]);

      delete s_feeTokenConfigs[tokenAddresses[i]];

      emit FeeTokenRemoved(tokenAddresses[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function setServiceFeeConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata getBenchmarkUsdFees,
    uint256[] calldata getReportUsdFees,
    uint256[] calldata requestUpkeepUsdFees
  ) external authorizeFeeManagement {
    if (
      configIds.length != getBenchmarkUsdFees.length ||
      configIds.length != getReportUsdFees.length ||
      configIds.length != requestUpkeepUsdFees.length
    ) revert UnequalArrayLengths();

    for (uint256 i; i < configIds.length; i++) {
      s_serviceFeeConfigs[configIds[i]][Service.GetBenchmarks] = getBenchmarkUsdFees[i];
      s_serviceFeeConfigs[configIds[i]][Service.GetReports] = getReportUsdFees[i];
      s_serviceFeeConfigs[configIds[i]][Service.RequestUpkeep] = requestUpkeepUsdFees[i];

      emit ServiceFeeConfigSet(configIds[i], getBenchmarkUsdFees[i], getReportUsdFees[i], requestUpkeepUsdFees[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function setFeedServiceFees(bytes32 configId, bytes32[] calldata feedIds) external authorizeFeeManagement {
    for (uint256 i; i < feedIds.length; i++) {
      s_feedServiceFeeConfigs[feedIds[i]] = configId;

      emit FeedServiceFeesSet(feedIds[i], configId);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function setFeeTokenDiscountConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata discounts
  ) external authorizeFeeManagement {
    if (configIds.length != discounts.length) revert UnequalArrayLengths();

    for (uint256 i; i < configIds.length; i++) {
      s_feeTokenDiscountConfigs[configIds[i]] = discounts[i];

      emit SetFeeTokenDiscountConfig(configIds[i], discounts[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function setFeeTokenDiscounts(bytes32 configId, address[] calldata tokenAddresses) external authorizeFeeManagement {
    for (uint256 i; i < tokenAddresses.length; i++) {
      s_feeTokenConfigs[tokenAddresses[i]].feeTokenDiscountConfigId = configId;

      emit SetFeeTokenDiscount(tokenAddresses[i], configId);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function setSenderDiscountConfigs(
    bytes32[] calldata configIds,
    uint256[] calldata discounts
  ) external authorizeFeeManagement {
    if (configIds.length != discounts.length) revert UnequalArrayLengths();

    for (uint256 i; i < configIds.length; i++) {
      s_senderDiscountConfigs[configIds[i]] = discounts[i];

      emit SetSenderDiscountConfig(configIds[i], discounts[i]);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function setSenderDiscounts(bytes32 configId, address[] calldata senders) external authorizeFeeManagement {
    for (uint256 i; i < senders.length; i++) {
      s_senderConfigs[senders[i]].senderDiscountConfigId = configId;

      emit SetSenderDiscount(senders[i], configId);
    }
  }

  /// @inheritdoc IDataFeedsFeeManager
  function withdraw(
    address[] calldata tokenAddresses,
    uint256[] calldata quantities,
    address recipientAddress
  ) external authorizeFeeManagement {
    if (tokenAddresses.length != quantities.length) revert UnequalArrayLengths();

    for (uint256 i; i < tokenAddresses.length; i++) {
      address tokenAddress = tokenAddresses[i];
      if (tokenAddress == address(0)) {
        payable(recipientAddress).transfer(quantities[i]);
      } else {
        IERC20(tokenAddress).transfer(recipientAddress, quantities[i]);
      }
    }
  }
}
