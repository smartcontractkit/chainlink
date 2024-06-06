// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_GetNodesTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);

    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapabilities(capabilities);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](2);
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
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    changePrank(NODE_OPERATOR_ONE_ADMIN);

    s_capabilityRegistry.addNodes(nodes);
  }

  function test_CorrectlyFetchesNodes() public view {
    (CapabilityRegistry.NodeInfo[] memory nodes, uint32[] memory configCounts) = s_capabilityRegistry.getNodes();
    assertEq(nodes.length, 2);

    assertEq(nodes[0].nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(nodes[0].signer, NODE_OPERATOR_ONE_SIGNER_ADDRESS);
    assertEq(nodes[0].p2pId, P2P_ID);
    assertEq(nodes[0].hashedCapabilityIds.length, 2);
    assertEq(nodes[0].hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(nodes[0].hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(configCounts[0], 1);

    assertEq(nodes[1].nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(nodes[1].signer, NODE_OPERATOR_TWO_SIGNER_ADDRESS);
    assertEq(nodes[1].p2pId, P2P_ID_TWO);
    assertEq(nodes[1].hashedCapabilityIds.length, 2);
    assertEq(nodes[1].hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(nodes[1].hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(configCounts[1], 1);
  }

  function test_DoesNotIncludeRemovedNodes() public {
    changePrank(ADMIN);
    bytes32[] memory nodesToRemove = new bytes32[](1);
    nodesToRemove[0] = P2P_ID_TWO;
    s_capabilityRegistry.removeNodes(nodesToRemove);

    (CapabilityRegistry.NodeInfo[] memory nodes, uint32[] memory configCounts) = s_capabilityRegistry.getNodes();
    assertEq(nodes.length, 1);

    assertEq(nodes[0].nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(nodes[0].signer, NODE_OPERATOR_ONE_SIGNER_ADDRESS);
    assertEq(nodes[0].p2pId, P2P_ID);
    assertEq(nodes[0].hashedCapabilityIds.length, 2);
    assertEq(nodes[0].hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(nodes[0].hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(configCounts[0], 1);
  }
}
