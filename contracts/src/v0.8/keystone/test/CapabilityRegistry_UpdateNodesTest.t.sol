// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_UpdateNodesTest is BaseTest {
  event NodeUpdated(bytes32 p2pId, uint256 nodeOperatorId, address signer);

  uint256 private constant TEST_NODE_OPERATOR_ONE_ID = 0;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 1;
  bytes32 private constant INVALID_P2P_ID = bytes32("fake-p2p");

  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);
    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);

    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    changePrank(NODE_OPERATOR_ONE_ADMIN);

    s_capabilityRegistry.addNodes(nodes);
  }

  function test_RevertWhen_CalledByNonNodeOperatorAdminAndNonOwner() public {
    changePrank(STRANGER);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(CapabilityRegistry.AccessForbidden.selector);
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_NodeDoesNotExist() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: INVALID_P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, INVALID_P2P_ID));
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_P2PIDEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: bytes32(""),
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeP2PId.selector, bytes32("")));
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_SignerAddressEmpty() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: address(0),
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeSigner.selector));
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_UpdatingNodeWithoutCapabilities() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](0);

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeCapabilities.selector, hashedCapabilityIds));
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_RevertWhen_AddingNodeWithInvalidCapability() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);
    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);

    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_nonExistentHashedCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.InvalidNodeCapabilities.selector, hashedCapabilityIds));
    s_capabilityRegistry.updateNodes(nodes);
  }

  function test_UpdatesNode() public {
    changePrank(NODE_OPERATOR_ONE_ADMIN);

    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeUpdated(P2P_ID, TEST_NODE_OPERATOR_ONE_ID, NODE_OPERATOR_TWO_SIGNER_ADDRESS);
    s_capabilityRegistry.updateNodes(nodes);

    CapabilityRegistry.Node memory node = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.signer, NODE_OPERATOR_TWO_SIGNER_ADDRESS);
    assertEq(node.supportedHashedCapabilityIds.length, 1);
    assertEq(node.supportedHashedCapabilityIds[0], s_basicHashedCapabilityId);
  }

  function test_OwnerCanUpdateNodes() public {
    changePrank(ADMIN);

    CapabilityRegistry.Node[] memory nodes = new CapabilityRegistry.Node[](1);
    bytes32[] memory hashedCapabilityIds = new bytes32[](1);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[0] = CapabilityRegistry.Node({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      supportedHashedCapabilityIds: hashedCapabilityIds
    });

    vm.expectEmit(address(s_capabilityRegistry));
    emit NodeUpdated(P2P_ID, TEST_NODE_OPERATOR_ONE_ID, NODE_OPERATOR_TWO_SIGNER_ADDRESS);
    s_capabilityRegistry.updateNodes(nodes);

    CapabilityRegistry.Node memory node = s_capabilityRegistry.getNode(P2P_ID);
    assertEq(node.nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(node.p2pId, P2P_ID);
    assertEq(node.signer, NODE_OPERATOR_TWO_SIGNER_ADDRESS);
    assertEq(node.supportedHashedCapabilityIds.length, 1);
    assertEq(node.supportedHashedCapabilityIds[0], s_basicHashedCapabilityId);
  }
}
