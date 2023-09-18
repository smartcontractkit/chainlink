// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsRouter} from "../../dev/v1_0_0/FunctionsRouter.sol";
import {FunctionsSubscriptions} from "../../dev/v1_0_0/FunctionsSubscriptions.sol";
import {FunctionsCoordinator} from "../../dev/v1_0_0/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/v1_0_0/FunctionsBilling.sol";
import {FunctionsRequest} from "../../dev/v1_0_0/libraries/FunctionsRequest.sol";
import {FunctionsResponse} from "../../dev/v1_0_0/libraries/FunctionsResponse.sol";
import {FunctionsCoordinatorTestHelper} from "./testhelpers/FunctionsCoordinatorTestHelper.sol";
import {FunctionsClientTestHelper} from "./testhelpers/FunctionsClientTestHelper.sol";

import {FunctionsRouterSetup, FunctionsRoutesSetup, FunctionsSubscriptionSetup, FunctionsClientRequestSetup} from "./Setup.t.sol";

import "forge-std/Vm.sol";

// ================================================================
// |                        Functions Router                      |
// ================================================================

/// @notice #constructor
contract FunctionsRouter_Constructor is FunctionsRouterSetup {
  function test_Constructor_Success() public {
    assertEq(s_functionsRouter.typeAndVersion(), "Functions Router v1.0.0");
    assertEq(s_functionsRouter.owner(), OWNER_ADDRESS);
  }
}

/// @notice #getConfig
contract FunctionsRouter_GetConfig is FunctionsRouterSetup {
  function test_GetConfig_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    FunctionsRouter.Config memory config = s_functionsRouter.getConfig();
    assertEq(config.maxConsumersPerSubscription, getRouterConfig().maxConsumersPerSubscription);
    assertEq(config.adminFee, getRouterConfig().adminFee);
    assertEq(config.handleOracleFulfillmentSelector, getRouterConfig().handleOracleFulfillmentSelector);
    assertEq(config.maxCallbackGasLimits[0], getRouterConfig().maxCallbackGasLimits[0]);
    assertEq(config.maxCallbackGasLimits[1], getRouterConfig().maxCallbackGasLimits[1]);
    assertEq(config.maxCallbackGasLimits[2], getRouterConfig().maxCallbackGasLimits[2]);
    assertEq(config.gasForCallExactCheck, getRouterConfig().gasForCallExactCheck);
    assertEq(config.subscriptionDepositMinimumRequests, getRouterConfig().subscriptionDepositMinimumRequests);
    assertEq(config.subscriptionDepositJuels, getRouterConfig().subscriptionDepositJuels);
  }
}

/// @notice #updateConfig
contract FunctionsRouter_UpdateConfig is FunctionsRouterSetup {
  FunctionsRouter.Config internal configToSet;

  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();

    uint32[] memory maxCallbackGasLimits = new uint32[](4);
    maxCallbackGasLimits[0] = 300_000;
    maxCallbackGasLimits[1] = 500_000;
    maxCallbackGasLimits[2] = 1_000_000;
    maxCallbackGasLimits[3] = 3_000_000;

    configToSet = FunctionsRouter.Config({
      maxConsumersPerSubscription: s_maxConsumersPerSubscription,
      adminFee: s_adminFee,
      handleOracleFulfillmentSelector: s_handleOracleFulfillmentSelector,
      maxCallbackGasLimits: maxCallbackGasLimits,
      gasForCallExactCheck: 5000,
      subscriptionDepositMinimumRequests: 10,
      subscriptionDepositJuels: 5 * 1e18
    });
  }

  function test_UpdateConfig_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.updateConfig(configToSet);
  }

  event ConfigUpdated(FunctionsRouter.Config config);

  function test_UpdateConfig_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit ConfigUpdated(configToSet);

    s_functionsRouter.updateConfig(configToSet);

    FunctionsRouter.Config memory config = s_functionsRouter.getConfig();
    assertEq(config.maxConsumersPerSubscription, configToSet.maxConsumersPerSubscription);
    assertEq(config.adminFee, configToSet.adminFee);
    assertEq(config.handleOracleFulfillmentSelector, configToSet.handleOracleFulfillmentSelector);
    assertEq(config.maxCallbackGasLimits[0], configToSet.maxCallbackGasLimits[0]);
    assertEq(config.maxCallbackGasLimits[1], configToSet.maxCallbackGasLimits[1]);
    assertEq(config.maxCallbackGasLimits[2], configToSet.maxCallbackGasLimits[2]);
    assertEq(config.maxCallbackGasLimits[3], configToSet.maxCallbackGasLimits[3]);
    assertEq(config.gasForCallExactCheck, configToSet.gasForCallExactCheck);
  }
}

