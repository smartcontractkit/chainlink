// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddNodesTest is BaseTest {
  event NodeAdded(bytes32 p2pId, uint32 indexed nodeOperatorId, bytes32 signer);

  function setUp() public override {
    BaseTest.setUp();
    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapabilities(capabilities);
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdminAndNonOwner() public {
    changePrank(STRANGER);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.AccessForbidden.selector, STRANGER));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_AddingNodeWithInvalidNodeOperator() public {
    changePrank(ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    uint32 invalidNodeOperatorId = 10000;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: invalidNodeOperatorId, // Invalid NOP
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(
      abi.encodeWithSelector(CapabilityRegistry.NodeOperatorDoesNotExist.selector, invalidNodeOperatorId)
    );
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_SignerAddressEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: bytes32(""),
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeSigner.selector));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_SignerAddressNotUnique() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    s_capabilityRegistry.addNodes(nodes);

    changePrank(NODE_OPERATOR_TWO_ADMIN);

    // Try adding another node with the same signer address
    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_TWO_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeSigner.selector));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_AddingDuplicateP2PId() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    s_capabilityRegistry.addNodes(nodes);

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodeAlreadyExists.selector, P2P_ID));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_P2PIDEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: bytes32(""),
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, bytes32("")));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_AddingNodeWithoutCapabilities() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](0);

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeCapabilities.selector, hashedCapabilityIds));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_AddingNodeWithInvalidCapability() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_nonExistentHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeCapabilities.selector, hashedCapabilityIds));
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_AddsNodeInfo() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

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

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeAdded(P2P_ID, TEST_NODE_OPERATOR_ONE_ID, NODE_OPERATOR_ONE_SIGNER_ADDRESS);
    s_capabilityRegistry.addNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.hashedCapabilityIds.length, 2);
    assertEq(node.hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(node.hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(configCount, 1);
  }

  function test_OwnerCanAddNodes() public {
    changePrank(ADMIN);

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

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeAdded(P2P_ID, TEST_NODE_OPERATOR_ONE_ID, NODE_OPERATOR_ONE_SIGNER_ADDRESS);
    s_capabilityRegistry.addNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.hashedCapabilityIds.length, 2);
    assertEq(node.hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(node.hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(configCount, 1);
  }
}
