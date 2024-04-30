// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddNodesTest is BaseTest {
  event NodeAdded(bytes p2pId, uint256 nodeOperatorId);

  uint256 private constant TEST_NODE_OPERATOR_ONE_ID = 0;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 1;

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdmin() public {
    changePrank(STRANGER);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    string[] memory capabilityIds = new string[](1);
    capabilityIds[0] = "ccip-exec-0.0.1";

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      supportedCapabilityIds: capabilityIds
    });

    vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_P2PIDEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    string[] memory capabilityIds = new string[](1);
    capabilityIds[0] = "ccip-exec-0.0.1";

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: bytes(""),
      supportedCapabilityIds: capabilityIds
    });

    vm.expectRevert(CapabilityRegistry.InvalidNodeP2PId.selector);
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_AddsNode() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);
    string[] memory capabilityIds = new string[](1);
    capabilityIds[0] = "ccip-exec-0.0.1";

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      supportedCapabilityIds: capabilityIds
    });

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeAdded(P2P_ID, TEST_NODE_OPERATOR_ONE_ID);
    s_capabilityRegistry.addNodes(nodes);

    CapabilityRegistry.Node memory node = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.supportedCapabilityIds.length, 1);
    assertEq(node.supportedCapabilityIds[0], "ccip-exec-0.0.1");
  }
}