/// @notice #isValidCallbackGasLimit
contract FunctionsRouter_IsValidCallbackGasLimit is FunctionsSubscriptionSetup {
  function test_IsValidCallbackGasLimit_RevertInvalidConfig() public {
    // Set an invalid maxCallbackGasLimit flag
    bytes32 flagsToSet = 0x5a00000000000000000000000000000000000000000000000000000000000000;
    s_functionsRouter.setFlags(s_subscriptionId, flagsToSet);

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.InvalidGasFlagValue.selector, 90));
    s_functionsRouter.isValidCallbackGasLimit(s_subscriptionId, 0);
  }

  function test_IsValidCallbackGasLimit_RevertGasLimitTooBig() public {
    uint8 MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX = 0;
    bytes32 subscriptionFlags = s_functionsRouter.getFlags(s_subscriptionId);
    uint8 callbackGasLimitsIndexSelector = uint8(subscriptionFlags[MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX]);

    FunctionsRouter.Config memory config = s_functionsRouter.getConfig();
    uint32[] memory _maxCallbackGasLimits = config.maxCallbackGasLimits;
    uint32 maxCallbackGasLimit = _maxCallbackGasLimits[callbackGasLimitsIndexSelector];

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.GasLimitTooBig.selector, maxCallbackGasLimit));
    s_functionsRouter.isValidCallbackGasLimit(s_subscriptionId, maxCallbackGasLimit + 1);
  }

  function test_IsValidCallbackGasLimit_Success() public view {
    uint8 MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX = 0;
    bytes32 subscriptionFlags = s_functionsRouter.getFlags(s_subscriptionId);
    uint8 callbackGasLimitsIndexSelector = uint8(subscriptionFlags[MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX]);

    FunctionsRouter.Config memory config = s_functionsRouter.getConfig();
    uint32[] memory _maxCallbackGasLimits = config.maxCallbackGasLimits;
    uint32 maxCallbackGasLimit = _maxCallbackGasLimits[callbackGasLimitsIndexSelector];

    s_functionsRouter.isValidCallbackGasLimit(s_subscriptionId, maxCallbackGasLimit);
  }
}

/// @notice #getAdminFee
contract FunctionsRouter_GetAdminFee is FunctionsRouterSetup {
  function test_GetAdminFee_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    uint72 adminFee = s_functionsRouter.getAdminFee();
    assertEq(adminFee, getRouterConfig().adminFee);
  }
}

/// @notice #getAllowListId
contract FunctionsRouter_GetAllowListId is FunctionsRouterSetup {
  function test_GetAllowListId_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 defaultAllowListId = bytes32(0);

    bytes32 allowListId = s_functionsRouter.getAllowListId();
    assertEq(allowListId, defaultAllowListId);
  }
}

/// @notice #setAllowListId
contract FunctionsRouter_SetAllowListId is FunctionsRouterSetup {
  function test_UpdateConfig_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 routeIdToSet = bytes32("allowList");

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.setAllowListId(routeIdToSet);
  }

  function test_SetAllowListId_Success() public {
    bytes32 routeIdToSet = bytes32("allowList");
    s_functionsRouter.setAllowListId(routeIdToSet);
    bytes32 allowListId = s_functionsRouter.getAllowListId();
    assertEq(allowListId, routeIdToSet);
  }
}

/// @notice #_getMaxConsumers
contract FunctionsRouter__GetMaxConsumers is FunctionsRouterSetup {
  // TODO: make contract internal function helper
}

