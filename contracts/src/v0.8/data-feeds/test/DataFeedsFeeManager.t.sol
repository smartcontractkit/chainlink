// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {DataFeedsTestSetup} from "./DataFeedsTestSetup.t.sol";
import {DataFeedsFeeManager} from "../dev/DataFeedsFeeManager.sol";
import {IDataFeedsFeeManager} from "../dev/interfaces/IDataFeedsFeeManager.sol";
import {IDataFeedsRouter} from "../dev/interfaces/IDataFeedsRouter.sol";
import {Math} from "../../vendor/openzeppelin-solidity/v4.8.3/contracts/utils/math/Math.sol";

contract DataFeedsFeemanagerTest is DataFeedsTestSetup {
  DataFeedsFeeManager internal dataFeedsFeeManager;
  address internal constant ROUTER = address(10001);
  address[] internal FIN_ADMINS = [address(10002)];
  address internal constant SENDER = address(10003);
  address internal constant BALANCE_ADDRESS = address(10004);
  bytes32[] internal FEED_IDS = [reportStructBasic.feedId, reportStructPremium.feedId];
  int256[] internal BENCHMARKS = [reportStructBasic.benchmarkPrice, reportStructPremium.benchmarkPrice];
  uint256[] internal TIMESTAMPS = [reportStructBasic.observationsTimestamp, reportStructPremium.observationsTimestamp];
  address internal constant NATIVE_FEE_TOKEN = address(0);
  address[] internal FEE_TOKENS = new address[](2);
  bytes32[] internal FEE_TOKEN_FEED_IDS = [reportStructBasic.feedId, reportStructBasic.feedId];
  int256[] internal FEE_BENCHMARK = [BENCHMARKS[0]];
  uint256[] internal FEE_TIMESTAMP = [TIMESTAMPS[0]];
  bytes32[] internal SERVICE_FEE_CONFIG_ID = [bytes32("3")];
  uint256[] internal GET_BENCHMARK_FEE = [12345];
  uint256[] internal GET_REPORT_FEE = [23456];
  uint256[] internal REQUEST_UPKEEP_FEE = [34567];

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

  function setUp() public virtual override {
    DataFeedsTestSetup.setUp();

    dataFeedsFeeManager = new DataFeedsFeeManager(ROUTER);

    dataFeedsFeeManager.addFinanceAdmins(FIN_ADMINS);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setServiceFeeConfigs(
      SERVICE_FEE_CONFIG_ID,
      GET_BENCHMARK_FEE,
      GET_REPORT_FEE,
      REQUEST_UPKEEP_FEE
    );

    dataFeedsFeeManager.setFeedServiceFees(SERVICE_FEE_CONFIG_ID[0], FEED_IDS);

    FEE_TOKENS[0] = NATIVE_FEE_TOKEN;
    FEE_TOKENS[1] = address(link);

    dataFeedsFeeManager.addFeeTokens(FEE_TOKENS, FEE_TOKEN_FEED_IDS, bytes32(""));

    vm.startPrank(BALANCE_ADDRESS);

    address[] memory senders = new address[](1);
    senders[0] = SENDER;

    dataFeedsFeeManager.enableSpendingAddresses(senders);
  }

  function test_processFeeNative() public {
    uint256 totalFee = Math.ceilDiv(GET_BENCHMARK_FEE[0] * 2 * PERCENTAGE_SCALAR, uint256(FEE_BENCHMARK[0]));
    bytes32[] memory feeTokenFeedIds = new bytes32[](1);
    feeTokenFeedIds[0] = FEE_TOKEN_FEED_IDS[0];

    vm.deal(ROUTER, totalFee);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IDataFeedsRouter.getBenchmarksNonbillable.selector, feeTokenFeedIds),
      abi.encode(FEE_BENCHMARK, FEE_TIMESTAMP)
    );

    vm.expectEmit();
    emit FeeProcessed(SENDER, SENDER, NATIVE_FEE_TOKEN, totalFee);

    vm.startPrank(ROUTER);

    dataFeedsFeeManager.processFee{value: totalFee}(
      SENDER,
      IDataFeedsFeeManager.Service.GetBenchmarks,
      FEED_IDS,
      abi.encode(NATIVE_FEE_TOKEN, NATIVE_FEE_TOKEN)
    );
  }

  function test_processFeeERC20() public {
    uint256 totalFee = Math.ceilDiv(GET_REPORT_FEE[0] * 2 * PERCENTAGE_SCALAR, uint256(FEE_BENCHMARK[0]));
    bytes32[] memory feeTokenFeedIds = new bytes32[](1);
    feeTokenFeedIds[0] = FEE_TOKEN_FEED_IDS[1];

    link.mint(BALANCE_ADDRESS, totalFee);

    vm.startPrank(BALANCE_ADDRESS);

    link.approve(address(dataFeedsFeeManager), totalFee);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IDataFeedsRouter.getBenchmarksNonbillable.selector, feeTokenFeedIds),
      abi.encode(FEE_BENCHMARK, FEE_TIMESTAMP)
    );

    vm.expectEmit();
    emit FeeProcessed(SENDER, BALANCE_ADDRESS, address(link), totalFee);

    vm.startPrank(ROUTER);

    dataFeedsFeeManager.processFee(
      SENDER,
      IDataFeedsFeeManager.Service.GetReports,
      FEED_IDS,
      abi.encode(address(link), BALANCE_ADDRESS)
    );
  }

  function test_processFeeRevertUnauthorizedFeeProcessing() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeProcessing.selector));

    vm.startPrank(SENDER);

    dataFeedsFeeManager.processFee(SENDER, IDataFeedsFeeManager.Service.GetBenchmarks, FEED_IDS, bytes(""));
  }

  function test_processFeeRevertInsufficientNative() public {
    uint256 totalFee = Math.ceilDiv(GET_BENCHMARK_FEE[0] * 2 * PERCENTAGE_SCALAR, uint256(FEE_BENCHMARK[0]));
    bytes32[] memory feeTokenFeedIds = new bytes32[](1);
    feeTokenFeedIds[0] = FEE_TOKEN_FEED_IDS[0];

    vm.deal(ROUTER, totalFee - 1);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IDataFeedsRouter.getBenchmarksNonbillable.selector, feeTokenFeedIds),
      abi.encode(FEE_BENCHMARK, FEE_TIMESTAMP)
    );

    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.InsufficientNative.selector));

    vm.startPrank(ROUTER);

    dataFeedsFeeManager.processFee{value: totalFee - 1}(
      SENDER,
      IDataFeedsFeeManager.Service.GetBenchmarks,
      FEED_IDS,
      abi.encode(NATIVE_FEE_TOKEN, NATIVE_FEE_TOKEN)
    );
  }

  function test_processFeeRevertUnauthorizedBalanceAddress() public {
    uint256 totalFee = Math.ceilDiv(GET_REPORT_FEE[0] * 2 * PERCENTAGE_SCALAR, uint256(FEE_BENCHMARK[0]));
    bytes32[] memory feeTokenFeedIds = new bytes32[](1);
    feeTokenFeedIds[0] = FEE_TOKEN_FEED_IDS[1];

    address unauthorizedBalanceAddress = address(10005);

    link.mint(unauthorizedBalanceAddress, totalFee);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IDataFeedsRouter.getBenchmarksNonbillable.selector, feeTokenFeedIds),
      abi.encode(FEE_BENCHMARK, FEE_TIMESTAMP)
    );

    vm.expectRevert(
      abi.encodeWithSelector(
        DataFeedsFeeManager.UnauthorizedBalanceAddress.selector,
        SENDER,
        unauthorizedBalanceAddress
      )
    );

    vm.startPrank(ROUTER);

    dataFeedsFeeManager.processFee(
      SENDER,
      IDataFeedsFeeManager.Service.GetReports,
      FEED_IDS,
      abi.encode(address(link), unauthorizedBalanceAddress)
    );
  }

  function test_processFeeRevertAmountExceedsBalance() public {
    uint256 totalFee = Math.ceilDiv(GET_REPORT_FEE[0] * 2 * PERCENTAGE_SCALAR, uint256(FEE_BENCHMARK[0]));
    bytes32[] memory feeTokenFeedIds = new bytes32[](1);
    feeTokenFeedIds[0] = FEE_TOKEN_FEED_IDS[1];

    vm.startPrank(BALANCE_ADDRESS);

    link.approve(address(dataFeedsFeeManager), totalFee);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IDataFeedsRouter.getBenchmarksNonbillable.selector, feeTokenFeedIds),
      abi.encode(FEE_BENCHMARK, FEE_TIMESTAMP)
    );

    vm.expectRevert("ERC20: transfer amount exceeds balance");

    vm.startPrank(ROUTER);

    dataFeedsFeeManager.processFee(
      SENDER,
      IDataFeedsFeeManager.Service.GetReports,
      FEED_IDS,
      abi.encode(address(link), BALANCE_ADDRESS)
    );
  }

  function test_getFee() public {
    bytes32[] memory feeTokenFeedIds = new bytes32[](1);
    feeTokenFeedIds[0] = FEE_TOKEN_FEED_IDS[0];

    bytes32[] memory feeTokenDiscountConfigs = new bytes32[](1);
    feeTokenDiscountConfigs[0] = bytes32("5");

    uint256[] memory feeTokenDiscounts = new uint256[](1);
    feeTokenDiscounts[0] = 1e17;

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setFeeTokenDiscountConfigs(feeTokenDiscountConfigs, feeTokenDiscounts);

    dataFeedsFeeManager.setFeeTokenDiscounts(feeTokenDiscountConfigs[0], FEE_TOKENS);

    bytes32[] memory senderDiscountConfigs = new bytes32[](1);
    senderDiscountConfigs[0] = bytes32("6");

    uint256[] memory senderDiscounts = new uint256[](1);
    senderDiscounts[0] = 2e17;

    address[] memory senders = new address[](1);
    senders[0] = SENDER;

    dataFeedsFeeManager.setSenderDiscountConfigs(senderDiscountConfigs, senderDiscounts);

    dataFeedsFeeManager.setSenderDiscounts(senderDiscountConfigs[0], senders);

    uint256 totalFee = Math.ceilDiv(
      GET_BENCHMARK_FEE[0] *
        2 *
        PERCENTAGE_SCALAR *
        (PERCENTAGE_SCALAR - feeTokenDiscounts[0]) *
        (PERCENTAGE_SCALAR - senderDiscounts[0]),
      uint256(FEE_BENCHMARK[0]) * PERCENTAGE_SCALAR * PERCENTAGE_SCALAR
    );

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IDataFeedsRouter.getBenchmarksNonbillable.selector, feeTokenFeedIds),
      abi.encode(FEE_BENCHMARK, FEE_TIMESTAMP)
    );

    uint256 fee = dataFeedsFeeManager.getFee(
      SENDER,
      IDataFeedsFeeManager.Service.GetBenchmarks,
      FEED_IDS,
      abi.encode(NATIVE_FEE_TOKEN, bytes32(""))
    );

    assertEq(totalFee, fee);
  }

  function test_getFeeRevertNoBenchmarkReturned() public {
    bytes32[] memory feeTokenFeedIds = new bytes32[](1);
    feeTokenFeedIds[0] = FEE_TOKEN_FEED_IDS[0];

    uint256[] memory benchmarks = new uint256[](0);
    uint256[] memory timestamps = new uint256[](0);

    vm.mockCall(
      ROUTER,
      abi.encodeWithSelector(IDataFeedsRouter.getBenchmarksNonbillable.selector, feeTokenFeedIds),
      abi.encode(benchmarks, timestamps)
    );

    vm.expectRevert(DataFeedsFeeManager.NoBenchmarkReturned.selector);

    dataFeedsFeeManager.getFee(
      SENDER,
      IDataFeedsFeeManager.Service.GetBenchmarks,
      FEED_IDS,
      abi.encode(NATIVE_FEE_TOKEN, bytes32(""))
    );
  }

  function test_enableSpendingAddresses() public {
    address[] memory senders = new address[](2);
    senders[0] = address(10005);
    senders[1] = address(10006);

    vm.expectEmit();
    emit SpendingAddressEnabled(BALANCE_ADDRESS, senders[0]);

    vm.expectEmit();
    emit SpendingAddressEnabled(BALANCE_ADDRESS, senders[1]);

    vm.startPrank(BALANCE_ADDRESS);

    dataFeedsFeeManager.enableSpendingAddresses(senders);
  }

  function test_disableSpendingAddresses() public {
    address[] memory senders = new address[](2);
    senders[0] = address(10005);
    senders[1] = address(10006);

    vm.expectEmit();
    emit SpendingAddressDisabled(BALANCE_ADDRESS, senders[0]);

    vm.expectEmit();
    emit SpendingAddressDisabled(BALANCE_ADDRESS, senders[1]);

    vm.startPrank(BALANCE_ADDRESS);

    dataFeedsFeeManager.disableSpendingAddresses(senders);
  }

  function test_addFinanceAdmins() public {
    address[] memory financeAdmins = new address[](2);
    financeAdmins[0] = address(10005);
    financeAdmins[1] = address(10006);

    vm.expectEmit();
    emit FinanceAdminAdded(financeAdmins[0]);

    vm.expectEmit();
    emit FinanceAdminAdded(financeAdmins[1]);

    vm.startPrank(OWNER);

    dataFeedsFeeManager.addFinanceAdmins(financeAdmins);
  }

  function test_addFinanceAdminsRevertOnlyOwner() public {
    address[] memory financeAdmins = new address[](2);
    financeAdmins[0] = address(10005);
    financeAdmins[1] = address(10006);

    vm.expectRevert("Only callable by owner");

    dataFeedsFeeManager.addFinanceAdmins(financeAdmins);
  }

  function test_removeFinanceAdmins() public {
    address[] memory financeAdmins = new address[](2);
    financeAdmins[0] = address(10005);
    financeAdmins[1] = address(10006);

    vm.expectEmit();
    emit FinanceAdminRemoved(financeAdmins[0]);

    vm.expectEmit();
    emit FinanceAdminRemoved(financeAdmins[1]);

    vm.startPrank(OWNER);

    dataFeedsFeeManager.removeFinanceAdmins(financeAdmins);
  }

  function test_removeFinanceAdminsRevertOnlyOwner() public {
    address[] memory financeAdmins = new address[](2);
    financeAdmins[0] = address(10005);
    financeAdmins[1] = address(10006);

    vm.expectRevert("Only callable by owner");

    dataFeedsFeeManager.removeFinanceAdmins(financeAdmins);
  }

  function test_addFeeTokens() public {
    address[] memory tokenAddresses = new address[](2);
    tokenAddresses[0] = address(10005);
    tokenAddresses[1] = address(10006);

    bytes32[] memory priceFeedIds = new bytes32[](2);
    priceFeedIds[0] = FEE_TOKEN_FEED_IDS[0];
    priceFeedIds[1] = FEE_TOKEN_FEED_IDS[1];

    bytes32 feeTokenDiscountConfigId = bytes32("");

    vm.expectEmit();
    emit FeeTokenAdded(tokenAddresses[0], priceFeedIds[0], feeTokenDiscountConfigId);

    vm.expectEmit();
    emit FeeTokenAdded(tokenAddresses[1], priceFeedIds[1], feeTokenDiscountConfigId);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.addFeeTokens(tokenAddresses, priceFeedIds, feeTokenDiscountConfigId);
  }

  function test_addFeeTokensRevertUnauthorizedFeeManagement() public {
    address[] memory tokenAddresses = new address[](1);
    tokenAddresses[0] = address(10005);

    bytes32[] memory priceFeedIds = new bytes32[](1);
    priceFeedIds[0] = FEE_TOKEN_FEED_IDS[0];

    bytes32 feeTokenDiscountConfigId = bytes32("");

    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.addFeeTokens(tokenAddresses, priceFeedIds, feeTokenDiscountConfigId);
  }

  function test_addFeeTokensRevertUnequalArrayLengths() public {
    address[] memory tokenAddresses = new address[](2);
    tokenAddresses[0] = address(10005);
    tokenAddresses[1] = address(10006);

    bytes32[] memory priceFeedIds = new bytes32[](1);
    priceFeedIds[0] = FEE_TOKEN_FEED_IDS[0];

    bytes32 feeTokenDiscountConfigId = bytes32("");

    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.addFeeTokens(tokenAddresses, priceFeedIds, feeTokenDiscountConfigId);
  }

  function test_removeFeeTokens() public {
    address[] memory tokenAddresses = new address[](2);
    tokenAddresses[0] = address(10005);
    tokenAddresses[1] = address(10006);

    vm.expectEmit();
    emit FeeTokenRemoved(tokenAddresses[0]);

    vm.expectEmit();
    emit FeeTokenRemoved(tokenAddresses[1]);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.removeFeeTokens(tokenAddresses);
  }

  function test_removeFeeTokensRevertUnauthorizedFeeManagement() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.removeFeeTokens(FEE_TOKENS);
  }

  function test_setServiceFeeConfigs() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = bytes32("3");
    configIds[0] = bytes32("4");

    uint256[] memory getBenchmarkUsdFees = new uint256[](2);
    getBenchmarkUsdFees[0] = 12345;
    getBenchmarkUsdFees[1] = 23456;

    uint256[] memory getReportUsdFees = new uint256[](2);
    getReportUsdFees[0] = 34567;
    getReportUsdFees[1] = 45678;

    uint256[] memory requestUpkeepUsdFees = new uint256[](2);
    requestUpkeepUsdFees[0] = 56789;
    requestUpkeepUsdFees[1] = 67890;

    vm.expectEmit();
    emit ServiceFeeConfigSet(configIds[0], getBenchmarkUsdFees[0], getReportUsdFees[0], requestUpkeepUsdFees[0]);

    vm.expectEmit();
    emit ServiceFeeConfigSet(configIds[1], getBenchmarkUsdFees[1], getReportUsdFees[1], requestUpkeepUsdFees[1]);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setServiceFeeConfigs(configIds, getBenchmarkUsdFees, getReportUsdFees, requestUpkeepUsdFees);
  }

  function test_setServiceFeeConfigsRevertUnauthorizedFeeManagement() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.setServiceFeeConfigs(
      SERVICE_FEE_CONFIG_ID,
      GET_BENCHMARK_FEE,
      GET_REPORT_FEE,
      REQUEST_UPKEEP_FEE
    );
  }

  function test_setServiceFeeConfigsRevertUnequalArrayLengths() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = bytes32("3");
    configIds[0] = bytes32("4");

    uint256[] memory getBenchmarkUsdFees = new uint256[](2);
    getBenchmarkUsdFees[0] = 12345;
    getBenchmarkUsdFees[1] = 23456;

    uint256[] memory getReportUsdFees = new uint256[](2);
    getReportUsdFees[0] = 34567;
    getReportUsdFees[1] = 45678;

    uint256[] memory requestUpkeepUsdFees = new uint256[](1);
    requestUpkeepUsdFees[0] = 56789;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnequalArrayLengths.selector));

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setServiceFeeConfigs(configIds, getBenchmarkUsdFees, getReportUsdFees, requestUpkeepUsdFees);
  }

  function test_setFeedServiceFees() public {
    vm.expectEmit();
    emit FeedServiceFeesSet(FEED_IDS[0], SERVICE_FEE_CONFIG_ID[0]);

    vm.expectEmit();
    emit FeedServiceFeesSet(FEED_IDS[1], SERVICE_FEE_CONFIG_ID[0]);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setFeedServiceFees(SERVICE_FEE_CONFIG_ID[0], FEED_IDS);
  }

  function test_setFeedServiceFeesRevertUnauthorizedFeeManagement() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.setFeedServiceFees(SERVICE_FEE_CONFIG_ID[0], FEED_IDS);
  }

  function test_setFeeTokenDiscountConfigs() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = bytes32("3");
    configIds[0] = bytes32("4");

    uint256[] memory discounts = new uint256[](2);
    discounts[0] = 12345;
    discounts[1] = 23456;

    vm.expectEmit();
    emit SetFeeTokenDiscountConfig(configIds[0], discounts[0]);

    vm.expectEmit();
    emit SetFeeTokenDiscountConfig(configIds[1], discounts[1]);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setFeeTokenDiscountConfigs(configIds, discounts);
  }

  function test_setFeeTokenDiscountConfigsRevertUnauthorizedFeeManagement() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.setFeeTokenDiscountConfigs(SERVICE_FEE_CONFIG_ID, GET_BENCHMARK_FEE);
  }

  function test_setFeeTokenDiscountConfigsRevertUnequalArrayLengths() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = bytes32("3");
    configIds[0] = bytes32("4");

    uint256[] memory discounts = new uint256[](1);
    discounts[0] = 12345;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnequalArrayLengths.selector));

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setFeeTokenDiscountConfigs(configIds, discounts);
  }

  function test_setFeeTokenDiscounts() public {
    vm.expectEmit();
    emit SetFeeTokenDiscount(FEE_TOKENS[0], SERVICE_FEE_CONFIG_ID[0]);

    vm.expectEmit();
    emit SetFeeTokenDiscount(FEE_TOKENS[1], SERVICE_FEE_CONFIG_ID[0]);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setFeeTokenDiscounts(SERVICE_FEE_CONFIG_ID[0], FEE_TOKENS);
  }

  function test_setFeeTokenDiscountsRevertUnauthorizedFeeManagement() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.setFeeTokenDiscounts(SERVICE_FEE_CONFIG_ID[0], FEE_TOKENS);
  }

  function test_setSenderDiscountConfigs() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = bytes32("3");
    configIds[0] = bytes32("4");

    uint256[] memory discounts = new uint256[](2);
    discounts[0] = 12345;
    discounts[1] = 23456;

    vm.expectEmit();
    emit SetSenderDiscountConfig(configIds[0], discounts[0]);

    vm.expectEmit();
    emit SetSenderDiscountConfig(configIds[1], discounts[1]);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setSenderDiscountConfigs(configIds, discounts);
  }

  function test_setSenderDiscountConfigsRevertUnauthorizedFeeManagement() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.setSenderDiscountConfigs(SERVICE_FEE_CONFIG_ID, GET_BENCHMARK_FEE);
  }

  function test_setSenderDiscountConfigsRevertUnequalArrayLengths() public {
    bytes32[] memory configIds = new bytes32[](2);
    configIds[0] = bytes32("3");
    configIds[0] = bytes32("4");

    uint256[] memory discounts = new uint256[](1);
    discounts[0] = 12345;

    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnequalArrayLengths.selector));

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setSenderDiscountConfigs(configIds, discounts);
  }

  function test_setSenderDiscounts() public {
    address[] memory senders = new address[](2);
    senders[0] = address(12345);
    senders[1] = address(23456);

    vm.expectEmit();
    emit SetSenderDiscount(senders[0], SERVICE_FEE_CONFIG_ID[0]);

    vm.expectEmit();
    emit SetSenderDiscount(senders[1], SERVICE_FEE_CONFIG_ID[0]);

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.setSenderDiscounts(SERVICE_FEE_CONFIG_ID[0], senders);
  }

  function test_setSenderDiscountsRevertUnauthorizedFeeManagement() public {
    address[] memory senders = new address[](2);
    senders[0] = address(12345);
    senders[1] = address(23456);

    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.setSenderDiscounts(SERVICE_FEE_CONFIG_ID[0], senders);
  }

  function test_withdraw() public {
    uint256[] memory amounts = new uint256[](2);
    amounts[0] = 12345;
    amounts[1] = 23456;

    vm.deal(address(dataFeedsFeeManager), amounts[0]);

    link.mint(address(dataFeedsFeeManager), amounts[1]);

    vm.startPrank(FIN_ADMINS[0]);

    assertEq(address(dataFeedsFeeManager).balance, amounts[0]);
    assertEq(BALANCE_ADDRESS.balance, 0);
    assertEq(link.balanceOf(address(dataFeedsFeeManager)), amounts[1]);
    assertEq(link.balanceOf(BALANCE_ADDRESS), 0);

    dataFeedsFeeManager.withdraw(FEE_TOKENS, amounts, BALANCE_ADDRESS);

    assertEq(address(dataFeedsFeeManager).balance, 0);
    assertEq(BALANCE_ADDRESS.balance, amounts[0]);
    assertEq(link.balanceOf(address(dataFeedsFeeManager)), 0);
    assertEq(link.balanceOf(BALANCE_ADDRESS), amounts[1]);
  }

  function test_withdrawRevertUnauthorizedFeeManagement() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnauthorizedFeeManagement.selector));

    dataFeedsFeeManager.withdraw(FEE_TOKENS, GET_BENCHMARK_FEE, BALANCE_ADDRESS);
  }

  function test_withdrawRevertUnequalArrayLengths() public {
    vm.expectRevert(abi.encodeWithSelector(DataFeedsFeeManager.UnequalArrayLengths.selector));

    vm.startPrank(FIN_ADMINS[0]);

    dataFeedsFeeManager.withdraw(FEE_TOKENS, GET_BENCHMARK_FEE, BALANCE_ADDRESS);
  }

  function test_typeAndVersion() public {
    string memory typeAndVersion = dataFeedsFeeManager.typeAndVersion();
    assertEq(typeAndVersion, "DataFeedsFeeManager 1.0.0", "typeAndVersion should match expected value");
  }
}
