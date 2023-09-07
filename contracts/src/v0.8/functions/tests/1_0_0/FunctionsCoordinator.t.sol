// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../dev/1_0_0/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/1_0_0/FunctionsBilling.sol";
import {FunctionsRequest} from "../../dev/1_0_0/libraries/FunctionsRequest.sol";

import {FunctionsSubscriptionSetup} from "./Setup.t.sol";

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

  function test_EstimateCost_Success() public view {
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

    s_functionsCoordinator.estimateCost(s_subscriptionId, requestData, callbackGasLimit, gasPriceWei);
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
contract FunctionsBilling_OracleWithdrawAll {

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