/// @notice #sendRequest
contract FunctionsRouter_SendRequest is FunctionsSubscriptionSetup {
  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Add sending wallet as a subscription consumer
    s_functionsRouter.addConsumer(s_subscriptionId, OWNER_ADDRESS);
  }

  function test_SendRequest_RevertIfInvalidDonId() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    bytes32 invalidDonId = bytes32("this does not exist");

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.RouteNotFound.selector, invalidDonId));
    s_functionsRouter.sendRequest(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5_000,
      invalidDonId
    );
  }

  function test_SendRequest_RevertIfIncorrectDonId() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    bytes32 incorrectDonId = s_functionsRouter.getAllowListId();

    // Low level revert from incorrect call
    vm.expectRevert();
    s_functionsRouter.sendRequest(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5_000,
      incorrectDonId
    );
  }

  function test_SendRequest_RevertIfPaused() public {
    s_functionsRouter.pause();

    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.sendRequest(s_subscriptionId, requestData, FunctionsRequest.REQUEST_DATA_VERSION, 5000, s_donId);
  }

  function test_SendRequest_RevertIfNoSubscription() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint64 invalidSubscriptionId = 123456789;

    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    s_functionsRouter.sendRequest(
      invalidSubscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5000,
      s_donId
    );
  }

  function test_SendRequest_RevertIfConsumerNotAllowed() public {
    // Remove sending wallet as a subscription consumer
    s_functionsRouter.removeConsumer(s_subscriptionId, OWNER_ADDRESS);

    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    vm.expectRevert(FunctionsSubscriptions.InvalidConsumer.selector);
    s_functionsRouter.sendRequest(s_subscriptionId, requestData, FunctionsRequest.REQUEST_DATA_VERSION, 5000, s_donId);
  }

  function test_SendRequest_RevertIfInvalidCallbackGasLimit() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint8 MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX = 0;
    bytes32 subscriptionFlags = s_functionsRouter.getFlags(s_subscriptionId);
    uint8 callbackGasLimitsIndexSelector = uint8(subscriptionFlags[MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX]);

    FunctionsRouter.Config memory config = s_functionsRouter.getConfig();
    uint32[] memory _maxCallbackGasLimits = config.maxCallbackGasLimits;
    uint32 maxCallbackGasLimit = _maxCallbackGasLimits[callbackGasLimitsIndexSelector];

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.GasLimitTooBig.selector, maxCallbackGasLimit));
    s_functionsRouter.sendRequest(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      500_000,
      s_donId
    );
  }

  function test_SendRequest_RevertIfEmptyData() public {
    // Build invalid request data
    bytes memory emptyRequestData = new bytes(0);

    vm.expectRevert(FunctionsRouter.EmptyRequestData.selector);
    s_functionsRouter.sendRequest(
      s_subscriptionId,
      emptyRequestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5_000,
      s_donId
    );
  }

  function test_SendRequest_RevertIfInsufficientSubscriptionBalance() public {
    // Create new subscription that does not have any funding
    uint64 subscriptionId = s_functionsRouter.createSubscription();
    s_functionsRouter.addConsumer(subscriptionId, address(OWNER_ADDRESS));

    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint32 callbackGasLimit = 5000;
    vm.expectRevert(FunctionsBilling.InsufficientBalance.selector);

    s_functionsRouter.sendRequest(
      subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      s_donId
    );
  }

  function test_SendRequest_RevertIfDuplicateRequestId() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    uint32 callbackGasLimit = 5_000;
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    // Send a first request that will remain pending
    bytes32 requestId = s_functionsRouter.sendRequest(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      s_donId
    );

    // Mock the Coordinator to always give back the first requestId
    FunctionsResponse.Commitment memory mockCommitment = FunctionsResponse.Commitment({
      adminFee: s_adminFee,
      coordinator: address(s_functionsCoordinator),
      client: OWNER_ADDRESS,
      subscriptionId: s_subscriptionId,
      callbackGasLimit: callbackGasLimit,
      estimatedTotalCostJuels: 0,
      timeoutTimestamp: uint32(block.timestamp + getCoordinatorConfig().requestTimeoutSeconds),
      requestId: requestId,
      donFee: s_donFee,
      gasOverheadBeforeCallback: getCoordinatorConfig().gasOverheadBeforeCallback,
      gasOverheadAfterCallback: getCoordinatorConfig().gasOverheadAfterCallback
    });

    vm.mockCall(
      address(s_functionsCoordinator),
      abi.encodeWithSelector(FunctionsCoordinator.startRequest.selector),
      abi.encode(mockCommitment)
    );

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.DuplicateRequestId.selector, requestId));
    s_functionsRouter.sendRequest(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      s_donId
    );
  }

  event RequestStart(
    bytes32 indexed requestId,
    bytes32 indexed donId,
    uint64 indexed subscriptionId,
    address subscriptionOwner,
    address requestingContract,
    address requestInitiator,
    bytes data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    uint96 estimatedTotalCostJuels
  );

  function test_SendRequest_Success() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint32 callbackGasLimit = 5000;

    bytes32 expectedRequestId = keccak256(
      abi.encode(address(s_functionsCoordinator), OWNER_ADDRESS, s_subscriptionId, 1)
    );

    uint96 costEstimate = s_functionsCoordinator.estimateCost(
      s_subscriptionId,
      requestData,
      callbackGasLimit,
      tx.gasprice
    );

    vm.recordLogs();

    // topic0 (function signature, always checked), topic1 (true), topic2 (true), topic3 (true), and data (true).
    bool checkTopic1 = true;
    bool checkTopic2 = true;
    bool checkTopic3 = true;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit RequestStart({
      requestId: expectedRequestId,
      donId: s_donId,
      subscriptionId: s_subscriptionId,
      subscriptionOwner: OWNER_ADDRESS,
      requestingContract: OWNER_ADDRESS,
      requestInitiator: OWNER_ADDRESS,
      data: requestData,
      dataVersion: FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit: callbackGasLimit,
      estimatedTotalCostJuels: costEstimate
    });

    bytes32 requestIdFromReturn = s_functionsRouter.sendRequest(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      s_donId
    );

    // Get requestId from RequestStart event log topic 1
    Vm.Log[] memory entries = vm.getRecordedLogs();
    bytes32 requestIdFromEvent = entries[2].topics[1];

    assertEq(requestIdFromReturn, requestIdFromEvent);
  }
}

