// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../dev/v1_X/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/v1_X/FunctionsBilling.sol";
import {FunctionsRequest} from "../../dev/v1_X/libraries/FunctionsRequest.sol";
import {FunctionsResponse} from "../../dev/v1_X/libraries/FunctionsResponse.sol";
import {FunctionsSubscriptions} from "../../dev/v1_X/FunctionsSubscriptions.sol";
import {Routable} from "../../dev/v1_X/Routable.sol";

import {FunctionsRouterSetup, FunctionsSubscriptionSetup, FunctionsClientRequestSetup, FunctionsFulfillmentSetup, FunctionsMultipleFulfillmentsSetup} from "./Setup.t.sol";

import {FunctionsBillingConfig} from "../../dev/v1_X/interfaces/IFunctionsBilling.sol";

/// @notice #constructor
contract FunctionsBilling_Constructor is FunctionsSubscriptionSetup {
  function test_Constructor_Success() public {
    assertEq(address(s_functionsRouter), s_functionsCoordinator.getRouter_HARNESS());
    assertEq(address(s_linkEthFeed), s_functionsCoordinator.getLinkToNativeFeed_HARNESS());
  }
}

/// @notice #getConfig
contract FunctionsBilling_GetConfig is FunctionsRouterSetup {
  function test_GetConfig_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    FunctionsBillingConfig memory config = s_functionsCoordinator.getConfig();
    assertEq(config.feedStalenessSeconds, getCoordinatorConfig().feedStalenessSeconds);
    assertEq(config.gasOverheadBeforeCallback, getCoordinatorConfig().gasOverheadBeforeCallback);
    assertEq(config.gasOverheadAfterCallback, getCoordinatorConfig().gasOverheadAfterCallback);
    assertEq(config.requestTimeoutSeconds, getCoordinatorConfig().requestTimeoutSeconds);
    assertEq(config.donFeeCentsUsd, getCoordinatorConfig().donFeeCentsUsd);
    assertEq(config.maxSupportedRequestDataVersion, getCoordinatorConfig().maxSupportedRequestDataVersion);
    assertEq(config.fulfillmentGasPriceOverEstimationBP, getCoordinatorConfig().fulfillmentGasPriceOverEstimationBP);
    assertEq(config.fallbackNativePerUnitLink, getCoordinatorConfig().fallbackNativePerUnitLink);
  }
}

/// @notice #updateConfig
contract FunctionsBilling_UpdateConfig is FunctionsRouterSetup {
  FunctionsBillingConfig internal configToSet;

  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();

    // Multiply all config values by 2 to confirm that they change
    configToSet = FunctionsBillingConfig({
      feedStalenessSeconds: getCoordinatorConfig().feedStalenessSeconds * 2,
      gasOverheadAfterCallback: getCoordinatorConfig().gasOverheadAfterCallback * 2,
      gasOverheadBeforeCallback: getCoordinatorConfig().gasOverheadBeforeCallback * 2,
      requestTimeoutSeconds: getCoordinatorConfig().requestTimeoutSeconds * 2,
      donFeeCentsUsd: getCoordinatorConfig().donFeeCentsUsd * 2,
      operationFeeCentsUsd: getCoordinatorConfig().operationFeeCentsUsd * 2,
      maxSupportedRequestDataVersion: getCoordinatorConfig().maxSupportedRequestDataVersion * 2,
      fulfillmentGasPriceOverEstimationBP: getCoordinatorConfig().fulfillmentGasPriceOverEstimationBP * 2,
      fallbackNativePerUnitLink: getCoordinatorConfig().fallbackNativePerUnitLink * 2,
      fallbackUsdPerUnitLink: getCoordinatorConfig().fallbackUsdPerUnitLink * 2,
      fallbackUsdPerUnitLinkDecimals: getCoordinatorConfig().fallbackUsdPerUnitLinkDecimals * 2,
      minimumEstimateGasPriceWei: getCoordinatorConfig().minimumEstimateGasPriceWei * 2
    });
  }

  function test_UpdateConfig_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsCoordinator.updateConfig(configToSet);
  }

  event ConfigUpdated(FunctionsBillingConfig config);

  function test_UpdateConfig_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit ConfigUpdated(configToSet);

    s_functionsCoordinator.updateConfig(configToSet);

    FunctionsBillingConfig memory config = s_functionsCoordinator.getConfig();
    assertEq(config.feedStalenessSeconds, configToSet.feedStalenessSeconds);
    assertEq(config.gasOverheadAfterCallback, configToSet.gasOverheadAfterCallback);
    assertEq(config.gasOverheadBeforeCallback, configToSet.gasOverheadBeforeCallback);
    assertEq(config.requestTimeoutSeconds, configToSet.requestTimeoutSeconds);
    assertEq(config.donFeeCentsUsd, configToSet.donFeeCentsUsd);
    assertEq(config.operationFeeCentsUsd, configToSet.operationFeeCentsUsd);
    assertEq(config.maxSupportedRequestDataVersion, configToSet.maxSupportedRequestDataVersion);
    assertEq(config.fulfillmentGasPriceOverEstimationBP, configToSet.fulfillmentGasPriceOverEstimationBP);
    assertEq(config.fallbackNativePerUnitLink, configToSet.fallbackNativePerUnitLink);
    assertEq(config.minimumEstimateGasPriceWei, configToSet.minimumEstimateGasPriceWei);
    assertEq(config.fallbackUsdPerUnitLink, configToSet.fallbackUsdPerUnitLink);
    assertEq(config.fallbackUsdPerUnitLinkDecimals, configToSet.fallbackUsdPerUnitLinkDecimals);
  }
}

