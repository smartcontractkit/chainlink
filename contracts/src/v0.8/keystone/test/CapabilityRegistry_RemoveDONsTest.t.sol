// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";

import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_RemoveDONsTest is BaseTest {
  event ConfigSet(uint32 donId, uint32 configCount);

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

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    bytes32[] memory nodeIds = new bytes32[](2);
    nodeIds[0] = P2P_ID;
    nodeIds[1] = P2P_ID_TWO;

    changePrank(ADMIN);
    s_capabilityRegistry.addDON(nodeIds, capabilityConfigs, true, true, 1);
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

    CapabilityRegistry.DONInfo memory donInfo = s_capabilityRegistry.getDON(DON_ID);
    assertEq(donInfo.id, 0);
    assertEq(donInfo.configCount, 0);
    assertEq(donInfo.isPublic, false);
    assertEq(donInfo.capabilityConfigurations.length, 0);

    (bytes memory capabilityRegistryDONConfig, bytes memory capabilityConfigContractConfig) = s_capabilityRegistry
      .getCapabilityConfigs(DON_ID, s_basicHashedCapabilityId);

    assertEq(capabilityRegistryDONConfig, bytes(""));
    assertEq(capabilityConfigContractConfig, bytes(""));
    assertEq(donInfo.nodeP2PIds.length, 0);
  }
}
