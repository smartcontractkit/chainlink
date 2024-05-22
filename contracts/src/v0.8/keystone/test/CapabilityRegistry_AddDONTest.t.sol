// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {BaseTest} from "./BaseTest.t.sol";
import {ICapabilityConfiguration} from "../interfaces/ICapabilityConfiguration.sol";
import {CapabilityRegistry} from "../CapabilityRegistry.sol";

contract CapabilityRegistry_AddDONTest is BaseTest {
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

    changePrank(ADMIN);
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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);
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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);
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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);
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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);
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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);
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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);
  }

  function test_AddDON() public {
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

    vm.expectEmit(true, true, true, true, address(s_capabilityRegistry));
    emit ConfigSet(DON_ID, 1);
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
    s_capabilityRegistry.addDON(nodes, capabilityConfigs, true);

    CapabilityRegistry.DONParams memory donParams = s_capabilityRegistry.getDON(DON_ID);
    assertEq(donParams.id, DON_ID);
    assertEq(donParams.configCount, 1);
    assertEq(donParams.isPublic, true);
    assertEq(donParams.capabilityConfigurations.length, capabilityConfigs.length);
    assertEq(donParams.capabilityConfigurations[0].capabilityId, s_basicHashedCapabilityId);
    assertEq(s_capabilityRegistry.getDONCapabilityConfig(DON_ID, s_basicHashedCapabilityId), BASIC_CAPABILITY_CONFIG);

    assertEq(donParams.nodeP2PIds.length, nodes.length);
    assertEq(donParams.nodeP2PIds[0], P2P_ID);
  }
}