/// @notice #getDONFee
contract FunctionsBilling_GetDONFeeJuels is FunctionsRouterSetup {
  function test_GetDONFeeJuels_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint72 donFee = s_functionsCoordinator.getDONFeeJuels(new bytes(0));
    uint72 expectedDonFee = uint72(((s_donFee * 10 ** (18 + LINK_USD_DECIMALS)) / uint256(LINK_USD_RATE)) / 100);
    assertEq(donFee, expectedDonFee);
  }
}

/// @notice #getOperationFee
contract FunctionsBilling_GetOperationFee is FunctionsRouterSetup {
  function test_GetOperationFeeJuels_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint72 operationFee = s_functionsCoordinator.getOperationFeeJuels();
    uint72 expectedOperationFee = uint72(
      ((s_operationFee * 10 ** (18 + LINK_USD_DECIMALS)) / uint256(LINK_USD_RATE)) / 100
    );
    assertEq(operationFee, expectedOperationFee);
  }
}

/// @notice #getAdminFee
contract FunctionsBilling_GetAdminFeeJuels is FunctionsRouterSetup {
  function test_GetAdminFeeJuels_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint72 adminFee = s_functionsCoordinator.getAdminFeeJuels();
    assertEq(adminFee, s_adminFee);
  }
}

/// @notice #getWeiPerUnitLink
contract FunctionsBilling_GetWeiPerUnitLink is FunctionsRouterSetup {
  function test_GetWeiPerUnitLink_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint256 weiPerUnitLink = s_functionsCoordinator.getWeiPerUnitLink();
    assertEq(weiPerUnitLink, uint256(LINK_ETH_RATE));
  }
}

