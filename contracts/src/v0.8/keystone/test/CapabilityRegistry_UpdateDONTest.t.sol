// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {ICapabilityConfiguration} from "../interfaces/ICapabilityConfiguration.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_UpdateDONTest is BaseTest {
  event ConfigSet(uint32 donId, uint32 configCount);

  uint32 private constant DON_ID = 1;
  uint32 private constant TEST_NODE_OPERATOR_ONE_ID = 1;
  uint256 private constant TEST_NODE_OPERATOR_TWO_ID = 2;
  bytes32 private constant INVALID_P2P_ID = bytes32("fake-p2p");
  bytes private constant BASIC_CAPABILITY_CONFIG = bytes("basic-capability-config");
  bytes private constant CONFIG_CAPABILITY_CONFIG = bytes("config-capability-config");

  function setUp() public override {
    BaseTest.setUp();

    s_capabilityRegistry.addNodeOperators(_getNodeOperators());
    s_capabilityRegistry.addCapability(s_basicCapability);
    s_capabilityRegistry.addCapability(s_capabilityWithConfigurationContract);

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

    changePrank(ADMIN);

    bytes32[] memory donNodes = new bytes32[](1);
    donNodes[0] = P2P_ID;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    s_capabilityRegistry.addDON(donNodes, capabilityConfigs, true);
  }

  function test_RevertWhen_CalledByNonAdmin() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    bytes32[] memory nodes = new bytes32[](1);
    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);

    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    s_capabilityRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true);
  }

  function test_RevertWhen_NodeDoesNotSupportCapability() public {
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID_TWO;
    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_capabilityWithConfigurationContractId,
      config: CONFIG_CAPABILITY_CONFIG
    });
    vm.expectRevert(
      abi.encodeWithSelector(
        CapabilityRegistry.NodeDoesNotSupportCapability.selector,
        P2P_ID_TWO,
        s_capabilityWithConfigurationContractId
      )
    );
    s_capabilityRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true);
  }

  function test_RevertWhen_DONDoesNotExist() public {
    uint32 nonExistentDONId = 10;
    bytes32[] memory nodes = new bytes32[](1);
    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.DONDoesNotExist.selector, nonExistentDONId));
    s_capabilityRegistry.updateDON(nonExistentDONId, nodes, capabilityConfigs, true);
  }

  function test_RevertWhen_CapabilityDoesNotExist() public {
    bytes32[] memory nodes = new bytes32[](1);
    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_nonExistentHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    vm.expectRevert(
      abi.encodeWithSelector(CapabilityRegistry.CapabilityDoesNotExist.selector, s_nonExistentHashedCapabilityId)
    );
    s_capabilityRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true);
  }

  function test_RevertWhen_DuplicateCapabilityAdded() public {
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](2);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    capabilityConfigs[1] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    vm.expectRevert(
      abi.encodeWithSelector(CapabilityRegistry.DuplicateDONCapability.selector, 1, s_basicHashedCapabilityId)
    );
    s_capabilityRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true);
  }

  function test_RevertWhen_DeprecatedCapabilityAdded() public {
    bytes32 capabilityId = s_basicHashedCapabilityId;
    s_capabilityRegistry.deprecateCapability(capabilityId);

    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: capabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });

    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.CapabilityIsDeprecated.selector, capabilityId));
    s_capabilityRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true);
  }

  function test_RevertWhen_DuplicateNodeAdded() public {
    bytes32[] memory nodes = new bytes32[](2);
    nodes[0] = P2P_ID;
    nodes[1] = P2P_ID;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](1);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    vm.expectRevert(abi.encodeWithSelector(CapabilityRegistry.DuplicateDONNode.selector, 1, P2P_ID));
    s_capabilityRegistry.updateDON(DON_ID, nodes, capabilityConfigs, true);
  }

  function test_UpdatesDON() public {
    bytes32[] memory nodes = new bytes32[](1);
    nodes[0] = P2P_ID;

    CapabilityRegistry.CapabilityConfiguration[]
      memory capabilityConfigs = new CapabilityRegistry.CapabilityConfiguration[](2);
    capabilityConfigs[0] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_basicHashedCapabilityId,
      config: BASIC_CAPABILITY_CONFIG
    });
    capabilityConfigs[1] = CapabilityRegistry.CapabilityConfiguration({
      capabilityId: s_capabilityWithConfigurationContractId,
      config: CONFIG_CAPABILITY_CONFIG
    });

    CapabilityRegistry.DONInfo memory oldDONInfo = s_capabilityRegistry.getDON(DON_ID);

    bool expectedDONIsPublic = false;
    uint32 expectedConfigCount = oldDONInfo.configCount + 1;

    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit ConfigSet(DON_ID, expectedConfigCount);
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
    s_capabilityRegistry.updateDON(DON_ID, nodes, capabilityConfigs, expectedDONIsPublic);

    CapabilityRegistry.DONInfo memory DONInfo = s_capabilityRegistry.getDON(DON_ID);
    assertEq(DONInfo.id, DON_ID);
    assertEq(DONInfo.configCount, expectedConfigCount);
    assertEq(DONInfo.isPublic, false);
    assertEq(DONInfo.capabilityConfigurations.length, capabilityConfigs.length);
    assertEq(DONInfo.capabilityConfigurations[0].capabilityId, s_basicHashedCapabilityId);
    assertEq(s_capabilityRegistry.getDONCapabilityConfig(DON_ID, s_basicHashedCapabilityId), BASIC_CAPABILITY_CONFIG);

    assertEq(DONInfo.nodeP2PIds.length, nodes.length);
    assertEq(DONInfo.nodeP2PIds[0], P2P_ID);
  }
}
