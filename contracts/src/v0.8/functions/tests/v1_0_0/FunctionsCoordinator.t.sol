// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {FunctionsCoordinator} from "../../dev/v1_0_0/FunctionsCoordinator.sol";
import {FunctionsBilling} from "../../dev/v1_0_0/FunctionsBilling.sol";
import {FunctionsRequest} from "../../dev/v1_0_0/libraries/FunctionsRequest.sol";
import {FunctionsRouter} from "../../dev/v1_0_0/FunctionsRouter.sol";

import {FunctionsRouterSetup} from "./Setup.t.sol";

/// @notice #constructor
contract FunctionsCoordinator_Constructor is FunctionsRouterSetup {
  function test_Constructor_Success() public {
    assertEq(s_functionsCoordinator.typeAndVersion(), "Functions Coordinator v1.0.0");
    assertEq(s_functionsCoordinator.owner(), OWNER_ADDRESS);
  }
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
  // TODO: make contract internal function helper
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
  // TODO: make contract internal function helper
}

/// @notice #_getTransmitters
contract FunctionsCoordinator__GetTransmitters {
  // TODO: make contract internal function helper
}

/// @notice #_report
contract FunctionsCoordinator__Report {

}

/// @notice #_onlyOwner
contract FunctionsCoordinator__OnlyOwner {

}