/// @notice #estimateCost
contract FunctionsBilling_EstimateCost is FunctionsSubscriptionSetup {
  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Get cost estimate as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);
  }

  uint256 private constant REASONABLE_GAS_PRICE_CEILING = 1_000_000_000_000_000; // 1 million gwei

  function test_EstimateCost_RevertsIfGasPriceAboveCeiling() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest._initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest._encodeCBOR(request);

    uint32 callbackGasLimit = 5_500;
    uint256 gasPriceWei = REASONABLE_GAS_PRICE_CEILING + 1;

    vm.expectRevert(FunctionsBilling.InvalidCalldata.selector);

    s_functionsCoordinator.estimateCost(s_subscriptionId, requestData, callbackGasLimit, gasPriceWei);
  }

  function test_EstimateCost_SuccessLowGasPrice() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest._initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest._encodeCBOR(request);

    uint32 callbackGasLimit = 5_500;
    uint256 gasPriceWei = 1;

    uint96 costEstimate = s_functionsCoordinator.estimateCost(
      s_subscriptionId,
      requestData,
      callbackGasLimit,
      gasPriceWei
    );
    uint96 expectedCostEstimate = 51110500000000000 +
      s_adminFee +
      s_functionsCoordinator.getDONFeeJuels(requestData) +
      s_functionsCoordinator.getOperationFeeJuels();
    assertEq(costEstimate, expectedCostEstimate);
  }

  function test_EstimateCost_Success() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest._initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest._encodeCBOR(request);

    uint32 callbackGasLimit = 5_500;
    uint256 gasPriceWei = 5000000000; // 5 gwei

    uint96 costEstimate = s_functionsCoordinator.estimateCost(
      s_subscriptionId,
      requestData,
      callbackGasLimit,
      gasPriceWei
    );
    uint96 expectedCostEstimate = 255552500000000000 +
      s_adminFee +
      s_functionsCoordinator.getDONFeeJuels(requestData) +
      s_functionsCoordinator.getOperationFeeJuels();
    assertEq(costEstimate, expectedCostEstimate);
  }
}

/// @notice #_calculateCostEstimate
contract FunctionsBilling__CalculateCostEstimate {
  // TODO: make contract internal function helper
}

/// @notice #_startBilling
contract FunctionsBilling__StartBilling is FunctionsFulfillmentSetup {
  function test__FulfillAndBill_HasUniqueGlobalRequestId() public {
    // Variables that go into a requestId:
    // - Coordinator address
    // - Consumer contract
    // - Subscription ID,
    // - Consumer initiated requests
    // - Request data
    // - Request data version
    // - Request callback gas limit
    // - Estimated total cost in Juels
    // - Request timeout timestamp
    // - tx.origin

    // Request #1 has already been fulfilled by the test setup

    // Reset the nonce (initiatedRequests) by removing and re-adding the consumer
    s_functionsRouter.removeConsumer(s_subscriptionId, address(s_functionsClient));
    assertEq(s_functionsRouter.getSubscription(s_subscriptionId).consumers.length, 0);
    s_functionsRouter.addConsumer(s_subscriptionId, address(s_functionsClient));
    assertEq(s_functionsRouter.getSubscription(s_subscriptionId).consumers[0], address(s_functionsClient));

    // Make Request #2
    _sendAndStoreRequest(
      2,
      s_requests[1].requestData.sourceCode,
      s_requests[1].requestData.secrets,
      s_requests[1].requestData.args,
      s_requests[1].requestData.bytesArgs,
      s_requests[1].requestData.callbackGasLimit
    );

    // Request #1 and #2 should have different request IDs, because the request timeout timestamp has advanced.
    // A request cannot be fulfilled in the same block, which prevents removing a consumer in the same block
    assertNotEq(s_requests[1].requestId, s_requests[2].requestId);
  }
}

/// @notice #_fulfillAndBill
contract FunctionsBilling__FulfillAndBill is FunctionsClientRequestSetup {
  function test__FulfillAndBill_RevertIfInvalidCommitment() public {
    vm.expectRevert();
    s_functionsCoordinator.fulfillAndBill_HARNESS(
      s_requests[1].requestId,
      new bytes(0),
      new bytes(0),
      new bytes(0), // malformed commitment data
      new bytes(0),
      1
    );
  }

  event RequestBilled(
    bytes32 indexed requestId,
    uint96 juelsPerGas,
    uint256 l1FeeShareWei,
    uint96 callbackCostJuels,
    uint72 donFee,
    uint72 adminFee,
    uint72 operationFee
  );

  function test__FulfillAndBill_Success() public {
    uint96 juelsPerGas = uint96((1e18 * TX_GASPRICE_START) / uint256(LINK_ETH_RATE));
    uint96 callbackCostGas = 5072; // Taken manually
    uint96 callbackCostJuels = juelsPerGas * callbackCostGas;

    // topic0 (function signature, always checked), check topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit RequestBilled(
      s_requests[1].requestId,
      juelsPerGas,
      0,
      callbackCostJuels,
      s_functionsCoordinator.getDONFeeJuels(new bytes(0)),
      s_adminFee,
      s_functionsCoordinator.getOperationFeeJuels()
    );

    FunctionsResponse.FulfillResult resultCode = s_functionsCoordinator.fulfillAndBill_HARNESS(
      s_requests[1].requestId,
      new bytes(0),
      new bytes(0),
      abi.encode(s_requests[1].commitment),
      new bytes(0),
      1
    );

    assertEq(uint256(resultCode), uint256(FunctionsResponse.FulfillResult.FULFILLED));
  }
}

