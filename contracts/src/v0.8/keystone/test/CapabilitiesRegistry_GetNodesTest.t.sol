// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract CapabilitiesRegistry_GetNodesTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
    changePrank(ADMIN);

    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_CapabilitiesRegistry.addNodeOperators(_getNodeOperators());
    s_CapabilitiesRegistry.addCapabilities(capabilities);

    CapabilitiesRegistry.NodeParams[] memory nodes = new CapabilitiesRegistry.NodeParams[](2);
    bytes32[] memory hashedCapabilityIds = new bytes32[](2);
    hashedCapabilityIds[0] = s_basicHashedCapabilityId;
    hashedCapabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    nodes[1] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: hashedCapabilityIds
    });

    changePrank(NODE_OPERATOR_ONE_ADMIN);

    s_CapabilitiesRegistry.addNodes(nodes);
  }

  function test_CorrectlyFetchesNodes() public view {
    CapabilitiesRegistry.NodeInfo[] memory nodes = s_CapabilitiesRegistry.getNodes();
    assertEq(nodes.length, 2);

    assertEq(nodes[0].nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(nodes[0].signer, NODE_OPERATOR_ONE_SIGNER_ADDRESS);
    assertEq(nodes[0].p2pId, P2P_ID);
    assertEq(nodes[0].hashedCapabilityIds.length, 2);
    assertEq(nodes[0].hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(nodes[0].hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(nodes[0].configCount, 1);

    assertEq(nodes[1].nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(nodes[1].signer, NODE_OPERATOR_TWO_SIGNER_ADDRESS);
    assertEq(nodes[1].p2pId, P2P_ID_TWO);
    assertEq(nodes[1].hashedCapabilityIds.length, 2);
    assertEq(nodes[1].hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(nodes[1].hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(nodes[1].configCount, 1);
  }

  function test_DoesNotIncludeRemovedNodes() public {
    changePrank(ADMIN);
    bytes32[] memory nodesToRemove = new bytes32[](1);
    nodesToRemove[0] = P2P_ID_TWO;
    s_CapabilitiesRegistry.removeNodes(nodesToRemove);

    CapabilitiesRegistry.NodeInfo[] memory nodes = s_CapabilitiesRegistry.getNodes();
    assertEq(nodes.length, 1);

    assertEq(nodes[0].nodeOperatorId, TEST_NODE_OPERATOR_ONE_ID);
    assertEq(nodes[0].signer, NODE_OPERATOR_ONE_SIGNER_ADDRESS);
    assertEq(nodes[0].p2pId, P2P_ID);
    assertEq(nodes[0].hashedCapabilityIds.length, 2);
    assertEq(nodes[0].hashedCapabilityIds[0], s_basicHashedCapabilityId);
    assertEq(nodes[0].hashedCapabilityIds[1], s_capabilityWithConfigurationContractId);
    assertEq(nodes[0].configCount, 1);
  }
}
