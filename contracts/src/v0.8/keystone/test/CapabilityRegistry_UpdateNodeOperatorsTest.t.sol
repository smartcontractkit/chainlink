// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_UpdateNodeOperatorTest is BaseTest {
  event NodeOperatorUpdated(uint32 indexed nodeOperatorId, address indexed admin, string name);

  uint32 private constant TEST_NODE_OPERATOR_ID = 1;
  address private constant NEW_NODE_OPERATOR_ADMIN = address(3);
  string private constant NEW_NODE_OPERATOR_NAME = "new-node-operator";

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_RevertWhen_CalledByNonAdminAndNonOwner() public {
    changePrank(STRANGER);

    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NEW_NODE_OPERATOR_ADMIN, name: NEW_NODE_OPERATOR_NAME});

    uint32[] memory nodeOperatorIds = new uint32[](1);
    nodeOperatorIds[0] = TEST_NODE_OPERATOR_ID;

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.AccessForbidden.selector, STRANGER));
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);
  }

  function test_RevertWhen_NodeOperatorAdminIsZeroAddress() public {
    changePrank(ADMIN);
    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: address(0), name: NEW_NODE_OPERATOR_NAME});

    uint32[] memory nodeOperatorIds = new uint32[](1);
    nodeOperatorIds[0] = TEST_NODE_OPERATOR_ID;

    vm.expectRevert(CapabilityRegistry.InvalidNodeOperatorAdmin.selector);
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);
  }

  function test_RevertWhen_NodeOperatorIdAndParamLengthsMismatch() public {
    changePrank(ADMIN);
    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NEW_NODE_OPERATOR_ADMIN, name: NEW_NODE_OPERATOR_NAME});

    uint32 invalidNodeOperatorId = 10000;
    uint32[] memory nodeOperatorIds = new uint32[](2);
    nodeOperatorIds[0] = invalidNodeOperatorId;
    vm.expectRevert(
      abi.encodeWithSelector(CapabilityRegistry.LengthMismatch.selector, nodeOperatorIds.length, nodeOperators.length)
    );
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);
  }

  function test_RevertWhen_NodeOperatorDoesNotExist() public {
    changePrank(ADMIN);
    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NEW_NODE_OPERATOR_ADMIN, name: NEW_NODE_OPERATOR_NAME});

    uint32 invalidNodeOperatorId = 10000;
    uint32[] memory nodeOperatorIds = new uint32[](1);
    nodeOperatorIds[0] = invalidNodeOperatorId;
    vm.expectRevert(
      abi.encodeWithSelector(CapabilityRegistry.NodeOperatorDoesNotExist.selector, invalidNodeOperatorId)
    );
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);
  }

  function test_UpdatesNodeOperator() public {
    changePrank(ADMIN);

    CapabilityRegistry.NodeOperator[] memory nodeOperators = new CapabilityRegistry.NodeOperator[](1);
    nodeOperators[0] = CapabilityRegistry.NodeOperator({admin: NEW_NODE_OPERATOR_ADMIN, name: NEW_NODE_OPERATOR_NAME});

    uint32[] memory nodeOperatorIds = new uint32[](1);
    nodeOperatorIds[0] = TEST_NODE_OPERATOR_ID;

    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit NodeOperatorUpdated(TEST_NODE_OPERATOR_ID, NEW_NODE_OPERATOR_ADMIN, NEW_NODE_OPERATOR_NAME);
    s_capabilityRegistry.updateNodeOperators(nodeOperatorIds, nodeOperators);

    CapabilityRegistry.NodeOperator memory nodeOperator = s_capabilityRegistry.getNodeOperator(TEST_NODE_OPERATOR_ID);
    assertEq(nodeOperator.admin, NEW_NODE_OPERATOR_ADMIN);
    assertEq(nodeOperator.name, NEW_NODE_OPERATOR_NAME);
  }
}