/// @notice #deleteCommitment
contract FunctionsBilling_DeleteCommitment is FunctionsClientRequestSetup {
  function test_DeleteCommitment_RevertIfNotRouter() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(Routable.OnlyCallableByRouter.selector);
    s_functionsCoordinator.deleteCommitment(s_requests[1].requestId);
  }

  event CommitmentDeleted(bytes32 requestId);

  function test_DeleteCommitment_Success() public {
    // Send as Router
    vm.stopPrank();
    vm.startPrank(address(s_functionsRouter));

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit CommitmentDeleted(s_requests[1].requestId);

    s_functionsCoordinator.deleteCommitment(s_requests[1].requestId);
  }
}

/// @notice #oracleWithdraw
contract FunctionsBilling_OracleWithdraw is FunctionsMultipleFulfillmentsSetup {
  function test_OracleWithdraw_RevertWithNoBalance() public {
    uint256[4] memory transmitterBalancesBefore = _getTransmitterBalances();
    _assertTransmittersAllHaveBalance(transmitterBalancesBefore, 0);

    // Send as stranger, which has no balance
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);

    // Attempt to withdraw with no amount, which would withdraw the full balance
    s_functionsCoordinator.oracleWithdraw(STRANGER_ADDRESS, 0);

    uint256[4] memory transmitterBalancesAfter = _getTransmitterBalances();
    _assertTransmittersAllHaveBalance(transmitterBalancesAfter, 0);
  }

  function test_OracleWithdraw_RevertIfInsufficientBalance() public {
    // Send as transmitter 1, which has transmitted 1 report
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_1);

    vm.expectRevert(FunctionsBilling.InsufficientBalance.selector);

    // Attempt to withdraw more than the Coordinator has assigned
    s_functionsCoordinator.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, s_fulfillmentCoordinatorBalance + 1);
  }

  function test_OracleWithdraw_SuccessTransmitterWithBalanceValidAmountGiven() public {
    uint256[4] memory transmitterBalancesBefore = _getTransmitterBalances();
    _assertTransmittersAllHaveBalance(transmitterBalancesBefore, 0);

    // Send as transmitter 1, which has transmitted 1 report
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_1);

    uint96 expectedTransmitterBalance = s_fulfillmentCoordinatorBalance / s_requestsFulfilled;

    // Attempt to withdraw half of balance
    uint96 halfBalance = expectedTransmitterBalance / 2;
    s_functionsCoordinator.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, halfBalance);

    uint256[4] memory transmitterBalancesAfter = _getTransmitterBalances();
    assertEq(transmitterBalancesAfter[0], halfBalance);
    assertEq(transmitterBalancesAfter[1], 0);
    assertEq(transmitterBalancesAfter[2], 0);
    assertEq(transmitterBalancesAfter[3], 0);
  }

  function test_OracleWithdraw_SuccessTransmitterWithBalanceNoAmountGiven() public {
    uint256[4] memory transmitterBalancesBefore = _getTransmitterBalances();
    _assertTransmittersAllHaveBalance(transmitterBalancesBefore, 0);

    // Send as transmitter 1, which has transmitted 2 reports
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_1);

    // Attempt to withdraw with no amount, which will withdraw the full balance
    s_functionsCoordinator.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, 0);

    uint96 totalOperationFees = s_functionsCoordinator.getOperationFeeJuels() * s_requestsFulfilled;
    uint96 totalDonFees = s_functionsCoordinator.getDONFeeJuels(new bytes(0)) * s_requestsFulfilled;
    uint96 donFeeShare = totalDonFees / uint8(s_transmitters.length);
    uint96 expectedBalancePerFulfillment = ((s_fulfillmentCoordinatorBalance - totalOperationFees - totalDonFees) /
      s_requestsFulfilled);

    uint256[4] memory transmitterBalancesAfter = _getTransmitterBalances();
    // Transmitter 1 has transmitted twice
    assertEq(transmitterBalancesAfter[0], (expectedBalancePerFulfillment * 2) + donFeeShare);
    assertEq(transmitterBalancesAfter[1], 0);
    assertEq(transmitterBalancesAfter[2], 0);
    assertEq(transmitterBalancesAfter[3], 0);
  }

  function test_OracleWithdraw_SuccessCoordinatorOwner() public {
    // Send as Coordinator Owner
    address coordinatorOwner = s_functionsCoordinator.owner();
    vm.stopPrank();
    vm.startPrank(coordinatorOwner);

    uint256 coordinatorOwnerBalanceBefore = s_linkToken.balanceOf(coordinatorOwner);

    // Attempt to withdraw with no amount, which will withdraw the full balance
    s_functionsCoordinator.oracleWithdraw(coordinatorOwner, 0);

    // 4 report transmissions have been made
    uint96 totalOperationFees = s_functionsCoordinator.getOperationFeeJuels() * s_requestsFulfilled;

    uint256 coordinatorOwnerBalanceAfter = s_linkToken.balanceOf(coordinatorOwner);
    assertEq(coordinatorOwnerBalanceBefore + totalOperationFees, coordinatorOwnerBalanceAfter);
  }
}