/// @notice #sendRequestToProposed
contract FunctionsRouter_SendRequestToProposed is FunctionsSubscriptionSetup {
  FunctionsCoordinatorTestHelper internal s_functionsCoordinator2; // TODO: use actual FunctionsCoordinator instead of helper

  function setUp() public virtual override {
    FunctionsSubscriptionSetup.setUp();

    // Add sending wallet as a subscription consumer
    s_functionsRouter.addConsumer(s_subscriptionId, OWNER_ADDRESS);

    // Deploy new Coordinator contract
    s_functionsCoordinator2 = new FunctionsCoordinatorTestHelper(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed)
    );

    // Propose new Coordinator contract
    bytes32[] memory proposedContractSetIds = new bytes32[](1);
    proposedContractSetIds[0] = s_donId;
    address[] memory proposedContractSetAddresses = new address[](1);
    proposedContractSetAddresses[0] = address(s_functionsCoordinator2);

    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
  }

  function test_SendRequestToProposed_RevertIfInvalidDonId() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    bytes32 invalidDonId = bytes32("this does not exist");

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.RouteNotFound.selector, invalidDonId));
    s_functionsRouter.sendRequestToProposed(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5_000,
      invalidDonId
    );
  }

  function test_SendRequestToProposed_RevertIfIncorrectDonId() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    bytes32 incorrectDonId = s_functionsRouter.getAllowListId();

    // Low level revert from incorrect call
    vm.expectRevert();
    s_functionsRouter.sendRequestToProposed(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5_000,
      incorrectDonId
    );
  }

  function test_SendRequestToProposed_RevertIfPaused() public {
    s_functionsRouter.pause();

    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.sendRequestToProposed(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5000,
      s_donId
    );
  }

  function test_SendRequestToProposed_RevertIfNoSubscription() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint64 invalidSubscriptionId = 123456789;

    vm.expectRevert(FunctionsSubscriptions.InvalidSubscription.selector);
    s_functionsRouter.sendRequestToProposed(
      invalidSubscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5000,
      s_donId
    );
  }

  function test_SendRequestToProposed_RevertIfConsumerNotAllowed() public {
    // Remove sending wallet as a subscription consumer
    s_functionsRouter.removeConsumer(s_subscriptionId, OWNER_ADDRESS);

    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    vm.expectRevert(FunctionsSubscriptions.InvalidConsumer.selector);
    s_functionsRouter.sendRequestToProposed(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5000,
      s_donId
    );
  }

  function test_SendRequestToProposed_RevertIfInvalidCallbackGasLimit() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint8 MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX = 0;
    bytes32 subscriptionFlags = s_functionsRouter.getFlags(s_subscriptionId);
    uint8 callbackGasLimitsIndexSelector = uint8(subscriptionFlags[MAX_CALLBACK_GAS_LIMIT_FLAGS_INDEX]);

    FunctionsRouter.Config memory config = s_functionsRouter.getConfig();
    uint32[] memory _maxCallbackGasLimits = config.maxCallbackGasLimits;
    uint32 maxCallbackGasLimit = _maxCallbackGasLimits[callbackGasLimitsIndexSelector];

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.GasLimitTooBig.selector, maxCallbackGasLimit));
    s_functionsRouter.sendRequestToProposed(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      500_000,
      s_donId
    );
  }

  function test_SendRequestToProposed_RevertIfEmptyData() public {
    // Build invalid request data
    bytes memory emptyRequestData = new bytes(0);

    vm.expectRevert(FunctionsRouter.EmptyRequestData.selector);
    s_functionsRouter.sendRequestToProposed(
      s_subscriptionId,
      emptyRequestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      5_000,
      s_donId
    );
  }

  function test_SendRequest_RevertIfInsufficientSubscriptionBalance() public {
    // Create new subscription that does not have any funding
    uint64 subscriptionId = s_functionsRouter.createSubscription();
    s_functionsRouter.addConsumer(subscriptionId, address(OWNER_ADDRESS));

    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint32 callbackGasLimit = 5000;
    vm.expectRevert(FunctionsBilling.InsufficientBalance.selector);

    s_functionsRouter.sendRequestToProposed(
      subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      s_donId
    );
  }

  event RequestStart(
    bytes32 indexed requestId,
    bytes32 indexed donId,
    uint64 indexed subscriptionId,
    address subscriptionOwner,
    address requestingContract,
    address requestInitiator,
    bytes data,
    uint16 dataVersion,
    uint32 callbackGasLimit,
    uint96 estimatedTotalCostJuels
  );

  function test_SendRequestToProposed_Success() public {
    // Build minimal valid request data
    string memory sourceCode = "return 'hello world';";
    FunctionsRequest.Request memory request;
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint32 callbackGasLimit = 5000;

    bytes32 expectedRequestId = keccak256(
      abi.encode(address(s_functionsCoordinator2), OWNER_ADDRESS, s_subscriptionId, 1)
    );

    uint96 costEstimate = s_functionsCoordinator2.estimateCost(
      s_subscriptionId,
      requestData,
      callbackGasLimit,
      tx.gasprice
    );

    vm.recordLogs();

    // topic0 (function signature, always checked), topic1 (true), topic2 (true), topic3 (true), and data (true).
    bool checkTopic1 = true;
    bool checkTopic2 = true;
    bool checkTopic3 = true;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit RequestStart({
      requestId: expectedRequestId,
      donId: s_donId,
      subscriptionId: s_subscriptionId,
      subscriptionOwner: OWNER_ADDRESS,
      requestingContract: OWNER_ADDRESS,
      requestInitiator: OWNER_ADDRESS,
      data: requestData,
      dataVersion: FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit: callbackGasLimit,
      estimatedTotalCostJuels: costEstimate
    });

    bytes32 requestIdFromReturn = s_functionsRouter.sendRequestToProposed(
      s_subscriptionId,
      requestData,
      FunctionsRequest.REQUEST_DATA_VERSION,
      callbackGasLimit,
      s_donId
    );

    // Get requestId from RequestStart event log topic 1
    Vm.Log[] memory entries = vm.getRecordedLogs();
    bytes32 requestIdFromEvent = entries[2].topics[1];

    assertEq(requestIdFromReturn, requestIdFromEvent);
  }
}

/// @notice #_sendRequest
contract FunctionsRouter__SendRequest is FunctionsRouterSetup {
  // TODO: make contract internal function helper
}

