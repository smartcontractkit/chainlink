// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_UpdateNodesTest is BaseTest {
  event NodeUpdated(bytes32 p2pId, uint32 indexed nodeOperatorId, bytes32 signer);

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapabilities(capabilities);

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

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_TWO_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    changePrank(NODE_OPERATOR_TWO_ADMIN);
    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdminAndNonOwner() public {
    changePrank(STRANGER);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.AccessForbidden.selector, STRANGER));
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_NodeDoesNotExist() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: INVALID_P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodeDoesNotExist.selector, INVALID_P2P_ID));
    s_capabilityRegistry.updateNodes(nodes);
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

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.NodeDoesNotExist.selector, bytes32("")));
    s_capabilityRegistry.updateNodes(nodes);
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
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_NodeSignerAlreadyAssignedToAnotherNode() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(CapabilityRegistry.InvalidNodeSigner.selector);
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_UpdatingNodeWithoutCapabilities() public {
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
    s_capabilityRegistry.updateNodes(nodes);
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
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_CanUpdateParamsIfNodeSignerAddressNoLongerUsed() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    // Set node one's signer to another address
    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: bytes32(abi.encodePacked(address(6666))),
      hashedCapabilityIds: hashedCapabilityIds
    });

    s_capabilityRegistry.updateNodes(nodes);

    // Set node two's signer to node one's signer
    changePrank(NODE_OPERATOR_TWO_ADMIN);
    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_TWO_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });
    s_capabilityRegistry.updateNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, ) = s_capabilityRegistry.getNode(P2P_ID_TWO);
    assertEq(node.signer, NODE_OPERATOR_ONE_SIGNER_ADDRESS);
  }

  function test_UpdatesNodeInfo() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NEW_NODE_SIGNER,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeUpdated(P2P_ID, TEST_NODE_OPERATOR_ONE_ID, NEW_NODE_SIGNER);
    s_capabilityRegistry.updateNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.signer, NEW_NODE_SIGNER);
    assertEq(node.hashedCapabilityIds.length, 1);
    assertEq(node.hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(configCount, 2);
  }

  function test_OwnerCanUpdateNodes() public {
    changePrank(ADMIN);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NEW_NODE_SIGNER,
      hashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeUpdated(P2P_ID, TEST_NODE_OPERATOR_ONE_ID, NEW_NODE_SIGNER);
    s_capabilityRegistry.updateNodes(nodes);

    (CapabilityRegistry.NodeInfo memory node, uint32 configCount) = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.signer, NEW_NODE_SIGNER);
    assertEq(node.hashedCapabilityIds.length, 1);
    assertEq(node.hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(configCount, 2);
  }
}
