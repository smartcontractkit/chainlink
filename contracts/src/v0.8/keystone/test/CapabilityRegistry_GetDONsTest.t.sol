// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_GetDONsTest is BaseTest {
  event ConfigSet(uint32 donId, uint32 configCount);

  uint32 private constant DON_ID_ONE = 1;
  uint32 private constant DON_ID_TWO = 2;
  uint32 private constant TEST_NODE_OPERATOR_ONE_ID = 1;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 2;
  bytes32 private constant INVALID_P2P_ID = bytes32("fake-p2p");
  bytes private constant CONFIG = bytes("onchain-config");
  CapabilityRegistry.CapabilityConfiguration[] private s_capabilityConfigs;

  function setUp() public override {
    BaseTest.setUp();

    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);

    CapabilityRegistry.NodeParams[] memory nodes = new CapabilityRegistry.NodeParams[](2);
    bytes32[] memory capabilityIds = new bytes32[](2);
    capabilityIds[0] = s_basicHashedCapabilityId;
    capabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilityRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: capabilityIds
    });

    bytes32[] memory nodeTwoCapabilityIds = new bytes32[](1);
    nodeTwoCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[1] = CapabilityRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: nodeTwoCapabilityIds
    });

    changePrank(NODE_OPERATOR_ONE_ADMIN);
    s_capabilityRegistry.addNodes(nodes);

    s_capabilityConfigs.push(
      CapabilityRegistry.CapabilityConfiguration({capabilityId: s_basicHashedCapabilityId, config: CONFIG})
    );

    bytes32[] memory nodeIds = new bytes32[](2);
    nodeIds[0] = P2P_ID;
    nodeIds[1] = P2P_ID_TWO;

    changePrank(ADMIN);
    s_capabilityRegistry.addDON(nodeIds, s_capabilityConfigs, true);
    s_capabilityRegistry.addDON(nodeIds, s_capabilityConfigs, false);
  }

  function test_CorrectlyFetchesDONs() public view {
    CapabilityRegistry.DONParams[] memory dons = s_capabilityRegistry.getDONs();
    assertEq(dons.length, 2);
    assertEq(dons[0].id, DON_ID_ONE);
    assertEq(dons[0].configCount, 1);
    assertEq(dons[0].isPublic, true);
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
    removedDONIDs[0] = DON_ID_ONE;
    s_capabilityRegistry.removeDONs(removedDONIDs);

    CapabilityRegistry.DONParams[] memory dons = s_capabilityRegistry.getDONs();
    assertEq(dons.length, 1);
    assertEq(dons[0].id, DON_ID_TWO);
    assertEq(dons[0].configCount, 1);
    assertEq(dons[0].isPublic, false);
    assertEq(dons[0].capabilityConfigurations.length, s_capabilityConfigs.length);
    assertEq(dons[0].capabilityConfigurations[0].capabilityId, s_basicHashedCapabilityId);
  }
}
