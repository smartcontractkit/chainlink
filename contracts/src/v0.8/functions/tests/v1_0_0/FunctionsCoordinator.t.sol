// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../dev/v1_0_0/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/v1_0_0/FunctionsBilling.sol";
import {FunctionsRequest} from "../../dev/v1_0_0/libraries/FunctionsRequest.sol";

import {FunctionsSubscriptionSetup, FunctionsMultipleFulfillmentsSetup} from "./Setup.t.sol";

// ================================================================
// |                       Functions Coordinator                  |
// ================================================================

/// @notice #constructor
contract FunctionsCoordinator_Constructor {

}

/// @notice #getThresholdPublicKey
contract FunctionsCoordinator_GetThresholdPublicKey {

}

/// @notice #setThresholdPublicKey
contract FunctionsCoordinator_SetThresholdPublicKey {

}

/// @notice #getDONPublicKey
contract FunctionsCoordinator_GetDONPublicKey {

}

/// @notice #setDONPublicKey
contract FunctionsCoordinator__SetDONPublicKey {

}

/// @notice #_isTransmitter
contract FunctionsCoordinator_IsTransmitter {

}

/// @notice #setNodePublicKey
contract FunctionsCoordinator_SetNodePublicKey {

}

/// @notice #deleteNodePublicKey
contract FunctionsCoordinator_DeleteNodePublicKey {

}

/// @notice #getAllNodePublicKeys
contract FunctionsCoordinator_GetAllNodePublicKeys {

}

/// @notice #startRequest
contract FunctionsCoordinator_StartRequest {

}

/// @notice #_beforeSetConfig
contract FunctionsCoordinator__BeforeSetConfig {

}

/// @notice #_getTransmitters
contract FunctionsCoordinator__GetTransmitters {

}

/// @notice #_report
contract FunctionsCoordinator__Report {

}

/// @notice #_onlyOwner
contract FunctionsCoordinator__OnlyOwner {

}

// ================================================================
// |                        Functions Billing                     |
// ================================================================

/// @notice #constructor
contract FunctionsBilling_Constructor {

}

/// @notice #getConfig
contract FunctionsBilling_GetConfig {

}

/// @notice #updateConfig
contract FunctionsBilling_UpdateConfig {

}

/// @notice #getDONFee
contract FunctionsBilling_GetDONFee {

}

/// @notice #getAdminFee
contract FunctionsBilling_GetAdminFee {

}

/// @notice #getWeiPerUnitLink
contract FunctionsBilling_GetWeiPerUnitLink {

}

/// @notice #_getJuelsPerGas
contract FunctionsBilling__GetJuelsPerGas {

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
    FunctionsRequest.initializeRequest(
      request,
      FunctionsRequest.Location.Inline,
      FunctionsRequest.CodeLanguage.JavaScript,
      sourceCode
    );
    bytes memory requestData = FunctionsRequest.encodeCBOR(request);

    uint32 callbackGasLimit = 5_500;
    uint256 gasPriceWei = REASONABLE_GAS_PRICE_CEILING + 1;

    vm.expectRevert(FunctionsBilling.InvalidCalldata.selector);

    s_functionsCoordinator.estimateCost(s_subscriptionId, requestData, callbackGasLimit, gasPriceWei);
  }

  function test_EstimateCost_Success() public {
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

    uint32 callbackGasLimit = 5_500;
    uint256 gasPriceWei = 1;

    uint96 costEstimate = s_functionsCoordinator.estimateCost(
      s_subscriptionId,
      requestData,
      callbackGasLimit,
      gasPriceWei
    );
    uint96 expectedCostEstimate = 10873200;
    assertEq(costEstimate, expectedCostEstimate);
  }
}

/// @notice #_calculateCostEstimate
contract FunctionsBilling__CalculateCostEstimate {

}

/// @notice #_startBilling
contract FunctionsBilling__StartBilling {

}

/// @notice #_computeRequestId
contract FunctionsBilling__ComputeRequestId {

}

/// @notice #_fulfillAndBill
contract FunctionsBilling__FulfillAndBill {

}

/// @notice #deleteCommitment
contract FunctionsBilling_DeleteCommitment {

}

/// @notice #oracleWithdraw
contract FunctionsBilling_OracleWithdraw {

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

/// @notice #_getTransmitters
contract FunctionsBilling__GetTransmitters {

}

/// @notice #_disperseFeePool
contract FunctionsBilling__DisperseFeePool {

}

// ================================================================
// |                           OCR2Base                           |
// ================================================================

/// @notice #constructor
contract OCR2Base_Constructor {

}

/// @notice #checkConfigValid
contract OCR2Base_CheckConfigValid {

}

/// @notice #latestConfigDigestAndEpoch
contract OCR2Base_LatestConfigDigestAndEpoch {

}

/// @notice #setConfig
contract OCR2Base_SetConfig {

}

/// @notice #configDigestFromConfigData
contract OCR2Base_ConfigDigestFromConfigData {

}

/// @notice #latestConfigDetails
contract OCR2Base_LatestConfigDetails {

}

/// @notice #transmitters
contract OCR2Base_Transmitters {

}

/// @notice #_report
contract OCR2Base__Report {

}

/// @notice #requireExpectedMsgDataLength
contract OCR2Base_RequireExpectedMsgDataLength {

}

/// @notice #transmit
contract OCR2Base_Transmit {

}
