// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_UpdateNodeOperatorTest is BaseTest {
  event NodeOperatorUpdated(uint256 nodeOperatorId, address indexed admin, string name);

  uint256 private constant TEST_NODE_OPERATOR_ID = 0;
  address private constant NEW_NODE_OPERATOR_ADMIN = address(3);
  string private constant NEW_NODE_OPERATOR_NAME = "new-node-operator";

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_RevertWhen_CalledByNonAdminAndNonOwner() public {
    changePrank(STRANGER);
    vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);

    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NEW_NODE_OPERATOR_ADMIN, name: NEW_NODE_OPERATOR_NAME});

    uint256[] memory nodeOperatorIds = new uint256[](1);
    nodeOperatorIds[0] = TEST_NODE_OPERATOR_ID;
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);
  }

  function test_RevertWhen_NodeOperatorAdminIsZeroAddress() public {
    changePrank(ADMIN);
    vm.expectRevert(CapabilityRegistry.InvalidNodeOperatorAdmin.selector);
    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: address(0), name: NEW_NODE_OPERATOR_NAME});

    uint256[] memory nodeOperatorIds = new uint256[](1);
    nodeOperatorIds[0] = TEST_NODE_OPERATOR_ID;
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);
  }

  function test_UpdatesNodeOperator() public {
    changePrank(ADMIN);

    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NEW_NODE_OPERATOR_ADMIN, name: NEW_NODE_OPERATOR_NAME});

    uint256[] memory nodeOperatorIds = new uint256[](1);
    nodeOperatorIds[0] = TEST_NODE_OPERATOR_ID;

    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit NodeOperatorUpdated(TEST_NODE_OPERATOR_ID, NEW_NODE_OPERATOR_ADMIN, NEW_NODE_OPERATOR_NAME);
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);

    CapabilityRegistry.NodeOperator memory nodeOperator = s_capabilityRegistry.getNodeOperator(0);
    assertEq(nodeOperator.admin, NEW_NODE_OPERATOR_ADMIN);
    assertEq(nodeOperator.name, NEW_NODE_OPERATOR_NAME);
  }
}
