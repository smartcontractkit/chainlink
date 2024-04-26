// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddNodeOperatorsTest is BaseTest {
  event NodeOperatorAdded(uint256 nodeOperatorId, address indexed admin, string name);

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_RevertWhen_NodeOperatorAdminAddressZero() public {
    changePrank(ADMIN);
    CapabilityRegistry.NodeOperator[] memory nodeOperators = _getNodeOperators();
    nodeOperators[0].admin = address(0);
    vm.expectRevert(CapabilityRegistry.InvalidNodeOperatorAdmin.selector);
    s_capabilityRegistry.addNodeOperators(nodeOperators);
  }

  function test_AddNodeOperators() public {
    changePrank(ADMIN);

    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit NodeOperatorAdded(0, NODE_OPERATOR_ONE_ADMIN, NODE_OPERATOR_ONE_NAME);
    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit NodeOperatorAdded(1, NODE_OPERATOR_TWO_ADMIN, NODE_OPERATOR_TWO_NAME);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());

    CapabilityRegistry.NodeOperator memory nodeOperatorOne = s_capabilityRegistry.getNodeOperator(0);
    assertEq(nodeOperatorOne.admin, NODE_OPERATOR_ONE_ADMIN);
    assertEq(nodeOperatorOne.name, NODE_OPERATOR_ONE_NAME);

    CapabilityRegistry.NodeOperator memory nodeOperatorTwo = s_capabilityRegistry.getNodeOperator(1);
    assertEq(nodeOperatorTwo.admin, NODE_OPERATOR_TWO_ADMIN);
    assertEq(nodeOperatorTwo.name, NODE_OPERATOR_TWO_NAME);
  }
}
