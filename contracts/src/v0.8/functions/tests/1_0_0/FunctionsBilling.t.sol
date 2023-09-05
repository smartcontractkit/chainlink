// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsBilling} from "../../dev/1_0_0/FunctionsBilling.sol";

import {FunctionsRouterSetup} from "./Setup.t.sol";

/// @notice #constructor
contract FunctionsBilling_Constructor {

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
  // TODO: make contract internal function helper
}

/// @notice #estimateCost
contract FunctionsBilling_EstimateCost {

}

/// @notice #_calculateCostEstimate
contract FunctionsBilling__CalculateCostEstimate {
  // TODO: make contract internal function helper
}

/// @notice #_startBilling
contract FunctionsBilling__StartBilling {
  // TODO: make contract internal function helper
}

/// @notice #_computeRequestId
contract FunctionsBilling__ComputeRequestId {
  // TODO: make contract internal function helper
}

/// @notice #_fulfillAndBill
contract FunctionsBilling__FulfillAndBill {
  // TODO: make contract internal function helper
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
  // TODO: make contract internal function helper
}

/// @notice #_disperseFeePool
contract FunctionsBilling__DisperseFeePool {
  // TODO: make contract internal function helper
}
