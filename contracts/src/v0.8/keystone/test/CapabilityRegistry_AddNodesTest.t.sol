// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddNodesTest is BaseTest {
  event NodeAdded(bytes32 p2pId, uint256 nodeOperatorId);

  uint256 private constant TEST_NODE_OPERATOR_ONE_ID = 0;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 1;

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdmin() public {
    changePrank(STRANGER);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory capabilityIds = new bytes32[](1);
    capabilityIds[0] = s_basicCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedCapabilityIds: capabilityIds
    });

    vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_AddingDuplicateP2PId() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory capabilityIds = new bytes32[](1);
    capabilityIds[0] = s_basicCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedCapabilityIds: capabilityIds
    });
    s_capabilityRegistry.addNodes(nodes);

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, P2P_ID));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_P2PIDEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory capabilityIds = new bytes32[](1);
    capabilityIds[0] = s_basicCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: bytes32(""),
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedCapabilityIds: capabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, bytes32("")));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_AddingNodeWithoutCapabilities() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory capabilityIds = new bytes32[](0);

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedCapabilityIds: capabilityIds
    });
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeCapabilities.selector, capabilityIds));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_AddingNodeWithInvalidCapability() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory capabilityIds = new bytes32[](1);
    capabilityIds[0] = s_nonExistentCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedCapabilityIds: capabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeCapabilities.selector, capabilityIds));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_AddsNode() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);
    bytes32[] memory capabilityIds = new bytes32[](2);
    capabilityIds[0] = s_basicCapabilityId;
    capabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedCapabilityIds: capabilityIds
    });

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeAdded(P2P_ID, TEST_NODE_OPERATOR_ONE_ID);
    s_capabilityRegistry.addNodes(nodes);

    CapabilityRegistry.Node memory node = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.supportedCapabilityIds.length, 2);
    assertEq(node.supportedCapabilityIds[0], s_basicCapabilityId);
    assertEq(node.supportedCapabilityIds[1], s_capabilityWithConfigurationContractId);
  }
}
