// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../dev/v1_X/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/v1_X/FunctionsBilling.sol";
import {FunctionsRequest} from "../../dev/v1_X/libraries/FunctionsRequest.sol";
import {FunctionsSubscriptions} from "../../dev/v1_X/FunctionsSubscriptions.sol";
import {Routable} from "../../dev/v1_X/Routable.sol";

import {FunctionsRouterSetup, FunctionsSubscriptionSetup, FunctionsClientRequestSetup, FunctionsMultipleFulfillmentsSetup} from "./Setup.t.sol";

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

    FunctionsBilling.Config memory config = s_functionsCoordinator.getConfig();
    assertEq(config.feedStalenessSeconds, getCoordinatorConfig().feedStalenessSeconds);
    assertEq(config.gasOverheadBeforeCallback, getCoordinatorConfig().gasOverheadBeforeCallback);
    assertEq(config.gasOverheadAfterCallback, getCoordinatorConfig().gasOverheadAfterCallback);
    assertEq(config.requestTimeoutSeconds, getCoordinatorConfig().requestTimeoutSeconds);
    assertEq(config.donFee, getCoordinatorConfig().donFee);
    assertEq(config.maxSupportedRequestDataVersion, getCoordinatorConfig().maxSupportedRequestDataVersion);
    assertEq(config.fulfillmentGasPriceOverEstimationBP, getCoordinatorConfig().fulfillmentGasPriceOverEstimationBP);
    assertEq(config.fallbackNativePerUnitLink, getCoordinatorConfig().fallbackNativePerUnitLink);
  }
}

/// @notice #updateConfig
contract FunctionsBilling_UpdateConfig is FunctionsRouterSetup {
  FunctionsBilling.Config internal configToSet;

  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();

    configToSet = FunctionsBilling.Config({
      feedStalenessSeconds: getCoordinatorConfig().feedStalenessSeconds * 2,
      gasOverheadAfterCallback: getCoordinatorConfig().gasOverheadAfterCallback * 2,
      gasOverheadBeforeCallback: getCoordinatorConfig().gasOverheadBeforeCallback * 2,
      requestTimeoutSeconds: getCoordinatorConfig().requestTimeoutSeconds * 2,
      donFee: getCoordinatorConfig().donFee * 2,
      maxSupportedRequestDataVersion: getCoordinatorConfig().maxSupportedRequestDataVersion * 2,
      fulfillmentGasPriceOverEstimationBP: getCoordinatorConfig().fulfillmentGasPriceOverEstimationBP * 2,
      fallbackNativePerUnitLink: getCoordinatorConfig().fallbackNativePerUnitLink * 2,
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

  event ConfigUpdated(FunctionsBilling.Config config);

  function test_UpdateConfig_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit ConfigUpdated(configToSet);

    s_functionsCoordinator.updateConfig(configToSet);

    FunctionsBilling.Config memory config = s_functionsCoordinator.getConfig();
    assertEq(config.feedStalenessSeconds, configToSet.feedStalenessSeconds);
    assertEq(config.gasOverheadAfterCallback, configToSet.gasOverheadAfterCallback);
    assertEq(config.gasOverheadBeforeCallback, configToSet.gasOverheadBeforeCallback);
    assertEq(config.requestTimeoutSeconds, configToSet.requestTimeoutSeconds);
    assertEq(config.donFee, configToSet.donFee);
    assertEq(config.maxSupportedRequestDataVersion, configToSet.maxSupportedRequestDataVersion);
    assertEq(config.fulfillmentGasPriceOverEstimationBP, configToSet.fulfillmentGasPriceOverEstimationBP);
    assertEq(config.fallbackNativePerUnitLink, configToSet.fallbackNativePerUnitLink);
    assertEq(config.minimumEstimateGasPriceWei, configToSet.minimumEstimateGasPriceWei);
  }
}

/// @notice #getDONFee
contract FunctionsBilling_GetDONFee is FunctionsRouterSetup {
  function test_GetDONFee_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint72 donFee = s_functionsCoordinator.getDONFee(new bytes(0));
    assertEq(donFee, s_donFee);
  }
}

/// @notice #getAdminFee
contract FunctionsBilling_GetAdminFee is FunctionsRouterSetup {
  function test_GetAdminFee_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint72 adminFee = s_functionsCoordinator.getAdminFee();
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
    uint96 expectedCostEstimate = 51110500000000200;
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
    uint96 expectedCostEstimate = 255552500000000200;
    assertEq(costEstimate, expectedCostEstimate);
  }
}

/// @notice #_calculateCostEstimate
contract FunctionsBilling__CalculateCostEstimate {
  // TODO: make contract internal function helper
}

/// @notice #_startBilling
contract FunctionsBilling__StartBilling {
  // TODO: make contract internal function helper
}

