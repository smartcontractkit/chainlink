// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import {BaseTest} from "./BaseTest.t.sol";
import {ICapabilityConfiguration} from "../interfaces/ICapabilityConfiguration.sol";
import {CapabilitiesRegistry} from "../CapabilitiesRegistry.sol";

contract CapabilitiesRegistry_UpdateDONTest is BaseTest {
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
      signer: NODE_OPERATOR_ONE_SIGNER_ADDRESS,
      hashedCapabilityIds: capabilityIds
    });

    bytes32[] memory nodeTwoCapabilityIds = new bytes32[](1);
    nodeTwoCapabilityIds[0] = s_basicHashedCapabilityId;

    nodes[1] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_TWO_ID,
      p2pId: P2P_ID_TWO,
      signer: NODE_OPERATOR_TWO_SIGNER_ADDRESS,
      hashedCapabilityIds: nodeTwoCapabilityIds
    });

    nodes[2] = CapabilitiesRegistry.NodeParams({
      nodeOperatorId: TEST_NODE_OPERATOR_THREE_ID,
      p2pId: P2P_ID_THREE,
      signer: NODE_OPERATOR_THREE_SIGNER_ADDRESS,
      hashedCapabilityIds: capabilityIds
    });

    s_CapabilitiesRegistry.addNodes(nodes);

    bytes32[] memory donNodes = new bytes32[](2);
    donNodes[0] = P2P_ID;
    donNodes[1] = P2P_ID_TWO;

    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    s_CapabilitiesRegistry.addDON(donNodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    bytes32[] memory nodes = new bytes32[](1);
    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](1);

    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    s_CapabilitiesRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_RevertWhen_NodeDoesNotSupportCapability() public {
    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;
    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_capabilityWithConfigurationContractId,
      config: CONFIG_CAPABILITY_CONFIG
    });
    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilitiesRegistry.NodeDoesNotSupportCapability.selector,
        P2P_ID_TWO,
        s_capabilityWithConfigurationContractId
      )
    );
    s_CapabilitiesRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_RevertWhen_DONDoesNotExist() public {
    uint32 nonExistentDONId = 10;
    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;
    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    vm.expectRevert(abi.encodeWithSelector(CapabilitiesRegistry.DONDoesNotExist.selector, nonExistentDONId));
    s_CapabilitiesRegistry.updateDON(nonExistentDONId, nodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_RevertWhen_CapabilityDoesNotExist() public {
    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;
    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_nonExistentHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    vm.expectRevert(
      abi.encodeWithSelector(CapabilitiesRegistry.CapabilityDoesNotExist.selector, s_nonExistentHashedCapabilityId)
    );
    s_CapabilitiesRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_RevertWhen_DuplicateCapabilityAdded() public {
    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;

    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](2);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    capabilityConfigs[1] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    vm.expectRevert(
      abi.encodeWithSelector(CapabilitiesRegistry.DuplicateDONCapability.selector, 1, s_basicHashedCapabilityId)
    );
    s_CapabilitiesRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_RevertWhen_DeprecatedCapabilityAdded() public {
    bytes32 capabilityId = s_basicHashedCapabilityId;
    bytes32[] memory deprecatedCapabilities = new bytes32[](1);
    deprecatedCapabilities[0] = capabilityId;
    s_CapabilitiesRegistry.deprecateCapabilities(deprecatedCapabilities);

    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID_TWO;

    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: capabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilitiesRegistry.CapabilityIsDeprecated.selector, capabilityId));
    s_CapabilitiesRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_RevertWhen_DuplicateNodeAdded() public {
    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID;

    CapabilitiesRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilitiesRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilitiesRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    vm.expectRevert(abi.encodeWithSelector(CapabilitiesRegistry.DuplicateDONNode.selector, 1, P2P_ID));
    s_CapabilitiesRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true, true, F_VALUE);
  }

  function test_UpdatesDON() public {
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

    CapabilitiesRegistry.DONInfo memory oldDONInfo = s_CapabilitiesRegistry.getDON(DON_ID);

    bool expectedDONIsPublic = false;
    uint32 expectedConfigCount = oldDONInfo.configCount + 1;

    vm.expectEmit(true, true, true, true, address(s_CapabilitiesRegistry));
    emit CapabilitiesRegistry.ConfigSet(DON_ID, expectedConfigCount);
    vm.expectCall(
      address(s_capabilityConfigurationContract),
      abi.encodeWithSelector(
        ICapabilityConfiguration.beforeCapabilityConfigSet.selector,
        nodes,
        CONFIG_CAPABILITY_CONFIG,
        expectedConfigCount,
        DON_ID
      ),
      1
    );
    s_CapabilitiesRegistry.updateDON(DON_ID, nodes, capabilityConfigs, expectedDONIsPublic, true, F_VALUE);

    CapabilitiesRegistry.DONInfo memory donInfo = s_CapabilitiesRegistry.getDON(DON_ID);
    assertEq(donInfo.id, DON_ID);
    assertEq(donInfo.configCount, expectedConfigCount);
    assertEq(donInfo.isPublic, false);
    assertEq(donInfo.capabilityConfigurations.length, capabilityConfigs.length);
    assertEq(donInfo.capabilityConfigurations[0].capabilityId, s_basicHashedCapabilityId);

    (bytes memory CapabilitiesRegistryDONConfig, bytes memory capabilityConfigContractConfig) = s_CapabilitiesRegistry
      .getCapabilityConfigs(DON_ID, s_basicHashedCapabilityId);
    assertEq(CapabilitiesRegistryDONConfig, BASIC_CAPABILITY_CONFIG);
    assertEq(capabilityConfigContractConfig, bytes(""));

    assertEq(donInfo.nodeP2PIds.length, nodes.length);
    assertEq(donInfo.nodeP2PIds[0], P2P_ID);
    assertEq(donInfo.nodeP2PIds[1], P2P_ID_THREE);
  }
}
