// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.24;

import {INodeInfoProvider} from "../../../../keystone/interfaces/INodeInfoProvider.sol";

import {CCIPHome} from "../../../capability/CCIPHome.sol";
import {CCIPHomeHelper} from "../../helpers/CCIPHomeHelper.sol";

import {CCIPHomeTestSetup} from "./CCIPHomeTestSetup.t.sol";

contract CCIPHome_applyChainConfigUpdates is CCIPHomeTestSetup {
  function setUp() public virtual override {
    s_ccipHome = new CCIPHomeHelper(CAPABILITIES_REGISTRY);
  }

  function test_applyChainConfigUpdates_addChainConfigs() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });
    INodeInfoProvider.NodeInfo[] memory nodeInfos = new INodeInfoProvider.NodeInfo[](1);
    nodeInfos[0] = INodeInfoProvider.NodeInfo({
      nodeOperatorId: 1,
      signer: bytes32(uint256(1)),
      p2pId: chainReaders[0],
      encryptionPublicKey: keccak256("encryptionPublicKey"),
      hashedCapabilityIds: new bytes32[](0),
      configCount: uint32(1),
      workflowDONId: uint32(1),
      capabilitiesDONIds: new uint256[](0)
    });
    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(INodeInfoProvider.getNodesByP2PIds.selector, chainReaders),
      abi.encode(nodeInfos)
    );
    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(2, adds[1].chainConfig);
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);

    CCIPHome.ChainConfigArgs[] memory configs = s_ccipHome.getAllChainConfigs(0, 2);
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");
    assertEq(s_ccipHome.getNumChainConfigurations(), 2, "total chain configs must be 2");
  }

  function test_getPaginatedCCIPHomes() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });
    INodeInfoProvider.NodeInfo[] memory nodeInfos = new INodeInfoProvider.NodeInfo[](1);
    nodeInfos[0] = INodeInfoProvider.NodeInfo({
      nodeOperatorId: 1,
      signer: bytes32(uint256(1)),
      p2pId: chainReaders[0],
      encryptionPublicKey: keccak256("encryptionPublicKey"),
      hashedCapabilityIds: new bytes32[](0),
      configCount: uint32(1),
      workflowDONId: uint32(1),
      capabilitiesDONIds: new uint256[](0)
    });
    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(INodeInfoProvider.getNodesByP2PIds.selector, chainReaders),
      abi.encode(nodeInfos)
    );

    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);

    CCIPHome.ChainConfigArgs[] memory configs = s_ccipHome.getAllChainConfigs(0, 2);
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");

    configs = s_ccipHome.getAllChainConfigs(0, 1);
    assertEq(configs.length, 1, "chain configs length must be 1");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");

    configs = s_ccipHome.getAllChainConfigs(0, 10);
    assertEq(configs.length, 2, "chain configs length must be 2");
    assertEq(configs[0].chainSelector, 1, "chain selector must match");
    assertEq(configs[1].chainSelector, 2, "chain selector must match");

    configs = s_ccipHome.getAllChainConfigs(1, 1);
    assertEq(configs.length, 1, "chain configs length must be 1");

    configs = s_ccipHome.getAllChainConfigs(1, 2);
    assertEq(configs.length, 0, "chain configs length must be 0");
  }

  function test_applyChainConfigUpdates_removeChainConfigs() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));

    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config2")})
    });

    INodeInfoProvider.NodeInfo[] memory nodeInfos = new INodeInfoProvider.NodeInfo[](1);
    nodeInfos[0] = INodeInfoProvider.NodeInfo({
      nodeOperatorId: 1,
      signer: bytes32(uint256(1)),
      p2pId: chainReaders[0],
      encryptionPublicKey: keccak256("encryptionPublicKey"),
      hashedCapabilityIds: new bytes32[](0),
      configCount: uint32(1),
      workflowDONId: uint32(1),
      capabilitiesDONIds: new uint256[](0)
    });

    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(INodeInfoProvider.getNodesByP2PIds.selector, chainReaders),
      abi.encode(nodeInfos)
    );

    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(1, adds[0].chainConfig);
    vm.expectEmit();
    emit CCIPHome.ChainConfigSet(2, adds[1].chainConfig);

    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);

    assertEq(s_ccipHome.getNumChainConfigurations(), 2, "total chain configs must be 2");

    assertEq(s_ccipHome.getChainConfig(adds[0].chainSelector).config, adds[0].chainConfig.config);
    assertEq(s_ccipHome.getChainConfig(adds[1].chainSelector).config, adds[1].chainConfig.config);

    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectEmit();
    emit CCIPHome.ChainConfigRemoved(1);
    s_ccipHome.applyChainConfigUpdates(removes, new CCIPHome.ChainConfigArgs[](0));

    assertEq(s_ccipHome.getNumChainConfigurations(), 1, "total chain configs must be 1");
  }

  // Reverts.

  function test_RevertWhen_applyChainConfigUpdates_selectorNotFound() public {
    uint64[] memory removes = new uint64[](1);
    removes[0] = uint64(1);

    vm.expectRevert(abi.encodeWithSelector(CCIPHome.ChainSelectorNotFound.selector, 1));
    s_ccipHome.applyChainConfigUpdates(removes, new CCIPHome.ChainConfigArgs[](0));
  }

  function test_RevertWhen_applyChainConfigUpdates_nodeNotInRegistry() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](1);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: abi.encode(1, 2, 3)})
    });

    vm.mockCallRevert(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(INodeInfoProvider.getNodesByP2PIds.selector, chainReaders),
      abi.encodeWithSelector(INodeInfoProvider.NodeDoesNotExist.selector, chainReaders[0])
    );

    vm.expectRevert(abi.encodeWithSelector(INodeInfoProvider.NodeDoesNotExist.selector, chainReaders[0]));
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);
  }

  function test_RevertWhen__applyChainConfigUpdates_FChainNotPositive() public {
    bytes32[] memory chainReaders = new bytes32[](1);
    chainReaders[0] = keccak256(abi.encode(1));
    CCIPHome.ChainConfigArgs[] memory adds = new CCIPHome.ChainConfigArgs[](2);
    adds[0] = CCIPHome.ChainConfigArgs({
      chainSelector: 1,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 1, config: bytes("config1")})
    });
    adds[1] = CCIPHome.ChainConfigArgs({
      chainSelector: 2,
      chainConfig: CCIPHome.ChainConfig({readers: chainReaders, fChain: 0, config: bytes("config2")}) // bad fChain
    });
    INodeInfoProvider.NodeInfo[] memory nodeInfos = new INodeInfoProvider.NodeInfo[](1);
    nodeInfos[0] = INodeInfoProvider.NodeInfo({
      nodeOperatorId: 1,
      signer: bytes32(uint256(1)),
      p2pId: chainReaders[0],
      encryptionPublicKey: keccak256("encryptionPublicKey"),
      hashedCapabilityIds: new bytes32[](0),
      configCount: uint32(1),
      workflowDONId: uint32(1),
      capabilitiesDONIds: new uint256[](0)
    });
    vm.mockCall(
      CAPABILITIES_REGISTRY,
      abi.encodeWithSelector(INodeInfoProvider.getNodesByP2PIds.selector, chainReaders),
      abi.encode(nodeInfos)
    );

    vm.expectRevert(CCIPHome.FChainMustBePositive.selector);
    s_ccipHome.applyChainConfigUpdates(new uint64[](0), adds);
  }
}
