// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_GetDONsTest is BaseTest {
  event ConfigSet(uint32 donId, uint32 configCount);

  CapabilityRegistry.CapabilityConfiguration[] private s_capabilityConfigs;

  function setUp() public override {
    BaseTest.setUp();

    CapabilityRegistry.Capability[] memory capabilities = new CapabilityRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapabilities(capabilities);

    CapabilityRegistry.NodeInfo[] memory nodes = new CapabilityRegistry.NodeInfo[](2);
    bytes32[] memory capabilityIds = new bytes32[](2);
    capabilityIds[0] = s_basicHashedCapabilityId;
    capabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: capabilityIds
    });

    bytes32[] memory nodeTwoCapabilityIds = new bytes32[](1);
    nodeTwoCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[1] = CapabilityRegistry.NodeInfo({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: nodeTwoCapabilityIds
    });

    changePrank(NODE_OPERATOR_ONE_ADMIN);
    s_capabilityRegistry.addNodes(nodes);

    s_capabilityConfigs.push(
      CapabilityRegistry.CapabilityConfiguration({
        capabilityId: s_basicHashedCapabilityId,
        config: BASIC_CAPABILITY_CONFIG
      })
    );

    bytes32[] memory nodeIds = new bytes32[](2);
    nodeIds[0] = P2P_ID;
    nodeIds[1] = P2P_ID_TWO;

    changePrank(ADMIN);
    s_capabilityRegistry.addDON(nodeIds, s_capabilityConfigs, true, true, 1);
    s_capabilityRegistry.addDON(nodeIds, s_capabilityConfigs, false, false, 1);
  }

  function test_CorrectlyFetchesDONs() public view {
    CapabilityRegistry.DONInfo[] memory dons = s_capabilityRegistry.getDONs();
    assertEq(dons.length, 2);
    assertEq(dons[0].id, DON_ID);
    assertEq(dons[0].configCount, 1);
    assertEq(dons[0].isPublic, true);
    assertEq(dons[0].acceptsWorkflows, true);
    assertEq(dons[0].f, 1);
    assertEq(dons[0].capabilityConfigurations.length, s_capabilityConfigs.length);
    assertEq(dons[0].capabilityConfigurations[0].capabilityId, s_basicHashedCapabilityId);

    assertEq(dons[1].id, DON_ID_TWO);
    assertEq(dons[1].configCount, 1);
    assertEq(dons[1].isPublic, false);
    assertEq(dons[1].capabilityConfigurations.length, s_capabilityConfigs.length);
    assertEq(dons[1].capabilityConfigurations[0].capabilityId, s_basicHashedCapabilityId);
  }

  function test_DoesNotIncludeRemovedDONs() public {
    uint32[] memory removedDONIDs = new uint32[](1);
    removedDONIDs[0] = DON_ID;
    s_capabilityRegistry.removeDONs(removedDONIDs);

    CapabilityRegistry.DONInfo[] memory dons = s_capabilityRegistry.getDONs();
    assertEq(dons.length, 1);
    assertEq(dons[0].id, DON_ID_TWO);
    assertEq(dons[0].configCount, 1);
    assertEq(dons[0].isPublic, false);
    assertEq(dons[0].acceptsWorkflows, false);
    assertEq(dons[0].f, 1);
    assertEq(dons[0].capabilityConfigurations.length, s_capabilityConfigs.length);
    assertEq(dons[0].capabilityConfigurations[0].capabilityId, s_basicHashedCapabilityId);
  }
}