/// @notice #_fulfillAndBill
contract FunctionsBilling__FulfillAndBill {
  // TODO: make contract internal function helper
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
    uint256 transmitter1BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceBefore, 0);
    uint256 transmitter2BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceBefore, 0);
    uint256 transmitter3BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceBefore, 0);
    uint256 transmitter4BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceBefore, 0);

    // Send as stranger, which has no balance
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(FunctionsSubscriptions.InvalidCalldata.selector);

    // Attempt to withdraw with no amount, which would withdraw the full balance
    s_functionsCoordinator.oracleWithdraw(STRANGER_ADDRESS, 0);

    uint96 expectedTransmitterBalance = 0;

    uint256 transmitter1BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceAfter, expectedTransmitterBalance);
    uint256 transmitter2BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceAfter, expectedTransmitterBalance);
    uint256 transmitter3BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceAfter, expectedTransmitterBalance);
    uint256 transmitter4BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceAfter, expectedTransmitterBalance);
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
    uint256 transmitter1BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceBefore, 0);
    uint256 transmitter2BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceBefore, 0);
    uint256 transmitter3BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceBefore, 0);
    uint256 transmitter4BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceBefore, 0);

    // Send as transmitter 1, which has transmitted 1 report
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_1);

    uint96 expectedTransmitterBalance = s_fulfillmentCoordinatorBalance / 3;

    // Attempt to withdraw half of balance
    uint96 halfBalance = expectedTransmitterBalance / 2;
    s_functionsCoordinator.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, halfBalance);

    uint256 transmitter1BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceAfter, halfBalance);
    uint256 transmitter2BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceAfter, 0);
    uint256 transmitter3BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceAfter, 0);
    uint256 transmitter4BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceAfter, 0);
  }

  function test_OracleWithdraw_SuccessTransmitterWithBalanceNoAmountGiven() public {
    uint256 transmitter1BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceBefore, 0);
    uint256 transmitter2BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceBefore, 0);
    uint256 transmitter3BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceBefore, 0);
    uint256 transmitter4BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceBefore, 0);

    // Send as transmitter 1, which has transmitted 1 report
    vm.stopPrank();
    vm.startPrank(NOP_TRANSMITTER_ADDRESS_1);

    // Attempt to withdraw with no amount, which will withdraw the full balance
    s_functionsCoordinator.oracleWithdraw(NOP_TRANSMITTER_ADDRESS_1, 0);

    // 3 report transmissions have been made
    uint96 totalDonFees = s_donFee * 3;
    // 4 transmitters will share the DON fees
    uint96 donFeeShare = totalDonFees / 4;
    uint96 expectedTransmitterBalance = ((s_fulfillmentCoordinatorBalance - totalDonFees) / 3) + donFeeShare;

    uint256 transmitter1BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceAfter, expectedTransmitterBalance);
    uint256 transmitter2BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceAfter, 0);
    uint256 transmitter3BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceAfter, 0);
    uint256 transmitter4BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceAfter, 0);
  }
}

/// @notice #oracleWithdrawAll
contract FunctionsBilling_OracleWithdrawAll is FunctionsMultipleFulfillmentsSetup {
  function setUp() public virtual override {
    // Use no DON fee so that a transmitter has a balance of 0
    s_donFee = 0;

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
    uint256 transmitter1BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceBefore, 0);
    uint256 transmitter2BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceBefore, 0);
    uint256 transmitter3BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceBefore, 0);
    uint256 transmitter4BalanceBefore = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceBefore, 0);

    s_functionsCoordinator.oracleWithdrawAll();

    uint96 expectedTransmitterBalance = s_fulfillmentCoordinatorBalance / 3;

    uint256 transmitter1BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(transmitter1BalanceAfter, expectedTransmitterBalance);
    uint256 transmitter2BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_2);
    assertEq(transmitter2BalanceAfter, expectedTransmitterBalance);
    uint256 transmitter3BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_3);
    assertEq(transmitter3BalanceAfter, expectedTransmitterBalance);
    // Transmitter 4 has no balance
    uint256 transmitter4BalanceAfter = s_linkToken.balanceOf(NOP_TRANSMITTER_ADDRESS_4);
    assertEq(transmitter4BalanceAfter, 0);
  }
}

/// @notice #_disperseFeePool
contract FunctionsBilling__DisperseFeePool is FunctionsRouterSetup {
  function test__DisperseFeePool_RevertIfNotSet() public {
    // Manually set s_feePool (at slot 11) to 1 to get past first check in _disperseFeePool
    vm.store(address(s_functionsCoordinator), bytes32(uint256(11)), bytes32(uint256(1)));

    vm.expectRevert(FunctionsBilling.NoTransmittersSet.selector);
    s_functionsCoordinator.disperseFeePool_HARNESS();
  }
}
