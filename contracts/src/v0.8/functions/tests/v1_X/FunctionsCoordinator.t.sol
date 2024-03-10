// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../dev/v1_X/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/v1_X/FunctionsBilling.sol";
import {FunctionsRequest} from "../../dev/v1_X/libraries/FunctionsRequest.sol";
import {FunctionsResponse} from "../../dev/v1_X/libraries/FunctionsResponse.sol";
import {FunctionsRouter} from "../../dev/v1_X/FunctionsRouter.sol";
import {Routable} from "../../dev/v1_X/Routable.sol";

import {BaseTest} from "./BaseTest.t.sol";
import {FunctionsRouterSetup, FunctionsDONSetup, FunctionsSubscriptionSetup} from "./Setup.t.sol";

/// @notice #constructor
contract FunctionsCoordinator_Constructor is FunctionsRouterSetup {
  function test_Constructor_Success() public {
    assertEq(s_functionsCoordinator.typeAndVersion(), "Functions Coordinator v1.3.0");
    assertEq(s_functionsCoordinator.owner(), OWNER_ADDRESS);
  }
}

/// @notice #getThresholdPublicKey
contract FunctionsCoordinator_GetThresholdPublicKey is FunctionsDONSetup {
  function test_GetThresholdPublicKey_RevertIfEmpty() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    // Reverts when empty
    vm.expectRevert(FunctionsCoordinator.EmptyPublicKey.selector);
    s_functionsCoordinator.getThresholdPublicKey();
  }

  function test_GetThresholdPublicKey_Success() public {
    s_functionsCoordinator.setThresholdPublicKey(s_thresholdKey);

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes memory thresholdKey = s_functionsCoordinator.getThresholdPublicKey();
    assertEq(thresholdKey, s_thresholdKey);
  }
}

/// @notice #setThresholdPublicKey
contract FunctionsCoordinator_SetThresholdPublicKey is FunctionsDONSetup {
  function test_SetThresholdPublicKey_RevertNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    bytes memory newThresholdKey = new bytes(0);
    s_functionsCoordinator.setThresholdPublicKey(newThresholdKey);
  }

  function test_SetThresholdPublicKey_Success() public {
    s_functionsCoordinator.setThresholdPublicKey(s_thresholdKey);

    bytes memory thresholdKey = s_functionsCoordinator.getThresholdPublicKey();

    assertEq(thresholdKey, s_thresholdKey);
  }
}

/// @notice #getDONPublicKey
contract FunctionsCoordinator_GetDONPublicKey is FunctionsDONSetup {
  function test_GetDONPublicKey_RevertIfEmpty() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    // Reverts when empty
    vm.expectRevert(FunctionsCoordinator.EmptyPublicKey.selector);
    s_functionsCoordinator.getDONPublicKey();
  }

  function test_GetDONPublicKey_Success() public {
    s_functionsCoordinator.setDONPublicKey(s_donKey);

    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    bytes memory donKey = s_functionsCoordinator.getDONPublicKey();
    assertEq(donKey, s_donKey);
  }
}

/// @notice #setDONPublicKey
contract FunctionsCoordinator_SetDONPublicKey is FunctionsDONSetup {
  function test_SetDONPublicKey_RevertNotOwner() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert("Only callable by owner");
    s_functionsCoordinator.setDONPublicKey(s_donKey);
  }

  function test_SetDONPublicKey_Success() public {
    s_functionsCoordinator.setDONPublicKey(s_donKey);

    bytes memory donKey = s_functionsCoordinator.getDONPublicKey();
    assertEq(donKey, s_donKey);
  }
}

/// @notice #_isTransmitter
contract FunctionsCoordinator__IsTransmitter is FunctionsDONSetup {
  function test__IsTransmitter_SuccessFound() public {
    bool isTransmitter = s_functionsCoordinator.isTransmitter_HARNESS(NOP_TRANSMITTER_ADDRESS_1);
    assertEq(isTransmitter, true);
  }

  function test__IsTransmitter_SuccessNotFound() public {
    bool isTransmitter = s_functionsCoordinator.isTransmitter_HARNESS(STRANGER_ADDRESS);
    assertEq(isTransmitter, false);
  }
}