/// @notice #oracleWithdrawAll
contract FunctionsBilling_OracleWithdrawAll is FunctionsMultipleFulfillmentsSetup {
  function setUp() public virtual override {
    FunctionsMultipleFulfillmentsSetup.setUp();
  }

  function test_OracleWithdrawAll_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsCoordinator.oracleWithdrawAll();
  }

  function test_OracleWithdrawAll_SuccessPaysTransmittersWithBalance() public {
    uint256[4] memory transmitterBalancesBefore = _getTransmitterBalances();
    _assertTransmittersAllHaveBalance(transmitterBalancesBefore, 0);

    s_functionsCoordinator.oracleWithdrawAll();

    uint96 totalOperationFees = s_functionsCoordinator.getOperationFeeJuels() * s_requestsFulfilled;
    uint96 totalDonFees = s_functionsCoordinator.getDONFeeJuels(new bytes(0)) * s_requestsFulfilled;
    uint96 donFeeShare = totalDonFees / uint8(s_transmitters.length);
    uint96 expectedBalancePerFulfillment = ((s_fulfillmentCoordinatorBalance - totalOperationFees - totalDonFees) /
      s_requestsFulfilled);

    uint256[4] memory transmitterBalancesAfter = _getTransmitterBalances();
    // Transmitter 1 has transmitted twice
    assertEq(transmitterBalancesAfter[0], (expectedBalancePerFulfillment * 2) + donFeeShare);
    // Transmitter 2 and 3 have transmitted once
    assertEq(transmitterBalancesAfter[1], expectedBalancePerFulfillment + donFeeShare);
    assertEq(transmitterBalancesAfter[2], expectedBalancePerFulfillment + donFeeShare);
    // Transmitter 4 only not transmitted, it only has its share of the DON fees
    assertEq(transmitterBalancesAfter[3], donFeeShare);
  }
}

/// @notice #_disperseFeePool
contract FunctionsBilling__DisperseFeePool is FunctionsRouterSetup {
  function test__DisperseFeePool_RevertIfNotSet() public {
    // Manually set s_feePool (at slot 12) to 1 to get past first check in _disperseFeePool
    vm.store(address(s_functionsCoordinator), bytes32(uint256(12)), bytes32(uint256(1)));

    vm.expectRevert(FunctionsBilling.NoTransmittersSet.selector);
    s_functionsCoordinator.disperseFeePool_HARNESS();
  }
}
