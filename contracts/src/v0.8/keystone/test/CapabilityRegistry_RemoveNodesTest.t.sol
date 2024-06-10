// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_RemoveNodesTest is BaseTest {
  event NodeRemoved(bytes32 p2pId);

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapabilities(capabilities);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](3);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    nodes[1] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_TWO_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    nodes[2] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_THREE_ID,
      p2pId: P2P_ID_THREE,
      signer: NODE_OPERATOR_THREE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    changePrank(ADMIN);

    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdminAndNonOwner() public {
    changePrank(STRANGER);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.AccessForbidden.selector, STRANGER));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_NodeDoesNotExist() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = INVALID_P2P_ID;

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodeDoesNotExist.selector, INVALID_P2P_ID));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_P2PIDEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = bytes32("");

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodeDoesNotExist.selector, bytes32("")));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_RevertWhen_NodePartOfDON() public {
    changePrank(ADMIN);
    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true, true, F_VALUE);

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodePartOfDON.selector, P2P_ID));
    s_capabilityRegistry.removeNodes(nodes);
  }

  function test_CanRemoveWhenDONDeleted() public {
    changePrank(ADMIN);

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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true, true, F_VALUE);

    // Try remove nodes
    bytes32[] memory removedNodes = new bytes32[](1);
    removedNodes[0] = P2P_ID;
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodePartOfDON.selector, P2P_ID));
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

    bytes32[] memory nodes = new bytes32[](3);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;
    nodes[2] = P2P_ID_THREE;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    // Add DON
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true, true, F_VALUE);

    // Try remove nodes
    bytes32[] memory removedNodes = new bytes32[](1);
    removedNodes[0] = P2P_ID_TWO;
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodePartOfDON.selector, P2P_ID_TWO));
    s_capabilityRegistry.removeNodes(removedNodes);

    // Update nodes in DON
    bytes32[] memory updatedNodes = new bytes32[](2);
    updatedNodes[0] = P2P_ID;
    updatedNodes[1] = P2P_ID_THREE;
    s_capabilityRegistry.updateDON(DON_ID, updatedNodes, capabilityConfigs, true, true, F_VALUE);

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
