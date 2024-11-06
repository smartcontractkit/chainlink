// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "./BaseTest.t.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";
import {ICapabilityConfiguration} from "../interfaces/ICapabilityConfiguration.sol";

contract CapabilitiesRegistry_GetNextDONIdTest is BaseTest {
  function setUp() public override {
    BaseTest.setUp();
    CapabilitiesRegistry.Capability[] memory capabilities = new CapabilitiesRegistry.Capability[](2);
    capabilities[0] = s_basicCapability;
    capabilities[1] = s_capabilityWithConfigurationContract;

    s_CapabilitiesRegistry.addNodeOperators(_getNodeOperators());
    s_CapabilitiesRegistry.addCapabilities(capabilities);

    CapabilitiesRegistry.NodeParams[] memory nodes = new CapabilitiesRegistry.NodeParams[](3);
    bytes32[] memory capabilityIds = new bytes32[](2);
    capabilityIds[0] = s_basicHashedCapabilityId;
    capabilityIds[1] = s_capabilityWithConfigurationContractId;

    nodes[0] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_ONE_ID,
      p2pId: P2P_ID,
      encryptionPublicKey: TEST_ENCRYPTION_PUBLIC_KEY,
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: capabilityIds
    });

    bytes32[] memory nodeTwoCapabilityIds = new bytes32[](1);
    nodeTwoCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[1] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_TWO_ID,
      p2pId: P2P_ID_TWO,
      encryptionPublicKey: TEST_ENCRYPTION_PUBLIC_KEY_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: nodeTwoCapabilityIds
    });

    nodes[2] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_THREE_ID,
      p2pId: P2P_ID_THREE,
      encryptionPublicKey: TEST_ENCRYPTION_PUBLIC_KEY_THREE,
      signer: NODE_OPERATOR_THREE_SIGNER_ADDRESS,
      hashedCapabilityIds: capabilityIds
    });

    s_CapabilitiesRegistry.addNodes(nodes);

    changePrank(ADMIN);
  }

  function test_CorrectlyFetchesNextDONId() public {
    uint32 nextDONId = s_CapabilitiesRegistry.getNextDONId();
    assertEq(nextDONId, 1); // Expecting the first DON ID since no DONs have been added yet

    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_THREE;

    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](2);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    capabilityConfigs[1] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_capabilityWithConfigurationContractId,
      config: CONFIG_CAPABILITY_CONFIG
    });

    vm.expectEmit(true, true, true, true, address(s_CapabilitiesRegistry));
    emit CapabilitiesRegistry.ConfigSet(DON_ID, 1);
    vm.expectCall(
      address(s_capabilityConfigurationContract),
      abi.encodeWithSelector(
        ICapabilityConfiguration.beforeCapabilityConfigSet.selector,
        nodes,
        CONFIG_CAPABILITY_CONFIG,
        1,
        DON_ID
      ),
      1
    );
    s_CapabilitiesRegistry.addDON(nodes, capabilityConfigs, true, true, F_VALUE);

    nextDONId = s_CapabilitiesRegistry.getNextDONId();
    assertEq(nextDONId, 2); // After adding one DON, the next DON ID should be 2
  }
}
