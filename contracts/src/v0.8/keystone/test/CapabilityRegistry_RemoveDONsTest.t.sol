// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_RemoveDONsTest is BaseTest {
  event ConfigSet(uint32 donId, uint32 configCount);

  uint32 private constant DON_ID = 1;
  uint32 private constant TEST_NODE_OPERATOR_ONE_ID = 0;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 1;
  bytes32 private constant INVALID_P2P_ID = bytes32("fake-p2p");
  bytes private constant CONFIG = bytes("onchain-config");

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

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: CONFIG
    });

    bytes32[] memory nodeIds = new bytes32[](2);
    nodeIds[0] = P2P_ID;
    nodeIds[1] = P2P_ID_TWO;

    changePrank(ADMIN);
    s_capabilityRegistry.addDON(nodeIds, capabilityConfigs, true);
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    uint32[] memory donIDs = new uint32[](1);
    donIDs[0] = 1;
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_capabilityRegistry.removeDONs(donIDs);
  }

  function test_RevertWhen_DONDoesNotExist() public {
    uint32 invalidDONId = 10;
    uint32[] memory donIDs = new uint32[](1);
    donIDs[0] = invalidDONId;
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.DONDoesNotExist.selector, invalidDONId));
    s_capabilityRegistry.removeDONs(donIDs);
  }

  function test_RemovesDON() public {
    uint32[] memory donIDs = new uint32[](1);
    donIDs[0] = DON_ID;
    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit ConfigSet(DON_ID, 0);
    s_capabilityRegistry.removeDONs(donIDs);

    (
      uint32 id,
      uint32 configCount,
      bool isPublic,
      bytes32[] memory donNodes,
      CapabilityRegistry.CapabilityConfiguration[] memory donCapabilityConfigs
    ) = s_capabilityRegistry.getDON(DON_ID);
    assertEq(id, 0);
    assertEq(configCount, 0);
    assertEq(isPublic, false);
    assertEq(donCapabilityConfigs.length, 0);
    assertEq(s_capabilityRegistry.getDONCapabilityConfig(DON_ID, s_basicHashedCapabilityId), bytes(""));
    assertEq(donNodes.length, 0);
  }
}
