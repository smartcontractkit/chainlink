// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract CapabilitiesRegistry_AddNodeOperatorsTest is BaseTest {
  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_CapabilitiesRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_RevertWhen_NodeOperatorAdminAddressZero() public {
    changePrank(ADMIN);
    CapabilitiesRegistry.NodeOperator[] memory nodeOperators = _getNodeOperators();
    nodeOperators[0].admin = address(0);
    vm.expectRevert(CapabilitiesRegistry.InvalidNodeOperatorAdmin.selector);
    s_CapabilitiesRegistry.addNodeOperators(nodeOperators);
  }

  function test_AddNodeOperators() public {
    changePrank(ADMIN);

    vm.expectEmit(true, true, true, true, address(s_CapabilitiesRegistry));
    emit CapabilitiesRegistry.NodeOperatorAdded(
      TEST_NODE_OPERATOR_ONE_ID,
      NODE_OPERATOR_ONE_ADMIN,
      NODE_OPERATOR_ONE_NAME
    );
    vm.expectEmit(true, true, true, true, address(s_CapabilitiesRegistry));
    emit CapabilitiesRegistry.NodeOperatorAdded(
      TEST_NODE_OPERATOR_TWO_ID,
      NODE_OPERATOR_TWO_ADMIN,
      NODE_OPERATOR_TWO_NAME
    );
    s_CapabilitiesRegistry.addNodeOperators(_getNodeOperators());

    CapabilitiesRegistry.NodeOperator memory nodeOperatorOne = s_CapabilitiesRegistry.getNodeOperator(
      TEST_NODE_OPERATOR_ONE_ID
    );
    assertEq(nodeOperatorOne.admin, NODE_OPERATOR_ONE_ADMIN);
    assertEq(nodeOperatorOne.name, NODE_OPERATOR_ONE_NAME);

    CapabilitiesRegistry.NodeOperator memory nodeOperatorTwo = s_CapabilitiesRegistry.getNodeOperator(
      TEST_NODE_OPERATOR_TWO_ID
    );
    assertEq(nodeOperatorTwo.admin, NODE_OPERATOR_TWO_ADMIN);
    assertEq(nodeOperatorTwo.name, NODE_OPERATOR_TWO_NAME);
  }
}