/// @notice #fulfill
contract FunctionsRouter_Fulfill is FunctionsClientRequestSetup {
  function test_Fulfill_RevertIfPaused() public {
    s_functionsRouter.pause();

    uint256 requestToFulfill = 1;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;

    string[] memory results = new string[](1);
    string memory response = "hello world!";
    results[0] = response;

    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    vm.expectRevert("Pausable: paused");
    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, false);
  }

  function test_Fulfill_RevertIfNotCommittedCoordinator() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes memory response = bytes("hello world!");
    bytes memory err = new bytes(0);
    uint96 juelsPerGas = 0;
    uint96 costWithoutCallback = 0;
    address transmitter = NOP_TRANSMITTER_ADDRESS_1;
    FunctionsResponse.Commitment memory commitment = s_requests[1].commitment;

    vm.expectRevert(FunctionsRouter.OnlyCallableFromCoordinator.selector);
    s_functionsRouter.fulfill(response, err, juelsPerGas, costWithoutCallback, transmitter, commitment);
  }

  event RequestNotProcessed(
    bytes32 indexed requestId,
    address coordinator,
    address transmitter,
    FunctionsResponse.FulfillResult resultCode
  );

  function test_Fulfill_RequestNotProcessedInvalidRequestId() public {
    // Send as committed Coordinator
    vm.stopPrank();
    vm.startPrank(address(s_functionsCoordinator));

    bytes memory response = bytes("hello world!");
    bytes memory err = new bytes(0);
    uint96 juelsPerGas = 0;
    uint96 costWithoutCallback = 0;
    address transmitter = NOP_TRANSMITTER_ADDRESS_1;
    FunctionsResponse.Commitment memory commitment = s_requests[1].commitment;
    // Modify request commitment to have a invalid requestId
    bytes32 invalidRequestId = bytes32("this does not exist");
    commitment.requestId = invalidRequestId;

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit RequestNotProcessed({
      requestId: s_requests[1].requestId,
      coordinator: address(s_functionsCoordinator),
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.INVALID_REQUEST_ID
    });

    (FunctionsResponse.FulfillResult resultCode, uint96 callbackGasCostJuels) = s_functionsRouter.fulfill(
      response,
      err,
      juelsPerGas,
      costWithoutCallback,
      transmitter,
      commitment
    );

    assertEq(uint(resultCode), uint(FunctionsResponse.FulfillResult.INVALID_REQUEST_ID));
    assertEq(callbackGasCostJuels, 0);
  }

  function test_Fulfill_RequestNotProcessedInvalidCommitment() public {
    // Send as committed Coordinator
    vm.stopPrank();
    vm.startPrank(address(s_functionsCoordinator));

    bytes memory response = bytes("hello world!");
    bytes memory err = new bytes(0);
    uint96 juelsPerGas = 0;
    uint96 costWithoutCallback = 0;
    address transmitter = NOP_TRANSMITTER_ADDRESS_1;
    FunctionsResponse.Commitment memory commitment = s_requests[1].commitment;
    // Modify request commitment to have charge more than quoted
    commitment.estimatedTotalCostJuels = 10 * JUELS_PER_LINK; // 10 LINK

    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit RequestNotProcessed({
      requestId: s_requests[1].requestId,
      coordinator: address(s_functionsCoordinator),
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.INVALID_COMMITMENT
    });

    (FunctionsResponse.FulfillResult resultCode, uint96 callbackGasCostJuels) = s_functionsRouter.fulfill(
      response,
      err,
      juelsPerGas,
      costWithoutCallback,
      transmitter,
      commitment
    );

    assertEq(uint(resultCode), uint(FunctionsResponse.FulfillResult.INVALID_COMMITMENT));
    assertEq(callbackGasCostJuels, 0);
  }

  function test_Fulfill_RequestNotProcessedInsufficientGas() public {
    uint256 requestToFulfill = 1;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;

    string[] memory results = new string[](1);
    string memory response = "hello world!";
    results[0] = response;

    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    uint32 callbackGasLimit = s_requests[requestToFulfill].requestData.callbackGasLimit;
    // Coordinator sends enough gas that would get through callback and payment, but fail after
    uint256 gasToUse = getCoordinatorConfig().gasOverheadBeforeCallback + callbackGasLimit + 100000;

    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1RequestId = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1RequestId, checkTopic2, checkTopic3, checkData);
    emit RequestNotProcessed({
      requestId: s_requests[requestToFulfill].requestId,
      coordinator: address(s_functionsCoordinator),
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.INSUFFICIENT_GAS_PROVIDED
    });

    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, false, 1, gasToUse);
  }

  function test_Fulfill_RequestNotProcessedSubscriptionBalanceInvariant() public {
    // Find the storage slot that the Subscription is on
    vm.record();
    s_functionsRouter.getSubscription(s_subscriptionId);
    (bytes32[] memory reads, ) = vm.accesses(address(s_functionsRouter));
    // The first read is from '_isExistingSubscription' which checks Subscription.owner on slot 0
    // Slot 0 is shared with the Subscription.balance
    uint256 slot = uint256(reads[0]);

    // The request has already been initiated, forcibly lower the subscription's balance by clearing out slot 0
    uint96 balance = 1;
    address owner = address(0);
    bytes32 data = bytes32(abi.encodePacked(balance, owner)); // TODO: make this more accurate
    vm.store(address(s_functionsRouter), bytes32(uint256(slot)), data);

    uint256 requestToFulfill = 1;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;

    string[] memory results = new string[](1);
    string memory response = "hello world!";
    results[0] = response;

    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1RequestId = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1RequestId, checkTopic2, checkTopic3, checkData);
    emit RequestNotProcessed({
      requestId: s_requests[requestToFulfill].requestId,
      coordinator: address(s_functionsCoordinator),
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.SUBSCRIPTION_BALANCE_INVARIANT_VIOLATION
    });

    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, false);
  }

  function test_Fulfill_RequestNotProcessedCostExceedsCommitment() public {
    // Use higher juelsPerGas than request time
    // 10x the gas price
    vm.txGasPrice(TX_GASPRICE_START * 10);

    uint256 requestToFulfill = 1;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;

    string[] memory results = new string[](1);
    string memory response = "hello world!";
    results[0] = response;

    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1RequestId = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1RequestId, checkTopic2, checkTopic3, checkData);
    emit RequestNotProcessed({
      requestId: s_requests[requestToFulfill].requestId,
      coordinator: address(s_functionsCoordinator),
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.COST_EXCEEDS_COMMITMENT
    });

    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, false);
  }

  event RequestProcessed(
    bytes32 indexed requestId,
    uint64 indexed subscriptionId,
    uint96 totalCostJuels,
    address transmitter,
    FunctionsResponse.FulfillResult resultCode,
    bytes response,
    bytes err,
    bytes callbackReturnData
  );

  FunctionsClientTestHelper internal s_clientWithFailingCallback;

  function test_Fulfill_SuccessUserCallbackReverts() public {
    // Deploy Client with failing callback
    s_clientWithFailingCallback = new FunctionsClientTestHelper(address(s_functionsRouter));
    s_clientWithFailingCallback.setRevertFulfillRequest(true);

    // Add Client as a subscription consumer
    s_functionsRouter.addConsumer(s_subscriptionId, address(s_clientWithFailingCallback));

    // Send a minimal request
    uint256 requestKey = 99;

    string memory sourceCode = "return 'hello world';";
    uint32 callbackGasLimit = 5500;

    vm.recordLogs();
    bytes32 requestId = s_clientWithFailingCallback.sendSimpleRequestWithJavaScript(
      sourceCode,
      s_subscriptionId,
      s_donId,
      callbackGasLimit
    );

    // Get commitment data from OracleRequest event log
    Vm.Log[] memory entries = vm.getRecordedLogs();
    (, , , , , , , FunctionsResponse.Commitment memory _commitment) = abi.decode(
      entries[0].data,
      (address, uint64, address, bytes, uint16, bytes32, uint64, FunctionsResponse.Commitment)
    );

    s_requests[requestKey] = Request({
      requestData: RequestData({
        sourceCode: sourceCode,
        secrets: new bytes(0),
        args: new string[](0),
        bytesArgs: new bytes[](0),
        callbackGasLimit: callbackGasLimit
      }),
      requestId: requestId,
      commitment: _commitment
    });

    // Fulfill
    uint256 requestToFulfill = requestKey;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;

    string[] memory results = new string[](1);
    string memory response = "hello world";
    results[0] = response;

    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1RequestId = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1RequestId, checkTopic2, checkTopic3, checkData);
    emit RequestProcessed({
      requestId: requestId,
      subscriptionId: s_subscriptionId,
      totalCostJuels: _getExpectedCost(1379), // gasUsed is manually taken
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.USER_CALLBACK_ERROR,
      response: bytes(response),
      err: err,
      callbackReturnData: vm.parseBytes(
        "0x08c379a00000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000f61736b656420746f207265766572740000000000000000000000000000000000"
      ) // TODO: build this programatically
    });

    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, true, 1);
  }

  function test_Fulfill_SuccessUserCallbackRunsOutOfGas() public {
    // Send request #2 with no callback gas
    string memory sourceCode = "return 'hello world';";
    bytes memory secrets = new bytes(0);
    string[] memory args = new string[](0);
    bytes[] memory bytesArgs = new bytes[](0);
    uint32 callbackGasLimit = 0;
    _sendAndStoreRequest(2, sourceCode, secrets, args, bytesArgs, callbackGasLimit);

    uint256 requestToFulfill = 2;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;

    string[] memory results = new string[](1);
    string memory response = "hello world!";
    results[0] = response;

    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    // topic0 (function signature, always checked), topic1: request ID(true), NOT topic2 (false), NOT topic3 (false), and data (true).
    vm.expectEmit(true, false, false, true);
    emit RequestProcessed({
      requestId: s_requests[requestToFulfill].requestId,
      subscriptionId: s_subscriptionId,
      totalCostJuels: _getExpectedCost(137), // gasUsed is manually taken
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.USER_CALLBACK_ERROR,
      response: bytes(response),
      err: err,
      callbackReturnData: new bytes(0)
    });

    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, true, 1);
  }

  function test_Fulfill_SuccessClientNoLongerExists() public {
    // Delete the Client contract in the time between request and fulfillment
    vm.etch(address(s_functionsClient), new bytes(0));

    uint256 requestToFulfill = 1;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;

    string[] memory results = new string[](1);
    string memory response = "hello world!";
    results[0] = response;

    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    // topic0 (function signature, always checked), topic1 (true), topic2 (true), NOT topic3 (false), and data (true).
    bool checkTopic1RequestId = true;
    bool checkTopic2SubscriptionId = true;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1RequestId, checkTopic2SubscriptionId, checkTopic3, checkData);
    emit RequestProcessed({
      requestId: s_requests[requestToFulfill].requestId,
      subscriptionId: s_subscriptionId,
      totalCostJuels: _getExpectedCost(0), // gasUsed is manually taken
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.USER_CALLBACK_ERROR,
      response: bytes(response),
      err: err,
      callbackReturnData: new bytes(0)
    });

    _reportAndStore(requestNumberKeys, results, errors, NOP_TRANSMITTER_ADDRESS_1, true, 1);
  }

  function test_Fulfill_SuccessFulfilled() public {
    // Fulfill request 1
    uint256 requestToFulfill = 1;

    uint256[] memory requestNumberKeys = new uint256[](1);
    requestNumberKeys[0] = requestToFulfill;
    string[] memory results = new string[](1);
    string memory response = "hello world!";
    results[0] = response;
    bytes[] memory errors = new bytes[](1);
    bytes memory err = new bytes(0);
    errors[0] = err;

    // topic0 (function signature, always checked), topic1 (true), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1RequestId = true;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1RequestId, checkTopic2, checkTopic3, checkData);
    emit RequestProcessed({
      requestId: s_requests[requestToFulfill].requestId,
      subscriptionId: s_subscriptionId,
      totalCostJuels: _getExpectedCost(5371), // gasUsed is manually taken
      transmitter: NOP_TRANSMITTER_ADDRESS_1,
      resultCode: FunctionsResponse.FulfillResult.FULFILLED,
      response: bytes(response),
      err: err,
      callbackReturnData: new bytes(0)
    });
    _reportAndStore(requestNumberKeys, results, errors);
  }
}

