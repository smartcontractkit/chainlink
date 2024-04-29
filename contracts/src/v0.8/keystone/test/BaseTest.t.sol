// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {Test} from "forge-std/Test.sol";
import {Constants} from "./Constants.t.sol";
import {CapabilityConfigurationContract} from "./mocks/CapabilityConfigurationContract.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract BaseTest is Test, Constants {
  CapabilityRegistry internal s_capabilityRegistry;
  CapabilityConfigurationContract internal s_capabilityConfigurationContract;

  function setUp() public virtual {
    vm.startPrank(ADMIN);
    s_capabilityRegistry = new CapabilityRegistry();
    s_capabilityConfigurationContract = new CapabilityConfigurationContract();
  }

  function _getNodeOperators() internal view returns (CapabilityRegistry.NodeOperator[] memory) {
    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](2);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NODE_OPERATOR_ONE_ADMIN, name: NODE_OPERATOR_ONE_NAME});
    nodeOperators[1] = CapabilityRegistry.NodeOperator({admin: NODE_OPERATOR_TWO_ADMIN, name: NODE_OPERATOR_TWO_NAME});
    return nodeOperators;
  }
}