/// @notice #startRequest
contract FunctionsCoordinator_StartRequest is FunctionsSubscriptionSetup {
  function test_StartRequest_RevertIfNotRouter() public {
    // Send as stranger
    vm.stopPrank();
    vm.startPrank(STRANGER_ADDRESS);

    vm.expectRevert(Routable.OnlyCallableByRouter.selector);

    s_functionsCoordinator.startRequest(
      FunctionsResponse.RequestMeta({
        requestingContract: address(s_functionsClient),
        data: new bytes(0),
        subscriptionId: s_subscriptionId,
        dataVersion: FunctionsRequest.REQUEST_DATA_VERSION,
        flags: bytes32(0),
        callbackGasLimit: 5_500,
        adminFee: s_adminFee,
        initiatedRequests: 0,
        completedRequests: 0,
        availableBalance: s_subscriptionInitialFunding,
        subscriptionOwner: OWNER_ADDRESS
      })
    );
  }

  event OracleRequest(
    bytes32 indexed requestId,
    address indexed requestingContract,
    address requestInitiator,
    uint64 subscriptionId,
    address subscriptionOwner,
    bytes data,
    uint16 dataVersion,
    bytes32 flags,
    uint64 callbackGasLimit,
    FunctionsResponse.Commitment commitment
  );

  function test_StartRequest_Success() public {
    // Send as Router
    vm.stopPrank();
    vm.startPrank(address(s_functionsRouter));
    (, , address txOrigin) = vm.readCallers();

    bytes memory _requestData = new bytes(0);
    uint32 _callbackGasLimit = 5_500;
    uint96 costEstimate = s_functionsCoordinator.estimateCost(
      s_subscriptionId,
      _requestData,
      _callbackGasLimit,
      tx.gasprice
    );
    uint32 timeoutTimestamp = uint32(block.timestamp + getCoordinatorConfig().requestTimeoutSeconds);
    bytes32 expectedRequestId = keccak256(
      abi.encode(
        address(s_functionsCoordinator),
        address(s_functionsClient),
        s_subscriptionId,
        1,
        keccak256(_requestData),
        FunctionsRequest.REQUEST_DATA_VERSION,
        _callbackGasLimit,
        costEstimate,
        timeoutTimestamp,
        txOrigin
      )
    );

    // WARNING: Kludge in place. Remove in contracts v2.0.0
    FunctionsResponse.Commitment memory expectedComittment = FunctionsResponse.Commitment({
      adminFee: s_functionsCoordinator.getOperationFeeJuels(),
      coordinator: address(s_functionsCoordinator),
      client: address(s_functionsClient),
      subscriptionId: s_subscriptionId,
      callbackGasLimit: _callbackGasLimit,
      estimatedTotalCostJuels: costEstimate,
      timeoutTimestamp: timeoutTimestamp,
      requestId: expectedRequestId,
      donFee: s_functionsCoordinator.getDONFeeJuels(_requestData),
      gasOverheadBeforeCallback: getCoordinatorConfig().gasOverheadBeforeCallback,
      gasOverheadAfterCallback: getCoordinatorConfig().gasOverheadAfterCallback
    });

    // topic0 (function signature, always checked), topic1 (true), topic2 (true), NOT topic3 (false), and data (true).
    vm.expectEmit(true, true, false, true);
    emit OracleRequest({
      requestId: expectedRequestId,
      requestingContract: address(s_functionsClient),
      requestInitiator: txOrigin,
      subscriptionId: s_subscriptionId,
      subscriptionOwner: OWNER_ADDRESS,
      data: _requestData,
      dataVersion: FunctionsRequest.REQUEST_DATA_VERSION,
      flags: bytes32(0),
      callbackGasLimit: _callbackGasLimit,
      commitment: expectedComittment
    });

    s_functionsCoordinator.startRequest(
      FunctionsResponse.RequestMeta({
        requestingContract: address(s_functionsClient),
        data: _requestData,
        subscriptionId: s_subscriptionId,
        dataVersion: FunctionsRequest.REQUEST_DATA_VERSION,
        flags: bytes32(0),
        callbackGasLimit: 5_500,
        adminFee: s_adminFee,
        initiatedRequests: 0,
        completedRequests: 0,
        availableBalance: s_subscriptionInitialFunding,
        subscriptionOwner: OWNER_ADDRESS
      })
    );
  }
}

/// @notice #_beforeSetConfig
contract FunctionsCoordinator__BeforeSetConfig {
  // TODO: make contract internal function helper
}

/// @notice #_getTransmitters
contract FunctionsCoordinator__GetTransmitters {
  // TODO: make contract internal function helper
}

/// @notice #_report
contract FunctionsCoordinator__Report {
  // TODO: make contract internal function helper
}

/// @notice #_onlyOwner
contract FunctionsCoordinator__OnlyOwner {
  // TODO: make contract internal function helper
}