/// @notice #_callback
contract FunctionsRouter__Callback is FunctionsRouterSetup {
  // TODO: make contract internal function helper
}

/// @notice #getContractById
contract FunctionsRouter_GetContractById is FunctionsRoutesSetup {
  function test_GetContractById_RevertIfRouteDoesNotExist() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 invalidRouteId = bytes32("this does not exist");

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.RouteNotFound.selector, invalidRouteId));
    s_functionsRouter.getContractById(invalidRouteId);
  }

  function test_GetContractById_SuccessIfRouteExists() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    address routeDestination = s_functionsRouter.getContractById(s_donId);
    assertEq(routeDestination, address(s_functionsCoordinator));
  }
}

/// @notice #getProposedContractById
contract FunctionsRouter_GetProposedContractById is FunctionsRoutesSetup {
  FunctionsCoordinatorTestHelper internal s_functionsCoordinator2; // TODO: use actual FunctionsCoordinator instead of helper

  function setUp() public virtual override {
    FunctionsRoutesSetup.setUp();

    // Deploy new Coordinator contract
    s_functionsCoordinator2 = new FunctionsCoordinatorTestHelper(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed)
    );

    // Propose new Coordinator contract
    bytes32[] memory proposedContractSetIds = new bytes32[](1);
    proposedContractSetIds[0] = s_donId;
    address[] memory proposedContractSetAddresses = new address[](1);
    proposedContractSetAddresses[0] = address(s_functionsCoordinator2);

    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
  }

  function test_GetProposedContractById_RevertIfRouteDoesNotExist() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes32 invalidRouteId = bytes32("this does not exist");

    vm.expectRevert(abi.encodeWithSelector(FunctionsRouter.RouteNotFound.selector, invalidRouteId));
    s_functionsRouter.getProposedContractById(invalidRouteId);
  }

  function test_GetProposedContractById_SuccessIfRouteExists() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    address routeDestination = s_functionsRouter.getProposedContractById(s_donId);
    assertEq(routeDestination, address(s_functionsCoordinator2));
  }
}

