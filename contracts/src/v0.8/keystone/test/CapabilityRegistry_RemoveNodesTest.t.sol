// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_RemoveNodesTest is BaseTest {
  event NodeRemoved(bytes32 p2pId);

  uint32 private constant DON_ID = 1;
  uint32 private constant TEST_NODE_OPERATOR_ONE_ID = 1;
  uint32 private constant TEST_NODE_OPERATOR_TWO_ID = 2;
  bytes32 private constant INVALID_P2P_ID = bytes32("fake-p2p");
  bytes private constant BASIC_CAPABILITY_CONFIG = bytes("basic-capability-config");

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    changePrank(NODE_OPERATOR_ONE_ADMIN);

    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdminAndNonOwner() public {
    changePrank(STRANGER);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_NodeDoesNotExist() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = INVALID_P2P_ID;

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, INVALID_P2P_ID));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_P2PIDEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = bytes32("");

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, bytes32("")));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_NodePartOfDON() public {
    changePrank(ADMIN);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodePartOfDON.selector, P2P_ID, DON_ID));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_CanRemoveWhenDONDeleted() public {
    changePrank(ADMIN);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    // Add DON
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);

    // Try remove nodes
    bytes32[] memory removedNodes = new bytes32[](1);
    removedNodes[0] = P2P_ID;
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodePartOfDON.selector, P2P_ID, DON_ID));
    s_capabilityRegistry.removeNodes(removedNodes);

    // Remove DON
    uint32[] memory donIds = new uint32[](1);
    donIds[0] = DON_ID;
    s_capabilityRegistry.removeDONs(donIds);

    // Remove node
    s_capabilityRegistry.removeNodes(removedNodes);
    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, 0);
    assertEq(node.p2pId, bytes32(""));
    assertEq(node.signer, bytes32(""));
    assertEq(node.hashedCapabilityIds.length, 0);
    assertEq(configCount, 0);
  }

  function test_CanRemoveWhenNodeNoLongerPartOfDON() public {
    changePrank(ADMIN);
    CapabilityRegistry.NodeInfo[] memory newNodes = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    newNodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_TWO_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    s_capabilityRegistry.addNodes(newNodes);

    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    // Add DON
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);

    // Try remove nodes
    bytes32[] memory removedNodes = new bytes32[](1);
    removedNodes[0] = P2P_ID_TWO;
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodePartOfDON.selector, P2P_ID_TWO, DON_ID));
    s_capabilityRegistry.removeNodes(removedNodes);

    // Update nodes in DON
    bytes32[] memory updatedNodes = new bytes32[](1);
    updatedNodes[0] = P2P_ID;
    s_capabilityRegistry.updateDON(DON_ID, updatedNodes, capabilityConfigs, true);

    // Remove node
    s_capabilityRegistry.removeNodes(removedNodes);
    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID_TWO);
    assertEq(node.nodeOperatorId, 0);
    assertEq(node.p2pId, bytes32(""));
    assertEq(node.signer, bytes32(""));
    assertEq(node.hashedCapabilityIds.length, 0);
    assertEq(configCount, 0);
  }

  function test_RemovesNode() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeRemoved(P2P_ID);
    s_capabilityRegistry.removeNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, 0);
    assertEq(node.p2pId, bytes32(""));
    assertEq(node.signer, bytes32(""));
    assertEq(node.hashedCapabilityIds.length, 0);
    assertEq(configCount, 0);
  }

  function test_CanAddNodeWithSameSignerAddressAfterRemoving() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    s_capabilityRegistry.removeNodes(nodes);

    CapabilityRegistry.NodeInfo[] memory NodeInfo = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    NodeInfo[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    s_capabilityRegistry.addNodes(NodeInfo);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.hashedCapabilityIds.length, 2);
    assertEq(node.hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(node.hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(configCount, 1);
  }

  function test_OwnerCanRemoveNodes() public {
    changePrank(ADMIN);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeRemoved(P2P_ID);
    s_capabilityRegistry.removeNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, 0);
    assertEq(node.p2pId, bytes32(""));
    assertEq(node.signer, bytes32(""));
    assertEq(node.hashedCapabilityIds.length, 0);
    assertEq(configCount, 0);
  }
}