/// @notice #getProposedContractSet
contract FunctionsRouter_GetProposedContractSet is FunctionsRoutesSetup {
  FunctionsCoordinatorTestHelper internal s_functionsCoordinator2; // TODO: use actual FunctionsCoordinator instead of helper
  bytes32[] s_proposedContractSetIds;
  address[] s_proposedContractSetAddresses;

  function setUp() public virtual override {
    FunctionsRoutesSetup.setUp();

    // Deploy new Coordinator contract
    s_functionsCoordinator2 = new FunctionsCoordinatorTestHelper(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed)
    );

    // Propose new Coordinator contract
    s_proposedContractSetIds = new bytes32[](1);
    s_proposedContractSetIds[0] = s_donId;
    s_proposedContractSetAddresses = new address[](1);
    s_proposedContractSetAddresses[0] = address(s_functionsCoordinator2);

    s_functionsRouter.proposeContractsUpdate(s_proposedContractSetIds, s_proposedContractSetAddresses);
  }

  function test_GetProposedContractSet_Success() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    (bytes32[] memory proposedContractSetIds, address[] memory proposedContractSetAddresses) = s_functionsRouter
      .getProposedContractSet();

    assertEq(proposedContractSetIds.length, 1);
    assertEq(proposedContractSetIds[0], s_donId);
    assertEq(proposedContractSetIds.length, 1);
    assertEq(proposedContractSetAddresses[0], address(s_functionsCoordinator2));
  }
}

/// @notice #proposeContractsUpdate
contract FunctionsRouter_ProposeContractsUpdate is FunctionsRoutesSetup {
  FunctionsCoordinatorTestHelper internal s_functionsCoordinator2; // TODO: use actual FunctionsCoordinator instead of helper
  bytes32[] s_proposedContractSetIds;
  address[] s_proposedContractSetAddresses;

  function setUp() public virtual override {
    FunctionsRoutesSetup.setUp();

    // Deploy new Coordinator contract
    s_functionsCoordinator2 = new FunctionsCoordinatorTestHelper(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed)
    );

    // Propose new Coordinator contract
    s_proposedContractSetIds = new bytes32[](1);
    s_proposedContractSetIds[0] = s_donId;
    s_proposedContractSetAddresses = new address[](1);
    s_proposedContractSetAddresses[0] = address(s_functionsCoordinator2);
  }

  function test_ProposeContractsUpdate_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.proposeContractsUpdate(s_proposedContractSetIds, s_proposedContractSetAddresses);
  }

  function test_ProposeContractsUpdate_RevertIfLengthMismatch() public {
    bytes32[] memory proposedContractSetIds = new bytes32[](1);
    proposedContractSetIds[0] = s_donId;
    address[] memory proposedContractSetAddresses = new address[](1);

    vm.expectRevert(FunctionsRouter.InvalidProposal.selector);
    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
  }

  function test_ProposeContractsUpdate_RevertIfExceedsMaxProposal() public {
    uint8 MAX_PROPOSAL_SET_LENGTH = 8;
    uint8 INVALID_PROPOSAL_SET_LENGTH = MAX_PROPOSAL_SET_LENGTH + 1;

    // Generate some mock data
    bytes32[] memory proposedContractSetIds = new bytes32[](INVALID_PROPOSAL_SET_LENGTH);
    for (uint256 i = 0; i < INVALID_PROPOSAL_SET_LENGTH; ++i) {
      proposedContractSetIds[i] = bytes32(uint256(i + 111));
    }
    address[] memory proposedContractSetAddresses = new address[](INVALID_PROPOSAL_SET_LENGTH);
    for (uint256 i = 0; i < INVALID_PROPOSAL_SET_LENGTH; ++i) {
      proposedContractSetAddresses[i] = address(uint160(uint(keccak256(abi.encodePacked(i + 111)))));
    }

    vm.expectRevert(FunctionsRouter.InvalidProposal.selector);
    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
  }

  function test_ProposeContractsUpdate_RevertIfEmptyAddress() public {
    bytes32[] memory proposedContractSetIds = new bytes32[](1);
    proposedContractSetIds[0] = s_donId;
    address[] memory proposedContractSetAddresses = new address[](1);
    proposedContractSetAddresses[0] = address(0);

    vm.expectRevert(FunctionsRouter.InvalidProposal.selector);
    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
  }

  function test_ProposeContractsUpdate_RevertIfNotNewContract() public {
    bytes32[] memory proposedContractSetIds = new bytes32[](1);
    proposedContractSetIds[0] = s_donId;
    address[] memory proposedContractSetAddresses = new address[](1);
    proposedContractSetAddresses[0] = address(s_functionsCoordinator);

    vm.expectRevert(FunctionsRouter.InvalidProposal.selector);
    s_functionsRouter.proposeContractsUpdate(proposedContractSetIds, proposedContractSetAddresses);
  }

  event ContractProposed(
    bytes32 proposedContractSetId,
    address proposedContractSetFromAddress,
    address proposedContractSetToAddress
  );

  function test_ProposeContractsUpdate_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit ContractProposed({
      proposedContractSetId: s_proposedContractSetIds[0],
      proposedContractSetFromAddress: address(s_functionsCoordinator),
      proposedContractSetToAddress: s_proposedContractSetAddresses[0]
    });

    s_functionsRouter.proposeContractsUpdate(s_proposedContractSetIds, s_proposedContractSetAddresses);
  }
}

/// @notice #updateContracts
contract FunctionsRouter_UpdateContracts is FunctionsRoutesSetup {
  FunctionsCoordinatorTestHelper internal s_functionsCoordinator2; // TODO: use actual FunctionsCoordinator instead of helper
  bytes32[] s_proposedContractSetIds;
  address[] s_proposedContractSetAddresses;

  function setUp() public virtual override {
    FunctionsRoutesSetup.setUp();

    // Deploy new Coordinator contract
    s_functionsCoordinator2 = new FunctionsCoordinatorTestHelper(
      address(s_functionsRouter),
      getCoordinatorConfig(),
      address(s_linkEthFeed)
    );

    // Propose new Coordinator contract
    s_proposedContractSetIds = new bytes32[](1);
    s_proposedContractSetIds[0] = s_donId;
    s_proposedContractSetAddresses = new address[](1);
    s_proposedContractSetAddresses[0] = address(s_functionsCoordinator2);

    s_functionsRouter.proposeContractsUpdate(s_proposedContractSetIds, s_proposedContractSetAddresses);
  }

  function test_UpdateContracts_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.updateContracts();
  }

  event ContractUpdated(bytes32 id, address from, address to);

  function test_UpdateContracts_Success() public {
    // topic0 (function signature, always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    bool checkTopic1 = false;
    bool checkTopic2 = false;
    bool checkTopic3 = false;
    bool checkData = true;
    vm.expectEmit(checkTopic1, checkTopic2, checkTopic3, checkData);
    emit ContractUpdated({
      id: s_proposedContractSetIds[0],
      from: address(s_functionsCoordinator),
      to: s_proposedContractSetAddresses[0]
    });

    s_functionsRouter.updateContracts();

    (bytes32[] memory proposedContractSetIds, address[] memory proposedContractSetAddresses) = s_functionsRouter
      .getProposedContractSet();

    assertEq(proposedContractSetIds.length, 0);
    assertEq(proposedContractSetAddresses.length, 0);
  }
}

/// @notice #_whenNotPaused
contract FunctionsRouter__WhenNotPaused is FunctionsRouterSetup {
  // TODO: make contract internal function helper
}

/// @notice #_onlyRouterOwner
contract FunctionsRouter__OnlyRouterOwner is FunctionsRouterSetup {
  // TODO: make contract internal function helper
}

/// @notice #_onlySenderThatAcceptedToS
contract FunctionsRouter__OnlySenderThatAcceptedToS is FunctionsRouterSetup {
  // TODO: make contract internal function helper
}

/// @notice #pause
contract FunctionsRouter_Pause is FunctionsRouterSetup {
  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();
  }

  event Paused(address account);

  function test_Pause_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.pause();
  }

  function test_Pause_Success() public {
    // topic0 (always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    vm.expectEmit(false, false, false, true);
    emit Paused(OWNER_ADDRESS);

    s_functionsRouter.pause();

    bool isPaused = s_functionsRouter.paused();
    assertEq(isPaused, true);

    vm.expectRevert("Pausable: paused");
    s_functionsRouter.createSubscription();
  }
}

/// @notice #unpause
contract FunctionsRouter_Unpause is FunctionsRouterSetup {
  function setUp() public virtual override {
    FunctionsRouterSetup.setUp();
    s_functionsRouter.pause();
  }

  event Unpaused(address account);

  function test_Unpause_RevertIfNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsRouter.unpause();
  }

  function test_Unpause_Success() public {
    // topic0 (always checked), NOT topic1 (false), NOT topic2 (false), NOT topic3 (false), and data (true).
    vm.expectEmit(false, false, false, true);
    emit Unpaused(OWNER_ADDRESS);

    s_functionsRouter.unpause();

    bool isPaused = s_functionsRouter.paused();
    assertEq(isPaused, false);

    s_functionsRouter.createSubscription();
  }
}
